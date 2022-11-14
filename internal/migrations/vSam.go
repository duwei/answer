package migrations

import (
	"time"
	"xorm.io/xorm"
)

func addSamToUser(x *xorm.Engine) error {
	type User struct {
		SamId       int64     `xorm:"not null default 0 sam_id"`
		AccessToken string    `xorm:"not null default '' access_token"`
		ExpiredAt   time.Time `xorm:"TIMESTAMP expired_at"`
	}
	return x.Sync(new(User))
}
