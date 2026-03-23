import instance from "../axios";

export async function createVehicleAction({ request }) {
  const formData = await request.formData();
  const intent = formData.get("intent");

  if (intent !== "create-vehicle") {
    return null;
  }

  const organizationId = formData.get("organization_id");
  const plateNumber = formData.get("plate_number");
  const vehicleType = formData.get("vehicle_type");
  const capacity = Number.parseInt(formData.get("capacity"), 10);

  if (!organizationId || !plateNumber || !vehicleType || Number.isNaN(capacity)) {
    return {
      error: "Organization, plate number, type, and capacity are required.",
    };
  }

  try {
    await instance.post("/v1/fleet/vehicles", {
      organization_id: organizationId,
      plate_number: plateNumber,
      vehicle_type: vehicleType,
      capacity,
    });

    return null;
  } catch (err) {
    return {
      error:
        err.response?.data?.error ||
        "We couldn't create the vehicle right now. Please try again.",
    };
  }
}
