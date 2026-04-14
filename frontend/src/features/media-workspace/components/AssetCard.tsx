import {
  Badge,
  Card,
  Group,
  Image,
  Text,
  ActionIcon,
  Menu,
} from "@mantine/core";
import {
  IconDots,
  IconTrash,
  IconPhoto,
  IconVideo,
} from "@tabler/icons-react";
import type { MediaAsset } from "../types";

interface AssetCardProps {
  asset: MediaAsset;
  onClick: (id: number) => void;
  onDelete: (id: number) => void;
}

const statusColor: Record<string, string> = {
  draft: "gray",
  ready: "blue",
  published: "green",
};

export function AssetCard({ asset, onClick, onDelete }: AssetCardProps) {
  const currentVersion = asset.versions?.find(
    (v) => v.id === asset.current_version_id
  );
  const thumbnailUrl = currentVersion?.blob_url;
  const isVideo = asset.asset_type === "video";

  return (
    <Card
      shadow="sm"
      padding="sm"
      radius="md"
      withBorder
      style={{ cursor: "pointer" }}
      onClick={() => onClick(asset.id)}
    >
      <Card.Section>
        {thumbnailUrl && !isVideo ? (
          <Image src={thumbnailUrl} height={180} alt={asset.title} fit="cover" />
        ) : thumbnailUrl && isVideo ? (
          <video
            src={thumbnailUrl}
            height={180}
            style={{
              width: "100%",
              height: 180,
              objectFit: "cover",
              display: "block",
            }}
            muted
          />
        ) : (
          <div
            style={{
              height: 180,
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              backgroundColor: "var(--mantine-color-gray-1)",
            }}
          >
            {isVideo ? <IconVideo size={48} opacity={0.3} /> : <IconPhoto size={48} opacity={0.3} />}
          </div>
        )}
      </Card.Section>

      <Group justify="space-between" mt="sm" mb={4}>
        <Text fw={500} lineClamp={1} style={{ flex: 1 }}>
          {asset.title}
        </Text>
        <Menu position="bottom-end" withArrow>
          <Menu.Target>
            <ActionIcon
              variant="subtle"
              size="sm"
              onClick={(e) => e.stopPropagation()}
            >
              <IconDots size={16} />
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
              Delete
            </Menu.Item>
          </Menu.Dropdown>
        </Menu>
      </Group>

      <Group gap="xs">
        <Badge color={statusColor[asset.status] ?? "gray"} size="sm">
          {asset.status}
        </Badge>
        <Badge variant="light" size="sm">
          {asset.asset_type}
        </Badge>
        {asset.versions && (
          <Badge variant="outline" size="sm">
            {asset.versions.length} version{asset.versions.length !== 1 ? "s" : ""}
          </Badge>
        )}
      </Group>
    </Card>
  );
}
