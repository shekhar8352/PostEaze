export interface ApiSuccessEnvelope<T> {
  status: string;
  msg: string;
  data: T;
}

export interface MediaVersion {
  id: number;
  version_number: number;
  label: string;
  blob_url: string;
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
  current_version_id: number | null;
  versions?: MediaVersion[];
  created_at: string;
  updated_at: string;
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

export interface PublishPayload {
  channel_ids: number[];
  caption: string;
}
