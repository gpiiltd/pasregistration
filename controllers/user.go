package controllers

import (
	"encoding/json"
	"pasregistration/models"

	"github.com/astaxie/beego"
)

//UserController handles all about users.
// Handles all User controls
type UserController struct {
	beego.Controller
}

//RegisterUsers handles user registration.
// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router /register [post]
func (u *UserController) RegisterUsers() {
	var user models.User
	if err := u.ParseForm(&user); err != nil {
		u.Data["json"] = models.ErrorResponse(403, err.Error())
		u.ServeJSON()
	}
	u.Data["json"] = models.RegisterUsers(user)
	u.ServeJSON()
}

//Login handles user login
// @Title Login
// @Description Logs user into the system
// @Param	username		query 	string	true		"The username for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {object} models.SuccessData
// @Failure 403 user not exist
// @router /login [POST]
func (u *UserController) Login() {
	email := u.GetString("username")
	password := u.GetString("password")

	u.Data["json"] = models.Login(email, password)
	u.ServeJSON()
}

//RecoverPassword recovers lost password
// @Title Recover password
// @Description recover lost password
// @Success 200 {string} "success"
// @router /password/recover/:email [get]
func (u *UserController) RecoverPassword() {
	email := u.GetString(":email")
	u.Data["json"] = models.PasswordRecovery(email)
	u.ServeJSON()
}

//PasswordReset recovers resets password
// @Title Recover password
// @Description recover lost password
// @Success 200 {string} "success"
// @router /password/reset [post]
func (u *UserController) PasswordReset() {
	var recoveryData models.ResetPasswordData
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &recoveryData)
	if err != nil {
		u.Data["json"] = models.ErrorResponse(405, "Method Not Allowed")
		u.ServeJSON()
		return
	}
	u.Data["json"] = models.ResetPassword(recoveryData)
	u.ServeJSON()
}

//PasswordUpdate changes user password
// @Title PasswordUpdate
// @Description update user password
// @Success 200 {string} "success"
// @router /password/update [post]
func (u *UserController) PasswordUpdate() {
	var user models.User
	resCode, user := models.GetUserFromTokenString(u.Ctx.Input.Header("authorization"))
	if resCode != 200 {
		u.Data["json"] = models.ErrorResponse(403, "Unable to get token string")
		u.ServeJSON()
		return
	}
	var userProfile models.User
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &userProfile)
	if err != nil {
		u.Data["json"] = models.ErrorResponse(405, "Method Not Allowed")
		u.ServeJSON()
		return
	}
	u.Data["json"] = models.ChangePassword(userProfile, user)
	u.ServeJSON()
}

//GetAllUsers gets a list of all users
// @Title GetAllUsers
// @Description get all users on the system
// @Success 200 {object} models.ValidResponse
// @Failure 403 :uid is empty
// @router / [get]
func (u *UserController) GetAllUsers() {
	var allUser []models.User
	allUser = models.GetAllUsers()
	u.Data["json"] = models.ValidResponse(200, allUser, "success")
	u.ServeJSON()
}

//UpdateProfile updates a user profile
// @Title UpdateProfile
// @Description update the user profile
// @Param	uid		path 	string	true		"The uid you want to update"
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router /update/ [post]
func (u *UserController) UpdateProfile() {
	token := u.Ctx.Input.Header("authorization")
	var user models.User
	resCode, user := models.GetUserFromTokenString(token)
	if resCode != 200 {
		u.Data["json"] = models.ErrorResponse(403, "Unable to get token string")
		u.ServeJSON()
		return
	}
	var update models.User
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &update)
	if err != nil {
		u.Data["json"] = models.ErrorResponse(405, err.Error())
		u.ServeJSON()
		return
	}
	u.Data["json"] = models.UpdateProfile(update, user)
	u.ServeJSON()
}

//UpdateProfilePicture updates a user profile avatar
// @Title UpdateProfilePicture
// @Description update the user profile avatar
// @Param	uid		path 	string	true		"The uid you want to update"
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router /update/avatar/ [post]
func (u *UserController) UpdateProfilePicture() {
	token := u.Ctx.Input.Header("authorization")
	var user models.User
	resCode, user := models.GetUserFromTokenString(token)
	if resCode != 200 {
		u.Data["json"] = models.ErrorResponse(403, "Unable to get token string")
		u.ServeJSON()
		return
	}
	var update models.User
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &update)
	if err != nil {
		u.Data["json"] = models.ErrorResponse(405, err.Error())
		u.ServeJSON()
		return
	}
	u.Data["json"] = models.UpdateProfileImage(update, user)
	u.ServeJSON()
}

//GetAllFromDepartment gets all user data from db using a particular department
// @Title GetAll
// @Description get all Users
// @Param	body		body 	string	true		"departmentID"
// @Success 200 {object} []models.User
// @router /:departmentid [get]
func (u *UserController) GetAllFromDepartment() {
	department := u.GetString(":department")
	users := models.GetAllUsersFromDepartment(department)
	u.Data["json"] = users
	u.ServeJSON()
}

// // @Title Delete
// // @Description delete the user
// // @Param	uid		path 	string	true		"The uid you want to delete"
// // @Success 200 {string} delete success!
// // @Failure 403 uid is empty
// // @router /:uid [delete]
// func (u *UserController) Delete() {
// 	uid := u.GetString(":uid")
// 	models.DeleteUser(uid)
// 	u.Data["json"] = "delete success!"
// 	u.ServeJSON()
// }

// // @Title Login
// // @Description Logs user into the system
// // @Param	username		query 	string	true		"The username for login"
// // @Param	password		query 	string	true		"The password for login"
// // @Success 200 {string} login success
// // @Failure 403 user not exist
// // @router /login [get]
// func (u *UserController) Login() {
// 	username := u.GetString("username")
// 	password := u.GetString("password")
// 	if models.Login(username, password) {
// 		u.Data["json"] = "login success"
// 	} else {
// 		u.Data["json"] = "user not exist"
// 	}
// 	u.ServeJSON()
// }

// // @Title logout
// // @Description Logs out current logged in user session
// // @Success 200 {string} logout success
// // @router /logout [get]
// func (u *UserController) Logout() {
// 	u.Data["json"] = "logout success"
// 	u.ServeJSON()
// }
