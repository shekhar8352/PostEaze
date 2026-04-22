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
  Paper,
  Box,
  TextInput,
  Badge,
} from "@mantine/core";
import {
  IconPlus,
  IconUpload,
  IconSearch,
  IconCloudUpload,
} from "@tabler/icons-react";
import { useNavigate } from "react-router-dom";
import { notifications } from "@mantine/notifications";
import { useMediaAssets, useDeleteMediaAsset } from "../hooks/useMediaQueries";
import { AssetCard } from "../components/AssetCard";
import { MediaUploader } from "../components/MediaUploader";
import { useMediaWorkspaceSurfaces } from "../hooks/useMediaWorkspaceSurfaces";
import { LinkToPieceModal } from "@/features/studio/components/LinkToPieceModal";
import { useLinkAsset } from "@/features/studio/hooks/useStudioQueries";
import type { MediaAsset } from "../types";

const PAGE_SIZE = 20;

export default function MediaWorkspace() {
  const [statusFilter, setStatusFilter] = useState("all");
  const [page, setPage] = useState(1);
  const [uploadOpen, setUploadOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [linkAssetTarget, setLinkAssetTarget] = useState<MediaAsset | null>(
    null
  );
  const navigate = useNavigate();
  const deleteAsset = useDeleteMediaAsset();
  const linkAsset = useLinkAsset();
  const surfaces = useMediaWorkspaceSurfaces();

  const status = statusFilter === "all" ? undefined : statusFilter;
  const offset = (page - 1) * PAGE_SIZE;

  const { data, isLoading } = useMediaAssets(status, PAGE_SIZE, offset);

  const totalPages = data ? Math.ceil(data.total / PAGE_SIZE) : 0;

  const filteredAssets = data?.assets.filter((a) =>
    searchQuery
      ? a.title.toLowerCase().includes(searchQuery.toLowerCase())
      : true
  );

  return (
    <Stack gap="xl" p="md">
      {/* Page header */}
      <Group justify="space-between" align="flex-end">
        <Box>
          <Title
            order={2}
            style={{ letterSpacing: "-0.02em" }}
          >
            Media Workspace
          </Title>
          <Text size="sm" c="dimmed" mt={4}>
            Upload and version assets; schedule or publish from the Calendar
          </Text>
        </Box>
        <Button
          size="md"
          radius="md"
          leftSection={<IconUpload size={18} />}
          onClick={() => setUploadOpen(true)}
          variant="gradient"
          gradient={{ from: "blue", to: "cyan", deg: 135 }}
          style={{
            boxShadow: "0 4px 12px rgba(34,139,230,0.25)",
            transition: "transform 150ms ease, box-shadow 150ms ease",
          }}
          onMouseEnter={(e) => {
            e.currentTarget.style.transform = "translateY(-1px)";
            e.currentTarget.style.boxShadow =
              "0 6px 16px rgba(34,139,230,0.35)";
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.transform = "translateY(0)";
            e.currentTarget.style.boxShadow =
              "0 4px 12px rgba(34,139,230,0.25)";
          }}
        >
          Upload Media
        </Button>
      </Group>

      {/* Filter + Search bar */}
      <Paper p="sm" radius="lg" withBorder>
        <Group justify="space-between" wrap="wrap">
          <SegmentedControl
            value={statusFilter}
            onChange={(val) => {
              setStatusFilter(val);
              setPage(1);
            }}
            radius="md"
            data={[
              { label: "All", value: "all" },
              { label: "Drafts", value: "draft" },
              { label: "Ready", value: "ready" },
              { label: "Published", value: "published" },
            ]}
          />
          <TextInput
            placeholder="Search assets…"
            leftSection={<IconSearch size={16} />}
            radius="md"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.currentTarget.value)}
            style={{ minWidth: 240 }}
          />
        </Group>
      </Paper>

      {/* Summary badges */}
      {data && data.total > 0 && (
        <Group gap="sm">
          <Badge variant="light" color="gray" size="lg" radius="md">
            {data.total} asset{data.total !== 1 ? "s" : ""}
          </Badge>
          {statusFilter !== "all" && (
            <Badge variant="dot" color="blue" size="lg" radius="md">
              Filtered: {statusFilter}
            </Badge>
          )}
        </Group>
      )}

      {/* Content */}
      {isLoading ? (
        <Center mih={400}>
          <Stack align="center" gap="md">
            <Loader size="lg" type="dots" />
            <Text size="sm" c="dimmed">
              Loading your media…
            </Text>
          </Stack>
        </Center>
      ) : !filteredAssets || filteredAssets.length === 0 ? (
        <Center mih={400}>
          <Paper
            p="xl"
            radius="xl"
            withBorder
            style={{
              borderStyle: "dashed",
              borderWidth: 2,
              borderColor: "var(--mantine-color-gray-3)",
              textAlign: "center",
              maxWidth: 480,
            }}
          >
            <Stack align="center" gap="lg">
              <Box
                style={{
                  width: 100,
                  height: 100,
                  borderRadius: "50%",
                  ...surfaces.emptyStateIconOrb,
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                }}
              >
                <IconCloudUpload
                  size={48}
                  style={{
                    color: surfaces.isDark
                      ? "var(--mantine-color-blue-4)"
                      : "var(--mantine-color-blue-6)",
                  }}
                  stroke={1.5}
                />
              </Box>
              <Stack gap={4} align="center">
                <Text fw={600} size="lg">
                  {searchQuery
                    ? "No matching assets"
                    : "Your workspace is empty"}
                </Text>
                <Text size="sm" c="dimmed" maw={320}>
                  {searchQuery
                    ? `No assets match "${searchQuery}". Try a different search.`
                    : "Upload photos and videos to manage versions and compare edits. Use Calendar to schedule or publish."}
                </Text>
              </Stack>
              {!searchQuery && (
                <Button
                  variant="light"
                  size="md"
                  radius="md"
                  leftSection={<IconPlus size={18} />}
                  onClick={() => setUploadOpen(true)}
                >
                  Upload your first media
                </Button>
              )}
            </Stack>
          </Paper>
        </Center>
      ) : (
        <>
          <SimpleGrid
            cols={{ base: 1, xs: 2, sm: 2, md: 3, lg: 4 }}
            spacing="lg"
          >
            {filteredAssets.map((asset) => (
              <AssetCard
                key={asset.id}
                asset={asset}
                onClick={(id) => navigate(`/workspace/${id}`)}
                onDelete={(id) => deleteAsset.mutate(id)}
                onLinkToPiece={(target) => setLinkAssetTarget(target)}
              />
            ))}
          </SimpleGrid>

          {totalPages > 1 && (
            <Center mt="lg">
              <Pagination
                value={page}
                onChange={setPage}
                total={totalPages}
                radius="md"
                withEdges
              />
            </Center>
          )}
        </>
      )}

      <MediaUploader
        opened={uploadOpen}
        onClose={() => setUploadOpen(false)}
      />

      <LinkToPieceModal
        opened={linkAssetTarget !== null}
        onClose={() => setLinkAssetTarget(null)}
        title="Link asset to a Piece"
        description={
          linkAssetTarget
            ? `Pick a Piece to attach "${linkAssetTarget.title}" to. You can change or remove the link later from the Piece detail view.`
            : undefined
        }
        ctaLabel="Link asset"
        isSubmitting={linkAsset.isPending}
        onSelect={async (piece) => {
          if (!linkAssetTarget) return;
          try {
            await linkAsset.mutateAsync({
              pieceId: piece.id,
              body: { media_asset_id: linkAssetTarget.id },
            });
            notifications.show({
              color: "green",
              title: "Linked to Piece",
              message: `"${linkAssetTarget.title}" is now attached to ${piece.title}.`,
            });
            setLinkAssetTarget(null);
          } catch {
            // useLinkAsset surfaces its own error notification
          }
        }}
      />
    </Stack>
  );
}
