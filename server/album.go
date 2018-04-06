package main

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/gorilla/mux"
)

func createAlbumHandler(w http.ResponseWriter, r *http.Request) {
	type albumMeta struct {
		Title    string   `json:"title"`
		Pictures []string `json:"pictureIDs"`
	}
	auth := r.Header.Get("token")

	s, ok := getSession(auth)
	if !ok {
		writeError(&w, "invalid session, please reload your page", 401)
		return
	}

	albumMetaData := &albumMeta{}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&albumMetaData); err != nil {
		writeError(&w, "invalid json format", 400)
		return
	}

	// TODO: Check for mask collisions
	mask := generateRandomString(32)

	pictures, err := getPictures(s.user, albumMetaData.Pictures, false)
	if err != nil {
		writeError(&w, "internal error", 500)
		l.Println(err)
		return
	}

	album := &Album{
		UserID: s.user.ID,
		Title:  albumMetaData.Title,
		Mask:   mask,
		// Pictures: pictures,
	}

	createAlbum(album, pictures)

	pictureMasks := make([]string, len(album.Pictures))
	for i, picture := range album.Pictures {
		pictureMasks[i] = picture.Mask
	}
	resp := struct {
		*jsonResponse
		Title    string   `json:"title"`
		Mask     string   `json:"albumID"`
		Pictures []string `json:"pictures"`
	}{
		&jsonResponse{
			Success: true,
		},
		album.Title,
		album.Mask,
		pictureMasks,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func getAlbumHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("token")
	vars := mux.Vars(r)
	albumID, ok := vars["albumID"]
	if !ok {
		writeError(&w, "no album ID provided", 400)
		return
	}
	s, ok := getSession(auth)
	if !ok {
		writeError(&w, "invalid session, please reload your page", 401)
		return
	}
	album, err := getAlbum(s.user, albumID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			writeError(&w, "album not found for your user or your album ID", 404)
			return
		}
		writeError(&w, "internal error, please try again later", 500)
		l.Print(err)
		return
	}

	type picResponse struct {
		Mask string `json:"pictureID"`
		URL  string `json:"url"`
	}
	albumPictures := make([]picResponse, len(album.Pictures))
	for i, pic := range album.Pictures {
		albumPictures[i].Mask = pic.Mask
		albumPictures[i].URL = pic.ValidURL
	}
	resp := struct {
		*jsonResponse
		Title    string        `json:"title"`
		Mask     string        `json:"albumID"`
		Pictures []picResponse `json:"pictures"`
	}{
		&jsonResponse{
			Success: true,
		},
		album.Title,
		album.Mask,
		albumPictures,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func deleteAlbumHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("token")
	vars := mux.Vars(r)
	albumID, ok := vars["albumID"]
	if !ok {
		writeError(&w, "no album ID provided", 400)
		return
	}
	s, ok := getSession(auth)
	if !ok {
		writeError(&w, "invalid session, please reload your page", 401)
		return
	}
	album, err := getAlbum(s.user, albumID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			writeError(&w, "album not found for your user or your album ID", 404)
			return
		}
		writeError(&w, "internal error, please try again later", 500)
		l.Print(err)
		return
	}
	if err := deleteAlbum(album); err != nil {
		writeError(&w, "error deleting album, please try again later", 500)
		l.Print(err)
		return
	}
	resp := struct {
		*jsonResponse
	}{
		&jsonResponse{
			Success: true,
			Message: "deleted album",
		},
	}
	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func getAllAlbumsHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("token")

	s, ok := getSession(auth)
	if !ok {
		writeError(&w, "invalid session, please reload your page", 401)
		return
	}
	albums, err := getAllAlbums(s.user)
	if err != nil {
		writeError(&w, "internal error retrieving album, try again later", 500)
		l.Println(err)
		return
	}
	type albumMeta struct {
		Title string `json:"title"`
		Mask  string `json:"albumID"`
	}
	meta := make([]albumMeta, len(albums))
	for i, album := range albums {
		meta[i].Title = album.Title
		meta[i].Mask = album.Mask
	}

	resp := struct {
		*jsonResponse
		Meta []albumMeta
	}{
		&jsonResponse{
			Success: true,
		},
		meta,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
