// src/components/Navbar.js
import React, { useContext } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import AuthContext from '../context/AuthContext';
import './Navbar.css';

function Navbar() {
  const { isAuthenticated, userRole, logout } = useContext(AuthContext);
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <nav>
      <ul>
        <li><Link to="/">Home</Link></li>
        {isAuthenticated && <li><Link to="/mis-cursos">Mis Cursos</Link></li>}
        {isAuthenticated && userRole === 'admin' && <li><Link to="/crear-curso">Crear Curso</Link></li>}
        {!isAuthenticated ? (
          <>
            <li><Link to="/login">Login</Link></li>
            <li><Link to="/register">Registro</Link></li>
          </>
        ) : (
          <li><button className="logout-button" onClick={handleLogout}>Cerrar Sesión</button></li>
        )}
      </ul>
    </nav>
  );
}

export default Navbar;