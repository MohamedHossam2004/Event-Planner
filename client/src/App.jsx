import React, { useState } from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import { Header } from "./components/Header";
import { HeroPage } from "./pages/HeroPage";
import { Login } from "./pages/Login";
import { Signup } from "./pages/Signup";
import { CreateEventOverlay } from "./components/CreateEventOverlay";
import { AuthProvider } from "./contexts/AuthContext";
//import {RegisterdEvents} from './components/RegisterdEvents'
import MyEvents from "./components/myevents";
import EventApplications from "./components/eventApplications";

function App() {
  const [showCreateEventOverlay, setShowCreateEventOverlay] = useState(false);

  const handleCreateEvent = () => {
    setShowCreateEventOverlay(true);
  };

  const handleCloseCreateEventOverlay = () => {
    setShowCreateEventOverlay(false);
  };

  return (
    <AuthProvider>
      <Router>
        <div className="min-h-screen bg-gray-50">
          <Header onCreateEvent={handleCreateEvent} />

          <Routes>
            <Route path="/" element={<HeroPage />} />
            <Route path="/login" element={<Login />} />
            <Route path="/signup" element={<Signup />} />
             <Route path="/myevents" element={<MyEvents />} /> 
             <Route path="/eventApplications" element={<EventApplications />} />
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

          {showCreateEventOverlay && (
            <CreateEventOverlay onClose={handleCloseCreateEventOverlay} />
          )}
        </div>
      </Router>
    </AuthProvider>
  );
}

export default App;
