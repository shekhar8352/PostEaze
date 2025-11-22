import { Navigate } from "react-router-dom";
import { useAuth } from "@/features/auth";
import { getRedirectPath } from "@/features/auth/utils";
import { LoadingFallback } from "./LoadingFallback";

interface PublicRouteProps {
    element: React.JSX.Element;
}

/**
 * PublicRoute component for auth pages (login, register, etc.)
 * Redirects authenticated users to dashboard or intended destination
 */
const PublicRoute = ({ element }: PublicRouteProps) => {
    const { isAuthenticated, loading } = useAuth();

    // Show loading while checking authentication
    if (loading) {
        return <LoadingFallback />;
    }

    // Redirect authenticated users away from auth pages
    if (isAuthenticated) {
        const redirectPath = getRedirectPath();
        return <Navigate to={redirectPath} replace />;
    }

    return element;
};

export default PublicRoute;
