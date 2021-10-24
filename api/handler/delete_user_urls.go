package handler

import (
	"encoding/json"
	"fmt"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
	"log"
	"net/http"
	"strconv"
)

func (h *URLShortenerHandler) handleDeleteUserURLs(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		msg := fmt.Sprintf("Expected Content-Type: 'application/json', but got [%s]", contentType)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	userID := userID(r)

	dec := json.NewDecoder(r.Body)
	urlIDsStr := make([]string, 0)
	if errDec := dec.Decode(&urlIDsStr); errDec != nil {
		msg := fmt.Sprintf("Cannot decode request body: %v", errDec)
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	urlsToDelete := make([]model.URLToDelete, 0, len(urlIDsStr))
	for _, urlIdStr := range urlIDsStr {
		urlID, err := strconv.Atoi(urlIdStr)
		if err != nil {
			msg := fmt.Sprintf("Cannot parse URL ID [%s]: %v", urlIdStr, err.Error())
			log.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		urlsToDelete = append(urlsToDelete, model.URLToDelete{UserID: userID, ID: urlID})
	}

	for _, u := range urlsToDelete {
		h.Service.ScheduleDeletion(u)
		log.Printf("Scheduled URL for deletion: %v", u)
	}

	w.WriteHeader(http.StatusAccepted)
}
