package models

type User struct {
  Id int `db:"id" json:"id"`
  Username string `db:"username" json:"username"`
}

type Post struct {
  Id int `db:"id" json:"id"`
  AuthorId int `db:"author_id" json:"author_id"`
  Text string `db:"text" json:"text"  binding:"required" form:"text"`
}

type Thread struct {
  Id int `db:"id" json:"id"`
  Name string `db:"name" json:"name"  binding:"required" form:"name"`
  CreatorId int `db:"creator_id" json:"creator_id"`
}
