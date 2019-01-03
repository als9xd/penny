package main

import (
  "log"

  "github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
  "github.com/jmoiron/sqlx"

  "github.com/als9xd/penny/server/handlers"
)

var schema = `
CREATE TABLE IF NOT EXISTS thread (
  id   SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL
)`

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

  router.Run(":3000")
}
