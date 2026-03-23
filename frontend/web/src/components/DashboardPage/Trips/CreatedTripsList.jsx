import { Building2, CalendarDays, Truck, UserRound } from "lucide-react";
import { useRouteLoaderData } from "react-router-dom";

export const CreatedTripsList = () => {
  const { trips, organizations, drivers, vehicles } = useRouteLoaderData("trip");

  const organizationNames = Object.fromEntries(
    organizations.map((organization) => [
      organization.organization_id,
      organization.name,
    ]),
  );
  const driverNames = Object.fromEntries(
    drivers.map((driver) => [
      driver.driver_id,
      `${driver.first_name} ${driver.last_name}`,
    ]),
  );
  const vehicleNames = Object.fromEntries(
    vehicles.map((vehicle) => [vehicle.vehicle_id, vehicle.plate_number]),
  );

  return (
    <section className="rounded-3xl border border-slate-800 bg-slate-900 p-6">
      <div className="flex items-center justify-between gap-4">
        <div>
          <p className="text-xs font-semibold uppercase tracking-[0.3em] text-slate-500">
            Created Trips
          </p>
          <h2 className="mt-2 text-xl font-semibold text-white">
            Ready to dispatch
          </h2>
        </div>
        <p className="text-sm text-slate-400">{trips.length} total</p>
      </div>

      <div className="mt-6 space-y-4">
        {trips.length === 0 ? (
          <div className="rounded-2xl border border-dashed border-slate-700 bg-slate-950/50 p-8 text-center">
            <p className="text-base font-medium text-white">No trips created yet.</p>
            <p className="mt-2 text-sm text-slate-400">
              Use the create trip flow above to dispatch your first trip.
            </p>
          </div>
        ) : (
          trips.map((trip) => (
            <article
              key={trip.trip_id}
              className="rounded-2xl border border-slate-800 bg-slate-950/70 p-5"
            >
              <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
                <div className="space-y-3">
                  <div className="flex flex-wrap items-center gap-3">
                    <h3 className="text-lg font-semibold text-white">{trip.trip_id}</h3>
                    <StatusBadge status={trip.status} />
                  </div>

                  <div className="grid grid-cols-1 gap-2 text-sm text-slate-400 sm:grid-cols-2">
                    <InfoRow
                      icon={Building2}
                      value={
                        organizationNames[trip.organization_id] ?? trip.organization_id
                      }
                    />
                    <InfoRow
                      icon={UserRound}
                      value={driverNames[trip.driver_id] ?? trip.driver_id}
                    />
                    <InfoRow
                      icon={Truck}
                      value={vehicleNames[trip.vehicle_id] ?? trip.vehicle_id}
                    />
                    <InfoRow
                      icon={CalendarDays}
                      value={new Date(trip.created_at * 1000).toLocaleString()}
                    />
                  </div>
                </div>
              </div>
            </article>
          ))
        )}
      </div>
    </section>
  );
};

const InfoRow = ({ icon: Icon, value }) => {
  return (
    <div className="flex items-center gap-2">
      <Icon className="h-4 w-4 text-slate-500" />
      <span className="truncate">{value}</span>
    </div>
  );
};

const StatusBadge = ({ status }) => {
  const badgeClass =
    status === "completed"
      ? "bg-emerald-500/15 text-emerald-300"
      : status === "active"
        ? "bg-cyan-500/15 text-cyan-300"
        : "bg-amber-500/15 text-amber-300";

  return (
    <span className={`rounded-full px-3 py-1 text-xs font-medium ${badgeClass}`}>
      {status}
    </span>
  );
};
