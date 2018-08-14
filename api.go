package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type PostData struct {
	Key    string
	Target string
}

func DatabaseAPI(db ShortlinkDB) http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/manage/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				return
			}

			var data PostData
			err = json.Unmarshal(body, &data)
			if err != nil || data.Key == "" || data.Target == "" {
				http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
				return
			}

			if err = db.Add(data.Key, data.Target); err != nil {
				http.Error(w, err.Error(), http.StatusPreconditionFailed)
			} else {
				w.WriteHeader(http.StatusCreated)
			}

		} else {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})
	return router
}
