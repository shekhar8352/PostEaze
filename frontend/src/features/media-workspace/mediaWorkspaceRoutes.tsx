import { lazy, Suspense } from "react";
import type { RouteObject } from "react-router-dom";
import { LoadingFallback } from "@/app/routes/LoadingFallback";

const MediaWorkspace = lazy(() => import("./pages/MediaWorkspace"));
const MediaDetail = lazy(() => import("./pages/MediaDetail"));

const mediaWorkspaceRoutes: RouteObject[] = [
  {
    path: "/workspace",
    element: (
      <Suspense fallback={<LoadingFallback />}>
        <MediaWorkspace />
      </Suspense>
    ),
  },
  {
    path: "/workspace/:assetId",
    element: (
      <Suspense fallback={<LoadingFallback />}>
        <MediaDetail />
      </Suspense>
    ),
  },
];

export default mediaWorkspaceRoutes;
