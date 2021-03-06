// TODO: On reboot, reload sessionMap based on database for persistency

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/gorilla/mux"
	"github.com/paul-io/go-expiring-map"
)

type jsonResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message,omitempty"`
	ResponseCode int    `json:"-"`
}

type session struct {
	user *User
}

const (
	sessionTimeToExpire = 2 // 2 hours
)

var (
	// sessionMap = make(map[string]session)
	sessionMap *expiring.ThreadSafeMap
)

type routes []route
type route struct {
	Method      string
	Pattern     string
	Validation  bool
	HandlerFunc http.HandlerFunc
}

func initServer() {
	// initialize expiring map
	sessionMap = expiring.NewThreadSafe(sessionTimeToExpire, time.Hour, true)

	r := mux.NewRouter().StrictSlash(true)
	for _, route := range apiRoutes {
		handler := route.HandlerFunc
		if route.Validation {
			handler = validationHandler(handler)
		}
		handler = logHandler(handler)

		r.Methods(route.Method).Path(route.Pattern).Handler(handler)
	}

	log.Printf("Listening on port %d\n", config.Port)
	go http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r)
}

var apiRoutes = routes{
	// User routes
	route{
		Method:      "POST",
		Pattern:     "/user/login",
		Validation:  false,
		HandlerFunc: loginHandler,
	},
	route{
		Method:      "POST",
		Pattern:     "/user/register",
		Validation:  false,
		HandlerFunc: registrationHandler,
	},
	route{
		Method:      "POST",
		Pattern:     "/user/logout",
		Validation:  true,
		HandlerFunc: logoutHandler,
	},
	route{
		Method:      "GET",
		Pattern:     "/user/test",
		Validation:  true,
		HandlerFunc: validateUserLoggedIn,
	},
	route{
		Method:      "GET",
		Pattern:     "/user/pictures",
		Validation:  true,
		HandlerFunc: getUsersPicturesHandler,
	},

	// Picture routes
	route{
		Method:      "POST",
		Pattern:     "/picture/upload",
		Validation:  true,
		HandlerFunc: pictureUploadHandler,
	},
	route{
		Method:      "GET",
		Pattern:     "/picture/{pictureID}",
		Validation:  true,
		HandlerFunc: pictureRetrievalHandler,
	},
	route{
		Method:      "DELETE",
		Pattern:     "/picture/{pictureID}",
		Validation:  true,
		HandlerFunc: pictureDeletionHandler,
	},
	route{
		Method:      "POST",
		Pattern:     "/picture/deletes",
		Validation:  true,
		HandlerFunc: massPictureDeletionHandler,
	},

	// Tag routes
	route{
		Method:      "POST",
		Pattern:     "/picture/{pictureID}/tags",
		Validation:  true,
		HandlerFunc: addTagHandler,
	},
	route{
		Method:      "DELETE",
		Pattern:     "/picture/{pictureID}/tags",
		Validation:  true,
		HandlerFunc: removeTagHandler,
	},
	route{
		Method:      "GET",
		Pattern:     "/picture/{pictureID}/tags",
		Validation:  true,
		HandlerFunc: getTagHandler,
	},
	route{
		Method:      "POST",
		Pattern:     "/picture/search",
		Validation:  true,
		HandlerFunc: tagFuzzySearch,
	},

	// Album routes
	route{
		Method:      "POST",
		Pattern:     "/album",
		Validation:  true,
		HandlerFunc: createAlbumHandler,
	},
	route{
		Method:      "GET",
		Pattern:     "/album/{albumID}",
		Validation:  true,
		HandlerFunc: getAlbumHandler,
	},
	route{
		Method:      "DELETE",
		Pattern:     "/album/{albumID}",
		Validation:  true,
		HandlerFunc: deleteAlbumHandler,
	},
	route{
		Method:      "GET",
		Pattern:     "/album",
		Validation:  true,
		HandlerFunc: getAllAlbumsHandler,
	},
}

func logHandler(h http.HandlerFunc) http.HandlerFunc {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in [%s]\n", getFuncName(h))
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		h.ServeHTTP(w, r)
		log.Printf("[%s] %s: %s\n", getFuncName(h), r.RequestURI, time.Since(t1))
	})
}

func validationHandler(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("token")
		if len(auth) == 0 {
			errorHandler(h, "authorization required", 401).ServeHTTP(w, r)
			return
		}
		// TODO: something with this session (pass on the user through the middleware)
		_, ok := getSession(auth)
		if !ok {
			errorHandler(h, "invalid token", 401).ServeHTTP(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func errorHandler(h http.HandlerFunc, message string, code int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeError(&w, message, code)
	})
}

func writeError(w *http.ResponseWriter, message string, code int) {
	resp := jsonResponse{
		Success:      false,
		Message:      message,
		ResponseCode: code,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	http.Error(*w, string(data), code)
}

func getSession(key string) (sess *session, ok bool) {
	s, ok := sessionMap.Get(key)
	if !ok {
		u, err := getUserFromSession(key)
		if err != nil {
			return nil, false
		}
		ok = true
		sess = &session{
			user: u,
		}
		sessionMap.Set(key, sess)
		return
	}
	sess = s.(*session)
	return
}

func reloadUserSessionFromToken(key string) {
	if sessionMap.Remove(key) {
		getSession(key)
	}
}

func getFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
