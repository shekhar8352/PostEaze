import { useMemo } from "react";
import { useComputedColorScheme } from "@mantine/core";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler,
} from "chart.js";

let registered = false;

export function registerChartJs(): void {
  if (registered) return;
  ChartJS.register(
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    BarElement,
    ArcElement,
    Title,
    Tooltip,
    Legend,
    Filler
  );
  registered = true;
}

function getCssVar(name: string): string {
  return getComputedStyle(document.documentElement).getPropertyValue(name).trim();
}

export const chartAccents = {
  accent: "#1d4ed8",
  accent2: "#7c3aed",
  accent3: "#0d9488",
  accent4: "#ea580c",
  accent5: "#db2777",
  accent6: "#4f46e5",
} as const;

/**
 * Returns chart grid / text / muted colors that follow the active color scheme.
 * Re-computes when Mantine's color scheme changes.
 */
export function useChartTheme() {
  const scheme = useComputedColorScheme("light", { getInitialValueInEffect: true });

  return useMemo(() => {
    const grid = getCssVar("--pe-border") || (scheme === "dark" ? "#2a2d38" : "#e2e8f0");
    const text = getCssVar("--pe-text-muted") || (scheme === "dark" ? "#7e849a" : "#64748b");
    const muted = getCssVar("--pe-text-muted") || text;

    return { ...chartAccents, grid, text, muted };
  }, [scheme]);
}

/** @deprecated Use useChartTheme() hook instead for dark-mode support */
export const chartColors = {
  ...chartAccents,
  muted: "#94a3b8",
  grid: "#e2e8f0",
  text: "#64748b",
};
