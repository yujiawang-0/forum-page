package auth

// authMiddleware enforces authentication on other endpoints
// therefore it runs before every protected request (create, delete, edit post/topics/comments)

import (
	"context"
	//"fmt"
	"net/http"
	"strings"

	//"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/yujiawang-0/forum-page/internal/api"
)

// Define a custom type for context keys to avoid collisions
type contextKey string

const (
	userIDKey contextKey = "userID"
	RoleKey contextKey = "role"
)

// requireAuth is a middleware that checks for a valid user session.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		// In a real application, validate a token (e.g., from a cookie or header)
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			api.ErrorJSON(w, errors.New("missing authorisation header"), http.StatusUnauthorized)
			return 
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			api.ErrorJSON(w, errors.New("invalid authorization header"), http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// and fetch user details from a database.
		claims, err := verifyToken(tokenString)
		if err != nil {
			api.ErrorJSON(w, errors.New("invalid or expired token"), http.StatusUnauthorized)
			return 
		}

		// Create a new context with the user and replace the request context
		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		ctx = context.WithValue(ctx, RoleKey, claims.Role)

		r = r.WithContext(ctx)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// GetUserFromContext is a helper to retrieve the user from the context.
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(userIDKey).(int)
	return id, ok
}

// // profileHandler uses the authenticated user from the context.
// func profileHandler(w http.ResponseWriter, r *http.Request) {
// 	user := GetUserFromContext(r.Context())
// 	if user == nil {
// 		http.Error(w, "User not found in context", http.StatusInternalServerError)
// 		return
// 	}
// 	w.Write([]byte(fmt.Sprintf("Welcome, %s (%s)!", user.Email, user.ID)))
// }

// func main() {
// 	r := chi.NewRouter()

// 	// Public routes
// 	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("Public homepage"))
// 	})

// 	// Protected routes group
// 	r.Group(func(r chi.Router) {
// 		// Apply the requireAuth middleware
// 		r.Use(requireAuth)

// 		r.Get("/profile", profileHandler)
// 		// Add other protected routes
// 	})

// 	http.ListenAndServe(":8080", r)
//}
