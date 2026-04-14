import { useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import {
  ActionIcon,
  Badge,
  Box,
  Button,
  Center,
  Grid,
  Group,
  Image,
  Loader,
  Stack,
  Text,
  Title,
  Tooltip,
} from "@mantine/core";
import {
  IconArrowLeft,
  IconPlus,
  IconColumns,
  IconSend,
  IconTrash,
} from "@tabler/icons-react";
import {
  useMediaAsset,
  useDeleteMediaAsset,
  useSetCurrentVersion,
  useDeleteVersion,
} from "../hooks/useMediaQueries";
import { VersionTimeline } from "../components/VersionTimeline";
import { VersionCompare } from "../components/VersionCompare";
import { AddVersionModal } from "../components/AddVersionModal";
import { PublishDialog } from "../components/PublishDialog";
import type { MediaVersion } from "../types";

export default function MediaDetail() {
  const { assetId } = useParams<{ assetId: string }>();
  const navigate = useNavigate();
  const id = Number(assetId);

  const { data: asset, isLoading } = useMediaAsset(id);
  const deleteAsset = useDeleteMediaAsset();
  const setCurrent = useSetCurrentVersion();
  const deleteVer = useDeleteVersion();

  const [selectedVersion, setSelectedVersion] = useState<MediaVersion | null>(null);
  const [compareMode, setCompareMode] = useState(false);
  const [compareLeft, setCompareLeft] = useState<MediaVersion | null>(null);
  const [addVersionOpen, setAddVersionOpen] = useState(false);
  const [publishOpen, setPublishOpen] = useState(false);

  if (isLoading) {
    return (
      <Center mih={400}>
        <Loader />
      </Center>
    );
  }

  if (!asset) {
    return (
      <Center mih={400}>
        <Stack align="center">
          <Text c="dimmed">Asset not found</Text>
          <Button variant="light" onClick={() => navigate("/workspace")}>
            Back to Workspace
          </Button>
        </Stack>
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
      <Group justify="space-between">
        <Group>
          <ActionIcon variant="subtle" onClick={() => navigate("/workspace")}>
            <IconArrowLeft size={20} />
          </ActionIcon>
          <Title order={3}>{asset.title}</Title>
          <Badge>{asset.status}</Badge>
          <Badge variant="light">{asset.asset_type}</Badge>
        </Group>
        <Group>
          <Tooltip label={compareMode ? "Exit compare" : "Compare versions"}>
            <Button
              variant={compareMode ? "filled" : "light"}
              size="sm"
              leftSection={<IconColumns size={16} />}
              onClick={toggleCompare}
            >
              Compare
            </Button>
          </Tooltip>
          <Button
            variant="light"
            size="sm"
            leftSection={<IconPlus size={16} />}
            onClick={() => setAddVersionOpen(true)}
          >
            Add Version
          </Button>
          <Button
            size="sm"
            leftSection={<IconSend size={16} />}
            onClick={() => setPublishOpen(true)}
            disabled={asset.status === "published"}
          >
            Publish
          </Button>
          <Tooltip label="Delete asset">
            <ActionIcon
              variant="light"
              color="red"
              onClick={() => {
                deleteAsset.mutate(id, {
                  onSuccess: () => navigate("/workspace"),
                });
              }}
            >
              <IconTrash size={18} />
            </ActionIcon>
          </Tooltip>
        </Group>
      </Group>

      <Grid gutter="lg">
        <Grid.Col span={{ base: 12, md: 8 }}>
          {compareMode ? (
            <VersionCompare
              left={compareLeft}
              right={selectedVersion}
            />
          ) : currentVersion ? (
            <Box>
              {isVideo ? (
                <video
                  key={currentVersion.id}
                  src={currentVersion.blob_url}
                  controls
                  style={{
                    width: "100%",
                    maxHeight: 500,
                    borderRadius: 8,
                    background: "#000",
                  }}
                />
              ) : (
                <Image
                  key={currentVersion.id}
                  src={currentVersion.blob_url}
                  alt={currentVersion.label}
                  fit="contain"
                  mah={500}
                  radius="md"
                  style={{ background: "var(--mantine-color-gray-1)" }}
                />
              )}
              <Group mt="xs" gap="xs">
                <Text size="sm" fw={500}>
                  v{currentVersion.version_number} — {currentVersion.label}
                </Text>
                <Text size="xs" c="dimmed">
                  {currentVersion.file_name} &middot;{" "}
                  {(currentVersion.file_size / (1024 * 1024)).toFixed(1)} MB
                </Text>
              </Group>
              {currentVersion.notes && (
                <Text size="sm" c="dimmed" mt={4}>
                  {currentVersion.notes}
                </Text>
              )}
            </Box>
          ) : (
            <Center mih={300}>
              <Text c="dimmed">No versions yet. Upload one to get started.</Text>
            </Center>
          )}
        </Grid.Col>

        <Grid.Col span={{ base: 12, md: 4 }}>
          <Stack gap="md">
            <Group justify="space-between">
              <Text fw={600}>Versions</Text>
              <Text size="xs" c="dimmed">
                {versions.length} version{versions.length !== 1 ? "s" : ""}
              </Text>
            </Group>

            {versions.length === 0 ? (
              <Text size="sm" c="dimmed">
                No versions uploaded.
              </Text>
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
        </Grid.Col>
      </Grid>

      <AddVersionModal
        opened={addVersionOpen}
        onClose={() => setAddVersionOpen(false)}
        assetId={id}
      />

      <PublishDialog
        opened={publishOpen}
        onClose={() => setPublishOpen(false)}
        assetId={id}
      />
    </Stack>
  );
}
