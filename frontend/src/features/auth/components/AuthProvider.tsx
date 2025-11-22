import { useEffect } from 'react';
import { useAuth } from '../hooks/useAuth';

interface AuthProviderProps {
    children: React.ReactNode;
}

/**
 * AuthProvider component that initializes authentication state on mount
 * This should wrap the entire app to ensure auth state is hydrated from localStorage
 */
export const AuthProvider = ({ children }: AuthProviderProps) => {
    const { initializeAuth } = useAuth();

    useEffect(() => {
        // Initialize auth state from localStorage when app mounts
        initializeAuth();
    }, [initializeAuth]);

    return <>{children}</>;
};
