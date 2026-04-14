import { useState } from "react";
import {
  Button,
  Group,
  Modal,
  TextInput,
  Select,
  Stack,
  Text,
  rem,
} from "@mantine/core";
import { Dropzone, IMAGE_MIME_TYPE } from "@mantine/dropzone";
import { IconUpload, IconPhoto, IconX } from "@tabler/icons-react";
import { useCreateMediaAsset } from "../hooks/useMediaQueries";

const VIDEO_MIME_TYPES = ["video/mp4", "video/quicktime", "video/webm"];
const ALL_MIME_TYPES = [...IMAGE_MIME_TYPE, ...VIDEO_MIME_TYPES];
const MAX_SIZE = 50 * 1024 * 1024; // 50 MB

interface MediaUploaderProps {
  opened: boolean;
  onClose: () => void;
}

export function MediaUploader({ opened, onClose }: MediaUploaderProps) {
  const [file, setFile] = useState<File | null>(null);
  const [title, setTitle] = useState("");
  const [assetType, setAssetType] = useState<string | null>(null);
  const [label, setLabel] = useState("");
  const createAsset = useCreateMediaAsset();

  const reset = () => {
    setFile(null);
    setTitle("");
    setAssetType(null);
    setLabel("");
  };

  const handleClose = () => {
    reset();
    onClose();
  };

  const handleDrop = (files: File[]) => {
    const f = files[0];
    if (!f) return;
    setFile(f);
    if (!title) setTitle(f.name.replace(/\.[^.]+$/, ""));
    if (!assetType) {
      setAssetType(f.type.startsWith("video/") ? "video" : "photo");
    }
  };

  const handleSubmit = () => {
    if (!file || !title || !assetType) return;
    const fd = new FormData();
    fd.append("file", file);
    fd.append("title", title);
    fd.append("asset_type", assetType);
    if (label) fd.append("label", label);

    createAsset.mutate(fd, {
      onSuccess: handleClose,
    });
  };

  return (
    <Modal
      opened={opened}
      onClose={handleClose}
      title="Upload Media"
      size="lg"
    >
      <Stack gap="md">
        {!file ? (
          <Dropzone
            onDrop={handleDrop}
            maxSize={MAX_SIZE}
            accept={ALL_MIME_TYPES}
            multiple={false}
          >
            <Group
              justify="center"
              gap="xl"
              mih={180}
              style={{ pointerEvents: "none" }}
            >
              <Dropzone.Accept>
                <IconUpload
                  style={{ width: rem(52), height: rem(52) }}
                  stroke={1.5}
                />
              </Dropzone.Accept>
              <Dropzone.Reject>
                <IconX
                  style={{ width: rem(52), height: rem(52) }}
                  stroke={1.5}
                  color="red"
                />
              </Dropzone.Reject>
              <Dropzone.Idle>
                <IconPhoto
                  style={{ width: rem(52), height: rem(52) }}
                  stroke={1.5}
                  opacity={0.4}
                />
              </Dropzone.Idle>

              <div>
                <Text size="xl" inline>
                  Drag a photo or video here
                </Text>
                <Text size="sm" c="dimmed" inline mt={7}>
                  Max 50 MB. JPEG, PNG, WebP, MP4, MOV, WebM.
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
            <Button
              variant="subtle"
              size="xs"
              color="red"
              onClick={() => setFile(null)}
            >
              Remove
            </Button>
          </Group>
        )}

        <TextInput
          label="Title"
          placeholder="e.g. Summer Campaign Photo"
          value={title}
          onChange={(e) => setTitle(e.currentTarget.value)}
          required
        />

        <Select
          label="Type"
          data={[
            { value: "photo", label: "Photo" },
            { value: "video", label: "Video" },
          ]}
          value={assetType}
          onChange={setAssetType}
          required
        />

        <TextInput
          label="Version Label"
          placeholder="raw (default)"
          value={label}
          onChange={(e) => setLabel(e.currentTarget.value)}
        />

        <Group justify="flex-end">
          <Button variant="default" onClick={handleClose}>
            Cancel
          </Button>
          <Button
            onClick={handleSubmit}
            loading={createAsset.isPending}
            disabled={!file || !title || !assetType}
          >
            Upload & Create
          </Button>
        </Group>
      </Stack>
    </Modal>
  );
}
