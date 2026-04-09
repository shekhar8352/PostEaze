import { ColorSchemeScript, MantineProvider } from "@mantine/core";
import { DatesProvider } from "@mantine/dates";
import type { ReactNode } from "react";
import { theme } from "../theme";

const AppMantineProvider = ({ children }: { children: ReactNode }) => {
  return (
    <>
      <ColorSchemeScript defaultColorScheme="auto" />
      <MantineProvider theme={theme} defaultColorScheme="auto">
        <DatesProvider settings={{ locale: "en", firstDayOfWeek: 0, consistentWeeks: true }}>
          {children}
        </DatesProvider>
      </MantineProvider>
    </>
  );
};

export default AppMantineProvider;
