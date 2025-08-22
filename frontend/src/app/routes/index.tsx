import { useRoutes } from "react-router-dom";
import { authRoutes } from "@/features/auth";
import { landingRoutes } from "@/features/landing";
import {ProtectedLayout} from "./ProtectedRoute";

const AppRoutes = () => {
  const routes = [
    ...authRoutes,
    
    {
      element: <ProtectedLayout />,
      children: [
         // Write all protected routes here
        ...landingRoutes
      ],
    },
  ];

  return useRoutes(routes);
};

export default AppRoutes;
