import instance from "../axios";
import { redirect } from "react-router-dom";

export async function tripLoader() {
  try {
    const [organizationsRes, tripsRes] = await Promise.all([
      instance.get("/v1/organizations"),
      instance.get("/v1/fleet/trips/all"),
    ]);

    const organizations = organizationsRes.data.data[0]?.organizations ?? [];
    const trips = tripsRes.data.data[0]?.trips ?? [];

    const fleetData = await Promise.all(
      organizations.map(async (organization) => {
        const organizationId = organization.organization_id;

        const [driversRes, vehiclesRes] = await Promise.all([
          instance.get("/v1/fleet/drivers", {
            params: { organization_id: organizationId },
          }),
          instance.get("/v1/fleet/vehicles", {
            params: { organization_id: organizationId },
          }),
        ]);

        return {
          drivers: driversRes.data.data[0]?.drivers ?? [],
          vehicles: vehiclesRes.data.data[0]?.vehicles ?? [],
        };
      }),
    );

    return {
      organizations,
      drivers: fleetData.flatMap((entry) => entry.drivers),
      vehicles: fleetData.flatMap((entry) => entry.vehicles),
      trips,
    };
  } catch (err) {
    if (err.response?.status === 401) {
      throw redirect("/auth/login?reason=session_expired");
    }

    throw new Response("Failed to load trip data", { status: 500 });
  }
}
