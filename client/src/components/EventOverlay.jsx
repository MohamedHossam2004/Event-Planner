import React from "react";
import { X, Calendar, Clock, MapPin, Users } from "lucide-react";
import { formatDate, formatTime } from "../services/helpers";

export const EventOverlay = ({ event, onClose }) => {
  if (!event) return null;

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div className="relative top-20 mx-auto p-5 border w-full max-w-4xl shadow-lg rounded-md bg-white">
        <div className="relative">
          <img
            src={event.imageUrl}
            alt={event.name}
            className="w-full h-64 object-cover rounded-t-md"
          />
          <button
            onClick={onClose}
            className="absolute top-4 right-4 p-2 rounded-full bg-white hover:bg-gray-100 transition-colors"
            aria-label="Close"
          >
            <X className="h-6 w-6 text-gray-600" />
          </button>
        </div>

        <div className="mt-4">
          <h2 className="text-2xl font-bold text-gray-800">{event.name}</h2>
          <p className="text-gray-600 mt-1">{event.organizers[0].name}</p>
          <div className="mt-2 flex flex-wrap gap-4 text-sm text-gray-600">
            <div className="flex items-center gap-1">
              <Calendar className="h-5 w-5" />
              <span>{formatDate(event.date)}</span>
            </div>
            <div className="flex items-center gap-1">
              <Clock className="h-5 w-5" />
              <span>{formatTime(event.date)}</span>
            </div>
            <div className="flex items-center gap-1">
              <MapPin className="h-5 w-5" />
              <span>{`${event.location.address}, ${event.location.city}, ${event.location.state}, ${event.location.country}`}</span>
            </div>
            <div className="flex items-center gap-1">
              <Users className="h-5 w-5" />
              <span>{event.number_of_applications} spots left</span>
            </div>
          </div>
        </div>

        <div className="mt-4">
          <h3 className="font-semibold text-gray-800">About the event</h3>
          <p className="mt-2 text-sm text-gray-600">{event.description}</p>
        </div>

        <div className="mt-6 flex items-center justify-between">
          <span className="px-2 py-1 rounded-full text-xs font-semibold bg-purple-100 text-purple-800">
            {event.type}
          </span>
          <div className="space-x-2">
            <button className="px-4 py-2 border border-purple-500 text-purple-500 rounded hover:bg-purple-100 transition-colors">
              Share Event
            </button>
            <button className="px-4 py-2 bg-purple-500 text-white rounded hover:bg-purple-600 transition-colors">
              Register Now
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};
