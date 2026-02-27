import React from "react";
import { AccountToggle } from "./AccountToggle";
import { Search } from "./Search";
import { RouteSelect } from "./RouteSelect";

export const Sidebar = ({ collapsed, setCollapsed }) => {
  return (
    <aside
      className={`transition-all duration-300 bg-slate-950
      ${collapsed ? "w-15" : "w-64"}`}
    >
      <div className="mt-3 h-full overflow-y-auto p-3 space-y-6">
        <AccountToggle collapsed={collapsed} setCollapsed={setCollapsed} />
        {!collapsed && (
          <>
            <Search />
            <RouteSelect />
          </>
        )}
      </div>
    </aside>
  );
};
