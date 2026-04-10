import { useEffect, useMemo, useState } from "react";
import dayjs from "dayjs";
import {
  Alert,
  Badge,
  Button,
  Group,
  Modal,
  MultiSelect,
  SegmentedControl,
  SimpleGrid,
  Stack,
  Text,
  Textarea,
  TextInput,
  Title,
} from "@mantine/core";
import { useReducedMotion } from "@mantine/hooks";
import { DateTimePicker } from "@mantine/dates";
import {
  IconBrandInstagram,
  IconCircleCheck,
  IconClock,
  IconPhotoPlus,
  IconUsers,
} from "@tabler/icons-react";
import type { BaseChannelDisplay } from "@/features/channels/types/base.types";
import { useCreateScheduledPost } from "../hooks/useScheduledPostsQueries";
import type { PostType } from "../types";
import styles from "./SchedulePostModal.module.css";

const POST_TYPES: { value: PostType; label: string }[] = [
  { value: "image", label: "Image" },
  { value: "video", label: "Video" },
  { value: "carousel", label: "Carousel" },
];

const STEP_DEF = [
  { label: "Timing", hint: "Now or schedule", Icon: IconClock },
  { label: "Channels", hint: "Where it goes live", Icon: IconUsers },
  { label: "Content", hint: "Media & caption", Icon: IconPhotoPlus },
  { label: "Review", hint: "Confirm & send", Icon: IconCircleCheck },
] as const;

