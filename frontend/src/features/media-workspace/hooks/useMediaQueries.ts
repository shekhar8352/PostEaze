import {
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";
import { notifications } from "@mantine/notifications";
import { mediaApi } from "../services/mediaApi";
import type { UpdateMediaAssetPayload } from "../types";

export const mediaKeys = {
  all: ["media-assets"] as const,
  list: (status?: string, limit?: number, offset?: number) =>
    [...mediaKeys.all, "list", status ?? "all", limit, offset] as const,
  detail: (id: number) => [...mediaKeys.all, "detail", id] as const,
};

export function useMediaAssets(status?: string, limit = 20, offset = 0) {
  return useQuery({
    queryKey: mediaKeys.list(status, limit, offset),
    queryFn: () => mediaApi.list(status, limit, offset),
  });
}

export function useMediaAsset(id: number) {
  return useQuery({
    queryKey: mediaKeys.detail(id),
    queryFn: () => mediaApi.get(id),
    enabled: id > 0,
  });
}

export function useCreateMediaAsset() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (formData: FormData) => mediaApi.create(formData),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: mediaKeys.all });
      notifications.show({
        title: "Created",
        message: "Media asset created successfully.",
        color: "green",
      });
    },
    onError: (e: Error) => {
      notifications.show({ title: "Error", message: e.message, color: "red" });
    },
  });
}

export function useUpdateMediaAsset() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      id,
      ...body
    }: UpdateMediaAssetPayload & { id: number }) => mediaApi.update(id, body),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: mediaKeys.all });
    },
    onError: (e: Error) => {
      notifications.show({ title: "Error", message: e.message, color: "red" });
    },
  });
}

export function useDeleteMediaAsset() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => mediaApi.remove(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: mediaKeys.all });
      notifications.show({
        title: "Deleted",
        message: "Asset deleted.",
        color: "green",
      });
    },
    onError: (e: Error) => {
      notifications.show({ title: "Error", message: e.message, color: "red" });
    },
  });
}

export function useAddVersion() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ assetId, formData }: { assetId: number; formData: FormData }) =>
      mediaApi.addVersion(assetId, formData),
    onSuccess: (_data, vars) => {
      void qc.invalidateQueries({
        queryKey: mediaKeys.detail(vars.assetId),
      });
      notifications.show({
        title: "Version added",
        message: "New version uploaded.",
        color: "green",
      });
    },
    onError: (e: Error) => {
      notifications.show({ title: "Error", message: e.message, color: "red" });
    },
  });
}

export function useDeleteVersion() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      assetId,
      versionId,
    }: {
      assetId: number;
      versionId: number;
    }) => mediaApi.deleteVersion(assetId, versionId),
    onSuccess: (_data, vars) => {
      void qc.invalidateQueries({
        queryKey: mediaKeys.detail(vars.assetId),
      });
      notifications.show({
        title: "Version deleted",
        message: "Version removed.",
        color: "green",
      });
    },
    onError: (e: Error) => {
      notifications.show({ title: "Error", message: e.message, color: "red" });
    },
  });
}

export function useSetCurrentVersion() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      assetId,
      versionId,
    }: {
      assetId: number;
      versionId: number;
    }) => mediaApi.setCurrentVersion(assetId, versionId),
    onSuccess: (_data, vars) => {
      void qc.invalidateQueries({
        queryKey: mediaKeys.detail(vars.assetId),
      });
    },
    onError: (e: Error) => {
      notifications.show({ title: "Error", message: e.message, color: "red" });
    },
  });
}

