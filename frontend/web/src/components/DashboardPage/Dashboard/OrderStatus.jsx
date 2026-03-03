import { EllipsisVertical } from "lucide-react";
import { Pie, PieChart, ResponsiveContainer, Cell } from "recharts";

export const OrderStatus = () => {
  const orderStatusData = [
    { name: "On the way", value: 642, color: "#065f46" },
    { name: "Delivered", value: 257, color: "#16a34a" },
    { name: "Delayed", value: 193, color: "#4ade80" },
    { name: "Waiting", value: 128, color: "#14b8a6" },
    { name: "Canceled", value: 64, color: "#0284c7" },
  ];

  return (
    <div className="p-4 col-span-1 lg:col-span-6 rounded-xl border border-slate-700 shadow-2xl h-[420px] w-full flex flex-col gap-4">
      
      {/* title */}
      <div className="flex items-center justify-between">
        <span className="font-semibold">Current loading status</span>
        <div className="text-slate-500 text-sm transition-colors duration-300 hover:text-slate-200">
          <EllipsisVertical className="w-4 h-4" />
        </div>
      </div>

      {/* content */}
      <div className="flex gap-4 flex-1">
        
        {/* chart */}
        <div className="w-3/5 h-full p-4">
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie
                data={orderStatusData}
                innerRadius={70}
                outerRadius={100}
                paddingAngle={4}
                cornerRadius={6}
                dataKey="value"
                isAnimationActive
              >
                {orderStatusData.map((entry, index) => (
                  <Cell key={index} fill={entry.color} />
                ))}
              </Pie>
            </PieChart>
          </ResponsiveContainer>
        </div>

        {/* legend */}
        <div className="w-2/5 h-full p-4 mr-10 flex flex-col justify-center gap-3">
          {orderStatusData.map((item, index) => (
            <div key={index} className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <span
                  className="w-3 h-3 rounded-full"
                  style={{ backgroundColor: item.color }}
                />
                <span className="text-lg text-slate-300">{item.name}</span>
              </div>
              <span className="text-lg font-medium">{item.value}</span>
            </div>
          ))}
        </div>

      </div>
    </div>
  );
};