import { lazy, Suspense } from "react";
import type { RouteObject } from "react-router-dom";
import { LoadingFallback } from "@/app/routes/LoadingFallback";

const CalendarPage = lazy(() => import("./pages/CalendarPage"));

const calendarRoutes: RouteObject[] = [
  {
    path: "/calendar",
    element: (
      <Suspense fallback={<LoadingFallback />}>
        <CalendarPage />
      </Suspense>
    ),
  },
];

export default calendarRoutes;
