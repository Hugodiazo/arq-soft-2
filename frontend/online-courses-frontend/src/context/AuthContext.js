// src/context/AuthContext.js
import React, { createContext, useState, useEffect } from 'react';

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [userRole, setUserRole] = useState(null);

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (token) {
      setIsAuthenticated(true);
      const payload = JSON.parse(atob(token.split('.')[1])); // Decodificar el token JWT
      console.log("Rol del usuario:", payload.role); // Agrega este console.log
      setUserRole(payload.role); // Obtener y guardar el rol del usuario
    }
  }, []);

  const login = (token) => {
    localStorage.setItem('token', token);
    setIsAuthenticated(true);
    const payload = JSON.parse(atob(token.split('.')[1])); // Decodificar el token JWT
    setUserRole(payload.role);
  };

  const logout = () => {
    localStorage.removeItem('token');
    setIsAuthenticated(false);
    setUserRole(null);
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, userRole, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export default AuthContext;