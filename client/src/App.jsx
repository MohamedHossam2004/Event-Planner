import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import { Header } from "./components/Header";
import { HeroPage } from "./pages/HeroPage";
import { Login } from "./pages/Login";
import { Signup } from "./pages/Signup";

function App() {
  return (
    <Router>
      <div className="min-h-screen bg-gray-50">
        <Header />

        <Routes>
          <Route path="/" element={<HeroPage />} />
          <Route path="/login" element={<Login />} />
          <Route path="/signup" element={<Signup />} />
        </Routes>

        <footer className="bg-white border-t mt-16 py-8">
          <div className="max-w-7xl mx-auto px-4 text-center text-gray-600">
            <p>Event Hub</p>
            <p className="mt-2">
              Copyright {new Date().getFullYear()} Event Hub. All rights
              reserved.
            </p>
          </div>
        </footer>
      </div>
    </Router>
  );
}

export default App;
