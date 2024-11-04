import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import API_URL from '../config';
import './CourseDetail.css';

const CourseDetail = () => {
  const { id } = useParams();
  const [course, setCourse] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchCourse = async () => {
      try {
        const response = await fetch(`${API_URL}/courses/${id}`);
        if (!response.ok) {
          throw new Error('Error al obtener detalles del curso');
        }
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
      const token = localStorage.getItem('token');
      if (!token) {
        alert('Debes iniciar sesión para inscribirte en un curso');
        return;
      }

      const response = await fetch(`${API_URL}/courses/enroll`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ course_id: id }),
      });

      if (response.ok) {
        alert('¡Inscripción exitosa!');
        navigate('/mis-cursos');
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