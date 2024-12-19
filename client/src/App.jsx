import { useState, useEffect } from "react";
import { Header } from "./components/Header";
import { Stats } from "./components/Stats";
import { CategoryFilter } from "./components/CategoryFilter";
import { EventList } from "./components/EventList";
import { EventOverlay } from "./components/EventOverlay";
import { getEvents } from "./services/api";

function App() {
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [selectedCategory, setSelectedCategory] = useState("All");
  const [selectedEvent, setSelectedEvent] = useState(null);

  const handleCategorySelect = (category) => {
    setSelectedCategory(category);
  };

  const handleEventSelect = (event) => {
    setSelectedEvent(event);
  };

  const handleCloseOverlay = () => {
    setSelectedEvent(null);
  };

  useEffect(() => {
    const fetchEvents = async () => {
      try {
        const data = await getEvents();
        setEvents(data);
        setLoading(false);
      } catch {
        setError("Failed to fetch events. Please try again later.");
        setLoading(false);
      }
    };

    fetchEvents();
  }, []);

  const filteredEvents =
    selectedCategory === "All"
      ? events
      : events.filter((event) => event.type === selectedCategory);

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />

      <main className="py-8">
        <div className="max-w-7xl mx-auto px-4">
          <h1 className="text-4xl font-bold text-center text-purple-800">
            Discover Amazing Events
          </h1>
          <p className="text-center text-gray-600 mt-2">
            Join exciting events and connect with like-minded people
          </p>
        </div>

        <Stats />
        <CategoryFilter
          selectedCategory={selectedCategory}
          onCategorySelect={handleCategorySelect}
        />

        {loading ? (
          <div className="text-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600 mx-auto"></div>
          </div>
        ) : error ? (
          <div className="text-center py-12 text-red-600">{error}</div>
        ) : (
          <EventList
            events={filteredEvents}
            onEventSelect={handleEventSelect}
          />
        )}
      </main>

      <footer className="bg-white border-t mt-16 py-8">
        <div className="max-w-7xl mx-auto px-4 text-center text-gray-600">
          <p>Event Hub</p>
          <p className="mt-2">
            Copyright {new Date().getFullYear()} Event Hub. All rights reserved.
          </p>
        </div>
      </footer>

      {selectedEvent && (
        <EventOverlay event={selectedEvent} onClose={handleCloseOverlay} />
      )}
    </div>
  );
}

export default App;
