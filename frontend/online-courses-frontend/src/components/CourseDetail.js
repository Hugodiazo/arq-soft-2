import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import API_URL from '../config';
import './CourseDetail.css';

const CourseDetail = () => {
  const { id } = useParams(); // Obtener el ID del curso desde la URL
  const [course, setCourse] = useState(null);

  useEffect(() => {
    const fetchCourse = async () => {
      try {
        const response = await fetch(`${API_URL}/courses/${id}`);
        const data = await response.json();
        setCourse(data);
      } catch (error) {
        console.error('Error al obtener detalles del curso:', error);
      }
    };

    fetchCourse();
  }, [id]);

  const handleEnrollment = async () => {
    try {
      const token = localStorage.getItem('token'); // Asegúrate de que el usuario esté autenticado
      const response = await fetch(`${API_URL}/courses/enroll`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`, // Usar el token de autenticación
        },
        body: JSON.stringify({ course_id: id }), // Enviar el ID del curso
      });

      if (response.ok) {
        alert('¡Inscripción exitosa!');
      } else {
        alert('Error al inscribirse en el curso');
      }
    } catch (error) {
      console.error('Error al inscribirse:', error);
    }
  };

  if (!course) return <p>Cargando detalles del curso...</p>;

  return (
    <div className="course-detail">
      <h2>{course.title}</h2>
      <p>{course.description}</p>
      <p>Instructor: {course.instructor}</p>
      <p>Duración: {course.duration} horas</p>
      <button onClick={handleEnrollment}>Inscribirme</button>
    </div>
  );
};

export default CourseDetail;