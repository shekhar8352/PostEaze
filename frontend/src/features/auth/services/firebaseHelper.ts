// src/features/auth/services/firebaseHelper.ts
import {
  createUserWithEmailAndPassword,
  signInWithEmailAndPassword,
  signInWithPopup,
  GoogleAuthProvider,
  FacebookAuthProvider,
  sendEmailVerification,
  sendPasswordResetEmail,
  signOut,
  updateProfile,
  reload,
} from "firebase/auth";
import { auth } from "./firebase.config";

const googleProvider = new GoogleAuthProvider();
const facebookProvider = new FacebookAuthProvider();

googleProvider.setCustomParameters({
  prompt: "select_account",
});

class FirebaseHelper {
  // Email/Password Login - ONLY ALLOWS VERIFIED USERS
  async loginWithEmail(email: string, password: string) {
    try {
      const userCredential = await signInWithEmailAndPassword(
        auth,
        email,
        password
      );
      const user = userCredential.user;

      // CRITICAL: Check if email is verified
      if (!user.emailVerified) {
        // Sign out the user immediately
        await signOut(auth);
        throw new Error("EMAIL_NOT_VERIFIED");
      }

      const token = await user.getIdToken();

      return {
        firebase_uid: user.uid,
        firebase_token: token,
        display_name: user.displayName || "",
        email_verified: true as const,
        provider: "email" as const,
      };
    } catch (error: any) {
      if (error.message === "EMAIL_NOT_VERIFIED") {
        throw new Error(
          "Please verify your email before logging in. Check your inbox for the verification link."
        );
      }
      throw new Error(this.mapError(error));
    }
  }

  // Email/Password Registration - SENDS VERIFICATION EMAIL
  async registerWithEmail(email: string, password: string, name: string) {
    try {
      const userCredential = await createUserWithEmailAndPassword(
        auth,
        email,
        password
      );
      const user = userCredential.user;

      // Update profile with display name
      await updateProfile(user, { displayName: name });

      // Send email verification IMMEDIATELY
      await sendEmailVerification(user);

      // Sign out the user - they can't proceed until verified
      await signOut(auth);

      return {
        email: user.email!,
        name: name,
        emailSent: true,
      };
    } catch (error: any) {
      throw new Error(this.mapError(error));
    }
  }

  // Google Login - GOOGLE EMAILS ARE PRE-VERIFIED
  async loginWithGoogle() {
    try {
      const userCredential = await signInWithPopup(auth, googleProvider);
      const user = userCredential.user;
      const token = await user.getIdToken();

      return {
        firebase_uid: user.uid,
        firebase_token: token,
        display_name: user.displayName || "",
        email_verified: true as const, // Google emails are pre-verified
        provider: "google.com" as const,
      };
    } catch (error: any) {
      throw new Error(this.mapError(error));
    }
  }

  // Facebook Login - FACEBOOK EMAILS ARE PRE-VERIFIED
  async loginWithFacebook() {
    try {
      const userCredential = await signInWithPopup(auth, facebookProvider);
      const user = userCredential.user;
      const token = await user.getIdToken();

      return {
        firebase_uid: user.uid,
        firebase_token: token,
        display_name: user.displayName || "",
        email_verified: true as const, // Facebook emails are pre-verified
        provider: "facebook.com" as const,
      };
    } catch (error: any) {
      throw new Error(this.mapError(error));
    }
  }

  // Resend Email Verification
  async resendEmailVerification(
    email: string,
    password: string
  ): Promise<void> {
    try {
      // Sign in to get the user
      const userCredential = await signInWithEmailAndPassword(
        auth,
        email,
        password
      );
      const user = userCredential.user;

      if (user.emailVerified) {
        await signOut(auth);
        throw new Error("Email is already verified. You can now log in.");
      }

      // Send verification email
      await sendEmailVerification(user);

      // Sign out immediately
      await signOut(auth);
    } catch (error: any) {
      if (error.message.includes("already verified")) {
        throw error;
      }
      throw new Error(this.mapError(error));
    }
  }

  // Check email verification status
  async checkEmailVerification(
    email: string,
    password: string
  ): Promise<boolean> {
    try {
      const userCredential = await signInWithEmailAndPassword(
        auth,
        email,
        password
      );
      const user = userCredential.user;

      // Reload user to get latest verification status
      await reload(user);

      const isVerified = user.emailVerified;

      // Sign out immediately
      await signOut(auth);

      return isVerified;
    } catch (error: any) {
      throw new Error(this.mapError(error));
    }
  }

  // Password Reset
  async resetPassword(email: string): Promise<void> {
    try {
      await sendPasswordResetEmail(auth, email);
    } catch (error: any) {
      throw new Error(this.mapError(error));
    }
  }

  // Logout
  async logout(): Promise<void> {
    try {
      await signOut(auth);
    } catch (error: any) {
      throw new Error(this.mapError(error));
    }
  }

  private mapError(error: any): string {
    switch (error.code) {
      case "auth/user-not-found":
        return "No account found with this email address";
      case "auth/wrong-password":
        return "Invalid password. Please try again.";
      case "auth/email-already-in-use":
        return "An account with this email already exists";
      case "auth/weak-password":
        return "Password must be at least 6 characters";
      case "auth/invalid-email":
        return "Please enter a valid email address";
      case "auth/too-many-requests":
        return "Too many failed attempts. Please try again later";
      case "auth/popup-closed-by-user":
        return "Sign-in was cancelled by user";
      case "auth/cancelled-popup-request":
        return "Sign-in was cancelled";
      default:
        return error.message || "Authentication failed. Please try again.";
    }
  }
}

export const firebaseHelper = new FirebaseHelper();
