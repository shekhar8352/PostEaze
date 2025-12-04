import { useAuthInitialization } from '../hooks/useAuthInitialization';

interface AuthProviderProps {
    children: React.ReactNode;
}

/**
 * Initializes authentication state on app startup
 * Fetches user data from backend and shows loading state
 */
export const AuthProvider = ({ children }: AuthProviderProps) => {
    const { isLoading } = useAuthInitialization();

    if (isLoading) {
        return (
            <div style={{
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                height: '100vh'
            }}>
                Loading...
            </div>
        );
    }

    return <>{children}</>;
};
