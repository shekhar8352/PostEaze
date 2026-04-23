import { useState } from "react";
import {
  Button,
  ColorInput,
  Group,
  Modal,
  NumberInput,
  Select,
  Stack,
  TextInput,
} from "@mantine/core";
import type { CreatePhasePayload, PhaseKind } from "../types";

interface CreatePhaseModalProps {
  opened: boolean;
  onClose: () => void;
  onCreate: (body: CreatePhasePayload) => Promise<void> | void;
  submitting?: boolean;
}

const KIND_OPTIONS: { value: PhaseKind; label: string }[] = [
  { value: "idea", label: "Idea" },
  { value: "script", label: "Script" },
  { value: "shoot", label: "Shoot" },
  { value: "edit", label: "Edit" },
  { value: "review", label: "Review" },
  { value: "scheduled", label: "Scheduled" },
  { value: "published", label: "Published" },
  { value: "custom", label: "Custom" },
];

/**
 * CreatePhaseModal - Adds a new phase to the end of the pipeline.
 * Kind governs automation hooks, so users can rename phases without
 * breaking behavior.
 */
export function CreatePhaseModal({
  opened,
  onClose,
  onCreate,
  submitting,
}: CreatePhaseModalProps) {
  const [name, setName] = useState("");
  const [kind, setKind] = useState<PhaseKind>("custom");
  const [color, setColor] = useState("#6b7280");
  const [wipLimit, setWipLimit] = useState<number | "">("");

  const reset = () => {
    setName("");
    setKind("custom");
    setColor("#6b7280");
    setWipLimit("");
  };

  const handleClose = () => {
    reset();
    onClose();
  };

  const handleSubmit = async () => {
    const trimmed = name.trim();
    if (!trimmed) return;
    await onCreate({
      name: trimmed,
      kind,
      color,
      wip_limit: wipLimit === "" ? null : Number(wipLimit),
    });
    reset();
  };

  return (
    <Modal
      opened={opened}
      onClose={handleClose}
      title="Add phase"
      centered
      size="md"
    >
      <Stack gap="sm">
        <TextInput
          label="Name"
          placeholder="e.g. Graphics review"
          value={name}
          onChange={(e) => setName(e.currentTarget.value)}
          data-autofocus
          required
        />
        <Select
          label="Kind"
          description="Used by automations; safe to leave as Custom."
          data={KIND_OPTIONS}
          value={kind}
          onChange={(v) => v && setKind(v as PhaseKind)}
          allowDeselect={false}
        />
        <Group grow align="flex-start">
          <ColorInput
            label="Color"
            value={color}
            onChange={setColor}
            format="hex"
            swatches={[
              "#6366f1",
              "#8b5cf6",
              "#ec4899",
              "#f59e0b",
              "#10b981",
              "#0ea5e9",
              "#22c55e",
              "#ef4444",
              "#6b7280",
            ]}
          />
          <NumberInput
            label="WIP limit"
            description="Optional"
            value={wipLimit}
            onChange={(v) => setWipLimit(typeof v === "number" ? v : "")}
            min={0}
            max={999}
            allowDecimal={false}
          />
        </Group>
        <Group justify="flex-end" gap="xs" mt="xs">
          <Button variant="subtle" color="gray" onClick={handleClose}>
            Cancel
          </Button>
          <Button
            onClick={handleSubmit}
            disabled={!name.trim()}
            loading={submitting}
          >
            Add phase
          </Button>
        </Group>
      </Stack>
    </Modal>
  );
}
