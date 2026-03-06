import { useState } from "react";
import { Dashboard } from "./Dashboard/Dashboard";
import { Sidebar } from "./Sidebar/Sidebar";

export const RootDashboard = () => {
  const [collapsed, setCollapsed] = useState(false);

  return (
    <div className="flex h-screen bg-slate-950">
      <Sidebar collapsed={collapsed} setCollapsed={setCollapsed} />

      <main className="flex-1 overflow-y-auto">
        <Dashboard />
      </main>
    </div>
  );
};
