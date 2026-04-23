import { Badge, Card, Group, Stack, Text } from "@mantine/core";
import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import type { CSSProperties } from "react";
import type { Piece } from "../types";

interface PieceCardProps {
  piece: Piece;
  /** Rendered inside DragOverlay as a static preview. */
  dragging?: boolean;
  onClick?: (piece: Piece) => void;
}

export function PieceCard({ piece, dragging, onClick }: PieceCardProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: piece.id,
    disabled: dragging,
  });

  const style: CSSProperties = {
    transform: CSS.Translate.toString(transform),
    transition,
    opacity: isDragging ? 0.4 : 1,
    cursor: "grab",
  };

  return (
    <Card
      ref={dragging ? undefined : setNodeRef}
      withBorder
      radius="md"
      p="sm"
      shadow={dragging ? "md" : "xs"}
      style={dragging ? { cursor: "grabbing" } : style}
      {...(dragging ? {} : attributes)}
      {...(dragging ? {} : listeners)}
      onClick={(e) => {
        // Skip click while dragging to avoid opening drawer on drop.
        if (isDragging) return;
        e.stopPropagation();
        onClick?.(piece);
      }}
    >
      <Stack gap={4}>
        <Text fw={500} size="sm" lineClamp={2}>
          {piece.title}
        </Text>
        <Group gap="xs">
          <Badge size="xs" variant="light" color="blue">
            {piece.content_type}
          </Badge>
          {piece.due_at && (
            <Text size="xs" c="dimmed">
              Due {new Date(piece.due_at).toLocaleDateString()}
            </Text>
          )}
        </Group>
      </Stack>
    </Card>
  );
}
