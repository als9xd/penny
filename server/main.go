package main

import (
  "log"

  "github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
  "github.com/jmoiron/sqlx"

  "github.com/als9xd/penny/server/handlers"
)

var userSchema = `
CREATE TABLE IF NOT EXISTS profile (
  id       SERIAL PRIMARY KEY,
  username VARCHAR(100) NOT NULL,
  email    VARCHAR(100)
)`

var threadSchema = `
CREATE TABLE IF NOT EXISTS thread (
  id         SERIAL PRIMARY KEY,
  name       VARCHAR(100) NOT NULL,
  profile_id NUMERIC REFERENCES profile (id)
)`

var postSchema = `
CREATE TABLE IF NOT EXISTS post (
  id         SERIAL PRIMARY KEY,
  comment    TEXT,
  profile_id NUMERIC REFERENCES profile (id),
  thread_id  NUMERIC REFERENCES thread (id)
)`


func main(){
  db, err := sqlx.Connect("postgres", "host=db user=docker password=docker dbname=penny sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

  db.MustExec(userSchema)
  db.MustExec(threadSchema)
  db.MustExec(postSchema)

  router := gin.Default()

  apiv1 := router.Group("/api/v1")

  apiv1.GET("/t/:id", func(c *gin.Context) {handlers.GetThread(db,c)})
  apiv1.GET("/t", func(c *gin.Context) {handlers.GetThreads(db,c)})
  apiv1.POST("/t", func(c *gin.Context) {handlers.CreateThread(db,c)})
  apiv1.PUT("/t/:id", func(c *gin.Context) {handlers.UpdateThread(db,c)})
  apiv1.DELETE("/t/:id", func(c *gin.Context) {handlers.DeleteThread(db,c)})

  apiv1.GET("/p/:id", func(c *gin.Context) {handlers.GetPost(db,c)})
  apiv1.GET("/p", func(c *gin.Context) {handlers.GetPosts(db,c)})
  apiv1.POST("/p", func(c *gin.Context) {handlers.CreatePost(db,c)})
  apiv1.PUT("/p/:id", func(c *gin.Context) {handlers.UpdatePost(db,c)})
  apiv1.DELETE("/p/:id", func(c *gin.Context) {handlers.DeletePost(db,c)})

  apiv1.GET("/u/:id", func(c *gin.Context) {handlers.GetProfile(db,c)})
  apiv1.GET("/u", func(c *gin.Context) {handlers.GetProfiles(db,c)})
  apiv1.POST("/u", func(c *gin.Context) {handlers.CreateProfile(db,c)})
  apiv1.PUT("/u/:id", func(c *gin.Context) {handlers.UpdateProfile(db,c)})
  apiv1.DELETE("/u/:id", func(c *gin.Context) {handlers.DeleteProfile(db,c)})

  router.Run(":3000")
}
