import { BaseService } from "@/services/base/BaseService";
import apiClient from "@/services/api/client";
import { type ApiResponse } from "@/services/api/types";
import { type User, type LoginRequest, type RegisterRequest } from "../types";
import { firebaseHelper } from "./firebaseHelper";

class AuthService extends BaseService {
  constructor() {
    super("/v1/auth");
  }

  // Auth-specific methods that don't follow CRUD pattern
  async login(
    data: LoginRequest
  ): Promise<{ user: User; access_token: string; refresh_token: string }> {
    const response = await apiClient.post<ApiResponse<any>>(
      `${this.endpoint}/authenticate`,
      data
    );
    return response.data.data;
  }

  async register(
    data: RegisterRequest
  ): Promise<{ user: User; access_token: string; refresh_token: string }> {
    const response = await apiClient.post<ApiResponse<any>>(
      `${this.endpoint}/authenticate`,
      data
    );
    return response.data.data;
  }

  async logout(): Promise<void> {
    const refreshToken = localStorage.getItem("refresh_token");
    await apiClient.post(`${this.endpoint}/logout`, {
      refresh_token: refreshToken,
    });
  }

  async getCurrentUser(): Promise<User> {
    const response = await apiClient.get<ApiResponse<User>>(
      `${this.endpoint}/me`
    );
    return response.data.data;
  }

  async forgotPassword(email: string): Promise<void> {
    await apiClient.post(`${this.endpoint}/forgot-password`, { email });
  }

  async resetPassword(token: string, password: string): Promise<void> {
    await apiClient.post(`${this.endpoint}/reset-password`, {
      token,
      password,
    });
  }

  // Business logic methods
  async loginUser(credentials: { email: string; password: string }) {
    try {
      // Step 1: Authenticate with Firebase (throws error if not verified)
      const firebaseData = await firebaseHelper.loginWithEmail(
        credentials.email,
        credentials?.password
      );

      // Step 2: Send to backend (only verified users reach here)
      const loginData: LoginRequest = {
        firebase_id: firebaseData.firebase_uid,
        firebase_token: firebaseData.firebase_token,
        platform: "email",
      };

      const response = await this.login(loginData);

      // Store your backend tokens
      localStorage.setItem("auth_token", response.access_token);
      localStorage.setItem("refresh_token", response.refresh_token);
      localStorage.setItem("user", JSON.stringify(response.user));

      return response;
    } catch (error: any) {
      throw new Error(error.message || "Login failed");
    }
  }

  async registerUser(data: { name: string; email: string; password: string }) {
    try {
      // Step 1: Register with Firebase (sends verification email)
      const result = await firebaseHelper.registerWithEmail(
        data.email,
        data.password!,
        data.name
      );

      // Return result indicating email was sent
      return {
        emailSent: true,
        email: result.email,
        name: result.name,
        message:
          "Verification email sent. Please check your inbox and verify your email before logging in.",
      };
    } catch (error: any) {
      throw new Error(error.message || "Registration failed");
    }
  }

  // Complete registration after email verification
  async completeRegistration(email: string, password: string) {
    try {
      // Step 1: Check if email is now verified
      const isVerified = await firebaseHelper.checkEmailVerification(
        email,
        password
      );

      if (!isVerified) {
        throw new Error(
          "Email is not yet verified. Please check your inbox and click the verification link."
        );
      }

      // Step 2: Login with Firebase to get verified user data
      const firebaseData = await firebaseHelper.loginWithEmail(email, password);

      // Step 3: Send verified user data to backend
      const registerData: RegisterRequest = {
        firebase_id: firebaseData.firebase_uid,
        firebase_token: firebaseData.firebase_token,
        platform: "email",
      };

      const response = await this.register(registerData);

      // Store tokens
      localStorage.setItem("auth_token", response.access_token);
      localStorage.setItem("refresh_token", response.refresh_token);
      localStorage.setItem("user", JSON.stringify(response.user));

      return response;
    } catch (error: any) {
      throw new Error(error.message || "Registration completion failed");
    }
  }

  // Resend verification email
  async resendVerificationEmail(email: string, password: string) {
    try {
      await firebaseHelper.resendEmailVerification(email, password);
      return {
        success: true,
        message: "Verification email sent. Please check your inbox.",
      };
    } catch (error: any) {
      throw new Error(error.message || "Failed to send verification email");
    }
  }

  // Check verification status
  async checkEmailVerificationStatus(email: string, password: string) {
    try {
      const isVerified = await firebaseHelper.checkEmailVerification(
        email,
        password
      );
      return { isVerified };
    } catch (error: any) {
      throw new Error(error.message || "Failed to check verification status");
    }
  }

  // Social Auth - ALREADY VERIFIED
  async googleAuth() {
    try {
      const firebaseData = await firebaseHelper.loginWithGoogle();

      const loginData: LoginRequest = {
        firebase_id: firebaseData.firebase_uid,
        firebase_token: firebaseData.firebase_token,
        platform: "google",
      };

      const response = await this.login(loginData);

      localStorage.setItem("auth_token", response.access_token);
      localStorage.setItem("refresh_token", response.refresh_token);
      localStorage.setItem("user", JSON.stringify(response.user));

      return response;
    } catch (error: any) {
      throw new Error(error.message || "Google authentication failed");
    }
  }

  // Facebook Auth - ALREADY VERIFIED
  async facebookAuth() {
    try {
      const firebaseData = await firebaseHelper.loginWithFacebook();

      const loginData: LoginRequest = {
        firebase_id: firebaseData.firebase_uid,
        firebase_token: firebaseData.firebase_token,
        platform: "facebook",
      };

      const response = await this.login(loginData);

      localStorage.setItem("auth_token", response.access_token);
      localStorage.setItem("refresh_token", response.refresh_token);
      localStorage.setItem("user", JSON.stringify(response.user));

      return response;
    } catch (error: any) {
      throw new Error(error.message || "Facebook authentication failed");
    }
  }

  // Firebase Password Reset
  async firebasePasswordReset(email: string): Promise<void> {
    try {
      await firebaseHelper.resetPassword(email);
    } catch (error: any) {
      throw new Error(error.message || "Password reset failed");
    }
  }

  // Firebase Logout
  async logoutUser() {
    try {
      await this.logout();
    } catch (error) {
      console.error("Backend logout error:", error);
    }

    try {
      await firebaseHelper.logout();
    } catch (error) {
      console.error("Firebase logout error:", error);
    } finally {
      localStorage.clear();
    }
  }

  isAuthenticated(): boolean {
    return !!localStorage.getItem("auth_token");
  }
}

export const authService = new AuthService();
