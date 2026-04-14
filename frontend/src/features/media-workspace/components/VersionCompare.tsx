import { Badge, Group, Image, Stack, Text } from "@mantine/core";
import type { MediaVersion } from "../types";

interface VersionCompareProps {
  left: MediaVersion | null;
  right: MediaVersion | null;
}

function VersionPanel({ version, label }: { version: MediaVersion | null; label: string }) {
  if (!version) {
    return (
      <Stack
        align="center"
        justify="center"
        style={{
          flex: 1,
          minHeight: 300,
          background: "var(--mantine-color-gray-1)",
          borderRadius: 8,
        }}
      >
        <Text c="dimmed" size="sm">
          Select a version for {label}
        </Text>
      </Stack>
    );
  }

  const isVideo = version.content_type.startsWith("video/");

  return (
    <Stack style={{ flex: 1 }} gap="xs">
      <Group gap="xs">
        <Badge size="sm">{label}</Badge>
        <Text size="sm" fw={500}>
          v{version.version_number} — {version.label}
        </Text>
      </Group>

      {isVideo ? (
        <video
          src={version.blob_url}
          controls
          style={{
            width: "100%",
            maxHeight: 400,
            borderRadius: 8,
            background: "#000",
          }}
        />
      ) : (
        <Image
          src={version.blob_url}
          alt={version.label}
          fit="contain"
          h={400}
          radius="md"
          style={{ background: "var(--mantine-color-gray-1)" }}
        />
      )}

      <Text size="xs" c="dimmed">
        {version.file_name} &middot; {(version.file_size / (1024 * 1024)).toFixed(1)} MB
      </Text>
    </Stack>
  );
}

export function VersionCompare({ left, right }: VersionCompareProps) {
  return (
    <Group grow align="flex-start" gap="md">
      <VersionPanel version={left} label="Before" />
      <VersionPanel version={right} label="After" />
    </Group>
  );
}
