import { EllipsisVertical, Icon, Package } from "lucide-react";

export const OprationsOverview = () => {
  const stats = [
    {
      color: "green",
      title: "Total shipments today",
      value: "1,284",
      description: "From all shipments",
      delta: "+12.5%",
    },
    {
      color: "green",
      title: "Active trucks",
      value: "482",
      description: "Currently on route",
      delta: "+12.5%",
    },
    {
      color: "green",
      title: "Avg. Delivery time",
      value: "4h 12m",
      description: "Faster than avg",
      delta: "-18m",
    },
    {
      color: "red",
      title: "On-time delivery rate",
      value: "98.2%",
      description: "Target: 99.0%",
      delta: "-0.4%",
    },
  ];
  return (
    <div className="p-4 col-span-1 lg:col-span-12 rounded-xl border border-slate-700 shadow-2xl">
      <div className="flex flex-col">
        <div className="flex items-center justify-between">
          <span className="font-semibold">Operations overview</span>
          <div className="flex items-center justify-center gap-2 text-slate-500 text-sm">
            <span>Last Updated: Today, 14:32PM</span>
            <EllipsisVertical className="w-4 h-4" />
          </div>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 mt-6">
          {stats.map((stat, index) => {
            return <Stat key={index} {...stat} />;
          })}
        </div>
      </div>
    </div>
  );
};

const Stat = ({
  color = "green",
  title,
  value,
  description,
  delta,
}) => {
  const isPositive = color === "green";

  const accentColor = isPositive
    ? "border-green-400"
    : "border-red-400";

  const badgeStyles = isPositive
    ? "bg-green-500/20 text-green-400"
    : "bg-red-500/20 text-red-400";

  return (
    <div
      className={`h-full border-l-2 ${accentColor} px-4 flex flex-col mt-2`}
    >
      {/* Title */}
      <div className="flex items-center gap-2 text-sm text-slate-400 mb-5">
        <Package className="w-4 h-4" />
        <span className="whitespace-nowrap">
          {title}
        </span>
      </div>

      {/* Value Row */}
      <div className="flex items-center justify-between">
        <span className="text-3xl font-semibold">
          {value}
        </span>

        <div
          className={`${badgeStyles} text-xs font-medium px-3 py-1 rounded-full whitespace-nowrap`}
        >
          {delta}
        </div>
      </div>

      {/* Description Row */}
      <div className="mt-1 flex items-center justify-between text-sm text-slate-500">
        <span>{description}</span>
        <span className="whitespace-nowrap">
          vs yesterday
        </span>
      </div>
    </div>
  );
};
