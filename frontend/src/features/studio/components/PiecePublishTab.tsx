import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import {
  ActionIcon,
  Badge,
  Button,
  Card,
  Group,
  Loader,
  Modal,
  ScrollArea,
  Stack,
  Text,
  ThemeIcon,
  Tooltip,
} from "@mantine/core";
import { useQueries, useQueryClient } from "@tanstack/react-query";
import dayjs from "dayjs";
import { notifications } from "@mantine/notifications";
import { Icons } from "@/app/theme";
import { scheduledPostService } from "@/features/calendar/services/scheduledPostService";
import { scheduledPostKeys } from "@/features/calendar/hooks/useScheduledPostsQueries";
import type { ScheduledPostListItem } from "@/features/calendar/types";
import type { PieceScheduledPostLink } from "../types";
import {
  useLinkScheduledPost,
  useUnlinkScheduledPost,
} from "../hooks/useStudioQueries";

interface PiecePublishTabProps {
  pieceId: number;
  links: PieceScheduledPostLink[];
}

/** Cache key for an individual scheduled post lookup. Kept in sync with
 *  `scheduledPostKeys` so future mutations in the calendar feature can
 *  invalidate the same entries. */
const detailKey = (id: number) =>
  [...scheduledPostKeys.all, "detail", id] as const;

export function PiecePublishTab({ pieceId, links }: PiecePublishTabProps) {
  const [pickerOpen, setPickerOpen] = useState(false);
  const unlink = useUnlinkScheduledPost();

  const postQueries = useQueries({
    queries: links.map((link) => ({
      queryKey: detailKey(link.scheduled_post_id),
      queryFn: () => scheduledPostService.get(link.scheduled_post_id),
    })),
  });

  const loading = postQueries.some((q) => q.isLoading);

  const resolved = useMemo(
    () =>
      links.map((link, i) => ({
        link,
        post: postQueries[i]?.data,
        error: postQueries[i]?.error as Error | undefined,
      })),
    [links, postQueries]
  );

  const handleUnlink = (scheduledPostId: number) => {
    unlink.mutate({ pieceId, scheduledPostId });
  };

  return (
    <Stack gap="md">
      <Group justify="space-between" align="center">
        <Text size="sm" c="dimmed">
          {links.length === 0
            ? "No scheduled posts linked yet."
            : `${links.length} post${links.length === 1 ? "" : "s"} linked`}
        </Text>
        <Group gap="xs">
          <Button
            size="xs"
            variant="light"
            leftSection={<Icons.Plus size={14} />}
            onClick={() => setPickerOpen(true)}
          >
            Link existing
          </Button>
          <Button
            component={Link}
            to="/calendar"
            size="xs"
            leftSection={<Icons.Send size={14} />}
            rightSection={<Icons.ExternalLink size={12} />}
          >
            Schedule new
          </Button>
        </Group>
      </Group>

      {loading && links.length > 0 && (
        <Group gap="xs">
          <Loader size="xs" />
          <Text size="xs" c="dimmed">
            Loading scheduled posts…
          </Text>
        </Group>
      )}

      {resolved.length === 0 ? (
        <Card withBorder radius="md" p="lg">
          <Stack align="center" gap="xs">
            <ThemeIcon variant="light" color="gray" size={40} radius="xl">
              <Icons.Calendar size={20} />
            </ThemeIcon>
            <Text size="sm" c="dimmed" ta="center">
              Schedule a post from the Calendar, then come back here to attach
              it to this piece — or link an existing scheduled post below.
            </Text>
          </Stack>
        </Card>
      ) : (
        <Stack gap="xs">
          {resolved.map(({ link, post, error }) => (
            <ScheduledPostRow
              key={`${link.scheduled_post_id}:${link.role}`}
              link={link}
              post={post}
              error={error}
              onUnlink={() => handleUnlink(link.scheduled_post_id)}
              unlinking={unlink.isPending}
            />
          ))}
        </Stack>
      )}

      <LinkScheduledPostModal
        pieceId={pieceId}
        opened={pickerOpen}
        onClose={() => setPickerOpen(false)}
        existingIds={links.map((l) => l.scheduled_post_id)}
        onLinked={() => {
          notifications.show({
            title: "Linked",
            message: "Scheduled post attached to this piece.",
            color: "green",
          });
        }}
      />
    </Stack>
  );
}

// ---------------------------------------------------------------------------
// Row
// ---------------------------------------------------------------------------

function ScheduledPostRow({
  link,
  post,
  error,
  onUnlink,
  unlinking,
}: {
  link: PieceScheduledPostLink;
  post?: ScheduledPostListItem;
  error?: Error;
  onUnlink: () => void;
  unlinking: boolean;
}) {
  const when = post?.scheduled_at
    ? dayjs(post.scheduled_at).format("MMM D, YYYY · HH:mm")
    : "—";
  const { color, label } = statusBadge(post?.status);

  return (
    <Card withBorder radius="md" p="sm">
      <Group wrap="nowrap" align="flex-start" gap="sm">
        <ThemeIcon variant="light" color={color} size={36} radius="md">
          <Icons.Calendar size={18} />
        </ThemeIcon>
        <Stack gap={4} style={{ flex: 1, minWidth: 0 }}>
          <Group gap="xs" wrap="wrap">
            <Badge size="sm" variant="light" color={color}>
              {label}
            </Badge>
            <Badge size="sm" variant="outline" color="gray">
              {link.role}
            </Badge>
            {post?.post_type && (
              <Badge size="sm" variant="outline" color="blue">
                {post.post_type}
              </Badge>
            )}
            {(post?.platforms ?? []).map((p) => (
              <Badge key={p} size="sm" variant="dot" color="grape">
                {p}
              </Badge>
            ))}
          </Group>
          <Group gap={6}>
            <Icons.Clock size={12} />
            <Text size="xs" c="dimmed">
              {when}
            </Text>
          </Group>
          {post?.caption && (
            <Text size="xs" lineClamp={2}>
              {post.caption}
            </Text>
          )}
          {error && (
            <Text size="xs" c="red">
              Failed to load post #{link.scheduled_post_id}
            </Text>
          )}
        </Stack>
        <Tooltip label="Unlink from piece">
          <ActionIcon
            variant="subtle"
            color="red"
            onClick={onUnlink}
            loading={unlinking}
            aria-label="Unlink scheduled post"
          >
            <Icons.Trash size={14} />
          </ActionIcon>
        </Tooltip>
      </Group>
    </Card>
  );
}

