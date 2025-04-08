package middleware

import "net/http"


func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value("role").(string)
		if role != "admin" {
			http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}