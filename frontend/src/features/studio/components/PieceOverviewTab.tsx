import { useEffect, useState } from "react";
import {
  ActionIcon,
  Badge,
  Button,
  Group,
  Select,
  Stack,
  Text,
  Textarea,
  TextInput,
  Tooltip,
} from "@mantine/core";
import { DateTimePicker } from "@mantine/dates";
import dayjs from "dayjs";
import { notifications } from "@mantine/notifications";
import { Icons } from "@/app/theme";
import type { Piece, PieceContentType, Phase } from "../types";
import {
  useDeletePiece,
  useSetPieceStatus,
  useUpdatePiece,
} from "../hooks/useStudioQueries";

interface PieceOverviewTabProps {
  piece: Piece;
  phase?: Phase;
  studioId: number;
  onClose: () => void;
}

const CONTENT_TYPE_OPTIONS: { value: PieceContentType; label: string }[] = [
  { value: "post", label: "Post" },
  { value: "reel", label: "Reel" },
  { value: "story", label: "Story" },
  { value: "video", label: "Video" },
  { value: "carousel", label: "Carousel" },
  { value: "other", label: "Other" },
];

export function PieceOverviewTab({
  piece,
  phase,
  studioId,
  onClose,
}: PieceOverviewTabProps) {
  const updatePiece = useUpdatePiece(studioId);
  const setStatus = useSetPieceStatus(studioId);
  const deletePiece = useDeletePiece(studioId);

  const [title, setTitle] = useState(piece.title);
  const [description, setDescription] = useState(piece.description ?? "");
  const [contentType, setContentType] = useState<PieceContentType>(
    piece.content_type
  );
  const [dueAt, setDueAt] = useState<string | null>(
    piece.due_at ? dayjs(piece.due_at).format("YYYY-MM-DD HH:mm") : null
  );

  // Reset local state when switching between pieces.
  useEffect(() => {
    setTitle(piece.title);
    setDescription(piece.description ?? "");
    setContentType(piece.content_type);
    setDueAt(
      piece.due_at ? dayjs(piece.due_at).format("YYYY-MM-DD HH:mm") : null
    );
  }, [piece.id, piece.title, piece.description, piece.content_type, piece.due_at]);

  const dueAtIso = dueAt && dayjs(dueAt).isValid() ? dayjs(dueAt).toISOString() : null;

  const dirty =
    title.trim() !== piece.title ||
    description !== (piece.description ?? "") ||
    contentType !== piece.content_type ||
    dueAtIso !== (piece.due_at ?? null);

  const handleSave = async () => {
    const trimmed = title.trim();
    if (!trimmed) {
      notifications.show({
        title: "Title required",
        message: "A piece needs a title.",
        color: "red",
      });
      return;
    }
    await updatePiece.mutateAsync({
      pieceId: piece.id,
      body: {
        title: trimmed,
        description,
        content_type: contentType,
        due_at: dueAtIso,
      },
    });
    notifications.show({
      title: "Saved",
      message: "Piece updated.",
      color: "green",
    });
  };

  const handleArchiveToggle = async () => {
    const next = piece.status === "archived" ? "active" : "archived";
    await setStatus.mutateAsync({ pieceId: piece.id, status: next });
    notifications.show({
      title: next === "archived" ? "Archived" : "Restored",
      message:
        next === "archived"
          ? "Moved to archive."
          : "Returned to the board.",
      color: "green",
    });
  };

  const handleDelete = async () => {
    if (!window.confirm("Delete this piece? This cannot be undone.")) return;
    await deletePiece.mutateAsync(piece.id);
    onClose();
  };

  return (
    <Stack gap="md">
      <TextInput
        label="Title"
        value={title}
        onChange={(e) => setTitle(e.currentTarget.value)}
        required
      />

      <Textarea
        label="Description"
        placeholder="Script notes, hook, CTA, references…"
        value={description}
        onChange={(e) => setDescription(e.currentTarget.value)}
        minRows={3}
        autosize
      />

      <Group grow align="flex-start">
        <Select
          label="Content type"
          value={contentType}
          onChange={(v) =>
            v && setContentType(v as PieceContentType)
          }
          data={CONTENT_TYPE_OPTIONS}
          allowDeselect={false}
        />
        <DateTimePicker
          label="Due"
          value={dueAt}
          onChange={setDueAt}
          valueFormat="YYYY-MM-DD HH:mm"
          clearable
          placeholder="No due date"
          popoverProps={{ withinPortal: true }}
        />
      </Group>

      <Group gap="xs" align="center">
        <Text size="sm" c="dimmed">
          Phase:
        </Text>
        {phase ? (
          <Badge
            variant="light"
            color="gray"
            styles={{ root: { borderLeft: `3px solid ${phase.color || "#868e96"}` } }}
          >
            {phase.name}
          </Badge>
        ) : (
          <Badge variant="light" color="gray">
            Unknown
          </Badge>
        )}
        <Text size="xs" c="dimmed">
          (drag on the board to change)
        </Text>
      </Group>

      <Group justify="space-between" mt="xs">
        <Group gap="xs">
          <Button
            onClick={handleSave}
            disabled={!dirty}
            loading={updatePiece.isPending}
          >
            Save changes
          </Button>
          <Button
            variant="subtle"
            color="gray"
            onClick={handleArchiveToggle}
            loading={setStatus.isPending}
          >
            {piece.status === "archived" ? "Restore" : "Archive"}
          </Button>
        </Group>
        <Tooltip label="Delete piece">
          <ActionIcon
            variant="subtle"
            color="red"
            onClick={handleDelete}
            loading={deletePiece.isPending}
            aria-label="Delete piece"
          >
            <Icons.Trash size={16} />
          </ActionIcon>
        </Tooltip>
      </Group>
    </Stack>
  );
}
