import {
  Badge,
  Box,
  Center,
  Group,
  Image,
  Paper,
  Stack,
  Text,
  ThemeIcon,
  useMantineTheme,
} from "@mantine/core";
import { IconArrowRight, IconPhoto } from "@tabler/icons-react";
import type { MediaVersion } from "../types";
import { versionMediaUrl } from "../types";
import { useMediaWorkspaceSurfaces } from "../hooks/useMediaWorkspaceSurfaces";

interface VersionCompareProps {
  left: MediaVersion | null;
  right: MediaVersion | null;
}

function VersionPanel({
  version,
  label,
  color,
}: {
  version: MediaVersion | null;
  label: string;
  color: "orange" | "teal";
}) {
  const theme = useMantineTheme();
  const surfaces = useMediaWorkspaceSurfaces();
  const tint = color === "orange" ? theme.colors.orange : theme.colors.teal;

  if (!version) {
    return (
      <Paper
        p="xl"
        radius="lg"
        withBorder
        style={{
          flex: 1,
          minHeight: 340,
          borderStyle: "dashed",
          borderWidth: 2,
          borderColor: surfaces.isDark ? tint[7] : tint[3],
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          background: surfaces.isDark ? theme.colors.dark[7] : tint[0],
        }}
      >
        <Stack align="center" gap="sm">
          <ThemeIcon variant="light" color={color} size={48} radius="xl">
            <IconPhoto size={24} />
          </ThemeIcon>
          <Text c="dimmed" size="sm" fw={500}>
            Select a version for {label}
          </Text>
          <Text size="xs" c="dimmed">
            Click a version in the timeline
          </Text>
        </Stack>
      </Paper>
    );
  }

  const isVideo = version.content_type.startsWith("video/");

  return (
    <Paper
      radius="lg"
      withBorder
      style={{ flex: 1, overflow: "hidden" }}
    >
      {/* Header */}
      <Group
        p="xs"
        px="sm"
        gap="xs"
        style={{
          borderBottom: `1px solid ${
            surfaces.isDark ? theme.colors.dark[4] : theme.colors.gray[2]
          }`,
        }}
      >
        <Badge size="sm" variant="filled" color={color} radius="sm">
          {label}
        </Badge>
        <Text size="sm" fw={600}>
          v{version.version_number}
        </Text>
        <Badge size="xs" variant="light" radius="sm">
          {version.label}
        </Badge>
      </Group>

      {/* Media */}
      {isVideo ? (
        <video
          src={versionMediaUrl(version) ?? undefined}
          controls
          style={{
            width: "100%",
            maxHeight: 380,
            display: "block",
            background: "#000",
          }}
        />
      ) : (
        <Box
          style={{
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            minHeight: 300,
            maxHeight: 380,
            overflow: "hidden",
            background: surfaces.checkerboardCss,
          }}
        >
          <Image
            src={versionMediaUrl(version) ?? undefined}
            alt={version.label}
            fit="contain"
            mah={380}
            style={{ display: "block" }}
          />
        </Box>
      )}

      {/* Footer metadata */}
      <Group
        p="xs"
        px="sm"
        gap="xs"
        style={{
          borderTop: surfaces.previewMetaBorderTop,
          background: surfaces.previewMetaBg,
        }}
      >
        <Text size="xs" c="dimmed" lineClamp={1} style={{ flex: 1 }}>
          {version.file_name}
        </Text>
        <Text size="xs" c="dimmed">
          {(version.file_size / (1024 * 1024)).toFixed(1)} MB
        </Text>
      </Group>
    </Paper>
  );
}

export function VersionCompare({ left, right }: VersionCompareProps) {
  return (
    <Stack gap="sm">
      <Group gap="xs" justify="center">
        <Badge variant="light" color="orange" radius="sm">
          Compare Mode
        </Badge>
        <Text size="xs" c="dimmed">
          Select two versions from the timeline to compare
        </Text>
      </Group>

      <Group grow align="flex-start" gap="md" wrap="nowrap">
        <VersionPanel version={left} label="Before" color="orange" />

        <Center style={{ flexShrink: 0, alignSelf: "center" }}>
          <ThemeIcon variant="light" color="gray" size="lg" radius="xl">
            <IconArrowRight size={18} />
          </ThemeIcon>
        </Center>

        <VersionPanel version={right} label="After" color="teal" />
      </Group>
    </Stack>
  );
}
