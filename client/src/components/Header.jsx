import React, { useContext } from "react";
import { Link } from "react-router-dom";
import { Search, CalendarPlus } from "lucide-react";
import { AuthContext } from "../contexts/AuthContext";
import { logout } from "../services/api";

export const Header = ({ onCreateEvent }) => {
  const { user, setUser } = useContext(AuthContext);

  const handleLogout = async () => {
    await logout();
    setUser(null);
    window.location.href = "/";
  };

  return (
    <header className="bg-gradient-to-r from-purple-600 to-blue-600 p-4">
      <div className="max-w-7xl mx-auto flex items-center justify-between">
        <Link to="/" className="text-white text-2xl font-bold">
          Event Hub
        </Link>

        <div className="flex-1 max-w-2xl mx-8">
          <div className="relative">
            <input
              type="text"
              placeholder="Search Events..."
              className="w-full px-4 py-2 rounded-lg bg-white/90 focus:outline-none focus:ring-2 focus:ring-purple-300"
            />
            <Search
              className="absolute right-3 top-2.5 text-gray-500"
              size={20}
            />
          </div>
        </div>

        <nav className="flex items-center space-x-4">
          {user ? (
            <>
              <span className="text-white">Welcome, {user.name}</span>
              <button
                onClick={handleLogout}
                className="text-white hover:text-purple-200"
              >
                Logout
              </button>
              {user.isAdmin && (
                <button
                  onClick={onCreateEvent}
                  className="bg-white text-purple-600 px-4 py-2 rounded-lg flex items-center gap-2 hover:bg-purple-50 transition-colors"
                >
                  <CalendarPlus size={20} />
                  Create Event
                </button>
              )}
            </>
          ) : (
            <>
              <Link to="/login" className="text-white hover:text-purple-200">
                Login
              </Link>
              <Link to="/signup" className="text-white hover:text-purple-200">
                Sign Up
              </Link>
            </>
          )}
        </nav>
      </div>
    </header>
  );
};
