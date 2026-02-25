import { useEffect, useState } from "react";
import {
  GitGraph,
  Play,
  ArrowRight,
  ChevronDown,
  Van,
  MapPin,
  ChartColumn,
} from "lucide-react";
import fleetMap from "../../assets/fleetMap.gif";
import analyticsGraph from "../../assets/analyticsGraph.gif";
import geofenceMap from "../../assets/geofenceMap.gif";

const tabs = {
  "Fleet Tracking": {
    image: fleetMap,
    image_alt: "Fleet Map",
    icon: <Van />,
    title: "Live Fleet Overview",
    content:
      "Track vehicle locations, speed, route progress, and status in real time. Detect delays instantly and respond before they impact operations.",
    bgColor: "bg-blue-500/20",
    iconColor: "text-blue-400",
    textColor: "text-blue-900",
    contentColor: "text-black",
  },

  Geofencing: {
    image: geofenceMap,
    image_alt: "Geofencing",
    icon: <MapPin />,
    title: "Intelligent Zone Monitoring",
    content:
      "Create custom geofences and receive instant alerts when vehicles enter or exit defined areas. Enforce compliance without manual oversight.",
    bgColor: "bg-purple-700/20",
    iconColor: "text-purple-900",
    textColor: "text-purple-900",
    contentColor: "text-black",
  },

  Analytics: {
    image: analyticsGraph,
    image_alt: "Analytics",
    icon: <ChartColumn />,
    title: "Operational Performance Dashboard",
    content:
      "Analyze trip history, idle time, fuel efficiency, and route optimization metrics to reduce costs and improve delivery reliability.",
    bgColor: "bg-emerald-500/20",
    iconColor: "text-emerald-400",
    textColor: "text-emerald-400",
    contentColor: "text-emerald-400",
  },
};

