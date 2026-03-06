import { PanelRightOpen } from "lucide-react";

export const AccountToggle = ({ collapsed, setCollapsed }) => {
  return (
    <div className="pb-1">
      <div className="flex items-center justify-between">
        {!collapsed && (
          <span className="text-lg font-semibold">RoutePulse</span>
        )}

        <button
          onClick={() => setCollapsed(!collapsed)}
          className="p-2 rounded-md bg-slate-800 hover:bg-slate-700 transition"
        >
          <PanelRightOpen className="w-5 h-5" />
        </button>
      </div>
    </div>
  );
};
