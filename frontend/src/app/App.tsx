import AppMantineProvider from "./providers/MantineProvider";
import StoreProvider from "./providers/StoreProvider";
import TanstackProvider from "./providers/TanstackProvider";
import AppRoutes from "./routes";
import React from "react";
import { BrowserRouter } from "react-router-dom";
import "@mantine/core/styles.css";

export default function App() {
  return (
    <React.StrictMode>
      <StoreProvider>
        <AppMantineProvider>
          <TanstackProvider>
            <BrowserRouter>
              <AppRoutes />
            </BrowserRouter>
          </TanstackProvider>
        </AppMantineProvider>
      </StoreProvider>
    </React.StrictMode>
  );
}
