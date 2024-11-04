// src/App.js
import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Navbar from './components/Navbar';
import Home from './pages/Home';
import MyCourses from './pages/MyCourses';
import Login from './components/Login';
import Register from './components/Register';
import CourseDetail from './components/CourseDetail';
import { AuthProvider } from './context/AuthContext'; // Asegúrate de importar y envolver con AuthProvider
import PrivateRoute from './components/PrivateRoute'; // Importa el PrivateRoute

function App() {
  return (
    <AuthProvider>
      <Router>
        <Navbar />
        <Routes>
          <Route path="/" element={<Home />} />
          {/* Rutas protegidas */}
          <Route
            path="/courses/:id"
            element={
              <PrivateRoute>
                <CourseDetail />
              </PrivateRoute>
            }
          />
          <Route
            path="/mis-cursos"
            element={
              <PrivateRoute>
                <MyCourses />
              </PrivateRoute>
            }
          />
          {/* Rutas públicas */}
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;