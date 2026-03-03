import { ActiveFleet } from "./ActiveFleet";
import { DeliveryPerformance } from "./DeliveryPerformance";
import { DeliveryPerformanceMap } from "./DeliveryPerformanceMap";
import { OprationsOverview } from "./OprationsOverview";
import { OrderStatus } from "./OrderStatus";

export const Grid = () => {
  return (
    <div className="px-4 grid gap-3 grid-cols-1 lg:grid-cols-12 h-auto">
      {/* top row */}
      <OprationsOverview />

      {/* 2nd row (2 tiles) */}
      <DeliveryPerformance />
      <DeliveryPerformanceMap />

      {/* 3rd row (3 tiles) */}
      <ActiveFleet />
      <OrderStatus />
    </div>
  );
};
