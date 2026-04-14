import { useState } from "react";
import {
  Button,
  Group,
  Modal,
  TextInput,
  Select,
  Stack,
  Text,
  Box,
  Image,
  ThemeIcon,
  Paper,
  CloseButton,
  Badge,
  Progress,
} from "@mantine/core";
import { Dropzone, IMAGE_MIME_TYPE } from "@mantine/dropzone";
import {
  IconUpload,
  IconX,
  IconCloudUpload,
  IconFile,
  IconVideo,
  IconSparkles,
} from "@tabler/icons-react";
import { useCreateMediaAsset } from "../hooks/useMediaQueries";
import { useMediaWorkspaceSurfaces } from "../hooks/useMediaWorkspaceSurfaces";

const VIDEO_MIME_TYPES = ["video/mp4", "video/quicktime", "video/webm"];
const ALL_MIME_TYPES = [...IMAGE_MIME_TYPE, ...VIDEO_MIME_TYPES];
const MAX_SIZE = 50 * 1024 * 1024;

interface MediaUploaderProps {
  opened: boolean;
  onClose: () => void;
}

function formatFileSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(0)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

export function MediaUploader({ opened, onClose }: MediaUploaderProps) {
  const [file, setFile] = useState<File | null>(null);
  const [preview, setPreview] = useState<string | null>(null);
  const [title, setTitle] = useState("");
  const [assetType, setAssetType] = useState<string | null>(null);
  const [label, setLabel] = useState("");
  const createAsset = useCreateMediaAsset();
  const surfaces = useMediaWorkspaceSurfaces();

  const reset = () => {
    setFile(null);
    setPreview(null);
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
    if (f.type.startsWith("image/")) {
      const url = URL.createObjectURL(f);
      setPreview(url);
    }
  };

  const handleRemoveFile = () => {
    if (preview) URL.revokeObjectURL(preview);
    setFile(null);
    setPreview(null);
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

  const isVideo = file?.type.startsWith("video/");

  return (
    <Modal
      opened={opened}
      onClose={handleClose}
      title={
        <Group gap="xs">
          <ThemeIcon variant="light" color="blue" size="sm" radius="xl">
            <IconCloudUpload size={14} />
          </ThemeIcon>
          <Text fw={600}>Upload Media</Text>
        </Group>
      }
      size="lg"
      radius="lg"
      overlayProps={{ backgroundOpacity: 0.4, blur: 4 }}
    >
      <Stack gap="lg">
        {/* Dropzone or file preview */}
        {!file ? (
          <Dropzone
            onDrop={handleDrop}
            maxSize={MAX_SIZE}
            accept={ALL_MIME_TYPES}
            multiple={false}
            radius="lg"
            style={surfaces.dropzoneUpload}
          >
            <Stack
              align="center"
              justify="center"
              gap="md"
              mih={200}
              style={{ pointerEvents: "none" }}
            >
              <Dropzone.Accept>
                <ThemeIcon
                  variant="light"
                  color="blue"
                  size={72}
                  radius="xl"
                >
                  <IconUpload size={36} stroke={1.5} />
                </ThemeIcon>
              </Dropzone.Accept>
              <Dropzone.Reject>
                <ThemeIcon variant="light" color="red" size={72} radius="xl">
                  <IconX size={36} stroke={1.5} />
                </ThemeIcon>
              </Dropzone.Reject>
              <Dropzone.Idle>
                <ThemeIcon
                  variant="light"
                  color="blue"
                  size={72}
                  radius="xl"
                  style={{ opacity: 0.8 }}
                >
                  <IconCloudUpload size={36} stroke={1.5} />
                </ThemeIcon>
              </Dropzone.Idle>

              <Stack gap={4} align="center">
                <Text size="lg" fw={600} c="var(--mantine-color-text)">
                  Drag & drop your file here
                </Text>
                <Text size="sm" c="dimmed">
                  or{" "}
                  <Text span c="blue" fw={500}>
                    click to browse
                  </Text>
                </Text>
                <Group gap="xs" mt={4}>
                  {(["JPEG", "PNG", "WebP", "MP4", "MOV", "WebM"] as const).map(
                    (fmt) => (
                      <Badge
                        key={fmt}
                        size="xs"
                        variant="outline"
                        color="blue"
                      >
                        {fmt}
                      </Badge>
                    )
                  )}
                </Group>
                <Text size="xs" c="dimmed" mt={2}>
                  Max file size: 50 MB
                </Text>
              </Stack>
            </Stack>
          </Dropzone>
        ) : (
          <Paper p="md" radius="lg" withBorder>
            <Group wrap="nowrap" gap="md">
              {/* File thumbnail */}
              {preview ? (
                <Image
                  src={preview}
                  alt="Preview"
                  w={80}
                  h={80}
                  radius="md"
                  fit="cover"
                />
              ) : (
                <Box
                  style={{
                    width: 80,
                    height: 80,
                    borderRadius: 12,
                    background:
                      "linear-gradient(135deg, var(--mantine-color-violet-1), var(--mantine-color-violet-2))",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    flexShrink: 0,
                  }}
                >
                  {isVideo ? (
                    <IconVideo size={32} color="var(--mantine-color-violet-5)" />
                  ) : (
                    <IconFile size={32} color="var(--mantine-color-violet-5)" />
                  )}
                </Box>
              )}

              {/* File info */}
              <Stack gap={2} style={{ flex: 1, minWidth: 0 }}>
                <Text size="sm" fw={600} lineClamp={1}>
                  {file.name}
                </Text>
                <Group gap="xs">
                  <Badge size="xs" variant="light" color={isVideo ? "violet" : "blue"}>
                    {isVideo ? "Video" : "Image"}
                  </Badge>
                  <Text size="xs" c="dimmed">
                    {formatFileSize(file.size)}
                  </Text>
                </Group>
                {createAsset.isPending && (
                  <Progress
                    size="xs"
                    radius="xl"
                    animated
                    value={100}
                    mt={4}
                  />
                )}
              </Stack>

              <CloseButton
                onClick={handleRemoveFile}
                variant="subtle"
                color="gray"
                size="sm"
              />
            </Group>
          </Paper>
        )}

        {/* Form fields */}
        <Stack gap="sm">
          <TextInput
            label="Title"
            placeholder="e.g. Summer Campaign Photo"
            value={title}
            onChange={(e) => setTitle(e.currentTarget.value)}
            required
            radius="md"
          />

          <Group grow>
            <Select
              label="Type"
              data={[
                { value: "photo", label: "Photo" },
                { value: "video", label: "Video" },
              ]}
              value={assetType}
              onChange={setAssetType}
              required
              radius="md"
            />

            <TextInput
              label="Version Label"
              placeholder="raw (default)"
              value={label}
              onChange={(e) => setLabel(e.currentTarget.value)}
              radius="md"
            />
          </Group>
        </Stack>

        {/* Actions */}
        <Group justify="flex-end" mt="xs">
          <Button variant="subtle" color="gray" onClick={handleClose} radius="md">
            Cancel
          </Button>
          <Button
            onClick={handleSubmit}
            loading={createAsset.isPending}
            disabled={!file || !title || !assetType}
            radius="md"
            leftSection={<IconSparkles size={16} />}
            variant="gradient"
            gradient={{ from: "blue", to: "cyan", deg: 135 }}
          >
            Upload & Create
          </Button>
        </Group>
      </Stack>
    </Modal>
  );
}
