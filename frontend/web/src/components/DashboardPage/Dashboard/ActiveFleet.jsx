import { Dot, EllipsisVertical, Truck } from "lucide-react";

export const ActiveFleet = () => {
  const fleetCardStats = [
    {
      id: 1,
      truckNumber: "TRC-001",
      eta: "14:30",
      status: "In transit",
      startLoc: "New York",
      startLocCode: "NY",
      endLoc: "New Jersey",
      endLocCode: "NJ",
      progress: 55,
    },
    {
      id: 2,
      truckNumber: "TRC-042",
      eta: "18:15",
      status: "Delayed",
      startLoc: "Los Angeles",
      startLocCode: "LA",
      endLoc: "San Francisco",
      endLocCode: "SF",
      progress: 35,
    },
    {
      id: 3,
      truckNumber: "TRC-108",
      eta: "10:00",
      status: "Stationary",
      startLoc: "Chicago",
      startLocCode: "CHI",
      endLoc: "Detroit",
      endLocCode: "DET",
      progress: 10,
    },
    {
      id: 4,
      truckNumber: "TRC-220",
      eta: "21:45",
      status: "In transit",
      startLoc: "Dallas",
      startLocCode: "DAL",
      endLoc: "Houston",
      endLocCode: "HOU",
      progress: 75,
    },
    {
      id: 5,
      truckNumber: "TRC-315",
      eta: "16:20",
      status: "Delayed",
      startLoc: "Seattle",
      startLocCode: "SEA",
      endLoc: "Portland",
      endLocCode: "POR",
      progress: 60,
    },
    {
      id: 6,
      truckNumber: "TRC-501",
      eta: "09:50",
      status: "In transit",
      startLoc: "Miami",
      startLocCode: "MIA",
      endLoc: "Orlando",
      endLocCode: "ORL",
      progress: 25,
    },
  ];

  return (
    <div className="p-4 col-span-1 lg:col-span-6 rounded-xl border border-slate-700 shadow-2xl h-100 w-full flex flex-col">
      {/* title */}
      <div className="flex items-center justify-between">
        <span className="font-semibold">Active fleet</span>
        <div className="text-slate-500 text-sm transition-colors duration-300 hover:text-slate-200">
          <EllipsisVertical className="w-4 h-4" />
        </div>
      </div>

      {/* main card content */}
      <div className="w-full h-full grid grid-cols-1 py-4 overflow-y-scroll gap-4">
        {fleetCardStats.map((fleet) => (
          <FleetCard key={fleet.id} data={fleet} />
        ))}
      </div>
    </div>
  );
};

const FleetCard = ({ data }) => {
  const {
    truckNumber,
    eta,
    status,
    startLoc,
    startLocCode,
    endLoc,
    endLocCode,
    progress,
  } = data;

  const statusStyles = {
    "In transit": "bg-green-500/20 text-green-400 border-green-600",
    Delayed: "bg-yellow-500/20 text-yellow-400 border-yellow-600",
    Stationary: "bg-orange-500/20 text-orange-400 border-orange-600",
  };

  return (
    <div className="w-full h-30 col-span-1 border border-slate-700 rounded-xl">
      {/* top header section */}
      <div className="py-3 px-4 flex items-center justify-between">
        <div className="flex items-center text-sm tracking-tighter">
          <span>{truckNumber}</span>
          <Dot className="w-6 h-6" />
          <span>ETA {eta}</span>
        </div>
        <div
          className={`rounded-xl text-xs px-2 py-1 border ${statusStyles[status]}`}
        >
          <span>{status}</span>
        </div>
      </div>

      {/* line divider */}
      <div className="w-full border-t border-dashed border-slate-600" />

      {/* truck progress section */}
      <div className="py-2 px-3 flex items-center w-full h-auto gap-3">
        {/* truck icon */}
        <div className="p-2 bg-slate-700 rounded-xl">
          <Truck className="w-8 h-8" />
        </div>
        <div className="flex flex-col items-center">
          <span className="text-sm font-semibold tracking-tighter">
            {startLocCode}
          </span>
          <span className="text-xs tracking-tighter">
            <span className="text-xs tracking-tighter">
              {startLoc.length > 7 ? startLoc.slice(0, 7) + "..." : startLoc}
            </span>
          </span>
        </div>

        {/* progress line */}
        <div className="flex-1 px-2">
          <div className="relative w-full h-1 bg-slate-700 rounded-full">
            {/* filled progress */}
            <div
              className="absolute top-0 left-0 h-1 bg-white rounded-full"
              style={{ width: `${progress}%` }}
            />

            {/* moving indicator */}
            <div
              className="absolute top-1/2 -translate-y-1/2 w-6 h-6  flex items-center justify-center"
              style={{ left: `calc(${progress}% - 12px)` }}
            >
              <span className="text-white text-sm">▶</span>
            </div>
          </div>
        </div>

        {/* endlocation */}
        <div className="flex flex-col items-center">
          <span className="text-sm font-semibold tracking-tighter">
            {endLocCode}
          </span>
          <span className="text-xs tracking-tighter">
            {" "}
            {endLoc.length > 7 ? endLoc.slice(0, 7) + "..." : endLoc}
          </span>
        </div>
      </div>
    </div>
  );
};
