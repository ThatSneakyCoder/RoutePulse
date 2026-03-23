import instance from "../axios";
import { redirect } from "react-router-dom";

export async function createTripAction({ request }) {
  const formData = await request.formData();
  const intent = formData.get("intent");

  if (intent !== "create-trip") {
    return null;
  }

  const organizationId = formData.get("organization_id");
  const driverId = formData.get("driver_id");
  const vehicleId = formData.get("vehicle_id");
  const startLatitude = Number.parseFloat(formData.get("start_latitude"));
  const startLongitude = Number.parseFloat(formData.get("start_longitude"));
  const endLatitude = Number.parseFloat(formData.get("end_latitude"));
  const endLongitude = Number.parseFloat(formData.get("end_longitude"));

  if (
    !organizationId ||
    !driverId ||
    !vehicleId ||
    Number.isNaN(startLatitude) ||
    Number.isNaN(startLongitude) ||
    Number.isNaN(endLatitude) ||
    Number.isNaN(endLongitude)
  ) {
    return {
      error: "Fill in the organization, driver, vehicle, and both coordinates.",
    };
  }

  try {
    await instance.post("/v1/fleet/trips", {
      organization_id: organizationId,
      driver_id: driverId,
      vehicle_id: vehicleId,
      start_latitude: startLatitude,
      start_longitude: startLongitude,
      end_latitude: endLatitude,
      end_longitude: endLongitude,
    });

    return redirect("/dashboard/trip");
  } catch (err) {
    return {
      error:
        err.response?.data?.error ||
        "We couldn't create the trip right now. Please try again.",
    };
  }
}
