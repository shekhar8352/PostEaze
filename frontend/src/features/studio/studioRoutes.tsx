import { lazy, Suspense } from "react";
import type { RouteObject } from "react-router-dom";
import { LoadingFallback } from "@/app/routes/LoadingFallback";

const StudioBoard = lazy(() => import("./pages/StudioBoard"));
const StudioSettings = lazy(() => import("./pages/StudioSettings"));

const studioRoutes: RouteObject[] = [
  {
    path: "/studio",
    element: (
      <Suspense fallback={<LoadingFallback />}>
        <StudioBoard />
      </Suspense>
    ),
  },
  {
    path: "/studio/settings",
    element: (
      <Suspense fallback={<LoadingFallback />}>
        <StudioSettings />
      </Suspense>
    ),
  },
];

export default studioRoutes;
