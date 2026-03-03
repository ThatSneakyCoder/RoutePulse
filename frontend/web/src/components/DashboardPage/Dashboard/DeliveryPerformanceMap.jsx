import { Search } from "lucide-react";
import {
  MapContainer,
  Marker,
  Popup,
  TileLayer,
  useMap,
  Tooltip,
} from "react-leaflet";

export const DeliveryPerformanceMap = () => {
  return (
    <div className="p-4 col-span-1 lg:col-span-7 rounded-xl border border-slate-700 shadow-2xl h-auto w-full flex flex-col gap-4">
      <div className="flex items-center justify-between">
        <span className="font-semibold">Delivery performance</span>
        <div className="flex items-center gap-2 border border-slate-600 rounded-lg py-2 px-4 hover:bg-slate-600/40 transition-colors duration-300">
          <Search className="w-4 h-4" />
          <input
            type="text"
            placeholder="Search for fleet"
            className="bg-transparent outline-none text-sm placeholder:text-slate-400 w-full"
          />
        </div>
      </div>
      <div className="border border-slate-700 w-full h-80 md:h-110 lg:h-full rounded-lg overflow-hidden relative">
        <Map />
        {/* bottom overlay panel */}
        <div className="absolute bottom-4 left-1/2 -translate-x-1/2 bg-slate-900/80 backdrop-blur-lg border border-slate-700 rounded-xl px-6 py-3 flex items-center gap-10 shadow-lg z-1000">
          <LegendItem color="bg-green-500" label="In transit" />
          <LegendItem color="bg-cyan-400" label="Pick-up" />
          <LegendItem color="bg-yellow-400" label="Stationary" />
          <LegendItem color="bg-red-500" label="Delayed" />
        </div>
      </div>
    </div>
  );
};

const LegendItem = ({ color, label }) => {
  return (
    <div className="flex items-center gap-2 text-sm text-slate-300 whitespace-nowrap">
      <span className={`w-2.5 h-2.5 rounded-full ${color}`} />
      {label}
    </div>
  );
};

import L from "leaflet";

const getTruckIcon = (color = "#22c55e") =>
  L.divIcon({
    className: "fleet-marker",
    html: `
      <div style="
        width:14px;
        height:14px;
        background:${color};
        border-radius:50%;
        border:2px solid white;
        box-shadow:0 2px 8px rgba(0,0,0,0.4);
      "></div>
    `,
    iconSize: [14, 14],
    iconAnchor: [7, 7],
  });

const fleet = [
  {
    id: "TRC-042",
    position: [51.505, -0.09],
    status: "in-transit",
  },
  {
    id: "TRC-078",
    position: [51.51, -0.09],
    status: "delayed",
  },
  {
    id: "TRC-103",
    position: [51.498, -0.13],
    status: "stationary",
  },
];

const Map = () => {
  return (
    <MapContainer
      center={[51.505, -0.09]}
      zoom={13}
      zoomControl={true}
      attributionControl={false}
      scrollWheelZoom={false}
      className="h-full w-full"
    >
      <TileLayer
        attribution="&copy; OpenStreetMap contributors &copy; CARTO"
        url="https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{r}.png"
      />
      {fleet.map((truck) => (
        <Marker
          riseOnHover={true}
          key={truck.id}
          position={truck.position}
          icon={getTruckIcon("green")}
        >
          <Tooltip
            permanent
            direction="top"
            offset={[0, -25]}
            className="fleet-label"
          >
            {truck.id} <br />
            Status: {truck.status}
          </Tooltip>
        </Marker>
      ))}
    </MapContainer>
  );
};
