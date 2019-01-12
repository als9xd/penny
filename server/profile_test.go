package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/als9xd/penny/server/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type ProfileResponseDataPrivate struct {
	Profile models.ProfilePrivate
}

type ProfileResponsePrivate struct {
	Error ApiError
	Data  ProfileResponseDataPrivate
}

type ProfilesResponseDataPrivate struct {
	Profiles []models.ProfilePrivate
}
type ProfilesResponsePrivate struct {
	Error ApiError
	Data  ProfilesResponseDataPrivate
}

type ProfileResponseDataPublic struct {
	Profile models.ProfilePublic
}
type ProfileResponsePublic struct {
	Error ApiError
	Data  ProfileResponseDataPublic
}

type ProfilesResponseDataPublic struct {
	Profiles []models.ProfilePublic
}
type ProfilesResponsePublic struct {
	Error ApiError
	Data  ProfilesResponseDataPublic
}

func testCreateProfile(t *testing.T, profile *models.Profile, router *gin.Engine) {
	expectedResponse := ProfileResponsePrivate{
		Data: ProfileResponseDataPrivate{
			Profile: *profile.ToPrivate(),
		},
	}

	w := httptest.NewRecorder()

	data := url.Values{}
	data.Set("username", "test_username")
	data.Set("password", profile.Password)
	data.Set("email", "test@test.com")

	req, _ := http.NewRequest("POST", "/api/v1/u", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	response := ProfileResponsePrivate{}

	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, response)
}

func testGetProfile(t *testing.T, profile *models.Profile, router *gin.Engine) {
	expectedResponse := ProfileResponsePublic{
		Data: ProfileResponseDataPublic{
			Profile: *profile.ToPublic(),
		},
	}

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/u/%d", profile.Id), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	response := ProfileResponsePublic{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)

	assert.Equal(t, expectedResponse, response)
}

func testGetMissingProfile(t *testing.T, profile *models.Profile, router *gin.Engine) {

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/u/%d", profile.Id), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func testGetProfiles(t *testing.T, profile *models.Profile, router *gin.Engine) {
	expectedResponse := ProfilesResponsePublic{
		Data: ProfilesResponseDataPublic{
			Profiles: []models.ProfilePublic{*profile.ToPublic()},
		},
	}

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/v1/u", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	response := ProfilesResponsePublic{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)

	assert.Equal(t, expectedResponse, response)
}

func testUpdateProfile(t *testing.T, profile *models.Profile, router *gin.Engine) {
	expectedResponse := ProfileResponsePrivate{
		Data: ProfileResponseDataPrivate{
			Profile: *profile.ToPrivate(),
		},
	}

	w := httptest.NewRecorder()

	data := url.Values{}
	data.Set("username", profile.Username)
	data.Set("password", profile.Password)
	data.Set("email", profile.Email)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/u/%d", profile.Id), strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	response := ProfileResponsePrivate{}

	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)

	assert.Equal(t, expectedResponse, response)
}

type ProfileDeleteResponseData struct {
	Profile models.ProfilePrivate
}
type ProfileDeleteResponse struct {
	Error ApiError
	Data  ProfileDeleteResponseData
}

func testDeleteProfile(t *testing.T, profile *models.Profile, router *gin.Engine) {
	expectedResponse := ProfileDeleteResponse{
		Data: ProfileDeleteResponseData{
			Profile: *profile.ToPrivate(),
		},
	}

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/u/%d", profile.Id), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	response := ProfileDeleteResponse{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)

	assert.Equal(t, expectedResponse, response)
}

func testMissingParamProfile(t *testing.T, profile *models.Profile, router *gin.Engine) {

	w := httptest.NewRecorder()

	data := url.Values{}
	// Missing Param
	// data.Set("username", expectedResponse.Data.Profile.Username)
	data.Set("password", profile.Password)
	data.Set("email", profile.Email)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/t/%d", profile.Id), strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProfileEndpoints(t *testing.T) {
	db := ConnectDb()
	ResetDb(db)
	defer ResetDb(db)
	CreateDb(db)

	router := GetRouter(db)

	profile := models.Profile{
		Id:       1,
		Username: "test_username",
		Email:    "test@test.com",
		Password: "test_password",
	}

	testCreateProfile(t, &profile, router)
	testGetProfile(t, &profile, router)

	testGetProfiles(t, &profile, router)

	testMissingParamProfile(t, &profile, router)

	profile.Username = "test_username_updated"
	testUpdateProfile(t, &profile, router)
	testDeleteProfile(t, &profile, router)

	testGetMissingProfile(t, &profile, router)
}
