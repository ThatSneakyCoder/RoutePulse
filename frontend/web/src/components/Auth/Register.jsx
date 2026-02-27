import Bookingdotcom from "../../assets/swiper-logos/Bookingdotcom.png";
import Googledotcom from "../../assets/swiper-logos/Googledotcom.png";
import Amazondotcom from "../../assets/swiper-logos/Amazondotcom.png";
import Netflixdotcom from "../../assets/swiper-logos/Netflixdotcom.png";

import { ArrowBigLeft, Key, Lock, LogIn, Mail } from "lucide-react";
import { FcGoogle } from "react-icons/fc";
import { FaMeta } from "react-icons/fa6";
import { FaApple } from "react-icons/fa";
import { Link } from "react-router";

// Import Swiper React components
import { Swiper, SwiperSlide } from "swiper/react";
import { Autoplay } from "swiper/modules";

// Import Swiper styles
import "swiper/css";

export const Register = () => {
  return (
    <div className="pt-16 sm:pt-18 px-4 sm:px-6 lg:px-8 h-full overflow-hidden flex items-center justify-center">
      {/* divs for decoration */}
      <div className="absolute bottom-5 left-4 sm:left-40 w-48 sm:w-72 h-48 sm:h-72 bg-blue-500/10 rounded-full blur-3xl animate-pulse" />
      <div className="absolute top-20 right-4 sm:right-10 w-64 sm:w-96 h-64 sm:h-96 bg-cyan-500/10 rounded-full blur-3xl animate-pulse delay-300" />
      <div className="absolute -top-30 -left-30 sm:right-10 w-64 sm:w-96 h-64 sm:h-96 bg-cyan-200/10 rounded-full blur-3xl animate-pulse delay-700" />

      <div className="flex max-w-6xl w-full gap-4 justify-center">
        {/* left side info */}
        <div className="hidden lg:flex z-10 w-full lg:max-w-2xl rounded-2xl p-12 flex-col items-start justify-between">
          <Link to="/" className="flex items-center gap-2">
            <ArrowBigLeft size={20} className="h-7 w-7 text-blue-400" />
            <span className="text-gray-400">Go back</span>
          </Link>
          <h2 className="text-5xl font-semibold text-gray-200">
            The simplest way to manage your fleet
          </h2>
          <p className="text-gray-400 leading-relaxed mb-2 text-2xl">
            Track deliveries, manage teams, and get real-time updates from a
            single intuitive dashboard. Boost productivity and streamline your
            entire logistics pipeline.
          </p>

          {/* swiper slider */}
          <div className="w-full max-w-2xl overflow-hidden">
            <SwiperSlider />
          </div>
        </div>

        {/* right side form */}
        <div className="relative z-10 w-full max-w-md rounded-2xl bg-slate-900/80 backdrop-blur-xl border border-white/10 p-8 shadow-2xl flex flex-col items-center justify-center gap-2">
          <div className="p-2 bg-gray-600 rounded-xl">
            <LogIn size={36} />
          </div>
          <div className="flex flex-col items-center justify-center text-center">
            <span className="text-2xl text-gray-300">
              Register with email
            </span>
            <span className="text-md text-gray-600 leading-relaxed">
              Make a new account to start tracking your delivery system
              efficiently
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
                className="w-full pl-10 pr-4 py-3 rounded-xl bg-slate-800 border border-white/10 hover:border-blue-500/40 focus:outline-none focus:ring-2 focus:ring-blue-500/40"
              />
            </div>
            <div>
              <span className="w-full text-gray-400 text-sm">
                Type the password again
              </span>
            </div>
            <div className="relative">
              <Key className="absolute w-4 h-4 text-gray-400 left-4 top-1/2 -translate-y-1/2" />
              <input
                type="password"
                placeholder="password"
                className="w-full pl-10 pr-4 py-3 bg-slate-800 border border-white/10 rounded-xl hover:border-blue-500/40 focus:outline-none focus:ring-2 focus:ring-blue-500/40"
              />
            </div>
            <Link
              to="/auth/forgotPassword"
              className="flex items-center justify-end"
            >
              <span className="text-sm text-gray-400">Forgot Password?</span>
            </Link>
            {/* TODO: make a backend call */}
            <Link
              to=""
              className="w-full py-3 flex items-center justify-center bg-gray-700/80 rounded-xl text-gray-300/80 hover:scale-102 duration-300"
            >
              <span>Register</span>
            </Link>
            {/* Oauth buttons */}
            <div className="flex items-center gap-4 w-full">
              <div className="flex-1 h-px bg-gray-600" />
              <span className="text-gray-400 text-sm">Or register with</span>
              <div className="flex-1 h-px bg-gray-600" />
            </div>
            {/* Oauth tabs */}
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
          <div>
            <span className="text-gray-400">Already have an account?{" "}<Link to="/auth/login" className="text-blue-400 hover:underline">Sign In</Link></span>
          </div>
        </div>
      </div>
    </div>
  );
};

const SwiperSlider = () => {
  return (
    <Swiper
      modules={[Autoplay]}
      spaceBetween={20}
      slidesPerView={3}
      loop={true}
      autoplay={{
        delay: 0,
        disableOnInteraction: false,
      }}
      speed={3000}
      allowTouchMove={false}
    >
      <SwiperSlide className="flex items-center justify-center p-4">
        <img
          src={Bookingdotcom}
          alt="booking.com"
          className="h-10 w-full object-contain opacity-80"
        />
      </SwiperSlide>
      <SwiperSlide className="flex items-center justify-center p-4">
        <img
          src={Googledotcom}
          alt="google.com"
          className="h-10 w-full object-contain opacity-80"
        />
      </SwiperSlide>
      <SwiperSlide className="flex items-center justify-center p-4">
        <img
          src={Amazondotcom}
          alt="amazon.com"
          className="h-10 w-full object-contain opacity-80"
        />
      </SwiperSlide>
      <SwiperSlide className="flex items-center justify-center p-4">
        <img
          src={Netflixdotcom}
          alt="netflix.com"
          className="h-10 w-full object-contain opacity-80"
        />
      </SwiperSlide>
    </Swiper>
  );
};
