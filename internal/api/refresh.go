package api

import (
	"encoding/json"
	"io"
	"net"
	"net/http"

	"github.com/fedya-eremin/medods-trials/internal/contextkeys"
	"github.com/fedya-eremin/medods-trials/internal/dto"
)

// Refresh godoc
//
//	@Summary		Refresh handler
//	@Description	exchange refresh token for access and refresh token pair. if ip differs from initial, webhook request issued
//	@Tags			api
//	@Accept			json
//	@Return			json
//	@Param			user	body		dto.RefreshRequest	true	"request with refresh token"
//	@Success		201		{object}	dto.TokenPair
//	@Failure		403		{string}	string	"token is used, expired, or User-Agent is different"
//	@Failure		409		{string}	string	"wrong format of refresh token"
//	@Failure		422		{string}	string	"body is unprocessable"
//	@Failure		500		{string}	string
//	@Router			/api/refresh [post]
func (s *State) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := contextkeys.GetLogger(ctx)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Error("error", "err", err)
		return
	}
	defer r.Body.Close()
	var refreshRequest dto.RefreshRequest
	err = json.Unmarshal(body, &refreshRequest)
	if err != nil {
		http.Error(w, "Unprocessable entity", http.StatusUnprocessableEntity)
		logger.Error("error", "err", err)
		return
	}

	clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Error("error", "err", err)
		return
	}
	err = s.JWTService.UseRefreshToken(ctx, refreshRequest.RefreshToken, r.Header.Get("User-Agent"), clientIP)
	switch err.(type) {
	case *dto.TokenUsedError:
		http.Error(w, "Token has been already used", http.StatusForbidden)
		return
	case *dto.TokenExpiredError:
		http.Error(w, "Token has been expired", http.StatusForbidden)
		return
	case *dto.WrongUserAgent:
		http.Error(w, "Wrong User-Agent", http.StatusForbidden)
		return
	case error:
		http.Error(w, "Wrong refresh token format", http.StatusConflict)
		return
	}

	user, err := s.JWTService.GetUserByRefreshToken(ctx, refreshRequest.RefreshToken)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Error("error", "err", err)
		return
	}

	tokenPair, err := s.JWTService.GenerateTokenPair(ctx, user.ID, r.Header["User-Agent"][0], clientIP)
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
