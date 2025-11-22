import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useDispatch } from "react-redux";
import { authService } from "./authService";
import { type LoginFormData, type RegisterFormData } from "../types";
import { setUser, clearUser } from "../authSlice";

export const useLogin = () => {
  const queryClient = useQueryClient();
  const dispatch = useDispatch();
  
  return useMutation({
    mutationFn: (data: LoginFormData) => {
      const loginRequest: { email: string; password: string } = {
        email: data.email,
        password: data.password,
      };

      return authService.loginUser(loginRequest);
    },
    onSuccess: (data) => {
      // Update Redux store with user data
      dispatch(setUser(data.user));
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
      const registerRequest: { name: string; email: string; password: string } =
        {
          name: data.name,
          email: data.email,
          password: data.password,
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
  const dispatch = useDispatch();
  
  return useMutation({
    mutationFn: ({ email, password }: { email: string; password: string }) =>
      authService.completeRegistration(email, password),
    onSuccess: (data) => {
      // Update Redux store with user data
      dispatch(setUser(data.user));
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
  const dispatch = useDispatch();
  
  return useMutation({
    mutationFn: () => authService.googleAuth(),
    onSuccess: (data) => {
      // Update Redux store with user data
      dispatch(setUser(data.user));
      queryClient.invalidateQueries({ queryKey: ["currentUser"] });
    },
    onError: (error: any) => {
      console.error("Google auth failed:", error.message);
    },
  });
};

export const useFacebookAuth = () => {
  const queryClient = useQueryClient();
  const dispatch = useDispatch();
  
  return useMutation({
    mutationFn: () => authService.facebookAuth(),
    onSuccess: (data) => {
      // Update Redux store with user data
      dispatch(setUser(data.user));
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
  const dispatch = useDispatch();
  
  return useMutation({
    mutationFn: () => authService.logoutUser(),
    onSuccess: () => {
      // Clear Redux store
      dispatch(clearUser());
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
