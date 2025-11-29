import { MantineProvider } from "@mantine/core";
import type { ReactNode } from "react";
import { theme } from "../theme";

const AppMantineProvider = ({ children }: { children: ReactNode }) => {
  return (
    <MantineProvider theme={theme}>
      {children}
    </MantineProvider>
  );
};

export default AppMantineProvider;
