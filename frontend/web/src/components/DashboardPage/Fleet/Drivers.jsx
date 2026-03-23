import { Truck, UserPlus, UserRound } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import {
  Form,
  useActionData,
  useNavigation,
  useRouteLoaderData,
} from "react-router-dom";

export const Drivers = () => {
  const { drivers, organizations, vehicles } = useRouteLoaderData("fleet");
  const actionData = useActionData();
  const navigation = useNavigation();
  const isSubmitting = navigation.state === "submitting";

  const [organizationId, setOrganizationId] = useState(
    organizations[0]?.organization_id ?? "",
  );

  useEffect(() => {
    if (
      !organizations.some(
        (organization) => organization.organization_id === organizationId,
      )
    ) {
      setOrganizationId(organizations[0]?.organization_id ?? "");
    }
  }, [organizationId, organizations]);

  const orgVehicles = useMemo(
    () =>
      vehicles.filter(
        (vehicle) =>
          vehicle.organization_id === organizationId &&
          vehicle.status === "active",
      ),
    [organizationId, vehicles],
  );

  const organizationNames = Object.fromEntries(
    organizations.map((organization) => [
      organization.organization_id,
      organization.name,
    ]),
  );
  const vehicleNames = Object.fromEntries(
    vehicles.map((vehicle) => [vehicle.vehicle_id, vehicle.plate_number]),
  );

  const sortedDrivers = [...drivers].sort((left, right) =>
    `${left.first_name} ${left.last_name}`.localeCompare(
      `${right.first_name} ${right.last_name}`,
    ),
  );

  return (
    <section className="p-8 space-y-8">
      <header className="rounded-3xl border border-slate-800 bg-slate-900 p-8">
        <div className="flex flex-col gap-6 xl:flex-row xl:items-start xl:justify-between">
          <div className="max-w-2xl">
            <p className="text-xs font-semibold uppercase tracking-[0.35em] text-cyan-400/80">
              Drivers
            </p>
            <h1 className="mt-3 text-3xl font-semibold text-white">
              Manage drivers and assign them to vehicles.
            </h1>
            <p className="mt-3 text-sm leading-6 text-slate-400">
              Create drivers for any organization you manage and optionally link
              them to an active vehicle during setup.
            </p>
          </div>

          <div className="grid grid-cols-2 gap-4 sm:w-auto">
            <StatCard label="Total Drivers" value={String(drivers.length)} />
            <StatCard
              label="Active Drivers"
              value={String(drivers.filter((driver) => driver.status === "active").length)}
            />
          </div>
        </div>
      </header>

      <div className="grid gap-6 xl:grid-cols-[minmax(0,0.9fr)_minmax(0,1.1fr)]">
        <section className="rounded-3xl border border-slate-800 bg-slate-900 p-6 sm:p-8">
          <div className="flex items-center gap-2 text-white">
            <UserPlus className="h-4 w-4 text-cyan-400" />
            <h2 className="text-lg font-semibold">Create Driver</h2>
          </div>

          <Form method="post" className="mt-6 space-y-5">
            <input type="hidden" name="intent" value="create-driver" />

            <label className="block space-y-2">
              <span className="text-sm font-medium text-slate-200">
                Organization
              </span>
              <select
                name="organization_id"
                value={organizationId}
                onChange={(event) => setOrganizationId(event.target.value)}
                className="w-full rounded-2xl border border-slate-700 bg-slate-950 px-4 py-3 text-sm text-white outline-none transition focus:border-cyan-500/50"
                required
              >
                {organizations.map((organization) => (
                  <option
                    key={organization.organization_id}
                    value={organization.organization_id}
                  >
                    {organization.name}
                  </option>
                ))}
              </select>
            </label>

            <div className="grid gap-5 md:grid-cols-2">
              <label className="block space-y-2">
                <span className="text-sm font-medium text-slate-200">
                  First Name
                </span>
                <input
                  name="first_name"
                  placeholder="Aarav"
                  className="w-full rounded-2xl border border-slate-700 bg-slate-950 px-4 py-3 text-sm text-white outline-none transition focus:border-cyan-500/50"
                  required
                />
              </label>

              <label className="block space-y-2">
                <span className="text-sm font-medium text-slate-200">
                  Last Name
                </span>
                <input
                  name="last_name"
                  placeholder="Singh"
                  className="w-full rounded-2xl border border-slate-700 bg-slate-950 px-4 py-3 text-sm text-white outline-none transition focus:border-cyan-500/50"
                  required
                />
              </label>
            </div>

            <label className="block space-y-2">
              <span className="text-sm font-medium text-slate-200">
                Linked Vehicle
              </span>
              <select
                name="vehicle_id"
                className="w-full rounded-2xl border border-slate-700 bg-slate-950 px-4 py-3 text-sm text-white outline-none transition focus:border-cyan-500/50"
                defaultValue=""
              >
                <option value="">No vehicle assigned</option>
                {orgVehicles.map((vehicle) => (
                  <option key={vehicle.vehicle_id} value={vehicle.vehicle_id}>
                    {vehicle.plate_number}
                    {vehicle.vehicle_type ? ` • ${vehicle.vehicle_type}` : ""}
                  </option>
                ))}
              </select>
            </label>

            {actionData?.error ? (
              <p className="rounded-2xl border border-rose-500/30 bg-rose-500/10 px-4 py-3 text-sm text-rose-200">
                {actionData.error}
              </p>
            ) : null}

            <button
              type="submit"
              disabled={isSubmitting || organizations.length === 0}
              className="inline-flex items-center justify-center gap-2 rounded-2xl border border-cyan-500/30 bg-cyan-500/10 px-5 py-3 text-sm font-medium text-cyan-200 transition hover:border-cyan-400/50 hover:bg-cyan-500/15 disabled:cursor-not-allowed disabled:border-slate-700 disabled:bg-slate-800 disabled:text-slate-500"
            >
              <UserPlus className="h-4 w-4" />
              {isSubmitting ? "Creating Driver..." : "Create Driver"}
            </button>
          </Form>
        </section>

        <section className="flex min-h-0 flex-col rounded-3xl border border-slate-800 bg-slate-900 p-6 sm:p-8 xl:max-h-[28rem]">
          <div className="flex items-center justify-between gap-4">
            <div>
              <p className="text-xs font-semibold uppercase tracking-[0.3em] text-slate-500">
                Driver Directory
              </p>
              <h2 className="mt-2 text-xl font-semibold text-white">
                Existing Drivers
              </h2>
            </div>
            <p className="text-sm text-slate-400">{sortedDrivers.length} total</p>
          </div>

          <div className="mt-6 min-h-0 flex-1 space-y-4 overflow-y-auto pr-2">
            {sortedDrivers.length === 0 ? (
              <div className="rounded-2xl border border-dashed border-slate-700 bg-slate-950/50 p-8 text-center">
                <p className="text-base font-medium text-white">No drivers yet.</p>
                <p className="mt-2 text-sm text-slate-400">
                  Create your first driver from the panel on the left.
                </p>
              </div>
            ) : (
              sortedDrivers.map((driver) => (
                <article
                  key={driver.driver_id}
                  className="rounded-2xl border border-slate-800 bg-slate-950/70 p-5"
                >
                  <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
                    <div className="space-y-3">
                      <div className="flex flex-wrap items-center gap-3">
                        <h3 className="text-lg font-semibold text-white">
                          {driver.first_name} {driver.last_name}
                        </h3>
                        <StatusBadge status={driver.status} />
                      </div>

                      <div className="grid gap-2 text-sm text-slate-400 sm:grid-cols-2">
                        <InfoRow
                          icon={UserRound}
                          value={organizationNames[driver.organization_id] ?? driver.organization_id}
                        />
                        <InfoRow
                          icon={Truck}
                          value={vehicleNames[driver.vehicle_id] ?? "No vehicle linked"}
                        />
                      </div>
                    </div>

                    <p className="text-xs text-slate-500">
                      Created {new Date(driver.created_at * 1000).toLocaleDateString()}
                    </p>
                  </div>
                </article>
              ))
            )}
          </div>
        </section>
      </div>
    </section>
  );
};

const StatCard = ({ label, value }) => {
  return (
    <div className="rounded-2xl border border-slate-800 bg-slate-950/70 px-5 py-4">
      <p className="text-sm text-slate-400">{label}</p>
      <p className="mt-2 text-2xl font-semibold text-white">{value}</p>
    </div>
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
    status === "active"
      ? "bg-emerald-500/15 text-emerald-300"
      : "bg-rose-500/15 text-rose-300";

  return (
    <span className={`rounded-full px-3 py-1 text-xs font-medium ${classes}`}>
      {status}
    </span>
  );
};
