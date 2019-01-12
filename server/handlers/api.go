package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin/binding"

	"net/http"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"

	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/als9xd/penny/server/models"
	uuid "github.com/satori/go.uuid"
)

/* Threads */

func GetThread(db *sqlx.DB, c *gin.Context) {
	thread := models.Thread{}
	err := db.Get(&thread, `
    SELECT * FROM thread WHERE id = $1;
  `, c.Param("id"))
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    http.StatusNotFound,
				"message": fmt.Sprintf("No thread with id '%s'", c.Param("id")),
			},
		})
		return
	}
	if err != nil {
		log.Fatal(err)
	}

	posts := []models.Post{}
	err = db.Select(&posts, `
    SELECT post.*,profile.avatar as profile_avatar,profile.username as profile_username FROM post INNER JOIN profile ON post.profile_id = profile.id WHERE thread_id = $1;
  `, c.Param("id"))
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"thread": thread,
		"posts":  posts,
	}})
}

func GetThreads(db *sqlx.DB, c *gin.Context) {
	threads := []models.Thread{}
	err := db.Select(&threads, `
    SELECT thread.*,COUNT(post.thread_id) as post_count FROM thread LEFT JOIN post ON thread.id = post.thread_id GROUP BY thread.id;
  `)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"threads": threads,
	}})
}

func CreateThread(db *sqlx.DB, c *gin.Context) {
	claims := jwt.ExtractClaims(c)

	var newThread models.Thread
	err := c.MustBindWith(&newThread, binding.FormMultipart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			},
		})
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil && err.Error() != "http: no such file" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			},
		})
		return
	}

	var fileUUID string
	if file != nil {
		fileUUID = uuid.Must(uuid.NewV4()).String()
		filepath := fmt.Sprintf("./uploads/%s", fileUUID)
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    http.StatusBadRequest,
					"message": fmt.Sprintf("File upload error: %v", err.Error()),
				},
			})
			return
		}
	}

	thread := models.Thread{}
	err = db.Get(&thread, `
    INSERT INTO thread (name,description,profile_id,avatar,created)
    VALUES ($1,$2,$3,$4,$5)
    RETURNING thread.*;
  `, newThread.Name, newThread.Description, claims["id"], fileUUID, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"thread": thread}})
}

func UpdateThread(db *sqlx.DB, c *gin.Context) {
	var b models.Thread
	err := c.ShouldBind(&b)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			},
		})
		return
	}
	thread := models.Thread{}
	err = db.Get(&thread, `
    UPDATE thread SET name = $1
    WHERE id = $2
    RETURNING id,name,profile_id;
  `, c.PostForm("name"), c.Param("id"))
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    http.StatusNotFound,
				"message": fmt.Sprintf("No thread with id '%s'", c.Param("id")),
			},
		})
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"thread": thread}})
}

func DeleteThread(db *sqlx.DB, c *gin.Context) {
	thread := models.Thread{}
	err := db.Get(&thread, `
    DELETE FROM thread
    WHERE thread.id = $1 RETURNING thread.*;
  `, c.Param("id"))
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    http.StatusNotFound,
				"message": fmt.Sprintf("No thread with id '%s'", c.Param("id")),
			},
		})
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"thread": thread}})
}

/* Posts */

func GetPost(db *sqlx.DB, c *gin.Context) {
	post := models.Post{}
	err := db.Get(&post, `SELECT post.*,profile.avatar as profile_avatar,profile.username as profile_username,thread.id as thread_id,thread.name as thread_name,thread.avatar as thread_avatar FROM post INNER JOIN profile ON post.profile_id = profile.id INNER JOIN thread ON post.thread_id = thread.id WHERE post.id =  $1`, c.Param("id"))
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    http.StatusNotFound,
				"message": fmt.Sprintf("No post with id '%s'", c.Param("id")),
			},
		})
		return
	}
	if err != nil {
		log.Fatal(err)
	}

	comments := []models.Comment{}
	err = db.Select(&comments, `SELECT comment.*,profile.id as profile_id,profile.username as profile_username,profile.avatar as profile_avatar FROM comment INNER JOIN profile ON profile.id = comment.profile_id WHERE comment.post_id = $1`, c.Param("id"))
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"post": post, "comments": comments}})
}

