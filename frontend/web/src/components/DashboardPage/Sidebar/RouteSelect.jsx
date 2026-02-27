import {
  FiHome,
  FiTruck,
  FiChevronDown,
  FiChevronRight,
  FiUsers,
  FiHelpCircle,
  FiPackage,
  FiBox,
  FiMap,
  FiBarChart2,
  FiTool,
  FiDroplet,
  FiUserCheck,
} from "react-icons/fi";

import { useState } from "react";

export const RouteSelect = () => {
  const [selected, setSelected] = useState("Dashboard");
  const [fleetOpen, setFleetOpen] = useState(true);
  const [shipmentOpen, setShipmentOpen] = useState(false);

  return (
    <div className="space-y-8">
      {/* MAIN MENU */}
      <div>
        <p className="text-xs font-semibold text-slate-500 mb-2 px-2 tracking-wider">
          MAIN MENU
        </p>

        <Route
          Icon={FiHome}
          title="Dashboard"
          selected={selected === "Dashboard"}
          onClick={() => setSelected("Dashboard")}
        />

        {/* Fleet Management */}
        <CollapsibleRoute
          Icon={FiTruck}
          title="Fleet Management"
          open={fleetOpen}
          toggle={() => setFleetOpen(!fleetOpen)}
        >
          <SubRoute
            title="All Vehicles"
            selected={selected === "All Vehicles"}
            onClick={() => setSelected("All Vehicles")}
          />
          <SubRoute
            title="Active Trucks"
            selected={selected === "Active Trucks"}
            onClick={() => setSelected("Active Trucks")}
          />
          <SubRoute
            title="Maintenance & Repairs"
            badge="2"
            selected={selected === "Maintenance & Repairs"}
            onClick={() => setSelected("Maintenance & Repairs")}
          />
          <SubRoute
            title="Fuel & Telematics"
            selected={selected === "Fuel & Telematics"}
            onClick={() => setSelected("Fuel & Telematics")}
          />
          <SubRoute
            title="Drivers Directory"
            selected={selected === "Drivers Directory"}
            onClick={() => setSelected("Drivers Directory")}
          />
        </CollapsibleRoute>

        {/* Shipments */}
        <CollapsibleRoute
          Icon={FiPackage}
          title="Shipments & Deliveries"
          open={shipmentOpen}
          toggle={() => setShipmentOpen(!shipmentOpen)}
        >
          <SubRoute
            title="Active Shipments"
            selected={selected === "Active Shipments"}
            onClick={() => setSelected("Active Shipments")}
          />
          <SubRoute
            title="Completed Shipments"
            selected={selected === "Completed Shipments"}
            onClick={() => setSelected("Completed Shipments")}
          />
        </CollapsibleRoute>

        <Route
          Icon={FiBox}
          title="Warehouse Operations"
          selected={selected === "Warehouse Operations"}
          onClick={() => setSelected("Warehouse Operations")}
        />

        <Route
          Icon={FiMap}
          title="Route Planning"
          badge="3"
          selected={selected === "Route Planning"}
          onClick={() => setSelected("Route Planning")}
        />

        <Route
          Icon={FiBarChart2}
          title="Analytics & Reports"
          selected={selected === "Analytics & Reports"}
          onClick={() => setSelected("Analytics & Reports")}
        />
      </div>

      {/* SUPPORT */}
      <div>
        <p className="text-xs font-semibold text-slate-500 mb-2 px-2 tracking-wider">
          SUPPORT
        </p>

        <Route
          Icon={FiUsers}
          title="Customers"
          selected={selected === "Customers"}
          onClick={() => setSelected("Customers")}
        />

        <Route
          Icon={FiUserCheck}
          title="Team & Roles"
          selected={selected === "Team & Roles"}
          onClick={() => setSelected("Team & Roles")}
        />

        <Route
          Icon={FiHelpCircle}
          title="Help & Support"
          selected={selected === "Help & Support"}
          onClick={() => setSelected("Help & Support")}
        />
      </div>
    </div>
  );
};

const Route = ({ selected, Icon, title, onClick }) => {
  return (
    <button
      onClick={onClick}
      className={`flex items-center gap-3 w-full px-3 py-2 rounded-lg text-sm transition
      ${
        selected
          ? "bg-slate-800 text-white font-medium"
          : "text-slate-400 hover:bg-slate-900 hover:text-white"
      }`}
    >
      <Icon className="text-base" />
      {title}
    </button>
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

const SubRoute = ({ selected, title, onClick }) => {
  return (
    <button
      onClick={onClick}
      className={`block w-full text-left px-2 py-1.5 rounded-md text-sm transition
      ${
        selected ? "text-white bg-slate-800" : "text-slate-500 hover:text-white"
      }`}
    >
      {title}
    </button>
  );
};
