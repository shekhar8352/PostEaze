import { useState } from "react";
import {
  Button,
  Checkbox,
  Group,
  Modal,
  Stack,
  Text,
  Textarea,
} from "@mantine/core";
import { IconSend } from "@tabler/icons-react";
import { useChannels } from "@/features/channels/services/channelQueries";
import { usePublishMediaAsset } from "../hooks/useMediaQueries";

interface PublishDialogProps {
  opened: boolean;
  onClose: () => void;
  assetId: number;
}

export function PublishDialog({ opened, onClose, assetId }: PublishDialogProps) {
  const [selectedChannels, setSelectedChannels] = useState<number[]>([]);
  const [caption, setCaption] = useState("");
  const { data: channels, isLoading: channelsLoading } = useChannels();
  const publish = usePublishMediaAsset();

  const handleClose = () => {
    setSelectedChannels([]);
    setCaption("");
    onClose();
  };

  const toggleChannel = (channelId: number) => {
    setSelectedChannels((prev) =>
      prev.includes(channelId)
        ? prev.filter((id) => id !== channelId)
        : [...prev, channelId]
    );
  };

  const handlePublish = () => {
    if (selectedChannels.length === 0) return;
    publish.mutate(
      { assetId, channel_ids: selectedChannels, caption },
      { onSuccess: handleClose }
    );
  };

  return (
    <Modal
      opened={opened}
      onClose={handleClose}
      title="Publish to Channels"
      size="md"
    >
      <Stack gap="md">
        <Text size="sm" c="dimmed">
          The active version will be published to the selected channels.
        </Text>

        <Stack gap="xs">
          <Text fw={500} size="sm">
            Select Channels
          </Text>
          {channelsLoading ? (
            <Text size="sm" c="dimmed">Loading channels...</Text>
          ) : !channels || channels.length === 0 ? (
            <Text size="sm" c="dimmed">
              No channels connected. Connect a channel first.
            </Text>
          ) : (
            channels.map((ch) => (
              <Checkbox
                key={ch.channel_id}
                label={`${ch.channelName} (${ch.provider})`}
                checked={selectedChannels.includes(ch.channel_id)}
                onChange={() => toggleChannel(ch.channel_id)}
              />
            ))
          )}
        </Stack>

        <Textarea
          label="Caption"
          placeholder="Write a caption for your post..."
          value={caption}
          onChange={(e) => setCaption(e.currentTarget.value)}
          minRows={3}
        />

        <Group justify="flex-end">
          <Button variant="default" onClick={handleClose}>
            Cancel
          </Button>
          <Button
            leftSection={<IconSend size={16} />}
            onClick={handlePublish}
            loading={publish.isPending}
            disabled={selectedChannels.length === 0}
          >
            Publish Now
          </Button>
        </Group>
      </Stack>
    </Modal>
  );
}
