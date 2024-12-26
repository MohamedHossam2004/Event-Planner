import { useState, useContext } from "react";
import { X } from "lucide-react";
import { AuthContext } from "../contexts/AuthContext";
import { createEvent } from "../services/api";

export const CreateEventOverlay = ({ onClose }) => {
  const { showMessage } = useContext(AuthContext);
  const [eventData, setEventData] = useState({
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
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    if (name.includes(".")) {
      const [parent, child] = name.split(".");
      setEventData((prev) => ({
        ...prev,
        [parent]: {
          ...prev[parent],
          [child]: value,
        },
      }));
    } else if (name === "organizers") {
      setEventData((prev) => ({
        ...prev,
        organizers: value.split(",").map((name) => ({ name: name.trim() })),
      }));
    } else {
      setEventData((prev) => ({ ...prev, [name]: value }));
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const formattedEventTime = new Date(eventData.date).toISOString();

    const dataToSend = {
      ...eventData,
      date: formattedEventTime,
      min_capacity: Number.parseInt(eventData.min_capacity, 10),
      max_capacity: Number.parseInt(eventData.max_capacity, 10),
      organizers: eventData.organizers, // Include organizers in dataToSend
    };

    try {
      await createEvent(dataToSend);
      showMessage("Event created successfully!", "success");
      onClose();
    } catch (error) {
      showMessage("Failed to create event. Please try again.", "error");
    }
  };

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50 flex items-center justify-center">
      <div className="relative p-8 bg-white w-full max-w-4xl m-auto rounded-2xl shadow-lg">
        <button
          onClick={onClose}
          className="absolute top-4 right-4 p-2 rounded-full bg-white hover:bg-gray-100 transition-colors"
        >
          <X size={24} className="text-gray-600" />
        </button>
        <h2 className="text-3xl font-bold mb-6 text-purple-800">
          Create New Event
        </h2>
        <form
          onSubmit={handleSubmit}
          className="grid grid-cols-1 md:grid-cols-2 gap-6"
        >
          <div>
            <label
              htmlFor="name"
              className="block text-gray-700 font-bold mb-2"
            >
              Event Name
            </label>
            <input
              type="text"
              id="name"
              name="name"
              value={eventData.name}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-purple-500 focus:border-purple-500"
              required
            />
          </div>
          <div>
            <label
              htmlFor="type"
              className="block text-gray-700 font-bold mb-2"
            >
              Event Type
            </label>
            <select
              id="type"
              name="type"
              value={eventData.type}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-purple-500 focus:border-purple-500"
              required
            >
              <option value="Conference">Conference</option>
              <option value="Workshop">Workshop</option>
              <option value="Meetup">Meetup</option>
              <option value="Social">Social</option>
              <option value="Career Fair">Career Fair</option>
              <option value="Graduation">Graduation</option>
              <option value="Other">Other</option>
            </select>
          </div>
          <div>
            <label
              htmlFor="date"
              className="block text-gray-700 font-bold mb-2"
            >
              Date and Time
            </label>
            <input
              type="datetime-local"
              id="date"
              name="date"
              value={eventData.date}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-purple-500 focus:border-purple-500"
              required
            />
          </div>
          <div>
            <label
              htmlFor="min_capacity"
              className="block text-gray-700 font-bold mb-2"
            >
              Minimum Capacity
            </label>
            <input
              type="number"
              id="min_capacity"
              name="min_capacity"
              value={eventData.min_capacity}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-purple-500 focus:border-purple-500"
              required
              min="0"
            />
          </div>
          <div>
            <label
              htmlFor="max_capacity"
              className="block text-gray-700 font-bold mb-2"
            >
              Maximum Capacity
            </label>
            <input
              type="number"
              id="max_capacity"
              name="max_capacity"
              value={eventData.max_capacity}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-purple-500 focus:border-purple-500"
              required
              min={eventData.min_capacity}
            />
          </div>
          <div className="md:col-span-2">
            <label
              htmlFor="description"
              className="block text-gray-700 font-bold mb-2"
            >
              Description
            </label>
            <textarea
              id="description"
              name="description"
              value={eventData.description}
              onChange={handleChange}
              rows="4"
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-purple-500 focus:border-purple-500"
              required
            ></textarea>
          </div>
          <div>
            <label
              htmlFor="location.address"
              className="block text-gray-700 font-bold mb-2"
            >
              Address
            </label>
            <input
              type="text"
              id="location.address"
              name="location.address"
              value={eventData.location.address}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-purple-500 focus:border-purple-500"
              required
            />
          </div>
          <div>
            <label
              htmlFor="location.city"
              className="block text-gray-700 font-bold mb-2"
            >
              City
            </label>
            <input
              type="text"
              id="location.city"
              name="location.city"
              value={eventData.location.city}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-purple-500 focus:border-purple-500"
              required
            />
          </div>
          <div>
            <label
              htmlFor="location.state"
              className="block text-gray-700 font-bold mb-2"
            >
              State
            </label>
            <input
              type="text"
              id="location.state"
              name="location.state"
              value={eventData.location.state}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-purple-500 focus:border-purple-500"
              required
            />
          </div>
          <div>
            <label
              htmlFor="location.country"
              className="block text-gray-700 font-bold mb-2"
            >
              Country
            </label>
            <input
              type="text"
              id="location.country"
              name="location.country"
              value={eventData.location.country}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-purple-500 focus:border-purple-500"
              required
            />
          </div>
          <div>
            <label
              htmlFor="organizers"
              className="block text-gray-700 font-bold mb-2"
            >
              Organizers (comma-separated)
            </label>
            <input
              type="text"
              id="organizers"
              name="organizers"
              value={eventData.organizers.map((org) => org.name).join(", ")}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-purple-500 focus:border-purple-500"
              required
            />
          </div>
          <div className="md:col-span-2">
            <button
              type="submit"
              className="w-full flex justify-center py-3 px-4 border border-transparent rounded-lg shadow-sm text-lg font-medium text-white bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-700 hover:to-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500"
            >
              Create Event
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
