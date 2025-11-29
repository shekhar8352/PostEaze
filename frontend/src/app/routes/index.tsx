import { Suspense } from "react";
import { useRoutes } from "react-router-dom";
import { authRoutes } from "@/features/auth";
import { landingRoutes } from "@/features/landing";
import { dashboardRoutes } from "@/features/dashboard";
import { channelRoutes } from "@/features/channels";
import { MainLayout } from "@/features/layout";
import { ProtectedLayout } from "./ProtectedRoute";
import { NotFound } from "./NotFound";
import { LoadingFallback } from "./LoadingFallback";

const AppRoutes = () => {
  const routes = [
    ...authRoutes,

    // Protected routes wrapped with authentication check
    {
      element: <ProtectedLayout />,
      children: [
        {
          element: <MainLayout />,
          children: [
            // All protected routes here
            ...dashboardRoutes,
            ...landingRoutes,
            ...channelRoutes,
          ],
        },
      ],
    },

    // 404 Not Found - must be last
    {
      path: "*",
      element: <NotFound />,
    },
  ];

  return (
    <Suspense fallback={<LoadingFallback />}>
      {useRoutes(routes)}
    </Suspense>
  );
};

export default AppRoutes;

