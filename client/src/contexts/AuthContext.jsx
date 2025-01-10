import React, { createContext, useState, useCallback } from "react";
import { getCookie, decodeToken } from "../services/api";

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(() => {
    const token = getCookie("token");
    if (token) {
      const decodedToken = decodeToken(token);
      if (decodedToken) {
        return {
          name: decodedToken.name,
          isAdmin: decodedToken.isAdmin,
          isActive: decodedToken.isActivated,
        };
      }
    }
    return null;
  });
  const [message, setMessage] = useState(null);

  const updateUser = useCallback((token) => {
    if (token) {
      const decodedToken = decodeToken(token);
      if (decodedToken) {
        setUser({
          name: decodedToken.name,
          isAdmin: decodedToken.isAdmin,
          isActive: decodedToken.isActivated,
        });
      }
    } else {
      setUser(null);
    }
  }, []);

  const showMessage = useCallback((text, type) => {
    // First clear any existing message
    setMessage(null);

    // Set new message after a brief delay
    setTimeout(() => {
      setMessage({ text, type });
    }, 100);
  }, []);

  return (
    <AuthContext.Provider
      value={{ user, setUser, updateUser, message, setMessage, showMessage }}
    >
      {children}
    </AuthContext.Provider>
  );
};
