// HTTP handlers for pictures, including uploading, deletion, and tagging

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	minio "github.com/minio/minio-go"
)

const (
	tempFilePath = "temp_uploads/"
)

var (
	urlExpirationDuration = time.Minute * 60 * 24
)

// Protected POST
func pictureUploadHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("token")

	s, ok := getSession(auth)
	if !ok {
		writeError(&w, "invalid session, please reload your page", 401)
		return
	}

	err := r.ParseMultipartForm(2 << 21) // Superpage! :)
	if err != nil {
		l.Println(err.Error())
		writeError(&w, "invalid form data", 400)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		l.Println(err.Error())
		writeError(&w, "error getting photo data", 400)
		return
	}

	defer file.Close()

	f, err := os.OpenFile(handler.Filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	io.Copy(f, file)
	defer func() {
		path := f.Name()
		f.Close()
		os.Remove(path)
	}()

	// Upload to S3, then return a message w/ link
	url, bucketPath, err := uploadToS3(f.Name(), s)
	if err != nil {
		// TODO: make this graceful
		panic(err)
	}
	// TODO: check for mask collisions and retry with a new random generation!
	mask := generateRandomString(32)
	picture := Picture{
		ImagePath:      bucketPath,
		Mask:           mask,
		ValidURL:       url,
		ExpirationTime: time.Now().Add(urlExpirationDuration),
	}
	s.user.Pictures = append(s.user.Pictures, picture)
	saveUser(s.user)

	resp := struct {
		*jsonResponse
		URL       string `json:"url"`
		PictureID string `json:"pictureID"`
	}{
		&jsonResponse{
			Success: true,
		},
		url,
		mask,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Protected GET
func pictureRetrievalHandler(w http.ResponseWriter, r *http.Request) {
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

	pic, err := getPicture(s.user, pictureID)
	if err != nil {
		writeError(&w, "this picture does not exist", 410)
		return
	}

	resp := struct {
		*jsonResponse
		URL string `json:"url"`
	}{
		&jsonResponse{
			Success: true,
		},
		pic.ValidURL,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Protected DELETE
func pictureDeletionHandler(w http.ResponseWriter, r *http.Request) {
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

	pic, err := getPicture(s.user, pictureID)
	if err != nil {
		writeError(&w, "you cannot access this resource", 401)
		return
	}

	deleteFromS3(pic.ImagePath)
	deletePicture(pic.Mask)

	resp := jsonResponse{
		Success: true,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Protected DELETE
func massPictureDeletionHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: this
	type pictures struct {
		Pictures []string
	}

	auth := r.Header.Get("token")
	_, ok := getSession(auth)
	if !ok {
		writeError(&w, "invalid session, please reload your page", 403)
		return
	}

}

func deleteFromS3(path string) error {
	s3Client, err := minio.New(config.S3Endpoint, config.S3Key, config.S3Secret, true)
	if err != nil {
		return err
	}
	return s3Client.RemoveObject(config.S3SpaceName, path)
}

// TODO: change filename on the server in case duplicate names get uploaded!
func uploadToS3(fileName string, s *session) (string, string, error) {
	f, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		return "", "", err
	}
	defer f.Close()
	buf := make([]byte, 512)
	_, err = f.Read(buf)
	if err != nil && err != io.EOF {
		return "", "", err
	}

	contentType := http.DetectContentType(buf)

	s3Client, err := minio.New(config.S3Endpoint, config.S3Key, config.S3Secret, true)
	if err != nil {
		return "", "", err
	}

	exists, err := s3Client.BucketExists(config.S3SpaceName)
	if err != nil {
		return "", "", err
	}

	if !exists {
		if err = s3Client.MakeBucket(config.S3SpaceName, config.S3Location); err != nil {
			return "", "", err
		}
	}

	bucketPath := fmt.Sprintf("users/%d/%s", s.user.ID, f.Name())
	_, err = s3Client.FPutObject(config.S3SpaceName, bucketPath, f.Name(), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", "", err
	}

	url, err := s3Client.PresignedGetObject(config.S3SpaceName, bucketPath, urlExpirationDuration, nil)
	if err != nil {
		return "", "", err
	}

	return url.String(), bucketPath, nil
}

func refreshURL(pic *Picture) (string, error) {
	s3Client, err := minio.New(config.S3Endpoint, config.S3Key, config.S3Secret, true)
	if err != nil {
		return "", err
	}

	url, err := s3Client.PresignedGetObject(config.S3SpaceName, pic.ImagePath, urlExpirationDuration, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func init() {
	if _, err := os.Stat(tempFilePath); os.IsNotExist(err) {
		os.Mkdir(tempFilePath, os.ModePerm)
	}
}
