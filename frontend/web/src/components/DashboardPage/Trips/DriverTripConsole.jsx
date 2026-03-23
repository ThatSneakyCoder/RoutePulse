import { CalendarDays, Route, Truck, UserRound } from "lucide-react";
import { Form, Link, useActionData, useNavigation, useRouteLoaderData } from "react-router-dom";

export const DriverTripConsole = () => {
  const { trips, organizations, drivers, vehicles } = useRouteLoaderData("trip");
  const actionData = useActionData();
  const navigation = useNavigation();
  const isSubmitting = navigation.state === "submitting";

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

  const sortedTrips = [...trips].sort((left, right) => right.created_at - left.created_at);
  const activeTrips = sortedTrips.filter((trip) => trip.status === "active");
  const createdTrips = sortedTrips.filter((trip) => trip.status === "created");

  return (
    <section className="p-8 space-y-8">
      <header className="rounded-3xl border border-slate-800 bg-slate-900 p-8">
        <div className="flex flex-col gap-6 xl:flex-row xl:items-start xl:justify-between">
          <div className="max-w-2xl">
            <p className="text-xs font-semibold uppercase tracking-[0.35em] text-cyan-400/80">
              Driver Console
            </p>
            <h1 className="mt-3 text-3xl font-semibold text-white">
              Start trips and prepare for live tracking.
            </h1>
            <p className="mt-3 text-sm leading-6 text-slate-400">
              This page is the operational handoff for drivers. Start or complete
              trips here now, and the live tracking panel is already structured
              for the WebSocket and tracking APIs you’ll wire next.
            </p>
          </div>

          <div className="flex flex-wrap gap-3">
            <Link
              to="/dashboard/trip"
              className="inline-flex items-center justify-center gap-2 rounded-2xl border border-slate-700 bg-slate-950 px-4 py-3 text-sm font-medium text-slate-200 transition hover:border-slate-600 hover:text-white"
            >
              Back to trips
            </Link>
            <Link
              to="/dashboard/trip/create"
              className="inline-flex items-center justify-center gap-2 rounded-2xl border border-cyan-500/30 bg-cyan-500/10 px-4 py-3 text-sm font-medium text-cyan-200 transition hover:border-cyan-400/50 hover:bg-cyan-500/15"
            >
              Create Trip
            </Link>
          </div>
        </div>
      </header>

      <section className="rounded-3xl border border-slate-800 bg-slate-900 p-6 sm:p-8">
          <div className="flex items-center justify-between gap-4">
            <div>
              <p className="text-xs font-semibold uppercase tracking-[0.3em] text-slate-500">
                Trip Queue
              </p>
              <h2 className="mt-2 text-xl font-semibold text-white">
                Assigned and Active Trips
              </h2>
            </div>
            <p className="text-sm text-slate-400">{sortedTrips.length} total</p>
          </div>

          {actionData?.error ? (
            <p className="mt-6 rounded-2xl border border-rose-500/30 bg-rose-500/10 px-4 py-3 text-sm text-rose-200">
              {actionData.error}
            </p>
          ) : null}

          <div className="mt-6 space-y-4">
            {sortedTrips.length === 0 ? (
              <div className="rounded-2xl border border-dashed border-slate-700 bg-slate-950/50 p-8 text-center">
                <p className="text-base font-medium text-white">No trips available.</p>
                <p className="mt-2 text-sm text-slate-400">
                  Create a trip first, then drivers can start it from here.
                </p>
              </div>
            ) : (
              sortedTrips.map((trip) => (
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

                      <div className="grid gap-2 text-sm text-slate-400 sm:grid-cols-2">
                        <InfoRow
                          icon={Route}
                          value={organizationNames[trip.organization_id] ?? trip.organization_id}
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

                    <div className="flex flex-wrap gap-3">
                      {trip.status === "created" ? (
                        <TripActionButton
                          intent="start-trip"
                          tripId={trip.trip_id}
                          disabled={isSubmitting}
                          label="Start Trip"
                        />
                      ) : null}
                      {trip.status === "active" ? (
                        <TripActionButton
                          intent="complete-trip"
                          tripId={trip.trip_id}
                          disabled={isSubmitting}
                          label="Complete Trip"
                        />
                      ) : null}
                    </div>
                  </div>
                </article>
              ))
            )}
          </div>
      </section>

      <div className="rounded-3xl border border-slate-800 bg-slate-900 p-6 text-sm text-cyan-100">
        Current active trips: {activeTrips.length}. Ready to start: {createdTrips.length}.
      </div>
    </section>
  );
};

const TripActionButton = ({ intent, tripId, disabled, label }) => {
  return (
    <Form method="post">
      <input type="hidden" name="intent" value={intent} />
      <input type="hidden" name="trip_id" value={tripId} />
      <button
        type="submit"
        disabled={disabled}
        className="inline-flex items-center justify-center rounded-xl border border-cyan-500/30 bg-cyan-500/10 px-4 py-2.5 text-sm font-medium text-cyan-200 transition hover:border-cyan-400/50 hover:bg-cyan-500/15 disabled:cursor-not-allowed disabled:border-slate-700 disabled:bg-slate-800 disabled:text-slate-500"
      >
        {label}
      </button>
    </Form>
  );
};

const InfoRow = ({ icon: Icon, value }) => {
  return (
    <div className="flex items-center gap-2">
      <Icon className="h-4 w-4 text-slate-500" />
      <span>{value}</span>
    </div>
  );
};

const StatusBadge = ({ status }) => {
  const classes =
    status === "completed"
      ? "bg-emerald-500/15 text-emerald-300"
      : status === "active"
        ? "bg-cyan-500/15 text-cyan-300"
        : "bg-amber-500/15 text-amber-300";

  return (
    <span className={`rounded-full px-3 py-1 text-xs font-medium ${classes}`}>
      {status}
    </span>
  );
};
