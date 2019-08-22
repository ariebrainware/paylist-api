package configdb

import (
	"github.com/jinzhu/gorm"
	"encoding/json"
	"os"
	"flag"
	"log"
	"fmt"
	"github.com/ariebrainware/paylist-api/model"
	_ "github.com/jinzhu/gorm/dialects/mysql"	
)

type Config struct {
	Db struct {
		Host string
		User string
		Password string
		Database string
	}
	Listen struct {
		Address string
		Port string
	}
}

var (
	config Config
	DB *gorm.DB
) 

//func Conf Database configuration using json file
func Conf() {
	c := flag.String("c","configdb/config.json", "Specify the file configuration.")
	flag.Parse()
	file, err := os.Open(*c)
	if err !=nil {
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
		fmt.Println(err.Error())
	}
	DB.AutoMigrate(&model.Paylist{})
	DB.AutoMigrate(&model.User{})
	fmt.Println("Schema migrated!!")
}