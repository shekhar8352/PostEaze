import { Provider } from "react-redux";
import { QueryClientProvider } from "@tanstack/react-query";
import { MantineProvider, createTheme } from "@mantine/core";
import { store } from "./store/store";
import { queryClient } from "./queryClient";


const theme = createTheme({
  primaryColor: "indigo",
  fontFamily: "Inter, sans-serif",
  headings: { fontFamily: "Inter, sans-serif" },
  // Add more theme customizations here
});


export const AppProviders = ({ children }: { children: React.ReactNode }) => (
  <Provider store={store}>
    <QueryClientProvider client={queryClient}>
      <MantineProvider theme={theme}>
        {children}
      </MantineProvider>
    </QueryClientProvider>
  </Provider>
);
