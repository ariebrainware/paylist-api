package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ariebrainware/paylist-api/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Config struct {
	Db struct {
		Host     string
		User     string
		Password string
		Database string
	}
	Listen struct {
		Address string
		Port    string
	}
}

var (
	config Config
	DB     *gorm.DB
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
	log.Println(config.Db.Database)
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", config.Db.User, config.Db.Password, config.Db.Host, config.Listen.Port, config.Db.Database)
	DB, err = gorm.Open("mysql", connString)
	if err != nil {
		panic("failed connect to database")
	}
	DB.AutoMigrate(&model.Paylist{})
	DB.AutoMigrate(&model.User{})
	fmt.Println("Schema migrated!!")
}
