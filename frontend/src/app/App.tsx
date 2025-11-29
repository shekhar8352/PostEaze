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
import "./styles/global.css";

export default function App() {
  return (
    <React.StrictMode>
      <StoreProvider>
        <AuthProvider>
          <AppMantineProvider>
            <Notifications position="top-right" limit={5} />
            <TanstackProvider>
              <BrowserRouter>
                <AppRoutes />
              </BrowserRouter>
            </TanstackProvider>
          </AppMantineProvider>
        </AuthProvider>
      </StoreProvider>
    </React.StrictMode>
  );
}
