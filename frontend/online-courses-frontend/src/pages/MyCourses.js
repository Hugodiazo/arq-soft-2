import React, { useEffect, useState } from 'react';
import API_URL from '../config';

const MyCourses = () => {
  const [courses, setCourses] = useState([]);

  useEffect(() => {
    const fetchMyCourses = async () => {
      try {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_URL}/enrollments`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        const data = await response.json();
        setCourses(data);
      } catch (error) {
        console.error('Error al obtener mis cursos:', error);
      }
    };

    fetchMyCourses();
  }, []);

  if (courses.length === 0) return <p>No estás inscrito en ningún curso.</p>;

  return (
    <div className="my-courses">
      <h2>Mis Cursos</h2>
      {courses.map((course) => (
        <div key={course.id}>
          <h3>{course.title}</h3>
          <p>{course.description}</p>
        </div>
      ))}
    </div>
  );
};

export default MyCourses;