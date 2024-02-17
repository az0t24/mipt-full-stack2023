package handler

import (
	"net/http"
	"scoreboardpro/internal/entity"
	"strings"
)

const AUTHZ_HEADER_NAME = "Authorization"

func (h *Handler) getClaimsFromAuthHeader(r *http.Request) (*map[string]string, error) {
	jwt_claims := &map[string]string{}

	authzHeader := r.Header.Get(AUTHZ_HEADER_NAME)
	if authzHeader == "" {
		return jwt_claims, entity.ErrEmptyAuthHeader
	}

	headerParts := strings.Split(authzHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return jwt_claims, entity.ErrInvalidAuthHeader
	}

	jwt_claims, err := h.auth.FetchAuthn(headerParts[1])
	return jwt_claims, err
}