function statusBadge(status?: string): { color: string; label: string } {
  switch (status) {
    case "published":
      return { color: "green", label: "Published" };
    case "scheduled":
      return { color: "blue", label: "Scheduled" };
    case "publishing":
      return { color: "yellow", label: "Publishing" };
    case "failed":
      return { color: "red", label: "Failed" };
    case "cancelled":
    case "canceled":
      return { color: "gray", label: "Cancelled" };
    default:
      return { color: "gray", label: status ? status : "Unknown" };
  }
}

// ---------------------------------------------------------------------------
// Link-existing modal
// ---------------------------------------------------------------------------

interface LinkScheduledPostModalProps {
  pieceId: number;
  opened: boolean;
  onClose: () => void;
  existingIds: number[];
  onLinked?: () => void;
}

function LinkScheduledPostModal({
  pieceId,
  opened,
  onClose,
  existingIds,
  onLinked,
}: LinkScheduledPostModalProps) {
  // Show a rolling ±60 day window so users can pick from recently scheduled
  // or already-published posts. The calendar API requires a range.
  const { from, to } = useMemo(() => {
    const now = dayjs();
    return {
      from: now.subtract(60, "day").toISOString(),
      to: now.add(60, "day").toISOString(),
    };
  }, []);

  const qc = useQueryClient();
  const linkMut = useLinkScheduledPost();
  const [posts, setPosts] = useState<ScheduledPostListItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [loadError, setLoadError] = useState<string | null>(null);

  // Load lazily when the modal opens; avoid firing network on mount.
  useEffect(() => {
    if (!opened) return;
    let cancelled = false;
    setLoading(true);
    setLoadError(null);
    scheduledPostService
      .list(from, to)
      .then((res) => {
        if (!cancelled) setPosts(res.posts ?? []);
      })
      .catch((e: Error) => {
        if (!cancelled) setLoadError(e.message);
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [opened, from, to]);

  const existingSet = useMemo(() => new Set(existingIds), [existingIds]);
  const available = posts.filter((p) => !existingSet.has(p.id));

  const handleLink = async (p: ScheduledPostListItem) => {
    await linkMut.mutateAsync({
      pieceId,
      body: { scheduled_post_id: p.id, role: "primary" },
    });
    // Seed the per-post cache so the detail row renders instantly.
    qc.setQueryData(detailKey(p.id), p);
    onLinked?.();
  };

  return (
    <Modal
      opened={opened}
      onClose={onClose}
      title="Link scheduled post"
      size="lg"
      centered
    >
      <Stack gap="md">
        {loading ? (
          <Group>
            <Loader size="sm" />
            <Text size="sm" c="dimmed">
              Loading scheduled posts…
            </Text>
          </Group>
        ) : loadError ? (
          <Text size="sm" c="red">
            {loadError}
          </Text>
        ) : available.length === 0 ? (
          <Text size="sm" c="dimmed">
            No eligible scheduled posts found in the last/next 60 days. Create
            one from the Calendar and return here to link it.
          </Text>
        ) : (
          <ScrollArea h={420} type="auto">
            <Stack gap="xs">
              {available.map((post) => (
                <ScheduledPostPickerRow
                  key={post.id}
                  post={post}
                  onLink={() => handleLink(post)}
                  linking={linkMut.isPending}
                />
              ))}
            </Stack>
          </ScrollArea>
        )}
      </Stack>
    </Modal>
  );
}

function ScheduledPostPickerRow({
  post,
  onLink,
  linking,
}: {
  post: ScheduledPostListItem;
  onLink: () => void;
  linking: boolean;
}) {
  const when = dayjs(post.scheduled_at).format("MMM D, YYYY · HH:mm");
  const { color, label } = statusBadge(post.status);
  return (
    <Card withBorder radius="md" p="xs">
      <Group wrap="nowrap" align="center" gap="sm">
        <ThemeIcon variant="light" color={color} size={32} radius="md">
          <Icons.Calendar size={16} />
        </ThemeIcon>
        <Stack gap={2} style={{ flex: 1, minWidth: 0 }}>
          <Group gap="xs" wrap="wrap">
            <Text size="sm" fw={500} lineClamp={1}>
              {post.caption?.trim() || `Post #${post.id}`}
            </Text>
            <Badge size="xs" variant="light" color={color}>
              {label}
            </Badge>
          </Group>
          <Group gap={6}>
            <Icons.Clock size={12} />
            <Text size="xs" c="dimmed">
              {when}
            </Text>
            {post.platforms?.map((p) => (
              <Badge key={p} size="xs" variant="dot" color="grape">
                {p}
              </Badge>
            ))}
          </Group>
        </Stack>
        <Button size="xs" variant="light" onClick={onLink} loading={linking}>
          Link
        </Button>
      </Group>
    </Card>
  );
}
