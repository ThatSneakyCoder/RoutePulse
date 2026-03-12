import { redirect } from "react-router";
import instance from "../axios";

export async function organizationMutationAction({ request }) {
  const formData = await request.formData();
  const intent = formData.get("intent");

   if (intent === "create-org") {
    const name = formData.get("name");
    const res = await instance.post("/v1/organizations", {
        name,
    });

    const org = res.data.data[0].organization;

    return redirect(`/dashboard/organization/${org.organization_id}`);
   }

  return null;
//   we can also return data and use it with useActionData
}