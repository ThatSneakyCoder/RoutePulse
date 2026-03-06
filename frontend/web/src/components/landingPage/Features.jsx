import realTimeFleet from "../../assets/realTimeFleet.svg";
import realTimeAnalytics from "../../assets/realTimeAnalytics.svg";
import geoFence from "../../assets/geoFence.svg";
import mapsvg from "../../assets/map.svg";

const features = [
  {
    title: "Real-Time Fleet Tracking",
    description:
      "Monitor vehicle locations, live speed, route history, and trip progress across your entire fleet from a centralized dashboard with real-time visibility and operational control.",
    visual: realTimeFleet,
    imagePosition: "left",
  },
  {
    title: "Geofencing & Instant Alerts",
    description:
      "Create custom zones around warehouses, delivery areas, or restricted locations and receive instant alerts when vehicles enter or exit those boundaries, helping improve compliance and security.",
    visual: geoFence,
    imagePosition: "right",
  },
  {
    title: "Advanced Fleet Analytics",
    description:
      "Gain insights into fuel usage, idle time, route efficiency, driver behavior, and overall fleet performance with detailed analytics, historical reports, and actionable performance metrics.",
    visual: realTimeAnalytics,
    imagePosition: "left",
  },
  {
    title: "Smart Route Optimization",
    description:
      "Automatically generate the most efficient delivery routes based on real-time traffic conditions, distance, fuel efficiency, and operational priorities to reduce delays, costs, and unnecessary mileage.",
    visual: mapsvg,
    imagePosition: "right",
  },
];

export default function Features() {
  return (
    <section
      id="features"
      className="py-16 sm:py-20 px-10 sm:px-6 lg:px-8 relative"
    >
      <div className="max-w-7xl mx-auto">
        {/* heading for cards section*/}
        <div className="text-center mb-12 sm:mb-16 lg:mb-20">
          <h2 className="text-5xl sm:text-4xl md:text-5xl lg:text-6xl font-bold mb-4 sm:mb-6">
            <span className="bg-linear-to-b from-white to-gray-300 bg-clip-text text-transparent">
              Your Complete Fleet Management
            </span>
            <br />
            <span className="bg-linear-to-b from-blue-400 to-cyan-400 bg-clip-text text-transparent">
              Workflow
            </span>
          </h2>
        </div>

        {/* actual cards start here */}
        <div className="space-y-16 sm:space-y-20 lg:space-y-32">
          {features.map((feature, key) => (
            <div
              key={key}
              className={`flex flex-col lg:flex-row items-center  gap-8 sm:gap-12 ${
                feature.imagePosition === "right" ? "lg:flex-row-reverse" : ""
              }`}
            >
              <div className="flex-1 max-w-4xl">
                <div className="relative group">
                  <div className="absolute inset-0 bg-linear-to-br from-blue-500/20 to-purple-500/20 rounded-xl sm:rounded-2xl transition-all durtaion-500" />
                  <div
                    className="relative bg-gray-900/50 backdrop-blur-sm border border-gray-700/50 
                  rounded-xl sm:rounded-2xl p-4 sm:p-6 overflow-hidden group-hover:border 
                  group-hover:border-blue-600/50 transition-all duration-300"
                  >
                    <div className="bg-gray-950 rounded-lg p-3 sm:p-4 font-mono text-xs sm:text-sm">
                      <div className="flex items-center space-x-1 sm:space-x-2 mb-3 sm:mb-4">
                        {/* rgb command buttons on the top */}
                        <div className="flex items-center space-x-1 sm:space-x-2">
                          <div className="w-2 h-2 sm:w-3 sm:h-3 rounded-full bg-red-500" />
                          <div className="w-2 h-2 sm:w-3 sm:h-3 rounded-full bg-yellow-500" />
                          <div className="w-2 h-2 sm:w-3 sm:h-3 rounded-full bg-green-500" />
                        </div>
                        <span className="text-gray-400 ml-2 sm:ml-4 text-xs sm:text-sm">
                          {feature.title}
                        </span>
                      </div>
                      {/* the Image needs to be put here */}
                      <div>
                        <img
                          src={feature.visual}
                          alt={feature.title}
                          className="w-full max-w-md mx-auto"
                        />
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              {/* card text goes here */}
              <div className="flex-1 w-full">
                <div className="max-w-lg mx-auto text-center lg:text-left">
                  <h3 className="text-4xl sm:text-3xl lg:text-4xl font-bold mb-4 sm:mb-6 text-white">
                    {feature.title}
                  </h3>
                  <p className="text-gray-300 text-base sm:text-lg lg:text-xl leading-relaxed">
                    {feature.description}
                  </p>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
