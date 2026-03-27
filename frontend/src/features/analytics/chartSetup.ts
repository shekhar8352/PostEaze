import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
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
    Title,
    Tooltip,
    Legend,
    Filler
  );
  registered = true;
}

export const chartColors = {
  accent: "#1d4ed8",
  accent2: "#7c3aed",
  accent3: "#0d9488",
  accent4: "#ea580c",
  accent5: "#db2777",
  accent6: "#4f46e5",
  muted: "#94a3b8",
  grid: "#e2e8f0",
  text: "#64748b",
};
