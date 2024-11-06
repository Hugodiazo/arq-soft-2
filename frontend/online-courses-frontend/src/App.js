import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Navbar from './components/Navbar';
import Home from './pages/Home';
import MyCourses from './pages/MyCourses';
import Login from './components/Login';
import Register from './components/Register';
import CourseDetail from './components/CourseDetail';
import Search from './components/Search'; 
import CreateCourse from './components/CreateCourse'; 
import { AuthProvider } from './context/AuthContext';
import PrivateRoute from './components/PrivateRoute';

function App() {
  return (
    <AuthProvider>
      <Router>
        <Navbar />
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/courses/:id" element={<PrivateRoute><CourseDetail /></PrivateRoute>} />
          <Route path="/mis-cursos" element={<PrivateRoute><MyCourses /></PrivateRoute>} />
          <Route path="/crear-curso" element={<PrivateRoute><CreateCourse /></PrivateRoute>} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;