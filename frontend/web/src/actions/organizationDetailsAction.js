import instance from "../axios";

export async function organizationDetailsAction({ request, params }) {
  const { orgId } = params;

  const formData = await request.formData();
  const intent = formData.get("intent");

  try {
    if (intent === "invite-member") {
      const email = formData.get("email");
      const role = formData.get("role");

      await instance.post(`/v1/organizations/${orgId}/invite`, {
        email,
        role,
      });

      return null;
    }

    if (intent === "remove-member") {
      const userId = formData.get("userId");

      await instance.delete(
        `/v1/organizations/${orgId}/members/${userId}`
      );

      return null;
    }

    if (intent === "update-role") {
      const userId = formData.get("userId");
      const role = formData.get("role");

      await instance.put(
        `/v1/organizations/${orgId}/members/${userId}/role`,
        { role }
      );

      return null;
    }

    return null;
  } catch (err) {
    console.error("Organization action failed", err);
    throw err;
  }
}