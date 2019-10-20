package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/ariebrainware/paylist-api/model"
)

// Config is a configuration model
type Config struct {
	Host         string
	User         string
	Password     string
	Database     string
	Port         int
	JWTSignature string
}

var (
	config Config
	// DB is a exported connection
	DB *gorm.DB
)

// Conf Database configuration using json file
func Conf() {
	c := flag.String("c", "config/config.json", "Specify the file configuration.")
	flag.Parse()
	file, err := os.Open(*c)
	if err != nil {
		log.Fatal("can't open the file: ", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal("can't decode config json", err)
	}
	log.Println(config.Database)
	connString := fmt.Sprintf(`user=%s password=%s host=%s port=%d dbname=%s sslmode=disable`, config.User, config.Password, config.Host, config.Port, config.Database)
	DB, err = gorm.Open("postgres", connString)
	DB.LogMode(true)
	if err != nil {
		fmt.Println(err)
		panic("failed connect to database")
	}
	DB.AutoMigrate(&model.Paylist{}, &model.User{}, &model.Logging{})
	fmt.Println("Schema migrated!!")
}
