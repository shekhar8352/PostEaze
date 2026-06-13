import { Button, Group, Loader, Paper, Stack, Table, Text, Title } from '@mantine/core';
import {
  useGoogleDriveRevisions,
  useImportDriveRevision,
} from '@/features/integrations/hooks/useGoogleDriveQueries';

function formatBytes(n: number): string {
  if (n < 1024) return `${n} B`;
  if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} KB`;
  return `${(n / (1024 * 1024)).toFixed(1)} MB`;
}

interface DriveRevisionsPanelProps {
  assetId: number;
  driveFileId: string;
}

export function DriveRevisionsPanel({ assetId, driveFileId }: DriveRevisionsPanelProps) {
  const { data, isLoading } = useGoogleDriveRevisions(driveFileId);
  const importRevision = useImportDriveRevision();

  return (
    <Paper p="md" radius="lg" withBorder>
      <Stack gap="sm">
        <Title order={5}>Drive revisions</Title>
        <Text size="sm" c="dimmed">
          Import a native Google Drive revision as a new version in this asset.
        </Text>
        {isLoading ? (
          <Group justify="center" py="md">
            <Loader size="sm" />
          </Group>
        ) : !data?.revisions?.length ? (
          <Text size="sm" c="dimmed">
            No revisions found for this file.
          </Text>
        ) : (
          <Table striped>
            <Table.Thead>
              <Table.Tr>
                <Table.Th>Modified</Table.Th>
                <Table.Th>Size</Table.Th>
                <Table.Th />
              </Table.Tr>
            </Table.Thead>
            <Table.Tbody>
              {data.revisions.map((rev) => (
                <Table.Tr key={rev.id}>
                  <Table.Td>
                    <Text size="sm">
                      {rev.modified_time
                        ? new Date(rev.modified_time).toLocaleString()
                        : rev.id}
                    </Text>
                  </Table.Td>
                  <Table.Td>
                    <Text size="sm" c="dimmed">
                      {formatBytes(rev.size)}
                    </Text>
                  </Table.Td>
                  <Table.Td>
                    <Button
                      size="xs"
                      variant="light"
                      loading={importRevision.isPending}
                      onClick={() =>
                        importRevision.mutate({
                          assetId,
                          revision_id: rev.id,
                          label: `Drive rev ${rev.id.slice(0, 8)}`,
                        })
                      }
                    >
                      Import as version
                    </Button>
                  </Table.Td>
                </Table.Tr>
              ))}
            </Table.Tbody>
          </Table>
        )}
      </Stack>
    </Paper>
  );
}
