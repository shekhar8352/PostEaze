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

/** Instagram cURLs the URL; viewer pages (Drive, etc.) are HTML, not JPEG bytes. */
function instagramMediaUrlHint(url: string): string | null {
  try {
    const u = new URL(url);
    const h = u.hostname.toLowerCase();
    const p = u.pathname.toLowerCase();
    if (h.includes("drive.google.com") && (p.includes("/file/d/") || p.includes("/file/u/") || p.startsWith("/open"))) {
      return "This looks like a Google Drive share link. Instagram needs a direct HTTPS link that returns the image file (JPEG), not a preview page.";
    }
    if (h === "docs.google.com") {
      return "Google Docs links are not direct image URLs. Use a public URL whose response is the raw image.";
    }
    const isDropbox = h.endsWith(".dropbox.com") || h === "dropbox.com";
    if (isDropbox && p.includes("/s/") && !u.search.includes("raw=1") && !u.search.includes("dl=1")) {
      return "For Dropbox, use a direct link with ?raw=1 so Instagram receives file bytes.";
    }
  } catch {
    return null;
  }
  return null;
}

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

type TimingMode = "now" | "later";

export function SchedulePostModal({ opened, onClose, initialStart, channels }: SchedulePostModalProps) {
  const createMutation = useCreateScheduledPost();
  const [step, setStep] = useState(0);
  const [timingMode, setTimingMode] = useState<TimingMode>("later");
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
      setTimingMode("later");
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

  const canNextStep0 =
    timingMode === "now" || (scheduledAtStr != null && scheduledAtStr.length > 0 && dayjs(scheduledAtStr).isValid());
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
      const hint = instagramMediaUrlHint(items[0].url);
      if (hint) return hint;
    }
    if (postType === "video") {
      if (items.length !== 1 || !items[0].url.startsWith("https://")) return "One HTTPS video URL required";
      const hint = instagramMediaUrlHint(items[0].url);
      if (hint) return hint;
    }
    if (postType === "carousel") {
      if (items.length < 2 || items.length > 10) return "Enter 2–10 image URLs (one per line)";
      if (items.some((i) => !i.url.startsWith("https://"))) return "All carousel URLs must be HTTPS";
      for (const i of items) {
        const hint = instagramMediaUrlHint(i.url);
        if (hint) return hint;
      }
    }
    return null;
  };

  const handleSchedule = async () => {
    const err = validateStep2();
    if (err) {
      return;
    }
    if (timingMode === "later") {
      const scheduledAt = dayjs(scheduledAtStr);
      if (!scheduledAt.isValid()) return;
    }
    const channelIds = channelValues.map((v) => Number(v));
    const body =
      timingMode === "now"
        ? {
            channel_ids: channelIds,
            platforms: ["instagram"] as const,
            publish_now: true,
            post_type: postType,
            caption: caption.trim(),
            media: { items: buildMediaItems() },
          }
        : {
            channel_ids: channelIds,
            platforms: ["instagram"] as const,
            publish_now: false,
            scheduled_at: dayjs(scheduledAtStr!).toDate().toISOString(),
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
      title="Create post"
      size="lg"
      transitionProps={{ duration: 200 }}
    >
      <Stepper active={step} onStepClick={setStep} allowNextStepsSelect={false}>
        <Stepper.Step label="When" description="Now or later">
          <Stack gap="md" mt="md">
            <Text size="sm" fw={500}>
              When to publish
            </Text>
            <SegmentedControl
              fullWidth
              data={[
                { value: "now", label: "Post now" },
                { value: "later", label: "Schedule" },
              ]}
              value={timingMode}
              onChange={(v) => setTimingMode(v as TimingMode)}
            />
            {timingMode === "later" && (
              <DateTimePicker
                label="Scheduled time"
                placeholder="Pick date and time"
                value={scheduledAtStr}
                onChange={setScheduledAtStr}
                valueFormat="YYYY-MM-DD HH:mm"
                popoverProps={{ withinPortal: true }}
              />
            )}
            {timingMode === "later" ? (
              <Text size="xs" c="dimmed">
                Instagram requires between 10 minutes and 75 days from now (UTC).
              </Text>
            ) : (
              <Text size="xs" c="dimmed">
                Publishes immediately after you confirm (Instagram feed; same media rules apply).
              </Text>
            )}
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
              <Stack gap={4}>
                <TextInput
                  label="Image URL (HTTPS)"
                  placeholder="https://cdn.example.com/photo.jpg"
                  value={imageUrl}
                  onChange={(e) => setImageUrl(e.target.value)}
                />
                <Text size="xs" c="dimmed">
                  Must be a direct link Meta can fetch as JPEG (not Google Drive “view” pages). See{" "}
                  <a
                    href="https://developers.facebook.com/docs/instagram-platform/content-publishing/"
                    target="_blank"
                    rel="noreferrer"
                  >
                    Content publishing
                  </a>
                  .
                </Text>
              </Stack>
            )}
            {postType === "video" && (
              <Stack gap={4}>
                <TextInput
                  label="Video URL (HTTPS)"
                  placeholder="https://..."
                  value={videoUrl}
                  onChange={(e) => setVideoUrl(e.target.value)}
                />
                <Text size="xs" c="dimmed">
                  Same as images: a direct HTTPS URL to the video file, not a player or share page.
                </Text>
              </Stack>
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
              {timingMode === "now"
                ? "Immediately"
                : scheduledAtStr && dayjs(scheduledAtStr).isValid()
                  ? dayjs(scheduledAtStr).format("YYYY-MM-DD HH:mm")
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
                {timingMode === "now" ? "Publish now" : "Schedule"}
              </Button>
            </Group>
          </Stack>
        </Stepper.Step>
      </Stepper>
    </Modal>
  );
}
