import { useRouteLoaderData } from "react-router-dom";

export const Vehicles = () => {
  const { vehicles } = useRouteLoaderData("fleet");

  if (!vehicles || vehicles.length === 0) {
    return <div className="p-8 text-gray-400">No vehicles found</div>;
  }

  return (
    <div className="p-8">
      <h1 className="text-2xl font-semibold text-white mb-8">Vehicles</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
        {vehicles.map((vehicle) => (
          <div
            key={vehicle.vehicle_id}
            className="
              bg-slate-800
              border border-slate-700
              rounded-xl
              p-6
              shadow-lg
              hover:border-blue-500
              hover:shadow-blue-500/10
              transition
              duration-200
            "
          >
            {/* header */}
            <div className="flex justify-between items-center mb-5">
              <h2 className="text-lg font-semibold text-white">
                {vehicle.plate_number}
              </h2>

              <span
                className={`text-xs px-3 py-1 rounded-full font-medium
                  ${
                    vehicle.status === "active"
                      ? "bg-green-500/20 text-green-400"
                      : "bg-red-500/20 text-red-400"
                  }
                `}
              >
                {vehicle.status}
              </span>
            </div>

            {/* content */}
            <div className="space-y-3 text-sm text-gray-300">
              <div className="flex justify-between">
                <span className="text-gray-400">Type</span>
                <span className="font-medium">{vehicle.vehicle_type}</span>
              </div>

              <div className="flex justify-between">
                <span className="text-gray-400">Capacity</span>
                <span className="font-medium">{vehicle.capacity}</span>
              </div>

              <div className="flex justify-between">
                <span className="text-gray-400">Organization</span>
                <span className="font-mono text-xs truncate max-w-35">
                  {vehicle.organization_id}
                </span>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
