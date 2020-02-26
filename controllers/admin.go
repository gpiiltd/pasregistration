package controllers

import (
	"encoding/json"
	"log"
	"pasregistration/models"

	"github.com/astaxie/beego"
)

//AdminController handles all functions that belongs explicitly to the admin
type AdminController struct {
	beego.Controller
}

//GetAllFrontDeskOfficer gets all front desk officers
// @Title Create
// @Description gets the list of all front desk officers
// @Success 200 {object} []models.User
// @Failure 403 body is empty
// @router /frontdesk/ [get]
func (a *AdminController) GetAllFrontDeskOfficer() {
	a.Data["json"] = models.GetAllFrontDeskOfficer()
	a.ServeJSON()
}

//GetAllTeamLead gets all team lead in the system
// @Title GetAllTeamLead
// @Description gets the list of all team leads on the system
// @Success 200 {object} []models.User
// @Failure 403 body is empty
// @router /teamlead/ [get]
func (a *AdminController) GetAllTeamLead() {
	a.Data["json"] = models.GetTeamLeads()
	a.ServeJSON()
}

//AddFrontDeskOfficer add a new front desk user
// @Title Create
// @Description adds a new front desk officer using the user ID
// @Param	visitid		path 	string	true		"the id of the user you want to make a front desk officer"
// @Success 200 {string} id of the user
// @Failure 403 body is empty
// @router /frontdesk/:id [post]
func (a *AdminController) AddFrontDeskOfficer() {
	var frontDesk models.User
	frontDeskID := a.GetString(":id")
	frontDesk, err := models.GetDataFromIDString(frontDeskID)
	if err != nil {
		a.Data["json"] = models.ErrorResponse(404, "Font Desk Officer data does not exist")
		a.ServeJSON()
		return
	}
	updateFrontDesk := models.AddFrontDeskOfficer(frontDesk)
	a.Data["json"] = updateFrontDesk
	a.ServeJSON()
}

//DeleteFrontDeskOfficer deletes a front desk user
// @Title Delete
// @Description deletes a front desk officer using the user ID
// @Param	visitid		path 	string	true		"the id of the user you want to delete"
// @Success 200 {string} id of the user
// @Failure 403 body is empty
// @router /frontdesk/:id [delete]
func (a *AdminController) DeleteFrontDeskOfficer() {
	frontDeskID := a.GetString(":id")
	updateFrontDesk := models.DeleteFrontDeskOfficer(frontDeskID)
	a.Data["json"] = updateFrontDesk
	a.ServeJSON()
}

//AddNewTeamLead adds a new team lead to the system
// @Title AddNewTeamLead
// @Description adds a new new team lead to the system
// @Param	visitid		path 	string	true		"the id of the user you want to make a front desk officer"
// @Success 200 {string} id of the user
// @Failure 403 body is empty
// @router /teamlead/:id [post]
func (a *AdminController) AddNewTeamLead() {
	var teamLead models.User
	teamLeadID := a.GetString(":id")
	teamLead, err := models.GetDataFromIDString(teamLeadID)
	if err != nil {
		a.Data["json"] = models.ErrorResponse(404, "Team Lead data does not exist")
		a.ServeJSON()
		return
	}
	updateFrontDesk := models.AddTeamLead(teamLead)
	a.Data["json"] = updateFrontDesk
	a.ServeJSON()
}

//DeleteTeamLead deletes a team lead
// @Title Delete
// @Description deletes a team lead using the user ID
// @Param	visitid		path 	string	true		"the id of the user you want to delete"
// @Success 200 {string} id of the user
// @Failure 403 body is empty
// @router /teamlead/:id [delete]
func (a *AdminController) DeleteTeamLead() {
	type pasStruct struct {
		PasString string `json:"pastokenstring"`
	}
	var tokenString pasStruct
	err := json.Unmarshal(a.Ctx.Input.RequestBody, &tokenString)
	if err != nil {
		log.Println(err.Error())
		return
	}
	teamLeadID := a.GetString(":id")
	a.Data["json"] = models.DeleteTeamLeadOfficer(teamLeadID, tokenString.PasString)
	a.ServeJSON()
}

