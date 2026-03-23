import instance from "../axios";
import { redirect } from "react-router-dom";

export async function fleetLoader() {
  try {
    const [allVehiclesRes, organizationsRes] = await Promise.all([
      instance.get("/v1/fleet/vehicles/all"),
      instance.get("/v1/organizations"),
    ]);

    const organizations = organizationsRes.data.data[0]?.organizations ?? [];
    const vehicles = allVehiclesRes.data.data[0]?.vehicles ?? [];

    const driverGroups = await Promise.all(
      organizations.map(async (organization) => {
        const driversRes = await instance.get("/v1/fleet/drivers", {
          params: { organization_id: organization.organization_id },
        });

        return driversRes.data.data[0]?.drivers ?? [];
      }),
    );

    return {
      organizations,
      vehicles,
      drivers: driverGroups.flat(),
    };
  } catch (error) {
    if (error.response?.status === 401) {
      throw redirect("/auth/login?reason=session_expired");
    }

    throw new Response("Failed to load fleet data", { status: 500 });
  }
}
