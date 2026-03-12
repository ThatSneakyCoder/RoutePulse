import { useRouteError } from "react-router-dom";

export const DashboardError = () => {
  const error = useRouteError();

  return (
    <div className="p-8 text-red-400">
      <h2 className="text-lg font-semibold">Failed to load page</h2>

      <p className="text-sm mt-2 text-slate-400">
        {error?.statusText || error?.message || "Unknown error"}
      </p>
    </div>
  );
};
