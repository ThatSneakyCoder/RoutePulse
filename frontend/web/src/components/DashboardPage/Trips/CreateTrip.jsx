import { ArrowLeft, MapPin, Route, Truck, UserRound } from "lucide-react";
import L from "leaflet";
import { useEffect, useMemo, useState } from "react";
import {
  Form,
  Link,
  useActionData,
  useNavigation,
  useRouteLoaderData,
} from "react-router-dom";
import {
  MapContainer,
  Marker,
  Polyline,
  TileLayer,
  useMap,
  useMapEvents,
} from "react-leaflet";
import instance from "../../../axios";

const startMarkerIcon = L.divIcon({
  className: "trip-location-marker",
  html: `
    <div style="
      width:16px;
      height:16px;
      background:#22d3ee;
      border-radius:50%;
      border:3px solid white;
      box-shadow:0 2px 10px rgba(0,0,0,0.35);
    "></div>
  `,
  iconSize: [16, 16],
  iconAnchor: [8, 8],
});

const endMarkerIcon = L.divIcon({
  className: "trip-location-marker",
  html: `
    <div style="
      width:16px;
      height:16px;
      background:#f59e0b;
      border-radius:50%;
      border:3px solid white;
      box-shadow:0 2px 10px rgba(0,0,0,0.35);
    "></div>
  `,
  iconSize: [16, 16],
  iconAnchor: [8, 8],
});

const distributorLocation = {
  lat: 37.7749,
  lng: -122.4194,
};

