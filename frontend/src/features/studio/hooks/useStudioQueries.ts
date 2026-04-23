import {
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";
import { notifications } from "@mantine/notifications";
import { studioApi } from "../services/studioApi";
import type {
  CreateCommentPayload,
  CreatePhasePayload,
  CreatePiecePayload,
  LinkAssetPayload,
  LinkScheduledPostPayload,
  MovePiecePayload,
  UpdatePhasePayload,
  UpdatePiecePayload,
  UpdateStudioPayload,
} from "../types";

/**
 * Query key factory for the Studio feature. Kept centralized so
 * mutations can invalidate the exact slices of the cache they touch.
 */
export const studioKeys = {
  all: ["studio"] as const,
  default: () => [...studioKeys.all, "default"] as const,
  board: (studioId: number) =>
    [...studioKeys.all, "board", studioId] as const,
  phases: (studioId: number) =>
    [...studioKeys.all, "phases", studioId] as const,
  piece: (pieceId: number) =>
    [...studioKeys.all, "piece", pieceId] as const,
  activities: (pieceId: number) =>
    [...studioKeys.all, "activities", pieceId] as const,
  comments: (pieceId: number) =>
    [...studioKeys.all, "comments", pieceId] as const,
};

function notifyError(e: unknown) {
  const message = e instanceof Error ? e.message : "Request failed";
  notifications.show({ title: "Error", message, color: "red" });
}

// ---------------------------------------------------------------------------
// Studio
// ---------------------------------------------------------------------------

export function useDefaultStudio() {
  return useQuery({
    queryKey: studioKeys.default(),
    queryFn: () => studioApi.ensureDefault(),
  });
}

export function useStudioBoard(studioId: number | undefined) {
  return useQuery({
    queryKey: studioKeys.board(studioId ?? 0),
    queryFn: () => studioApi.getBoard(studioId as number),
    enabled: typeof studioId === "number" && studioId > 0,
  });
}

export function useUpdateStudio() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, body }: { id: number; body: UpdateStudioPayload }) =>
      studioApi.update(id, body),
    onSuccess: (_data, vars) => {
      void qc.invalidateQueries({ queryKey: studioKeys.default() });
      void qc.invalidateQueries({ queryKey: studioKeys.board(vars.id) });
    },
    onError: notifyError,
  });
}

// ---------------------------------------------------------------------------
// Phases
// ---------------------------------------------------------------------------

export function usePhases(studioId: number | undefined) {
  return useQuery({
    queryKey: studioKeys.phases(studioId ?? 0),
    queryFn: () => studioApi.listPhases(studioId as number),
    enabled: typeof studioId === "number" && studioId > 0,
  });
}

export function useCreatePhase(studioId: number) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: CreatePhasePayload) =>
      studioApi.createPhase(studioId, body),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: studioKeys.phases(studioId) });
      void qc.invalidateQueries({ queryKey: studioKeys.board(studioId) });
    },
    onError: notifyError,
  });
}

export function useReorderPhases(studioId: number) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (phaseIds: number[]) =>
      studioApi.reorderPhases(studioId, phaseIds),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: studioKeys.phases(studioId) });
      void qc.invalidateQueries({ queryKey: studioKeys.board(studioId) });
    },
    onError: notifyError,
  });
}

export function useUpdatePhase(studioId: number) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      phaseId,
      body,
    }: {
      phaseId: number;
      body: UpdatePhasePayload;
    }) => studioApi.updatePhase(phaseId, body),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: studioKeys.phases(studioId) });
      void qc.invalidateQueries({ queryKey: studioKeys.board(studioId) });
    },
    onError: notifyError,
  });
}

export function useDeletePhase(studioId: number) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (phaseId: number) => studioApi.deletePhase(phaseId),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: studioKeys.phases(studioId) });
      void qc.invalidateQueries({ queryKey: studioKeys.board(studioId) });
      notifications.show({
        title: "Phase deleted",
        message: "The phase has been removed.",
        color: "green",
      });
    },
    onError: notifyError,
  });
}

// ---------------------------------------------------------------------------
// Pieces
// ---------------------------------------------------------------------------

export function useCreatePiece(studioId: number) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: CreatePiecePayload) =>
      studioApi.createPiece(studioId, body),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: studioKeys.board(studioId) });
    },
    onError: notifyError,
  });
}

export function usePieceDetail(pieceId: number | undefined) {
  return useQuery({
    queryKey: studioKeys.piece(pieceId ?? 0),
    queryFn: () => studioApi.getPiece(pieceId as number),
    enabled: typeof pieceId === "number" && pieceId > 0,
  });
}

