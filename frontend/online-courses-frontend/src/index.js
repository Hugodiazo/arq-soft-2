// src/index.js
import React from 'react';
import ReactDOM from 'react-dom/client'; // Asegúrate de usar createRoot
import './index.css';
import App from './App';
import 'bootstrap/dist/css/bootstrap.min.css';

const root = ReactDOM.createRoot(document.getElementById('root')); // createRoot en lugar de render
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);