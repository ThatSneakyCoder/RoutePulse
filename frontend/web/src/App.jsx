import Navbar from "./components/landingPage/Navbar.jsx";
import Hero from "./components/landingPage/Hero.jsx";
import Features from "./components/landingPage/Features.jsx";
import Pricing from "./components/landingPage/Pricing.jsx";
import Testimonials from "./components/landingPage/Testimonials.jsx";
import Footer from "./components/landingPage/Footer.jsx";

function App() {
  return (
    <div className="min-h-screen bg-slate-950 text-white overflow-hidden">
      <Navbar />
      <Hero />
      <Features />
      <Pricing />
      <Testimonials />
      <Footer />
    </div>
  );
}

export default App;
