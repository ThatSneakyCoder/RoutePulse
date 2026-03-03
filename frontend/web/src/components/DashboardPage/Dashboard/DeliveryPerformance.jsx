import { CalendarDays, ChevronDown, Dot } from "lucide-react";
import {
  CartesianGrid,
  Legend,
  Line,
  LineChart,
  Tooltip,
  XAxis,
  YAxis,
  ResponsiveContainer,
} from "recharts";

export const DeliveryPerformance = () => {
  const data = [
    { time: "06:00", delivered: 20, delayed: 5 },
    { time: "08:00", delivered: 45, delayed: 10 },
    { time: "10:00", delivered: 130, delayed: 18 },
    { time: "12:00", delivered: 230, delayed: 30 },
    { time: "14:00", delivered: 340, delayed: 55 },
    { time: "16:00", delivered: 470, delayed: 75 },
    { time: "18:00", delivered: 580, delayed: 80 },
  ];

  return (
    <div className="p-4 col-span-1 lg:col-span-5 rounded-xl border border-slate-700 shadow-2xl h-auto w-full flex flex-col gap-4">
      <div className="flex items-center justify-between">
        <span className="font-semibold">Delivery performance</span>
        <div className="flex items-center gap-2 border border-slate-600 rounded-lg py-2 px-4 hover:bg-slate-600/40 transition-colors duration-300">
          <CalendarDays className="w-4 h-4" />
          <span>Weekly</span>
          <ChevronDown className="w-4 h-4" />
        </div>
      </div>
      <div className="border border-slate-700 flex items-center rounded-lg overflow-hidden">
        <div className="flex-1 flex items-center justify-center">
          <Dot className="w-8 h-8 text-green-400" />
          <span>Delivered</span>
        </div>
        {/* divider line */}
        <div className="w-px h-full bg-slate-700" />

        <div className="flex-1 flex items-center justify-center">
          <Dot className="w-8 h-8 text-yellow-400" />
          <span>Delayed</span>
        </div>
      </div>
      {/* graph */}
      <div className="w-full h-100 p-4">
        <ResponsiveContainer width="100%" height="100%">
          <LineChart
            width={600}
            height={300}
            data={data}
            margin={{ top: 0, right: 15, bottom: 5, left: -30 }}
          >
            <CartesianGrid stroke="#334155" strokeDasharray="3 3" />

            <XAxis
              dataKey="time"
              stroke="#94a3b8"
              tick={{ fill: "#94a3b8", fontSize: 12 }}
            />

            <YAxis stroke="#94a3b8" tick={{ fill: "#94a3b8", fontSize: 12 }} />

            <Tooltip
              contentStyle={{
                backgroundColor: "#0f172a",
                border: "1px solid #334155",
                borderRadius: "8px",
              }}
            />

            <Line
              type="monotone"
              dataKey="delivered"
              stroke="#22c55e"
              fill="#22c55e"
              strokeWidth={2}
              dot={{ r: 4 }}
              name="Delivered"
            />

            <Line
              type="monotone"
              dataKey="delayed"
              stroke="#facc15"
              fill="#facc15"
              strokeWidth={2}
              dot={{ r: 4 }}
              name="Delayed"
            />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
};
