package users

import (
	"net/http"
)

// @OAS handleCreateUser /users POST
func (s *service) handleCreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
