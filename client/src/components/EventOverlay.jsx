import { useState } from "react";
import { X, Calendar, Clock, MapPin, Users } from "lucide-react";
import { formatDate, formatTime } from "../services/helpers";
import { applyToEvent, updateEvent, deleteEvent, getEventById } from "../services/api";
import { useContext } from "react";
import { AuthContext } from "../contexts/AuthContext";

export const EventOverlay = ({ event, onClose }) => {
  const { user } = useContext(AuthContext);
  const { showMessage } = useContext(AuthContext);
  const [isUpdateOverlayOpen, setUpdateOverlayOpen] = useState(false);
  const [existingEventData, setExistingEventData] = useState({
    name: "",
    type: "Conference",
    date: "",
    description: "",
    location: {
      address: "",
      city: "",
      state: "",
      country: "",
    },
    organizers: [{ name: "" }],
    min_capacity: 0,
    max_capacity: 0,
    status: "",
    ushers: [],
  });
  const [updatedEventData, setUpdatedEventData] = useState({
    name: "",
    type: "Conference",
    date: "",
    location: {
      address: "",
      city: "",
      state: "",
      country: "",
    },
    min_capacity: 0,
    max_capacity: 0,
  });

  const handleUpdateChange = (e) => {
    const { name, value } = e.target;
    if (name.includes(".")) {
      const [parent, child] = name.split(".");
      setUpdatedEventData((prev) => ({
        ...prev,
        [parent]: {
          ...prev[parent],
          [child]: value,
        },
      }));
    } else {
      setUpdatedEventData((prev) => ({ ...prev, [name]: value }));
    }
  };
  
  const onUpdateClick = async (eventId) => {
      setUpdateOverlayOpen(true);
      const existingEvent = await handleGetEventById(eventId);
      setExistingEventData(existingEvent);
  };

  const handleGetEventById = async (eventId) => {
    try {
      const result = await getEventById(eventId);
      return result;
    } catch (error) {
      throw new Error(error.response?.data?.message || "Failed to fetch event");
    }
  };
  
  const handleUpdateSubmit = async (eventId) => {
    try {
      console.log(existingEventData);
      const {
        number_of_applications,
        description,
        organizers,
        status,
        ushers,
      } = existingEventData.event;

      const formattedDate = new Date(updatedEventData.date).toISOString();
  
      setUpdatedEventData({
        ...updatedEventData,
        date: formattedDate,
        min_capacity: Number.parseInt(updatedEventData.min_capacity, 10),
        max_capacity: Number.parseInt(updatedEventData.max_capacity, 10),
        number_of_applications,
        description,
        organizers,
        status,
        ushers,
      });
  
      console.log(updatedEventData);
      const result = await updateEvent(eventId, updatedEventData);
      setUpdateOverlayOpen(false);
      showMessage("Event updated successfully!", "success");
      onClose();
    } catch (error) {
      showMessage(error.message || "Failed to update event", "error");
    }
  };
  
  const onDeleteClick = async (eventId) => {
      const confirmDelete = window.confirm("Are you sure you want to delete this event?");
      if (confirmDelete) {
          try {
            const result = await deleteEvent(eventId);
            showMessage("Event deleted successfully!", "success");
            onClose();
          } catch (error) {
            showMessage(error.message || "Failed to delete event", "error");
          }
      }
  };
  
  const onRegisterClick = async (eventId) => {
    try {
      const result = await applyToEvent(eventId);
      showMessage("Successfully registered for the event!", "success");
      onClose();
    } catch (error) {
      showMessage(error.message || "Failed to register for event", "error");
    }
  };

  if (!event) return null;

  const totalSpots = event.max_capacity;
  const availableSpots = totalSpots - event.number_of_applications;
  const progressPercentage = (availableSpots / totalSpots) * 100;

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50 flex items-center justify-center">
            {isUpdateOverlayOpen && (
        <div className="overlay">
          <div className="overlay-content p-6 bg-white rounded-lg shadow-lg max-w-md mx-auto">
            <h2 className="text-xl font-semibold mb-4">Update Event</h2>
            <form onSubmit={(e) => {
              e.preventDefault();
              handleUpdateSubmit(event._id);
            }}>
              <label className="block mb-2">
                Date:
                <input
                  type="datetime-local"
                  name="date"
                  value={updatedEventData.date}
                  onChange={handleUpdateChange}
                  className="mt-1 block w-full border rounded-md p-2"
                />
              </label>
              <label className="block mb-2">
                Type:
                <input
                  type="text"
                  name="type"
                  value={updatedEventData.type}
                  onChange={handleUpdateChange}
                  className="mt-1 block w-full border rounded-md p-2"
                />
              </label>
              <label className="block mb-2">
                Name:
                <input
                  type="text"
                  name="name"
                  value={updatedEventData.name}
                  onChange={handleUpdateChange}
                  className="mt-1 block w-full border rounded-md p-2"
                />
              </label>
              <label className="block mb-2">
                Location:
                <input
                  type="text"
                  name="location.address"
                  value={updatedEventData.location.address}
                  onChange={handleUpdateChange}
                  placeholder="Address"
                  className="mt-1 block w-full border rounded-md p-2"
                />
                <input
                  type="text"
                  name="location.city"
                  value={updatedEventData.location.city}
                  onChange={handleUpdateChange}
                  placeholder="City"
                  className="mt-1 block w-full border rounded-md p-2"
                />
                <input
                  type="text"
                  name="location.state"
                  value={updatedEventData.location.state}
                  onChange={handleUpdateChange}
                  placeholder="State"
                  className="mt-1 block w-full border rounded-md p-2"
                />
                <input
                  type="text"
                  name="location.country"
                  value={updatedEventData.location.country}
                  onChange={handleUpdateChange}
                  placeholder="Country"
                  className="mt-1 block w-full border rounded-md p-2"
                />
              </label>
              <label className="block mb-2">
                Min Capacity:
                <input
                  type="number"
                  name="min_capacity"
                  value={updatedEventData.min_capacity}
                  onChange={handleUpdateChange}
                  className="mt-1 block w-full border rounded-md p-2"
                />
              </label>
              <label className="block mb-2">
                Max Capacity:
                <input
                  type="number"
                  name="max_capacity"
                  value={updatedEventData.max_capacity}
                  onChange={handleUpdateChange}
                  className="mt-1 block w-full border rounded-md p-2"
                />
              </label>
              <div className="flex justify-between mt-4">
                <button
                  type="submit"
                  className="bg-blue-600 text-white rounded-md px-4 py-2 hover:bg-blue-700"
                >
                  Submit
                </button>
                <button
                  type="button"
                  className="bg-gray-300 text-black rounded-md px-4 py-2 hover:bg-gray-400"
                  onClick={() => setUpdateOverlayOpen(false)}
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
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
              { user.isAdmin ? (
                  <div className="space-y-3">
                      <button
                        className="w-full py-3 px-4 bg-gradient-to-r from-purple-600 to-blue-600 text-white rounded-lg hover:opacity-90 transition-opacity font-medium"
                        onClick={() => onUpdateClick(event._id)}
                      >
                        Update Event
                      </button>
                      <button
                        className="w-full py-3 px-4 bg-gradient-to-r from-purple-600 to-blue-600 text-white rounded-lg hover:opacity-90 transition-opacity font-medium"
                        onClick={() => onDeleteClick(event._id)}
                      >
                        Delete Event
                      </button>
                  </div>
              ) : (
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
              ) }
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
