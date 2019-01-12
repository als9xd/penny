package models

import "time"

type ProfilePublic struct {
	Id       int       `db:"id" json:"id"`
	Username string    `db:"username" json:"username" binding:"required" form:"username"`
	Avatar   string    `db:"avatar" json:"avatar" form:"avatar"`
	Created  time.Time `db:"created" json:"created"`
}

type ProfilePrivate struct {
	Id       int       `db:"id" json:"id"`
	Username string    `db:"username" json:"username" binding:"required" form:"username"`
	Email    string    `db:"email" json:"email" form:"email"`
	Avatar   string    `db:"avatar" json:"avatar" form:"avatar"`
	Created  time.Time `db:"created" json:"created"`
}

type Profile struct {
	Id       int       `db:"id" json:"id"`
	Username string    `db:"username" json:"username" binding:"required" form:"username"`
	Password string    `db:"password" json:"password" binding:"required" form:"password"`
	Email    string    `db:"email" json:"email" form:"email"`
	Avatar   string    `db:"avatar" json:"avatar" form:"avatar"`
	Created  time.Time `db:"created" json:"created"`
}

func (pub Profile) ToPublic() *ProfilePublic {
	return &ProfilePublic{
		Id:       pub.Id,
		Username: pub.Username,
		Avatar:   pub.Avatar,
	}
}

func (pub Profile) ToPrivate() *ProfilePrivate {
	return &ProfilePrivate{
		Id:       pub.Id,
		Username: pub.Username,
		Email:    pub.Email,
		Avatar:   pub.Avatar,
	}
}

type Post struct {
	Id              int        `db:"id" json:"id"`
	Name            string     `db:"name" json:"name" binding:"required" form:"name"`
	Description     string     `db:"description" json:"description" form:"description"`
	ProfileId       int        `db:"profile_id" json:"profile_id"`
	ProfileAvatar   string     `db:"profile_avatar" json:"profile_avatar"`
	ProfileUsername string     `db:"profile_username" json:"profile_username"`
	ThreadId        int        `db:"thread_id" json:"thread_id" binding:"required" form:"thread_id"`
	ThreadName      string     `db:"thread_name" json:"thread_name"`
	ThreadAvatar    string     `db:"thread_avatar" json:"thread_avatar"`
	Avatar          string     `db:"avatar" json:"avatar" form:"avatar"`
	Created         time.Time  `db:"created" json:"created"`
	LastEdited      *time.Time `db:"last_edited" json:"last_edited"`
}

type Comment struct {
	Id              int        `db:"id" json:"id"`
	Value           string     `db:"value" json:"value" binding:"required" form:"value"`
	Level           int        `db:"level" json:"level"`
	ParentCommentId *int       `db:"parent_comment_id" json:"parent_comment_id" form:"parent_comment_id"`
	PostId          int        `db:"post_id" json:"post_id" binding:"required" form:"post_id"`
	ProfileId       int        `db:"profile_id" json:"profile_id"`
	ProfileAvatar   string     `db:"profile_avatar" json:"profile_avatar"`
	ProfileUsername string     `db:"profile_username" json:"profile_username"`
	Created         time.Time  `db:"created" json:"created"`
	LastEdited      *time.Time `db:"last_edited" json:"last_edited"`
}

type Thread struct {
	Id          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name" binding:"required" form:"name"`
	Description string    `db:"description" json:"description" form:"description"`
	ProfileId   int       `db:"profile_id" json:"profile_id"`
	PostCount   int       `db:"post_count" json:"post_count"`
	Avatar      string    `db:"avatar" json:"avatar" form:"avatar"`
	Created     time.Time `db:"created" json:"created"`
}
