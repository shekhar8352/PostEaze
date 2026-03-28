import { useEffect, useMemo, useState } from "react";
import dayjs from "dayjs";
import {
  Button,
  Group,
  Modal,
  MultiSelect,
  SegmentedControl,
  Stack,
  Stepper,
  Text,
  Textarea,
  TextInput,
} from "@mantine/core";
import { DateTimePicker } from "@mantine/dates";
import type { BaseChannelDisplay } from "@/features/channels/types/base.types";
import { useCreateScheduledPost } from "../hooks/useScheduledPostsQueries";
import type { PostType } from "../types";

const POST_TYPES: { value: PostType; label: string }[] = [
  { value: "image", label: "Image" },
  { value: "video", label: "Video" },
  { value: "carousel", label: "Carousel" },
];

function defaultSlotString(initial: Date | null): string {
  const d = initial ?? (() => {
    const x = new Date();
    x.setMinutes(0, 0, 0);
    x.setHours(x.getHours() + 1);
    return x;
  })();
  return dayjs(d).format("YYYY-MM-DD HH:mm");
}

export interface SchedulePostModalProps {
  opened: boolean;
  onClose: () => void;
  /** Prefill start when opening from a slot */
  initialStart: Date | null;
  channels: BaseChannelDisplay[];
}

