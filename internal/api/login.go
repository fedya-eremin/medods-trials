package api

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"

	"github.com/fedya-eremin/medods-trials/internal/contextkeys"
	"github.com/fedya-eremin/medods-trials/internal/dto"
	"github.com/jackc/pgx/v5"
)

// Login godoc
//
//	@Summary		Login handler
//	@Description	exchange UUID for access and refresh token pair
//	@Tags			api
//	@Accept			json
//	@Return			json
//	@Param			user	body		dto.User	true	"user uuid"
//	@Success		201		{object}	dto.TokenPair
//	@Failure		422		{string}	string	"body is unprocessable"
//	@Failure		500		{object}	string
//	@Router			/api/login [post]
func (a *State) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := contextkeys.GetLogger(ctx)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Error("error", "err", err)
		return
	}
	defer r.Body.Close()
	var user dto.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Unprocessable entity", http.StatusUnprocessableEntity)
		logger.Error("error", "err", err)
		return
	}
	dbUser, err := a.AuthService.GetUser(ctx, user.ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		dbUser, err = a.AuthService.CreateUser(ctx, user.ID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			logger.Error("error", "err", err)
			return
		}
	case err != nil:
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Error("error", "err", err)
		return
	}

	clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Error("error", "err", err)
		return
	}

	tokenPair, err := a.JWTService.GenerateTokenPair(
		ctx, dbUser.ID, r.Header.Get("User-Agent"), clientIP,
	)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Error("error", "err", err)
		return
	}

	reponseBody, err := json.Marshal(tokenPair)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Error("error", "err", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(reponseBody)
}
