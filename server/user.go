// HTTP handlers for User tasks, such as login, registration, and logout

package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Unprotected POST
func loginHandler(w http.ResponseWriter, r *http.Request) {
	type login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	l := login{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		writeError(&w, "invalid json format", 400)
		return
	}
	user, err := getUserFromEmail(l.Email)
	if err != nil {
		writeError(&w, err.Error(), 401)
		return
	}
	if hash(l.Password, user.Salt) == user.Password {
		// success
		sessionID := generateRandomString(tokenLength)
		resp := struct {
			*jsonResponse
			Token string `json:"token"`
		}{
			&jsonResponse{
				Success: true,
			},
			sessionID,
		}
		sessionMap[sessionID] = session{
			user: user,
		}
		setUserToken(user, sessionID)
		data, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	} else {
		// failure
		writeError(&w, "invalid email or password", 401)
	}
}

// Unprotected POST
func registrationHandler(w http.ResponseWriter, r *http.Request) {
	type registration struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	reg := registration{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		writeError(&w, "invalid json format", 400)
		return
	}
	if len(reg.Name) == 0 || len(reg.Email) == 0 || len(reg.Password) == 0 {
		writeError(&w, "invalid json format", 400)
		return
	}

	salt := generateRandomString(saltLength)
	user := &User{
		Name:     reg.Name,
		Email:    reg.Email,
		Password: hash(reg.Password, salt),
		Salt:     salt,
		Oauth:    0,
	}

	token, err := registerUser(user)
	if err != nil {
		writeError(&w, err.Error(), 400)
		return
	}

	resp := struct {
		*jsonResponse
		Token string `json:"token"`
	}{
		&jsonResponse{
			Success: true,
		},
		token,
	}
	sessionMap[token] = session{
		user: user,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Protected POST
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("token")

	s, ok := sessionMap[auth]
	if !ok {
		writeError(&w, "you are already logged out!", 403)
		return
	}
	delete(sessionMap, auth)
	logoutUser(s.user)
	resp := jsonResponse{
		Success: true,
		Message: "you have successfully logged out",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Protected GET
func getUsersPicturesHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("token")

	s, ok := sessionMap[auth]
	if !ok {
		writeError(&w, "invalid session, please reload your page", 403)
		return
	}
	const (
		maxPhotosPerPage = 30
	)
	var (
		page  = 1
		limit = maxPhotosPerPage
	)

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	if len(pageStr) != 0 {
		if p, err := strconv.Atoi(pageStr); err != nil {
			page = 1
		} else {
			page = p
		}
	}
	if len(limitStr) != 0 {
		if l, err := strconv.Atoi(limitStr); err == nil {
			if l < maxPhotosPerPage {
				limit = l
			}
		}
	}

	type paginatedPicture struct {
		Mask string `json:"pictureID"`
		URL  string `json:"url"`
	}

	pictures, page, maxPages := getUsersPicturesAndRefreshURL(s.user, limit, page)
	paginated := make([]paginatedPicture, len(pictures))

	var i int
	for _, p := range pictures {
		paginated[i].Mask = p.Mask
		paginated[i].URL = p.ValidURL
		i++
	}

	resp := struct {
		*jsonResponse
		Pictures []paginatedPicture `json:"pictures"`
		Count    int                `json:"count"`
		Page     int                `json:"page"`
		MaxPage  int                `json:"maxPage"`
	}{
		&jsonResponse{
			Success: true,
		},
		paginated,
		len(pictures),
		page,
		maxPages,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Protected (can assume that header.Get("token") will always be valid)
func validateUserLoggedIn(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("token")

	s, ok := sessionMap[auth]
	if !ok {
		writeError(&w, "invalid session, please reload your page", 403)
		return
	}
	resp := jsonResponse{
		Success: true,
		Message: fmt.Sprintf("Good day, %s!", s.user.Name),
	}

	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func hash(password, salt string) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s+%s", password, salt)))
	return fmt.Sprintf("%x", h.Sum(nil))
}
