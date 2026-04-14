import { useMantineColorScheme, useMantineTheme } from "@mantine/core";
import { useMemo } from "react";

/**
 * Theme-aware surfaces for media workspace (fixes hardcoded light-mode colors in dark mode).
 */
export function useMediaWorkspaceSurfaces() {
  const theme = useMantineTheme();
  const { colorScheme } = useMantineColorScheme();
  const isDark = colorScheme === "dark";

  return useMemo(() => {
    const dropzoneUpload = {
      borderWidth: 2,
      borderStyle: "dashed" as const,
      borderColor: isDark ? theme.colors.blue[6] : theme.colors.blue[3],
      background: isDark
        ? `linear-gradient(135deg, ${theme.colors.dark[7]} 0%, ${theme.colors.dark[6]} 100%)`
        : `linear-gradient(135deg, ${theme.colors.blue[0]} 0%, ${theme.colors.cyan[0]} 100%)`,
      transition: "border-color 200ms ease, background 200ms ease",
    };

    const dropzoneVersion = {
      borderWidth: 2,
      borderStyle: "dashed" as const,
      borderColor: isDark ? theme.colors.violet[6] : theme.colors.violet[3],
      background: isDark
        ? `linear-gradient(135deg, ${theme.colors.dark[7]} 0%, ${theme.colors.dark[6]} 100%)`
        : `linear-gradient(135deg, ${theme.colors.violet[0]} 0%, ${theme.colors.blue[0]} 100%)`,
      transition: "border-color 200ms ease, background 200ms ease",
    };

    const checkerboardCss = isDark
      ? `repeating-conic-gradient(${theme.colors.dark[5]} 0% 25%, ${theme.colors.dark[7]} 0% 50%) 50% / 20px 20px`
      : `repeating-conic-gradient(${theme.colors.gray[2]} 0% 25%, transparent 0% 50%) 50% / 20px 20px`;

    const publishInfoBanner = {
      background: isDark
        ? `linear-gradient(135deg, ${theme.colors.dark[7]} 0%, ${theme.colors.dark[6]} 100%)`
        : `linear-gradient(135deg, ${theme.colors.green[0]} 0%, ${theme.colors.teal[0]} 100%)`,
      border: isDark ? `1px solid ${theme.colors.dark[4]}` : undefined,
    };

    const emptyStateIconOrb = {
      background: isDark
        ? `linear-gradient(135deg, ${theme.colors.dark[6]} 0%, ${theme.colors.dark[5]} 100%)`
        : `linear-gradient(135deg, ${theme.colors.blue[0]} 0%, ${theme.colors.cyan[0]} 100%)`,
    };

    /** Hover background for timeline rows and similar */
    const subtleHoverBg = isDark ? theme.colors.dark[5] : theme.colors.gray[0];

    return {
      isDark,
      dropzoneUpload,
      dropzoneVersion,
      checkerboardCss,
      publishInfoBanner,
      emptyStateIconOrb,
      subtleHoverBg,
      /** Subtle strip under preview metadata (optional; prefer default Paper when possible) */
      previewMetaBg: isDark ? theme.colors.dark[6] : theme.colors.gray[0],
      previewMetaBorderTop: `1px solid ${isDark ? theme.colors.dark[4] : theme.colors.gray[2]}`,
    };
  }, [theme, isDark]);
}
