import { useState } from "react";
import { X, Calendar, Clock, MapPin, Users } from "lucide-react";
import { formatDate, formatTime } from "../services/helpers";
import { applyToEvent } from "../services/api";

export const EventOverlay = ({ event, onClose }) => {
  const [errorMessage, setErrorMessage] = useState("");

  const onRegisterClick = async (eventId) => {
    try {
      setErrorMessage("");
      const result = await applyToEvent(eventId);
    } catch (error) {
      setErrorMessage(error.message);
    }
  };

  if (!event) return null;

  const totalSpots = 500; // Example total capacity
  const availableSpots = event.number_of_applications;
  const progressPercentage = (availableSpots / totalSpots) * 100;

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50 flex items-center justify-center">
      <div className="relative p-6 border w-full max-w-4xl shadow-lg rounded-2xl bg-white">
        {/* Close Button */}
        <button
          onClick={onClose}
          className="absolute top-4 right-4 p-2 rounded-full bg-white hover:bg-gray-100 transition-colors z-10 close-button"
          aria-label="Close"
        >
          <X className="h-6 w-6 text-gray-600" />
        </button>

        <div className="flex flex-col">
          {errorMessage && (
            <div className="text-red-600 text-center mb-4">{errorMessage}</div>
          )}
          {/* Title */}
          <h2 className="text-3xl font-bold text-gray-800 mb-4">
            {event.name}
          </h2>

          {/* Image */}
          <div className="relative w-full h-64 mb-6">
            <img
              src={event.imageUrl}
              alt={event.name}
              className="w-full h-full object-cover rounded-lg"
            />
          </div>

          {/* Content below image */}
          <div className="flex flex-col md:flex-row gap-6">
            {/* Main Content */}
            <div className="flex-1">
              {/* About Section */}
              <div className="mb-6">
                <h3 className="text-xl font-semibold mb-3">About the event</h3>
                <p className="text-gray-600">{event.description}</p>
              </div>

              {/* Divider */}
              <hr className="border-gray-200 my-6" />

              {/* Event Details Grid */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="flex items-center gap-2 text-gray-600">
                  <Calendar className="h-5 w-5 text-purple-500" />
                  <span>{formatDate(event.date)}</span>
                </div>
                <div className="flex items-center gap-2 text-gray-600">
                  <Clock className="h-5 w-5 text-purple-500" />
                  <span>{formatTime(event.date)}</span>
                </div>
                <div className="flex items-center gap-2 text-gray-600">
                  <MapPin className="h-5 w-5 text-purple-500" />
                  <span>{`${event.location.address}, ${event.location.city}, ${event.location.state}, ${event.location.country}`}</span>
                </div>
                <div className="flex items-center gap-2 text-gray-600">
                  <Users className="h-5 w-5 text-purple-500" />
                  <span>Organized by {event.organizers[0].name}</span>
                </div>
              </div>
            </div>

            {/* Right Side Panel */}
            <div className="md:w-72 space-y-6">
              {/* Available Spots */}
              <div className="bg-gray-50 p-4 rounded-lg">
                <h4 className="font-semibold mb-2">Available Spots</h4>
                <div className="w-full h-2 bg-gray-200 rounded-full overflow-hidden">
                  <div
                    className="h-full bg-gradient-to-r from-purple-600 to-blue-600"
                    style={{ width: `${progressPercentage}%` }}
                  />
                </div>
                <p className="text-sm text-gray-600 mt-2">
                  {availableSpots} spots remaining
                </p>
              </div>

              {/* Action Buttons */}
              <div className="space-y-3">
                <button
                  className="w-full py-3 px-4 bg-gradient-to-r from-purple-600 to-blue-600 text-white rounded-lg hover:opacity-90 transition-opacity font-medium event-register"
                  onClick={() => onRegisterClick(event._id)}
                >
                  Register Now
                </button>
                <button className="w-full py-3 px-4 bg-white text-purple-600 border border-purple-600 rounded-lg hover:bg-purple-50 transition-colors font-medium">
                  Apply as Usher
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
