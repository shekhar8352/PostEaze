import { useMemo, useState } from "react";
import {
  ActionIcon,
  Badge,
  Button,
  Card,
  Group,
  Image,
  Loader,
  Modal,
  ScrollArea,
  Select,
  SimpleGrid,
  Stack,
  Text,
  ThemeIcon,
  Tooltip,
} from "@mantine/core";
import { useQueries, useQueryClient } from "@tanstack/react-query";
import { Icons } from "@/app/theme";
import type { MediaAsset } from "@/features/media-workspace/types";
import { mediaApi } from "@/features/media-workspace/services/mediaApi";
import {
  mediaKeys,
  useMediaAssets,
} from "@/features/media-workspace/hooks/useMediaQueries";
import type { PieceAssetLink, PieceAssetRole } from "../types";
import { useLinkAsset, useUnlinkAsset } from "../hooks/useStudioQueries";

interface PieceAssetsTabProps {
  pieceId: number;
  assets: PieceAssetLink[];
}

const ROLE_OPTIONS: { value: PieceAssetRole; label: string }[] = [
  { value: "final", label: "Final" },
  { value: "edit", label: "Edit" },
  { value: "raw", label: "Raw" },
  { value: "thumbnail", label: "Thumbnail" },
  { value: "script", label: "Script" },
  { value: "attachment", label: "Attachment" },
];

const ROLE_COLOR: Record<PieceAssetRole, string> = {
  final: "green",
  edit: "blue",
  raw: "gray",
  thumbnail: "violet",
  script: "yellow",
  attachment: "cyan",
};

export function PieceAssetsTab({ pieceId, assets }: PieceAssetsTabProps) {
  const qc = useQueryClient();
  const unlink = useUnlinkAsset();
  const [pickerOpen, setPickerOpen] = useState(false);

  // Fetch details for each linked asset in parallel. Cached entries are
  // reused across renders so toggling tabs doesn't re-hit the API.
  const assetQueries = useQueries({
    queries: assets.map((link) => ({
      queryKey: mediaKeys.detail(link.media_asset_id),
      queryFn: () => mediaApi.get(link.media_asset_id),
    })),
  });

  const loading = assetQueries.some((q) => q.isLoading);

  const resolved = useMemo(
    () =>
      assets.map((link, i) => ({
        link,
        asset: assetQueries[i]?.data,
        error: assetQueries[i]?.error as Error | undefined,
      })),
    [assets, assetQueries]
  );

  const handleUnlink = (mediaAssetId: number, role: PieceAssetRole) => {
    unlink.mutate(
      { pieceId, mediaAssetId, role },
      {
        onSuccess: () => {
          void qc.invalidateQueries({ queryKey: mediaKeys.detail(mediaAssetId) });
        },
      }
    );
  };

  return (
    <Stack gap="md">
      <Group justify="space-between" align="center">
        <Text size="sm" c="dimmed">
          {assets.length === 0
            ? "No assets linked yet."
            : `${assets.length} asset${assets.length === 1 ? "" : "s"} linked`}
        </Text>
        <Button
          size="xs"
          leftSection={<Icons.Plus size={14} />}
          onClick={() => setPickerOpen(true)}
        >
          Link asset
        </Button>
      </Group>

      {loading && assets.length > 0 && (
        <Group gap="xs">
          <Loader size="xs" />
          <Text size="xs" c="dimmed">
            Loading assets…
          </Text>
        </Group>
      )}

      {resolved.length === 0 ? (
        <Card withBorder radius="md" p="lg">
          <Stack align="center" gap="xs">
            <ThemeIcon variant="light" color="gray" size={40} radius="xl">
              <Icons.Photo size={20} />
            </ThemeIcon>
            <Text size="sm" c="dimmed" ta="center">
              Link photos, videos, scripts, and other files to keep every
              deliverable for this piece in one place.
            </Text>
          </Stack>
        </Card>
      ) : (
        <SimpleGrid cols={{ base: 1, sm: 2 }} spacing="sm">
          {resolved.map(({ link, asset, error }) => (
            <LinkedAssetCard
              key={`${link.media_asset_id}:${link.role}`}
              link={link}
              asset={asset}
              error={error}
              onUnlink={() => handleUnlink(link.media_asset_id, link.role)}
              unlinking={unlink.isPending}
            />
          ))}
        </SimpleGrid>
      )}

      <AssetPickerModal
        pieceId={pieceId}
        opened={pickerOpen}
        onClose={() => setPickerOpen(false)}
        existingIds={assets.map((a) => a.media_asset_id)}
      />
    </Stack>
  );
}

interface LinkedAssetCardProps {
  link: PieceAssetLink;
  asset?: MediaAsset;
  error?: Error;
  onUnlink: () => void;
  unlinking: boolean;
}

