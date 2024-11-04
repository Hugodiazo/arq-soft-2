// src/pages/MyCourses.js
import React, { useEffect, useState } from 'react';
import API_URL from '../config';

const MyCourses = () => {
  const [courses, setCourses] = useState([]);

  useEffect(() => {
    const fetchMyCourses = async () => {
      try {
        const token = localStorage.getItem('token');
        console.log("Token enviado en la solicitud:", token);

        const response = await fetch(`${API_URL}/enrollments`, {
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
          },
        });

        if (!response.ok) {
          throw new Error('Error al obtener tus cursos');
        }

        const data = await response.json();
        setCourses(data);
      } catch (error) {
        console.error('Error al obtener tus cursos:', error);
      }
    };

    fetchMyCourses();
  }, []);

  const handleUnenroll = async (courseId) => {
    try {
      const token = localStorage.getItem('token');
      const response = await fetch(`${API_URL}/courses/unenroll?course_id=${courseId}`, {
        method: 'DELETE',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
  
      if (response.ok) {
        alert('Desinscripción exitosa');
        setCourses(courses.filter(course => course.id !== courseId)); // Actualiza los cursos
      } else {
        alert('Error al desinscribirse del curso');
      }
    } catch (error) {
      console.error('Error al desinscribirse:', error);
    }
  };

  return (
    <div className="my-courses">
      {courses.length > 0 ? (
        courses.map((course, index) => (
          <div key={`${course.id}-${index}`} className="course-item">
            <h3>{course.title}</h3>
            <p>{course.description}</p>
            <button onClick={() => handleUnenroll(course.id)}>Desinscribirse</button>
          </div>
        ))
      ) : (
        <p>No estás inscrito en ningún curso.</p>
      )}
    </div>
  );
};

export default MyCourses;