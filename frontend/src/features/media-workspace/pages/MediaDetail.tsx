import { useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import {
  ActionIcon,
  Badge,
  Box,
  Button,
  Center,
  Divider,
  Grid,
  Group,
  Image,
  Loader,
  Menu,
  Paper,
  Stack,
  Text,
  ThemeIcon,
  Title,
  Tooltip,
} from "@mantine/core";
import {
  IconArrowLeft,
  IconPlus,
  IconColumns,
  IconCalendar,
  IconTrash,
  IconDots,
  IconPhoto,
  IconVideo,
  IconDownload,
} from "@tabler/icons-react";
import { Link } from "react-router-dom";
import {
  useMediaAsset,
  useDeleteMediaAsset,
  useSetCurrentVersion,
  useDeleteVersion,
} from "../hooks/useMediaQueries";
import { useMediaWorkspaceSurfaces } from "../hooks/useMediaWorkspaceSurfaces";
import { VersionTimeline } from "../components/VersionTimeline";
import { VersionCompare } from "../components/VersionCompare";
import { AddVersionModal } from "../components/AddVersionModal";
import { DriveRevisionsPanel } from "../components/DriveRevisionsPanel";
import type { MediaVersion } from "../types";
import { versionMediaUrl } from "../types";

const statusConfig: Record<string, { color: string; label: string }> = {
  draft: { color: "gray", label: "Draft" },
  ready: { color: "blue", label: "Ready" },
  published: { color: "green", label: "Published" },
};

export default function MediaDetail() {
  const { assetId } = useParams<{ assetId: string }>();
  const navigate = useNavigate();
  const id = Number(assetId);

  const { data: asset, isLoading } = useMediaAsset(id);
  const deleteAsset = useDeleteMediaAsset();
  const setCurrent = useSetCurrentVersion();
  const deleteVer = useDeleteVersion();
  const surfaces = useMediaWorkspaceSurfaces();

  const [selectedVersion, setSelectedVersion] =
    useState<MediaVersion | null>(null);
  const [compareMode, setCompareMode] = useState(false);
  const [compareLeft, setCompareLeft] = useState<MediaVersion | null>(null);
  const [addVersionOpen, setAddVersionOpen] = useState(false);

  if (isLoading) {
    return (
      <Center mih={500}>
        <Stack align="center" gap="md">
          <Loader size="lg" type="dots" />
          <Text size="sm" c="dimmed">
            Loading asset…
          </Text>
        </Stack>
      </Center>
    );
  }

  if (!asset) {
    return (
      <Center mih={500}>
        <Paper p="xl" radius="xl" withBorder style={{ textAlign: "center" }}>
          <Stack align="center" gap="md">
            <ThemeIcon variant="light" color="gray" size={64} radius="xl">
              <IconPhoto size={32} />
            </ThemeIcon>
            <Text fw={600}>Asset not found</Text>
            <Text size="sm" c="dimmed">
              This asset may have been deleted or doesn't exist.
            </Text>
            <Button
              variant="light"
              onClick={() => navigate("/workspace")}
              radius="md"
            >
              Back to Workspace
            </Button>
          </Stack>
        </Paper>
      </Center>
    );
  }

  const versions = asset.versions ?? [];
  const currentVersion =
    selectedVersion ??
    versions.find((v) => v.id === asset.current_version_id) ??
    versions[versions.length - 1] ??
    null;
  const isVideo = asset.asset_type === "video";
  const cfg = statusConfig[asset.status] ?? statusConfig.draft;

  const handleVersionSelect = (v: MediaVersion) => {
    if (compareMode) {
      if (!compareLeft) {
        setCompareLeft(v);
      } else {
        setSelectedVersion(v);
      }
    } else {
      setSelectedVersion(v);
    }
  };

  const toggleCompare = () => {
    if (compareMode) {
      setCompareMode(false);
      setCompareLeft(null);
    } else {
      setCompareMode(true);
      setCompareLeft(currentVersion);
      setSelectedVersion(null);
    }
  };

  return (
    <Stack gap="lg" p="md">
      {/* Navigation breadcrumb bar */}
      <Group justify="space-between" wrap="wrap">
        <Group gap="sm">
          <ActionIcon
            variant="light"
            color="gray"
            size="lg"
            radius="md"
            onClick={() => navigate("/workspace")}
          >
            <IconArrowLeft size={18} />
          </ActionIcon>
          <Stack gap={0}>
            <Group gap="xs">
              <Title order={3} style={{ letterSpacing: "-0.01em" }}>
                {asset.title}
              </Title>
              <Badge
                variant="light"
                color={cfg.color}
                radius="sm"
                size="sm"
              >
                {cfg.label}
              </Badge>
              <Badge
                variant="dot"
                color={isVideo ? "violet" : "blue"}
                size="sm"
              >
                {asset.asset_type}
              </Badge>
            </Group>
            {currentVersion && (
              <Text size="xs" c="dimmed">
                v{currentVersion.version_number} — {currentVersion.label} —{" "}
                {new Date(currentVersion.created_at).toLocaleDateString()}
              </Text>
            )}
          </Stack>
        </Group>

        {/* Action buttons */}
        <Group gap="xs">
          <Tooltip label={compareMode ? "Exit compare" : "Compare versions"}>
            <Button
              variant={compareMode ? "filled" : "light"}
              size="sm"
              radius="md"
              leftSection={<IconColumns size={16} />}
              onClick={toggleCompare}
            >
              Compare
            </Button>
          </Tooltip>
          <Button
            variant="light"
            size="sm"
            radius="md"
            leftSection={<IconPlus size={16} />}
            onClick={() => setAddVersionOpen(true)}
          >
            Add Version
          </Button>
          <Button
            component={Link}
            to="/calendar"
            size="sm"
            radius="md"
            variant="light"
            leftSection={<IconCalendar size={16} />}
          >
            Schedule in Calendar
          </Button>
          <Menu position="bottom-end" withArrow shadow="md">
            <Menu.Target>
              <ActionIcon variant="light" color="gray" size="lg" radius="md">
                <IconDots size={18} />
              </ActionIcon>
            </Menu.Target>
            <Menu.Dropdown>
              {currentVersion && versionMediaUrl(currentVersion) && (
                <Menu.Item
                  leftSection={<IconDownload size={14} />}
                  component="a"
                  href={versionMediaUrl(currentVersion)!}
                  target="_blank"
                  rel="noopener"
                >
                  Download current
                </Menu.Item>
              )}
              <Menu.Divider />
              <Menu.Item
                color="red"
                leftSection={<IconTrash size={14} />}
                onClick={() => {
                  deleteAsset.mutate(id, {
                    onSuccess: () => navigate("/workspace"),
                  });
                }}
              >
                Delete asset
              </Menu.Item>
            </Menu.Dropdown>
          </Menu>
        </Group>
      </Group>

      <Divider />

      {/* Main content grid */}
      <Grid gutter="xl">
        {/* Preview area */}
        <Grid.Col span={{ base: 12, md: 8 }}>
          {compareMode ? (
            <VersionCompare left={compareLeft} right={selectedVersion} />
          ) : currentVersion ? (
            <Paper radius="lg" withBorder style={{ overflow: "hidden" }}>
              {isVideo ? (
                <video
                  key={currentVersion.id}
                  src={versionMediaUrl(currentVersion) ?? undefined}
                  controls
                  style={{
                    width: "100%",
                    maxHeight: 560,
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
                    minHeight: 400,
                    maxHeight: 560,
                    overflow: "hidden",
                    background: surfaces.checkerboardCss,
                  }}
                >
                  <Image
                    key={currentVersion.id}
                    src={versionMediaUrl(currentVersion) ?? undefined}
                    alt={currentVersion.label}
                    fit="contain"
                    mah={560}
                    style={{ display: "block" }}
                  />
                </Box>
              )}

              {/* Info bar below preview — theme-aware strip (no fixed light gray) */}
              <Box
                style={{
                  borderTop: surfaces.previewMetaBorderTop,
                  background: surfaces.previewMetaBg,
                }}
              >
                <Group p="sm" justify="space-between" wrap="wrap">
                  <Group gap="xs">
                    <Badge size="sm" variant="light" color="blue">
                      v{currentVersion.version_number}
                    </Badge>
                    {currentVersion.storage_provider === "google_drive" && (
                      <Badge size="sm" variant="light" color="violet">
                        Drive
                      </Badge>
                    )}
                    <Text size="sm" fw={500} c="var(--mantine-color-text)">
                      {currentVersion.label}
                    </Text>
                  </Group>
                  <Group gap="xs">
                    <Text size="xs" c="dimmed">
                      {currentVersion.file_name}
                    </Text>
                    <Text size="xs" c="dimmed">
                      ·
                    </Text>
                    <Text size="xs" c="dimmed">
                      {(currentVersion.file_size / (1024 * 1024)).toFixed(1)}{" "}
                      MB
                    </Text>
                  </Group>
                </Group>
              </Box>
              {currentVersion.notes && (
                <Text size="sm" c="dimmed" px="sm" pb="sm">
                  {currentVersion.notes}
                </Text>
              )}
            </Paper>
          ) : (
            <Paper
              p="xl"
              radius="lg"
              style={{
                borderStyle: "dashed",
                borderWidth: 2,
                borderColor: "var(--mantine-color-gray-3)",
                textAlign: "center",
              }}
              withBorder
            >
              <Center mih={300}>
                <Stack align="center" gap="md">
                  <ThemeIcon
                    variant="light"
                    color="gray"
                    size={64}
                    radius="xl"
                  >
                    {isVideo ? (
                      <IconVideo size={32} />
                    ) : (
                      <IconPhoto size={32} />
                    )}
                  </ThemeIcon>
                  <Text fw={500}>No versions yet</Text>
                  <Text size="sm" c="dimmed">
                    Upload a file to create the first version of this asset.
                  </Text>
                  <Button
                    variant="light"
                    radius="md"
                    leftSection={<IconPlus size={16} />}
                    onClick={() => setAddVersionOpen(true)}
                  >
                    Add first version
                  </Button>
                </Stack>
              </Center>
            </Paper>
          )}
        </Grid.Col>

        {/* Version sidebar */}
        <Grid.Col span={{ base: 12, md: 4 }}>
          <Paper p="md" radius="lg" withBorder>
            <Stack gap="md">
              <Group justify="space-between">
                <Text fw={700} size="sm" tt="uppercase" c="dimmed" style={{ letterSpacing: "0.05em" }}>
                  Version History
                </Text>
                <Badge size="sm" variant="light" color="gray" radius="sm">
                  {versions.length}
                </Badge>
              </Group>

              {versions.length === 0 ? (
                <Center mih={120}>
                  <Text size="sm" c="dimmed">
                    No versions uploaded.
                  </Text>
                </Center>
              ) : (
                <VersionTimeline
                  versions={versions}
                  currentVersionId={asset.current_version_id}
                  onSelect={handleVersionSelect}
                  onSetCurrent={(versionId) =>
                    setCurrent.mutate({ assetId: id, versionId })
                  }
                  onDelete={(versionId) =>
                    deleteVer.mutate({ assetId: id, versionId })
                  }
                />
              )}
            </Stack>
          </Paper>
        </Grid.Col>
      </Grid>

      {asset.drive_file_id && (
        <DriveRevisionsPanel assetId={id} driveFileId={asset.drive_file_id} />
      )}

      <AddVersionModal
        opened={addVersionOpen}
        onClose={() => setAddVersionOpen(false)}
        assetId={id}
      />

    </Stack>
  );
}
