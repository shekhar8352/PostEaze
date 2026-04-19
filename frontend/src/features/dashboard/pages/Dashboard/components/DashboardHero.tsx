import { Box, Button, Group, Paper, Text, Title } from '@mantine/core';
import { Link } from 'react-router-dom';
import { Icons } from '@/app/theme';
import dashStyles from '../DashboardPage.module.css';

interface DashboardHeroProps {
    displayName: string;
}

export function DashboardHero({ displayName }: DashboardHeroProps) {
    return (
        <Paper className={dashStyles.hero} shadow="none" radius="lg" withBorder={false}>
            <Box className={dashStyles.heroInner}>
                <Title order={1} className={dashStyles.heroTitle}>
                    Welcome back, {displayName}
                </Title>
                <Text className={dashStyles.heroSubtitle}>
                    Publishing pipeline, delivery health, and Instagram performance at a glance. Connect channels and
                    keep your calendar filled—then tune reach and engagement from Analytics.
                </Text>
                <Group gap="sm" mt="md" wrap="wrap">
                    <Button component={Link} to="/calendar" leftSection={<Icons.Calendar size={18} />} variant="light">
                        Open calendar
                    </Button>
                    <Button component={Link} to="/analytics" leftSection={<Icons.ChartBar size={18} />} variant="default">
                        Full analytics
                    </Button>
                </Group>
            </Box>
        </Paper>
    );
}
