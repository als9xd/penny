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

type PostResponseData struct {
	Post  models.Post
	Posts []models.Post
}

type PostResponse struct {
	Error ApiError
	Data  PostResponseData
}

type PostsResponseData struct {
	Posts []models.Post
}

type PostsResponse struct {
	Error ApiError
	Data  PostsResponseData
}

type PostDeleteResponseData struct {
	Post models.Post
}

type PostDeleteResponse struct {
	Error ApiError
	Data  PostDeleteResponseData
}

func testGetPost(t *testing.T, post *models.Post, router *gin.Engine) {
	expectedResponse := PostResponse{
		Data: PostResponseData{
			Post: *post,
		},
	}

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/p/%d", post.Id), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	response := PostResponse{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)

	assert.Equal(t, expectedResponse, response)
}

func testGetMissingPost(t *testing.T, post *models.Post, router *gin.Engine) {

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/p/%d", post.Id), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func testGetPosts(t *testing.T, posts *[]models.Post, router *gin.Engine) {
	expectedResponse := PostsResponse{
		Data: PostsResponseData{
			Posts: *posts,
		},
	}

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/v1/p", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	response := PostsResponse{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)

	assert.Equal(t, expectedResponse, response)
}

func testCreatePost(t *testing.T, post *models.Post, router *gin.Engine) {
	expectedResponse := PostResponse{
		Data: PostResponseData{
			Post: *post,
		},
	}

	w := httptest.NewRecorder()

	data := url.Values{}
	data.Set("comment", expectedResponse.Data.Post.Comment)
	data.Set("profile_id", strconv.Itoa(expectedResponse.Data.Post.ProfileId))
	data.Set("thread_id", strconv.Itoa(expectedResponse.Data.Post.ProfileId))

	req, _ := http.NewRequest("POST", "/api/v1/p", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	response := PostResponse{}

	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)

	assert.Equal(t, expectedResponse, response)
}

func testUpdatePost(t *testing.T, post *models.Post, router *gin.Engine) {
	expectedResponse := PostResponse{
		Data: PostResponseData{
			Post: *post,
		},
	}

	w := httptest.NewRecorder()

	data := url.Values{}
	data.Set("comment", expectedResponse.Data.Post.Comment)
	data.Set("profile_id", strconv.Itoa(expectedResponse.Data.Post.ProfileId))
	data.Set("thread_id", strconv.Itoa(expectedResponse.Data.Post.ProfileId))

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/p/%d", expectedResponse.Data.Post.Id), strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	response := PostResponse{}

	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)

	assert.Equal(t, expectedResponse, response)
}

func testDeletePost(t *testing.T, post *models.Post, router *gin.Engine) {
	expectedResponse := PostDeleteResponse{
		Data: PostDeleteResponseData{
			Post: *post,
		},
	}

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/p/%d", expectedResponse.Data.Post.Id), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	response := PostDeleteResponse{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)

	assert.Equal(t, expectedResponse, response)
}

func testMissingParamPost(t *testing.T, post *models.Post, router *gin.Engine) {
	expectedResponse := PostResponse{
		Data: PostResponseData{
			Post: *post,
		},
	}

	w := httptest.NewRecorder()

	data := url.Values{}
	// Missing Param
	// data.Set("comment", expectedResponse.Data.Post.Comment)
	data.Set("profile_id", strconv.Itoa(expectedResponse.Data.Post.ProfileId))
	data.Set("thread_id", strconv.Itoa(expectedResponse.Data.Post.ProfileId))

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/p/%d", expectedResponse.Data.Post.Id), strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostEndpoints(t *testing.T) {
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

	post := models.Post{
		Id:        1,
		Comment:   "test_comment",
		ProfileId: profile.Id,
		ThreadId:  thread.Id,
	}

	testCreatePost(t, &post, router)
	testGetPost(t, &post, router)

	posts := []models.Post{post}
	testGetPosts(t, &posts, router)

	testMissingParamPost(t, &post, router)

	post.Comment = "test_comment_updated"
	testUpdatePost(t, &post, router)
	testDeletePost(t, &post, router)

	testGetMissingPost(t, &post, router)
}
