import { Mail, Calendar, ShieldCheck, User } from "lucide-react";
import { useRouteLoaderData } from "react-router-dom";
import { useState } from "react";
import axios from "../../../axios";

export const Profile = () => {
  const loaderUser = useRouteLoaderData("dashboard");
  const [{ user }, setUser] = useState(loaderUser);

  const [sendingVerify, setSendingVerify] = useState(false);
  const [resendingToken, setResendingToken] = useState(false);
  const [token, setToken] = useState("");

  const submitVerification = async () => {
    try {
      setSendingVerify(true);

      const res = await axios.post("/v1/authentication/verify-email", {
        email: user.email,
        verify_token: token.trim(),
      });

      setUser((prev) => ({
        ...prev,
        is_verified: res.data.is_verified,
      }));

      setToken("");
    } catch {
      alert("Invalid or expired verification token");
    } finally {
      setSendingVerify(false);
    }
  };

  const resendVerification = async () => {
    try {
      setResendingToken(true);

      await axios.post("/v1/authentication/resend-email-verification", {
        email: user.email,
      });

      alert("Verification email sent.");
    } finally {
      setResendingToken(false);
    }
  };

  return (
    <section className="h-full w-full p-8">
      {/* Page Header */}
      <div className="mb-10">
        <h1 className="text-xl font-semibold text-slate-100">
          User Profile
        </h1>
        <p className="text-sm text-slate-400">
          Manage your account information
        </p>
      </div>

      {/* Profile Card */}
      <div className="bg-slate-900 border border-slate-800 rounded-xl p-8">

        {/* Avatar + Name */}
        <div className="flex items-center gap-6 pb-8 border-b border-slate-800">

          <div className="w-16 h-16 flex items-center justify-center
          rounded-full bg-slate-800 border border-slate-700">

            <User className="w-7 h-7 text-slate-300" />

          </div>

          <div>
            <h2 className="text-lg font-semibold text-white">
              {user.first_name} {user.last_name}
            </h2>

            <p className="text-sm text-slate-400">
              {user.email}
            </p>
          </div>

        </div>

        {/* Account Info */}
        <div className="mt-8 grid grid-cols-1 md:grid-cols-2 gap-8">

          <ProfileField
            icon={<Mail size={16} />}
            label="Email"
            value={user.email}
          />

          <ProfileField
            icon={<Calendar size={16} />}
            label="Account Created"
            value={new Date(user.created_at).toLocaleDateString()}
          />

          <ProfileField
            icon={<ShieldCheck size={16} />}
            label="Account Status"
            value={
              <StatusBadge active={user.is_active}>
                {user.is_active ? "Active" : "Inactive"}
              </StatusBadge>
            }
          />

          <ProfileField
            label="User ID"
            value={
              <span className="text-xs text-slate-400 break-all">
                {user.id}
              </span>
            }
          />

        </div>

        {/* Email Verification Section */}
        <div className="mt-10 pt-8 border-t border-slate-800">

          <div className="flex items-center gap-2 mb-3 text-sm text-slate-400">
            Email Verification
          </div>

          <div className="flex items-center gap-4">

            <StatusBadge verified={user.is_verified}>
              {user.is_verified ? "Verified" : "Not Verified"}
            </StatusBadge>

            {!user.is_verified && (
              <button
                onClick={resendVerification}
                disabled={resendingToken}
                className="text-xs text-blue-400 hover:text-blue-300"
              >
                {resendingToken
                  ? "Sending..."
                  : "Resend verification email"}
              </button>
            )}

          </div>

          {!user.is_verified && (
            <div className="flex items-center gap-3 mt-5">

              <input
                type="text"
                maxLength={6}
                value={token}
                onChange={(e) =>
                  setToken(e.target.value.replace(/\D/g, ""))
                }
                placeholder="Enter 6-digit code"
                className="bg-slate-800 border border-slate-700
                rounded-md px-3 py-2 text-sm text-white w-40"
              />

              <button
                onClick={submitVerification}
                disabled={sendingVerify || token.length !== 6}
                className="px-3 py-2 text-xs font-medium
                bg-green-600 hover:bg-green-500
                text-white rounded-md disabled:opacity-50"
              >
                {sendingVerify ? "Verifying..." : "Verify"}
              </button>

            </div>
          )}

        </div>

      </div>
    </section>
  );
};

const ProfileField = ({ icon, label, value }) => (
  <div>
    <div className="flex items-center gap-2 text-xs uppercase tracking-wide text-slate-500 mb-1">
      {icon}
      {label}
    </div>

    <div className="text-sm text-slate-200 font-medium">
      {value}
    </div>
  </div>
);

const StatusBadge = ({ children, active, verified }) => {
  const color =
    active !== undefined
      ? active
        ? "bg-green-500/10 text-green-300 border-green-500/20"
        : "bg-red-500/10 text-red-300 border-red-500/20"
      : verified
      ? "bg-green-500/10 text-green-300 border-green-500/20"
      : "bg-yellow-500/10 text-yellow-300 border-yellow-500/20";

  return (
    <span
      className={`text-xs px-2 py-1 rounded-md border ${color}`}
    >
      {children}
    </span>
  );
};