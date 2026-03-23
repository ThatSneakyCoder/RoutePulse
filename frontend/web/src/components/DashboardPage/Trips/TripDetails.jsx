import { ArrowLeft, MapPin, Truck, UserRound } from "lucide-react";
import L from "leaflet";
import { useMemo } from "react";
import { Link, Navigate, useParams } from "react-router-dom";
import { MapContainer, Marker, Polyline, TileLayer } from "react-leaflet";
import { getTripById } from "./tripData.js";

const tripMarker = L.divIcon({
  className: "trip-location-marker",
  html: `
    <div style="
      width:14px;
      height:14px;
      background:#22d3ee;
      border-radius:50%;
      border:2px solid white;
      box-shadow:0 2px 8px rgba(0,0,0,0.35);
    "></div>
  `,
  iconSize: [14, 14],
  iconAnchor: [7, 7],
});

export const TripDetails = () => {
  const { tripId } = useParams();
  const trip = getTripById(tripId);

  const markerPosition = useMemo(() => {
    if (!trip) return null;

    return trip.location.path[trip.location.path.length - 1];
  }, [trip]);

  if (!trip) {
    return <Navigate to="/dashboard/trip" replace />;
  }

  return (
    <section className="p-8 space-y-6">
      <div className="flex items-center gap-3">
        <Link
          to="/dashboard/trip"
          className="inline-flex items-center gap-2 rounded-xl border border-slate-800 bg-slate-900 px-4 py-2 text-sm text-slate-300 transition hover:border-cyan-500/30 hover:text-white"
        >
          <ArrowLeft className="h-4 w-4" />
          Back to trips
        </Link>
      </div>

      <header className="rounded-3xl border border-slate-800 bg-slate-900 p-8">
        <div className="flex flex-col gap-5 lg:flex-row lg:items-start lg:justify-between">
          <div>
            <p className="text-xs font-semibold uppercase tracking-[0.35em] text-cyan-400/80">
              Trip Details
            </p>
            <h1 className="mt-3 text-3xl font-semibold text-white">{trip.id}</h1>
            <p className="mt-2 text-sm text-slate-400">{trip.routeName}</p>
          </div>

          <span
            className={`rounded-full px-3 py-1 text-xs font-medium ${
              trip.status === "Started"
                ? "bg-emerald-500/15 text-emerald-300"
                : "bg-amber-500/15 text-amber-300"
            }`}
          >
            {trip.status}
          </span>
        </div>
      </header>

      <div className="grid grid-cols-1 gap-6 xl:grid-cols-[minmax(0,0.9fr)_minmax(0,1.1fr)]">
        <section className="rounded-3xl border border-slate-800 bg-slate-900 p-6">
          <h2 className="text-lg font-semibold text-white">Overview</h2>

          <div className="mt-5 space-y-4 text-sm text-slate-300">
            <DetailRow icon={UserRound} label="Driver" value={trip.driver} />
            <DetailRow icon={Truck} label="Vehicle" value={trip.vehicle} />
            <DetailRow icon={MapPin} label="Origin" value={trip.origin} />
            <DetailRow
              icon={MapPin}
              label="Destination"
              value={trip.destination}
            />
            <DetailRow label="Date" value={trip.date} />
          </div>

          <p className="mt-6 rounded-2xl border border-slate-800 bg-slate-950/70 p-4 text-sm leading-6 text-slate-400">
            {trip.notes}
          </p>
        </section>

        <section className="rounded-3xl border border-slate-800 bg-slate-900 p-6">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold text-white">Location</h2>
            <p className="text-sm text-slate-400">Placeholder geometry</p>
          </div>

          <div className="mt-5 h-[24rem] overflow-hidden rounded-2xl border border-slate-800">
            <MapContainer
              center={trip.location.center}
              zoom={12}
              attributionControl={false}
              scrollWheelZoom={false}
              className="h-full w-full"
            >
              <TileLayer url="https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{r}.png" />
              <Polyline
                positions={trip.location.path}
                pathOptions={{ color: "#22d3ee", weight: 4, opacity: 0.9 }}
              />
              {markerPosition ? (
                <Marker position={markerPosition} icon={tripMarker} />
              ) : null}
            </MapContainer>
          </div>
        </section>
      </div>
    </section>
  );
};

const DetailRow = ({ icon: Icon, label, value }) => {
  return (
    <div className="flex items-start justify-between gap-4 border-b border-slate-800 pb-4 last:border-b-0 last:pb-0">
      <div className="flex items-center gap-2 text-slate-500">
        {Icon ? <Icon className="h-4 w-4" /> : null}
        <span>{label}</span>
      </div>
      <span className="max-w-[16rem] text-right text-slate-200">{value}</span>
    </div>
  );
};
