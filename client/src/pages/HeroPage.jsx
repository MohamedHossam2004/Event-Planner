import { useState, useEffect,useContext } from "react";
import { Stats } from "../components/Stats";
import { CategoryFilter } from "../components/CategoryFilter";
import { EventList } from "../components/EventList";
import { EventOverlay } from "../components/EventOverlay";
import { getEvents, getUnsubedEvents } from "../services/api";
import { subscribeMailingList} from "../services/api";
import { AuthContext } from "../contexts/AuthContext";

export const HeroPage = ({ events, setEvents }) => {
  const { user, setUser } = useContext(AuthContext);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [selectedCategory, setSelectedCategory] = useState("All");
  const [selectedEvent, setSelectedEvent] = useState(null);
  const [message, setMessage] = useState("");
  const [messageType, setMessageType] = useState(""); // "success" or "error"

  const handleCategorySelect = (category) => {
    setSelectedCategory(category);
  };

  const handleEventSelect = (event) => {
    setSelectedEvent(event);
  };

  const handleCloseOverlay = () => {
    setSelectedEvent(null);
  };

  const handleSubscribe = async (category) => {
    if (category === "All") {
      category = "general";
    }
    try {
      const response = await subscribeMailingList(category);
      setMessage(response.message);
      setMessageType("success"); // Set message type as success
    } catch (error) {
      setMessage(error.response.message)
      setMessageType("error"); // Set message type as error
    }
  };

  useEffect(() => {
    const fetchEvents = async () => {

      

      //console.log(user.isAdmin)
      if(user&&user.isAdmin){
       const data = await getEvents();
      if (data.events != null) {
        setEvents(data.events);
        setLoading(false);
      } else {
        setError("Failed to fetch events. Please try again later.");
        setLoading(false);
      }
    } 
      else{
        const data = await getUnsubedEvents();
        if (data.unsubscribed_events != null) {
          setEvents(data.unsubscribed_events);
          setLoading(false);
        } else {
          setError("Failed to fetch events. Please try again later.");
          setLoading(false);
        }


      }
      
    };

    fetchEvents();
  }, []);

  const categories = events.map(event => event.type);
  categories.unshift("All");
  const filteredEvents =
    selectedCategory === "All"
      ? events
      : events.filter((event) => event.type === selectedCategory);

  return (
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

      {/* Centering the Category Filter */}
      <div className="flex justify-center mt-6">
        <CategoryFilter
        categories={categories}
          selectedCategory={selectedCategory}
          onCategorySelect={handleCategorySelect}
        />
      </div>

      <div className="text-center mt-6">
        <button
          className="px-6 py-3 bg-purple-600 text-white rounded-full font-semibold hover:bg-purple-700"
          onClick={() => handleSubscribe(selectedCategory)}
        >
          Subscribe to {selectedCategory} Mailing List
        </button>
      </div>

      {/* Styled Message */}
      {message && (
        <div
          className={`mt-4 py-2 px-4 rounded-lg text-center ${
            messageType == "success"
              ? "bg-purple-100 text-green-800 border"
              : "bg-purple-100 text-red-800 border"
          } transition-all`}
        >
          {message}
        </div>
      )}

      {loading ? (
        <div className="text-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600 mx-auto"></div>
        </div>
      ) : error ? (
        <div className="text-center py-12 text-red-600">{error}</div>
      ) : (
        events && (
          <>
            <EventList
              events={filteredEvents}
              onEventSelect={handleEventSelect}
            />
          </>
        )
      )}

      {selectedEvent && (
        <EventOverlay event={selectedEvent} onClose={handleCloseOverlay} />
      )}
    </main>
  );
};
