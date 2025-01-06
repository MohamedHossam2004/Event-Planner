import React, { createContext, useState, useEffect } from "react";
import { getCookie, decodeToken } from "../services/api";

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(() => {
    // Initialize user state from token in cookies
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

  // Optional: Refresh token check
  useEffect(() => {
    const token = getCookie("token");
    if (!user && token) {
      const decodedToken = decodeToken(token);
      if (decodedToken) {
        setUser({
          name: decodedToken.name,
          isAdmin: decodedToken.isAdmin,
          isActive: decodedToken.isActivated,
        });
      }
    }
  }, []);

  const showMessage = (text, type) => {
    setMessage({ text, type });
    setTimeout(() => setMessage(null), 5000);
  };

  return (
    <AuthContext.Provider value={{ user, setUser, message, showMessage }}>
      {children}
    </AuthContext.Provider>
  );
};
