package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"scoreboardpro/internal/entity"
	"scoreboardpro/internal/entity/model"
	"strconv"

	"clevergo.tech/jsend"
	"github.com/gorilla/mux"
)

func (h *Handler) getAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		users, err := h.userService.GetAll()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error when getting all users, Error: %v\n", err.Error())
			return
		}

		jsend.Success(w, users, http.StatusOK)
	} else if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		authClaims, _ := h.getClaimsFromAuthHeader(r)
		authUserID, _ := strconv.ParseUint((*authClaims)["sub"], 10, 32)

		vars := mux.Vars(r)
		userID, err := strconv.ParseUint(vars["user_id"], 10, 32)
		if err != nil {
			log.Printf("Error when parsing user_id to uint, Error: %v", err.Error())
			return
		}

		if authUserID != userID {
			jsend.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		user, err := h.userService.Get(uint(userID))
		if err != nil {
			switch {
			case errors.Is(err, entity.ErrUserNotFound):
				jsend.Error(w, err.Error(), http.StatusNotFound)
			default:
				jsend.Error(w, err.Error(), http.StatusInternalServerError)
				log.Println("Error occureed in handler.getUserByID, Error:", err.Error())
			}
			return
		}

		jsend.Success(w, user, http.StatusOK)
	} else if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) registerUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var userReg model.UserRegister
		err := json.NewDecoder(r.Body).Decode(&userReg)
		if err != nil {
			jsend.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("Error occured in handler.registerUser when decoding userReg:", err.Error())
			return
		}

		err = h.userService.Register(&userReg)
		if err != nil {
			switch {
			case errors.Is(err, entity.ErrUserExists) ||
				errors.Is(err, entity.ErrInvalidEmail) ||
				errors.Is(err, entity.ErrInvalidPassword):
				jsend.Error(w, err.Error(), http.StatusBadRequest)
				return
			default:
				jsend.Error(w, err.Error(), http.StatusBadRequest)
				log.Println("Error occured in handler.registerUser when register userReg:", err.Error())
				return
			}
		}

		json_msg := map[string]string{
			"message": "successful registration",
		}
		jsend.Success(w, json_msg, http.StatusCreated)
	} else if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		userLogin := model.UserLogin{}
		err := json.NewDecoder(r.Body).Decode(&userLogin)
		if err != nil {
			jsend.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("Error when trying to decode userLogin to login an user, Error: %v\n", err.Error())
			return
		}

		userID, err := h.userService.Login(&userLogin)
		if err != nil {
			jsend.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("Error occured in handler.loginUser when login userLogin:", err.Error())
			return
		}

		tokenString, err := h.auth.MakeAuthn(userID)
		if err != nil {
			jsend.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Error occured in handler.loginUser when make authentication:", err.Error())
			return
		}

		w.Header().Add("Authorization", "Bearer "+tokenString)

		json_msg := map[string]string{
			"message": "successful login",
		}
		jsend.Success(w, json_msg, http.StatusOK)
	} else if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) markAsFavorite(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		authClaims, _ := h.getClaimsFromAuthHeader(r)
		authUserID, _ := strconv.ParseUint((*authClaims)["sub"], 10, 32)

		vars := mux.Vars(r)
		userID, err := strconv.ParseUint(vars["user_id"], 10, 32)
		if err != nil {
			log.Printf("Error when parsing user_id to uint, Error: %v", err.Error())
			return
		}

		if authUserID != userID {
			jsend.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		tableId, err := strconv.ParseUint(vars["table_id"], 10, 32)
		if err != nil {
			log.Printf("Error when parsing table_id to uint, Error: %v", err.Error())
			return
		}

		err = h.userService.MarkAsFavorite(uint(userID), uint(tableId))
		if err != nil {
			jsend.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Error occureed in handler.markAsFavorite, Error:", err.Error())
			return
		}

		json_msg := map[string]string{
			"message": "successful marking",
		}
		jsend.Success(w, json_msg, http.StatusOK)
	} else if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getFavorites(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		authClaims, _ := h.getClaimsFromAuthHeader(r)
		authUserID, _ := strconv.ParseUint((*authClaims)["sub"], 10, 32)

		vars := mux.Vars(r)
		userID, err := strconv.ParseUint(vars["user_id"], 10, 32)
		if err != nil {
			log.Printf("Error when parsing user_id to uint, Error: %v", err.Error())
			return
		}

		if authUserID != userID {
			jsend.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		tables, err := h.userService.GetFavorites(uint(userID))
		if err != nil {
			jsend.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Error occureed in handler.markAsFavorite, Error:", err.Error())
			return
		}

		jsend.Success(w, tables, http.StatusOK)
	} else if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
