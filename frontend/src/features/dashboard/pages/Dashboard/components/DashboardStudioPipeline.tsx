import { useMemo } from 'react';
import {
    Badge,
    Box,
    Button,
    Group,
    Paper,
    Skeleton,
    Stack,
    Text,
    Title,
} from '@mantine/core';
import dayjs from 'dayjs';
import { Icons } from '@/app/theme';
import {
    useDefaultStudio,
    useStudioBoard,
} from '@/features/studio/hooks/useStudioQueries';
import type { Phase, Piece, Studio } from '@/features/studio/types';
import dashStyles from '../DashboardPage.module.css';

interface DashboardStudioPipelineProps {
    onNavigateStudio: () => void;
}

interface PhaseBucket {
    phase: Phase;
    pieces: Piece[];
}

const IN_FLIGHT_TERMINAL_KINDS = new Set(['published']);

function pieceLabel(studio: Studio | undefined, plural: boolean): string {
    const base = studio?.piece_label?.trim() || 'Piece';
    if (!plural) return base;
    return /s$/i.test(base) ? base : `${base}s`;
}

function formatDue(dueAt: string | null): string | null {
    if (!dueAt) return null;
    const d = dayjs(dueAt);
    if (!d.isValid()) return null;
    return d.format('MMM D');
}

