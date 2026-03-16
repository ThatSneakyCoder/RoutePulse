import {
  Building2,
  CalendarClock,
  Truck,
  UserRound,
} from "lucide-react";
import { useRouteLoaderData } from "react-router-dom";

export const ActiveShipments = () => {
  const { activeShipments } = useRouteLoaderData("shipments");

  if (!activeShipments || activeShipments.length === 0) {
    return (
      <section className="p-8">
        <div className="rounded-xl border border-slate-800 bg-slate-900 p-8 text-slate-400">
          No active shipments found
        </div>
      </section>
    );
  }

  return (
    <section className="p-8 space-y-8">
      <div>
        <h1 className="text-2xl font-semibold text-white">Active Shipments</h1>
        <p className="mt-2 text-sm text-slate-400">
          Track shipments that are newly created or currently in progress.
        </p>
      </div>

      <div className="grid grid-cols-1 gap-6 md:grid-cols-2 xl:grid-cols-3">
        {activeShipments.map((shipment) => (
          <ShipmentCard key={shipment.trip_id} shipment={shipment} />
        ))}
      </div>
    </section>
  );
};

const ShipmentCard = ({ shipment }) => {
  const isCreated = shipment.status === "created";

  return (
    <article className="rounded-xl border border-slate-700 bg-slate-800 p-6 shadow-lg transition duration-200 hover:border-blue-500 hover:shadow-blue-500/10">
      <div className="mb-5 flex items-start justify-between gap-4">
        <div>
          <p className="text-xs uppercase tracking-[0.2em] text-slate-400">
            Shipment
          </p>
          <h2 className="mt-2 text-lg font-semibold text-white">
            {shipment.trip_id}
          </h2>
        </div>

        <span
          className={`rounded-full px-3 py-1 text-xs font-medium ${
            isCreated
              ? "bg-amber-500/20 text-amber-300"
              : "bg-green-500/20 text-green-400"
          }`}
        >
          {isCreated ? "Created" : "Active"}
        </span>
      </div>

      <div className="mb-5 grid gap-3 text-sm text-slate-300">
        <InfoRow
          icon={Truck}
          label="Vehicle"
          value={shipment.vehicle_id}
          mono
        />
        <InfoRow
          icon={UserRound}
          label="Driver"
          value={shipment.driver_id}
          mono
        />
        <InfoRow
          icon={Building2}
          label="Organization"
          value={shipment.organization_id}
          mono
        />
        <InfoRow
          icon={CalendarClock}
          label="Created"
          value={formatDate(shipment.created_at)}
        />
      </div>

      <div className="border-t border-slate-700 pt-4">
        <div className="flex items-center justify-between text-xs text-slate-400">
          <span>Current status</span>
          <span className="font-medium uppercase tracking-wide text-slate-200">
            {shipment.status}
          </span>
        </div>
      </div>
    </article>
  );
};

const InfoRow = ({ icon: Icon, label, value, mono = false }) => (
  <div className="flex items-center justify-between gap-3">
    <span className="flex items-center gap-2 text-slate-400">
      <Icon className="h-4 w-4" />
      {label}
    </span>
    <span className={`max-w-[12rem] truncate font-medium ${mono ? "font-mono text-xs" : ""}`}>
      {value || "Unavailable"}
    </span>
  </div>
);

const formatDate = (value) => {
  if (!value) return "Unavailable";

  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return date.toLocaleString();
};
