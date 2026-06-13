export interface ApiSuccessEnvelope<T> {
  status: string;
  msg: string;
  data: T;
}

export type StorageProvider = "blob" | "google_drive";

export interface MediaVersion {
  id: number;
  version_number: number;
  label: string;
  storage_provider?: StorageProvider;
  blob_url: string;
  stream_url?: string;
  drive_file_id?: string;
  drive_revision_id?: string;
  file_name: string;
  content_type: string;
  file_size: number;
  metadata: Record<string, unknown>;
  notes: string;
  created_at: string;
}

export interface MediaAsset {
  id: number;
  title: string;
  asset_type: "photo" | "video";
  status: "draft" | "ready" | "published";
  drive_file_id?: string;
  current_version_id: number | null;
  versions?: MediaVersion[];
  created_at: string;
  updated_at: string;
}

/** Public HTTPS URL for publishing or preview (blob or signed stream). */
export function versionMediaUrl(v: MediaVersion | null | undefined): string | null {
  if (!v) return null;
  const u = (v.stream_url?.trim() || v.blob_url?.trim()) ?? "";
  return u.startsWith("https://") ? u : null;
}

export function currentVersionForAsset(asset: MediaAsset): MediaVersion | null {
  const versions = asset.versions ?? [];
  return (
    versions.find((x) => x.id === asset.current_version_id) ??
    versions[versions.length - 1] ??
    null
  );
}

export interface MediaAssetListResponse {
  assets: MediaAsset[];
  total: number;
  limit: number;
  offset: number;
}

export interface UploadResponse {
  url: string;
  pathname: string;
  content_type: string;
  file_size: number;
}

export interface CreateMediaAssetPayload {
  file: File;
  title: string;
  asset_type: "photo" | "video";
  label?: string;
}

export interface AddVersionPayload {
  assetId: number;
  file: File;
  label?: string;
  notes?: string;
}

export interface UpdateMediaAssetPayload {
  title?: string;
  status?: "draft" | "ready";
}
