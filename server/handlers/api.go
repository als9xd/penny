package handlers

import (
  "fmt"
  "log"

  "github.com/gin-gonic/gin"
  "net/http"

  "database/sql"
  _ "github.com/lib/pq"
  "github.com/jmoiron/sqlx"

  "github.com/als9xd/penny/server/models"
)

func GetThread(db *sqlx.DB,c *gin.Context) {
  thread := models.Thread{}
  err := db.Get(&thread,`
    SELECT id,name FROM thread WHERE id=$1;
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
  c.JSON(http.StatusOK,thread)
}
