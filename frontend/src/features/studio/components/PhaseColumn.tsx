import { useState } from "react";
import {
  ActionIcon,
  Badge,
  Box,
  Button,
  Group,
  Stack,
  Text,
  TextInput,
  Tooltip,
  useMantineColorScheme,
  useMantineTheme,
} from "@mantine/core";
import { useDroppable } from "@dnd-kit/core";
import { Icons } from "@/app/theme";
import type { Phase, Piece } from "../types";
import { useCreatePiece } from "../hooks/useStudioQueries";
import { PieceCard } from "./PieceCard";

interface PhaseColumnProps {
  studioId: number;
  phase: Phase;
  pieces: Piece[];
  pieceLabel: string;
  onPieceClick?: (piece: Piece) => void;
}

export function PhaseColumn({
  studioId,
  phase,
  pieces,
  pieceLabel,
  onPieceClick,
}: PhaseColumnProps) {
  const { setNodeRef, isOver } = useDroppable({
    id: `phase-column-${phase.id}`,
  });

  const { colorScheme } = useMantineColorScheme();
  const theme = useMantineTheme();
  const isDark = colorScheme === "dark";

  const overLimit =
    phase.wip_limit != null && pieces.length > phase.wip_limit;

  const [adding, setAdding] = useState(false);
  const [draftTitle, setDraftTitle] = useState("");
  const createPiece = useCreatePiece(studioId);

  const borderColor = isOver
    ? theme.colors.blue[isDark ? 4 : 5]
    : isDark
      ? theme.colors.dark[4]
      : theme.colors.gray[3];
  const backgroundColor = isOver
    ? isDark
      ? theme.colors.dark[5]
      : theme.colors.blue[0]
    : isDark
      ? theme.colors.dark[6]
      : theme.colors.gray[0];
  const headerBorderColor = isDark
    ? theme.colors.dark[4]
    : theme.colors.gray[3];

  const handleCreate = async () => {
    const title = draftTitle.trim();
    if (!title) return;
    await createPiece.mutateAsync({ phase_id: phase.id, title });
    setDraftTitle("");
    setAdding(false);
  };

  return (
    <Box
      ref={setNodeRef}
      style={{
        width: 300,
        minWidth: 300,
        height: "100%",
        alignSelf: "stretch",
        display: "flex",
        flexDirection: "column",
        borderRadius: 8,
        border: `1px solid ${borderColor}`,
        background: backgroundColor,
        transition: "background 120ms ease, border-color 120ms ease",
      }}
    >
      <Group
        justify="space-between"
        p="sm"
        style={{
          borderBottom: `1px solid ${headerBorderColor}`,
          borderLeft: `4px solid ${phase.color || "#868e96"}`,
          flexShrink: 0,
        }}
      >
        <Stack gap={2}>
          <Text fw={600} size="sm">
            {phase.name}
          </Text>
          <Text size="xs" c="dimmed" tt="uppercase">
            {phase.kind}
          </Text>
        </Stack>
        <Group gap="xs">
          <Badge
            variant={overLimit ? "filled" : "light"}
            color={overLimit ? "red" : "gray"}
            size="sm"
          >
            {pieces.length}
            {phase.wip_limit != null ? ` / ${phase.wip_limit}` : ""}
          </Badge>
          <Tooltip label={`New ${pieceLabel.toLowerCase()}`}>
            <ActionIcon
              variant="subtle"
              size="sm"
              onClick={() => setAdding((v) => !v)}
              aria-label={`Add ${pieceLabel.toLowerCase()}`}
            >
              <Icons.Plus size={14} />
            </ActionIcon>
          </Tooltip>
        </Group>
      </Group>

      <Stack
        gap="xs"
        p="sm"
        style={{ flex: 1, minHeight: 0, overflowY: "auto" }}
      >
        {pieces.map((piece) => (
          <PieceCard
            key={piece.id}
            piece={piece}
            accentColor={phase.color || undefined}
            onClick={onPieceClick}
          />
        ))}

        {adding && (
          <Stack gap="xs">
            <TextInput
              autoFocus
              size="xs"
              placeholder={`${pieceLabel} title`}
              value={draftTitle}
              onChange={(e) => setDraftTitle(e.currentTarget.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter") void handleCreate();
                if (e.key === "Escape") {
                  setAdding(false);
                  setDraftTitle("");
                }
              }}
            />
            <Group gap="xs" justify="flex-end">
              <Button
                size="xs"
                variant="subtle"
                onClick={() => {
                  setAdding(false);
                  setDraftTitle("");
                }}
              >
                Cancel
              </Button>
              <Button
                size="xs"
                loading={createPiece.isPending}
                onClick={handleCreate}
                disabled={!draftTitle.trim()}
              >
                Add
              </Button>
            </Group>
          </Stack>
        )}

        {pieces.length === 0 && !adding && (
          <Text size="xs" c="dimmed" ta="center" py="md">
            No {pieceLabel.toLowerCase()}s yet
          </Text>
        )}
      </Stack>
    </Box>
  );
}
