import { useParams } from "react-router-dom";
import { Users, Plus, Trash2 } from "lucide-react";
import { useState } from "react";

export const OrganizationDetails = () => {
  const { orgId } = useParams();

  const [inviteEmail, setInviteEmail] = useState("");
  const [role, setRole] = useState("Dispatcher");

  const [members, setMembers] = useState([
    {
      id: 1,
      first_name: "John",
      last_name: "Doe",
      email: "owner@routepulse.com",
      role: "Owner",
      joined_at: "2025-01-10",
    },
    {
      id: 2,
      first_name: "Sarah",
      last_name: "Lee",
      email: "ops@routepulse.com",
      role: "Ops",
      joined_at: "2025-02-02",
    },
  ]);

  const inviteUser = () => {
    if (!inviteEmail.trim()) return;

    const newMember = {
      id: Date.now(),
      first_name: "New",
      last_name: "User",
      email: inviteEmail,
      role,
      joined_at: new Date().toISOString(),
    };

    setMembers((prev) => [...prev, newMember]);
    setInviteEmail("");
  };

  const removeUser = (id) => {
    setMembers((prev) => prev.filter((m) => m.id !== id));
  };

  return (
    <section className="p-8 space-y-8">
      {/* Header */}
      <div>
        <h1 className="text-xl font-semibold text-white">
          Organization #{orgId}
        </h1>
        <p className="text-sm text-slate-400">
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

        <div className="flex gap-3 flex-wrap">
          <input
            value={inviteEmail}
            onChange={(e) => setInviteEmail(e.target.value)}
            placeholder="user@email.com"
            className="bg-slate-800 border border-slate-700 rounded-md px-3 py-2 text-sm text-white w-64"
          />

          <select
            value={role}
            onChange={(e) => setRole(e.target.value)}
            className="bg-slate-800 border border-slate-700 rounded-md px-2 py-2 text-sm text-white"
          >
            <option>Dispatcher</option>
            <option>Ops</option>
            <option>Owner</option>
          </select>

          <button
            onClick={inviteUser}
            className="px-4 py-2 text-sm font-medium bg-blue-600 hover:bg-blue-500 text-white rounded-md transition"
          >
            Send Invite
          </button>
        </div>
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
                <span className="text-xs text-slate-400">{m.email}</span>
              </div>

              <RoleBadge role={m.role} />

              <div className="text-xs text-slate-400">
                {new Date(m.joined_at).toLocaleDateString()}
              </div>

              <button
                onClick={() => removeUser(m.id)}
                className="text-red-400 hover:text-red-300 transition"
              >
                <Trash2 size={16} />
              </button>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};

const RoleBadge = ({ role }) => {
  const colors = {
    Owner: "bg-purple-500/10 text-purple-300 border-purple-500/20",
    Ops: "bg-blue-500/10 text-blue-300 border-blue-500/20",
    Dispatcher: "bg-green-500/10 text-green-300 border-green-500/20",
  };

  return (
    <span className={`text-xs px-2 py-1 rounded border w-fit ${colors[role]}`}>
      {role}
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
