import { ArrowBigLeft, Key, KeySquare, Lock, LogIn, Mail } from "lucide-react";
import { FcGoogle } from "react-icons/fc";
import { FaMeta } from "react-icons/fa6";
import { FaApple } from "react-icons/fa";
import { Link } from "react-router";

export const Login = () => {
  return (
    <div className="pt-16 sm:pt-18 px-4 sm:px-6 lg:px-8 h-full overflow-hidden flex items-center justify-center">
      {/* background blobs */}
      <div className="absolute bottom-5 left-4 sm:left-40 w-48 sm:w-72 h-48 sm:h-72 bg-blue-500/10 rounded-full blur-3xl animate-pulse" />
      <div className="absolute top-20 right-4 sm:right-10 w-64 sm:w-96 h-64 sm:h-96 bg-cyan-500/10 rounded-full blur-3xl animate-pulse delay-300" />
      <div className="absolute -top-30 -left-30 sm:right-10 w-64 sm:w-96 h-64 sm:h-96 bg-cyan-200/10 rounded-full blur-3xl animate-pulse delay-700" />

      <div className="flex max-w-6xl w-full gap-4 justify-center">
        {/* left side content */}
        <div className="hidden lg:flex z-10 w-full max-w-2xl rounded-2xl p-12 flex-col items-start justify-between">
          {/* Back Button */}
          <Link to="/" className="flex items-center gap-2">
            <ArrowBigLeft size={20} className="h-7 w-7 text-blue-400" />
            <span className="text-gray-400">Go back</span>
          </Link>

          {/* Headline */}
          <h2 className="text-5xl font-semibold text-gray-200 leading-tight my-8">
            Welcome back to your workspace
          </h2>

          {/* Subtext */}
          <p className="text-gray-400 mb-8 text-2xl leading-relaxed">
            Log in to continue managing deliveries, viewing analytics, and
            monitoring your fleet's performance in real time.
          </p>

          {/* Highlight Box */}
          <div className="mt-4 p-6 rounded-2xl bg-slate-900/40 border border-white/5 backdrop-blur-xl">
            <h3 className="text-xl font-medium text-gray-200 mb-3">
              What you can do:
            </h3>
            <ul className="space-y-3 text-gray-400">
              <li className="flex items-start gap-2">
                <span className="w-2 h-2 mt-2 rounded-full bg-blue-500"></span>
                Live-track active deliveries with precise updates.
              </li>
              <li className="flex items-start gap-2">
                <span className="w-2 h-2 mt-2 rounded-full bg-blue-500"></span>
                Manage teams and assign delivery tasks effortlessly.
              </li>
              <li className="flex items-start gap-2">
                <span className="w-2 h-2 mt-2 rounded-full bg-blue-500"></span>
                Explore insights and detect bottlenecks instantly.
              </li>
            </ul>
          </div>
        </div>

        {/* right side content */}
        <div className="relative z-10 w-full max-w-md rounded-2xl bg-slate-900/80 backdrop-blur-xl border border-white/10 p-8 shadow-2xl flex flex-col items-center justify-center gap-5">
          <div className="p-2 bg-gray-600 rounded-xl">
            <KeySquare size={36} />
          </div>

          <div className="flex flex-col items-center justify-center text-center">
            <span className="text-2xl text-gray-300">Sign in with email</span>
            <span className="text-md text-gray-600 leading-relaxed">
              Sign in to access your dashboard.
            </span>
          </div>

          <div className="max-w-md w-full space-y-4">
            <div className="relative">
              <Mail className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
              <input
                type="email"
                placeholder="Email"
                className="w-full pl-10 pr-4 py-3 rounded-lg bg-slate-800 border border-white/10 hover:border-blue-500/40 
                focus:outline-none focus:ring-2 focus:ring-blue-500/40"
              />
            </div>

            <div className="relative">
              <Lock className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
              <input
                type="password"
                placeholder="Password"
                className="w-full pl-10 pr-4 py-3 rounded-xl bg-slate-800 border border-white/10 hover:border-blue-500/40 
                focus:outline-none focus:ring-2 focus:ring-blue-500/40"
              />
            </div>

            <Link
              to="/auth/forgotPassword"
              className="flex items-center justify-end"
            >
              <span className="text-sm text-gray-400">Forgot Password?</span>
            </Link>

            {/* signin button */}
            <Link
              to=""
              className="w-full py-3 flex items-center justify-center bg-gray-700/80 rounded-xl text-gray-300/80 hover:scale-102 duration-300"
            >
              <span>Sign In</span>
            </Link>

            {/* OAuth section */}
            <div className="flex items-center gap-4 w-full">
              <div className="flex-1 h-px bg-gray-600" />
              <span className="text-gray-400 text-sm">Or sign in with</span>
              <div className="flex-1 h-px bg-gray-600" />
            </div>

            <div className="flex items-center justify-center gap-6 w-full">
              <button className="p-3 rounded-xl bg-slate-800 border border-white/10 hover:bg-slate-700 hover:scale-110 duration-300">
                <FcGoogle size={24} />
              </button>
              <button className="p-3 rounded-xl bg-slate-800 border border-white/10 hover:bg-slate-700 hover:scale-110 duration-300">
                <FaMeta size={24} />
              </button>
              <button className="p-3 rounded-xl bg-slate-800 border border-white/10 hover:bg-slate-700 hover:scale-110 duration-300">
                <FaApple size={24} />
              </button>
            </div>
          </div>

          {/* bottom redirect */}
          <div>
            <span className="text-gray-400">
              Don't have an account?{" "}
              <Link
                to="/auth/register"
                className="text-blue-400 hover:underline"
              >
                Register
              </Link>
            </span>
          </div>
        </div>
      </div>
    </div>
  );
};
