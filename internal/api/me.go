package api

import (
	"encoding/json"
	"net/http"

	"github.com/fedya-eremin/medods-trials/internal/contextkeys"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Me godoc
//
//	@Summary		Me handler
//	@Description	get user by his access token
//	@Tags			api
//	@Security		Bearer
//	@Success		200	{object}	dto.User
//	@Failure		500	{string}	string
//	@Router			/api/me [get]
func (s *State) MeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := contextkeys.GetLogger(ctx)
	jwtClaims, _ := contextkeys.GetContextValue[*jwt.RegisteredClaims](ctx, contextkeys.JWTClaimsKey)

	userID, err := uuid.Parse(jwtClaims.Subject)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Error("error", "err", err)
		return
	}
	user, err := s.AuthService.GetUser(ctx, userID)
	reponseBody, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Error("error", "err", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(reponseBody)
}
