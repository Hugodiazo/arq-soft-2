// src/components/Navbar.js
import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './Navbar.css';

function Navbar() {
  const [logoutMessage, setLogoutMessage] = useState('');
  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem('token'); // Elimina el token
    setLogoutMessage('Has cerrado sesión con éxito'); // Muestra el mensaje de cierre de sesión
    navigate('/login'); // Redirige al usuario al inicio de sesión

    // Oculta el mensaje después de 3 segundos
    setTimeout(() => {
      setLogoutMessage('');
    }, 3000);
  };

  return (
    <nav>
      <ul>
        <li><Link to="/">Home</Link></li>
        <li><Link to="/search">Buscar Cursos</Link></li>
        <li><Link to="/mis-cursos">Mis Cursos</Link></li>
        <li><Link to="/login">Login</Link></li>
        <li><Link to="/register">Registro</Link></li>
      </ul>
      <button className="logout-button" onClick={handleLogout}>Cerrar Sesión</button>

      {/* Mostrar el mensaje de cierre de sesión si existe */}
      {logoutMessage && <p className="logout-message">{logoutMessage}</p>}
    </nav>
  );
}

export default Navbar;