function LinkedAssetCard({
  link,
  asset,
  error,
  onUnlink,
  unlinking,
}: LinkedAssetCardProps) {
  const thumb = resolveThumbnail(asset);
  const isVideo = asset?.asset_type === "video";
  const title = asset?.title ?? (error ? "Unavailable" : `Asset #${link.media_asset_id}`);

  return (
    <Card withBorder radius="md" p="xs">
      <Group wrap="nowrap" align="flex-start" gap="sm">
        <ThumbBox src={thumb} isVideo={isVideo} />
        <Stack gap={4} style={{ flex: 1, minWidth: 0 }}>
          <Text fw={500} size="sm" lineClamp={1}>
            {title}
          </Text>
          <Group gap="xs">
            <Badge
              size="xs"
              variant="light"
              color={ROLE_COLOR[link.role] ?? "gray"}
            >
              {link.role}
            </Badge>
            {asset && (
              <Badge size="xs" variant="outline" color="gray">
                {asset.asset_type}
              </Badge>
            )}
          </Group>
          {error && (
            <Text size="xs" c="red">
              Failed to load
            </Text>
          )}
        </Stack>
        <Tooltip label="Unlink from piece">
          <ActionIcon
            variant="subtle"
            color="red"
            onClick={onUnlink}
            loading={unlinking}
            aria-label="Unlink asset"
          >
            <Icons.Trash size={14} />
          </ActionIcon>
        </Tooltip>
      </Group>
    </Card>
  );
}

function ThumbBox({ src, isVideo }: { src?: string; isVideo?: boolean }) {
  if (src && !isVideo) {
    return (
      <Image
        src={src}
        alt=""
        w={64}
        h={64}
        radius="sm"
        fit="cover"
        fallbackSrc="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='1' height='1'%3E%3Crect fill='%23f1f3f5' width='1' height='1'/%3E%3C/svg%3E"
      />
    );
  }
  if (src && isVideo) {
    return (
      <video
        src={src}
        muted
        playsInline
        style={{
          width: 64,
          height: 64,
          objectFit: "cover",
          borderRadius: 4,
          background: "var(--mantine-color-gray-1)",
        }}
      />
    );
  }
  return (
    <ThemeIcon
      variant="light"
      color="gray"
      size={64}
      radius="sm"
      style={{ flexShrink: 0 }}
    >
      {isVideo ? <Icons.Video size={24} /> : <Icons.Photo size={24} />}
    </ThemeIcon>
  );
}

function resolveThumbnail(asset: MediaAsset | undefined) {
  if (!asset?.versions?.length) return undefined;
  const current =
    asset.versions.find((v) => v.id === asset.current_version_id) ??
    asset.versions[asset.versions.length - 1];
  return current?.blob_url;
}

interface AssetPickerModalProps {
  pieceId: number;
  opened: boolean;
  onClose: () => void;
  existingIds: number[];
}

function AssetPickerModal({
  pieceId,
  opened,
  onClose,
  existingIds,
}: AssetPickerModalProps) {
  const { data, isLoading } = useMediaAssets(undefined, 50, 0);
  const link = useLinkAsset();
  const [role, setRole] = useState<PieceAssetRole>("final");

  const existingSet = useMemo(() => new Set(existingIds), [existingIds]);
  const available = (data?.assets ?? []).filter(
    (a) => !existingSet.has(a.id)
  );

  const handleLink = async (asset: MediaAsset) => {
    await link.mutateAsync({
      pieceId,
      body: { media_asset_id: asset.id, role },
    });
  };

  return (
    <Modal
      opened={opened}
      onClose={onClose}
      title="Link existing asset"
      size="lg"
      centered
    >
      <Stack gap="md">
        <Select
          label="Role for linked assets"
          value={role}
          onChange={(v) => v && setRole(v as PieceAssetRole)}
          data={ROLE_OPTIONS}
          allowDeselect={false}
        />

        {isLoading ? (
          <Group>
            <Loader size="sm" />
            <Text size="sm" c="dimmed">
              Loading library…
            </Text>
          </Group>
        ) : available.length === 0 ? (
          <Text size="sm" c="dimmed">
            No unlinked assets available. Upload one from the Media
            workspace, then return here.
          </Text>
        ) : (
          <ScrollArea h={400} type="auto">
            <Stack gap="xs">
              {available.map((asset) => (
                <PickerRow
                  key={asset.id}
                  asset={asset}
                  onLink={() => handleLink(asset)}
                  linking={link.isPending}
                />
              ))}
            </Stack>
          </ScrollArea>
        )}
      </Stack>
    </Modal>
  );
}

function PickerRow({
  asset,
  onLink,
  linking,
}: {
  asset: MediaAsset;
  onLink: () => void;
  linking: boolean;
}) {
  const thumb = resolveThumbnail(asset);
  const isVideo = asset.asset_type === "video";
  return (
    <Card withBorder radius="md" p="xs">
      <Group wrap="nowrap" align="center" gap="sm">
        <ThumbBox src={thumb} isVideo={isVideo} />
        <Stack gap={2} style={{ flex: 1, minWidth: 0 }}>
          <Text fw={500} size="sm" lineClamp={1}>
            {asset.title}
          </Text>
          <Group gap="xs">
            <Badge size="xs" variant="outline" color="gray">
              {asset.asset_type}
            </Badge>
            <Badge size="xs" variant="light" color="gray">
              {asset.status}
            </Badge>
          </Group>
        </Stack>
        <Button size="xs" variant="light" onClick={onLink} loading={linking}>
          Link
        </Button>
      </Group>
    </Card>
  );
}
