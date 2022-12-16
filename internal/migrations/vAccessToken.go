package migrations

import (
	"xorm.io/xorm"
)

func increaseAccessToken(x *xorm.Engine) error {
	type User struct {
		AccessToken string `xorm:"varchar(512) not null default '' access_token"`
	}
	return x.Sync(new(User))
}
