import React, { useContext } from "react";
import { AuthContext } from "../contexts/AuthContext";

export const AlertMessage = () => {
  const { message } = useContext(AuthContext);

  if (!message) return null;

  return (
    <div
      className={`fixed top-0 left-0 right-0 ${
        message.type === "success" ? "bg-green-500" : "bg-red-500"
      } text-white p-4 text-center transition-all duration-300 ease-in-out transform z-50 ${
        message ? "translate-y-0" : "-translate-y-full"
      }`}
    >
      {message.text}
    </div>
  );
};
