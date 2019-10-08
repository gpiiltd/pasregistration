package models

import (
	"pasregistration/controllers/mailer"

	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
)

//RegisterUsers registers new User
func RegisterUsers(u User) interface{} {
	userExist := UserExist(u)
	if userExist == true {
		return ErrorResponse(403, "User already exists. Kindly Login.")
	}
	hashPassword, _ := HashPassword(u.Password)
	u.Password = hashPassword
	Conn.Create(&u)
	SetupRole(u)
	SendRegistrationEmail(u)
	u.Password = ""
	tokenString := GetTokenString(u.Email)
	var roles []Roles
	roles = AppAdminRoles()

	if u.Email != beego.AppConfig.String("adminmail") {
		Conn.Where("user_id = ?", u.ID).Find(&roles)
	}
	return StructureLoginData(200, u, roles, tokenString)
	// return ValidResponse(200, tokenString, "success")
}

//UpdateProfile updates a user profile
func UpdateProfile(update User, authenticatedUser User) interface{} {
	verifiObject := VerifiUpdateProfile(update)
	if verifiObject != true {
		return ErrorResponse(403, "User object is not valid")
	}
	if update.Email != authenticatedUser.Email {
		return ErrorResponse(403, "Unauthorized access")
	}
	var subsidiary Subsidiaries
	subsidiary, err := GetSubsidiaryByID(update.SubsidiaryID)
	if err != nil {
		return ErrorResponse(403, "Invalid Subsidiary ID")
	}
	var department Departments
	department, err = GetDepartmentByID(update.DepartmentID)
	if err != nil {
		return ErrorResponse(403, "Invalid Department ID")
	}
	var u User
	Conn.Model(&u).Where("id = ?", authenticatedUser.ID).Updates(User{JobTitle: update.JobTitle, Number: update.Number, Location: update.Location, Subsidiary: subsidiary.Subsidiary, SubsidiaryID: subsidiary.ID, Department: department.Department, DepartmentID: department.ID, LinkedIn: update.LinkedIn})
	return ValidResponse(200, u, "User profile changed succesfully")
}

//UpdateProfileImage updates a user profile avatar
func UpdateProfileImage(update User, authenticatedUser User) interface{} {
	if update.Email == "" || update.Image == "" {
		return ErrorResponse(403, "User object is not valid")
	}
	if update.Email != authenticatedUser.Email {
		return ErrorResponse(403, "Unauthorized access")
	}
	var u User
	Conn.Model(&u).Where("id = ?", authenticatedUser.ID).Updates(User{Image: update.Image})
	return ValidResponse(200, u, "User profile avatar changed succesfully")
}

//VerifiUpdateProfile verifis that the fullname and the email is not left blank
func VerifiUpdateProfile(u User) bool {
	if u.FullName == "" {
		return false
	}
	if u.Email == "" {
		return false
	}
	return true
}

//SetupRole sets up a user role
func SetupRole(u User) {
	var role Roles
	role.User = u.FullName
	role.UserID = u.ID
	role.Code = u.Role
	role.Role = GetRoleFromCode(u.Role)

	if findRole := Conn.Where("code = ? AND user_id = ?").Find(&role); findRole.Error != nil {
		Conn.Create(&role)
	}

	return
}

//Login logs a user in
func Login(email, password string) interface{} {
	var u User
	var roles []Roles
	if email == beego.AppConfig.String("adminmail") && password == beego.AppConfig.String("adminpassword") {
		roles = AppAdminRoles()
		tokenString := GetTokenString(email)
		return StructureLoginData(200, u, roles, tokenString)
	}
	var loginAttempt AttemptedLogin
	u.Email = email
	loginAttempt = RecordLoginAttempt(email)
	u, err := GetUserDataEmail(email)
	if err != nil {
		UpdateLoginAttempt(loginAttempt, "Failed", "User Does not exist")
		return ErrorResponse(403, "User Does not exist. Kindly Sign up.")
	}
	//Paswords are usually converted to all lowercase. From the pas server. So password will be replaced with loweredpasswords
	// loweredPassword := strings.ToLower(password)
	passwordMatch := CheckPasswordHash(password, u.Password)

	if passwordMatch != true {
		UpdateLoginAttempt(loginAttempt, "Failed", "Invalid Login Credentials")
		return ErrorResponse(403, "Invalid Login Details")
	}

	u.Password = ""
	if findRole := Conn.Where("user_id = ?", u.ID).Find(&roles); findRole.Error != nil {
		return ErrorResponse(403, "Unable to get user roles")
	}
	tokenString := GetTokenString(u.Email)
	UpdateLoginAttempt(loginAttempt, "Success", "Success")
	return StructureLoginData(200, u, roles, tokenString)
}

