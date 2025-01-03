import React, { createContext, useState, useEffect } from "react";
import { getCookie, decodeToken } from "../services/api";

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [message, setMessage] = useState(null);

  useEffect(() => {
    if (!user) {
      const token = getCookie("token");
      if (token) {
        const decodedToken = decodeToken(token);
        if (decodedToken) {
          setUser({
            name: decodedToken.name,
            isAdmin: decodedToken.isAdmin,
            isActive: decodedToken.isActivated,
          });
        }
      }
    }
  }, [user]);

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