export function SchedulePostModal({ opened, onClose, initialStart, channels }: SchedulePostModalProps) {
  const createMutation = useCreateScheduledPost();
  const [step, setStep] = useState(0);
  const [scheduledAtStr, setScheduledAtStr] = useState<string | null>(null);
  const [channelValues, setChannelValues] = useState<string[]>([]);
  const [postType, setPostType] = useState<PostType>("image");
  const [caption, setCaption] = useState("");
  const [imageUrl, setImageUrl] = useState("");
  const [videoUrl, setVideoUrl] = useState("");
  const [carouselUrls, setCarouselUrls] = useState("");

  useEffect(() => {
    if (opened) {
      setStep(0);
      setScheduledAtStr(defaultSlotString(initialStart));
      setChannelValues([]);
      setPostType("image");
      setCaption("");
      setImageUrl("");
      setVideoUrl("");
      setCarouselUrls("");
    }
  }, [opened, initialStart]);

  const igChannels = useMemo(
    () => channels.filter((c) => c.provider === "instagram"),
    [channels]
  );

  const channelOptions = useMemo(
    () =>
      igChannels.map((c) => ({
        value: String(c.channel_id),
        label: c.channelName || `Channel ${c.channel_id}`,
      })),
    [igChannels]
  );

  const canNextStep0 = scheduledAtStr != null && scheduledAtStr.length > 0;
  const canNextStep1 = channelValues.length > 0;

  const buildMediaItems = () => {
    if (postType === "image") {
      return [{ url: imageUrl.trim(), kind: "image" as const }];
    }
    if (postType === "video") {
      return [{ url: videoUrl.trim(), kind: "video" as const }];
    }
    const lines = carouselUrls
      .split("\n")
      .map((s) => s.trim())
      .filter(Boolean);
    return lines.map((url) => ({ url, kind: "image" as const }));
  };

  const validateStep2 = (): string | null => {
    const items = buildMediaItems();
    if (postType === "image") {
      if (items.length !== 1 || !items[0].url.startsWith("https://")) return "One HTTPS image URL required";
    }
    if (postType === "video") {
      if (items.length !== 1 || !items[0].url.startsWith("https://")) return "One HTTPS video URL required";
    }
    if (postType === "carousel") {
      if (items.length < 2 || items.length > 10) return "Enter 2–10 image URLs (one per line)";
      if (items.some((i) => !i.url.startsWith("https://"))) return "All carousel URLs must be HTTPS";
    }
    return null;
  };

  const handleSchedule = async () => {
    const err = validateStep2();
    if (err) {
      return;
    }
    const scheduledAt = dayjs(scheduledAtStr);
    if (!scheduledAt.isValid()) return;
    const channelIds = channelValues.map((v) => Number(v));
    const body = {
      channel_ids: channelIds,
      platforms: ["instagram"],
      scheduled_at: scheduledAt.toDate().toISOString(),
      post_type: postType,
      caption: caption.trim(),
      media: { items: buildMediaItems() },
    };
    try {
      await createMutation.mutateAsync(body);
      onClose();
    } catch {
      /* notification in hook */
    }
  };

  return (
    <Modal
      opened={opened}
      onClose={onClose}
      title="Schedule post"
      size="lg"
      transitionProps={{ duration: 200 }}
    >
      <Stepper active={step} onStepClick={setStep} allowNextStepsSelect={false}>
        <Stepper.Step label="When" description="Date & time">
          <Stack gap="md" mt="md">
            <DateTimePicker
              label="Scheduled time"
              placeholder="Pick date and time"
              value={scheduledAtStr}
              onChange={setScheduledAtStr}
              valueFormat="YYYY-MM-DD HH:mm"
              popoverProps={{ withinPortal: true }}
            />
            <Text size="xs" c="dimmed">
              Instagram requires between 10 minutes and 75 days from now (UTC).
            </Text>
            <Group justify="flex-end">
              <Button variant="default" onClick={onClose}>
                Cancel
              </Button>
              <Button onClick={() => setStep(1)} disabled={!canNextStep0}>
                Next
              </Button>
            </Group>
          </Stack>
        </Stepper.Step>

        <Stepper.Step label="Channels" description="Where to publish">
          <Stack gap="md" mt="md">
            <MultiSelect
              label="Channels"
              placeholder={channelOptions.length ? "Select one or more" : "Connect Instagram first"}
              data={channelOptions}
              value={channelValues}
              onChange={setChannelValues}
              searchable
              nothingFoundMessage="No channels"
            />
            <Group justify="space-between">
              <Button variant="default" onClick={() => setStep(0)}>
                Back
              </Button>
              <Button onClick={() => setStep(2)} disabled={!canNextStep1}>
                Next
              </Button>
            </Group>
          </Stack>
        </Stepper.Step>

        <Stepper.Step label="Content" description="Type & media URLs">
          <Stack gap="md" mt="md">
            <Text size="sm" fw={500}>
              Post type
            </Text>
            <SegmentedControl
              fullWidth
              data={POST_TYPES.map((p) => ({ value: p.value, label: p.label }))}
              value={postType}
              onChange={(v) => setPostType(v as PostType)}
            />
            <Textarea label="Caption" placeholder="Optional" minRows={2} value={caption} onChange={(e) => setCaption(e.target.value)} />
            {postType === "image" && (
              <TextInput
                label="Image URL (HTTPS)"
                placeholder="https://..."
                value={imageUrl}
                onChange={(e) => setImageUrl(e.target.value)}
              />
            )}
            {postType === "video" && (
              <TextInput
                label="Video URL (HTTPS)"
                placeholder="https://..."
                value={videoUrl}
                onChange={(e) => setVideoUrl(e.target.value)}
              />
            )}
            {postType === "carousel" && (
              <Textarea
                label="Image URLs (one per line, 2–10, HTTPS)"
                placeholder="https://example.com/a.jpg&#10;https://example.com/b.jpg"
                minRows={4}
                value={carouselUrls}
                onChange={(e) => setCarouselUrls(e.target.value)}
              />
            )}
            <Group justify="space-between">
              <Button variant="default" onClick={() => setStep(1)}>
                Back
              </Button>
              <Button
                onClick={() => {
                  const e = validateStep2();
                  if (e) {
                    return;
                  }
                  setStep(3);
                }}
              >
                Next
              </Button>
            </Group>
          </Stack>
        </Stepper.Step>

        <Stepper.Step label="Summary" description="Confirm">
          <Stack gap="sm" mt="md">
            <Text size="sm">
              <strong>When:</strong>{" "}
              {scheduledAtStr
                ? dayjs(scheduledAtStr).isValid()
                  ? dayjs(scheduledAtStr).format("YYYY-MM-DD HH:mm")
                  : "—"
                : "—"}
            </Text>
            <Text size="sm">
              <strong>Channels:</strong>{" "}
              {channelValues.map((id) => channelOptions.find((o) => o.value === id)?.label ?? id).join(", ") || "—"}
            </Text>
            <Text size="sm">
              <strong>Type:</strong> {postType}
            </Text>
            <Text size="sm">
              <strong>Caption:</strong> {caption.trim() || "—"}
            </Text>
            <Text size="sm">
              <strong>Media:</strong>{" "}
              {postType === "carousel" ? `${buildMediaItems().length} images` : buildMediaItems()[0]?.url ?? "—"}
            </Text>
            {validateStep2() && (
              <Text size="sm" c="red">
                {validateStep2()}
              </Text>
            )}
            <Group justify="space-between" mt="md">
              <Button variant="default" onClick={() => setStep(2)}>
                Back
              </Button>
              <Button loading={createMutation.isPending} onClick={handleSchedule} disabled={!!validateStep2()}>
                Schedule
              </Button>
            </Group>
          </Stack>
        </Stepper.Step>
      </Stepper>
    </Modal>
  );
}
