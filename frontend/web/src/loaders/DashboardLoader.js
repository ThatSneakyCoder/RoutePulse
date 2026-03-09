import axios from "../axios";
import { redirect } from "react-router-dom";

export async function dashboardLoader() {
  try {
    const [userRes, orgRes] = await Promise.all([
      axios.get("/v1/user/me"),
      axios.get("/v1/organizations"),
    ]);

    const user = userRes.data.data[0];
    const orgList = orgRes.data.data[0]?.organizations ?? [];

    const organizations = orgList.map((org) => ({
      organization_id: org.organization_id,
      name: org.name,
      created_at: new Date(org.created_at * 1000),
      members: 0,
      vehicles: 0,
      drivers: 0,
    }));

    return {
      user,
      organizations,
    };
  } catch (err) {
    if (err.response?.status === 401) {
      throw redirect("/auth/login?reason=session_expired");
    }

    throw new Response("Failed to load dashboard data", { status: 500 });
  }
}