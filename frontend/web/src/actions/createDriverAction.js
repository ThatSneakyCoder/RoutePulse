import instance from "../axios";

export async function createDriverAction({ request }) {
  const formData = await request.formData();
  const intent = formData.get("intent");

  if (intent !== "create-driver") {
    return null;
  }

  const organizationId = formData.get("organization_id");
  const firstName = formData.get("first_name");
  const lastName = formData.get("last_name");
  const vehicleId = formData.get("vehicle_id");

  if (!organizationId || !firstName || !lastName) {
    return {
      error: "Organization, first name, and last name are required.",
    };
  }

  try {
    await instance.post("/v1/fleet/drivers", {
      organization_id: organizationId,
      first_name: firstName,
      last_name: lastName,
      vehicle_id: vehicleId || undefined,
    });

    return null;
  } catch (err) {
    return {
      error:
        err.response?.data?.error ||
        "We couldn't create the driver right now. Please try again.",
    };
  }
}
