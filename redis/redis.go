package redis
import (
	"github.com/garyburd/redigo/redis"
	"os"
	"time"
)

var (
	pool *redis.Pool
)

func init() {
	pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", os.Getenv("REDIS_DATA_SERVER"))
			if err != nil {
				return nil, err
			}

			return c, nil
		},
	}
}

func Scan(cursor int) (results ScanResults, err error) {
	client := pool.Get()
	data, err := redis.MultiBulk(client.Do("SCAN", cursor, "COUNT", "1000"))

	if err != nil {
		return results, err
	}

	results.Iterator, _ = redis.Int(data[0], nil)
	results.Keys, _ = redis.Strings(data[1], nil)

	return results, err
}

func Length(key string) (length int, err error) {
	client := pool.Get()
	length, err = redis.Int(client.Do("STRLEN", key))

	return length, err
}

func LengthMultiple(keys []string) (lengths []int, err error) {
	client := pool.Get()

	for _, key := range keys {
		client.Send("STRLEN", key)
	}

	client.Flush()

	keyCount := len(keys)
	for i := 0; i < keyCount; i++ {
		length, err := redis.Int(client.Receive())
		if err != nil {
			return lengths, err
		}

		lengths = append(lengths, length)
	}

	return
}

func Delete(key string) error {
	client := pool.Get()

	_, err := client.Do("DEL", key)
	if err != nil {
		return err
	}

	return nil
}

func MassDelete(keys []string) error {
	client := pool.Get()
	defer client.Close()

	for _, key := range keys {
		client.Send("DEL", key)
	}

	client.Flush()

	for _, _ = range keys {
		_, err := client.Receive()
		if err != nil {
			return err
		}
	}

	return nil
}