import {
  Building2,
  CalendarCheck2,
  CircleCheckBig,
  Truck,
  UserRound,
} from "lucide-react";
import { useRouteLoaderData } from "react-router-dom";

export const CompletedShipments = () => {
  const { completedShipments } = useRouteLoaderData("shipments");

  if (!completedShipments || completedShipments.length === 0) {
    return (
      <section className="p-8">
        <div className="rounded-xl border border-slate-800 bg-slate-900 p-8 text-slate-400">
          No completed shipments found
        </div>
      </section>
    );
  }

  return (
    <section className="p-8 space-y-8">
      <div>
        <h1 className="text-2xl font-semibold text-white">
          Completed Shipments
        </h1>
        <p className="mt-2 text-sm text-slate-400">
          Review deliveries that have already reached their final status.
        </p>
      </div>

      <div className="grid grid-cols-1 gap-6 md:grid-cols-2 xl:grid-cols-3">
        {completedShipments.map((shipment) => (
          <CompletedShipmentCard key={shipment.trip_id} shipment={shipment} />
        ))}
      </div>
    </section>
  );
};

const CompletedShipmentCard = ({ shipment }) => (
  <article className="flex h-full min-w-0 flex-col rounded-2xl border border-slate-700/80 bg-slate-800/95 p-5 shadow-lg transition duration-200 hover:border-blue-500 hover:shadow-blue-500/10">
    <div className="mb-5 flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
      <div className="min-w-0 flex-1">
        <p className="text-[11px] uppercase tracking-[0.32em] text-slate-400">
          Delivered Shipment
        </p>
        <h2 className="mt-2.5 break-all text-2xl font-semibold leading-tight text-white sm:text-[1.65rem]">
          {shipment.trip_id}
        </h2>
      </div>

      <span className="inline-flex shrink-0 items-center gap-2 self-start rounded-full bg-emerald-500/15 px-3 py-1.5 text-xs font-medium text-emerald-300 ring-1 ring-emerald-500/20">
        <CircleCheckBig className="h-3.5 w-3.5" />
        Completed
      </span>
    </div>

    <div className="grid gap-2.5">
      <InfoRow icon={Truck} label="Vehicle" value={shipment.vehicle_id} mono />
      <InfoRow icon={UserRound} label="Driver" value={shipment.driver_id} mono />
      <InfoRow
        icon={Building2}
        label="Organization"
        value={shipment.organization_id}
        mono
      />
      <InfoRow
        icon={CalendarCheck2}
        label="Created"
        value={formatDate(shipment.created_at)}
      />
    </div>

    <div className="mt-5 border-t border-slate-700/80 pt-3.5">
      <div className="flex items-center justify-between gap-4 text-xs text-slate-400">
        <span>Final status</span>
        <span className="rounded-full bg-emerald-500/10 px-2.5 py-1 font-medium uppercase tracking-[0.2em] text-emerald-300">
          {shipment.status}
        </span>
      </div>
    </div>
  </article>
);

const InfoRow = ({ icon: Icon, label, value, mono = false }) => (
  <div className="grid min-w-0 grid-cols-[auto_minmax(0,1fr)] items-center gap-3 rounded-xl border border-slate-700/60 bg-slate-900/30 px-4 py-2.5">
    <span className="flex min-w-0 items-center gap-2 text-slate-400">
      <Icon className="h-4 w-4 shrink-0" />
      <span className="w-24 truncate">{label}</span>
    </span>
    <span
      className={`min-w-0 text-right font-medium text-slate-200 ${
        mono ? "truncate font-mono text-sm leading-normal" : "truncate"
      }`}
      title={value || "Unavailable"}
    >
      {value || "Unavailable"}
    </span>
  </div>
);

const formatDate = (value) => {
  if (!value) return "Unavailable";

  const numericValue = Number(value);
  const parsedValue = Number.isNaN(numericValue)
    ? value
    : numericValue < 1_000_000_000_000
      ? numericValue * 1000
      : numericValue;

  const date = new Date(parsedValue);

  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return date.toLocaleString();
};