func GetPosts(db *sqlx.DB, c *gin.Context) {
	posts := []models.Post{}
	err := db.Select(&posts, `
    SELECT * FROM post;
  `)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"posts": posts}})
}

func CreatePost(db *sqlx.DB, c *gin.Context) {
	claims := jwt.ExtractClaims(c)

	var newPost models.Post
	err := c.MustBindWith(&newPost, binding.FormMultipart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			},
		})
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil && err.Error() != "http: no such file" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			},
		})
		return
	}

	var fileUUID string
	if file != nil {
		fileUUID = uuid.Must(uuid.NewV4()).String()
		filepath := fmt.Sprintf("./uploads/%s", fileUUID)
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    http.StatusBadRequest,
					"message": fmt.Sprintf("File upload error: %v", err.Error()),
				},
			})
			return
		}
	}

	post := models.Post{}
	err = db.Get(&post, `
    INSERT INTO post (name,description,avatar,profile_id,thread_id,created)
    VALUES ($1,$2,$3,$4,$5,$6)
    RETURNING post.*;
  `, newPost.Name, newPost.Description, fileUUID, claims["id"], newPost.ThreadId, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"post": post}})
}

func UpdatePost(db *sqlx.DB, c *gin.Context) {
	var b models.Post
	err := c.ShouldBind(&b)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			},
		})
		return
	}
	post := models.Post{}
	err = db.Get(&post, `
    UPDATE post
    SET comment  = $1,
        profile_id  = $2,
        thread_id = $3
    WHERE id = $4
    RETURNING id,comment,profile_id,thread_id;
  `, c.PostForm("comment"), c.PostForm("profile_id"), c.PostForm("thread_id"), c.Param("id"))
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    http.StatusNotFound,
				"message": fmt.Sprintf("No post with id '%s'", c.Param("id")),
			},
		})
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"post": post}})
}

func DeletePost(db *sqlx.DB, c *gin.Context) {
	post := models.Post{}
	err := db.Get(&post, `
    DELETE FROM post WHERE id = $1 RETURNING id,comment,profile_id,thread_id;
  `, c.Param("id"))
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    http.StatusNotFound,
				"message": fmt.Sprintf("No post with id '%s'", c.Param("id")),
			},
		})
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"post": post}})
}

/* Profiles */

func GetProfile(db *sqlx.DB, c *gin.Context) {
	profile := models.ProfilePublic{}
	err := db.Get(&profile, `
    SELECT id,username,avatar,created FROM profile
    WHERE id = $1;
  `, c.Param("id"))
	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    http.StatusBadRequest,
				"message": fmt.Sprintf("No profile with id '%s'", c.Param("id")),
			},
		})
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"profile": profile}})
}

func GetPrivateProfile(db *sqlx.DB, c *gin.Context) {
	claims := jwt.ExtractClaims(c)

	profile := models.ProfilePrivate{}
	err := db.Get(&profile, `
    SELECT id,username,email,avatar,created FROM profile
    WHERE id = $1;
  `, claims["id"])
	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    http.StatusBadRequest,
				"message": fmt.Sprintf("No profile with id '%d'", claims["id"]),
			},
		})
		return
	}
	if err != nil {
		log.Fatal(err)
	}

	profileSubscriptions := []models.ProfilePublic{}
	err = db.Select(&profileSubscriptions, `
	SELECT profile.id,profile.username,profile.avatar,profile.created FROM profile_subscription
	INNER JOIN profile ON profile_subscription.to_profile_id = profile.id
	WHERE profile_subscription.profile_id = $1;
	`, claims["id"])
	if err != nil {
		log.Fatal(err)
	}

	threadSubscriptions := []models.Thread{}
	err = db.Select(&threadSubscriptions, `
	SELECT thread.* FROM thread_subscription
	INNER JOIN thread ON thread_subscription.thread_id = thread.id
	WHERE thread_subscription.profile_id = $1;
	`, claims["id"])
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"profile": profile, "subscriptions": gin.H{"profiles": profileSubscriptions, "threads": threadSubscriptions}}})
}

