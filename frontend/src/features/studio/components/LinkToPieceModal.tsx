import { useEffect, useMemo, useState } from "react";
import {
  Alert,
  Badge,
  Button,
  Center,
  Group,
  Loader,
  Modal,
  ScrollArea,
  Stack,
  Text,
  TextInput,
  UnstyledButton,
} from "@mantine/core";
import { Link } from "react-router-dom";
import { Icons } from "@/app/theme";
import {
  useDefaultStudio,
  useStudioBoard,
} from "../hooks/useStudioQueries";
import type { Phase, Piece } from "../types";

interface LinkToPieceModalProps {
  opened: boolean;
  onClose: () => void;
  /**
   * Invoked with the chosen Piece so the parent can run the actual
   * link mutation (e.g. link media asset or scheduled post).
   */
  onSelect: (piece: Piece) => Promise<void> | void;
  title?: string;
  /** Rendered below the title for extra context. */
  description?: string;
  /** Label for the primary select button. Defaults to "Link". */
  ctaLabel?: string;
  /** When true, disables action buttons (parent is busy). */
  isSubmitting?: boolean;
}

function phaseLookup(phases: Phase[] | undefined) {
  const map = new Map<number, Phase>();
  (phases ?? []).forEach((p) => map.set(p.id, p));
  return map;
}

export function LinkToPieceModal({
  opened,
  onClose,
  onSelect,
  title = "Link to a Piece",
  description,
  ctaLabel = "Link",
  isSubmitting = false,
}: LinkToPieceModalProps) {
  const { data: studio, isLoading: studioLoading } = useDefaultStudio();
  const studioId = studio?.id;
  const {
    data: board,
    isLoading: boardLoading,
    isError: boardError,
  } = useStudioBoard(studioId);

  const [query, setQuery] = useState("");
  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [localSubmitting, setLocalSubmitting] = useState(false);

  useEffect(() => {
    if (!opened) {
      setQuery("");
      setSelectedId(null);
      setLocalSubmitting(false);
    }
  }, [opened]);

  const phases = useMemo(() => phaseLookup(board?.phases), [board?.phases]);
  const pieces = useMemo(() => {
    const list = board?.pieces ?? [];
    const q = query.trim().toLowerCase();
    if (!q) return list;
    return list.filter((p) => p.title.toLowerCase().includes(q));
  }, [board?.pieces, query]);

  const selected = useMemo(
    () => pieces.find((p) => p.id === selectedId) ?? null,
    [pieces, selectedId]
  );

  const loading = studioLoading || boardLoading;
  const busy = isSubmitting || localSubmitting;

  const handleConfirm = async () => {
    if (!selected || busy) return;
    try {
      setLocalSubmitting(true);
      await onSelect(selected);
      onClose();
    } finally {
      setLocalSubmitting(false);
    }
  };

  return (
    <Modal
      opened={opened}
      onClose={onClose}
      title={<Text fw={600}>{title}</Text>}
      size="lg"
      radius="md"
    >
      <Stack gap="md">
        {description ? (
          <Text size="sm" c="dimmed">
            {description}
          </Text>
        ) : null}

        <TextInput
          placeholder={`Search ${studio?.piece_label ?? "piece"}s by title…`}
          leftSection={<Icons.Search size={16} />}
          value={query}
          onChange={(e) => setQuery(e.currentTarget.value)}
          disabled={loading || busy}
        />

        <ScrollArea h={320} type="auto" offsetScrollbars>
          {loading ? (
            <Center py="xl">
              <Loader size="sm" />
            </Center>
          ) : boardError ? (
            <Alert color="red" variant="light">
              Could not load your Studio board. Try again.
            </Alert>
          ) : pieces.length === 0 ? (
            <Center py="xl">
              <Stack align="center" gap={6}>
                <Text size="sm" c="dimmed">
                  {board?.pieces.length
                    ? "No matching items."
                    : `No ${studio?.piece_label.toLowerCase() ?? "piece"}s in your Studio yet.`}
                </Text>
                {!board?.pieces.length && (
                  <Button
                    component={Link}
                    to="/studio"
                    variant="subtle"
                    size="xs"
                    onClick={onClose}
                  >
                    Open Studio
                  </Button>
                )}
              </Stack>
            </Center>
          ) : (
            <Stack gap={4}>
              {pieces.map((piece) => {
                const phase = phases.get(piece.phase_id);
                const active = piece.id === selectedId;
                return (
                  <UnstyledButton
                    key={piece.id}
                    onClick={() => setSelectedId(piece.id)}
                    disabled={busy}
                    style={(theme) => ({
                      padding: "8px 12px",
                      borderRadius: theme.radius.sm,
                      border: `1px solid ${
                        active
                          ? "var(--mantine-color-blue-5)"
                          : "var(--mantine-color-gray-3)"
                      }`,
                      background: active
                        ? "var(--mantine-color-blue-0)"
                        : "transparent",
                      cursor: busy ? "not-allowed" : "pointer",
                      opacity: busy ? 0.7 : 1,
                    })}
                  >
                    <Group justify="space-between" wrap="nowrap">
                      <Stack gap={2} style={{ minWidth: 0 }}>
                        <Text size="sm" fw={600} lineClamp={1}>
                          {piece.title}
                        </Text>
                        <Text size="xs" c="dimmed" lineClamp={1}>
                          {piece.content_type}
                        </Text>
                      </Stack>
                      {phase ? (
                        <Badge
                          color={phase.color ? undefined : "gray"}
                          styles={
                            phase.color
                              ? {
                                  root: {
                                    backgroundColor: phase.color,
                                    color: "#fff",
                                  },
                                }
                              : undefined
                          }
                          variant={phase.color ? "filled" : "light"}
                          size="sm"
                          style={{ flexShrink: 0 }}
                        >
                          {phase.name}
                        </Badge>
                      ) : null}
                    </Group>
                  </UnstyledButton>
                );
              })}
            </Stack>
          )}
        </ScrollArea>

        <Group justify="space-between">
          <Button
            component={Link}
            to="/studio"
            variant="subtle"
            size="sm"
            onClick={onClose}
            leftSection={<Icons.Kanban size={14} />}
          >
            Open Studio
          </Button>
          <Group gap="xs">
            <Button variant="default" onClick={onClose} disabled={busy}>
              Cancel
            </Button>
            <Button
              onClick={handleConfirm}
              disabled={!selected || busy}
              loading={busy}
            >
              {ctaLabel}
            </Button>
          </Group>
        </Group>
      </Stack>
    </Modal>
  );
}

export default LinkToPieceModal;
