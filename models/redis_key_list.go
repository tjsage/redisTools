package models

import (
	. "github.com/tjsage/redisTools/database"
)

type RedisKeyList struct {
	Keys []RedisKey
	Count int
}

func (list *RedisKeyList) FindByPattern(pattern string) error {
	err := DB.Where("name like ?", pattern).Find(&list.Keys).Error
	if err != nil {
		return err
	}

	list.Count = len(list.Keys)
	return nil
}