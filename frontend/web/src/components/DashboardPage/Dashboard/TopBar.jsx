import { Bell, ChevronDown, MessageCircleMore, Plus, User } from "lucide-react";
import { useState } from "react";
import { MyDatePicker } from "./MyDayPicker";

export const TopBar = () => {
  const [isDateBoxOpened, setDateBoxOpened] = useState(false);

  return (
    <div className="flex items-center justify-between">
      
      <h1 className="text-lg font-semibold text-slate-100">
        Dashboard
      </h1>

      <div className="flex items-center gap-3">

        {/* Date Picker */}
        <div className="relative">
          <button
            onClick={() => setDateBoxOpened(!isDateBoxOpened)}
            className="flex items-center gap-2 px-3 py-2 text-sm font-medium 
                       bg-slate-800 border border-slate-700 
                       rounded-lg text-slate-200
                       hover:bg-slate-700 transition"
          >
            <span>Jan - Nov 2025</span>
            <ChevronDown
              className={`w-4 h-4 transition-transform ${
                isDateBoxOpened ? "rotate-180" : ""
              }`}
            />
          </button>

          {isDateBoxOpened && (
            <div className="absolute right-0 mt-2 
                            bg-slate-800 border border-slate-700 
                            rounded-lg p-3 shadow-lg z-50">
              <MyDatePicker />
            </div>
          )}
        </div>

        {/* Icon Buttons */}
        <IconButton>
          <MessageCircleMore className="w-4 h-4" />
        </IconButton>

        <IconButton>
          <Bell className="w-4 h-4" />
        </IconButton>

        <IconButton>
          <User className="w-4 h-4" />
        </IconButton>

        {/* Divider */}
        <div className="h-6 w-px bg-slate-700" />

        {/* Primary Action */}
        <button className="flex items-center gap-2 px-4 py-2 
                           bg-blue-600 hover:bg-blue-500 
                           text-white text-sm font-medium 
                           rounded-lg transition">
          <Plus className="w-4 h-4" />
          Quick Shipment
        </button>

      </div>
    </div>
  );
};

const IconButton = ({ children }) => {
  return (
    <button
      className="p-2 bg-slate-800 border border-slate-700 
                 rounded-lg text-slate-300
                 hover:bg-slate-700 hover:text-white
                 transition"
    >
      {children}
    </button>
  );
};