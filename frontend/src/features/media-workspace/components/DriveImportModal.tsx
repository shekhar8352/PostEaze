import { useState } from 'react';
import {
  ActionIcon,
  Badge,
  Button,
  Group,
  Loader,
  Modal,
  Stack,
  Table,
  Text,
  TextInput,
} from '@mantine/core';
import { IconArrowLeft, IconFolder, IconPhoto, IconVideo } from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';
import type { DriveFileItem } from '@/features/integrations/services/googleDriveApi';
import {
  useGoogleDriveFiles,
  useGoogleDriveStatus,
  useImportFromGoogleDrive,
} from '@/features/integrations/hooks/useGoogleDriveQueries';
import { GoogleDriveConnectCard } from '@/features/integrations/components/GoogleDriveConnectCard';

const BLOB_MAX = 50 * 1024 * 1024;

function formatBytes(n: number): string {
  if (n < 1024) return `${n} B`;
  if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} KB`;
  if (n < 1024 * 1024 * 1024) return `${(n / (1024 * 1024)).toFixed(1)} MB`;
  return `${(n / (1024 * 1024 * 1024)).toFixed(2)} GB`;
}

function isImportable(file: DriveFileItem): boolean {
  if (file.is_folder) return false;
  return (
    file.mime_type.startsWith('image/') || file.mime_type.startsWith('video/')
  );
}

interface DriveImportModalProps {
  opened: boolean;
  onClose: () => void;
}

export function DriveImportModal({ opened, onClose }: DriveImportModalProps) {
  const navigate = useNavigate();
  const { data: driveStatus } = useGoogleDriveStatus();
  const [folderStack, setFolderStack] = useState<{ id: string; name: string }[]>([]);
  const [pageToken, setPageToken] = useState<string | undefined>();
  const [search, setSearch] = useState('');
  const [importingId, setImportingId] = useState<string | null>(null);

  const currentFolderId = folderStack.length
    ? folderStack[folderStack.length - 1].id
    : undefined;

  const { data, isLoading, isFetching } = useGoogleDriveFiles(
    currentFolderId,
    pageToken,
    search.trim() || undefined,
    opened && !!driveStatus?.connected
  );

  const importMutation = useImportFromGoogleDrive();

  const enterFolder = (file: DriveFileItem) => {
    if (!file.is_folder) return;
    setFolderStack((s) => [...s, { id: file.id, name: file.name }]);
    setPageToken(undefined);
  };

  const goUp = () => {
    setFolderStack((s) => s.slice(0, -1));
    setPageToken(undefined);
  };

  const handleImport = async (file: DriveFileItem) => {
    setImportingId(file.id);
    try {
      const result = (await importMutation.mutateAsync({
        file_id: file.id,
        title: file.name.replace(/\.[^/.]+$/, ''),
        label: 'Drive import',
      })) as { id?: number };
      onClose();
      if (result?.id) navigate(`/workspace/${result.id}`);
    } finally {
      setImportingId(null);
    }
  };

  return (
    <Modal
      opened={opened}
      onClose={onClose}
      title="Import from Google Drive"
      size="xl"
    >
      <Stack gap="md">
        <GoogleDriveConnectCard />

        {!driveStatus?.connected ? (
          <Text size="sm" c="dimmed">
            Connect Google Drive above to browse and import files.
          </Text>
        ) : (
          <>
            <Group gap="sm">
              {folderStack.length > 0 && (
                <ActionIcon variant="light" onClick={goUp} aria-label="Up">
                  <IconArrowLeft size={16} />
                </ActionIcon>
              )}
              <Text size="sm" c="dimmed">
                {folderStack.length
                  ? folderStack.map((f) => f.name).join(' / ')
                  : 'My Drive'}
              </Text>
              <TextInput
                placeholder="Search files…"
                value={search}
                onChange={(e) => {
                  setSearch(e.currentTarget.value);
                  setPageToken(undefined);
                }}
                style={{ flex: 1 }}
                size="xs"
              />
            </Group>

            {isLoading ? (
              <Group justify="center" py="xl">
                <Loader size="sm" />
              </Group>
            ) : (
              <Table striped highlightOnHover>
                <Table.Thead>
                  <Table.Tr>
                    <Table.Th>Name</Table.Th>
                    <Table.Th>Size</Table.Th>
                    <Table.Th />
                  </Table.Tr>
                </Table.Thead>
                <Table.Tbody>
                  {(data?.files ?? []).map((file) => (
                    <Table.Tr key={file.id}>
                      <Table.Td>
                        <Group gap="xs">
                          {file.is_folder ? (
                            <IconFolder size={16} />
                          ) : file.mime_type.startsWith('video/') ? (
                            <IconVideo size={16} />
                          ) : (
                            <IconPhoto size={16} />
                          )}
                          {file.is_folder ? (
                            <Button
                              variant="subtle"
                              size="compact-sm"
                              onClick={() => enterFolder(file)}
                            >
                              {file.name}
                            </Button>
                          ) : (
                            <Text size="sm">{file.name}</Text>
                          )}
                          {!file.is_folder &&
                            file.size > BLOB_MAX &&
                            file.mime_type.startsWith('video/') && (
                              <Badge size="xs" variant="light" color="violet">
                                Kept in Drive
                              </Badge>
                            )}
                        </Group>
                      </Table.Td>
                      <Table.Td>
                        <Text size="sm" c="dimmed">
                          {file.is_folder ? '—' : formatBytes(file.size)}
                        </Text>
                      </Table.Td>
                      <Table.Td>
                        {isImportable(file) && (
                          <Button
                            size="xs"
                            variant="light"
                            loading={importingId === file.id}
                            onClick={() => handleImport(file)}
                          >
                            Import
                          </Button>
                        )}
                      </Table.Td>
                    </Table.Tr>
                  ))}
                </Table.Tbody>
              </Table>
            )}

            {data?.next_page_token && (
              <Button
                variant="default"
                loading={isFetching}
                onClick={() => setPageToken(data.next_page_token)}
              >
                Load more
              </Button>
            )}
          </>
        )}
      </Stack>
    </Modal>
  );
}
