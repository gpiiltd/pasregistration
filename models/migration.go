package models

import (
	"database/sql"
	"log"
	"os"

	//Mysql driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/astaxie/beego"
)

var OldDB *sql.DB

//SetupOldDatabase sets up the data needed to mine the old DB
func SetupOldDatabase() {
	var DB = new(DBConfig)
	DB.Host = beego.AppConfig.String("databaseHost")
	DB.User = beego.AppConfig.String("databaseUsername")
	DB.Password = beego.AppConfig.String("databasePassword")
	DB.Database = beego.AppConfig.String("olddatabaseName")
	conn, err := sql.Open("mysql", DB.User+":"+DB.Password+"@/"+DB.Database+"?parseTime=true")
	if err != nil {
		panic(err)
	}
	OldDB = conn

}

//StartMining starts the process of extracting the data from old DB
func StartMining() {
	CreateUsers()

	return
}

//CreateUsers creates new users for the new system
func CreateUsers() {
	// var allUsers []User
	_, err := OldDB.Query(`SELECT names, email from companyEmployees`)
	if err != nil {
		LogError("Unable to get name and email from company employees")
	}
}

//LogError logs the error in a file
func LogError(errorMessage string) {
	// message := errors.New(errorMessage)
	f, err := os.OpenFile(beego.AppConfig.String("errorlogfile"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	return

}
