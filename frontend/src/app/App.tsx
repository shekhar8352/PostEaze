import AppMantineProvider from "./providers/MantineProvider";
import StoreProvider from "./providers/StoreProvider";
import TanstackProvider from "./providers/TanstackProvider";
import AppRoutes from "./routes";
import React from "react";
import { BrowserRouter } from "react-router-dom";
import "@mantine/core/styles.css";
import { Notifications } from '@mantine/notifications';
import "@mantine/notifications/styles.css";
import { AuthProvider } from "@/features/auth";
import "./shell/tokens.css";
import "./styles/global.css";
// Initialize API interceptors for auth token handling
import "@/services/api/interceptors";

export default function App() {
  return (
    <React.StrictMode>
      <StoreProvider>
        <TanstackProvider>
          <AuthProvider>
            <AppMantineProvider>
              <Notifications position="top-right" limit={5} />
              <BrowserRouter>
                <AppRoutes />
              </BrowserRouter>
            </AppMantineProvider>
          </AuthProvider>
        </TanstackProvider>
      </StoreProvider>
    </React.StrictMode>
  );
}
