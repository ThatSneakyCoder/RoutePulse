import instance from "../axios";

export async function fleetLoader() {
    try {
        const [allVehicles] = await Promise.all([
            instance.get("/v1/fleet/vehicles/all"),
        ]);

        return {
            vehicles: allVehicles.data.data[0].vehicles || [],
        };
    } catch (error) {
        throw new Response("Failed to load vehicles", { status: 500 });
  }
}