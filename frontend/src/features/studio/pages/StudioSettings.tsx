import { useEffect, useMemo, useState } from "react";
import {
  Alert,
  Badge,
  Button,
  Card,
  Center,
  Group,
  Loader,
  Modal,
  Paper,
  Select,
  Stack,
  Table,
  Text,
  TextInput,
  Title,
} from "@mantine/core";
import { notifications } from "@mantine/notifications";
import {
  DndContext,
  PointerSensor,
  closestCenter,
  useSensor,
  useSensors,
  type DragEndEvent,
} from "@dnd-kit/core";
import {
  SortableContext,
  arrayMove,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { Link } from "react-router-dom";
import { Icons } from "@/app/theme";
import type { Phase, UpdatePhasePayload } from "../types";
import {
  useCreatePhase,
  useDefaultStudio,
  useDeletePhase,
  usePhases,
  useReorderPhases,
  useUpdatePhase,
  useUpdateStudio,
} from "../hooks/useStudioQueries";
import { PhaseRow } from "../components/PhaseRow";
import { CreatePhaseModal } from "../components/CreatePhaseModal";

const PIECE_LABEL_PRESETS = [
  { value: "Piece", label: "Piece" },
  { value: "Drop", label: "Drop" },
  { value: "Post", label: "Post" },
  { value: "Story", label: "Story" },
  { value: "Card", label: "Card" },
  { value: "custom", label: "Custom…" },
];

/**
 * StudioSettings - Manage studio-level general settings, piece label
 * terminology, and the phase pipeline (CRUD + drag-to-reorder).
 */
export default function StudioSettings() {
  const {
    data: studio,
    isLoading: studioLoading,
    error: studioError,
  } = useDefaultStudio();

  const {
    data: phases,
    isLoading: phasesLoading,
    error: phasesError,
  } = usePhases(studio?.id);

  if (studioLoading || phasesLoading) {
    return (
      <Center style={{ minHeight: 320 }}>
        <Loader />
      </Center>
    );
  }

  const error = studioError ?? phasesError;
  if (error) {
    return (
      <Alert
        color="red"
        icon={<Icons.AlertCircle size={18} />}
        title="Failed to load Studio settings"
        m="md"
      >
        {error instanceof Error ? error.message : "Something went wrong."}
      </Alert>
    );
  }

  if (!studio) return null;

  return (
    <Stack gap="lg" p="md">
      <Group justify="space-between" align="flex-start">
        <Stack gap={4}>
          <Title order={2}>Studio settings</Title>
          <Text size="sm" c="dimmed">
            Configure terminology and your content pipeline.
          </Text>
        </Stack>
        <Button
          component={Link}
          to="/studio"
          variant="subtle"
          leftSection={<Icons.ArrowLeft size={16} />}
        >
          Back to board
        </Button>
      </Group>

      <GeneralSection
        studioId={studio.id}
        name={studio.name}
        pieceLabel={studio.piece_label}
      />

      <PhasesSection studioId={studio.id} phases={phases ?? []} />
    </Stack>
  );
}

// ---------------------------------------------------------------------------
// General section
// ---------------------------------------------------------------------------

interface GeneralSectionProps {
  studioId: number;
  name: string;
  pieceLabel: string;
}

function GeneralSection({ studioId, name, pieceLabel }: GeneralSectionProps) {
  const update = useUpdateStudio();

  const initialPresetMatch = PIECE_LABEL_PRESETS.find(
    (p) => p.value === pieceLabel
  );
  const initialPreset = initialPresetMatch ? pieceLabel : "custom";

  const [studioName, setStudioName] = useState(name);
  const [labelPreset, setLabelPreset] = useState<string>(initialPreset);
  const [customLabel, setCustomLabel] = useState(
    initialPresetMatch ? "" : pieceLabel
  );

  useEffect(() => {
    setStudioName(name);
    const match = PIECE_LABEL_PRESETS.find((p) => p.value === pieceLabel);
    setLabelPreset(match ? pieceLabel : "custom");
    setCustomLabel(match ? "" : pieceLabel);
  }, [name, pieceLabel]);

  const effectiveLabel =
    labelPreset === "custom" ? customLabel.trim() : labelPreset;

  const dirty =
    studioName.trim() !== name ||
    (effectiveLabel !== pieceLabel && effectiveLabel.length > 0);

  const handleSave = async () => {
    const trimmedName = studioName.trim();
    if (!trimmedName) return;
    if (!effectiveLabel) return;
    await update.mutateAsync({
      id: studioId,
      body: { name: trimmedName, piece_label: effectiveLabel },
    });
    notifications.show({
      title: "Settings saved",
      message: "Studio updated.",
      color: "green",
    });
  };

  return (
    <Card withBorder radius="md" p="md">
      <Stack gap="md">
        <Group justify="space-between" align="flex-start">
          <Stack gap={2}>
            <Title order={4}>General</Title>
            <Text size="xs" c="dimmed">
              Shown in the app header and throughout the Studio board.
            </Text>
          </Stack>
          <Badge variant="light" color="gray">
            Label: {pieceLabel}
          </Badge>
        </Group>

        <Group grow align="flex-start">
          <TextInput
            label="Studio name"
            value={studioName}
            onChange={(e) => setStudioName(e.currentTarget.value)}
            required
          />
          <Select
            label="Item label"
            description='Singular noun used for cards (e.g. "Piece", "Drop").'
            data={PIECE_LABEL_PRESETS}
            value={labelPreset}
            onChange={(v) => v && setLabelPreset(v)}
            allowDeselect={false}
          />
        </Group>

        {labelPreset === "custom" && (
          <TextInput
            label="Custom label"
            placeholder="e.g. Episode"
            value={customLabel}
            onChange={(e) => setCustomLabel(e.currentTarget.value)}
            required
          />
        )}

        <Group justify="flex-end">
          <Button
            onClick={handleSave}
            disabled={!dirty || !studioName.trim() || !effectiveLabel}
            loading={update.isPending}
          >
            Save general
          </Button>
        </Group>
      </Stack>
    </Card>
  );
}

// ---------------------------------------------------------------------------
// Phases section
// ---------------------------------------------------------------------------

interface PhasesSectionProps {
  studioId: number;
  phases: Phase[];
}

function PhasesSection({ studioId, phases }: PhasesSectionProps) {
  const createPhase = useCreatePhase(studioId);
  const updatePhase = useUpdatePhase(studioId);
  const deletePhase = useDeletePhase(studioId);
  const reorderPhases = useReorderPhases(studioId);

  const [createOpen, setCreateOpen] = useState(false);
  const [localOrder, setLocalOrder] = useState<number[]>(() =>
    [...phases].sort((a, b) => a.order_index - b.order_index).map((p) => p.id)
  );
  const [pendingSaveId, setPendingSaveId] = useState<number | null>(null);
  const [pendingDeleteId, setPendingDeleteId] = useState<number | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<Phase | null>(null);

  // Keep local order in sync when server data changes (e.g. after create/delete).
  useEffect(() => {
    const sorted = [...phases]
      .sort((a, b) => a.order_index - b.order_index)
      .map((p) => p.id);
    setLocalOrder(sorted);
  }, [phases]);

  const phasesById = useMemo(() => {
    const map = new Map<number, Phase>();
    for (const p of phases) map.set(p.id, p);
    return map;
  }, [phases]);

  const orderedPhases = useMemo(
    () =>
      localOrder
        .map((id) => phasesById.get(id))
        .filter((p): p is Phase => Boolean(p)),
    [localOrder, phasesById]
  );

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } })
  );

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    if (!over || active.id === over.id) return;
    const oldIndex = localOrder.indexOf(Number(active.id));
    const newIndex = localOrder.indexOf(Number(over.id));
    if (oldIndex === -1 || newIndex === -1) return;

    const next = arrayMove(localOrder, oldIndex, newIndex);
    const previous = localOrder;
    setLocalOrder(next);

    reorderPhases.mutate(next, {
      onError: () => {
        // Roll back if the server refused.
        setLocalOrder(previous);
      },
    });
  };

  const handleSave = async (phaseId: number, body: UpdatePhasePayload) => {
    setPendingSaveId(phaseId);
    try {
      await updatePhase.mutateAsync({ phaseId, body });
      notifications.show({
        title: "Phase updated",
        message: "Your changes are saved.",
        color: "green",
      });
    } finally {
      setPendingSaveId(null);
    }
  };

  const confirmDelete = (phase: Phase) => {
    setDeleteTarget(phase);
  };

  const handleConfirmDelete = async () => {
    if (!deleteTarget) return;
    setPendingDeleteId(deleteTarget.id);
    try {
      await deletePhase.mutateAsync(deleteTarget.id);
      setDeleteTarget(null);
    } finally {
      setPendingDeleteId(null);
    }
  };

  const handleCreate = async (body: Parameters<typeof createPhase.mutateAsync>[0]) => {
    await createPhase.mutateAsync(body);
    setCreateOpen(false);
    notifications.show({
      title: "Phase added",
      message: "New phase appended to your pipeline.",
      color: "green",
    });
  };

  return (
    <Card withBorder radius="md" p="md">
      <Stack gap="md">
        <Group justify="space-between" align="flex-start">
          <Stack gap={2}>
            <Title order={4}>Phases</Title>
            <Text size="xs" c="dimmed">
              Drag to reorder. Click the pencil to edit name, kind, color,
              and WIP limit.
            </Text>
          </Stack>
          <Button
            leftSection={<Icons.Plus size={14} />}
            onClick={() => setCreateOpen(true)}
            variant="light"
          >
            Add phase
          </Button>
        </Group>

        {orderedPhases.length === 0 ? (
          <Paper withBorder p="md" radius="sm">
            <Text size="sm" c="dimmed" ta="center">
              No phases yet. Add your first phase to start organizing pieces.
            </Text>
          </Paper>
        ) : (
          <DndContext
            sensors={sensors}
            collisionDetection={closestCenter}
            onDragEnd={handleDragEnd}
          >
            <Table verticalSpacing="xs">
              <Table.Thead>
                <Table.Tr>
                  <Table.Th style={{ width: 40 }} />
                  <Table.Th>Name</Table.Th>
                  <Table.Th>Kind</Table.Th>
                  <Table.Th>WIP</Table.Th>
                  <Table.Th style={{ width: 96 }} />
                </Table.Tr>
              </Table.Thead>
              <Table.Tbody>
                <SortableContext
                  items={orderedPhases.map((p) => p.id)}
                  strategy={verticalListSortingStrategy}
                >
                  {orderedPhases.map((phase) => (
                    <PhaseRow
                      key={phase.id}
                      phase={phase}
                      onSave={(body) => handleSave(phase.id, body)}
                      onDelete={() => confirmDelete(phase)}
                      saving={pendingSaveId === phase.id}
                      deleting={pendingDeleteId === phase.id}
                    />
                  ))}
                </SortableContext>
              </Table.Tbody>
            </Table>
          </DndContext>
        )}

        {reorderPhases.isPending && (
          <Group gap="xs">
            <Loader size="xs" />
            <Text size="xs" c="dimmed">
              Saving order…
            </Text>
          </Group>
        )}
      </Stack>

      <CreatePhaseModal
        opened={createOpen}
        onClose={() => setCreateOpen(false)}
        onCreate={handleCreate}
        submitting={createPhase.isPending}
      />

      <Modal
        opened={deleteTarget !== null}
        onClose={() => setDeleteTarget(null)}
        title="Delete phase"
        centered
        size="sm"
      >
        <Stack gap="sm">
          <Text size="sm">
            Are you sure you want to delete{" "}
            <Text span fw={600}>
              {deleteTarget?.name}
            </Text>
            ? This action cannot be undone. Phases with active pieces cannot
            be deleted — move their pieces first.
          </Text>
          <Group justify="flex-end" gap="xs">
            <Button
              variant="subtle"
              color="gray"
              onClick={() => setDeleteTarget(null)}
            >
              Cancel
            </Button>
            <Button
              color="red"
              onClick={handleConfirmDelete}
              loading={pendingDeleteId === deleteTarget?.id}
            >
              Delete phase
            </Button>
          </Group>
        </Stack>
      </Modal>
    </Card>
  );
}
