import { useEffect, useState } from "react";
import {
  ActionIcon,
  Badge,
  Box,
  Button,
  ColorInput,
  Collapse,
  Group,
  NumberInput,
  Select,
  Stack,
  Table,
  Text,
  TextInput,
  Tooltip,
} from "@mantine/core";
import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import type { CSSProperties } from "react";
import { Icons } from "@/app/theme";
import type { Phase, PhaseKind, UpdatePhasePayload } from "../types";

interface PhaseRowProps {
  phase: Phase;
  onSave: (body: UpdatePhasePayload) => Promise<void> | void;
  onDelete: () => Promise<void> | void;
  saving?: boolean;
  deleting?: boolean;
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
 * PhaseRow - A single sortable row in the Studio settings phase list.
 * Click the edit action to reveal an inline form for name, kind, color,
 * and WIP limit. Drag the handle to reorder.
 */
export function PhaseRow({
  phase,
  onSave,
  onDelete,
  saving,
  deleting,
}: PhaseRowProps) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } =
    useSortable({ id: phase.id });

  const style: CSSProperties = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
    background: isDragging ? "var(--mantine-color-gray-0)" : undefined,
  };

  const [editing, setEditing] = useState(false);
  const [name, setName] = useState(phase.name);
  const [kind, setKind] = useState<PhaseKind>(phase.kind);
  const [color, setColor] = useState(phase.color || "#868e96");
  const [wipLimit, setWipLimit] = useState<number | "">(
    phase.wip_limit ?? ""
  );

  // Keep local form state in sync if the phase prop updates (e.g. after save).
  useEffect(() => {
    setName(phase.name);
    setKind(phase.kind);
    setColor(phase.color || "#868e96");
    setWipLimit(phase.wip_limit ?? "");
  }, [phase.id, phase.name, phase.kind, phase.color, phase.wip_limit]);

  const dirty =
    name.trim() !== phase.name ||
    kind !== phase.kind ||
    color !== (phase.color || "#868e96") ||
    (wipLimit === "" ? null : Number(wipLimit)) !== phase.wip_limit;

  const handleSave = async () => {
    const trimmed = name.trim();
    if (!trimmed) return;
    await onSave({
      name: trimmed,
      kind,
      color,
      wip_limit: wipLimit === "" ? null : Number(wipLimit),
    });
    setEditing(false);
  };

  const handleCancel = () => {
    setName(phase.name);
    setKind(phase.kind);
    setColor(phase.color || "#868e96");
    setWipLimit(phase.wip_limit ?? "");
    setEditing(false);
  };

  return (
    <>
      <Table.Tr ref={setNodeRef} style={style}>
        <Table.Td style={{ width: 40 }}>
          <ActionIcon
            variant="subtle"
            color="gray"
            size="sm"
            style={{ cursor: "grab" }}
            aria-label="Drag to reorder"
            {...attributes}
            {...listeners}
          >
            <Icons.GripVertical size={14} />
          </ActionIcon>
        </Table.Td>
        <Table.Td>
          <Group gap="xs" wrap="nowrap">
            <Box
              style={{
                width: 12,
                height: 12,
                borderRadius: 3,
                background: phase.color || "#868e96",
                flexShrink: 0,
              }}
            />
            <Text size="sm" fw={500}>
              {phase.name}
            </Text>
            {phase.is_default && (
              <Badge size="xs" variant="light" color="blue">
                Default
              </Badge>
            )}
          </Group>
        </Table.Td>
        <Table.Td>
          <Badge size="xs" variant="light" color="gray">
            {phase.kind}
          </Badge>
        </Table.Td>
        <Table.Td>
          <Text size="xs" c="dimmed">
            {phase.wip_limit != null ? phase.wip_limit : "—"}
          </Text>
        </Table.Td>
        <Table.Td style={{ width: 96 }}>
          <Group gap={4} justify="flex-end" wrap="nowrap">
            <Tooltip label={editing ? "Close" : "Edit phase"}>
              <ActionIcon
                variant="subtle"
                size="sm"
                color="gray"
                onClick={() => setEditing((v) => !v)}
                aria-label="Edit phase"
              >
                <Icons.Edit size={14} />
              </ActionIcon>
            </Tooltip>
            <Tooltip label="Delete phase">
              <ActionIcon
                variant="subtle"
                size="sm"
                color="red"
                onClick={onDelete}
                loading={deleting}
                aria-label="Delete phase"
              >
                <Icons.Trash size={14} />
              </ActionIcon>
            </Tooltip>
          </Group>
        </Table.Td>
      </Table.Tr>

      {/* Colspan row for the inline editor. */}
      <Table.Tr>
        <Table.Td
          colSpan={5}
          style={{
            padding: editing ? undefined : 0,
            borderTop: editing ? undefined : "none",
          }}
        >
          <Collapse in={editing}>
            <Stack gap="sm" py="xs">
              <Group grow align="flex-start">
                <TextInput
                  label="Name"
                  value={name}
                  onChange={(e) => setName(e.currentTarget.value)}
                  required
                />
                <Select
                  label="Kind"
                  value={kind}
                  onChange={(v) => v && setKind(v as PhaseKind)}
                  data={KIND_OPTIONS}
                  allowDeselect={false}
                  description="Powers automations (e.g. auto-move to Scheduled)"
                />
              </Group>
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
                  description="Empty for no limit"
                  value={wipLimit}
                  onChange={(v) =>
                    setWipLimit(typeof v === "number" ? v : "")
                  }
                  min={0}
                  max={999}
                  allowDecimal={false}
                  clampBehavior="strict"
                />
              </Group>
              <Group justify="flex-end" gap="xs">
                <Button variant="subtle" color="gray" onClick={handleCancel}>
                  Cancel
                </Button>
                <Button
                  onClick={handleSave}
                  disabled={!dirty || !name.trim()}
                  loading={saving}
                >
                  Save changes
                </Button>
              </Group>
            </Stack>
          </Collapse>
        </Table.Td>
      </Table.Tr>
    </>
  );
}
