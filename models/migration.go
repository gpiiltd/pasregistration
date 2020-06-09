package models

import (
	"database/sql"
	"log"

	//Mysql driver
	"github.com/astaxie/beego/logs"
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
		log.Println(err.Error())
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
	var allUsersArray []User
	stmt, err := OldDB.Query(`SELECT name, email, password, department from companyEmployees`)
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
	var users User
	for stmt.Next() {
		stmt.Scan(&users.FullName, &users.Email, &users.Password, &users.Department)
		allUsersArray = append(allUsersArray, users)
	}
	CreateUserDB(allUsersArray)

}

//GetDepartmentFromName gets the user department from department name
func GetDepartmentFromName(department string) Departments {
	var thisDepartment Departments
	Conn.Where("department = ?", department).Find(&thisDepartment)
	return thisDepartment
}

//CreateUserDB creates the user in New DB
func CreateUserDB(allUsers []User) {
	var counter uint64
	counter = 0
	var userDepartment Departments
	for _, user := range allUsers {
		if user.Email != "noreply@my-gpi.io" {
			counter = counter + 1
			userDepartment = GetDepartmentFromName(user.Department)
			user.Department = userDepartment.Department
			user.DepartmentID = userDepartment.ID
			user.Subsidiary = userDepartment.Subsidiary
			user.SubsidiaryID = userDepartment.SubsidiaryID
			Conn.Create(&user)
		}
	}
	return
}

//CreateTeamLead creates a team lead role
func CreateTeamLead() {
	var allTeamLead []User
	stmt, err := OldDB.Query(`SELECT lead from teams`)
	if err != nil {
		panic(err.Error())
	}
	var teamLead User
	for stmt.Next() {
		stmt.Scan(&teamLead.ID)
		allTeamLead = append(allTeamLead, teamLead)
	}

	var allUsersArray []User
	stmts, err := OldDB.Query(`SELECT ID, name from companyEmployees`)
	if err != nil {
		panic(err.Error())
	}
	var user User
	for stmts.Next() {
		stmts.Scan(&user.ID, &user.FullName)
		allUsersArray = append(allUsersArray, user)
	}

	var teamLeadArray []User

	for _, newUser := range allUsersArray {
		for _, newTeamLead := range allTeamLead {
			if newUser.ID == newTeamLead.ID {
				teamLeadArray = append(teamLeadArray, newUser)
				CreateteamLeadRole(newUser)
			}
		}
	}

	return
}

//CreateteamLeadRole creates the team lead role
func CreateteamLeadRole(lead User) {
	var lastUser User
	Conn.Where("full_name = ?", lead.FullName).Find(&lastUser)
	AddTeamLead(lastUser)
}

//LogError logs the error in a file
func LogError(errorMessage string) {
	log := logs.GetBeeLogger()
	logs.SetLogger(logs.AdapterFile, `{"filename": "./error-logs/log.log", "level": 7, "maxlines": 0, "maxsize":0,"daily":true,"maxdays":10,"color":true}`)
	log.Error(errorMessage, "error")
	return

}
