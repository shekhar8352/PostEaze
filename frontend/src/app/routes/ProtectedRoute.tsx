import { Navigate, Outlet, useLocation } from "react-router-dom";
import { useAuth } from "@/features/auth";
import { LoadingFallback } from "./LoadingFallback";
import { createLoginRedirect } from "@/features/auth/utils";

// For protecting a group of routes (nested under Outlet)
export const ProtectedLayout = () => {
  const { isAuthenticated, loading } = useAuth();
  const location = useLocation();

  // Show loading while checking authentication
  if (loading) {
    return <LoadingFallback />;
  }

  // Redirect to login with return URL if not authenticated
  if (!isAuthenticated) {
    const loginPath = createLoginRedirect(location.pathname + location.search);
    return <Navigate to={loginPath} replace />;
  }

  return <Outlet />; // renders child routes (MainLayout wraps these in routes/index.tsx)
};

// For protecting a single route
const ProtectedRoute = ({ element }: { element: React.JSX.Element }) => {
  const { isAuthenticated, loading } = useAuth();
  const location = useLocation();

  // Show loading while checking authentication
  if (loading) {
    return <LoadingFallback />;
  }

  // Redirect to login with return URL if not authenticated
  if (!isAuthenticated) {
    const loginPath = createLoginRedirect(location.pathname + location.search);
    return <Navigate to={loginPath} replace />;
  }

  return element;
};

export default ProtectedRoute;
