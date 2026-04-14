import {
  ActionIcon,
  Badge,
  Card,
  Group,
  Stack,
  Text,
  Timeline,
  Tooltip,
} from "@mantine/core";
import { IconCheck, IconTrash } from "@tabler/icons-react";
import type { MediaVersion } from "../types";

interface VersionTimelineProps {
  versions: MediaVersion[];
  currentVersionId: number | null;
  onSelect: (version: MediaVersion) => void;
  onSetCurrent: (versionId: number) => void;
  onDelete: (versionId: number) => void;
}

export function VersionTimeline({
  versions,
  currentVersionId,
  onSelect,
  onSetCurrent,
  onDelete,
}: VersionTimelineProps) {
  return (
    <Timeline active={-1} bulletSize={24} lineWidth={2}>
      {versions.map((v) => {
        const isCurrent = v.id === currentVersionId;
        return (
          <Timeline.Item
            key={v.id}
            bullet={isCurrent ? <IconCheck size={12} /> : undefined}
            color={isCurrent ? "blue" : "gray"}
          >
            <Card
              padding="xs"
              radius="sm"
              withBorder={isCurrent}
              style={{ cursor: "pointer" }}
              onClick={() => onSelect(v)}
            >
              <Group justify="space-between" wrap="nowrap">
                <Stack gap={2}>
                  <Group gap="xs">
                    <Text fw={600} size="sm">
                      v{v.version_number}
                    </Text>
                    <Badge size="xs" variant="light">
                      {v.label}
                    </Badge>
                    {isCurrent && (
                      <Badge size="xs" color="blue">
                        active
                      </Badge>
                    )}
                  </Group>
                  <Text size="xs" c="dimmed">
                    {new Date(v.created_at).toLocaleDateString()} &middot;{" "}
                    {(v.file_size / (1024 * 1024)).toFixed(1)} MB
                  </Text>
                  {v.notes && (
                    <Text size="xs" c="dimmed" lineClamp={1}>
                      {v.notes}
                    </Text>
                  )}
                </Stack>

                <Group gap={4}>
                  {!isCurrent && (
                    <Tooltip label="Set as active">
                      <ActionIcon
                        variant="light"
                        size="sm"
                        onClick={(e) => {
                          e.stopPropagation();
                          onSetCurrent(v.id);
                        }}
                      >
                        <IconCheck size={14} />
                      </ActionIcon>
                    </Tooltip>
                  )}
                  <Tooltip label="Delete version">
                    <ActionIcon
                      variant="light"
                      color="red"
                      size="sm"
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
            </Card>
          </Timeline.Item>
        );
      })}
    </Timeline>
  );
}
