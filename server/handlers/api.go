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

/* Threads */

func GetThread(db *sqlx.DB,c *gin.Context) {
  thread := models.Thread{}
  err := db.Get(&thread,`
    SELECT thread.id,thread.name,thread.creator_id,array_agg(posts.id)
    FROM thread INNER JOIN posts
    WHERE thread.id=$1;
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

func GetThreads(db *sqlx.DB,c *gin.Context) {
  threads := []models.Thread{}
  err := db.Select(&threads,`
    SELECT thread.id,thread.name,thread.creator_id,array_agg(posts.id)
    FROM thread
    INNER JOIN posts;
  `)
  if err != nil {
    log.Fatal(err)
  }
  c.JSON(http.StatusOK,gin.H{"data":threads})
}

func CreateThread(db *sqlx.DB,c *gin.Context) {
  var b models.Thread
  err := c.ShouldBind(&b)
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
    RETURNING id,name,creator_id;
  `,c.PostForm("name"))
  if err != nil {
    log.Fatal(err)
  }
  c.JSON(http.StatusOK, gin.H{"data":thread})
}

func UpdateThread(db *sqlx.DB,c *gin.Context) {
  var b models.Thread
  err := c.ShouldBind(&b)
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
    UPDATE thread SET name = $1
    WHERE id = $2
    RETURNING id,name,creator_id;
  `,c.PostForm("name"),c.Param("id"))
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
}

func DeleteThread(db *sqlx.DB,c *gin.Context) {
  thread := models.Thread{}
  err := db.Get(&thread,`
    DELETE FROM thread WHERE id = $1 RETURNING id,name,creator_id;
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
}

/* Posts */

func GetPost(db *sqlx.DB,c *gin.Context) {
  post := models.Post{}
  err := db.Get(&post,`
    SELECT * from post WHERE id = $1;
  `,c.Param("id"))
  if err == sql.ErrNoRows {
    c.JSON(http.StatusNotFound,gin.H{
      "error": gin.H{
        "code": http.StatusNotFound,
        "message": fmt.Sprintf("No post with id '%s'",c.Param("id")),
      },
    })
    return
  }
  if err != nil {
    log.Fatal(err)
  }
  c.JSON(http.StatusOK, gin.H{"data":post})
}

func GetPosts(db *sqlx.DB,c *gin.Context) {
  posts := []models.Post{}
  err := db.Select(&posts,`
    SELECT id,author_id,text FROM post;
  `)
  if err != nil {
    log.Fatal(err)
  }
  c.JSON(http.StatusOK,gin.H{"data":posts})
}

func CreatePost(db *sqlx.DB,c *gin.Context) {
  var b models.Post
  err := c.ShouldBind(&b)
  if err != nil {
    c.JSON(http.StatusBadRequest,gin.H{
      "error": gin.H{
        "code": http.StatusBadRequest,
        "message": err.Error(),
      },
    })
    return
  }
  post := models.Post{}
  err = db.Get(&post,`
    INSERT INTO post (name)
    VALUES ($1)
    RETURNING id,author_id,text;
  `,c.PostForm("name"))
  if err != nil {
    log.Fatal(err)
  }
  c.JSON(http.StatusOK, gin.H{"data":post})
}
