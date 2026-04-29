package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"apigateway/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

// LoginHandler handles POST /api/auth/login
type LoginHandler struct {
	cfg model.JWTConfig
}

func NewLoginHandler(cfg model.JWTConfig) *LoginHandler {
	return &LoginHandler{cfg: cfg}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.Username != h.cfg.AdminUser || req.Password != h.cfg.AdminPass {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	expiration := time.Duration(h.cfg.Expiration) * time.Second
	if expiration == 0 {
		expiration = 24 * time.Hour
	}
	expiresAt := time.Now().Add(expiration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": req.Username,
		"iat": time.Now().Unix(),
		"exp": expiresAt.Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.cfg.Secret))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to sign token"})
		return
	}

	writeJSON(w, http.StatusOK, loginResponse{
		Token:     tokenString,
		ExpiresIn: int64(expiration.Seconds()),
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
