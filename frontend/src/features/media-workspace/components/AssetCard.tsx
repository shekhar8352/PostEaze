import {
  Badge,
  Card,
  Group,
  Image,
  Overlay,
  Stack,
  Text,
  ActionIcon,
  Menu,
  ThemeIcon,
  Box,
} from "@mantine/core";
import {
  IconDots,
  IconTrash,
  IconPhoto,
  IconVideo,
  IconClock,
  IconVersions,
} from "@tabler/icons-react";
import type { MediaAsset } from "../types";

interface AssetCardProps {
  asset: MediaAsset;
  onClick: (id: number) => void;
  onDelete: (id: number) => void;
}

const statusConfig: Record<string, { color: string; label: string }> = {
  draft: { color: "gray", label: "Draft" },
  ready: { color: "blue", label: "Ready" },
  published: { color: "green", label: "Published" },
};

function formatFileSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(0)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

function formatRelativeDate(iso: string): string {
  const diff = Date.now() - new Date(iso).getTime();
  const mins = Math.floor(diff / 60000);
  if (mins < 1) return "Just now";
  if (mins < 60) return `${mins}m ago`;
  const hrs = Math.floor(mins / 60);
  if (hrs < 24) return `${hrs}h ago`;
  const days = Math.floor(hrs / 24);
  if (days < 7) return `${days}d ago`;
  return new Date(iso).toLocaleDateString(undefined, {
    month: "short",
    day: "numeric",
  });
}

export function AssetCard({ asset, onClick, onDelete }: AssetCardProps) {
  const versions = asset.versions;
  const currentVersion =
    versions?.find((v) => v.id === asset.current_version_id) ??
    (versions?.length ? versions[versions.length - 1] : undefined);
  const thumbnailUrl = currentVersion?.blob_url;
  const isVideo = asset.asset_type === "video";
  const cfg = statusConfig[asset.status] ?? statusConfig.draft;

  return (
    <Card
      shadow="sm"
      padding={0}
      radius="lg"
      withBorder
      style={{
        cursor: "pointer",
        overflow: "hidden",
        transition: "transform 200ms ease, box-shadow 200ms ease",
      }}
      onMouseEnter={(e) => {
        e.currentTarget.style.transform = "translateY(-4px)";
        e.currentTarget.style.boxShadow =
          "0 12px 24px rgba(0,0,0,0.1), 0 4px 8px rgba(0,0,0,0.06)";
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.transform = "translateY(0)";
        e.currentTarget.style.boxShadow = "";
      }}
      onClick={() => onClick(asset.id)}
    >
      {/* Thumbnail */}
      <Box pos="relative" style={{ overflow: "hidden" }}>
        {thumbnailUrl && !isVideo ? (
          <Image
            src={thumbnailUrl}
            height={200}
            alt={asset.title}
            fit="cover"
            fallbackSrc="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='1' height='1'%3E%3Crect fill='%23f1f3f5' width='1' height='1'/%3E%3C/svg%3E"
          />
        ) : thumbnailUrl && isVideo ? (
          <video
            src={thumbnailUrl}
            style={{
              width: "100%",
              height: 200,
              objectFit: "cover",
              display: "block",
            }}
            muted
            playsInline
          />
        ) : (
          <Box
            style={{
              height: 200,
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              background:
                "linear-gradient(135deg, var(--mantine-color-gray-1) 0%, var(--mantine-color-gray-2) 100%)",
            }}
          >
            <ThemeIcon
              variant="light"
              color="gray"
              size={64}
              radius="xl"
              style={{ opacity: 0.5 }}
            >
              {isVideo ? <IconVideo size={32} /> : <IconPhoto size={32} />}
            </ThemeIcon>
          </Box>
        )}

        {/* Hover overlay */}
        <Overlay
          gradient="linear-gradient(0deg, rgba(0,0,0,0.4) 0%, transparent 60%)"
          opacity={0}
          style={{
            transition: "opacity 200ms ease",
            pointerEvents: "none",
          }}
          className="asset-card-overlay"
        />

        {/* Type badge pinned top-left */}
        <Badge
          size="sm"
          variant="filled"
          color={isVideo ? "violet" : "blue"}
          style={{
            position: "absolute",
            top: 10,
            left: 10,
            textTransform: "uppercase",
            letterSpacing: "0.05em",
            fontWeight: 700,
            fontSize: 10,
          }}
          leftSection={
            isVideo ? <IconVideo size={12} /> : <IconPhoto size={12} />
          }
        >
          {asset.asset_type}
        </Badge>

        {/* Menu pinned top-right */}
        <Menu position="bottom-end" withArrow shadow="md">
          <Menu.Target>
            <ActionIcon
              variant="filled"
              color="gray"
              size="sm"
              radius="xl"
              style={{
                position: "absolute",
                top: 10,
                right: 10,
                opacity: 0.7,
                transition: "opacity 150ms",
              }}
              onClick={(e) => e.stopPropagation()}
              onMouseEnter={(e) => {
                e.currentTarget.style.opacity = "1";
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.opacity = "0.7";
              }}
            >
              <IconDots size={14} />
            </ActionIcon>
          </Menu.Target>
          <Menu.Dropdown>
            <Menu.Item
              color="red"
              leftSection={<IconTrash size={14} />}
              onClick={(e) => {
                e.stopPropagation();
                onDelete(asset.id);
              }}
            >
              Delete asset
            </Menu.Item>
          </Menu.Dropdown>
        </Menu>
      </Box>

      {/* Content area */}
      <Stack gap={6} p="sm">
        <Group justify="space-between" wrap="nowrap">
          <Text fw={600} size="sm" lineClamp={1} style={{ flex: 1 }}>
            {asset.title}
          </Text>
          <Badge
            size="xs"
            variant="light"
            color={cfg.color}
            radius="sm"
            style={{ flexShrink: 0 }}
          >
            {cfg.label}
          </Badge>
        </Group>

        <Group gap="lg" wrap="nowrap">
          <Group gap={4} wrap="nowrap">
            <IconClock
              size={13}
              style={{ opacity: 0.45, flexShrink: 0 }}
            />
            <Text size="xs" c="dimmed">
              {formatRelativeDate(asset.updated_at)}
            </Text>
          </Group>

          {currentVersion && (
            <Text size="xs" c="dimmed">
              {formatFileSize(currentVersion.file_size)}
            </Text>
          )}

          {asset.versions && asset.versions.length > 1 && (
            <Group gap={4} wrap="nowrap">
              <IconVersions
                size={13}
                style={{ opacity: 0.45, flexShrink: 0 }}
              />
              <Text size="xs" c="dimmed">
                {asset.versions.length}
              </Text>
            </Group>
          )}
        </Group>
      </Stack>
    </Card>
  );
}
