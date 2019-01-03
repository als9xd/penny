package models

type User struct {
  Id int `db:"id" json:"id"`
  Username string `db:"username" json:"username" binding:"required" form:"username"`
  Email string `db:"email" json:"email" form:"email"`
}

type Post struct {
  Id int `db:"id" json:"id"`
  Comment string `db:"comment" json:"comment" binding:"required" form:"comment"`
  UserId int `db:"user_id" json:"user_id" binding:"required" form:"user_id"`
  ThreadId int `db:"thread_id" json:"thread_id" binding:"required" form:"thread_id"`
}

type Thread struct {
  Id int `db:"id" json:"id"`
  Name string `db:"name" json:"name" binding:"required" form:"name"`
  UserId int `db:"user_id" json:"user_id" binding:"required" form:"user_id"`
}
