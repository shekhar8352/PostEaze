import { useState } from "react";
import {
  ActionIcon,
  Button,
  Group,
  Stack,
  Text,
  Textarea,
} from "@mantine/core";
import { Icons } from "@/app/theme";
import type { PieceComment } from "../types";
import {
  useCreatePieceComment,
  useDeletePieceComment,
  usePieceComments,
} from "../hooks/useStudioQueries";

interface PieceCommentsPanelProps {
  pieceId: number;
}

export function PieceCommentsPanel({ pieceId }: PieceCommentsPanelProps) {
  const { data: comments = [], isLoading } = usePieceComments(pieceId);
  const createComment = useCreatePieceComment();
  const deleteComment = useDeletePieceComment();
  const [body, setBody] = useState("");

  const submit = async () => {
    const trimmed = body.trim();
    if (!trimmed) return;
    await createComment.mutateAsync({ pieceId, body: { body: trimmed } });
    setBody("");
  };

  return (
    <Stack gap="sm">
      <Textarea
        placeholder="Add a comment…"
        value={body}
        onChange={(e) => setBody(e.currentTarget.value)}
        minRows={2}
        autosize
      />
      <Group justify="flex-end">
        <Button
          size="xs"
          onClick={submit}
          disabled={!body.trim()}
          loading={createComment.isPending}
        >
          Post comment
        </Button>
      </Group>

      {isLoading ? (
        <Text size="xs" c="dimmed">
          Loading comments…
        </Text>
      ) : comments.length === 0 ? (
        <Text size="xs" c="dimmed">
          No comments yet.
        </Text>
      ) : (
        <Stack gap="xs">
          {comments.map((c) => (
            <CommentRow
              key={c.id}
              comment={c}
              onDelete={() =>
                deleteComment.mutate({ pieceId, commentId: c.id })
              }
            />
          ))}
        </Stack>
      )}
    </Stack>
  );
}

function CommentRow({
  comment,
  onDelete,
}: {
  comment: PieceComment;
  onDelete: () => void;
}) {
  return (
    <Stack
      gap={2}
      p="xs"
      style={{
        border: "1px solid var(--mantine-color-gray-2)",
        borderRadius: 6,
      }}
    >
      <Group justify="space-between" wrap="nowrap">
        <Text size="xs" c="dimmed">
          {new Date(comment.created_at).toLocaleString()}
        </Text>
        <ActionIcon
          size="xs"
          variant="subtle"
          color="red"
          onClick={onDelete}
          aria-label="Delete comment"
        >
          <Icons.Trash size={12} />
        </ActionIcon>
      </Group>
      <Text size="sm" style={{ whiteSpace: "pre-wrap" }}>
        {comment.body}
      </Text>
    </Stack>
  );
}
