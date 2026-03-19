import instance from "../axios";

export async function analyticsLoader() {
  try {
    const res = await instance.get("/v1/analytics/total-members");

    return {
      totalMembers: res.data.data[0].count,
    };
  } catch (error) {
    console.error("Analytics loader error:", error);
    throw error;
  }
}