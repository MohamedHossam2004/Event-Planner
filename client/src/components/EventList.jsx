import React from "react";
import { EventCard } from "./EventCard";

export const EventList = ({ events, onEventSelect }) => {
  return (
    <div className="max-w-7xl mx-auto grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 px-4 mt-8">
      {events.map((event) => (
        <EventCard key={event.id} event={event} onSelect={onEventSelect} />
      ))}
    </div>
  );
};