export function useUpdatePiece(studioId: number) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      pieceId,
      body,
    }: {
      pieceId: number;
      body: UpdatePiecePayload;
    }) => studioApi.updatePiece(pieceId, body),
    onSuccess: (_data, vars) => {
      void qc.invalidateQueries({ queryKey: studioKeys.board(studioId) });
      void qc.invalidateQueries({
        queryKey: studioKeys.piece(vars.pieceId),
      });
    },
    onError: notifyError,
  });
}

export function useMovePiece(studioId: number) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      pieceId,
      body,
    }: {
      pieceId: number;
      body: MovePiecePayload;
    }) => studioApi.movePiece(pieceId, body),
    onSuccess: (_data, vars) => {
      void qc.invalidateQueries({ queryKey: studioKeys.board(studioId) });
      void qc.invalidateQueries({
        queryKey: studioKeys.piece(vars.pieceId),
      });
    },
    onError: notifyError,
  });
}

export function useSetPieceStatus(studioId: number) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      pieceId,
      status,
    }: {
      pieceId: number;
      status: "active" | "archived";
    }) => studioApi.setPieceStatus(pieceId, status),
    onSuccess: (_data, vars) => {
      void qc.invalidateQueries({ queryKey: studioKeys.board(studioId) });
      void qc.invalidateQueries({
        queryKey: studioKeys.piece(vars.pieceId),
      });
    },
    onError: notifyError,
  });
}

export function useDeletePiece(studioId: number) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (pieceId: number) => studioApi.deletePiece(pieceId),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: studioKeys.board(studioId) });
      notifications.show({
        title: "Piece deleted",
        message: "The piece has been removed.",
        color: "green",
      });
    },
    onError: notifyError,
  });
}

// ---------------------------------------------------------------------------
// Piece assets / scheduled posts
// ---------------------------------------------------------------------------

export function useLinkAsset() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      pieceId,
      body,
    }: {
      pieceId: number;
      body: LinkAssetPayload;
    }) => studioApi.linkAsset(pieceId, body),
    onSuccess: (_data, vars) => {
      void qc.invalidateQueries({ queryKey: studioKeys.piece(vars.pieceId) });
    },
    onError: notifyError,
  });
}

export function useUnlinkAsset() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      pieceId,
      mediaAssetId,
      role,
    }: {
      pieceId: number;
      mediaAssetId: number;
      role?: string;
    }) => studioApi.unlinkAsset(pieceId, mediaAssetId, role),
    onSuccess: (_data, vars) => {
      void qc.invalidateQueries({ queryKey: studioKeys.piece(vars.pieceId) });
    },
    onError: notifyError,
  });
}

export function useLinkScheduledPost() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      pieceId,
      body,
    }: {
      pieceId: number;
      body: LinkScheduledPostPayload;
    }) => studioApi.linkScheduledPost(pieceId, body),
    onSuccess: (_data, vars) => {
      void qc.invalidateQueries({ queryKey: studioKeys.piece(vars.pieceId) });
    },
    onError: notifyError,
  });
}

export function useUnlinkScheduledPost() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      pieceId,
      scheduledPostId,
    }: {
      pieceId: number;
      scheduledPostId: number;
    }) => studioApi.unlinkScheduledPost(pieceId, scheduledPostId),
    onSuccess: (_data, vars) => {
      void qc.invalidateQueries({ queryKey: studioKeys.piece(vars.pieceId) });
    },
    onError: notifyError,
  });
}

// ---------------------------------------------------------------------------
// Activities / Comments
// ---------------------------------------------------------------------------

export function usePieceActivities(
  pieceId: number | undefined,
  limit = 50
) {
  return useQuery({
    queryKey: studioKeys.activities(pieceId ?? 0),
    queryFn: () => studioApi.listActivities(pieceId as number, limit),
    enabled: typeof pieceId === "number" && pieceId > 0,
  });
}

export function usePieceComments(pieceId: number | undefined) {
  return useQuery({
    queryKey: studioKeys.comments(pieceId ?? 0),
    queryFn: () => studioApi.listComments(pieceId as number),
    enabled: typeof pieceId === "number" && pieceId > 0,
  });
}

export function useCreatePieceComment() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      pieceId,
      body,
    }: {
      pieceId: number;
      body: CreateCommentPayload;
    }) => studioApi.createComment(pieceId, body),
    onSuccess: (_data, vars) => {
      void qc.invalidateQueries({
        queryKey: studioKeys.comments(vars.pieceId),
      });
      void qc.invalidateQueries({
        queryKey: studioKeys.activities(vars.pieceId),
      });
    },
    onError: notifyError,
  });
}

export function useDeletePieceComment() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      pieceId,
      commentId,
    }: {
      pieceId: number;
      commentId: number;
    }) => studioApi.deleteComment(pieceId, commentId),
    onSuccess: (_data, vars) => {
      void qc.invalidateQueries({
        queryKey: studioKeys.comments(vars.pieceId),
      });
    },
    onError: notifyError,
  });
}