func GetProfiles(db *sqlx.DB, c *gin.Context) {
	profiles := []models.ProfilePublic{}
	err := db.Select(&profiles, `
    SELECT id,username,avatar,creataed FROM profile;
  `)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"profiles": profiles}})
}

func CreateProfile(db *sqlx.DB, c *gin.Context) {
	var newProfile models.Profile
	err := c.MustBindWith(&newProfile, binding.FormMultipart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			},
		})
		return
	}

	var usernameTaken bool
	err = db.Get(&usernameTaken, "SELECT 1 FROM profile WHERE username = $1;", newProfile.Username)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}
	if usernameTaken {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    http.StatusBadRequest,
				"message": "Username taken",
			},
		})
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil && err.Error() != "http: no such file" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			},
		})
		return
	}

	var fileUUID string
	if file != nil {
		fileUUID = uuid.Must(uuid.NewV4()).String()
		filepath := fmt.Sprintf("./uploads/%s", fileUUID)
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    http.StatusBadRequest,
					"message": fmt.Sprintf("File upload error: %v", err.Error()),
				},
			})
			return
		}
	}

	profile := models.ProfilePrivate{}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newProfile.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Get(&profile, `
    INSERT INTO profile (username,email,password,avatar,created)
    VALUES ($1,$2,$3,$4,$5)
    RETURNING id,username,email,avatar,created;
  `, newProfile.Username, newProfile.Email, passwordHash, fileUUID, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"profile": profile}})
}

func UpdateProfile(db *sqlx.DB, c *gin.Context) {
	var b models.ProfilePrivate
	err := c.ShouldBind(&b)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			},
		})
		return
	}
	profile := models.ProfilePrivate{}
	err = db.Get(&profile, `
    UPDATE profile
    SET username = $1,
        email = $2
    WHERE id = $3
    RETURNING id,username,email,avatar;
  `, c.PostForm("username"), c.PostForm("email"), c.Param("id"))
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    http.StatusNotFound,
				"message": fmt.Sprintf("No profile with id '%s'", c.Param("id")),
			},
		})
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"profile": profile}})
}

func DeleteProfile(db *sqlx.DB, c *gin.Context) {
	profile := models.ProfilePrivate{}
	err := db.Get(&profile, `
    DELETE FROM profile WHERE id = $1 RETURNING id,username,email,avatar;
  `, c.Param("id"))
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    http.StatusNotFound,
				"message": fmt.Sprintf("No profile with id '%s'", c.Param("id")),
			},
		})
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"profile": profile}})
}

func Search(db *sqlx.DB, c *gin.Context) {
	profiles := []models.ProfilePublic{}
	err := db.Select(&profiles, `SELECT id,username,avatar,created FROM profile WHERE username ILIKE $1`, fmt.Sprintf("%%%s%%", c.Query("search")))
	if err != nil {
		log.Fatal(err)
	}
	threads := []models.Thread{}
	err = db.Select(&threads, `SELECT * FROM thread WHERE name ILIKE $1`, fmt.Sprintf("%%%s%%", c.Query("search")))
	if err != nil {
		log.Fatal(err)
	}
	posts := []models.Post{}
	err = db.Select(&posts, `SELECT post.*,profile.avatar as profile_avatar,profile.username as profile_username,thread.id as thread_id,thread.name as thread_name,thread.avatar as thread_avatar FROM post INNER JOIN profile ON post.profile_id = profile.id INNER JOIN thread ON post.thread_id = thread.id WHERE post.name ILIKE $1`, fmt.Sprintf("%%%s%%", c.Query("search")))
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"profiles": profiles, "threads": threads, "posts": posts}})
}

