import React, { useContext, useEffect, useState } from "react";
import { AuthContext } from "../contexts/AuthContext";

export const AlertMessage = () => {
  const { message, setMessage } = useContext(AuthContext);
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    if (message) {
      setIsVisible(true);
      const timer = setTimeout(() => {
        setIsVisible(false);
        setTimeout(() => {
          setMessage(null);
        }, 500); // Wait for animation to complete before removing from DOM
      }, 3000);

      return () => clearTimeout(timer);
    }
  }, [message, setMessage]);

  if (!message) {
    return null;
  }

  const getBackgroundColor = (type) => {
    switch (type) {
      case "success":
        return "bg-green-500";
      case "error":
        return "bg-red-500";
      default:
        return "bg-gray-500";
    }
  };

  return (
    <div
      className={`fixed top-0 left-0 right-0 flex justify-center z-50 pointer-events-none transition-transform duration-500 ease-in-out ${
        isVisible ? "translate-y-0" : "-translate-y-full"
      }`}
    >
      <div
        className={`${getBackgroundColor(
          message.type,
        )} text-white px-6 py-2 rounded-b-lg shadow-lg max-w-md text-sm font-medium mt-0`}
      >
        {message.text}
      </div>
    </div>
  );
};
