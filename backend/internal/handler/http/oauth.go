package handler

import (
	"errors"
	"io"
	"log"
	"math/rand"
	"net/http"
	"scoreboardpro/internal/entity"
	"scoreboardpro/internal/entity/model"

	"clevergo.tech/jsend"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63n(int64(len(letterBytes)))]
	}
	return string(b)
}

func (h *Handler) authViaOauth(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Body.Close()
	log.Println(string(b))

	accessToken := h.oauth.GetOauthAccessToken()
	userOauth := h.oauth.GetUserInfoByOauthToken(accessToken)

	userDB, err := h.userService.GetByEmail((*userOauth)["Email"])
	if err == nil {
		tokenString, err := h.auth.MakeAuthn(userDB.ID)
		if err != nil {
			jsend.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Error occured in handler.authViaOauth when make authentication:", err.Error())
			return
		}
		w.Header().Add("Authorization", "Bearer "+tokenString)

		data := map[string]string{
			"message": "successful login",
		}
		jsend.Success(w, data)
		return

	} else if errors.Is(err, entity.ErrUserNotFound) {
		err1 := h.registerUserViaOauth(userOauth)
		if err1 != nil {
			jsend.Error(w, err1.Error(), http.StatusBadRequest)
			return
		} else {
			data := map[string]string{
				"message": "successful register",
			}
			jsend.Success(w, data, http.StatusCreated)
			return
		}

	} else {
		jsend.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Error occured in handler.authViaOauth when login:", err.Error())
		return
	}
}

func (h *Handler) registerUserViaOauth(userOauth *map[string]string) error {
	userReg := model.UserRegister{
		UserLogin: model.UserLogin{
			Email:    (*userOauth)["Email"],
			Password: randString(32),
		},
		FirstName: (*userOauth)["FirstName"],
		LastName:  (*userOauth)["LastName"],
	}

	err := h.userService.Register(&userReg)
	if err != nil {
		log.Println("Error when creating user via oauth, error:", err.Error())
		return err
	} else {
		return nil
	}
}

func (h *Handler) loginViaOauth(w http.ResponseWriter, r *http.Request) {
	reqURL := h.oauth.CreateLinkForOAuthToken()
	http.Redirect(w, r, reqURL, http.StatusSeeOther)
}
