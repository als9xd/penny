package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/als9xd/penny/server/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type ThreadResponseData struct {
	Thread models.Thread
	Posts  []models.Post
}

type ThreadResponse struct {
	Error ApiError
	Data  ThreadResponseData
}

type ThreadsResponseData struct {
	Threads []models.Thread
}

type ThreadsResponse struct {
	Error ApiError
	Data  ThreadsResponseData
}

type ThreadDeleteResponseData struct {
	Thread models.Thread
}

type ThreadDeleteResponse struct {
	Error ApiError
	Data  ThreadDeleteResponseData
}

func testGetThread(t *testing.T, thread *models.Thread, router *gin.Engine) {
	expectedResponse := ThreadResponse{
		Data: ThreadResponseData{
			Thread: *thread,
			Posts:  []models.Post{},
		},
	}

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/t/%d", thread.Id), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	response := ThreadResponse{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)

	assert.Equal(t, expectedResponse, response)
}

func testGetMissingThread(t *testing.T, thread *models.Thread, router *gin.Engine) {

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/t/%d", thread.Id), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func testGetThreads(t *testing.T, threads *[]models.Thread, router *gin.Engine) {
	expectedResponse := ThreadsResponse{
		Data: ThreadsResponseData{
			Threads: *threads,
		},
	}

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/v1/t", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	response := ThreadsResponse{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)

	assert.Equal(t, expectedResponse, response)
}

func testCreateThread(t *testing.T, thread *models.Thread, router *gin.Engine) {
	expectedResponse := ThreadResponse{
		Data: ThreadResponseData{
			Thread: *thread,
		},
	}

	w := httptest.NewRecorder()

	data := url.Values{}
	data.Set("name", expectedResponse.Data.Thread.Name)
	data.Set("profile_id", strconv.Itoa(expectedResponse.Data.Thread.ProfileId))

	req, _ := http.NewRequest("POST", "/api/v1/t", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	response := ThreadResponse{}

	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)

	assert.Equal(t, expectedResponse, response)
}

func testUpdateThread(t *testing.T, thread *models.Thread, router *gin.Engine) {
	expectedResponse := ThreadResponse{
		Data: ThreadResponseData{
			Thread: *thread,
		},
	}

	w := httptest.NewRecorder()

	data := url.Values{}
	data.Set("name", expectedResponse.Data.Thread.Name)
	data.Set("profile_id", strconv.Itoa(expectedResponse.Data.Thread.ProfileId))

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/t/%d", expectedResponse.Data.Thread.Id), strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	response := ThreadResponse{}

	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)

	assert.Equal(t, expectedResponse, response)
}

func testDeleteThread(t *testing.T, thread *models.Thread, router *gin.Engine) {
	expectedResponse := ThreadDeleteResponse{
		Data: ThreadDeleteResponseData{
			Thread: *thread,
		},
	}

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/t/%d", expectedResponse.Data.Thread.Id), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	response := ThreadDeleteResponse{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)

	assert.Equal(t, expectedResponse, response)
}

func testMissingParamThread(t *testing.T, thread *models.Thread, router *gin.Engine) {
	expectedResponse := ThreadResponse{
		Data: ThreadResponseData{
			Thread: *thread,
		},
	}

	w := httptest.NewRecorder()

	data := url.Values{}
	// Missing Param
	// data.Set("name", expectedResponse.Data.Thread.Name)
	data.Set("profile_id", strconv.Itoa(expectedResponse.Data.Thread.ProfileId))

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/t/%d", expectedResponse.Data.Thread.Id), strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestThreadEndpoints(t *testing.T) {
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

	thread := models.Thread{
		Id:        1,
		Name:      "test_name",
		ProfileId: profile.Id,
	}

	testCreateThread(t, &thread, router)
	testGetThread(t, &thread, router)

	threads := []models.Thread{thread}
	testGetThreads(t, &threads, router)

	testMissingParamThread(t, &thread, router)

	thread.Name = "test_name_updated"
	testUpdateThread(t, &thread, router)
	testDeleteThread(t, &thread, router)

	testGetMissingThread(t, &thread, router)
}
