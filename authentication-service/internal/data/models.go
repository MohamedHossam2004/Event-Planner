package data

import (
	"database/sql"
	"time"
)

const dbTimeout = 3 * time.Second

type Models struct {
	Tokens TokenModel
	Users  UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Tokens: TokenModel{DB: db},
		Users:  UserModel{DB: db},
	}
}
