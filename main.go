package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/tjsage/redisTools/models"
	"github.com/tjsage/redisTools/redis"
)

var functions =  map[string]func() {
	"keys": keys,
	"delete": delete,
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("You must provide the name of the function you want to execute")
		return
	}

	functionName := os.Args[1]

	fn, found := functions[functionName]
	if !found {
		fmt.Println("Didnt find the function: ", functionName)
		return
	}

	fn()
}

func keys() {

	var iterator = 0
	if len(os.Args) > 2 {
		iterator, _ = strconv.Atoi(os.Args[2])
	}

	for true {
		fmt.Println("Scanning for keys, iterator: ", iterator)

		// Scan our redis data
		results, err := redis.Scan(iterator)
		if err != nil {
			fmt.Println("You had an error: ", err.Error())
			return
		}

		iterator = results.Iterator


		// GetLengths
		lengths, err := redis.LengthMultiple(results.Keys)
		if err != nil {
			fmt.Println("Error getting lengths: ", err.Error())
			return
		}


		var records[] models.RedisKey
		for i, key := range results.Keys {
			var model models.RedisKey
			model.Name = key
			model.Size = lengths[i]

			records = append(records, model)
		}

		err = models.MassSave(records)
		if err != nil {
			fmt.Printf("Unable to mass save: %s\n", err.Error())
			return
		}

		if iterator == 0 {
			fmt.Println("We're done!")
			return
		}


		//time.Sleep(time.Second)
	}
}

func delete() {
	if len(os.Args) < 3 {
		fmt.Println("You must provide the delete pattern")
		return
	}

	fmt.Println("Started")
	deletePattern := os.Args[2]

	var list models.RedisKeyList
	err := list.FindByPattern(deletePattern)
	if err != nil {
		fmt.Println("Error deleting pattern: %s", err.Error())
		return
	}


	for i := 0; i < list.Count; i += 1000 {
		max := i + 1000
		if max > list.Count {
			max = list.Count
		}

		subKeys := list.Keys[i:max]

		// Extract key names and ids
		var ids []int
		var keyNames []string
		for _, key := range subKeys {
			keyNames = append(keyNames, key.Name)
			ids = append(ids, key.ID)
		}

		for _, key := range keyNames {
			fmt.Println("Deleting: ", key)
		}

		// Delete in redis
		err := redis.MassDelete(keyNames)
		if err != nil {
			fmt.Printf("Error running deleting keys: %s\n", err.Error())
			return
		}

		// Delete from database
		err = models.MassDelete(ids)
		if err != nil {
			fmt.Printf("Error mass deleting: %s", err.Error())
			return
		}

		time.Sleep(20 *time.Millisecond)
	}

	fmt.Println("Done")
}
