import L from "leaflet";
import {
  ArrowLeft,
  CalendarDays,
  MapPin,
  Route,
  Truck,
  UserRound,
  Wifi,
} from "lucide-react";
import { useEffect, useMemo } from "react";
import { MapContainer, Marker, Polyline, TileLayer, useMap } from "react-leaflet";
import {
  Form,
  Link,
  useLoaderData,
  useParams,
  useRouteLoaderData,
} from "react-router-dom";
import { useTrackingSocket } from "../../../hooks/useTrackingSocket";

const currentMarker = L.divIcon({
  className: "trip-location-marker",
  html: `
    <div style="
      width:18px;
      height:18px;
      background:#22d3ee;
      border-radius:50%;
      border:3px solid white;
      box-shadow:0 2px 12px rgba(0,0,0,0.35);
    "></div>
  `,
  iconSize: [18, 18],
  iconAnchor: [9, 9],
});

export const LiveTripTracking = () => {
  const { tripId } = useParams();
  const { trips, organizations, drivers, vehicles } =
    useRouteLoaderData("trip");
  const { currentLocation, locationHistory, geometry } = useLoaderData();

  const trip = trips.find(
    (entry) =>
      entry.trip_id === geometry.trip_id ||
      entry.trip_id === currentLocation?.trip_id,
  );

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

  const plannedPath = useMemo(
    () =>
      geometry.planned_geometry.map((point) => [
        point.latitude,
        point.longitude,
      ]),
    [geometry.planned_geometry],
  );
  const routePoints = useMemo(
    () => geometry.planned_geometry.map((point) => ({
      latitude: point.latitude,
      longitude: point.longitude,
    })),
    [geometry.planned_geometry],
  );
  const { connectionState, liveLocation, livePath, socketError } = useTrackingSocket({
    channel: "dispatch",
    tripId,
    driverId: trip?.driver_id,
    vehicleId: trip?.vehicle_id,
    routePoints,
  });
  const actualPath = useMemo(
    () =>
      (livePath.length > 0
        ? livePath
        : geometry.actual_geometry.length > 0
        ? geometry.actual_geometry
        : locationHistory.points
      ).map((point) => [point.latitude, point.longitude]),
    [geometry.actual_geometry, livePath, locationHistory.points],
  );
  const markerPosition = liveLocation
    ? [liveLocation.latitude, liveLocation.longitude]
    : currentLocation
      ? [currentLocation.latitude, currentLocation.longitude]
      : (actualPath.at(-1) ?? plannedPath.at(0) ?? [37.7749, -122.4194]);
  const initialMapCenter = plannedPath.at(0) ?? markerPosition;

  return (
    <section className="p-8 space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <Link
          to="/dashboard/trip/driver-console"
          className="inline-flex items-center gap-2 rounded-xl border border-slate-800 bg-slate-900 px-4 py-2 text-sm text-slate-300 transition hover:border-cyan-500/30 hover:text-white"
        >
          <ArrowLeft className="h-4 w-4" />
          Back to console
        </Link>

        {trip?.status === "active" ? (
          <Form method="post">
            <input type="hidden" name="intent" value="complete-trip" />
            <input type="hidden" name="trip_id" value={trip.trip_id} />
            <button
              type="submit"
              className="inline-flex items-center justify-center rounded-xl border border-cyan-500/30 bg-cyan-500/10 px-4 py-2.5 text-sm font-medium text-cyan-200 transition hover:border-cyan-400/50 hover:bg-cyan-500/15"
            >
              Complete Trip
            </button>
          </Form>
        ) : null}
      </div>

      <header className="rounded-3xl border border-slate-800 bg-slate-900 p-8">
        <div className="flex flex-col gap-5 xl:flex-row xl:items-start xl:justify-between">
          <div>
            <p className="text-xs font-semibold uppercase tracking-[0.35em] text-cyan-400/80">
              Live Trip Tracking
            </p>
            <h1 className="mt-3 text-3xl font-semibold text-white">
              {trip?.trip_id ?? geometry.trip_id ?? currentLocation?.trip_id}
            </h1>
            <p className="mt-3 text-sm leading-6 text-slate-400">
              Monitor current position, planned route geometry, and actual
              traveled path from one dedicated tracking page.
            </p>
          </div>

          <StatusBadge status={trip?.status ?? "active"} />
        </div>
      </header>

      <div className="grid gap-6 xl:grid-cols-[minmax(0,1.25fr)_minmax(0,0.75fr)]">
        <section className="rounded-3xl border border-slate-800 bg-slate-900 p-6">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold text-white">Tracking Map</h2>
            <p className="text-sm text-slate-400">
              Planned route: cyan · Actual path: emerald
            </p>
          </div>

          <div className="mt-5 h-[34rem] overflow-hidden rounded-2xl border border-slate-800">
            <MapContainer
              center={initialMapCenter}
              zoom={13}
              attributionControl={false}
              className="h-full w-full"
            >
              <TileLayer url="https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{r}.png" />
              {markerPosition ? (
                <FollowDriverMarker position={markerPosition} />
              ) : null}
              {plannedPath.length > 1 ? (
                <Polyline
                  positions={plannedPath}
                  pathOptions={{ color: "#22d3ee", weight: 5, opacity: 0.75 }}
                />
              ) : null}
              {actualPath.length > 1 ? (
                <Polyline
                  positions={actualPath}
                  pathOptions={{ color: "#10b981", weight: 5, opacity: 0.9 }}
                />
              ) : null}
              {markerPosition ? (
                <Marker position={markerPosition} icon={currentMarker} />
              ) : null}
            </MapContainer>
          </div>
        </section>

        <section className="space-y-6">
          <div className="rounded-3xl border border-slate-800 bg-slate-900 p-6">
            <h2 className="text-lg font-semibold text-white">Trip Overview</h2>

            <div className="mt-5 space-y-4 text-sm text-slate-300">
              <DetailRow
                icon={Route}
                label="Organization"
                value={
                  organizationNames[trip?.organization_id] ??
                  trip?.organization_id ??
                  "-"
                }
              />
              <DetailRow
                icon={UserRound}
                label="Driver"
                value={driverNames[trip?.driver_id] ?? trip?.driver_id ?? "-"}
              />
              <DetailRow
                icon={Truck}
                label="Vehicle"
                value={
                  vehicleNames[trip?.vehicle_id] ?? trip?.vehicle_id ?? "-"
                }
              />
              <DetailRow
                icon={CalendarDays}
                label="Created"
                value={
                  trip?.created_at
                    ? new Date(trip.created_at * 1000).toLocaleString()
                    : "-"
                }
              />
            </div>
          </div>

          <div className="rounded-3xl border border-slate-800 bg-slate-900 p-6">
            <h2 className="text-lg font-semibold text-white">
              Tracking Status
            </h2>

            <div className="mt-5 space-y-4 text-sm text-slate-300">
              <DetailRow
                icon={Wifi}
                label="Connection"
                value={currentLocation?.connection || connectionState}
              />
              <DetailRow
                icon={MapPin}
                label="Current Position"
                value={
                  liveLocation
                    ? `${liveLocation.latitude.toFixed(5)}, ${liveLocation.longitude.toFixed(5)}`
                    : currentLocation
                    ? `${currentLocation.latitude.toFixed(5)}, ${currentLocation.longitude.toFixed(5)}`
                    : "No live point yet"
                }
              />
              <DetailRow
                icon={Route}
                label="History Points"
                value={String(locationHistory.points.length)}
              />
              <DetailRow
                icon={Route}
                label="Planned Geometry"
                value={String(geometry.planned_geometry.length)}
              />
            </div>
          </div>
        </section>
      </div>
    </section>
  );
};

const FollowDriverMarker = ({ position }) => {
  const map = useMap();

  useEffect(() => {
    map.panTo(position, { animate: true, duration: 0.5 });
  }, [map, position]);

  return null;
};

const DetailRow = ({ icon: Icon, label, value }) => {
  return (
    <div className="flex items-start justify-between gap-4 border-b border-slate-800 pb-4 last:border-b-0 last:pb-0">
      <div className="flex items-center gap-2 text-slate-500">
        <Icon className="h-4 w-4" />
        <span>{label}</span>
      </div>
      <span className="max-w-[16rem] text-right text-slate-200">{value}</span>
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
