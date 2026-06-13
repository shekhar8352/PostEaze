import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { notifications } from '@mantine/notifications';
import { googleDriveApi } from '../services/googleDriveApi';
import { mediaKeys } from '@/features/media-workspace/hooks/useMediaQueries';

export const googleDriveKeys = {
  all: ['google-drive'] as const,
  status: () => [...googleDriveKeys.all, 'status'] as const,
  files: (folderId?: string, pageToken?: string, q?: string) =>
    [...googleDriveKeys.all, 'files', folderId ?? 'root', pageToken ?? '', q ?? ''] as const,
  revisions: (fileId: string) => [...googleDriveKeys.all, 'revisions', fileId] as const,
};

export function useGoogleDriveStatus() {
  return useQuery({
    queryKey: googleDriveKeys.status(),
    queryFn: () => googleDriveApi.status(),
  });
}

export function useConnectGoogleDrive() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (code: string) => googleDriveApi.connect(code),
    onSuccess: (data) => {
      void qc.invalidateQueries({ queryKey: googleDriveKeys.status() });
      notifications.show({
        title: 'Google Drive connected',
        message: data.email ? `Signed in as ${data.email}` : 'Drive is ready for imports',
        color: 'green',
      });
    },
    onError: (e: Error) => {
      notifications.show({ title: 'Connection failed', message: e.message, color: 'red' });
    },
  });
}

export function useDisconnectGoogleDrive() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: () => googleDriveApi.disconnect(),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: googleDriveKeys.all });
      notifications.show({
        title: 'Disconnected',
        message: 'Google Drive integration removed',
        color: 'green',
      });
    },
    onError: (e: Error) => {
      notifications.show({ title: 'Error', message: e.message, color: 'red' });
    },
  });
}

export function useGoogleDriveFiles(
  folderId?: string,
  pageToken?: string,
  q?: string,
  enabled = true
) {
  return useQuery({
    queryKey: googleDriveKeys.files(folderId, pageToken, q),
    queryFn: () => googleDriveApi.listFiles({ folderId, pageToken, q }),
    enabled,
  });
}

export function useGoogleDriveRevisions(fileId: string | null) {
  return useQuery({
    queryKey: googleDriveKeys.revisions(fileId ?? ''),
    queryFn: () => googleDriveApi.listRevisions(fileId!),
    enabled: !!fileId,
  });
}

export function useImportFromGoogleDrive() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: googleDriveApi.importFile,
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: mediaKeys.all });
      notifications.show({
        title: 'Imported',
        message: 'Media imported from Google Drive',
        color: 'green',
      });
    },
    onError: (e: Error) => {
      notifications.show({ title: 'Import failed', message: e.message, color: 'red' });
    },
  });
}

export function useImportDriveRevision() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      assetId,
      revision_id,
      label,
    }: {
      assetId: number;
      revision_id: string;
      label?: string;
    }) => googleDriveApi.importRevision(assetId, { revision_id, label }),
    onSuccess: (_data, vars) => {
      void qc.invalidateQueries({ queryKey: mediaKeys.detail(vars.assetId) });
      void qc.invalidateQueries({ queryKey: mediaKeys.all });
      notifications.show({
        title: 'Revision imported',
        message: 'New version created from Drive revision',
        color: 'green',
      });
    },
    onError: (e: Error) => {
      notifications.show({ title: 'Import failed', message: e.message, color: 'red' });
    },
  });
}
