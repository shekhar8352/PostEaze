
import AppMantineProvider from "./providers/MantineProvider";
import StoreProvider from "./providers/StoreProvider";
import TanstackProvider from "./providers/TanstackProvider";
import AppRoutes from "./routes";


export default function App() {
  return (
    <StoreProvider>
      <AppMantineProvider>
        <TanstackProvider>
          <AppRoutes />
        </TanstackProvider>
      </AppMantineProvider>
    </StoreProvider>
  );
}
