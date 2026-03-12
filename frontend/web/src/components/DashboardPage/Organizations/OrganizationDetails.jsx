import { useLoaderData, Form } from "react-router-dom";
import { Users, Plus, Trash2 } from "lucide-react";

export const OrganizationDetails = () => {
  const { organization, members } = useLoaderData();

  return (
    <section className="p-8 space-y-8">
      {/* Header */}
      <div>
        <h1 className="text-xl font-semibold text-white">
          {organization.name}
        </h1>

        <p className="text-xs text-slate-400">
          ID: {organization.organization_id}
        </p>

        <p className="text-sm text-slate-400 mt-1">
          Manage organization members, fleet and permissions
        </p>
      </div>

      {/* Organization Stats */}
      <div className="bg-slate-900 border border-slate-800 rounded-xl p-6">
        <h2 className="text-white font-medium mb-4">Organization Info</h2>

        <div className="grid grid-cols-3 gap-6 text-sm">
          <Stat label="Members" value={members.length} />
          <Stat label="Vehicles" value="24" />
          <Stat label="Drivers" value="8" />
        </div>
      </div>

      {/* Invite Member */}
      <div className="bg-slate-900 border border-slate-800 rounded-xl p-8">
        <div className="flex items-center gap-2 text-white font-medium mb-4">
          <Plus className="w-4 h-4" />
          Invite User
        </div>

        <Form method="post" className="flex gap-3 flex-wrap">
          <input
            name="email"
            placeholder="user@email.com"
            required
            className="bg-slate-800 border border-slate-700 rounded-md px-3 py-2 text-sm text-white w-64"
          />

          <select
            name="role"
            defaultValue="Dispatcher"
            className="bg-slate-800 border border-slate-700 rounded-md px-2 py-2 text-sm text-white"
          >
            <option>Dispatcher</option>
            <option>Ops</option>
            <option>Owner</option>
          </select>

          <input type="hidden" name="intent" value="invite-member" />

          <button
            type="submit"
            className="px-4 py-2 text-sm font-medium bg-blue-600 hover:bg-blue-500 text-white rounded-md transition"
          >
            Send Invite
          </button>
        </Form>
      </div>

      {/* Members */}
      <div className="bg-slate-900 border border-slate-800 rounded-xl p-8">
        <div className="flex items-center gap-2 text-white font-medium mb-6">
          <Users className="w-4 h-4" />
          Organization Members
        </div>

        <div className="grid grid-cols-[2fr_1fr_1fr_auto] text-xs text-slate-500 mb-3 px-4">
          <div>Member</div>
          <div>Role</div>
          <div>Joined</div>
          <div></div>
        </div>

        <div className="space-y-3">
          {members.map((m) => (
            <div
              key={m.id}
              className="grid grid-cols-[2fr_1fr_1fr_auto] items-center bg-slate-800 border border-slate-700 rounded-md px-4 py-4"
            >
              <div className="flex flex-col">
                <span className="text-white font-medium">
                  {m.first_name} {m.last_name}
                </span>

                <span className="text-xs text-slate-400">
                  {m.email}
                </span>
              </div>

              <RoleBadge role={m.role} />

              <div className="text-xs text-slate-400">
                {new Date(m.joined_at).toLocaleDateString()}
              </div>

              <Form method="post">
                <input
                  type="hidden"
                  name="intent"
                  value="remove-member"
                />

                <input
                  type="hidden"
                  name="userId"
                  value={m.id}
                />

                <button
                  type="submit"
                  className="text-red-400 hover:text-red-300 transition"
                >
                  <Trash2 size={16} />
                </button>
              </Form>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};

const RoleBadge = ({ role }) => {
  const normalizedRole =
    role.charAt(0).toUpperCase() + role.slice(1).toLowerCase();

  const colors = {
    Owner: "bg-purple-500/10 text-purple-300 border-purple-500/20",
    Ops: "bg-blue-500/10 text-blue-300 border-blue-500/20",
    Dispatcher: "bg-green-500/10 text-green-300 border-green-500/20",
  };

  return (
    <span
      className={`text-xs px-2 py-1 rounded border w-fit ${
        colors[normalizedRole] || "bg-slate-700 text-slate-300 border-slate-600"
      }`}
    >
      {normalizedRole}
    </span>
  );
};

const Stat = ({ label, value }) => {
  return (
    <div className="flex flex-col">
      <span className="text-xs text-slate-500">{label}</span>
      <span className="text-white font-medium">{value}</span>
    </div>
  );
};