export default function Hero() {
  const [mousePosition, setMousePosition] = useState({ x: 0, y: 0 });
  const [activeTab, setActiveTab] = useState("Fleet Tracking");

  useEffect(() => {
    function handleMouseMove(e) {
      setMousePosition({ x: e.clientX, y: e.clientY });
    }

    window.addEventListener("mousemove", handleMouseMove);

    return () => window.removeEventListener("mousemove", handleMouseMove);
  }, []);

  const currentFloatingCard = tabs[activeTab];
  return (
    <section className="relative min-h-screen flex items-center justify-center pt-16 sm:pt-20 px-4 sm:px-6 lg:px-8 overflow-hidden">
      {/* div that follows mouse */}
      <div
        className="absolute inset-0 opacity-30"
        style={{
          background: `radial-gradient(600px circle at ${mousePosition.x}px ${mousePosition.y}px, rgba(59, 130, 246, 0.15), transparent 40%)`,
        }}
      />

      {/* divs for decoration */}
      <div className="absolute top-20 left-4 sm:left-10 w-48 sm:w-72 h-48 sm:h-72 bg-blue-500/10 rounded-full blur-3xl animate-pulse" />
      <div className="absolute bottom-20 right-4 sm:right-10 w-64 sm:w-96 h-64 sm:h-96 bg-cyan-500/10 rounded-full blur-3xl animate-pulse delay-1000" />

      <div className="max-w-7xl mx-auto text-center relative w-full">
        <div className="max-w-7xl mx-auto flex flex-col lg:grid lg:grid-cols-2 text-center lg:text-left gap-6 sm:gap-8 lg:gap-12 items-center relative">
          {/* left side of hero section */}
          <div>
            <div className="inline-flex items-center space-x-2 px-3 sm:px-4 py-2 bg-blue-500/10 border border-blue-500/20 rounded-full mb-4 sm:mb-6 animate-in slide-in-from-bottom duration-700">
              <GitGraph className="w-4 h-4 text-blue-400" />
              <span className="text-xs sm:text-sm text-blue-300">
                Smart Fleet Operations Platform
              </span>
            </div>

            <h1 className="text-4xl sm:text-4xl md:text-4xl lg:text-4xl xl:text-5xl font-semibold mb-4 sm:mb-6 animate-in slide-in-from-bottom delay-300 duration-700 leading-tight">
              <span className="bg-gradient-to-r from-white via-blue-100 to-cyan-100 bg-clip-text text-transparent block mb-1 sm:mb-2">
                Real-Time Tracking
              </span>
              <span className="bg-gradient-to-b from-blue-400 via-cyan-400 to-blue-400 bg-clip-text text-transparent block mb-1 sm:mb-2">
                Smarter Fleet Control
              </span>
              <span className="bg-gradient-to-r from-white via-blue-100 to-cyan-100 bg-clip-text text-transparent block mb-1 sm:mb-2">
                Powered by RoutePulse
              </span>
            </h1>

            <p className="text-md sm:text-base lg:text-lg text-gray-400 max-w-2xl mx-auto lg:mx-0 mb-6 sm:mb-8 animate-in slide-in-from-bottom delay-500 duration-700 delay-200 leading-relaxed">
              RoutePulse helps logistics teams monitor vehicles in real time,
              enforce geofencing rules, and gain actionable insights across
              their fleet — all from a single, intuitive dashboard.
            </p>
            <div className="flex flex-col sm:flex-row items-center justify-center lg:justify-start gap-3 sm:gap-4 mb-8 sm:mb-12 animate-in slide-in-from-bottom duration-700 delay-700">
              <button className="group w-full sm:w-auto px-6 sm:px-8 py-3 sm:py-4 bg-gradient-to-b from-blue-600 to-blue-400 rounded-lg font-semibold text-sm sm:text-base transition-all duration-300 hover:scale-102 flex items-center justify-center space-x-2">
                <span>Start Monitoring</span>
                <ArrowRight className="w-4 h-4 sm:w-5 sm:h-5 group-hover:translate-x-1 transition-transform duration-300" />
              </button>

              <button className="group w-full sm:w-auto px-6 sm:px-8 py-3 sm:py-4 bg-white/5 backdrop-blur-sm border border-white/10 rounded-lg font-semibold text-sm sm:text-base transition-all duration-300 hover:bg-white/10 flex items-center justify-center space-x-2">
                <div className="p-2 bg-white/10 rounded-full  duration-700 group-hover:rotate-360">
                  <Play className="w-4 h-4 sm:w-5 sm:h-5 fill-white" />
                </div>
                <span>Watch Demo</span>
              </button>
            </div>
          </div>

          {/* right side of the hero section */}
          <div className="relative order-2 w-full">
            <div className="relative bg-white/5 rounded-xl sm:rounded-2xl p-3 sm:p-4 shadow-xl backdrop-blur-xl border border-white/10">
              <div className="flex flex-col bg-gradient-to-br from-gray-900/20 to-gray-800/20 backdrop-blur-sm rounded-lg overflow-hidden h-[400px] sm:h-[450px] lg:h-[450px] border border-white/10">
                {/* header for this section*/}
                <div className="flex items-center justify-between px-3 sm:px-4 py-2 sm:py-3 bg-white/5">
                  <div className="flex items-center space-x-3">
                    <div className="flex items-center space-x-1 sm:space-x-2">
                      <div className="w-2 h-2 sm:w-3 sm:h-3 rounded-full bg-red-500" />
                      <div className="w-2 h-2 sm:w-3 sm:h-3 rounded-full bg-yellow-500" />
                      <div className="w-2 h-2 sm:w-3 sm:h-3 rounded-full bg-green-500" />
                    </div>
                    <span className="text-xs sm:text-sm text-gray-300">
                      CodeFlow AI
                    </span>
                  </div>
                  <ChevronDown className="w-3 h-3 sm:w-4 sm:h-4 text-gray-400" />
                </div>

                {/* map image section */}
                <div className="p-3 sm:p-4 relative">
                  <div className="flex space-x-1 sm:space-x-2 mb-3 sm:mb-4 overflow-x-auto ">
                    {/* fleet tracking on map */}
                    <button
                      onClick={() => setActiveTab("Fleet Tracking")}
                      className={`px-3 py-2 backdrop-blur-sm tex-xs sm:text-sm rounded-t-lg border ${
                        activeTab === "Fleet Tracking"
                          ? "bg-blue-500/30 text-white border-blue-400/20"
                          : "bg-white/5 text-gray-300 border-white/10 hover:bg-white/10"
                      }  transition-all duration-200 whitespace-nowrap`}
                    >
                      Fleet Tracking
                    </button>
                    {/* geofencing enforcement image */}
                    <button
                      onClick={() => setActiveTab("Geofencing")}
                      className={`px-3 py-2 backdrop-blur-sm text-xs sm:text-sm rounded-t-lg border ${
                        activeTab === "Geofencing"
                          ? "bg-blue-500/30 text-white border-blue-400/20"
                          : "bg-white/5 text-gray-300 border-white/10 hover:bg-white/10"
                      } border-white/10 transition-all duration-200 whitespace-nowrap`}
                    >
                      Geofencing
                    </button>
                    <button
                      onClick={() => setActiveTab("Analytics")}
                      className={`px-3 py-2 backdrop-blur-sm tex-xs sm:text-sm rounded-t-lg border ${
                        activeTab === "Analytics"
                          ? "bg-blue-500/30 text-white border-blue-400/20"
                          : "bg-white/5 text-gray-300 border-white/10 hover:bg-white/10"
                      }  transition-all duration-200 whitespace-nowrap`}
                    >
                      Analytics
                    </button>
                  </div>
                </div>

                {/* Images */}
                <div className="relative flex-grow">
                  <img
                    src={currentFloatingCard.image}
                    alt={currentFloatingCard.image_alt}
                    className="w-full h-full object-cover"
                  />
                </div>
              </div>
            </div>

            {/* floating cards of right side */}
            <div
              key={activeTab}
              className={`hidden lg:block absolute bottom-4 right-4 transform translate-x-8 translate-y-8 w-72 ${currentFloatingCard.bgColor} backdrop-blur-xl rounded-lg p-4 border border-white/20 shadow-2xl   animate-in slide-in-from-bottom duration-700`}
            >
              <div className="flex items-center space-x-2 mb-2">
                <div
                  className={`w-6 h-6 ${currentFloatingCard.iconColor} flex items-center justify-center text-sm font-bold sm:w-8 sm:h-8`}
                >
                  {currentFloatingCard.icon}
                </div>
                <span
                  className={`text-sm font-medium ${currentFloatingCard.textColor}`}
                >
                  Live Fleet Overview
                </span>
              </div>

              <div
                className={`text-sm text-left ${currentFloatingCard.contentColor}`}
              >
                {currentFloatingCard.content}
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
