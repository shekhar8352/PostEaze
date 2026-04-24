import {
  Badge,
  Card,
  Group,
  Stack,
  Text,
  useMantineColorScheme,
  useMantineTheme,
} from "@mantine/core";
import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import type { CSSProperties } from "react";
import type { Piece } from "../types";

interface PieceCardProps {
  piece: Piece;
  /** Rendered inside DragOverlay as a static preview. */
  dragging?: boolean;
  /** Phase accent; ties the card visually to its column rail. */
  accentColor?: string;
  onClick?: (piece: Piece) => void;
}

export function PieceCard({
  piece,
  dragging,
  accentColor,
  onClick,
}: PieceCardProps) {
  const theme = useMantineTheme();
  const { colorScheme } = useMantineColorScheme();
  const isDark = colorScheme === "dark";

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

  const rail = accentColor || theme.colors.blue[isDark ? 4 : 6];
  const surfaceBg = isDark ? theme.colors.dark[7] : theme.white;
  const edge = isDark ? theme.colors.dark[4] : theme.colors.gray[3];
  const edgeHover = isDark ? theme.colors.dark[3] : theme.colors.gray[4];
  const liftShadow = isDark
    ? "0 4px 16px rgba(0, 0, 0, 0.38), 0 0 0 1px rgba(255, 255, 255, 0.05) inset"
    : "0 2px 10px rgba(0, 0, 0, 0.07), 0 0 0 1px rgba(255, 255, 255, 0.85) inset";
  const dragShadow = isDark
    ? "0 18px 40px rgba(0, 0, 0, 0.55), 0 0 0 1px rgba(255, 255, 255, 0.06) inset"
    : "0 14px 34px rgba(0, 0, 0, 0.14), 0 0 0 1px rgba(255, 255, 255, 0.9) inset";

  return (
    <Card
      ref={dragging ? undefined : setNodeRef}
      withBorder={false}
      radius="md"
      p="sm"
      bg={surfaceBg}
      style={{
        border: `1px solid ${edge}`,
        borderLeft: `4px solid ${rail}`,
        boxShadow: dragging ? dragShadow : liftShadow,
        ...(dragging ? { cursor: "grabbing" } : style),
      }}
      styles={{
        root: {
          "@media (hover: hover)": {
            "&:hover": {
              borderColor: edgeHover,
              boxShadow: isDark
                ? "0 8px 22px rgba(0, 0, 0, 0.42), 0 0 0 1px rgba(255, 255, 255, 0.06) inset"
                : "0 4px 16px rgba(0, 0, 0, 0.09), 0 0 0 1px rgba(255, 255, 255, 0.9) inset",
            },
          },
        },
      }}
      {...(dragging ? {} : attributes)}
      {...(dragging ? {} : listeners)}
      onClick={(e) => {
        // Skip click while dragging to avoid opening drawer on drop.
        if (isDragging) return;
        e.stopPropagation();
        onClick?.(piece);
      }}
    >
      <Stack gap={6}>
        <Text fw={600} size="sm" lineClamp={2} style={{ letterSpacing: "-0.01em" }}>
          {piece.title}
        </Text>
        <Group gap="xs" wrap="nowrap">
          <Badge
            size="xs"
            variant="filled"
            color="blue"
            radius="xl"
            styles={{
              root: {
                fontWeight: 650,
                letterSpacing: "0.04em",
                textTransform: "uppercase",
              },
            }}
          >
            {piece.content_type}
          </Badge>
          {piece.due_at && (
            <Text size="xs" c="dimmed" style={{ flexShrink: 0 }}>
              Due {new Date(piece.due_at).toLocaleDateString()}
            </Text>
          )}
        </Group>
      </Stack>
    </Card>
  );
}
