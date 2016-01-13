package models

import (
	. "github.com/tjsage/redisTools/database"
)

func init() {
	DB.AutoMigrate(&RedisKey{})
}