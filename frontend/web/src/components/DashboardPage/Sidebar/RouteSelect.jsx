import {
  FiHome,
  FiTruck,
  FiChevronDown,
  FiChevronRight,
  FiUsers,
  FiHelpCircle,
  FiPackage,
  FiBarChart2,
  FiUserCheck,
  FiUser,
} from "react-icons/fi";

import { GiPathDistance } from "react-icons/gi";

import { useState } from "react";
import { NavLink } from "react-router-dom";

export const RouteSelect = () => {
  const [fleetOpen, setFleetOpen] = useState(true);
  const [shipmentOpen, setShipmentOpen] = useState(false);
  const [tripOpen, setTripOpen] = useState(true);

  return (
    <div className="space-y-8">
      {/* MAIN MENU */}
      <div>
        <p className="text-xs font-semibold text-slate-500 mb-2 px-2 tracking-wider">
          MAIN MENU
        </p>

        <Route Icon={FiHome} title="Dashboard" to="/dashboard" end />

        <Route Icon={FiUser} title="My Organization" to="/dashboard/organization" />

        <CollapsibleRoute
          Icon={GiPathDistance}
          title="Trip"
          open={tripOpen}
          toggle={() => setTripOpen(!tripOpen)}
        >
          <SubRoute title="All Trips" to="/dashboard/trip" end />
          <SubRoute title="Create Trip" to="/dashboard/trip/create" />
          <SubRoute title="Driver Console" to="/dashboard/trip/driver-console" />
        </CollapsibleRoute>

        {/* Fleet Management */}
        <CollapsibleRoute
          Icon={FiTruck}
          title="Fleet Management"
          open={fleetOpen}
          toggle={() => setFleetOpen(!fleetOpen)}
        >
          <SubRoute title="All Vehicles" to="/dashboard/fleet/vehicles/all" />
          <SubRoute title="Drivers" to="/dashboard/fleet/drivers/all" />
        </CollapsibleRoute>

        {/* Shipments */}
        <CollapsibleRoute
          Icon={FiPackage}
          title="Shipments & Deliveries"
          open={shipmentOpen}
          toggle={() => setShipmentOpen(!shipmentOpen)}
        >
          <SubRoute title="Active Shipments" to="shipments/active" />
          <SubRoute title="Completed Shipments" to="shipments/complete" />
        </CollapsibleRoute>

        <Route Icon={FiBarChart2} title="Analytics & Reports" to="/dashboard/analytics" />
      </div>

      {/* SUPPORT */}
      <div>
        <p className="text-xs font-semibold text-slate-500 mb-2 px-2 tracking-wider">
          SUPPORT
        </p>

        <Route Icon={FiUsers} title="Customers" to="/customers" />

        <Route Icon={FiUserCheck} title="Team & Roles" to="/teams" />

        <Route Icon={FiHelpCircle} title="Help & Support" to="/help" />
      </div>
    </div>
  );
};

const Route = ({ Icon, title, to, end }) => {
  return (
    <NavLink
      to={to}
      end={end}
      className={({ isActive }) =>
        `flex items-center gap-3 w-full px-3 py-2 rounded-lg text-sm transition ${
          isActive
            ? "bg-slate-800 text-white font-medium"
            : "text-slate-400 hover:bg-slate-900 hover:text-white"
        }`
      }
    >
      <Icon className="text-base" />
      {title}
    </NavLink>
  );
};

const CollapsibleRoute = ({ Icon, title, open, toggle, children }) => {
  return (
    <div>
      <button
        onClick={toggle}
        className="flex items-center justify-between w-full px-3 py-2 rounded-lg text-sm text-slate-400 hover:bg-slate-900 hover:text-white transition"
      >
        <div className="flex items-center gap-3">
          <Icon className="text-base" />
          {title}
        </div>

        {open ? <FiChevronDown /> : <FiChevronRight />}
      </button>

      {open && (
        <div className="ml-5 mt-1 space-y-1 border-l border-slate-600 pl-4">
          {children}
        </div>
      )}
    </div>
  );
};

const SubRoute = ({ title, to, end }) => {
  return (
    <NavLink
      to={to}
      end={end}
      className={({ isActive }) =>
        `block w-full text-left px-2 py-1.5 rounded-md text-sm transition ${
          isActive
            ? "text-white bg-slate-800"
            : "text-slate-500 hover:text-white"
        }`
      }
    >
      {title}
    </NavLink>
  );
};
