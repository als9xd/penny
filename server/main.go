package main

import (
  "fmt"
  "log"

  "github.com/gin-gonic/gin"
  "net/http"

  "database/sql"
	_ "github.com/lib/pq"
  "github.com/jmoiron/sqlx"

  "github.com/als9xd/penny/server/handlers"
  "github.com/als9xd/penny/server/models"
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

  apiv1.GET("/t", func(c *gin.Context) {
    threads := []models.Thread{}
    err = db.Select(&threads,`
      SELECT id,name FROM thread;
    `)
    if err != nil {
      log.Fatal(err)
    }
		c.JSON(http.StatusOK,gin.H{"data":threads})
	})

  apiv1.GET("/t/:id", func(c *gin.Context) {
    handlers.GetThread(db,c)
  })

  apiv1.POST("/t", func(c *gin.Context) {
    var b models.Thread
    err = c.ShouldBind(&b)
    if err != nil {
      c.JSON(http.StatusBadRequest,gin.H{
        "error": gin.H{
          "code": http.StatusBadRequest,
          "message": err.Error(),
        },
      })
      return
    }
    thread := models.Thread{}
    err = db.Get(&thread,`
      INSERT INTO thread (name)
      VALUES ($1)
      RETURNING id,name;
    `,c.PostForm("name"))
    if err != nil {
      log.Fatal(err)
    }
    c.JSON(http.StatusOK, gin.H{"data":thread})
  })

  apiv1.DELETE("/t/:id", func(c *gin.Context) {
    thread := models.Thread{}
    err = db.Get(&thread,`
      DELETE FROM thread WHERE id = $1 RETURNING id,name;
    `,c.Param("id"))
    if err == sql.ErrNoRows {
      c.JSON(http.StatusNotFound,gin.H{
        "error": gin.H{
          "code": http.StatusNotFound,
          "message": fmt.Sprintf("No thread with id '%s'",c.Param("id")),
        },
      })
      return
    }
    if err != nil {
      log.Fatal(err)
    }
    c.JSON(http.StatusOK, gin.H{"data":thread})
  })


  router.Run(":3000")
}
