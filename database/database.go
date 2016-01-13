package database
import (
	"github.com/jinzhu/gorm"
	"os"
	"fmt"

	_"github.com/go-sql-driver/mysql"
)

var DB *gorm.DB

func init() {
	conn, err := gorm.Open("mysql", os.Getenv("DB_CONNECTION_STRING"))

	if err != nil {
		fmt.Printf("Error opening up connection to DB: %s \n", err.Error())
	}


	conn.LogMode(true)
	DB = &conn
}