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
import { organizationDetailsAction } from "./actions/organizationDetailsAction.js";
import { UserAnalytics } from "./components/DashboardPage/Analytics/UserAnalytics.jsx";
import { Vehicles } from "./components/DashboardPage/Fleet/Vehicles.jsx";
import { OrganizationDetails } from "./components/DashboardPage/Organizations/OrganizationDetails.jsx";
import { ActiveShipments } from "./components/DashboardPage/Shipments/ActiveShipments.jsx";
import { CompletedShipments } from "./components/DashboardPage/Shipments/CompletedShipments.jsx";
import { DashboardError } from "./errors/DashboardError";
import { analyticsLoader } from "./loaders/AnalyticsLoader.js";
import { dashboardLoader } from "./loaders/DashboardLoader.js";
import { fleetLoader } from "./loaders/FleetLoader.js";
import { organizationDetailsLoader } from "./loaders/OrganizationDetailsLoader.js";
import { shipmentsLoader } from "./loaders/ShipmentsLoader.js";

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
              {
                path: ":orgId",
                Component: OrganizationDetails,
                loader: organizationDetailsLoader,
                action: organizationDetailsAction,
              },
            ],
          },
          {
            id: "fleet",
            path: "fleet",
            loader: fleetLoader,
            children: [
              {
                path: "vehicles/all",
                Component: Vehicles,
              },
            ],
          },
          {
            id: "shipments",
            path: "shipments",
            loader: shipmentsLoader,
            children: [
              {
                path: "active",
                Component: ActiveShipments,
              },
              {
                path: "complete",
                Component: CompletedShipments,
              },
            ],
          },
          {
            id: "analytics",
            path: "analytics",
            loader: analyticsLoader,
            Component: UserAnalytics,
          },

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
