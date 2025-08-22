import { Navigate, Outlet } from "react-router-dom";
// import { useAuth } from "@/features/auth/hooks/useAuth";

// For protecting a group of routes (nested under Outlet)
export const ProtectedLayout = () => {
  // const { user } = useAuth();
  const user = null;
  if (!user) return <Navigate to="/login" replace />;

  return <Outlet />; // renders child routes
};


// For protecting a single route
const ProtectedRoute = ({ element }: { element: React.JSX.Element }) => {
  // const { user } = useAuth();
  const user = null;

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  return element;
};

export default ProtectedRoute;
