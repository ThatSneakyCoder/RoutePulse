import axios from "../axios";
import { redirect } from "react-router-dom";

export async function organizationDetailsLoader({ params }) {
  const { orgId } = params;

  try {
    const [orgRes, membersRes] = await Promise.all([
      axios.get(`/v1/organizations/${orgId}`),
      axios.get(`/v1/organizations/${orgId}/members`)
    ]);

    const organization = orgRes.data.data[0];

    const members =
      membersRes.data.data[0]?.members?.map((m) => ({
        id: m.user_id,
        first_name: m.first_name,
        last_name: m.last_name,
        email: m.email,
        role: m.role,
        joined_at: new Date(m.joined_at * 1000),
      })) ?? [];

    return {
      organization,
      members,
    };
  } catch (err) {
    if (err.response?.status === 401) {
      throw redirect("/auth/login?reason=session_expired");
    }

    throw new Response("Failed to load organization", { status: 500 });
  }
}