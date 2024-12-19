import { useState, useEffect } from "react";
import { Calendar, Users, TrendingUp } from "lucide-react";
import { getEvents } from "../services/api";

export const Stats = () => {
  const [eventCount, setEventCount] = useState(0);

  useEffect(() => {
    const fetchEventCount = async () => {
      try {
        const events = await getEvents();
        setEventCount(events.events.length);
      } catch (error) {
        console.error("Failed to fetch event count", error);
      }
    };

    fetchEventCount();
  }, []);

  return (
    <div className="grid grid-cols-3 gap-8 max-w-4xl mx-auto my-8">
      <div className="bg-white p-6 rounded-xl shadow-sm flex items-center gap-4">
        <Calendar className="text-purple-500" size={24} />
        <div>
          <p className="text-2xl font-bold">{eventCount}+</p>
          <p className="text-gray-600">Active Events</p>
        </div>
      </div>

      <div className="bg-white p-6 rounded-xl shadow-sm flex items-center gap-4">
        <Users className="text-purple-500" size={24} />
        <div>
          <p className="text-2xl font-bold">15,000+</p>
          <p className="text-gray-600">Community Members</p>
        </div>
      </div>

      <div className="bg-white p-6 rounded-xl shadow-sm flex items-center gap-4">
        <TrendingUp className="text-purple-500" size={24} />
        <div>
          <p className="text-2xl font-bold">5,000+</p>
          <p className="text-gray-600">Monthly Attendees</p>
        </div>
      </div>
    </div>
  );
};
