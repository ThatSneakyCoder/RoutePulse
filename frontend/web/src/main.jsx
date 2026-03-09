import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";

import App from "./App.jsx";
import { Dashboard } from "./components/DashboardPage/Dashboard/Dashboard.jsx";
import { Organization } from "./components/DashboardPage/Organizations/Organization.jsx";
import { Profile } from "./components/DashboardPage/Profile/Profile.jsx";
import { RootDashboard } from "./components/DashboardPage/RootDashboard.jsx";

import { createBrowserRouter, RouterProvider } from "react-router-dom";
import AuthLayout from "./components/Auth/AuthLayout.jsx";
import { Login } from "./components/Auth/Login.jsx";
import { Register } from "./components/Auth/Register.jsx";
import { HealthCheck } from "./components/HealthCheck.jsx";
import { NotFoundPage } from "./components/NotFoundPage.jsx";
import { Home } from "./components/landingPage/Home.jsx";

import { organizationMutationAction } from "./actions/createOrganizationAction.js";
import { OrganizationDetails } from "./components/DashboardPage/Organizations/OrganizationDetails.jsx";
import { DashboardError } from "./errors/DashboardError";
import { dashboardLoader } from "./loaders/DashboardLoader.js";

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
        id: "dashboard",
        path: "dashboard",
        Component: RootDashboard,
        loader: dashboardLoader,
        errorElement: <DashboardError />,
        children: [
          {
            index: true,
            Component: Dashboard,
          },
          { path: "profile", Component: Profile },
          {
            path: "organization",
            children: [
              {
                index: true,
                Component: Organization,
                action: organizationMutationAction,
              },
              { path: ":orgId", Component: OrganizationDetails },
            ],
          },

          // { path: "fleet/vehicles", Component: AllVehicles },
          // { path: "fleet/active", Component: ActiveFleet },
          // { path: "fleet/maintenance", Component: Maintenance },
          // { path: "fleet/telematics", Component: Telematics },
          // { path: "fleet/driversDirectory", Component: DriversDirectory },

          // { path: "shipments/active", Component: ActiveShipments },
          // { path: "shipments/complete", Component: CompletedShipments },

          // { path: "warehouse", Component: Warehouse },
          // { path: "routes", Component: RoutePlanning },
          // { path: "analytics", Component: Analytics },

          // { path: "customers", Component: Customers },
          // { path: "teams", Component: Teams },
          // { path: "help", Component: Help },
        ],
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
