package main

import (
  "log"

  "github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
  "github.com/jmoiron/sqlx"

  "github.com/als9xd/penny/server/handlers"
)

var schema = `
CREATE TABLE IF NOT EXISTS user (
  id       SERIAL PRIMARY KEY,
  username VARCHAR(100) NOT NULL,
  email    VARCHAR(100) NOT NULL
)

CREATE TABLE IF NOT EXISTS post (
  id        SERIAL PRIMARY KEY,
  comment   TEXT NOT NULL,
  user_id   INTEGER REFERENCES user (id),
  thread_id INTEGER REFERENCES thread (id)
)

CREATE TABLE IF NOT EXISTS thread (
  id      SERIAL PRIMARY KEY,
  name    VARCHAR(100) NOT NULL,
  user_id INTEGER REFERENCES user (id)
)
`

func main(){
  db, err := sqlx.Connect("postgres", "host=db user=docker password=docker dbname=penny sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

  db.MustExec(schema)

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

  apiv1.GET("/u/:id", func(c *gin.Context) {handlers.GetUser(db,c)})
  apiv1.GET("/u", func(c *gin.Context) {handlers.GetUsers(db,c)})
  apiv1.POST("/u", func(c *gin.Context) {handlers.CreateUser(db,c)})
  apiv1.PUT("/u/:id", func(c *gin.Context) {handlers.UpdateUser(db,c)})
  apiv1.DELETE("/u/:id", func(c *gin.Context) {handlers.DeleteUser(db,c)})

  router.Run(":3000")
}
