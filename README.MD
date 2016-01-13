#redisTools

redisTools is a program written in Go for analyzing all the items in a Redis server. I created this program as a way of figuring out what types of data were taking up the most space in a Redis server. 

## Setup
### Step 1: Setup a MySQL database
You'll need a MySQL database for the program to use. The program will automatically take care of any table creation. 

### Step 2: Configure environment variables
You'll need to setup the following environment variables:

```
export DB_CONNECTION_STRING=someuser:password@tcp(localhost:3306)/DatabaseName?charset=utf8&parseTime=true
export REDIS_DATA_SERVER=localhost:6379
```

### Step 3: Build
```
go build -o bin/redisTools
```

## Usage
```
redisTools keys iterator
redisTools delete pattern
```

### keys
This command fills the redis_keys table with the key and size of every item in your Redis server. 

***iterator*** - This represents the iterator that Redis returns when performing a `SCAN` command. You can use
this iterator as a way of picking up where you left off if you need to exit the program at any time. If you 
don't provide the iterator, then the program will start from the beginning key `REDIS SCAN 0`

***Example***

```
redisTools keys 3921
```
This would start the `SCAN` command at iterator `3921`

#### delete
This command allows you to delete all keys matching a pattern. The delete command does not rescan Redis. Instead,
it uses the keys that have been recorded into your `redis_keys` table. 

***pattern*** - The pattern you want used when searching `redis_keys` for a match. The pattern will be used in a like
operator against the `name` field in `redis_keys`.

***Example***: 

```
redisTools del '%profile'
```

This would delete every key ending in `profile`