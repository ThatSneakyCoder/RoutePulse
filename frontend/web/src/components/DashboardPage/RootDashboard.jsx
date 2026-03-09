import { useState } from "react";
import { Sidebar } from "./Sidebar/Sidebar";
import { Outlet } from "react-router-dom";

export const RootDashboard = () => {
  const [collapsed, setCollapsed] = useState(false);

  return (
    <div className="flex h-screen bg-slate-950">

      <Sidebar
        collapsed={collapsed}
        setCollapsed={setCollapsed}
      />

      <main className="flex-1 overflow-y-auto">
        <Outlet />
      </main>

    </div>
  );
};