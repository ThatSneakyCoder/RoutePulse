import React from "react";

export const HealthCheck = () => {
  return (
    <section className="min-h-screen bg-slate-950 text-white flex items-center justify-center p-6 sm:p-8 lg:px-8 overflow-hidden">
      <div className="max-w-7xl flex items-center justify-center">
        <div className="h-full w-full border border-gray-600 backdrop-blur-2xl bg-slate-700/40 rounded-lg py-4 px-2 flex items-center justify-center gap-2 text-green-500">
          <span className="text-5xl ">Health Check Successful</span>
          <span className="text-5xl  transition-transform animate-pulse">.</span>
          <span className="text-5xl animate-pulse delay-300">.</span>
          <span className="text-5xl animate-pulse delay-500">.</span>
        </div>
      </div>
    </section>
  );
};
