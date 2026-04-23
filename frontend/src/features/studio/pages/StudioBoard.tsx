import { useMemo, useState } from "react";
import {
  Alert,
  Box,
  Button,
  Center,
  Group,
  Loader,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { Link } from "react-router-dom";
import { useQueryClient } from "@tanstack/react-query";
import {
  DndContext,
  DragOverlay,
  PointerSensor,
  closestCenter,
  useSensor,
  useSensors,
  type DragEndEvent,
  type DragStartEvent,
} from "@dnd-kit/core";
import {
  SortableContext,
  arrayMove,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { Icons } from "@/app/theme";
import {
  studioKeys,
  useDefaultStudio,
  useMovePiece,
  useStudioBoard,
} from "../hooks/useStudioQueries";
import type { Phase, Piece, StudioBoard as StudioBoardData } from "../types";
import { PhaseColumn } from "../components/PhaseColumn";
import { PieceCard } from "../components/PieceCard";
import { PieceDetailDrawer } from "../components/PieceDetailDrawer";

/**
 * StudioBoard - Kanban-style board for the default studio.
 *
 * Supports drag-and-drop of pieces between and within phases using
 * @dnd-kit. Moves are applied optimistically: the cached board is
 * mutated immediately and then reconciled with the backend response.
 * If the server rejects the move the previous snapshot is restored.
 */
export default function StudioBoard() {
  const qc = useQueryClient();
  const {
    data: studio,
    isLoading: studioLoading,
    error: studioError,
  } = useDefaultStudio();

  const {
    data: board,
    isLoading: boardLoading,
    error: boardError,
  } = useStudioBoard(studio?.id);

  const movePiece = useMovePiece(studio?.id ?? 0);

  const [activePieceId, setActivePieceId] = useState<number | null>(null);
  const [detailPieceId, setDetailPieceId] = useState<number | null>(null);

  const { phases, piecesByPhase, allPieceIds } = useMemo(() => {
    const byPhase = new Map<number, Piece[]>();
    const ids: number[] = [];
    if (!board) {
      return { phases: [] as Phase[], piecesByPhase: byPhase, allPieceIds: ids };
    }
    for (const piece of board.pieces) {
      if (piece.status !== "active") continue;
      const arr = byPhase.get(piece.phase_id) ?? [];
      arr.push(piece);
      byPhase.set(piece.phase_id, arr);
      ids.push(piece.id);
    }
    for (const arr of byPhase.values()) {
      arr.sort((a, b) => (a.position < b.position ? -1 : 1));
    }
    const sortedPhases = board.phases
      .slice()
      .sort((a, b) => a.order_index - b.order_index);
    return {
      phases: sortedPhases,
      piecesByPhase: byPhase,
      allPieceIds: ids,
    };
  }, [board]);

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } })
  );

  if (studioLoading || boardLoading) {
    return (
      <Center style={{ minHeight: 320 }}>
        <Loader />
      </Center>
    );
  }

  const error = studioError ?? boardError;
  if (error) {
    return (
      <Alert
        color="red"
        icon={<Icons.AlertCircle size={18} />}
        title="Failed to load Studio"
      >
        {error instanceof Error ? error.message : "Something went wrong."}
      </Alert>
    );
  }

  if (!studio || !board) return null;

  const pieceLabel = studio.piece_label || "Piece";
  const activePiece = activePieceId
    ? board.pieces.find((p) => p.id === activePieceId) ?? null
    : null;

  const handleDragStart = (event: DragStartEvent) => {
    const id = Number(event.active.id);
    if (!Number.isNaN(id)) setActivePieceId(id);
  };

  const handleDragCancel = () => setActivePieceId(null);

  const handleDragEnd = (event: DragEndEvent) => {
    setActivePieceId(null);
    const { active, over } = event;
    if (!over) return;

    const activeId = Number(active.id);
    if (Number.isNaN(activeId)) return;

    const activePieceRecord = board.pieces.find((p) => p.id === activeId);
    if (!activePieceRecord) return;

    const { targetPhaseId, targetIndex } = resolveDropTarget({
      overId: String(over.id),
      piecesByPhase,
      phases,
    });
    if (targetPhaseId == null) return;

    const sourcePieces = piecesByPhase.get(activePieceRecord.phase_id) ?? [];
    const targetPiecesOriginal = piecesByPhase.get(targetPhaseId) ?? [];

    // Build the post-move target list so we can compute neighbors the
    // backend expects for fractional indexing.
    let targetPieces: Piece[];
    if (activePieceRecord.phase_id === targetPhaseId) {
      const fromIdx = sourcePieces.findIndex((p) => p.id === activeId);
      if (fromIdx === -1) return;
      // Clamp to valid range after removal
      const safeIndex = Math.min(
        Math.max(targetIndex, 0),
        sourcePieces.length - 1
      );
      if (fromIdx === safeIndex) return; // no-op
      targetPieces = arrayMove(sourcePieces, fromIdx, safeIndex);
    } else {
      const withoutActive = targetPiecesOriginal.slice();
      const insertAt = Math.min(
        Math.max(targetIndex, 0),
        withoutActive.length
      );
      withoutActive.splice(insertAt, 0, activePieceRecord);
      targetPieces = withoutActive;
    }

    const newIndex = targetPieces.findIndex((p) => p.id === activeId);
    const before = newIndex > 0 ? targetPieces[newIndex - 1] : null;
    const after =
      newIndex < targetPieces.length - 1 ? targetPieces[newIndex + 1] : null;

    // Optimistic cache update.
    const prevBoard = qc.getQueryData<StudioBoardData>(
      studioKeys.board(studio.id)
    );
    if (prevBoard) {
      qc.setQueryData<StudioBoardData>(studioKeys.board(studio.id), {
        ...prevBoard,
        pieces: prevBoard.pieces.map((p) => {
          if (p.id !== activeId) return p;
          return {
            ...p,
            phase_id: targetPhaseId,
            // Position is recomputed by the server; a placeholder
            // keeps client sort stable until the response arrives.
            position: computeInterimPosition(before, after),
          };
        }),
      });
    }

    movePiece.mutate(
      {
        pieceId: activeId,
        body: {
          phase_id: targetPhaseId,
          before_piece_id: before?.id ?? null,
          after_piece_id: after?.id ?? null,
        },
      },
      {
        onError: () => {
          if (prevBoard) {
            qc.setQueryData(studioKeys.board(studio.id), prevBoard);
          }
        },
      }
    );
  };

  return (
    <Stack
      gap="md"
      p="md"
      style={{
        height: "calc(100dvh - 60px - var(--mantine-spacing-md) * 2)",
        minHeight: 0,
      }}
    >
      <Group justify="space-between" align="center" style={{ flexShrink: 0 }}>
        <Stack gap={2}>
          <Title order={2}>{studio.name}</Title>
          <Text size="sm" c="dimmed">
            Content pipeline · {phases.length} phases ·{" "}
            {allPieceIds.length} active {pieceLabel.toLowerCase()}s
          </Text>
        </Stack>
        <Button
          component={Link}
          to="/studio/settings"
          variant="subtle"
          leftSection={<Icons.Settings size={16} />}
        >
          Settings
        </Button>
      </Group>

      <DndContext
        sensors={sensors}
        collisionDetection={closestCenter}
        onDragStart={handleDragStart}
        onDragEnd={handleDragEnd}
        onDragCancel={handleDragCancel}
      >
        <Box
          style={{
            flex: 1,
            minHeight: 0,
            overflowX: "auto",
            overflowY: "hidden",
          }}
        >
          <Group
            align="stretch"
            gap="md"
            wrap="nowrap"
            h="100%"
            style={{ minWidth: "min-content" }}
          >
            {phases.map((phase) => {
              const pieces = piecesByPhase.get(phase.id) ?? [];
              return (
                <SortableContext
                  key={phase.id}
                  id={`phase-${phase.id}`}
                  items={pieces.map((p) => p.id)}
                  strategy={verticalListSortingStrategy}
                >
                  <PhaseColumn
                    studioId={studio.id}
                    phase={phase}
                    pieces={pieces}
                    pieceLabel={pieceLabel}
                    onPieceClick={(p) => setDetailPieceId(p.id)}
                  />
                </SortableContext>
              );
            })}
            {phases.length === 0 && (
              <Box p="lg">
                <Text size="sm" c="dimmed">
                  This studio has no phases yet. Add one from Studio settings.
                </Text>
              </Box>
            )}
          </Group>
        </Box>

        <DragOverlay>
          {activePiece ? <PieceCard piece={activePiece} dragging /> : null}
        </DragOverlay>
      </DndContext>

      <PieceDetailDrawer
        pieceId={detailPieceId}
        studioId={studio.id}
        phases={phases}
        opened={detailPieceId != null}
        onClose={() => setDetailPieceId(null)}
      />
    </Stack>
  );
}

