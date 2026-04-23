import {
  Alert,
  Badge,
  Center,
  Divider,
  Drawer,
  Group,
  Loader,
  Stack,
  Tabs,
  Text,
} from "@mantine/core";
import { Icons } from "@/app/theme";
import type { Phase } from "../types";
import { usePieceDetail } from "../hooks/useStudioQueries";
import { PieceOverviewTab } from "./PieceOverviewTab";
import { PieceAssetsTab } from "./PieceAssetsTab";
import { PiecePublishTab } from "./PiecePublishTab";
import { PieceCommentsPanel } from "./PieceCommentsPanel";
import { PieceActivityFeed } from "./PieceActivityFeed";

interface PieceDetailDrawerProps {
  pieceId: number | null;
  studioId: number;
  phases: Phase[];
  opened: boolean;
  onClose: () => void;
}

/**
 * PieceDetailDrawer - Right-hand drawer showing a single piece with
 * Overview / Assets / Publish tabs and a shared comments + activity
 * feed below the tabs.
 */
export function PieceDetailDrawer({
  pieceId,
  studioId,
  phases,
  opened,
  onClose,
}: PieceDetailDrawerProps) {
  const { data, isLoading, error } = usePieceDetail(
    pieceId ?? undefined
  );

  const piece = data?.piece;
  const phase = piece ? phases.find((p) => p.id === piece.phase_id) : undefined;

  const title = piece ? (
    <Group gap="sm" align="center" wrap="nowrap">
      <Text fw={600} lineClamp={1} style={{ maxWidth: 420 }}>
        {piece.title}
      </Text>
      {phase && (
        <Badge
          variant="light"
          color="gray"
          styles={{
            root: { borderLeft: `3px solid ${phase.color || "#868e96"}` },
          }}
        >
          {phase.name}
        </Badge>
      )}
    </Group>
  ) : (
    "Loading…"
  );

  return (
    <Drawer
      opened={opened}
      onClose={onClose}
      title={title}
      position="right"
      size="xl"
      padding="md"
      keepMounted={false}
      overlayProps={{ backgroundOpacity: 0.3, blur: 1 }}
    >
      {isLoading && (
        <Center py="xl">
          <Loader size="sm" />
        </Center>
      )}

      {error && !isLoading && (
        <Alert
          color="red"
          icon={<Icons.AlertCircle size={16} />}
          title="Failed to load piece"
        >
          {error instanceof Error ? error.message : "Something went wrong."}
        </Alert>
      )}

      {data && piece && (
        <Stack gap="md">
          <Tabs defaultValue="overview" keepMounted={false}>
            <Tabs.List>
              <Tabs.Tab
                value="overview"
                leftSection={<Icons.FileText size={14} />}
              >
                Overview
              </Tabs.Tab>
              <Tabs.Tab
                value="assets"
                leftSection={<Icons.Photo size={14} />}
              >
                Assets
                {data.assets.length > 0 && (
                  <Badge size="xs" variant="light" ml="xs" color="gray">
                    {data.assets.length}
                  </Badge>
                )}
              </Tabs.Tab>
              <Tabs.Tab
                value="publish"
                leftSection={<Icons.Send size={14} />}
              >
                Publish
                {data.scheduled_posts.length > 0 && (
                  <Badge size="xs" variant="light" ml="xs" color="gray">
                    {data.scheduled_posts.length}
                  </Badge>
                )}
              </Tabs.Tab>
            </Tabs.List>

            <Tabs.Panel value="overview" pt="md">
              <PieceOverviewTab
                piece={piece}
                phase={phase}
                studioId={studioId}
                onClose={onClose}
              />
            </Tabs.Panel>

            <Tabs.Panel value="assets" pt="md">
              <PieceAssetsTab pieceId={piece.id} assets={data.assets} />
            </Tabs.Panel>

            <Tabs.Panel value="publish" pt="md">
              <PiecePublishTab
                pieceId={piece.id}
                links={data.scheduled_posts}
              />
            </Tabs.Panel>
          </Tabs>

          <Divider my="sm" label="Comments" labelPosition="left" />
          <PieceCommentsPanel pieceId={piece.id} />

          <Divider my="sm" label="Activity" labelPosition="left" />
          <PieceActivityFeed pieceId={piece.id} />
        </Stack>
      )}
    </Drawer>
  );
}
