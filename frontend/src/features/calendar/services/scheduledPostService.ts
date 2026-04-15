import apiClient from "@/services/api/client";
import type {
  ApiSuccessEnvelope,
  CreateScheduledPostRequest,
  CreateScheduledPostResponse,
  ListScheduledPostsResponse,
  ScheduledPostListItem,
} from "../types";

function qs(params: Record<string, string | number | undefined>): string {
  const sp = new URLSearchParams();
  Object.entries(params).forEach(([k, v]) => {
    if (v !== undefined && v !== "") sp.set(k, String(v));
  });
  const s = sp.toString();
  return s ? `?${s}` : "";
}

export const scheduledPostService = {
  async list(from: string, to: string, channelId?: number): Promise<ListScheduledPostsResponse> {
    const { data } = await apiClient.get<ApiSuccessEnvelope<ListScheduledPostsResponse>>(
      `/v1/scheduled-posts${qs({ from, to, channel_id: channelId })}`
    );
    return data.data;
  },

  async create(body: CreateScheduledPostRequest): Promise<{ response: CreateScheduledPostResponse; httpStatus: number }> {
    const res = await apiClient.post<ApiSuccessEnvelope<CreateScheduledPostResponse>>(
      "/v1/scheduled-posts",
      body,
      {
        // Instagram "publish now" can take longer while container processing completes.
        timeout: 120_000,
        validateStatus: () => true,
      }
    );
    if (res.status >= 400) {
      const msg = (res.data as { msg?: string })?.msg ?? "Request failed";
      throw new Error(msg);
    }
    const envelope = res.data as ApiSuccessEnvelope<CreateScheduledPostResponse>;
    return { response: envelope.data, httpStatus: res.status };
  },

  async get(id: number): Promise<ScheduledPostListItem> {
    const { data } = await apiClient.get<ApiSuccessEnvelope<ScheduledPostListItem>>(`/v1/scheduled-posts/${id}`);
    return data.data;
  },

  async cancel(id: number): Promise<void> {
    await apiClient.delete(`/v1/scheduled-posts/${id}`);
  },
};
