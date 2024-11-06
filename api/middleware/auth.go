package middleware

import (
	"net/http"

	"github.com/hugodiazo/arq-soft-2/api/users"
)

// CheckRole verifica si el usuario tiene el rol adecuado
func CheckRole(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extraer el userID del token
		userID, err := users.GetUserIDFromToken(r)
		if err != nil {
			http.Error(w, "No autorizado", http.StatusUnauthorized)
			return
		}

		// Obtener el usuario desde la base de datos
		user, err := users.GetUserByIDFromDB(userID)
		if err != nil || user.Role != requiredRole {
			http.Error(w, "No tienes permiso para acceder a esta ruta", http.StatusForbidden)
			return
		}

		// Si todo está bien, proceder al siguiente handler
		next(w, r)
	}
}
