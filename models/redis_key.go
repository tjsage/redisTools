package models

import (
	. "github.com/tjsage/redisTools/database"
	"fmt"
	"strings"
)

type RedisKey struct {
	ID int
	Name string
	Size int
}

func (key *RedisKey) Save() error {
	return DB.Save(&key).Error
}

func (key *RedisKey) Delete() error {
	return DB.Delete(&key).Error
}

func MassSave(keys []RedisKey) error {
	var subQueries []string
	for _, key := range keys {
		vals := fmt.Sprintf("('%s', %d)", strings.Replace(key.Name, "'", "''", -1), key.Size)
		subQueries = append(subQueries, vals)
	}

	finalQuery := fmt.Sprintf("INSERT INTO redis_keys (name, size) VALUES %s", strings.Join(subQueries, ","))
	return DB.Exec(finalQuery).Error
}

func MassDelete(ids []int) error {
	var query = `
		DELETE FROM redis_keys WHERE id in (?)
	`

	return DB.Exec(query, ids).Error
}