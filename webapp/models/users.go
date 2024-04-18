package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID             uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	GitHubUsername *string   `bun:"github_username,type:varchar(128),unique" json:"githubUsername"`
	Name           string    `bun:"name,type:varchar(128),notnull" json:"name"`
	AvatarUrl      string    `bun:"avatar_url,type:varchar(256),notnull,default:''" json:"avatarUrl"`
	CreatedAt      time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"createdAt"`
}

func NewUser(u GitHubUser) User {
	return User{GitHubUsername: &u.Login, Name: u.Name, AvatarUrl: u.AvatarUrl}
}

func NewAnonymousUser() User {
	return User{Name: "Anonymous", AvatarUrl: "https://avatars.githubusercontent.com/u/583231?v=4"}
}
