// src/pages/Home.js
import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom'; // Importa Link de react-router-dom
import './Home.css';

function Home() {
  const [query, setQuery] = useState('');
  const [courses, setCourses] = useState([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const fetchCourses = async () => {
      setLoading(true);
      try {
        const response = await fetch('http://localhost:8080/courses');
        if (!response.ok) {
          throw new Error('Error al obtener los cursos');
        }
        const data = await response.json();
        setCourses(data);
      } catch (error) {
        console.error('Error al obtener los cursos:', error);
        alert('Hubo un problema al obtener los cursos');
      } finally {
        setLoading(false);
      }
    };

    fetchCourses();
  }, []);

  const handleSearch = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      const response = query.trim()
        ? await fetch(`http://localhost:8080/search?q=${encodeURIComponent(query)}`)
        : await fetch('http://localhost:8080/courses');

      if (!response.ok) {
        throw new Error('Error al realizar la búsqueda');
      }

      const data = await response.json();
      setCourses(query.trim() ? data.response.docs : data);
    } catch (error) {
      console.error('Error al buscar cursos:', error);
      alert('Hubo un problema al realizar la búsqueda');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="home">
      <h2>Bienvenido a la Plataforma de Cursos</h2>
      <form onSubmit={handleSearch} className="search-form">
        <input
          type="text"
          placeholder="Buscar cursos..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
        <button type="submit">Buscar</button>
      </form>
      {loading && <p>Cargando...</p>}
      <div className="courses">
        {courses.length > 0 ? (
          courses.map((course) => (
            <div key={course.id} className="course-item">
              {/* Usa Link para navegar al detalle del curso */}
              <Link to={`/courses/${course.id}`}>
                <h3>{course.title}</h3>
              </Link>
              <p>{course.description}</p>
              <p>Instructor: {course.instructor}</p>
              <p>Duración: {course.duration} horas</p>
              <p>Nivel: {course.level}</p>
            </div>
          ))
        ) : (
          <p>No se encontraron cursos.</p>
        )}
      </div>
    </div>
  );
}

export default Home;