/** Shared SegmentedControl styling — selected segment must pop on dark UI */
const MODAL_SEGMENTED_CLASS_NAMES = {
  root: styles.segmentedRoot,
  indicator: styles.segmentedIndicator,
  control: styles.segmentedControl,
  label: styles.segmentedLabel,
} as const;

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
  const d =
    initial ??
    (() => {
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
  const reduceMotion = useReducedMotion();
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
  const [contentError, setContentError] = useState<string | null>(null);

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
      setContentError(null);
    }
  }, [opened, initialStart]);

  useEffect(() => {
    setContentError(null);
  }, [postType, imageUrl, videoUrl, carouselUrls]);

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

  const validationMessage = validateStep2();

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
            platforms: ["instagram"],
            publish_now: true,
            post_type: postType,
            caption: caption.trim(),
            media: { items: buildMediaItems() },
          }
        : {
            channel_ids: channelIds,
            platforms: ["instagram"],
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

  const goToContentNext = () => {
    const e = validateStep2();
    setContentError(e);
    if (!e) setStep(3);
  };

  const modalTitle = (
    <div className={styles.titleBlock}>
      <div className={styles.kicker}>
        <IconBrandInstagram size={14} stroke={1.75} aria-hidden />
        Instagram
      </div>
      <Title order={3} className={styles.title}>
        New post
      </Title>
      <Text className={styles.subtitle}>Schedule or publish feed content in four short steps—timing, channels, media, review.</Text>
    </div>
  );

  return (
    <Modal
      opened={opened}
      onClose={onClose}
      title={modalTitle}
      size="xl"
      padding="lg"
      transitionProps={{ duration: reduceMotion ? 0 : 220, timingFunction: "cubic-bezier(0.4, 0, 0.2, 1)" }}
      classNames={{
        content: styles.modalContent,
        header: styles.modalHeader,
        body: styles.modalBody,
      }}
    >
      <nav className={styles.stepNav} aria-label="Create post steps">
        <ol className={styles.stepList}>
          {STEP_DEF.map((s, i) => {
            const Icon = s.Icon;
            const isCurrent = step === i;
            const isDone = step > i;
            return (
              <li key={s.label} aria-current={isCurrent ? "step" : undefined}>
                <div
                  className={`${styles.stepItem} ${isCurrent ? styles.stepItemCurrent : ""} ${isDone ? styles.stepItemDone : ""} ${!isCurrent && !isDone ? styles.stepItemUpcoming : ""}`}
                >
                  <span className={styles.stepIcon} aria-hidden>
                    <Icon size={18} stroke={1.75} />
                  </span>
                  <div className={styles.stepMeta}>
                    <div className={styles.stepLabel}>{s.label}</div>
                    <div className={styles.stepHint}>{s.hint}</div>
                  </div>
                </div>
              </li>
            );
          })}
        </ol>
      </nav>

      {step === 0 && (
        <div className={styles.panel}>
          <Text className={styles.sectionLabel}>When to publish</Text>
          <Stack gap="md">
            <SegmentedControl
              fullWidth
              color="blue"
              classNames={MODAL_SEGMENTED_CLASS_NAMES}
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
              <Text size="sm" className={styles.hint}>
                Instagram requires between 10 minutes and 75 days from now (UTC).
              </Text>
            ) : (
              <Text size="sm" className={styles.hint}>
                Publishes immediately after you confirm. The same media rules apply as for scheduled posts.
              </Text>
            )}
            <Group justify="flex-end" mt="md">
              <Button variant="default" onClick={onClose}>
                Cancel
              </Button>
              <Button onClick={() => setStep(1)} disabled={!canNextStep0}>
                Continue
              </Button>
            </Group>
          </Stack>
        </div>
      )}

      {step === 1 && (
        <div className={styles.panel}>
          <Text className={styles.sectionLabel}>Destination</Text>
          <Stack gap="md">
            {channelOptions.length === 0 ? (
              <div className={styles.emptyChannels}>
                <Text size="sm" fw={600} c="var(--pe-text)">
                  No Instagram channels yet
                </Text>
                <Text size="sm" mt={6} className={styles.hint}>
                  Connect an Instagram account in Channels, then return here to schedule.
                </Text>
              </div>
            ) : (
              <MultiSelect
                label="Channels"
                description="Pick every account that should receive this post"
                placeholder="Select one or more"
                data={channelOptions}
                value={channelValues}
                onChange={setChannelValues}
                searchable
                nothingFoundMessage="No channels"
              />
            )}
            <Group justify="space-between" className={styles.footerActions} wrap="nowrap">
              <Button variant="default" onClick={() => setStep(0)}>
                Back
              </Button>
              <Button onClick={() => setStep(2)} disabled={!canNextStep1}>
                Continue
              </Button>
            </Group>
          </Stack>
        </div>
      )}

      {step === 2 && (
        <div className={styles.panel}>
          <Text className={styles.sectionLabel}>Post content</Text>
          <Stack gap="md">
            <div>
              <Text size="sm" fw={600} mb={8} c="var(--pe-text)">
                Format
              </Text>
              <SegmentedControl
                fullWidth
                color="blue"
                classNames={MODAL_SEGMENTED_CLASS_NAMES}
                data={POST_TYPES.map((p) => ({ value: p.value, label: p.label }))}
                value={postType}
                onChange={(v) => setPostType(v as PostType)}
              />
            </div>
            <Textarea
              label="Caption"
              description="Optional — appears below your media on Instagram"
              placeholder="Write something your audience will want to engage with…"
              minRows={3}
              autosize
              minLength={0}
              maxLength={2200}
              value={caption}
              onChange={(e) => setCaption(e.target.value)}
            />
            {postType === "image" && (
              <Stack gap={6}>
                <TextInput
                  label="Image URL (HTTPS)"
                  placeholder="https://cdn.example.com/photo.jpg"
                  value={imageUrl}
                  onChange={(e) => setImageUrl(e.target.value)}
                />
                <Text className={styles.hint}>
                  Use a direct link Meta can fetch as JPEG—not Google Drive “view” pages.{" "}
                  <a href="https://developers.facebook.com/docs/instagram-platform/content-publishing/" target="_blank" rel="noreferrer">
                    Content publishing docs
                  </a>
                </Text>
              </Stack>
            )}
            {postType === "video" && (
              <Stack gap={6}>
                <TextInput
                  label="Video URL (HTTPS)"
                  placeholder="https://…"
                  value={videoUrl}
                  onChange={(e) => setVideoUrl(e.target.value)}
                />
                <Text className={styles.hint}>Direct HTTPS URL to the video file—not a player or share page.</Text>
              </Stack>
            )}
            {postType === "carousel" && (
              <Textarea
                label="Image URLs (one per line, 2–10, HTTPS)"
                placeholder={"https://example.com/a.jpg\nhttps://example.com/b.jpg"}
                minRows={4}
                value={carouselUrls}
                onChange={(e) => setCarouselUrls(e.target.value)}
              />
            )}
            {contentError && (
              <Alert color="red" variant="light" title="Fix media URLs">
                {contentError}
              </Alert>
            )}
            <Group justify="space-between" className={styles.footerActions} wrap="nowrap">
              <Button variant="default" onClick={() => setStep(1)}>
                Back
              </Button>
              <Button onClick={goToContentNext}>Continue</Button>
            </Group>
          </Stack>
        </div>
      )}

      {step === 3 && (
        <div className={styles.panel}>
          <Text className={styles.sectionLabel}>Summary</Text>
          <Stack gap="md">
            <SimpleGrid cols={{ base: 1, sm: 2 }} spacing="sm">
              <div className={styles.summaryCard}>
                <div className={styles.summaryKey}>When</div>
                <div className={styles.summaryValue}>
                  {timingMode === "now"
                    ? "Immediately"
                    : scheduledAtStr && dayjs(scheduledAtStr).isValid()
                      ? dayjs(scheduledAtStr).format("YYYY-MM-DD HH:mm")
                      : "—"}
                </div>
              </div>
              <div className={styles.summaryCard}>
                <div className={styles.summaryKey}>Format</div>
                <div className={styles.summaryValue}>
                  <Badge variant="light" color="blue" size="md" tt="none">
                    {POST_TYPES.find((p) => p.value === postType)?.label ?? postType}
                  </Badge>
                </div>
              </div>
              <div className={`${styles.summaryCard} ${styles.summaryCardWide}`}>
                <div className={styles.summaryKey}>Channels</div>
                <div className={styles.summaryValue}>
                  {channelValues.map((id) => channelOptions.find((o) => o.value === id)?.label ?? id).join(", ") || "—"}
                </div>
              </div>
              <div className={`${styles.summaryCard} ${styles.summaryCardWide}`}>
                <div className={styles.summaryKey}>Caption</div>
                <div className={styles.summaryValue}>{caption.trim() || "—"}</div>
              </div>
              <div className={`${styles.summaryCard} ${styles.summaryCardWide}`}>
                <div className={styles.summaryKey}>Media</div>
                <div className={styles.summaryValue}>
                  {postType === "carousel" ? `${buildMediaItems().length} images` : buildMediaItems()[0]?.url ?? "—"}
                </div>
              </div>
            </SimpleGrid>
            {validationMessage && (
              <Alert color="red" variant="light">
                {validationMessage}
              </Alert>
            )}
            <Group justify="space-between" className={styles.footerActions} wrap="nowrap">
              <Button variant="default" onClick={() => setStep(2)}>
                Back
              </Button>
              <Button loading={createMutation.isPending} onClick={handleSchedule} disabled={!!validationMessage}>
                {timingMode === "now" ? "Publish now" : "Schedule post"}
              </Button>
            </Group>
          </Stack>
        </div>
      )}
    </Modal>
  );
}