//SendRegistrationEmail send confirmation registration email to the registered user.
func SendRegistrationEmail(u User) {
	path := beego.AppConfig.String("mailtemplatepath") + "registration.html"
	mailSubject := "Registration to PAS Successful"
	newRequest := mailer.NewRequest(u.Email, mailSubject)
	data := mailer.Data{}
	data.User = u.FullName

	go newRequest.Send(path, data)

	return
}

//StructureLoginData organizes login data to contain user roles and permission
func StructureLoginData(code int, u User, roles []Roles, tokenString string) interface{} {
	//LoginData structures user Login Information
	type loginDataStructure struct {
		Code   int     `json:"code"`
		User   User    `json:"user_data"`
		Role   []Roles `json:"user_roles"`
		Status string  `json:"token"`
	}

	var loginResponse loginDataStructure
	loginResponse.Code = code
	loginResponse.User = u
	loginResponse.Role = roles
	loginResponse.Status = tokenString

	return loginResponse
}

//UserExist checks if a user exists using email
func UserExist(u User) bool {
	var users User
	if findUser := Conn.Where("email = ?", u.Email).Find(&users); findUser.Error != nil {
		return false
	}
	return true
}

//IfEmailExist checks if a user email exists
func IfEmailExist(email string) bool {
	var u User
	if ifEmailExist := Conn.Where("email = ?", email).Find(&u); ifEmailExist.Error != nil {
		return false
	}
	return true
}

//GetUserDataEmail gets user data from email
func GetUserDataEmail(email string) (User, error) {
	var u User
	data := Conn.Where("email = ?", email).Find(&u)
	if data != nil && data.Error != nil {
		return u, data.Error
	}

	return u, nil
}

//GetUserFromNumber get user data from number
func GetUserFromNumber(number string) (User, error) {
	var u User
	data := Conn.Where("number = ?", number).Find(&u)
	if data != nil && data.Error != nil {
		return u, data.Error
	}

	return u, nil
}

//AppAdminRoles get the app admin roles into a Roles model
func AppAdminRoles() []Roles {
	var roles Roles
	var rolesArray []Roles
	roles.Role = "admin"
	roles.Code = 99
	roles.User = beego.AppConfig.String("adminmail")
	rolesArray = append(rolesArray, roles)

	return rolesArray
}

//PasswordRecovery recovers lost password
func PasswordRecovery(email string) interface{} {
	emailExist := IfEmailExist(email)
	if emailExist != true {
		return ErrorResponse(403, "User does not exist. Kindly Register.")
	}

	var passwordRecovery PasswordRecoveryData
	Conn.Where("email = ?", email).Delete(&passwordRecovery)
	passwordRecovery.Email = email
	code := GenerateRandCode(16)
	passwordRecovery.Code = code
	Conn.Create(&passwordRecovery)
	passwordRecovery.Code = beego.AppConfig.String("resetpasswordpage") + "?email=" + email + "&code=" + code
	SendPasswordReocoveryEmail(passwordRecovery)

	//use ErrorResponse since we're just returning a string
	return ErrorResponse(200, "Password recovery instruction has been sent to your email. Check spam too. ")
}