//AddNewHROfficer adds a new HR officer to the system
// @Title AddNewHROfficer
// @Description adds a new hr officer to the system
// @Param	userid		path 	string	true		"the id of the user you want to make an HR officer"
// @Success 200 {string} id of the user
// @Failure 403 body is empty
// @router /hro/:id [post]
func (a *AdminController) AddNewHROfficer() {
	var HRO models.User
	HROID := a.GetString(":id")
	HRO, err := models.GetDataFromIDString(HROID)
	if err != nil {
		a.Data["json"] = models.ErrorResponse(404, "Team Lead data does not exist")
		a.ServeJSON()
		return
	}
	updateHROfficer := models.AddHROfficer(HRO)
	a.Data["json"] = updateHROfficer
	a.ServeJSON()
}

//GetAllHRO gets all HR officers in the system
// @Title GetAllHRO
// @Description gets the list of all HR officers on the system
// @Success 200 {object} []models.User
// @Failure 403 body is empty
// @router /hro/ [get]
func (a *AdminController) GetAllHRO() {
	a.Data["json"] = models.GetAllHROs()
	a.ServeJSON()
}

//DeleteHRO deletes an HR officer
// @Title DeleteHRO
// @Description deletes an HRO using the user ID
// @Param	userid		path 	string	true		"the id of the user you want to delete"
// @Success 200 {string} id of the user
// @Failure 403 body is empty
// @router /hro/:id [delete]
func (a *AdminController) DeleteHRO() {
	HROID := a.GetString(":id")
	a.Data["json"] = models.DeleteHROfficer(HROID)
	a.ServeJSON()
}

//AddVMSAdmin adds a new VMS Admin to the system
// @Title AddVMSAdmin
// @Description adds a new vms officer to the system
// @Param	userid		path 	string	true		"the id of the user you want to make a vms officer"
// @Success 200 {string} id of the user
// @Failure 403 body is empty
// @router /vmsadmin/:id [post]
func (a *AdminController) AddVMSAdmin() {
	var VMSAdmin models.User
	vmsAdminID := a.GetString(":id")
	VMSAdmin, err := models.GetDataFromIDString(vmsAdminID)
	if err != nil {
		a.Data["json"] = models.ErrorResponse(404, "User data does not exist")
		a.ServeJSON()
		return
	}
	addVMSAdmin := models.AddVMSAdminOfficer(VMSAdmin)
	a.Data["json"] = addVMSAdmin
	a.ServeJSON()
}

//GetAllVMSAdmin gets all vms admin in the system
// @Title GetAllVMSAdmin
// @Description gets the list of all VMS admin officers on the system
// @Success 200 {object} []models.User
// @Failure 403 body is empty
// @router /vmsadmin/ [get]
func (a *AdminController) GetAllVMSAdmin() {
	a.Data["json"] = models.GetAllVMSAdmin()
	a.ServeJSON()
}

//DeleteVMSAdmin deletes a vms admin
// @Title DeleteVMSAdmin
// @Description deletes a vms admin using the user ID
// @Param	userid		path 	string	true		"the id of the user you want to delete"
// @Success 200 {string} id of the user
// @Failure 403 body is empty
// @router /vmsadmin/:id [delete]
func (a *AdminController) DeleteVMSAdmin() {
	adminID := a.GetString(":id")
	a.Data["json"] = models.DeletevmsAdmin(adminID)
	a.ServeJSON()
}

//AddTaskAdmin adds a new daily task Admin to the system
// @Title AddTaskAdmin
// @Description adds a new daily task admin to the system
// @Param	userid		path 	string	true		"the id of the user you want to make an admin."
// @Success 200 {string} id of the user
// @Failure 403 body is empty
// @router /taskadmin/:userid [post]
func (a *AdminController) AddTaskAdmin() {
	var taskAdmin models.User
	taskAdminID := a.GetString(":userid")
	taskAdmin, err := models.GetDataFromIDString(taskAdminID)
	if err != nil {
		a.Data["json"] = models.ErrorResponse(404, "User data does not exist")
		a.ServeJSON()
		return
	}
	addVMSAdmin := models.AddTaskAdminOfficer(taskAdmin)
	a.Data["json"] = addVMSAdmin
	a.ServeJSON()
}

//GetAllTaskAdmin gets all task admin in the system
// @Title GetAllTaskAdmin
// @Description gets the list of all task admin officers on the system
// @Success 200 {object} []models.User
// @Failure 403 body is empty
// @router /taskadmin/ [get]
func (a *AdminController) GetAllTaskAdmin() {
	allTaskAdmin, err := models.GetAllSystemTaskAdmin()
	if err != nil {
		a.Data["json"] = models.ErrorResponse(403, "Unable to get all system admin ")
		a.ServeJSON()
		return
	}
	a.Data["json"] = models.ValidResponse(200, allTaskAdmin, "success")
	a.ServeJSON()
}
