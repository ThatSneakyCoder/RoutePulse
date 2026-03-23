import { CalendarDays, Plus, Route } from "lucide-react";
import { Link, useRouteLoaderData } from "react-router-dom";
import { CreatedTripsList } from "./CreatedTripsList.jsx";

export const Trips = () => {
  const { trips } = useRouteLoaderData("trip");
  const stats = [
    {
      label: "Created trips",
      value: String(trips.length),
      icon: Route,
    },
    {
      label: "Pending start",
      value: String(trips.filter((trip) => trip.status === "created").length),
      icon: CalendarDays,
    },
  ];

  return (
    <section className="p-8 space-y-6">
      <header className="rounded-3xl border border-slate-800 bg-slate-900 p-8">
        <div className="flex flex-col gap-6 lg:flex-row lg:items-end lg:justify-between">
          <div className="max-w-2xl">
            <p className="text-xs font-semibold uppercase tracking-[0.35em] text-cyan-400/80">
              Trips
            </p>
            <h1 className="mt-3 text-3xl font-semibold text-white">
              Create and start trips from one place.
            </h1>
            <p className="mt-3 text-sm leading-6 text-slate-400">
              Dispatch new trips with the right organization, driver, vehicle,
              and route coordinates without leaving the dashboard.
            </p>
          </div>

          <Link
            to="driver-console"
            className="inline-flex items-center justify-center gap-2 rounded-2xl border border-slate-700 bg-slate-950 px-4 py-3 text-sm font-medium text-slate-200 transition hover:border-slate-600 hover:text-white"
          >
            Driver Console
          </Link>

          <Link
            to="create"
            className="inline-flex items-center justify-center gap-2 rounded-2xl border border-cyan-500/30 bg-cyan-500/10 px-4 py-3 text-sm font-medium text-cyan-200 transition hover:border-cyan-400/50 hover:bg-cyan-500/15"
          >
            <Plus className="h-4 w-4" />
            Create Trip
          </Link>
        </div>
      </header>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        {stats.map((item) => {
          const Icon = item.icon;

          return (
            <div
              key={item.label}
              className="rounded-2xl border border-slate-800 bg-slate-900 p-5"
            >
              <div className="flex items-center justify-between">
                <p className="text-sm text-slate-400">{item.label}</p>
                <Icon className="h-4 w-4 text-cyan-400" />
              </div>
              <p className="mt-3 text-3xl font-semibold text-white">{item.value}</p>
            </div>
          );
        })}
      </div>

      <CreatedTripsList />
    </section>
  );
};
