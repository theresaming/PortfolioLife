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
		writeError(&w, "invalid session, please reload your page", 401)
		return
	}
	pic, err := getPicture(s.user, pictureID, false)
	if err != nil {
		writeError(&w, "this picture does not exist for your account", 410)
		return
	}

	tags := make([]Tag, len(tagReqs.Tags))
	i := 0
	for _, tag := range tagReqs.Tags {
		tags[i] = Tag{Tag: tag}
		i++
	}
	createTags(pic, tags)
	resp := struct {
		*jsonResponse
	}{
		&jsonResponse{
			Success: true,
			Message: "added tags to photo",
		},
	}
	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Protected DELETE
func removeTagHandler(w http.ResponseWriter, r *http.Request) {
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
		writeError(&w, "invalid session, please reload your page", 401)
		return
	}
	pic, err := getPicture(s.user, pictureID, false)
	if err != nil {
		writeError(&w, "this picture does not exist for your account", 410)
		return
	}
	deleteTags(pic, tagReqs.Tags)
	resp := struct {
		*jsonResponse
	}{
		&jsonResponse{
			Success: true,
			Message: "Deleted tags from photo",
		},
	}
	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Protected GET
func getTagHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("token")
	vars := mux.Vars(r)
	pictureID, ok := vars["pictureID"]
	if !ok {
		writeError(&w, "no picture ID provided", 400)
		return
	}
	s, ok := getSession(auth)
	if !ok {
		writeError(&w, "invalid session, please reload your page", 401)
		return
	}
	pic, err := getPicture(s.user, pictureID, false)
	if err != nil {
		writeError(&w, "this picture does not exist for yor account", 410)
		return
	}
	tags, err := getTags(pic)
	if err != nil {
		writeError(&w, "unspecified error", 500)
		l.Println(err)
		return
	}
	tagResp := make([]string, len(tags))

	for i, tag := range tags {
		tagResp[i] = tag.Tag
	}
	resp := struct {
		*jsonResponse
		Tags []string `json:"tags"`
	}{
		&jsonResponse{
			Success: true,
		},
		tagResp,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func tagFuzzySearch(w http.ResponseWriter, r *http.Request) {
	type query struct {
		Search string `json:"search"`
	}
	tagReqs := new(query)
	if err := json.NewDecoder(r.Body).Decode(&tagReqs); err != nil {
		writeError(&w, "invalid json format", 400)
		return
	}
	auth := r.Header.Get("token")
	_, ok := getSession(auth)
	if !ok {
		writeError(&w, "invalid session, please reload your page", 401)
		return
	}
}
