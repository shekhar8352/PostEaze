import { useState } from "react";
import {
  Box,
  Button,
  CloseButton,
  Group,
  Image,
  Modal,
  Paper,
  Progress,
  Stack,
  Text,
  TextInput,
  Textarea,
  ThemeIcon,
} from "@mantine/core";
import { Dropzone, IMAGE_MIME_TYPE } from "@mantine/dropzone";
import {
  IconUpload,
  IconX,
  IconCloudUpload,
  IconFile,
  IconVideo,
  IconVersions,
} from "@tabler/icons-react";
import { useAddVersion } from "../hooks/useMediaQueries";
import { useMediaWorkspaceSurfaces } from "../hooks/useMediaWorkspaceSurfaces";

const VIDEO_MIME_TYPES = ["video/mp4", "video/quicktime", "video/webm"];
const ALL_MIME_TYPES = [...IMAGE_MIME_TYPE, ...VIDEO_MIME_TYPES];
const MAX_SIZE = 50 * 1024 * 1024;

interface AddVersionModalProps {
  opened: boolean;
  onClose: () => void;
  assetId: number;
}

function formatFileSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(0)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

export function AddVersionModal({
  opened,
  onClose,
  assetId,
}: AddVersionModalProps) {
  const [file, setFile] = useState<File | null>(null);
  const [preview, setPreview] = useState<string | null>(null);
  const [label, setLabel] = useState("");
  const [notes, setNotes] = useState("");
  const addVersion = useAddVersion();
  const surfaces = useMediaWorkspaceSurfaces();

  const reset = () => {
    if (preview) URL.revokeObjectURL(preview);
    setFile(null);
    setPreview(null);
    setLabel("");
    setNotes("");
  };

  const handleClose = () => {
    reset();
    onClose();
  };

  const handleDrop = (files: File[]) => {
    const f = files[0];
    if (!f) return;
    setFile(f);
    if (f.type.startsWith("image/")) {
      setPreview(URL.createObjectURL(f));
    }
  };

  const handleRemoveFile = () => {
    if (preview) URL.revokeObjectURL(preview);
    setFile(null);
    setPreview(null);
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

  const isVideo = file?.type.startsWith("video/");

  return (
    <Modal
      opened={opened}
      onClose={handleClose}
      title={
        <Group gap="xs">
          <ThemeIcon variant="light" color="violet" size="sm" radius="xl">
            <IconVersions size={14} />
          </ThemeIcon>
          <Text fw={600}>Add New Version</Text>
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
            style={surfaces.dropzoneVersion}
          >
            <Stack
              align="center"
              justify="center"
              gap="md"
              mih={160}
              style={{ pointerEvents: "none" }}
            >
              <Dropzone.Accept>
                <ThemeIcon variant="light" color="violet" size={64} radius="xl">
                  <IconUpload size={32} stroke={1.5} />
                </ThemeIcon>
              </Dropzone.Accept>
              <Dropzone.Reject>
                <ThemeIcon variant="light" color="red" size={64} radius="xl">
                  <IconX size={32} stroke={1.5} />
                </ThemeIcon>
              </Dropzone.Reject>
              <Dropzone.Idle>
                <ThemeIcon
                  variant="light"
                  color="violet"
                  size={64}
                  radius="xl"
                  style={{ opacity: 0.8 }}
                >
                  <IconCloudUpload size={32} stroke={1.5} />
                </ThemeIcon>
              </Dropzone.Idle>
              <Stack gap={4} align="center">
                <Text size="md" fw={600} c="var(--mantine-color-text)">
                  Drop the new version file
                </Text>
                <Text size="sm" c="dimmed">
                  or{" "}
                  <Text span c="violet" fw={500}>
                    browse files
                  </Text>{" "}
                  · Max 50 MB
                </Text>
              </Stack>
            </Stack>
          </Dropzone>
        ) : (
          <Paper p="md" radius="lg" withBorder>
            <Group wrap="nowrap" gap="md">
              {preview ? (
                <Image
                  src={preview}
                  alt="Preview"
                  w={72}
                  h={72}
                  radius="md"
                  fit="cover"
                />
              ) : (
                <Box
                  style={{
                    width: 72,
                    height: 72,
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
                    <IconVideo
                      size={28}
                      color="var(--mantine-color-violet-5)"
                    />
                  ) : (
                    <IconFile
                      size={28}
                      color="var(--mantine-color-violet-5)"
                    />
                  )}
                </Box>
              )}

              <Stack gap={2} style={{ flex: 1, minWidth: 0 }}>
                <Text size="sm" fw={600} lineClamp={1}>
                  {file.name}
                </Text>
                <Text size="xs" c="dimmed">
                  {formatFileSize(file.size)}
                </Text>
                {addVersion.isPending && (
                  <Progress
                    size="xs"
                    radius="xl"
                    animated
                    value={100}
                    color="violet"
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
        <Group grow>
          <TextInput
            label="Label"
            placeholder="e.g. edited, color_graded, final"
            value={label}
            onChange={(e) => setLabel(e.currentTarget.value)}
            radius="md"
          />
        </Group>

        <Textarea
          label="Notes"
          placeholder="What changed in this version?"
          value={notes}
          onChange={(e) => setNotes(e.currentTarget.value)}
          minRows={2}
          radius="md"
        />

        {/* Actions */}
        <Group justify="flex-end" mt="xs">
          <Button
            variant="subtle"
            color="gray"
            onClick={handleClose}
            radius="md"
          >
            Cancel
          </Button>
          <Button
            onClick={handleSubmit}
            loading={addVersion.isPending}
            disabled={!file}
            radius="md"
            variant="gradient"
            gradient={{ from: "violet", to: "blue", deg: 135 }}
          >
            Upload Version
          </Button>
        </Group>
      </Stack>
    </Modal>
  );
}
