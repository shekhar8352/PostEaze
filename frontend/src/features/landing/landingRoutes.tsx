import type { RouteObject } from "react-router-dom";
import LandingPage from "./pages/Landing Page";

const landingRoutes: RouteObject[] = [
  {
    path: "/",
    element: <LandingPage />,
  },
];

export default landingRoutes;
