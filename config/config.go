package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/tkanos/gonfig"

	"github.com/ariebrainware/paylist-api/model"
)

// Configuration is a configuration model
type Configuration struct {
	Host         string
	User         string
	Password     string
	Database     string
	DBPort       int
	JWTSignature string
	Port         int
}

var (
	config Configuration
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
	connString := fmt.Sprintf(`user=%s password=%s host=%s port=%d dbname=%s sslmode=disable`, config.User, config.Password, config.Host, config.DBPort, config.Database)
	DB, err = gorm.Open("postgres", connString)
	DB.LogMode(true)
	if err != nil {
		fmt.Println(err)
		// panic("failed connect to database")
	}
	DB.AutoMigrate(&model.Paylist{}, &model.User{}, &model.Logging{})
	fmt.Println("Schema migrated!!")
}

// Misc is a global variable to handle exported configuration
var Misc Configuration

// LoadConfiguration is a function to export value from config.json based on defined struct
func LoadConfiguration(path string) {
	Misc = Configuration{}
	err := gonfig.GetConf(path, &Misc)
	if err != nil {
		panic(err)
	}
}
