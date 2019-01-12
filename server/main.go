package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/als9xd/penny/server/handlers"
	"github.com/gin-contrib/static"
)

var profileSchema = `
CREATE TABLE IF NOT EXISTS profile (
  id       SERIAL PRIMARY KEY,
  username VARCHAR(100) UNIQUE NOT NULL,
  email    VARCHAR(100),
  password TEXT NOT NULL,
  avatar   TEXT,
  created TIMESTAMPTZ NOT NULL
)`

var threadSchema = `
CREATE TABLE IF NOT EXISTS thread (
  id          SERIAL PRIMARY KEY,
  name        VARCHAR(100) NOT NULL,
  description TEXT,
  profile_id  INTEGER REFERENCES profile (id),
  avatar      TEXT,
  created TIMESTAMPTZ NOT NULL
)`

var threadSubscriptionSchema = `
CREATE TABLE IF NOT EXISTS thread_subscription (
	profile_id         INTEGER REFERENCES profile (id),
	thread_id      INTEGER REFERENCES thread (id),
	UNIQUE(profile_id,thread_id)
)`

var postSchema = `
CREATE TABLE IF NOT EXISTS post (
  id         SERIAL PRIMARY KEY,
  name    VARCHAR(100) NOT NULL,
  description TEXT,
  avatar TEXT,
  profile_id INTEGER REFERENCES profile (id),
  thread_id  INTEGER REFERENCES thread (id),
  created TIMESTAMPTZ NOT NULL,
  last_edited TIMESTAMPTZ
)`

var commentSchema = `
CREATE TABLE IF NOT EXISTS comment (
  id         SERIAL PRIMARY KEY,
  value    TEXT NOT NULL,
  level INTEGER NOT NUll,
  parent_comment_id INTEGER,
  profile_id INTEGER REFERENCES profile (id),
  post_id  INTEGER REFERENCES post (id),
  created TIMESTAMPTZ NOT NULL,
  last_edited TIMESTAMPTZ
)`

var profileSubscriptionSchema = `
CREATE TABLE IF NOT EXISTS profile_subscription (
  profile_id         INTEGER REFERENCES profile (id),
  to_profile_id      INTEGER REFERENCES profile (id),
  UNIQUE(profile_id,to_profile_id)
)`

func ConnectDb() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "host=localhost user=postgres password=password dbname=penny sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func CreateDb(db *sqlx.DB) {
	db.MustExec(profileSchema)
	db.MustExec(threadSchema)
	db.MustExec(threadSubscriptionSchema)
	db.MustExec(postSchema)
	db.MustExec(profileSubscriptionSchema)
	db.MustExec(commentSchema)
}

func ResetDb(db *sqlx.DB) {
	db.MustExec("DROP TABLE IF EXISTS profile CASCADE")
	db.MustExec("DROP TABLE IF EXISTS thread CASCADE")
	db.MustExec("DROP TABLE IF EXISTS thread_subscription CASCADE")
	db.MustExec("DROP TABLE IF EXISTS post CASCADE")
	db.MustExec("DROP TABLE IF EXISTS profile_subscription CASCADE")
	db.MustExec("DROP TABLE IF EXISTS comment CASCADE")
}

func GetRouter(db *sqlx.DB) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	corsConfig := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		AllowHeaders:     []string{"Authorization", "Origin", "Access-Control-Allow-Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
	router.Use(corsConfig)

	router.Use(static.Serve("/", static.LocalFile("../client/dist", false)))
	router.Use(static.Serve("/uploads", static.LocalFile("./uploads", false)))

	apiv1 := router.Group("/api/v1")
	apiv1Protected := router.Group("/api/v1/protected")

	authMiddleware := handlers.AuthMiddleware(db)
	apiv1.POST("/login", authMiddleware.LoginHandler)

	apiv1Protected.Use(authMiddleware.MiddlewareFunc())
	{
		apiv1Protected.GET("/u", func(c *gin.Context) { handlers.GetPrivateProfile(db, c) })
		apiv1Protected.POST("/subscribe/u/:profile_id", func(c *gin.Context) { handlers.SubscribeProfile(db, c) })
		apiv1Protected.POST("/subscribe/t/:thread_id", func(c *gin.Context) { handlers.SubscribeThread(db, c) })
		apiv1Protected.POST("/t", func(c *gin.Context) { handlers.CreateThread(db, c) })
		apiv1Protected.POST("/p", func(c *gin.Context) { handlers.CreatePost(db, c) })
		apiv1Protected.POST("/c", func(c *gin.Context) { handlers.CreateComment(db, c) })
	}

	apiv1.GET("/t/:id", func(c *gin.Context) { handlers.GetThread(db, c) })
	apiv1.GET("/t", func(c *gin.Context) { handlers.GetThreads(db, c) })

	apiv1.GET("/p/:id", func(c *gin.Context) { handlers.GetPost(db, c) })
	apiv1.GET("/p", func(c *gin.Context) { handlers.GetPosts(db, c) })

	apiv1.GET("/u/:id", func(c *gin.Context) { handlers.GetProfile(db, c) })
	apiv1.GET("/u", func(c *gin.Context) { handlers.GetProfiles(db, c) })
	apiv1.POST("/u", func(c *gin.Context) { handlers.CreateProfile(db, c) })

	apiv1.GET("/search", func(c *gin.Context) { handlers.Search(db, c) })

	return router
}

func main() {
	db := ConnectDb()
	// ResetDb(db)
	CreateDb(db)

	GetRouter(db).Run(":3000")
}
