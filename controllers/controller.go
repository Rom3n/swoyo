package controllers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"swoyo/models"
	"swoyo/utils"

	"go.mongodb.org/mongo-driver/bson"
)

func EncodeUrls(w http.ResponseWriter, r *http.Request) {
	var url models.URL
	err := json.NewDecoder(r.Body).Decode(&url)
	if isError := utils.HandleHttpErrors(w, "", http.StatusInternalServerError, err); isError {
		return
	}

	err = utils.ValidateUrl(url.UrlString)
	if isError := utils.HandleHttpErrors(w, "invalid url", http.StatusOK, err); isError {
		return
	}

	collection, err := utils.GetCollection(utils.DatabaseClient, utils.ENCODED_URLS)
	if isError := utils.HandleHttpErrors(w, "could not fetch records", http.StatusInternalServerError, err); isError {
		return
	}

	var base64EncodedUrl string
	host := "http://localhost:8000/"
	base64EncodedUrl = base64.StdEncoding.EncodeToString([]byte(url.UrlString))

	encodedUrl := models.EncodeUrls{
		EncodeUrl: base64EncodedUrl,
		ShortUrl:  host + base64EncodedUrl[0:4],
	}

	_, err = collection.InsertOne(context.TODO(), encodedUrl)
	if isError := utils.HandleHttpErrors(w, "unable to insert", http.StatusInternalServerError, err); isError {
		return
	}

	var response struct {
		ShortUrl string `json:"shortUrl"`
		Message  string `json:"message"`
	}

	response.ShortUrl = encodedUrl.ShortUrl
	response.Message = "short url created successfully"
	err = json.NewEncoder(w).Encode(response)
	if isError := utils.HandleHttpErrors(w, "Unable to fetch records", http.StatusInternalServerError, err); isError {
		return
	}
}

func DecodeUrl(w http.ResponseWriter, r *http.Request) {
	var encodedUrl models.URL
	err := json.NewDecoder(r.Body).Decode(&encodedUrl)
	if isError := utils.HandleHttpErrors(w, "", http.StatusInternalServerError, err); isError {
		return
	}

	collection, err := utils.GetCollection(utils.DatabaseClient, utils.ENCODED_URLS)
	if isError := utils.HandleHttpErrors(w, "could not fetch records", http.StatusInternalServerError, err); isError {
		return
	}

	var encodingUrl models.EncodeUrls
	filter := bson.M{"shortUrl": encodedUrl.UrlString}
	doc := collection.FindOne(context.TODO(), filter)
	err = doc.Decode(&encodingUrl)
	if isError := utils.HandleHttpErrors(w, "unable to fetch record", http.StatusInternalServerError, err); isError {
		return
	}

	var response struct {
		OriginalUrl string `json:"originalUrl"`
	}

	originalUrl, err := base64.StdEncoding.DecodeString(encodingUrl.EncodeUrl)
	if isError := utils.HandleHttpErrors(w, "", http.StatusInternalServerError, err); isError {
		return
	}
	response.OriginalUrl = string(originalUrl)
	err = json.NewEncoder(w).Encode(response)
	if isError := utils.HandleHttpErrors(w, "Unable to fetch records", http.StatusInternalServerError, err); isError {
		return
	}
}
