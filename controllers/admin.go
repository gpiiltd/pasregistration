package controllers

import (
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
	teamLeadID := a.GetString(":id")
	a.Data["json"] = models.DeleteTeamLeadOfficer(teamLeadID)
	a.ServeJSON()
}
