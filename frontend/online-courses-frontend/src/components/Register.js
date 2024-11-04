// src/components/Register.js
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './Register.css';

function Register() {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [role, setRole] = useState('user'); // Establece "user" como valor por defecto
  const [message, setMessage] = useState('');
  const navigate = useNavigate();

  const handleRegister = async () => {
    try {
      const response = await fetch('http://localhost:8080/users/register', {
        method: 'POST', // Asegúrate de que el método sea POST
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name, email, password, role }),
      });
  
      if (!response.ok) {
        throw new Error(`Error: ${response.statusText}`);
      }
  
      const data = await response.json();
      console.log('Usuario registrado con éxito:', data);
      alert('Usuario registrado con éxito');
    } catch (error) {
      console.error('Error al registrar usuario:', error);
      alert('Error al conectar con el servidor');
    }
  };

  return (
    <div className="register-form">
      <h2>Registrarse</h2>
      <p className="register-message">{message}</p>
      <form>
        <input
            type="text"
            placeholder="Nombre"
            value={name}
            onChange={(e) => setName(e.target.value)}
        />
        <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
        />
        <input
            type="password"
            placeholder="Contraseña"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
        />
        <select value={role} onChange={(e) => setRole(e.target.value)}>
            <option value="user">Usuario</option>
            <option value="admin">Administrador</option>
        </select>
        <button type="button" onClick={handleRegister}>
            Registrarse
        </button>
      </form>
    </div>
  );
}

export default Register;