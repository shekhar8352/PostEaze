import { Box, Text, Stack } from '@mantine/core';
import type { ReactNode } from 'react';

type EmptyChannelPanelProps = {
  icon: ReactNode;
  title: string;
  hint: string;
  background?: string;
};

export function EmptyChannelPanel({
  icon,
  title,
  hint,
  background = 'var(--pe-bg-subtle)',
}: EmptyChannelPanelProps) {
  return (
    <Box
      p="xl"
      style={{
        textAlign: 'center',
        borderRadius: 'var(--pe-radius-md)',
        background,
        border: '1px dashed var(--pe-border)',
      }}
    >
      <Stack gap="sm" align="center">
        <Box c="dimmed" style={{ opacity: 0.45 }}>
          {icon}
        </Box>
        <Text fw={500} c="dimmed">
          {title}
        </Text>
        <Text size="sm" c="dimmed">
          {hint}
        </Text>
      </Stack>
    </Box>
  );
}
