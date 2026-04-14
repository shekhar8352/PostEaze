import { useState } from "react";
import {
  Button,
  Group,
  Modal,
  Stack,
  Text,
  TextInput,
  Textarea,
  rem,
} from "@mantine/core";
import { Dropzone, IMAGE_MIME_TYPE } from "@mantine/dropzone";
import { IconUpload, IconPhoto, IconX } from "@tabler/icons-react";
import { useAddVersion } from "../hooks/useMediaQueries";

const VIDEO_MIME_TYPES = ["video/mp4", "video/quicktime", "video/webm"];
const ALL_MIME_TYPES = [...IMAGE_MIME_TYPE, ...VIDEO_MIME_TYPES];
const MAX_SIZE = 50 * 1024 * 1024;

interface AddVersionModalProps {
  opened: boolean;
  onClose: () => void;
  assetId: number;
}

export function AddVersionModal({
  opened,
  onClose,
  assetId,
}: AddVersionModalProps) {
  const [file, setFile] = useState<File | null>(null);
  const [label, setLabel] = useState("");
  const [notes, setNotes] = useState("");
  const addVersion = useAddVersion();

  const reset = () => {
    setFile(null);
    setLabel("");
    setNotes("");
  };

  const handleClose = () => {
    reset();
    onClose();
  };

  const handleSubmit = () => {
    if (!file) return;
    const fd = new FormData();
    fd.append("file", file);
    if (label) fd.append("label", label);
    if (notes) fd.append("notes", notes);

    addVersion.mutate(
      { assetId, formData: fd },
      { onSuccess: handleClose }
    );
  };

  return (
    <Modal
      opened={opened}
      onClose={handleClose}
      title="Add New Version"
      size="lg"
    >
      <Stack gap="md">
        {!file ? (
          <Dropzone
            onDrop={(files) => setFile(files[0] ?? null)}
            maxSize={MAX_SIZE}
            accept={ALL_MIME_TYPES}
            multiple={false}
          >
            <Group
              justify="center"
              gap="xl"
              mih={140}
              style={{ pointerEvents: "none" }}
            >
              <Dropzone.Accept>
                <IconUpload style={{ width: rem(42), height: rem(42) }} stroke={1.5} />
              </Dropzone.Accept>
              <Dropzone.Reject>
                <IconX style={{ width: rem(42), height: rem(42) }} stroke={1.5} color="red" />
              </Dropzone.Reject>
              <Dropzone.Idle>
                <IconPhoto style={{ width: rem(42), height: rem(42) }} stroke={1.5} opacity={0.4} />
              </Dropzone.Idle>
              <div>
                <Text size="lg" inline>
                  Drop the new version file here
                </Text>
                <Text size="sm" c="dimmed" inline mt={4}>
                  Max 50 MB
                </Text>
              </div>
            </Group>
          </Dropzone>
        ) : (
          <Group gap="xs">
            <Text size="sm" fw={500}>
              {file.name}
            </Text>
            <Text size="xs" c="dimmed">
              ({(file.size / (1024 * 1024)).toFixed(1)} MB)
            </Text>
            <Button variant="subtle" size="xs" color="red" onClick={() => setFile(null)}>
              Remove
            </Button>
          </Group>
        )}

        <TextInput
          label="Label"
          placeholder="e.g. edited, color_graded, final"
          value={label}
          onChange={(e) => setLabel(e.currentTarget.value)}
        />

        <Textarea
          label="Notes"
          placeholder="What changed in this version?"
          value={notes}
          onChange={(e) => setNotes(e.currentTarget.value)}
          minRows={2}
        />

        <Group justify="flex-end">
          <Button variant="default" onClick={handleClose}>
            Cancel
          </Button>
          <Button
            onClick={handleSubmit}
            loading={addVersion.isPending}
            disabled={!file}
          >
            Upload Version
          </Button>
        </Group>
      </Stack>
    </Modal>
  );
}
