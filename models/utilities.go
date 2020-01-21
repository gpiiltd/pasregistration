package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	//Mysql driver
	_ "github.com/go-sql-driver/mysql"
)

//Conn db
var Conn *gorm.DB

//SetupDatabase  handles the connection to the database
func SetupDatabase() {

	//DB hold new instance of databse objects
	var DB = new(DBConfig)
	DB.Host = beego.AppConfig.String("databaseHost")
	DB.User = beego.AppConfig.String("databaseUsername")
	DB.Password = beego.AppConfig.String("databasePassword")
	DB.Database = beego.AppConfig.String("databaseName")
	conn, err := gorm.Open("mysql", DB.User+":"+DB.Password+"@/"+DB.Database+"?parseTime=true")
	if err != nil {
		panic(err)
	}
	Conn = conn
	SetupTables()
}

//SetupTables set up the database tables
func SetupTables() {
	Conn.AutoMigrate(&User{})
	Conn.AutoMigrate(&Subsidiaries{})
	Conn.AutoMigrate(&Departments{})
	Conn.AutoMigrate(&AttemptedLogin{})
	Conn.AutoMigrate(&Roles{})
	Conn.AutoMigrate(&PasswordRecoveryData{})
	if findSubsidiaries := Conn.Find(&Subsidiaries{}); findSubsidiaries.Error != nil {
		go SetupSubsidiaries()
	}
	PrepareDepartment()
	SetupOldDatabase()
	CreateTeamLead()
	// StartMining()
	// go SetupDepartments()
}

//PrepareDepartment ready the system to create departments
func PrepareDepartment() {
	cesl := "CESL"
	gpi := "GPI"
	ctes := "CTES"
	cepl := "CEPL"
	templatePath := beego.AppConfig.String("companydatapath")
	if findCESL := Conn.Where("subsidiary = ?", cesl).Find(&Departments{}); findCESL.Error != nil {
		SetupDepartments(templatePath+"company-department-cesl.json", cesl)
	}
	if findGPI := Conn.Where("subsidiary = ?", gpi).Find(&Departments{}); findGPI.Error != nil {
		SetupDepartments(templatePath+"company-department-gpi.json", gpi)
	}
	if findCTES := Conn.Where("subsidiary = ?", ctes).Find(&Departments{}); findCTES.Error != nil {
		SetupDepartments(templatePath+"company-department-ctes.json", ctes)
	}
	if findCPTL := Conn.Where("subsidiary = ?", cepl).Find(&Departments{}); findCPTL.Error != nil {
		SetupDepartments(templatePath+"company-department-cepl.json", cepl)
	}
}

//SetupSubsidiaries sets up the subsidiary of companies using the csv in the app
func SetupSubsidiaries() {
	type department struct {
		Name string `json:"departments"`
	}
	type subsidiaries struct {
		Subsidiary string       `json:"subsidiary"`
		Department []department `json:"deparments"`
	}
	type companyData struct {
		Subsidiary []subsidiaries `json:"subsidiaries"`
	}
	companyDataFile, _ := ioutil.ReadFile(beego.AppConfig.String("companydatapath") + "company-data.json")
	var subsidiaryArray []companyData
	err := json.Unmarshal(companyDataFile, &subsidiaryArray)
	if err != nil {
		log.Println(err.Error())
	}
	var allSubsidiary Subsidiaries
	// log.Println(subsidiaryArray)
	for _, subsid := range subsidiaryArray {
		var tempSubsid []subsidiaries
		tempSubsid = subsid.Subsidiary
		allSubsidiary.ID = 0
		for _, tempSub := range tempSubsid {
			allSubsidiary.ID++
			allSubsidiary.Subsidiary = tempSub.Subsidiary
			// log.Println(allSubsidiary)
			Conn.Create(&allSubsidiary)
		}
	}
}

//SetupDepartments sets up all the deparement in the system
func SetupDepartments(filePath string, subsidiaryName string) {
	type department struct {
		Name string `json:"name"`
	}
	type subsidiaries struct {
		Subsidiary string       `json:"subsidiary"`
		Department []department `json:"deparments"`
	}
	var subsidiaryArray subsidiaries
	companyDataFile, _ := ioutil.ReadFile(filePath)
	err := json.Unmarshal(companyDataFile, &subsidiaryArray)
	if err != nil {
		log.Println(err.Error())
	}
	departmentArray := subsidiaryArray.Department
	var allSubsidiaries []Subsidiaries
	Conn.Find(&allSubsidiaries)
	var tempDepartmentData Departments
	Conn.Last(&tempDepartmentData)
	var departments Departments
	for _, subsidiary := range allSubsidiaries {
		if subsidiary.Subsidiary == subsidiaryName {
			departments.ID = tempDepartmentData.ID
			for _, dept := range departmentArray {
				departments.ID = departments.ID + 1
				departments.Subsidiary = subsidiary.Subsidiary
				departments.SubsidiaryID = subsidiary.ID
				departments.Department = dept.Name
				Conn.Create(departments)
				// log.Println(departments)
			}
		}
	}
}