export const CreateTrip = () => {
  const { organizations, drivers, vehicles } = useRouteLoaderData("trip");
  const actionData = useActionData();
  const navigation = useNavigation();
  const isSubmitting = navigation.state === "submitting";

  const [organizationId, setOrganizationId] = useState(
    organizations[0]?.organization_id ?? "",
  );
  const [driverId, setDriverId] = useState("");
  const [vehicleId, setVehicleId] = useState("");
  const [startLocation] = useState(distributorLocation);
  const [endLocation, setEndLocation] = useState(null);
  const [routePreview, setRoutePreview] = useState(null);
  const [isRouteLoading, setIsRouteLoading] = useState(false);
  const [routeError, setRouteError] = useState("");

  const organizationDrivers = drivers.filter(
    (driver) =>
      driver.organization_id === organizationId && driver.status === "active",
  );
  const selectedDriver = organizationDrivers.find(
    (driver) => driver.driver_id === driverId,
  );
  const linkedVehicles = vehicles.filter(
    (vehicle) =>
      vehicle.organization_id === organizationId &&
      vehicle.status === "active" &&
      vehicle.vehicle_id === selectedDriver?.vehicle_id,
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

  useEffect(() => {
    if (!organizationDrivers.some((driver) => driver.driver_id === driverId)) {
      setDriverId(organizationDrivers[0]?.driver_id ?? "");
    }
  }, [driverId, organizationDrivers]);

  useEffect(() => {
    if (!linkedVehicles.some((vehicle) => vehicle.vehicle_id === vehicleId)) {
      setVehicleId(linkedVehicles[0]?.vehicle_id ?? "");
    }
  }, [linkedVehicles, vehicleId]);

  const mapCenter = useMemo(() => {
    if (endLocation) return [endLocation.lat, endLocation.lng];
    if (startLocation) return [startLocation.lat, startLocation.lng];
    return [distributorLocation.lat, distributorLocation.lng];
  }, [endLocation, startLocation]);

  const routeGeometry = useMemo(
    () =>
      routePreview?.geometry?.map((point) => [point.latitude, point.longitude]) ??
      [],
    [routePreview],
  );

  useEffect(() => {
    if (!endLocation) {
      setRoutePreview(null);
      setRouteError("");
      setIsRouteLoading(false);
      return;
    }

    let cancelled = false;

    async function previewRoute() {
      setIsRouteLoading(true);
      setRouteError("");

      try {
        const response = await instance.post("/v1/fleet/trips/preview-route", {
          start_latitude: startLocation.lat,
          start_longitude: startLocation.lng,
          end_latitude: endLocation.lat,
          end_longitude: endLocation.lng,
        });

        if (cancelled) return;

        setRoutePreview(response.data.data[0] ?? null);
      } catch (error) {
        if (cancelled) return;

        setRoutePreview(null);
        setRouteError(
          error.response?.data?.error ||
            "We couldn't generate the road route right now.",
        );
      } finally {
        if (!cancelled) {
          setIsRouteLoading(false);
        }
      }
    }

    previewRoute();

    return () => {
      cancelled = true;
    };
  }, [endLocation, startLocation]);

  return (
    <section className="w-full p-8">
      <div className="w-full space-y-6">
        <Link
          to="/dashboard/trip"
          className="inline-flex items-center gap-2 text-sm text-slate-400 transition hover:text-cyan-200"
        >
          <ArrowLeft className="h-4 w-4" />
          Back to trips
        </Link>

        <div className="rounded-3xl border border-slate-800 bg-slate-900 p-8">
          <div className="max-w-3xl">
            <p className="text-xs font-semibold uppercase tracking-[0.35em] text-cyan-400/80">
              New Trip
            </p>
            <h1 className="mt-3 text-3xl font-semibold text-white">
              Create a trip with the right assignment.
            </h1>
            <p className="mt-3 text-sm leading-6 text-slate-400">
              Choose an organization, pick an active driver, and mark the
              destination directly on the map. The distributor start location is
              fixed for every trip.
            </p>
          </div>
        </div>

        <Form
          method="post"
          className="rounded-3xl border border-slate-800 bg-slate-900 p-6 sm:p-8"
        >
          <input type="hidden" name="intent" value="create-trip" />
          <input
            type="hidden"
            name="start_latitude"
            value={startLocation?.lat ?? ""}
          />
          <input
            type="hidden"
            name="start_longitude"
            value={startLocation?.lng ?? ""}
          />
          <input
            type="hidden"
            name="end_latitude"
            value={endLocation?.lat ?? ""}
          />
          <input
            type="hidden"
            name="end_longitude"
            value={endLocation?.lng ?? ""}
          />

          <div className="grid gap-6 xl:grid-cols-[minmax(0,1.6fr)_minmax(20rem,0.72fr)]">
            <div className="space-y-6">
              <div className="grid gap-5 lg:grid-cols-2">
                <Field label="Organization" icon={Route}>
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
                </Field>

                <Field label="Driver" icon={UserRound}>
                  <select
                    name="driver_id"
                    value={driverId}
                    onChange={(event) => setDriverId(event.target.value)}
                    className="w-full rounded-2xl border border-slate-700 bg-slate-950 px-4 py-3 text-sm text-white outline-none transition focus:border-cyan-500/50"
                    required
                    disabled={organizationDrivers.length === 0}
                  >
                    {organizationDrivers.length === 0 ? (
                      <option value="">No active drivers available</option>
                    ) : (
                      organizationDrivers.map((driver) => (
                        <option key={driver.driver_id} value={driver.driver_id}>
                          {driver.first_name} {driver.last_name}
                        </option>
                      ))
                    )}
                  </select>
                </Field>
              </div>

              <Field label="Linked Vehicle" icon={Truck}>
                <select
                  name="vehicle_id"
                  value={vehicleId}
                  onChange={(event) => setVehicleId(event.target.value)}
                  className="w-full rounded-2xl border border-slate-700 bg-slate-950 px-4 py-3 text-sm text-white outline-none transition focus:border-cyan-500/50"
                  required
                  disabled={linkedVehicles.length === 0}
                >
                  {linkedVehicles.length === 0 ? (
                    <option value="">No linked active vehicle for this driver</option>
                  ) : (
                    linkedVehicles.map((vehicle) => (
                      <option key={vehicle.vehicle_id} value={vehicle.vehicle_id}>
                        {vehicle.plate_number}
                        {vehicle.vehicle_type ? ` • ${vehicle.vehicle_type}` : ""}
                      </option>
                    ))
                  )}
                </select>
              </Field>

              <div className="rounded-3xl border border-slate-800 bg-slate-950/60 p-5">
                <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
                  <div>
                    <p className="text-sm font-medium text-white">Location Picker</p>
                    <p className="mt-1 text-sm text-slate-400">
                      Click the map to place the destination marker. The start
                      point is already fixed to the distributor location.
                    </p>
                  </div>
                </div>

                <div className="mt-5 h-[26rem] overflow-hidden rounded-2xl border border-slate-800">
                  <MapContainer
                    center={mapCenter}
                    zoom={12}
                    attributionControl={false}
                    className="h-full w-full"
                  >
                    <TileLayer url="https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{r}.png" />
                    <MapClickHandler
                      setEndLocation={setEndLocation}
                    />
                    <RouteViewport
                      routeGeometry={routeGeometry}
                      fallbackCenter={mapCenter}
                    />
                    {startLocation ? (
                      <Marker
                        position={[startLocation.lat, startLocation.lng]}
                        icon={startMarkerIcon}
                      />
                    ) : null}
                    {endLocation ? (
                      <Marker
                        position={[endLocation.lat, endLocation.lng]}
                        icon={endMarkerIcon}
                      />
                    ) : null}
                    {routeGeometry.length > 1 ? (
                      <Polyline
                        positions={routeGeometry}
                        pathOptions={{ color: "#22d3ee", weight: 5, opacity: 0.85 }}
                      />
                    ) : null}
                  </MapContainer>
                </div>

                {isRouteLoading ? (
                  <p className="mt-4 text-sm text-cyan-200">
                    Building road route preview...
                  </p>
                ) : null}
                {routeError ? (
                  <p className="mt-4 text-sm text-rose-300">{routeError}</p>
                ) : null}

                <div className="mt-5 grid gap-4 lg:grid-cols-2">
                  <LocationCard
                    title="Start Location"
                    accent="cyan"
                    location={startLocation}
                    description="Fixed distributor location"
                  />
                  <LocationCard
                    title="End Location"
                    accent="amber"
                    location={endLocation}
                    active
                    description="Selected from the map"
                    actionLabel="Pick on map"
                    onClear={() => setEndLocation(null)}
                  />
                </div>
              </div>
            </div>

            <aside className="rounded-3xl border border-slate-800 bg-slate-950/70 p-5">
              <p className="text-xs font-semibold uppercase tracking-[0.3em] text-slate-500">
                Assignment
              </p>
              <div className="mt-5 space-y-4 text-sm text-slate-300">
                <SummaryRow
                  label="Organization"
                  value={
                    organizations.find(
                      (organization) =>
                        organization.organization_id === organizationId,
                    )?.name ?? "Select organization"
                  }
                />
                <SummaryRow
                  label="Driver"
                  value={
                    selectedDriver
                      ? `${selectedDriver.first_name} ${selectedDriver.last_name}`
                      : "Select driver"
                  }
                />
                <SummaryRow
                  label="Vehicle"
                  value={
                    linkedVehicles.find(
                      (vehicle) => vehicle.vehicle_id === vehicleId,
                    )?.plate_number ?? "No linked vehicle"
                  }
                />
                <SummaryRow
                  label="Start"
                  value={formatLocation(startLocation)}
                />
                <SummaryRow label="End" value={formatLocation(endLocation)} />
                <SummaryRow
                  label="Distance"
                  value={formatDistance(routePreview?.distance_meters)}
                />
                <SummaryRow
                  label="ETA"
                  value={formatDuration(routePreview?.duration_seconds)}
                />
              </div>

              <div className="mt-6 rounded-2xl border border-cyan-500/20 bg-cyan-500/10 p-4 text-sm text-cyan-100">
                Only vehicles linked to the selected driver are available here.
              </div>

              <div className="mt-4 rounded-2xl border border-slate-800 bg-slate-900/70 p-4 text-sm text-slate-300">
                Start point is fixed to the distributor location. Use the map to
                choose where the driver should go.
              </div>

              {actionData?.error ? (
                <p className="mt-4 rounded-2xl border border-rose-500/30 bg-rose-500/10 px-4 py-3 text-sm text-rose-200">
                  {actionData.error}
                </p>
              ) : null}
              {routeError ? (
                <p className="mt-4 rounded-2xl border border-rose-500/30 bg-rose-500/10 px-4 py-3 text-sm text-rose-200">
                  {routeError}
                </p>
              ) : null}

              <button
                type="submit"
                disabled={
                  isSubmitting ||
                  !organizationId ||
                  !driverId ||
                  !vehicleId ||
                  !startLocation ||
                  !endLocation
                }
                className="mt-6 inline-flex w-full items-center justify-center gap-2 rounded-2xl border border-cyan-500/30 bg-cyan-500/10 px-4 py-3 text-sm font-medium text-cyan-200 transition hover:border-cyan-400/50 hover:bg-cyan-500/15 disabled:cursor-not-allowed disabled:border-slate-700 disabled:bg-slate-800 disabled:text-slate-500"
              >
                <MapPin className="h-4 w-4" />
                {isSubmitting ? "Creating Trip..." : "Create Trip"}
              </button>
            </aside>
          </div>
        </Form>
      </div>
    </section>
  );
};

const MapClickHandler = ({
  setEndLocation,
}) => {
  useMapEvents({
    click(event) {
      const { lat, lng } = event.latlng;
      const nextLocation = {
        lat: Number(lat.toFixed(6)),
        lng: Number(lng.toFixed(6)),
      };

      setEndLocation(nextLocation);
    },
  });

  return null;
};

const RouteViewport = ({ routeGeometry, fallbackCenter }) => {
  const map = useMap();

  useEffect(() => {
    if (routeGeometry.length > 1) {
      map.fitBounds(routeGeometry, {
        padding: [40, 40],
      });
      return;
    }

    map.setView(fallbackCenter, 12);
  }, [fallbackCenter, map, routeGeometry]);

  return null;
};

const Field = ({ label, icon: Icon, children }) => {
  return (
    <label className="block space-y-2">
      <span className="inline-flex items-center gap-2 text-sm font-medium text-slate-200">
        <Icon className="h-4 w-4 text-cyan-400" />
        {label}
      </span>
      {children}
    </label>
  );
};

const LocationCard = ({
  title,
  accent,
  location,
  active,
  description,
  actionLabel,
  onClear,
}) => {
  const accentClass =
    accent === "cyan"
      ? "border-cyan-500/20 bg-cyan-500/8"
      : "border-amber-500/20 bg-amber-500/8";
  const accentText = accent === "cyan" ? "text-cyan-300" : "text-amber-300";

  return (
    <div
      className={`rounded-3xl border p-5 transition ${
        active ? accentClass : "border-slate-800 bg-slate-950/50"
      }`}
    >
      <div className="flex items-center justify-between gap-3">
        <div>
          <p className={`text-sm font-medium ${active ? accentText : "text-white"}`}>
            {title}
          </p>
          <p className="mt-1 text-sm text-slate-400">
            {location ? formatLocation(location) : "Not selected on map yet"}
          </p>
          {description ? (
            <p className="mt-2 text-xs uppercase tracking-[0.25em] text-slate-500">
              {description}
            </p>
          ) : null}
        </div>

        <div className="flex gap-2">
          {actionLabel ? (
            <div className="rounded-xl border border-slate-700 px-3 py-2 text-xs font-medium text-slate-200">
              {actionLabel}
            </div>
          ) : null}
          <button
            type="button"
            onClick={onClear}
            className="rounded-xl border border-slate-700 px-3 py-2 text-xs font-medium text-slate-400 transition hover:border-slate-600 hover:text-slate-200"
            disabled={!onClear}
          >
            Clear
          </button>
        </div>
      </div>
    </div>
  );
};

const SummaryRow = ({ label, value }) => {
  return (
    <div className="flex items-center justify-between gap-4 border-b border-slate-800 pb-3 last:border-b-0 last:pb-0">
      <span className="text-slate-500">{label}</span>
      <span className="text-right text-white">{value}</span>
    </div>
  );
};

function formatLocation(location) {
  if (!location) {
    return "Not selected";
  }

  return `${location.lat.toFixed(5)}, ${location.lng.toFixed(5)}`;
}

function formatDistance(distanceMeters) {
  if (!distanceMeters) {
    return "Waiting for route";
  }

  return `${(distanceMeters / 1000).toFixed(1)} km`;
}

function formatDuration(durationSeconds) {
  if (!durationSeconds) {
    return "Waiting for route";
  }

  return `${Math.round(durationSeconds / 60)} min`;
}
