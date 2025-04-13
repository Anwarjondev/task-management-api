package middleware

import (
	"net/http"

	"github.com/Anwarjondev/task-management-api/utils"
)


func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value("role").(string)
		if role != "admin" {
			utils.SendError(w, http.StatusForbidden, "Forbidden: Admins only")
			return
		}
		next.ServeHTTP(w, r)
	})
}