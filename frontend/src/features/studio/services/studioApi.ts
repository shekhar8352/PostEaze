import apiClient from "@/services/api/client";
import type {
  ApiSuccessEnvelope,
  CreateCommentPayload,
  CreatePhasePayload,
  CreatePiecePayload,
  LinkAssetPayload,
  LinkScheduledPostPayload,
  MovePiecePayload,
  Phase,
  Piece,
  PieceActivity,
  PieceComment,
  PieceDetail,
  Studio,
  StudioBoard,
  UpdatePhasePayload,
  UpdatePiecePayload,
  UpdateStudioPayload,
} from "../types";

/**
 * Thin wrapper around the /v1/studios, /v1/phases, and /v1/pieces
 * endpoints. Responses are unwrapped from the ApiSuccessEnvelope so
 * callers get domain objects directly.
 */
export const studioApi = {
  // ----- Studio -----

  async ensureDefault(): Promise<Studio> {
    const { data } = await apiClient.get<ApiSuccessEnvelope<Studio>>(
      "/v1/studios/default"
    );
    return data.data;
  },

  async update(id: number, body: UpdateStudioPayload): Promise<Studio> {
    const { data } = await apiClient.put<ApiSuccessEnvelope<Studio>>(
      `/v1/studios/${id}`,
      body
    );
    return data.data;
  },

  async getBoard(id: number): Promise<StudioBoard> {
    const { data } = await apiClient.get<ApiSuccessEnvelope<StudioBoard>>(
      `/v1/studios/${id}/board`
    );
    return data.data;
  },

  // ----- Phases -----

  async listPhases(studioId: number): Promise<Phase[]> {
    const { data } = await apiClient.get<ApiSuccessEnvelope<Phase[]>>(
      `/v1/studios/${studioId}/phases`
    );
    return data.data;
  },

  async createPhase(studioId: number, body: CreatePhasePayload): Promise<Phase> {
    const { data } = await apiClient.post<ApiSuccessEnvelope<Phase>>(
      `/v1/studios/${studioId}/phases`,
      body
    );
    return data.data;
  },

  async reorderPhases(studioId: number, phaseIds: number[]): Promise<void> {
    await apiClient.post(`/v1/studios/${studioId}/phases/reorder`, {
      phase_ids: phaseIds,
    });
  },

  async updatePhase(phaseId: number, body: UpdatePhasePayload): Promise<Phase> {
    const { data } = await apiClient.put<ApiSuccessEnvelope<Phase>>(
      `/v1/phases/${phaseId}`,
      body
    );
    return data.data;
  },

  async deletePhase(phaseId: number): Promise<void> {
    await apiClient.delete(`/v1/phases/${phaseId}`);
  },

  // ----- Pieces -----

  async createPiece(studioId: number, body: CreatePiecePayload): Promise<Piece> {
    const { data } = await apiClient.post<ApiSuccessEnvelope<Piece>>(
      `/v1/studios/${studioId}/pieces`,
      body
    );
    return data.data;
  },

  async getPiece(pieceId: number): Promise<PieceDetail> {
    const { data } = await apiClient.get<ApiSuccessEnvelope<PieceDetail>>(
      `/v1/pieces/${pieceId}`
    );
    return data.data;
  },

  async updatePiece(pieceId: number, body: UpdatePiecePayload): Promise<Piece> {
    const { data } = await apiClient.put<ApiSuccessEnvelope<Piece>>(
      `/v1/pieces/${pieceId}`,
      body
    );
    return data.data;
  },

  async movePiece(pieceId: number, body: MovePiecePayload): Promise<Piece> {
    const { data } = await apiClient.post<ApiSuccessEnvelope<Piece>>(
      `/v1/pieces/${pieceId}/move`,
      body
    );
    return data.data;
  },

  async setPieceStatus(
    pieceId: number,
    status: "active" | "archived"
  ): Promise<void> {
    await apiClient.put(`/v1/pieces/${pieceId}/status`, { status });
  },

  async deletePiece(pieceId: number): Promise<void> {
    await apiClient.delete(`/v1/pieces/${pieceId}`);
  },

  // ----- Piece assets / scheduled posts -----

  async linkAsset(pieceId: number, body: LinkAssetPayload): Promise<void> {
    await apiClient.post(`/v1/pieces/${pieceId}/assets`, body);
  },

  async unlinkAsset(
    pieceId: number,
    mediaAssetId: number,
    role?: string
  ): Promise<void> {
    const query = role ? `?role=${encodeURIComponent(role)}` : "";
    await apiClient.delete(
      `/v1/pieces/${pieceId}/assets/${mediaAssetId}${query}`
    );
  },

  async linkScheduledPost(
    pieceId: number,
    body: LinkScheduledPostPayload
  ): Promise<void> {
    await apiClient.post(`/v1/pieces/${pieceId}/scheduled-posts`, body);
  },

  async unlinkScheduledPost(
    pieceId: number,
    scheduledPostId: number
  ): Promise<void> {
    await apiClient.delete(
      `/v1/pieces/${pieceId}/scheduled-posts/${scheduledPostId}`
    );
  },

  // ----- Activities / Comments -----

  async listActivities(
    pieceId: number,
    limit = 50
  ): Promise<PieceActivity[]> {
    const { data } = await apiClient.get<ApiSuccessEnvelope<PieceActivity[]>>(
      `/v1/pieces/${pieceId}/activities?limit=${limit}`
    );
    return data.data;
  },

  async listComments(pieceId: number): Promise<PieceComment[]> {
    const { data } = await apiClient.get<ApiSuccessEnvelope<PieceComment[]>>(
      `/v1/pieces/${pieceId}/comments`
    );
    return data.data;
  },

  async createComment(
    pieceId: number,
    body: CreateCommentPayload
  ): Promise<PieceComment> {
    const { data } = await apiClient.post<ApiSuccessEnvelope<PieceComment>>(
      `/v1/pieces/${pieceId}/comments`,
      body
    );
    return data.data;
  },

  async deleteComment(pieceId: number, commentId: number): Promise<void> {
    await apiClient.delete(`/v1/pieces/${pieceId}/comments/${commentId}`);
  },
};
