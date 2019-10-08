package controllers

import (
	"encoding/json"
	"pasregistration/models"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

//ValidateToken validates token
var ValidateToken = func(ctx *context.Context) {
	filter := Filter(ctx)
	if filter == true {
		return
	}
	type unAuthorized struct {
		Code int    `json:"code"`
		Body string `json:"body"`
	}

	token := ctx.Input.Header("authorization")
	validToken := models.ValidToken(token)
	if validToken != true {
		var res unAuthorized
		res.Code = 403
		res.Body = "Unauthorized Connection. Invalid Token"
		ctx.Output.JSON(res, false, false)

		return
	}
	isNull := NullToken(token)
	if isNull == true {
		var res unAuthorized
		res.Code = 403
		res.Body = "Unauthorized Connection. Empty Token String"
		ctx.Output.JSON(res, false, false)

		return
	}
	if token == "" {
		var res unAuthorized
		res.Code = 403
		res.Body = "Unauthorized Connection. Empty Token"
		ctx.Output.JSON(res, false, false)

		return
	}
	isTokenExpired := models.TokenExpire(token)
	if isTokenExpired != true {
		var res unAuthorized
		res.Code = 401
		res.Body = "Token Expired, Kindly Login again."
		ctx.Output.JSON(res, false, false)
	}
	if strings.HasPrefix(ctx.Input.URL(), "/v1/user/validate") {
		return
	}
}

//Filter checks if there are endpoint that shouldn't contain token string
func Filter(ctx *context.Context) bool {
	if strings.HasPrefix(ctx.Input.URL(), "/v1/contact/") {
		return true
	}

	if strings.HasPrefix(ctx.Input.URL(), "/v1/user/login") {
		return true
	}

	if strings.HasPrefix(ctx.Input.URL(), "/v1/user/register") {
		return true
	}

	if strings.HasPrefix(ctx.Input.URL(), "/v1/user/password/recover") {
		return true
	}

	if strings.HasPrefix(ctx.Input.URL(), "/v1/user/password/reset") {
		return true
	}

	if strings.HasPrefix(ctx.Input.URL(), "/v1/validate/password/code") {
		return true
	}

	if strings.HasPrefix(ctx.Input.URL(), "/v1/validate/teamlead/") {
		return true
	}

	if strings.HasPrefix(ctx.Input.URL(), "/v1/utility/subsidiary/") {
		return true
	}

	if strings.HasPrefix(ctx.Input.URL(), "/v1/utility/departments/") {
		return true
	}

	if strings.HasPrefix(ctx.Input.URL(), "/v1/guest/") {
		return true
	}

	if strings.HasPrefix(ctx.Input.URL(), "/v1/contact/") {
		return true
	}
	return false
}

//NullToken checks if token is null
func NullToken(wholeToken string) bool {
	splitString := strings.Split(wholeToken, ",")
	if splitString[1] == "" {
		return true
	}

	return false
}

//TokenController handles all about tokens.
type TokenController struct {
	beego.Controller
}

//ValidateController handles all validation.
type ValidateController struct {
	beego.Controller
}

//UtilityController handles all extra utilities.
type UtilityController struct {
	beego.Controller
}

//UploadController handles uploads of images, documents, csvs etc
type UploadController struct {
	beego.Controller
}

//UploadImage uploads a business document to the back end.
// @Title Accepts Invitation Link
// @Description validates an invitation link and confirms it.
// @Param	body		body 	models.Invitation	true		"A json containing the role {int}, email {string} and code {string}"
// @Success 200 {string} "Invitation Url"
// @router /image [POST]
func (u *UploadController) UploadImage() {
	var user models.User
	code, user := models.GetUserFromTokenString(u.Ctx.Input.Header("authorization"))
	if code != 200 {
		u.Data["json"] = models.ErrorResponse(403, "Invalid token string")
		u.ServeJSON()
		return
	}
	file, header, _ := u.GetFile("image")
	if file == nil {
		u.Data["json"] = models.ErrorResponse(404, "Image not Found")
		u.ServeJSON()
		return
	}
	u.Data["json"] = models.UploadImage(file, header, user)
	u.ServeJSON()
}

//ValidateAttachedToken validates a user token.
// @Title ValidateAttachedToken
// @Description validates a user token and send a true or false response depending on the validity
// @Success 200 {int} models.ValidResponse
// @Failure 403 body is empty
// @router /validate [get]
func (t *TokenController) ValidateAttachedToken() {
	t.Data["json"] = models.ValidateTokenString(t.Ctx.Input.Header("authorization"))
	t.ServeJSON()
}

//ValidateResetPasswordCode validates a password recovery code.
// @Title ValidateResetPasswordCode
// @Description validates a user recovery email and send a true or false response depending on the validity
// @Success 200 {object} models.ValidResponse
// @Failure 403 body is empty
// @router /password/code [POST]
func (v *ValidateController) ValidateResetPasswordCode() {
	var validationInfo models.PasswordRecoveryData
	err := json.Unmarshal(v.Ctx.Input.RequestBody, &validationInfo)
	if err != nil {
		v.Data["json"] = models.ErrorResponse(405, "Method Not Allowed")
		v.ServeJSON()
		return
	}
	v.Data["json"] = models.ValidateRecoveryCode(validationInfo)
	v.ServeJSON()
}

//ValidateTeamLead validates if a user is a teamLead.
// @Title ValidateTeamLead
// @Description validates a user to see if he is a team lead. Returns true or false response
// @Success 200 {object} models.ValidResponse
// @Failure 403 body is empty
// @router /teamlead/ [GET]
func (v *ValidateController) ValidateTeamLead() {
	var user models.User
	err := json.Unmarshal(v.Ctx.Input.RequestBody, &user)
	if err != nil {
		v.Data["json"] = models.ErrorResponse(405, "Method Not Allowed")
		v.ServeJSON()
		return
	}
	v.Data["json"] = models.ValidateTeamLeadUser(user)
	v.ServeJSON()
}

//GetSubsidiaryList retirieves the list of all subsidiaries in the system.
// @Title GetSubsidiaryList
// @Description gets the list of all subsidiaries in the system
// @Success 200 {object} models.ValidResponse
// @Failure 403 body is empty
// @router /subsidiary/ [get]
func (util *UtilityController) GetSubsidiaryList() {
	subsidiaryList := models.GetSubsidiaries()
	util.Data["json"] = models.ValidResponse(200, subsidiaryList, "success")
	util.ServeJSON()
}

//GetDepartmentList retirieves the list of all departments in the system.
// @Title GetDepartmentList
// @Description gets the list of all departments in the system
// @Success 200 {object} models.ValidResponse
// @Failure 403 body is empty
// @router /departments/ [get]
func (util *UtilityController) GetDepartmentList() {
	departmentList := models.GetDeparments()
	util.Data["json"] = models.ValidResponse(200, departmentList, "success")
	util.ServeJSON()
}

//GetSubDepartmentList retirieves the list of all departments in the system.
// @Title GetSubDepartmentList
// @Description gets the list of all departments in a subsidiary
// @Success 200 {object} models.ValidResponse
// @Failure 403 body is empty
// @router /departments/:subsidiaryid [get]
func (util *UtilityController) GetSubDepartmentList() {
	subsidiaryID := util.GetString(":subsidiaryid")
	departmentList := models.GetSubsidiaryDepartments(subsidiaryID)
	util.Data["json"] = models.ValidResponse(200, departmentList, "success")
	util.ServeJSON()
}
