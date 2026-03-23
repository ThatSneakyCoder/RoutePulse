import instance from "../axios";
import { redirect } from "react-router-dom";

export async function tripOperationsAction({ request }) {
  const formData = await request.formData();
  const intent = formData.get("intent");
  const tripId = formData.get("trip_id");

  if (!tripId || !intent) {
    return null;
  }

  try {
    if (intent === "start-trip") {
      await instance.post(`/v1/fleet/trips/${tripId}/start`, {
        trip_id: tripId,
      });
      return redirect(`/dashboard/trip/${tripId}/live`);
    }

    if (intent === "complete-trip") {
      await instance.post(`/v1/fleet/trips/${tripId}/complete`, {
        trip_id: tripId,
      });
      return redirect("/dashboard/trip/driver-console");
    }

    return null;
  } catch (err) {
    return {
      error:
        err.response?.data?.error ||
        "We couldn't update the trip status right now. Please try again.",
    };
  }
}
