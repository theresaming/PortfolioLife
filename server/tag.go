package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Protected POST
func addTagHandler(w http.ResponseWriter, r *http.Request) {
	type tagReq struct {
		Tags []string `json:"tags"`
	}
	tagReqs := new(tagReq)
	if err := json.NewDecoder(r.Body).Decode(&tagReqs); err != nil {
		writeError(&w, "invalid json format", 400)
		return
	}
	auth := r.Header.Get("token")
	vars := mux.Vars(r)
	pictureID, ok := vars["pictureID"]
	if !ok {
		writeError(&w, "no picture ID provided", 400)
		return
	}
	s, ok := getSession(auth)
	if !ok {
		writeError(&w, "invalid session, please reload your page", 403)
		return
	}
	pic, err := getPicture(s.user, pictureID, false)
	if err != nil {
		// TODO: check this error code
		writeError(&w, "invalid picture ID", 403)
		return
	}
	for _, tag := range tagReqs.Tags {
		t := Tag{Tag: tag}
		pic.Tags = append(pic.Tags, t)
	}
	savePictureTags(pic)
}

// Protected DELETE
func removeTagHandler(w http.ResponseWriter, r *http.Request) {

}
