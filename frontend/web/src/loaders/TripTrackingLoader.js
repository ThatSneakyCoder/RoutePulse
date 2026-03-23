import instance from "../axios";
import { redirect } from "react-router-dom";

export async function tripTrackingLoader({ params }) {
  const { tripId } = params;

  try {
    const [currentRes, historyRes, geometryRes] = await Promise.all([
      instance
        .get(`/v1/tracking/trips/${tripId}/current`)
        .catch((err) => {
          if (err.response?.status === 404) {
            return null;
          }

          throw err;
        }),
      instance.get(`/v1/tracking/trips/${tripId}/history`, {
        params: { limit: 200 },
      }),
      instance.get(`/v1/tracking/trips/${tripId}/geometry`),
    ]);

    return {
      currentLocation: currentRes?.data?.data?.[0] ?? null,
      locationHistory: historyRes.data.data[0] ?? { trip_id: tripId, points: [] },
      geometry: geometryRes.data.data[0] ?? {
        trip_id: tripId,
        planned_geometry: [],
        actual_geometry: [],
      },
    };
  } catch (err) {
    if (err.response?.status === 401) {
      throw redirect("/auth/login?reason=session_expired");
    }

    throw new Response("Failed to load live trip tracking data", { status: 500 });
  }
}
