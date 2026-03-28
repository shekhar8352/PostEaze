import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { notifications } from "@mantine/notifications";
import { scheduledPostService } from "../services/scheduledPostService";
import type { CreateScheduledPostRequest } from "../types";

export const scheduledPostKeys = {
  all: ["scheduled-posts"] as const,
  range: (from: string, to: string, channelId?: number) =>
    [...scheduledPostKeys.all, "range", from, to, channelId ?? "all"] as const,
};

export function useScheduledPostsRange(from: string, to: string, channelId?: number) {
  return useQuery({
    queryKey: scheduledPostKeys.range(from, to, channelId),
    queryFn: () => scheduledPostService.list(from, to, channelId),
  });
}

export function useCreateScheduledPost() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: CreateScheduledPostRequest) => scheduledPostService.create(body),
    onSuccess: ({ response, httpStatus }) => {
      void qc.invalidateQueries({ queryKey: scheduledPostKeys.all });
      if (response.overall_status === "scheduled") {
        notifications.show({ title: "Scheduled", message: "Post submitted to Instagram.", color: "green" });
      } else if (response.overall_status === "partial_failure" || httpStatus === 207) {
        notifications.show({
          title: "Partial success",
          message: "Some channels failed. Check results in the response.",
          color: "yellow",
        });
      } else {
        notifications.show({
          title: "Scheduling failed",
          message: "No channels succeeded. You can retry from the calendar.",
          color: "red",
        });
      }
    },
    onError: (e: Error) => {
      notifications.show({ title: "Error", message: e.message, color: "red" });
    },
  });
}
