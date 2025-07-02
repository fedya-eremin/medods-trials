package api

import (
	"net/http"

	"example.com/m/internal/contextkeys"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Logout godoc
//
//	@Summary		Logout handler
//	@Description	logout via access token, i.e. remove refresh token entry from db
//	@Tags			api
//	@Security		Bearer
//	@Success		204
//	@Failure		500	{string}	string
//	@Router			/api/logout [post]
func (s *State) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := contextkeys.GetLogger(ctx)
	jwtClaims, _ := contextkeys.GetContextValue[*jwt.RegisteredClaims](ctx, contextkeys.JWTClaimsKey)

	jti, err := uuid.Parse(jwtClaims.ID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Error("error", "err", err)
		return
	}
	err = s.JWTService.DeleteRefreshTokenByID(ctx, jti)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Error("error", "err", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
