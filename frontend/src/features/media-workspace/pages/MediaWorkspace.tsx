import { useState } from "react";
import {
  Button,
  Group,
  Loader,
  SegmentedControl,
  SimpleGrid,
  Stack,
  Text,
  Title,
  Center,
  Pagination,
} from "@mantine/core";
import { IconPlus, IconPhotoVideo } from "@tabler/icons-react";
import { useNavigate } from "react-router-dom";
import { useMediaAssets, useDeleteMediaAsset } from "../hooks/useMediaQueries";
import { AssetCard } from "../components/AssetCard";
import { MediaUploader } from "../components/MediaUploader";

const PAGE_SIZE = 20;

export default function MediaWorkspace() {
  const [statusFilter, setStatusFilter] = useState("all");
  const [page, setPage] = useState(1);
  const [uploadOpen, setUploadOpen] = useState(false);
  const navigate = useNavigate();
  const deleteAsset = useDeleteMediaAsset();

  const status = statusFilter === "all" ? undefined : statusFilter;
  const offset = (page - 1) * PAGE_SIZE;

  const { data, isLoading } = useMediaAssets(status, PAGE_SIZE, offset);

  const totalPages = data ? Math.ceil(data.total / PAGE_SIZE) : 0;

  return (
    <Stack gap="lg" p="md">
      <Group justify="space-between">
        <Title order={2}>Media Workspace</Title>
        <Button
          leftSection={<IconPlus size={16} />}
          onClick={() => setUploadOpen(true)}
        >
          Upload Media
        </Button>
      </Group>

      <SegmentedControl
        value={statusFilter}
        onChange={(val) => {
          setStatusFilter(val);
          setPage(1);
        }}
        data={[
          { label: "All", value: "all" },
          { label: "Drafts", value: "draft" },
          { label: "Ready", value: "ready" },
          { label: "Published", value: "published" },
        ]}
      />

      {isLoading ? (
        <Center mih={300}>
          <Loader />
        </Center>
      ) : !data || data.assets.length === 0 ? (
        <Center mih={300}>
          <Stack align="center" gap="xs">
            <IconPhotoVideo size={64} opacity={0.2} />
            <Text c="dimmed">No media assets yet.</Text>
            <Button
              variant="light"
              onClick={() => setUploadOpen(true)}
            >
              Upload your first media
            </Button>
          </Stack>
        </Center>
      ) : (
        <>
          <SimpleGrid cols={{ base: 1, sm: 2, md: 3, lg: 4 }} spacing="md">
            {data.assets.map((asset) => (
              <AssetCard
                key={asset.id}
                asset={asset}
                onClick={(id) => navigate(`/workspace/${id}`)}
                onDelete={(id) => deleteAsset.mutate(id)}
              />
            ))}
          </SimpleGrid>

          {totalPages > 1 && (
            <Center>
              <Pagination
                value={page}
                onChange={setPage}
                total={totalPages}
              />
            </Center>
          )}
        </>
      )}

      <MediaUploader
        opened={uploadOpen}
        onClose={() => setUploadOpen(false)}
      />
    </Stack>
  );
}
