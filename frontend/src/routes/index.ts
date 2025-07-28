import { useRoutes } from "react-router-dom";
import { authRoutes } from "@/features/auth";
import { landingRoutes } from "@/features/landing";

const AppRoutes = () => {
  const routes = [...authRoutes, ...landingRoutes]; // Later add dashboard, posts, etc.
  return useRoutes(routes);
};

export default AppRoutes;
