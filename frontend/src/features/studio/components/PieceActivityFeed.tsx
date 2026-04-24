import { Stack, Text, ThemeIcon, Group, Loader, Timeline } from "@mantine/core";
import { Icons } from "@/app/theme";
import { usePieceActivities } from "../hooks/useStudioQueries";
import type { Phase, PieceActivity } from "../types";

interface PieceActivityFeedProps {
  pieceId: number;
  /** Used to show phase names for older activity rows that only stored phase ids. */
  phases?: Phase[];
}

export function PieceActivityFeed({ pieceId, phases }: PieceActivityFeedProps) {
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
            <ActivityPayload activity={a} phases={phases} />
          </Stack>
        </Timeline.Item>
      ))}
    </Timeline>
  );
}

function numericId(v: unknown): number | null {
  if (typeof v === "number" && Number.isFinite(v)) return v;
  if (typeof v === "string" && v.trim() !== "" && !Number.isNaN(Number(v))) {
    return Number(v);
  }
  return null;
}

function phaseNameFromId(
  id: unknown,
  phases: Phase[] | undefined
): string | null {
  const n = numericId(id);
  if (n == null || !phases?.length) return null;
  return phases.find((p) => p.id === n)?.name ?? null;
}

function ActivityPayload({
  activity,
  phases,
}: {
  activity: PieceActivity;
  phases?: Phase[];
}) {
  const payload = activity.payload;
  if (!payload || typeof payload !== "object") return null;

  if (activity.kind === "moved") {
    return (
      <MovedActivityPayload
        record={payload as Record<string, unknown>}
        phases={phases}
      />
    );
  }

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

function MovedActivityPayload({
  record,
  phases,
}: {
  record: Record<string, unknown>;
  phases?: Phase[];
}) {
  const fromStr =
    (typeof record.from_phase === "string" && record.from_phase) ||
    phaseNameFromId(record.from_phase_id, phases);
  const toStr =
    (typeof record.to_phase === "string" && record.to_phase) ||
    phaseNameFromId(record.to_phase_id, phases);
  const position =
    typeof record.position === "string" && record.position
      ? record.position
      : null;
  const reason =
    typeof record.reason === "string" && record.reason
      ? record.reason
      : null;

  const parts: string[] = [];
  if (fromStr && toStr) {
    parts.push(`From ${fromStr} to ${toStr}`);
  } else if (toStr) {
    parts.push(`To ${toStr}`);
  } else {
    // Fallback: show raw fields if we cannot resolve names
    const rest = { ...record };
    delete rest.from_phase;
    delete rest.to_phase;
    const entries = Object.entries(rest).filter(
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
  if (position) parts.push(`position: ${position}`);
  if (reason) parts.push(`reason: ${reason}`);

  if (parts.length === 0) return null;
  return (
    <Text size="xs" c="dimmed" style={{ whiteSpace: "pre-wrap" }}>
      {parts.join("  •  ")}
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