/**
 * Resolve where the piece should land given the drop target id.
 * @dnd-kit gives us either the id of another sortable item (a piece)
 * or, when the user drops onto empty space, the id of the droppable
 * column we exposed (`phase-column-<id>`).
 */
function resolveDropTarget({
  overId,
  piecesByPhase,
  phases,
}: {
  overId: string;
  piecesByPhase: Map<number, Piece[]>;
  phases: Phase[];
}): { targetPhaseId: number | null; targetIndex: number } {
  if (overId.startsWith("phase-column-")) {
    const phaseId = Number(overId.replace("phase-column-", ""));
    const pieces = piecesByPhase.get(phaseId) ?? [];
    return { targetPhaseId: phaseId, targetIndex: pieces.length };
  }

  const overPieceId = Number(overId);
  if (Number.isNaN(overPieceId)) return { targetPhaseId: null, targetIndex: 0 };

  for (const phase of phases) {
    const pieces = piecesByPhase.get(phase.id) ?? [];
    const idx = pieces.findIndex((p) => p.id === overPieceId);
    if (idx !== -1) {
      return { targetPhaseId: phase.id, targetIndex: idx };
    }
  }
  return { targetPhaseId: null, targetIndex: 0 };
}

/**
 * Generate a placeholder position string between two neighbors for
 * optimistic rendering. The server recomputes the canonical value;
 * this only needs to preserve ordering between `before` and `after`.
 */
function computeInterimPosition(
  before: Piece | null,
  after: Piece | null
): string {
  if (before && after) return `${before.position}~${after.position}`;
  if (before) return `${before.position}~`;
  if (after) return `~${after.position}`;
  return "a";
}
