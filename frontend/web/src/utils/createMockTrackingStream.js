const DEFAULT_ROUTE_POINTS = [
  { latitude: 37.7749, longitude: -122.4194 },
  { latitude: 37.7754, longitude: -122.4194 },
  { latitude: 37.7760, longitude: -122.4194 },
  { latitude: 37.7760, longitude: -122.4185 },
  { latitude: 37.7760, longitude: -122.4177 },
  { latitude: 37.7767, longitude: -122.4177 },
  { latitude: 37.7773, longitude: -122.4177 },
  { latitude: 37.7773, longitude: -122.4168 },
  { latitude: 37.7773, longitude: -122.4160 },
  { latitude: 37.7780, longitude: -122.4160 },
  { latitude: 37.7787, longitude: -122.4160 },
  { latitude: 37.7787, longitude: -122.4152 },
  { latitude: 37.7793, longitude: -122.4152 },
];

function normalizeRoutePoints(routePoints) {
  if (Array.isArray(routePoints) && routePoints.length > 1) {
    const validPoints = routePoints.filter((point) => (
      typeof point?.latitude === "number" && typeof point?.longitude === "number"
    ));

    if (validPoints.length > 1) {
      return validPoints;
    }
  }

  return DEFAULT_ROUTE_POINTS;
}

function interpolateRoutePoints(points) {
  const interpolatedPoints = [];

  for (let index = 0; index < points.length - 1; index += 1) {
    const currentPoint = points[index];
    const nextPoint = points[index + 1];
    const latitudeDelta = nextPoint.latitude - currentPoint.latitude;
    const longitudeDelta = nextPoint.longitude - currentPoint.longitude;
    const segmentDistance = Math.max(Math.abs(latitudeDelta), Math.abs(longitudeDelta));
    const steps = Math.max(6, Math.ceil(segmentDistance / 0.00008));

    for (let step = 0; step < steps; step += 1) {
      const progress = step / steps;

      interpolatedPoints.push({
        latitude: currentPoint.latitude + (latitudeDelta * progress),
        longitude: currentPoint.longitude + (longitudeDelta * progress),
      });
    }
  }

  interpolatedPoints.push(points.at(-1));

  return interpolatedPoints;
}

export function createMockTrackingStream({
  tripId,
  driverId,
  vehicleId,
  routePoints,
}) {
  const points = interpolateRoutePoints(normalizeRoutePoints(routePoints));
  let currentIndex = 0;

  return function getNextTrackingMessage() {
    const currentPoint = points[currentIndex % points.length];
    currentIndex += 1;

    return {
      type: "driver_location_update",
      data: {
        trip_id: tripId,
        driver_id: driverId ?? "simulated-driver",
        vehicle_id: vehicleId ?? "simulated-vehicle",
        latitude: currentPoint.latitude,
        longitude: currentPoint.longitude,
        recorded_at: new Date().toISOString(),
        sequence: currentIndex,
      },
    };
  };
}
