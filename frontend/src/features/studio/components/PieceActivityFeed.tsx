import { Stack, Text, ThemeIcon, Group, Loader, Timeline } from "@mantine/core";
import { Icons } from "@/app/theme";
import { usePieceActivities } from "../hooks/useStudioQueries";
import type { PieceActivity } from "../types";

interface PieceActivityFeedProps {
  pieceId: number;
}

export function PieceActivityFeed({ pieceId }: PieceActivityFeedProps) {
  const { data: activities = [], isLoading } = usePieceActivities(pieceId);

  if (isLoading) {
    return (
      <Group gap="xs">
        <Loader size="xs" />
        <Text size="xs" c="dimmed">
          Loading activity…
        </Text>
      </Group>
    );
  }

  if (activities.length === 0) {
    return (
      <Text size="xs" c="dimmed">
        No activity yet.
      </Text>
    );
  }

  return (
    <Timeline bulletSize={22} lineWidth={2} active={-1}>
      {activities.map((a) => (
        <Timeline.Item
          key={a.id}
          bullet={
            <ThemeIcon
              variant="light"
              color={iconColor(a.kind)}
              size={22}
              radius="xl"
            >
              {renderIcon(a.kind)}
            </ThemeIcon>
          }
          title={
            <Text size="sm" fw={500}>
              {humanizeKind(a.kind)}
            </Text>
          }
        >
          <Stack gap={2}>
            <Text size="xs" c="dimmed">
              {new Date(a.created_at).toLocaleString()}
            </Text>
            <ActivityPayload activity={a} />
          </Stack>
        </Timeline.Item>
      ))}
    </Timeline>
  );
}

function ActivityPayload({ activity }: { activity: PieceActivity }) {
  const payload = activity.payload;
  if (!payload || typeof payload !== "object") return null;
  const entries = Object.entries(payload as Record<string, unknown>).filter(
    ([, v]) => v !== null && v !== undefined && v !== ""
  );
  if (entries.length === 0) return null;
  return (
    <Text size="xs" c="dimmed" style={{ whiteSpace: "pre-wrap" }}>
      {entries
        .map(([k, v]) => `${k}: ${formatValue(v)}`)
        .join("  •  ")}
    </Text>
  );
}

function formatValue(v: unknown): string {
  if (typeof v === "string" || typeof v === "number" || typeof v === "boolean") {
    return String(v);
  }
  try {
    return JSON.stringify(v);
  } catch {
    return "";
  }
}

function humanizeKind(kind: string): string {
  return kind
    .replace(/[._]/g, " ")
    .replace(/\b\w/g, (c) => c.toUpperCase());
}

function iconColor(kind: string): string {
  if (kind.includes("publish")) return "green";
  if (kind.includes("schedule")) return "blue";
  if (kind.includes("asset")) return "grape";
  if (kind.includes("phase") || kind.includes("move")) return "indigo";
  if (kind.includes("comment")) return "yellow";
  if (kind.includes("delete") || kind.includes("archive")) return "red";
  return "gray";
}

function renderIcon(kind: string) {
  if (kind.includes("publish")) return <Icons.Send size={12} />;
  if (kind.includes("schedule")) return <Icons.Clock size={12} />;
  if (kind.includes("asset")) return <Icons.Photo size={12} />;
  if (kind.includes("phase") || kind.includes("move"))
    return <Icons.Kanban size={12} />;
  if (kind.includes("comment")) return <Icons.FileText size={12} />;
  if (kind.includes("delete") || kind.includes("archive"))
    return <Icons.Trash size={12} />;
  return <Icons.Dots size={12} />;
}
