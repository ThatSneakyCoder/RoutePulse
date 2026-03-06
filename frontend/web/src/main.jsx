import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";

import App from "./App.jsx";
import { Dashboard } from "./components/DashboardPage/Dashboard/Dashboard.jsx";

import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { Login } from "./components/Auth/Login.jsx";
import { Register } from "./components/Auth/Register.jsx";
import { NotFoundPage } from "./components/NotFoundPage.jsx";
import { Home } from "./components/landingPage/Home.jsx";
import AuthLayout from "./components/Auth/AuthLayout.jsx";
import { HealthCheck } from "./components/HealthCheck.jsx";
import { RootDashboard } from "./components/DashboardPage/RootDashboard.jsx";

const router = createBrowserRouter([
  {
    path: "/",
    Component: App,
    children: [
      { index: true, Component: Home },
      {
        path: "auth",
        Component: AuthLayout,
        children: [
          { path: "login", Component: Login },
          { path: "register", Component: Register },
        ],
      },
      {
        path: "dashboard",
        Component: RootDashboard,
        children: [{ index: true, Component: Dashboard }],
      },
      { path: "*", Component: NotFoundPage },
    ],
  },
  {
    path: "/health",
    Component: HealthCheck,
  },
]);

createRoot(document.getElementById("root")).render(
  <StrictMode>
    <RouterProvider router={router} />
  </StrictMode>,
);
