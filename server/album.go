package main

import (
	"encoding/json"
	"net/http"
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
		UserID:   s.user.ID,
		Title:    albumMetaData.Title,
		Mask:     mask,
		Pictures: pictures,
	}

	saveAlbum(album)

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
