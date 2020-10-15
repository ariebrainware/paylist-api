package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Configuration is a configuration model
type Configuration struct {
	Port int `json:"port"`

	Host         string `json:"host"`
	User         string `json:"user"`
	Database     string `json:"database"`
	DBPort       int    `json:"dbPort"`
	JWTSignature string `json:"jwtSignature"`
}

const password = "DB_PASSWORD"

var (
	Conf Configuration
	// DB is a exported connection
	DB *gorm.DB
)

// LoadConfiguration Database configuration using json file
func LoadConfiguration() {
	c := flag.String("c", "config/config.json", "Specify the file configuration.")
	flag.Parse()
	file, err := os.Open(*c)
	if err != nil {
		log.Fatal("can't open the file: ", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Conf)
	if err != nil {
		log.Fatal("can't decode config json", err)
	}
	connString := fmt.Sprintf(`user=%s password=%s host=%s port=%d dbname=%s sslmode=disable`, Conf.User, os.Getenv(password), Conf.Host, Conf.DBPort, Conf.Database)
	DB, err = gorm.Open("postgres", connString)
	DB.LogMode(true)
	if err != nil {
		fmt.Println(err)
		// panic("failed connect to database")
	}
}