export function DashboardStudioPipeline({ onNavigateStudio }: DashboardStudioPipelineProps) {
    const { data: studio, isLoading: studioLoading, isError: studioError } =
        useDefaultStudio();
    const studioId = studio?.id;
    const { data: board, isLoading: boardLoading, isError: boardError } =
        useStudioBoard(studioId);

    const isLoading = studioLoading || (!!studioId && boardLoading);
    const isError = studioError || boardError;

    const buckets = useMemo<PhaseBucket[]>(() => {
        if (!board) return [];
        const activePieces = board.pieces.filter(
            (p) => p.status !== 'archived'
        );
        return board.phases
            .filter((ph) => !IN_FLIGHT_TERMINAL_KINDS.has(ph.kind))
            .sort((a, b) => a.order_index - b.order_index)
            .map((phase) => ({
                phase,
                pieces: activePieces.filter((p) => p.phase_id === phase.id),
            }));
    }, [board]);

    const inFlightCount = useMemo(
        () => buckets.reduce((sum, b) => sum + b.pieces.length, 0),
        [buckets]
    );

    const dueSoon = useMemo<Piece[]>(() => {
        if (!board) return [];
        const now = dayjs();
        const horizon = now.add(7, 'day');
        const activePieces = board.pieces.filter(
            (p) => p.status !== 'archived'
        );
        return activePieces
            .filter((p) => {
                if (!p.due_at) return false;
                const d = dayjs(p.due_at);
                return d.isValid() && d.isAfter(now) && d.isBefore(horizon);
            })
            .sort((a, b) => dayjs(a.due_at!).valueOf() - dayjs(b.due_at!).valueOf())
            .slice(0, 5);
    }, [board]);

    const header = (
        <Group justify="space-between" align="flex-start" wrap="wrap" gap="md">
            <div style={{ minWidth: 0, flex: '1 1 16rem' }}>
                <Title order={2} className={dashStyles.sectionTitle} mb={6}>
                    Studio pipeline
                </Title>
                <Text className={dashStyles.sectionDesc}>
                    Active {pieceLabel(studio, true).toLowerCase()} still moving through your
                    phases, plus anything due in the next 7 days.
                </Text>
            </div>
            <Badge variant="light" color="gray" style={{ flexShrink: 0 }}>
                {isLoading
                    ? 'Loading…'
                    : `${inFlightCount} in flight`}
            </Badge>
        </Group>
    );

    const goToStudioButton = (
        <Button
            variant="light"
            rightSection={<Icons.ArrowRight size={16} />}
            onClick={onNavigateStudio}
        >
            Open Studio
        </Button>
    );

    let body: React.ReactNode;

    if (isError) {
        body = (
            <Stack gap="md">
                <Text size="sm" c="red">
                    Could not load your Studio right now. Try opening it directly.
                </Text>
                {goToStudioButton}
            </Stack>
        );
    } else if (isLoading) {
        body = (
            <Stack gap="md">
                <Skeleton height={18} width="60%" />
                <Skeleton height={56} radius="md" />
                <Skeleton height={56} radius="md" />
                <Skeleton height={40} radius="md" />
                {goToStudioButton}
            </Stack>
        );
    } else if (!board || buckets.length === 0) {
        body = (
            <Stack gap="md">
                <Text size="sm" c="dimmed">
                    No phases yet. Create your first phase in the Studio to start tracking
                    work.
                </Text>
                {goToStudioButton}
            </Stack>
        );
    } else {
        body = (
            <Stack gap="lg">
                <Stack gap="xs">
                    <Title order={3} className={dashStyles.subsectionHeading}>
                        By phase
                    </Title>
                    <div className={dashStyles.nextUpList}>
                        <Stack gap={0}>
                            {buckets.map(({ phase, pieces }) => (
                                <div key={phase.id} className={dashStyles.miniPostRow}>
                                    <div style={{ minWidth: 0 }}>
                                        <Group gap={8} align="center" wrap="nowrap">
                                            <span
                                                style={{
                                                    width: 8,
                                                    height: 8,
                                                    borderRadius: 999,
                                                    background: phase.color || 'var(--pe-accent)',
                                                    flexShrink: 0,
                                                    display: 'inline-block',
                                                }}
                                                aria-hidden
                                            />
                                            <Text
                                                size="sm"
                                                fw={600}
                                                lineClamp={1}
                                                c="var(--pe-text)"
                                            >
                                                {phase.name}
                                            </Text>
                                        </Group>
                                        <Text size="xs" c="dimmed" mt={2}>
                                            {pieces.length} {pieceLabel(studio, pieces.length !== 1).toLowerCase()}
                                        </Text>
                                    </div>
                                    <Badge
                                        size="sm"
                                        variant="light"
                                        color={pieces.length > 0 ? 'blue' : 'gray'}
                                        style={{ flexShrink: 0 }}
                                    >
                                        {pieces.length}
                                    </Badge>
                                </div>
                            ))}
                        </Stack>
                    </div>
                </Stack>

                <Stack gap="xs">
                    <Title order={3} className={dashStyles.subsectionHeading}>
                        Due in the next 7 days
                    </Title>
                    {dueSoon.length === 0 ? (
                        <Text size="sm" c="dimmed" py="xs">
                            Nothing due soon. You&apos;re clear.
                        </Text>
                    ) : (
                        <div className={dashStyles.nextUpList}>
                            <Stack gap={0}>
                                {dueSoon.map((p) => (
                                    <div key={p.id} className={dashStyles.miniPostRow}>
                                        <div style={{ minWidth: 0 }}>
                                            <Text
                                                size="sm"
                                                fw={600}
                                                lineClamp={1}
                                                c="var(--pe-text)"
                                            >
                                                {p.title?.trim() || `Untitled ${pieceLabel(studio, false).toLowerCase()}`}
                                            </Text>
                                            <Text size="xs" c="dimmed">
                                                {p.content_type}
                                                {formatDue(p.due_at)
                                                    ? ` · due ${formatDue(p.due_at)}`
                                                    : ''}
                                            </Text>
                                        </div>
                                        <Badge
                                            size="sm"
                                            variant="light"
                                            color="orange"
                                            style={{ flexShrink: 0 }}
                                        >
                                            Due
                                        </Badge>
                                    </div>
                                ))}
                            </Stack>
                        </div>
                    )}
                </Stack>

                <Box>{goToStudioButton}</Box>
            </Stack>
        );
    }

    return (
        <Paper p="xl" radius="lg" withBorder={false} className={dashStyles.section}>
            <Stack gap="lg" align="stretch">
                {header}
                {body}
            </Stack>
        </Paper>
    );
}
