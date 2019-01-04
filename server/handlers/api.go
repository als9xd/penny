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
    SELECT thread.id,thread.name,thread.creator_id,array_agg(post.id)
    FROM thread INNER JOIN post ON thread.id = post.thread_id
    WHERE thread.id = $1;
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
    SELECT thread.id,thread.name,thread.creator_id,array_agg(post.id)
    FROM thread
    INNER JOIN post ON thread.id = post.thread_id;
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
    INSERT INTO thread (name,creator_id)
    VALUES ($1)
    RETURNING id,name,creator_id;
  `,c.PostForm("name"),c.PostForm("creator_id"))
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
    DELETE FROM thread
    WHERE thread.id = $1 RETURNING thread.id,thread.name,thread.profile_id;
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
    SELECT id,comment,profile_id,thread_id from post WHERE id = $1;
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
    SELECT id,comment,profile_id,thread_id FROM post;
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
    INSERT INTO post (comment,profile_id,thread_id)
    VALUES ($1,$2,$3)
    RETURNING id,comment,profile_id,thread_id;
  `,c.PostForm("comment"),c.PostForm("profile_id"),c.PostForm("thread_id"))
  if err != nil {
    log.Fatal(err)
  }
  c.JSON(http.StatusOK, gin.H{"data":post})
}

func UpdatePost(db *sqlx.DB,c *gin.Context) {
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
    UPDATE post
    SET comment  = $1,
        profile_id  = $2,
        thread_id = $3
    WHERE id = $3
    RETURNING id,comment,profile_id,thread_id;
  `,c.PostForm("comment"),c.PostForm("profile_id"),c.Param("id"))
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

func DeletePost(db *sqlx.DB,c *gin.Context) {
  post := models.Post{}
  err := db.Get(&post,`
    DELETE FROM post WHERE id = $1 RETURNING id,comment,profile_id,thread_id;
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

/* Profiles */

func GetProfile(db *sqlx.DB,c *gin.Context) {
  profile := models.Profile{}
  err := db.Get(&profile,`
    SELECT id,username,email
    WHERE id = $1;
  `,c.Param("id"))
  if err == sql.ErrNoRows {
    c.JSON(http.StatusNotFound,gin.H{
      "error": gin.H{
        "code": http.StatusNotFound,
        "message": fmt.Sprintf("No profile with id '%s'",c.Param("id")),
      },
    })
    return
  }
  if err != nil {
    log.Fatal(err)
  }
  c.JSON(http.StatusOK,profile)
}

func GetProfiles(db *sqlx.DB,c *gin.Context) {
  profiles := []models.Profile{}
  err := db.Select(&profiles,`
    SELECT id,username,email FROM profile;
  `)
  if err != nil {
    log.Fatal(err)
  }
  c.JSON(http.StatusOK,gin.H{"data":profiles})
}

func CreateProfile(db *sqlx.DB,c *gin.Context) {
  var b models.Profile
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
  profile := models.Profile{}
  err = db.Get(&profile,`
    INSERT INTO profile (username,email)
    VALUES ($1,$2)
    RETURNING id,username,email;
  `,c.PostForm("username"),c.PostForm("email"))
  if err != nil {
    log.Fatal(err)
  }
  c.JSON(http.StatusOK, gin.H{"data":profile})
}

func UpdateProfile(db *sqlx.DB,c *gin.Context) {
  var b models.Profile
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
  profile := models.Profile{}
  err = db.Get(&profile,`
    UPDATE profile
    SET username = $1,
        email = $2
    WHERE id = $3
    RETURNING id,username,email;
  `,c.PostForm("username"),c.PostForm("email"),c.Param("id"))
  if err == sql.ErrNoRows {
    c.JSON(http.StatusNotFound,gin.H{
      "error": gin.H{
        "code": http.StatusNotFound,
        "message": fmt.Sprintf("No profile with id '%s'",c.Param("id")),
      },
    })
    return
  }
  if err != nil {
    log.Fatal(err)
  }
  c.JSON(http.StatusOK, gin.H{"data":profile})
}

func DeleteProfile(db *sqlx.DB,c *gin.Context) {
  profile := models.Profile{}
  err := db.Get(&profile,`
    DELETE FROM profile WHERE id = $1 RETURNING id,username,email;
  `,c.Param("id"))
  if err == sql.ErrNoRows {
    c.JSON(http.StatusNotFound,gin.H{
      "error": gin.H{
        "code": http.StatusNotFound,
        "message": fmt.Sprintf("No profile with id '%s'",c.Param("id")),
      },
    })
    return
  }
  if err != nil {
    log.Fatal(err)
  }
  c.JSON(http.StatusOK, gin.H{"data":profile})
}
