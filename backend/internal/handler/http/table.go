package handler

import (
	"clevergo.tech/jsend"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) getTables(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tables, err := h.tableService.GetAll()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error when getting all users, Error: %v\n", err.Error())
			return
		}

		jsend.Success(w, tables, http.StatusOK)
	} else if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getTableById(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		tableId, err := strconv.ParseUint(vars["table_id"], 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Error when parsing table_id to uint, Error: %v", err.Error())
			return
		}

		tables, err := h.tableService.Get(uint(tableId))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error when getting all tables, Error: %v\n", err.Error())
			return
		}

		jsend.Success(w, tables, http.StatusOK)
	} else if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getTableByNameAndSeason(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		champName := vars["champ_name"]
		season := vars["season"]

		table, err := h.tableService.GetByNameAndSeason(champName, season)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error when getting all tables, Error: %v\n", err.Error())
			return
		}

		jsend.Success(w, table, http.StatusOK)
	} else if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
