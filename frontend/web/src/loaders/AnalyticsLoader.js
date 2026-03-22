import instance from "../axios";

export async function analyticsLoader() {
  try {
    const [membersRes, activeUsersRes, activityRes] = await Promise.all([
      instance.get("/v1/analytics/total-members"),
      instance.get("/v1/analytics/active-users-today"),
      instance.get("/v1/analytics/recent-activity"),
    ]);

    return {
      totalMembers: membersRes.data.data[0].count,
      activeUsersToday: activeUsersRes.data.data[0].count,
      recentActivity: activityRes.data.data[0].events, // assuming array
    };
  } catch (error) {
    console.error("Analytics loader error:", error);
    throw error;
  }
}