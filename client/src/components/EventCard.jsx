import { MapPin, Clock, Users, Calendar } from "lucide-react";
import { formatDate, formatTime } from "../services/helpers";

export const EventCard = ({ event, onSelect }) => {
  return (
    <div
      className="bg-white rounded-xl overflow-hidden shadow-sm hover:shadow-md transition-shadow cursor-pointer event-card"
      onClick={() => onSelect(event)}
    >
      <img
        src={event.imageUrl}
        alt={event.name}
        className="w-full h-48 object-cover"
      />
      <div className="p-6">
        <span className="px-3 py-1 text-sm bg-purple-100 text-purple-700 rounded-full">
          {event.type}
        </span>
        <h3 className="text-xl font-semibold mt-3">{event.name}</h3>
        <p className="text-gray-600 mt-1">{event.organizers[0].name}</p>

        <div className="mt-4 space-y-2">
          <div className="flex items-center gap-2 text-gray-600">
            <Calendar size={18} />
            <span>{formatDate(event.date)}</span>
          </div>
          <div className="flex items-center gap-2 text-gray-600">
            <Clock size={18} />
            <span>{formatTime(event.date)}</span>
          </div>
          <div className="flex items-center gap-2 text-gray-600">
            <MapPin size={18} />
            <span>{`${event.location.address}, ${event.location.city}`}</span>
          </div>
        </div>

        <div className="mt-6 flex items-center justify-between">
          <div className="flex items-center gap-1 text-gray-600">
            <Users size={18} />
            <span>
              {event.max_capacity - event.number_of_applications} spots left
            </span>
          </div>
          <button
            className="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors"
            onClick={(e) => {
              e.stopPropagation();
              onSelect(event);
            }}
          >
            View Details
          </button>
        </div>
      </div>
    </div>
  );
};