type SubscribeRequest struct {
	ProfileId int `uri:"profile_id"  binding:"required"`
}

func SubscribeProfile(db *sqlx.DB, c *gin.Context) {
	var subscribeRequest SubscribeRequest
	if err := c.ShouldBindUri(&subscribeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    http.StatusBadRequest,
				"message": "Profile id should be an integer",
			},
		})
		return
	}

	claims := jwt.ExtractClaims(c)

	if claims["id"] == float64(subscribeRequest.ProfileId) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    http.StatusBadRequest,
				"message": "Cannot subscribe to yourself",
			},
		})
		return
	}

	var alreadySubscribed bool
	err := db.Get(&alreadySubscribed, "SELECT 1 FROM profile_subscription WHERE profile_id = $1 AND to_profile_id = $2;", claims["id"], subscribeRequest.ProfileId)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}

	if alreadySubscribed == true {
		db.MustExec(`DELETE FROM profile_subscription WHERE profile_id = $1 AND to_profile_id = $2;`, claims["id"], subscribeRequest.ProfileId)
		c.JSON(http.StatusOK, gin.H{"data": gin.H{"subscribed": false, "subscription": gin.H{"profile_id": claims["id"], "to_profile_id": subscribeRequest.ProfileId}}})
		return
	}

	db.MustExec(`INSERT INTO profile_subscription (profile_id,to_profile_id) VALUES ($1,$2);`, claims["id"], subscribeRequest.ProfileId)
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"subscribed": true, "subscription": gin.H{"profile_id": claims["id"], "to_profile_id": subscribeRequest.ProfileId}}})
}

type ThreadRequest struct {
	ThreadId int `uri:"thread_id"  binding:"required"`
}

func SubscribeThread(db *sqlx.DB, c *gin.Context) {
	var threadRequest ThreadRequest
	if err := c.ShouldBindUri(&threadRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    http.StatusBadRequest,
				"message": "Thread id should be an integer",
			},
		})
		return
	}

	claims := jwt.ExtractClaims(c)

	var alreadySubscribed bool
	err := db.Get(&alreadySubscribed, "SELECT 1 FROM thread_subscription WHERE profile_id = $1 AND thread_id = $2;", claims["id"], threadRequest.ThreadId)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}

	if alreadySubscribed == true {
		db.MustExec(`DELETE FROM thread_subscription WHERE profile_id = $1 AND thread_id = $2;`, claims["id"], threadRequest.ThreadId)
		c.JSON(http.StatusOK, gin.H{"data": gin.H{"subscribed": false, "subscription": gin.H{"profile_id": claims["id"], "thread_id": threadRequest.ThreadId}}})
		return
	}

	db.MustExec(`INSERT INTO thread_subscription (profile_id,thread_id) VALUES ($1,$2);`, claims["id"], threadRequest.ThreadId)
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"subscribed": true, "subscription": gin.H{"profile_id": claims["id"], "thread_id": threadRequest.ThreadId}}})
}

func CreateComment(db *sqlx.DB, c *gin.Context) {
	newComment := models.Comment{}
	if err := c.Bind(&newComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			},
		})
		return
	}

	claims := jwt.ExtractClaims(c)

	var parentLevel int
	err := db.Get(&parentLevel, "SELECT level FROM comment WHERE id = $1", newComment.ParentCommentId)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}

	comment := models.Comment{}
	err = db.Get(&comment, "INSERT INTO comment (value,level,parent_comment_id,profile_id,post_id,created) VALUES ($1,$2,$3,$4,$5,$6) RETURNING comment.*;", newComment.Value, parentLevel, newComment.ParentCommentId, claims["id"], newComment.PostId, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	err = db.Get(&comment, "SELECT comment.*,profile.id as profile_id,profile.avatar as profile_avatar,profile.username as profile_username FROM comment INNER JOIN profile ON comment.profile_id = profile.id WHERE comment.id = $1;", comment.Id)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"comment": comment}})
}
