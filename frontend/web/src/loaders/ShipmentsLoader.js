import instance from "../axios";

export async function shipmentsLoader() {
  try {
    const response = await instance.get("/v1/fleet/trips/all");

    const trips = response.data.data[0]?.trips || [];

    const activeShipments = trips.filter(
      (t) => t.status === "active" || t.status === "created"
    );

    const completedShipments = trips.filter(
      (t) => t.status === "completed"
    );

    return {
      activeShipments,
      completedShipments,
    };

  } catch (error) {
    throw new Response("Failed to load shipments", { status: 500 });
  }
}