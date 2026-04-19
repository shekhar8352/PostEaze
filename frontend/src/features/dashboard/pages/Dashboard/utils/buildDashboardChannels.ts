import { Icons } from '@/app/theme';
import type { BaseChannelDisplay } from '@/features/channels/types/base.types';
import type { DashboardChannelCard } from '../types';

export function buildDashboardChannelCards(
    connectedChannels: BaseChannelDisplay[] | undefined
): DashboardChannelCard[] {
    const list = connectedChannels ?? [];
    return [
        {
            name: 'Instagram',
            icon: Icons.Instagram,
            provider: 'instagram',
            path: '/channels/instagram',
            connected: list.some((ch) => ch.provider === 'instagram'),
        },
        {
            name: 'Facebook',
            icon: Icons.Facebook,
            provider: 'facebook',
            path: '/channels/facebook',
            connected: list.some((ch) => ch.provider === 'facebook'),
        },
        {
            name: 'YouTube',
            icon: Icons.YouTube,
            provider: 'youtube',
            path: '/channels/youtube',
            connected: list.some((ch) => ch.provider === 'youtube'),
        },
    ];
}
