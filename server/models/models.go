package models

type Thread struct {
  Id int `db:"id" json:"id"`
  Name string `db:"name" json:"name"  binding:"required" form:"name"`
}
