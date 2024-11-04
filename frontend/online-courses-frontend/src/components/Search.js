// src/components/Search.js
import React, { useState } from 'react';
import './Search.css';

function Search() {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState([]);
  const [loading, setLoading] = useState(false);

  const handleSearch = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      const response = await fetch(`http://localhost:8080/search?q=${encodeURIComponent(query)}`);
      if (!response.ok) {
        throw new Error('Error al realizar la búsqueda');
      }

      const data = await response.json();
      setResults(data.response.docs); // Ajusta esto según la estructura de tu respuesta de Solr
    } catch (error) {
      console.error('Error al buscar cursos:', error);
      alert('Hubo un problema al realizar la búsqueda');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="search">
      <h2>Buscar Cursos</h2>
      <form onSubmit={handleSearch}>
        <input
          type="text"
          placeholder="Buscar cursos..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
        <button type="submit">Buscar</button>
      </form>
      {loading && <p>Cargando resultados...</p>}
      <div className="search-results">
        {results.length > 0 ? (
          results.map((course) => (
            <div key={course.id} className="course-item">
              <h3>{course.title}</h3>
              <p>{course.description}</p>
              <p>Instructor: {course.instructor}</p>
              <p>Duración: {course.duration} horas</p>
              <p>Nivel: {course.level}</p>
            </div>
          ))
        ) : (
          <p>No se encontraron resultados</p>
        )}
      </div>
    </div>
  );
}

export default Search;