import type { RouteObject } from "react-router-dom";
import SignupPage from "./pages/SignUp";

const authRoutes: RouteObject[] = [
  { path: "/signup", element: <SignupPage /> },
];

export default authRoutes;
