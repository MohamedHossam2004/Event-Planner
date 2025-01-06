import React, { useState, useEffect } from "react";
import { getEventsForUser, unsubFromEvent } from "../services/api";

const MyEvents = () => {
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [message, setMessage] = useState("");

  useEffect(() => {
    const fetchUserEvents = async () => {
      try {
        const response = await getEventsForUser();
        setEvents(response.data.events);
        console.log(response.data.events);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchUserEvents();
  }, []);

  async function handleUnsub(id) {
    try {
      console.log("HERE");
      const response = await unsubFromEvent(id);
      setMessage(response.message);
    } catch (error) {
      setMessage(error.response?.message || "Error while unsubscribing");
    }
  }

  if (loading) {
    return <div style={styles.loading}>Loading events...</div>;
  }

  if (error) {
    return <div style={styles.error}>Error: {error}</div>;
  }

  return (
    <div style={styles.eventsContainer} className="event-application">
      <h1 style={styles.pageTitle}>Your Registered Events</h1>
      {events && events.length > 0 ? (
        <div style={styles.eventsList}>
          {events.map((event) => (
            <div style={styles.eventCard} key={event._id}>
              <div style={styles.eventHeader}>
                <h2 style={styles.eventTitle}>{event.name}</h2>
                <p style={styles.eventDate}>
                  {new Date(event.date).toLocaleDateString()}
                </p>
              </div>
              <div style={styles.eventLocation}>
                <p>
                  <strong>Location:</strong> {event.location.address}
                </p>
                <p>
                  {event.location.city}, {event.location.state},{" "}
                  {event.location.country}
                </p>
              </div>
              <div style={styles.eventOrganizers}>
                <h3 style={styles.subHeading}>Organizers:</h3>
                <ul>
                  {event.organizers.map((organizer, index) => (
                    <li key={index} style={styles.organizerItem}>
                      <p>
                        <strong>{organizer.name}</strong> - {organizer.role}
                      </p>
                      <p>Email: {organizer.email}</p>
                      <p>Phone: {organizer.phone}</p>
                    </li>
                  ))}
                </ul>
              </div>
              <div style={styles.eventCapacity}>
                <p>
                  <strong>Capacity:</strong> {event.min_capacity} -{" "}
                  {event.max_capacity}
                </p>
              </div>
              <button
                onClick={() => handleUnsub(event._id)}
                style={styles.unsubscribeButton}
              >
                Unsubscribe
              </button>
              <h1>{message}</h1>
            </div>
          ))}
        </div>
      ) : (
        <div>
          <h2 style={styles.noEvents}>You have no registered events.</h2>
          <p style={{ color: "#777", fontSize: "1rem" }}>
            Browse the event catalog to find something exciting!
          </p>
        </div>
      )}
    </div>
  );
};

const styles = {
  eventsContainer: {
    minHeight: "100vh",
    padding: "40px 20px",
    backgroundColor: "#f5f7fa",
    fontFamily: "Arial, sans-serif",
    textAlign: "center",
  },
  pageTitle: {
    fontSize: "2.5rem",
    fontWeight: "bold",
    color: "#333",
    marginBottom: "30px",
  },
  loading: {
    fontSize: "1.2rem",
    color: "#333",
    marginTop: "20px",
  },
  error: {
    fontSize: "1.2rem",
    color: "red",
    marginTop: "20px",
  },
  eventsList: {
    display: "flex",
    flexWrap: "wrap",
    gap: "20px",
    justifyContent: "center",
    marginTop: "20px",
  },
  eventCard: {
    backgroundColor: "white",
    borderRadius: "12px",
    boxShadow: "0 6px 12px rgba(0, 0, 0, 0.1)",
    width: "300px",
    padding: "20px",
    transition: "transform 0.3s ease, box-shadow 0.3s ease",
    cursor: "pointer",
    overflow: "hidden",
  },
  eventHeader: {
    marginBottom: "20px",
  },
  eventTitle: {
    fontSize: "1.5rem",
    fontWeight: "600",
    marginBottom: "10px",
  },
  eventDate: {
    color: "#666",
    fontSize: "1rem",
  },
  eventLocation: {
    marginBottom: "15px",
    fontSize: "1rem",
    color: "#333",
  },
  eventOrganizers: {
    marginBottom: "15px",
    textAlign: "left",
  },
  subHeading: {
    fontSize: "1.2rem",
    fontWeight: "bold",
    marginBottom: "10px",
  },
  organizerItem: {
    fontSize: "1rem",
    marginBottom: "10px",
    color: "#555",
  },
  eventCapacity: {
    marginTop: "10px",
    fontSize: "1rem",
    color: "#333",
  },
  unsubscribeButton: {
    marginTop: "15px",
    padding: "10px 15px",
    backgroundColor: "#ff4d4f",
    color: "#fff",
    border: "none",
    borderRadius: "5px",
    fontWeight: "bold",
  },
  noEvents: {
    fontSize: "1.5rem",
    color: "#555",
    marginTop: "20px",
    textAlign: "center",
  },
};

export default MyEvents;
