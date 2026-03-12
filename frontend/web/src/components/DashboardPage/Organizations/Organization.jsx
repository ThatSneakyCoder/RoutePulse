import { Building2, ChevronDown, ChevronUp, Plus } from "lucide-react";
import { useState } from "react";
import {
  Form,
  useNavigate,
  useRouteLoaderData,
} from "react-router-dom";

export const Organization = () => {
  const navigate = useNavigate();
  const { organizations } = useRouteLoaderData("dashboard");
  const [createOpen, setCreateOpen] = useState(false);
  const [orgName, setOrgName] = useState("");
  const [error, setError] = useState("");

  return (
    <section className="h-full w-full p-8 space-y-8">
      {/* Header */}
      <div>
        <h1 className="text-xl font-semibold text-white">Organization</h1>
        <p className="text-sm text-slate-400">
          Manage your organizations and teams
        </p>
      </div>

      {/* Create Organization */}
      <Form
        method="POST"
        className="bg-slate-900 border border-slate-800 rounded-xl p-6"
      >
        <button
          type="button"
          onClick={() => setCreateOpen(!createOpen)}
          className="flex items-center justify-between w-full text-white font-medium"
        >
          <span className="flex items-center gap-2">
            <Plus size={16} />
            Create Organization
          </span>

          {createOpen ? <ChevronUp size={16} /> : <ChevronDown size={16} />}
        </button>

        {createOpen && (
          <div className="flex gap-3 mt-4">
            <input
              name="name"
              value={orgName}
              onChange={(e) => setOrgName(e.target.value)}
              placeholder="Organization name"
              className="bg-slate-800 border border-slate-700 rounded-md px-3 py-2 text-sm text-white w-72 focus:outline-none focus:ring-0"
            />

            <button
              type="submit"
              name="intent"
              value="create-org"
              className="px-4 py-2 text-sm font-medium bg-blue-600 hover:bg-blue-500 text-white rounded-md transition"
            >
              Create
            </button>
          </div>
        )}
      </Form>

      {error && <p className="text-red-400 text-sm">{error}</p>}

      {/* Organizations */}
      <div className="grid md:grid-cols-2 gap-6">
        {organizations.map((org) => (
          <div
            key={org.organization_id}
            onClick={() =>
              navigate(`/dashboard/organization/${org.organization_id}`)
            }
            className="cursor-pointer bg-slate-900 border border-slate-800 rounded-xl p-6 hover:border-slate-700 transition hover:shadow-lg hover:shadow-black/20"
          >
            <div className="flex items-center justify-between mb-5">
              <div className="flex items-center gap-4">
                <div className="w-12 h-12 flex items-center justify-center rounded-lg bg-linear-to-br from-blue-500/20 to-purple-500/20 border border-slate-700">
                  <Building2 className="w-6 h-6 text-slate-200" />
                </div>

                <div>
                  <div className="text-white font-semibold">{org.name}</div>

                  <div className="text-xs text-slate-400">
                    Created {new Date(org.created_at).toLocaleDateString()}
                  </div>
                </div>
              </div>

              <span className="text-xs px-2 py-1 rounded-md bg-blue-500/10 text-blue-300 border border-blue-500/20">
                Org Plan
              </span>
            </div>

            <div className="border-t border-slate-800 mb-5"></div>

            <div className="flex items-center gap-10 text-sm">
              <Stat label="Members" value="12" />
              <Stat label="Vehicles" value="4" />
              <Stat label="Drivers" value="4" />
            </div>
          </div>
        ))}
      </div>
    </section>
  );
};

const Stat = ({ label, value }) => (
  <div className="flex flex-col">
    <span className="text-xs text-slate-500">{label}</span>
    <span className="text-white font-medium">{value}</span>
  </div>
);
