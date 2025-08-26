import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { authService } from "./authService";
import {
  type LoginRequest,
  type RegisterRequest,
  type LoginFormData,
  type RegisterFormData,
} from "../types";

export const useLogin = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: LoginFormData) => {
      const loginRequest: LoginRequest = {
        email: data.email,
        password: data.password,
        firebase_uid: "", // Temporary - will be overwritten
        firebase_token: "", // Temporary - will be overwritten
        email_verified: true, // Will be verified by Firebase
        provider: "email",
      };

      return authService.loginUser(loginRequest);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["currentUser"] });
    },
    onError: (error: any) => {
      console.error("Login failed:", error.message);
    },
  });
};

export const useRegister = () => {
  return useMutation({
    mutationFn: (data: RegisterFormData) => {
      const registerRequest: RegisterRequest = {
        name: data.name,
        email: data.email,
        password: data.password,
        confirmPassword: data.confirmPassword,
        terms: true, // You might want to add this to your form
        // These Firebase fields aren't used in the initial registration
        firebase_uid: "", // Temporary - not used in registerUser
        firebase_token: "", // Temporary - not used in registerUser
        display_name: data.name,
        email_verified: true, // Will be verified via email
        provider: "email",
      };

      return authService.registerUser(registerRequest);
    },
    onSuccess: () => {
      // Don't invalidate currentUser as user is not logged in yet
      console.log("Registration initiated, verification email sent");
    },
    onError: (error: any) => {
      console.error("Registration failed:", error.message);
    },
  });
};

// Complete registration after email verification
export const useCompleteRegistration = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      email,
      password,
      name,
    }: {
      email: string;
      password: string;
      name: string;
    }) => authService.completeRegistration(email, password, name),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["currentUser"] });
    },
    onError: (error: any) => {
      console.error("Registration completion failed:", error.message);
    },
  });
};

// Resend verification email
export const useResendVerificationEmail = () => {
  return useMutation({
    mutationFn: ({ email, password }: { email: string; password: string }) =>
      authService.resendVerificationEmail(email, password),
    onError: (error: any) => {
      console.error("Resend verification failed:", error.message);
    },
  });
};

// Check email verification status
export const useCheckEmailVerification = () => {
  return useMutation({
    mutationFn: ({ email, password }: { email: string; password: string }) =>
      authService.checkEmailVerificationStatus(email, password),
    onError: (error: any) => {
      console.error("Verification check failed:", error.message);
    },
  });
};

// Social Auth (already verified)
export const useGoogleAuth = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => authService.googleAuth(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["currentUser"] });
    },
    onError: (error: any) => {
      console.error("Google auth failed:", error.message);
    },
  });
};

export const useFacebookAuth = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => authService.facebookAuth(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["currentUser"] });
    },
    onError: (error: any) => {
      console.error("Facebook auth failed:", error.message);
    },
  });
};

// Password reset
export const useForgotPassword = () => {
  return useMutation({
    mutationFn: (email: string) => authService.firebasePasswordReset(email),
    onSuccess: () => {
      console.log("Password reset email sent!");
    },
    onError: (error: any) => {
      console.error("Password reset failed:", error.message);
    },
  });
};

// Logout
export const useLogout = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => authService.logoutUser(),
    onSuccess: () => {
      queryClient.clear();
    },
  });
};

// Get Current User (from your backend)
export const useCurrentUser = () => {
  return useQuery({
    queryKey: ["currentUser"],
    queryFn: () => authService.getCurrentUser(),
    enabled: authService.isAuthenticated(),
  });
};

// Backend password reset (if you have it)
export const useResetPassword = () => {
  return useMutation({
    mutationFn: ({ token, password }: { token: string; password: string }) =>
      authService.resetPassword(token, password),
    onSuccess: () => {
      console.log("Password reset successful!");
    },
    onError: (error: any) => {
      console.error("Password reset failed:", error.message);
    },
  });
};
