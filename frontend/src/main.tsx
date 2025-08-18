import React from "react";
import ReactDOM from "react-dom/client";
import { Provider } from "react-redux";
import { store } from "./app/store/store";
import { BrowserRouter } from "react-router-dom";
import App from "./app/App";
import "./App.css";
import "@mantine/core/styles.css";
import { MantineProvider, createTheme } from "@mantine/core";



ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <MantineProvider theme={theme}>
      <Provider store={store}>
        <BrowserRouter>
          <App />
        </BrowserRouter>
      </Provider>
    </MantineProvider>
  </React.StrictMode>
);