//SendPasswordReocoveryEmail sends a password reocovery email.
func SendPasswordReocoveryEmail(recovery PasswordRecoveryData) {
	path := beego.AppConfig.String("mailtemplatepath") + "recover-password.html"
	mailSubject := "Recover lost Password"
	newRequest := mailer.NewRequest(recovery.Email, mailSubject)
	data := mailer.Data{}
	data.Email = recovery.Email
	data.Code = recovery.Code

	go newRequest.Send(path, data)

	return
}

//SendResetPasswordEmail sends reset password email.
func SendResetPasswordEmail(u User) {
	path := beego.AppConfig.String("mailtemplatepath") + "reset-password.html"
	mailSubject := "Password changed successfully."
	newRequest := mailer.NewRequest(u.Email, mailSubject)
	data := mailer.Data{}
	data.User = u.FullName

	go newRequest.Send(path, data)

	return
}

//ResetPassword resets a user password
func ResetPassword(recoveryData ResetPasswordData) interface{} {
	if recoveryData.Email == "" || recoveryData.Code == "" || recoveryData.Password == "" {
		return ErrorResponse(403, "User not found.")
	}
	var verifyData PasswordRecoveryData
	if verifyCode := Conn.Where("email = ? AND code = ?", recoveryData.Email, recoveryData.Code).Find(&verifyData); verifyCode.Error != nil {
		return ErrorResponse(403, verifyCode.Error.Error())
	}
	var user User
	email := recoveryData.Email
	password := recoveryData.Password
	hashPassword, _ := HashPassword(password)
	user.Password = hashPassword
	Conn.Model(&user).Where("email = ?", email).Update("password", user.Password)
	Conn.Where("email = ?", email).Delete(&PasswordRecoveryData{})
	user, err := GetUserDataEmail(email)
	if err != nil {
		return ErrorResponse(403, "Invalid user data")
	}
	user.Password = ""
	SendResetPasswordEmail(user)
	var roles []Roles
	roles = AppAdminRoles()

	if user.Email != beego.AppConfig.String("adminmail") {
		Conn.Where("user_id = ?", user.ID).Find(&roles)
	}
	tokenString := GetTokenString(user.Email)
	return StructureLoginData(200, user, roles, tokenString)
}

//ChangePassword changes user password
func ChangePassword(update User, user User) interface{} {
	if user.Email != update.Email || update.Password == "" {
		return ErrorResponse(403, "User not found.")
	}
	hashPassword, _ := HashPassword(update.Password)
	user.Password = hashPassword
	var tempUser User
	Conn.Model(&tempUser).Where("email = ?", user.Email).Update("password", user.Password)
	SendResetPasswordEmail(user)
	return ErrorResponse(200, "Password Changed Successfully")
}

//GetUserFromTokenString get full user information from string
func GetUserFromTokenString(token string) (int, User) {
	var user User
	code, tokenString := SplitToken(token)
	if code != 200 {
		return 401, user
	}
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(beego.AppConfig.String("jwtkey")), nil
	})
	if err != nil {
		return 401, user
	}
	var email string
	for key, val := range claims {
		if key == "email" {
			email = val.(string)
		}
	}
	if email != "" {
		user, err = GetUserDataEmail(email)
		if err != nil {
			return 401, user
		}
	}
	return 200, user
}

//GetAllUsersFromDepartment retrieves all user data from db bases on department
func GetAllUsersFromDepartment(departmentID string) interface{} {
	var allUsersArray []User
	Conn.Where("department_id = ?", departmentID).Find(&allUsersArray)
	response := ValidResponse(200, allUsersArray, "Success")
	return response
}

//GetAllUsers gets a list of all users on the system
func GetAllUsers() []User {
	var users []User
	if getAllUsers := Conn.Find(&users); getAllUsers.Error != nil {
		return nil
	}
	return users
}

//GetDataFromIDString retrieves user data from ID string
func GetDataFromIDString(id string) (User, error) {
	var u User
	if getUserData := Conn.Where("id = ?", id).Find(&u); getUserData.Error != nil {
		return u, getUserData.Error
	}
	return u, nil
}