//ErrorResponse structures error messages
func ErrorResponse(code int, message string) interface{} {
	type errorResponseData struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	var response errorResponseData
	response.Code = code
	response.Message = message

	return response
}

//ValidationResponse returns responses for validation purposes
func ValidationResponse(code int, body bool) interface{} {
	type validationResponseData struct {
		Code int  `json:"code"`
		Body bool `json:"body"`
	}

	var response validationResponseData
	response.Code = code
	response.Body = body

	return response
}

//ValidResponse structures the data for all API response.
func ValidResponse(code int, body interface{}, message string) interface{} {
	type validResponseData struct {
		Code    int         `json:"code"`
		Body    interface{} `json:"body"`
		Message string      `json:"message"`
	}
	var response validResponseData
	response.Code = code
	response.Body = body
	response.Message = message

	return response
}

//HashPassword encrypts "Hash" a password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//CheckPasswordHash checks if password and hash is the same. Returns true or false.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//GetTokenString generates and returns a string.
func GetTokenString(email string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"expire": time.Now().Add(time.Minute * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte(beego.AppConfig.String("jwtkey")))
	if err != nil {
		panic(err)
	}

	return tokenString
}

//RecordLoginAttempt is triggered when a user tries to login
func RecordLoginAttempt(username string) AttemptedLogin {
	var atl AttemptedLogin
	atl.Username = username

	Conn.Create(&atl)

	return atl
}

//UpdateLoginAttempt is triggered at the end of login execution
func UpdateLoginAttempt(atl AttemptedLogin, status string, message string) {
	atl.Status = status
	atl.Message = message

	Conn.Model(&atl).Updates(AttemptedLogin{Status: atl.Status, Message: atl.Message})
}

//ValidToken checks if a token is valid
func ValidToken(wholeToken string) bool {
	splitString := strings.Split(wholeToken, ",")
	if splitString[0] != beego.AppConfig.String("tokenprefix") {
		return false
	}

	if splitString[1] == "" {
		return false
	}

	return true
}

//TokenExpire checks if the user token is valid and hasn't expired
func TokenExpire(tokenS string) bool {
	code, tokenString := SplitToken(tokenS)
	if code != 200 {
		log.Println("At split")
		return false
	}
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(beego.AppConfig.String("jwtkey")), nil
	})
	if err != nil {
		return false
	}
	var expireAt float64
	nowTime := time.Now().Add(time.Minute * 1).Unix()
	for key, val := range claims {
		if key == "expire" {
			expireAt = val.(float64)
		}
	}
	tm := float64(nowTime)
	diff := tm - expireAt
	if diff >= 360000 {
		return false
	}

	return true
}

//GetEmailFromToken returns the email used to generate a token string
func GetEmailFromToken(token string) string {
	wholeString := strings.Split(token, ",")
	tokenString := wholeString[1]
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(beego.AppConfig.String("jwtkey")), nil
	})
	if err != nil {
		return err.Error()
	}
	var email string
	for key, val := range claims {
		if key == "email" {
			email = val.(string)
		}
	}

	return email
}

//ValidateTokenString validates a token string and return true or false
func ValidateTokenString(token string) interface{} {
	type tokenReturn struct {
		Code uint64 `json:"code"`
		Body bool   `json:"body"`
		User User   `json:"user"`
	}
	code, _ := SplitToken(token)
	if code != 200 {
		return ValidResponse(200, "false", "Error when spliting user token")
	}
	isTokenExpired := TokenExpire(token)
	if isTokenExpired != true {
		return ValidResponse(200, "false", "Token Expired")
	}
	var u User
	code, u = GetUserFromTokenString(token)
	if code != 200 {
		return ValidResponse(200, "false", "Unable to get user information.")
	}
	var response tokenReturn
	response.Code = 200
	response.Body = true
	response.User = u
	return response
}

//SplitToken splits token token and
func SplitToken(wholeToken string) (int, string) {
	splitString := strings.Split(wholeToken, ",")
	if splitString[0] != beego.AppConfig.String("tokenprefix") {
		return 403, splitString[1]
	}

	return 200, splitString[1]
}

//GetRoleFromCode gets a user code string from role code
func GetRoleFromCode(code uint64) string {
	if code == 999 {
		return "the Application Owner"
	}
	if code == 99 {
		return "an Administrator"
	}
	if code == 88 {
		return "Front Desk Officer"
	}
	if code == 77 {
		return "Head of HR"
	}
	if code == 66 {
		return "Team Lead"
	}
	if code == 55 {
		return "Regular User"
	}
	if code == 0 {
		return "a User"
	}

	return "Invalid role Code"
}

