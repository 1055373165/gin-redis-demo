package model

import "time"

type Post struct {
	Id          int       `json:"id" db:"id"`
	PostId      string    `json:"post_id" db:"post_id"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	AuthorId    string    `json:"author_id" db:"author_id"`
	CommunityId string    `json:"community_id" db:"community_id"`
	Status      string    `json:"status" db:"status"`
	CreateTime  time.Time `json:"-"`
	UpdateTime  time.Time `json:"-"`
}

/* 帖子
缓存 id
帖子 id
帖子标题
帖子内容
作者 id
社区 id
状态
创建时间、更新时间
*/
