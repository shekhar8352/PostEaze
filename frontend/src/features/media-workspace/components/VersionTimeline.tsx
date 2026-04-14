import {
  ActionIcon,
  Badge,
  Group,
  Paper,
  Stack,
  Text,
  ThemeIcon,
  Timeline,
  Tooltip,
} from "@mantine/core";
import { IconCheck, IconTrash, IconClock } from "@tabler/icons-react";
import type { MediaVersion } from "../types";
import { useMediaWorkspaceSurfaces } from "../hooks/useMediaWorkspaceSurfaces";

interface VersionTimelineProps {
  versions: MediaVersion[];
  currentVersionId: number | null;
  onSelect: (version: MediaVersion) => void;
  onSetCurrent: (versionId: number) => void;
  onDelete: (versionId: number) => void;
}

function formatRelativeDate(iso: string): string {
  const diff = Date.now() - new Date(iso).getTime();
  const mins = Math.floor(diff / 60000);
  if (mins < 1) return "Just now";
  if (mins < 60) return `${mins}m ago`;
  const hrs = Math.floor(mins / 60);
  if (hrs < 24) return `${hrs}h ago`;
  const days = Math.floor(hrs / 24);
  if (days < 30) return `${days}d ago`;
  return new Date(iso).toLocaleDateString(undefined, {
    month: "short",
    day: "numeric",
  });
}

export function VersionTimeline({
  versions,
  currentVersionId,
  onSelect,
  onSetCurrent,
  onDelete,
}: VersionTimelineProps) {
  const { subtleHoverBg } = useMediaWorkspaceSurfaces();

  return (
    <Timeline active={-1} bulletSize={28} lineWidth={2}>
      {versions.map((v) => {
        const isCurrent = v.id === currentVersionId;
        return (
          <Timeline.Item
            key={v.id}
            bullet={
              isCurrent ? (
                <ThemeIcon color="blue" size={28} radius="xl">
                  <IconCheck size={14} />
                </ThemeIcon>
              ) : (
                <ThemeIcon variant="light" color="gray" size={28} radius="xl">
                  <IconClock size={14} />
                </ThemeIcon>
              )
            }
            color={isCurrent ? "blue" : "gray"}
          >
            <Paper
              p="xs"
              px="sm"
              radius="md"
              withBorder={isCurrent}
              shadow={isCurrent ? "xs" : undefined}
              style={{
                cursor: "pointer",
                transition:
                  "background 150ms ease, transform 150ms ease, box-shadow 150ms ease",
                borderColor: isCurrent
                  ? "var(--mantine-color-blue-4)"
                  : undefined,
              }}
              onMouseEnter={(e) => {
                if (!isCurrent) {
                  e.currentTarget.style.background = subtleHoverBg;
                }
                e.currentTarget.style.transform = "translateX(2px)";
              }}
              onMouseLeave={(e) => {
                if (!isCurrent) {
                  e.currentTarget.style.background = "";
                }
                e.currentTarget.style.transform = "translateX(0)";
              }}
              onClick={() => onSelect(v)}
            >
              <Group justify="space-between" wrap="nowrap">
                <Stack gap={2} style={{ minWidth: 0, flex: 1 }}>
                  <Group gap={6}>
                    <Text fw={600} size="sm">
                      v{v.version_number}
                    </Text>
                    <Badge
                      size="xs"
                      variant="light"
                      radius="sm"
                      color={isCurrent ? "blue" : "gray"}
                    >
                      {v.label}
                    </Badge>
                    {isCurrent && (
                      <Badge
                        size="xs"
                        variant="filled"
                        color="blue"
                        radius="sm"
                      >
                        Active
                      </Badge>
                    )}
                  </Group>
                  <Group gap={4}>
                    <Text size="xs" c="dimmed">
                      {formatRelativeDate(v.created_at)}
                    </Text>
                    <Text size="xs" c="dimmed">
                      ·
                    </Text>
                    <Text size="xs" c="dimmed">
                      {(v.file_size / (1024 * 1024)).toFixed(1)} MB
                    </Text>
                  </Group>
                  {v.notes && (
                    <Text size="xs" c="dimmed" lineClamp={1} mt={2}>
                      {v.notes}
                    </Text>
                  )}
                </Stack>

                <Group gap={4} style={{ flexShrink: 0 }}>
                  {!isCurrent && (
                    <Tooltip label="Set as active" withArrow>
                      <ActionIcon
                        variant="light"
                        color="blue"
                        size="sm"
                        radius="md"
                        onClick={(e) => {
                          e.stopPropagation();
                          onSetCurrent(v.id);
                        }}
                      >
                        <IconCheck size={14} />
                      </ActionIcon>
                    </Tooltip>
                  )}
                  <Tooltip label="Delete version" withArrow>
                    <ActionIcon
                      variant="subtle"
                      color="red"
                      size="sm"
                      radius="md"
                      onClick={(e) => {
                        e.stopPropagation();
                        onDelete(v.id);
                      }}
                    >
                      <IconTrash size={14} />
                    </ActionIcon>
                  </Tooltip>
                </Group>
              </Group>
            </Paper>
          </Timeline.Item>
        );
      })}
    </Timeline>
  );
}