//GenerateRandCode gets new random codes
func GenerateRandCode(length int) string {
	return StringWithCharset(length, charset)
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

//StringWithCharset chill
func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

//ValidateRecoveryCode checks if a password recovery code associated with a user is valid.
func ValidateRecoveryCode(validationObject PasswordRecoveryData) interface{} {
	var recovery PasswordRecoveryData
	if validate := Conn.Where("email = ? AND code = ?", validationObject.Email, validationObject.Code).Find(&recovery); validate.Error != nil {
		return ErrorResponse(404, validate.Error.Error())
	}
	return ErrorResponse(200, "Code Valid")
}

//ValidateUserRoles checks if a user code is valid.
func ValidateUserRoles(validationObject ValidateRole) interface{} {
	var role Roles
	if validate := Conn.Where("user_id = ? AND code = ?", validationObject.UserID, validationObject.RoleCode).Find(&role); validate.Error != nil {
		return ValidationResponse(200, false)
	}
	return ValidationResponse(200, true)
}

//ValidateTeamLeadUser checks if a user is a team lead
func ValidateTeamLeadUser(user User) interface{} {
	var role Roles
	if validate := Conn.Where("user_id = ?", user.ID).Find(&role); validate.Error != nil {
		return ValidationResponse(200, false)
	}
	return ValidationResponse(200, true)
}

//GetSubsidiaries gets the list of all subsidiaries
func GetSubsidiaries() []Subsidiaries {
	var subsidiaryList []Subsidiaries
	Conn.Find(&subsidiaryList)
	return subsidiaryList
}

//GetDeparments gets the list of all departments
func GetDeparments() []Departments {
	var departmentList []Departments
	Conn.Find(&departmentList)
	return departmentList
}

//GetSubsidiaryDepartments gets the list of all departments
func GetSubsidiaryDepartments(subsidiaryID string) []Departments {
	var departmentList []Departments
	if findDepartments := Conn.Where("subsidiary_id = ?", subsidiaryID).Find(&departmentList); findDepartments.Error != nil {
		return nil
	}
	return departmentList
}

//GetSubsidiaryByID gets a list of subsidiary by id
func GetSubsidiaryByID(subsidiaryID uint64) (Subsidiaries, error) {
	var subsidiary Subsidiaries
	if getSubsidiary := Conn.Where("id = ?", subsidiaryID).Find(&subsidiary); getSubsidiary != nil {
		return subsidiary, getSubsidiary.Error
	}
	return subsidiary, nil
}

//GetDepartmentByID gets a department by ID
func GetDepartmentByID(deptID uint64) (Departments, error) {
	var department Departments
	if getDepartment := Conn.Where("id = ?", deptID).Find(&department); getDepartment.Error != nil {
		return department, getDepartment.Error
	}
	return department, nil
}

//UploadImage uploads a document and return the document path.
func UploadImage(file multipart.File, handle *multipart.FileHeader, user User) interface{} {
	uploadpath := beego.AppConfig.String("uploadimage")
	path := beego.AppConfig.String("imagepath")
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return ErrorResponse(401, "Unable to read file")
	}
	nowTime := time.Now().Format("20060102150405")
	fileName := uploadpath + nowTime + "-" + user.Email + "-" + handle.Filename
	err = ioutil.WriteFile(fileName, data, 0666)
	if err != nil {
		return ErrorResponse(401, "Unable to save file")
	}
	fileName = path + nowTime + "-" + user.Email + "-" + handle.Filename
	returnString := beego.AppConfig.String("serverip") + fileName
	return ErrorResponse(200, returnString)
}

//IsFrontDesk checks if a user is a front desk officer
func IsFrontDesk(user User) bool {
	var frontDeskRole Roles
	if getRole := Conn.Where("code = 88 AND user_id = ?", user.ID).Find(&frontDeskRole); getRole.Error != nil {
		return false
	}
	return true
}

//IsTeamLead checks if a user is a team lead
func IsTeamLead(user User) bool {
	var teamLeadRole Roles
	if getRole := Conn.Where("code = 66 AND user_id = ?", user.ID).Find(&teamLeadRole); getRole.Error != nil {
		return false
	}
	return true
}

//IsVMSAdmin checks if a user is a vms admin
func IsVMSAdmin(user User) bool {
	var vmsAdminRole Roles
	if getRole := Conn.Where("code = 44 AND user_id = ?", user.ID).Find(&vmsAdminRole); getRole.Error != nil {
		return false
	}
	return true
}

//IsAdmin checks if a user is an admin
func IsAdmin(user User) bool {
	var adminRole Roles
	if getRole := Conn.Where("code = 99 AND user_id = ?", user.ID).Find(&adminRole); getRole.Error != nil {
		return false
	}
	return true
}

//IsHRO checks if a user is an HRO
func IsHRO(user User) bool {
	var HRRole Roles
	if getRole := Conn.Where("code = 77 AND user_id = ?", user.ID).Find(&HRRole); getRole.Error != nil {
		return false
	}
	return true
}
