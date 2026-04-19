const nfCompact = new Intl.NumberFormat(undefined, {
    notation: 'compact',
    maximumFractionDigits: 1,
});

export function formatDashboardInt(n: number): string {
    return nfCompact.format(n);
}
