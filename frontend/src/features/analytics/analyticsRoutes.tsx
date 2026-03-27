import { lazy, Suspense } from "react";
import type { RouteObject } from "react-router-dom";
import { LoadingFallback } from "@/app/routes/LoadingFallback";

const AnalyticsPage = lazy(() => import("./pages/AnalyticsPage"));

const analyticsRoutes: RouteObject[] = [
  {
    path: "/analytics",
    element: (
      <Suspense fallback={<LoadingFallback />}>
        <AnalyticsPage />
      </Suspense>
    ),
  },
];

export default analyticsRoutes;
