package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"scoreboardpro/internal/entity"
	"scoreboardpro/internal/service"
	"scoreboardpro/pkg/auth"
	"scoreboardpro/pkg/oauth"
)

type Handler struct {
	// service *service.Service
	userService  entity.UserService
	tableService entity.TableService
	auth         entity.AuthManager
	oauth        entity.OAuthManager
}

func NewHandler(s *service.Service, auth *auth.AuthManager, oauth *oauth.OAuthManager) *Handler {
	return &Handler{userService: s.User, tableService: s.Table, auth: auth, oauth: oauth}
}

func (h *Handler) NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", index).Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/user", h.getAllUsers).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/user/{user_id}", h.getUserByID).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/user/register", h.registerUser).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/user/login", h.loginUser).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/user/{user_id}/mark/{table_id}", h.markAsFavorite).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/user/{user_id}/favorites", h.getFavorites).Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/oauth/login", h.loginViaOauth).Methods(http.MethodGet)
	r.HandleFunc("/oauth/auth", h.authViaOauth).Methods(http.MethodPost)

	r.HandleFunc("/tables/all", h.getTables).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/tables/{table_id}", h.getTableById).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/tables/{champ_name}/{season}", h.getTableByNameAndSeason).Methods(http.MethodGet, http.MethodOptions)

	r.Use(LoggingMiddleware(r))
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(CustomCORSMiddleware(r))

	return r
}

func LoggingMiddleware(r *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			log.Printf("Origin: %s | Forwarded: %s | Method: %s | RequestURI: %s", req.Header.Get("Origin"), req.Header.Get("Forwarded"), req.Method, req.RequestURI)

			next.ServeHTTP(w, req)
		})
	}
}

func CustomCORSMiddleware(r *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Length, Content-Type, Authorization, Host, Origin, X-CSRF-Token")
			w.Header().Set("Access-Control-Expose-Headers", "Authorization")

			next.ServeHTTP(w, req)
		})
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!\n")
}
