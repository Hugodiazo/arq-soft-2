import React, { useState } from 'react';
import './Login.css';

function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');
  const [messageType, setMessageType] = useState('error');

  const handleSubmit = (e) => {
    e.preventDefault();
    fetch('http://localhost:8080/users/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.token) {
          localStorage.setItem('token', data.token);
          console.log("Token guardado:", data.token);
          setMessage('Inicio de sesión exitoso');
          setMessageType('success');
        } else {
          setMessage('Credenciales incorrectas');
          setMessageType('error');
        }
      })
      .catch((error) => {
        console.error('Error al iniciar sesión:', error);
        setMessage('Error al conectar con el servidor');
        setMessageType('error');
      });
  };

  return (
    <div className="login-form">
      <h2>Iniciar Sesión</h2>
      <p className={`login-message ${messageType}`}>{message}</p>
      <form onSubmit={handleSubmit}>
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
        <button type="submit">Iniciar sesión</button>
      </form>
    </div>
  );
}

export default Login;