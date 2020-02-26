package models

import (
	"encoding/json"
	"net/http"
	"pasregistration/controllers/mailer"

	"github.com/astaxie/beego"
)

//GetAllFrontDeskOfficer gets all front desk officers in the system
func GetAllFrontDeskOfficer() interface{} {
	var frontDeskOfficers []Roles
	if getAll := Conn.Where("code = 88").Find(&frontDeskOfficers); getAll.Error != nil {
		return ErrorResponse(401, "Unable to get front desk officers")
	}
	var userArray []User
	var u User
	for _, role := range frontDeskOfficers {
		u.ID = role.UserID
		if getUser := Conn.Where("id = ?", role.UserID).Find(&u); getUser.Error != nil {
			return ErrorResponse(401, "Error get user information from role email")
		}
		userArray = append(userArray, u)
	}
	return ValidResponse(200, userArray, "success")
}

//GetTeamLeads gets all team lead on the system
func GetTeamLeads() interface{} {
	var teamLeads []Roles
	if getAll := Conn.Where("code = 66").Find(&teamLeads); getAll.Error != nil {
		return ErrorResponse(401, "Unable to all Team Lead")
	}
	var userArray []User
	var u User
	for _, role := range teamLeads {
		u.ID = role.UserID
		if getUser := Conn.Where("id = ?", role.UserID).Find(&u); getUser.Error != nil {
			return ErrorResponse(401, "Error get user information from role email")
		}
		userArray = append(userArray, u)
	}
	return ValidResponse(200, userArray, "success")
}

//GetAllHROs gets all HR officer on the system
func GetAllHROs() interface{} {
	var allHRO []Roles
	if getAll := Conn.Where("code = 77").Find(&allHRO); getAll.Error != nil {
		return ErrorResponse(401, "Unable to all Team Lead")
	}
	var userArray []User
	var u User
	for _, role := range allHRO {
		u.ID = role.UserID
		if getUser := Conn.Where("id = ?", role.UserID).Find(&u); getUser.Error != nil {
			return ErrorResponse(401, "Error get user information from role email")
		}
		userArray = append(userArray, u)
	}
	return ValidResponse(200, userArray, "success")
}

//AddFrontDeskOfficer adds a front desk user to the system
func AddFrontDeskOfficer(frontDesk User) interface{} {
	isFrontDeskOfficer := IsFrontDesk(frontDesk)
	if isFrontDeskOfficer == true {
		return ErrorResponse(401, "User already a front desk officer")
	}
	var role Roles
	Conn.Where("user_id = ? AND code = ?", frontDesk.ID, 88).Delete(&Roles{})
	role.Code = 88
	role.Role = "Front Desk Officer"
	role.User = frontDesk.FullName
	role.UserID = frontDesk.ID
	if createRole := Conn.Create(&role); createRole.Error != nil {
		return ErrorResponse(403, "Unable to create User Role")
	}
	return ValidResponse(200, "Successfully added front desk officer", "success")
}

//AddTeamLead adds a team lead to the system
func AddTeamLead(teamLead User) interface{} {
	isTeamLead := IsTeamLead(teamLead)
	if isTeamLead == true {
		return ErrorResponse(401, "User already a team Lead")
	}
	var role Roles
	role.Code = 66
	role.Role = "PAS Team Lead"
	role.User = teamLead.FullName
	role.UserID = teamLead.ID
	if createRole := Conn.Create(&role); createRole.Error != nil {
		return ErrorResponse(403, "Unable to add a team lead")
	}
	return ValidResponse(200, "Successfully added team lead", "success")
}

//DeleteFrontDeskOfficer delete front desk officer
func DeleteFrontDeskOfficer(uid string) interface{} {
	var frontDesk User
	frontDesk, err := GetDataFromIDString(uid)
	if err != nil {
		return ErrorResponse(403, "User does not exist")
	}
	isFrontDeskOfficer := IsFrontDesk(frontDesk)
	if isFrontDeskOfficer == false {
		return ErrorResponse(403, "User is not a Front Desk Officer")
	}
	if deleteRole := Conn.Where("user_id = ?", uid).Delete(&Roles{}); deleteRole.Error != nil {
		return ErrorResponse(401, "Unable to delete Front Desk record")
	}
	return ValidResponse(200, "Delete Successful", "success")
}

//ValidationResponseData holds data that needs a true or false response
type ValidationResponseData struct {
	Code int  `json:"code"`
	Body bool `json:"body"`
}

//DeleteTeamLeadOfficer deletes team lead from the system
func DeleteTeamLeadOfficer(uid string, pasTokenString string) interface{} {
	var teamLead User
	teamLead, err := GetDataFromIDString(uid)
	if err != nil {
		return ErrorResponse(403, "User does not exist")
	}
	IsTeamLead := IsTeamLead(teamLead)
	if IsTeamLead == false {
		return ErrorResponse(403, "User is not a Team Lead")
	}
	// roleCode := 66
	client := &http.Client{}
	req, _ := http.NewRequest("GET", beego.AppConfig.String("pasapi")+"team/verifi/"+uid, nil)
	req.Header.Set("authorization", pasTokenString)
	response, _ := client.Do(req)

	var responseBody ValidationResponseData
	json.NewDecoder(response.Body).Decode(&responseBody)

	if responseBody.Body == true {
		return ErrorResponse(403, "Team Lead currently has an active team. Delete team before deleting lead")
	}
	roleCode := 66
	if deleteRole := Conn.Where("user_id = ? AND code = ?", uid, roleCode).Delete(&Roles{}); deleteRole.Error != nil {
		return ErrorResponse(401, "Unable to delete Team Lead record")
	}
	return ValidResponse(200, "Delete Successful", "success")
}

//AddHROfficer adds a new HR officer to the system
func AddHROfficer(hrOfficer User) interface{} {
	isAdmin := IsAdmin(hrOfficer)
	if isAdmin == true {
		return ErrorResponse(401, "User already has a classified role.")
	}
	var role Roles
	role.Code = 77
	role.Role = "HR Officer"
	role.User = hrOfficer.FullName
	role.UserID = hrOfficer.ID
	if createRole := Conn.Create(&role); createRole.Error != nil {
		return ErrorResponse(403, "Unable to add a HR Officer")
	}

	var teamLead Roles
	teamLead.Code = 66
	teamLead.Role = "PAS Team Lead"
	teamLead.User = hrOfficer.FullName
	teamLead.UserID = hrOfficer.ID
	if createRole := Conn.Create(&teamLead); createRole.Error != nil {
		return ErrorResponse(403, "Unable to add a team lead")
	}

	go SendTeamLeadMail(hrOfficer)

	return ValidResponse(200, "Successfully added HR Officer", "success")
}

//SendTeamLeadMail send email confirming users as the new Team Lead.
func SendTeamLeadMail(u User) {
	path := beego.AppConfig.String("mailtemplatepath") + "added-role.html"
	mailSubject := "New Role Added"
	newRequest := mailer.NewRequest(u.Email, mailSubject)
	data := mailer.Data{}
	data.User = u.FullName
	data.Role = "a Team Lead"
	data.Link = beego.AppConfig.String("loginpage")

	go newRequest.Send(path, data)

	return
}

//DeleteHROfficer deletes an HR officer from the system
func DeleteHROfficer(uid string) interface{} {
	var HRO User
	HRO, err := GetDataFromIDString(uid)
	if err != nil {
		return ErrorResponse(403, "User does not exist")
	}
	isHRO := IsHRO(HRO)
	if isHRO == false {
		return ErrorResponse(403, "User is not an HR Officer")
	}
	roleCode := 77
	if deleteRole := Conn.Where("user_id = ? AND code = ?", uid, roleCode).Delete(&Roles{}); deleteRole.Error != nil {
		return ErrorResponse(401, deleteRole.Error.Error())
	}
	return ValidResponse(200, "Delete Successful", "success")
}

//AddVMSAdminOfficer adds a new vms admin to the system
func AddVMSAdminOfficer(vmsAdmin User) interface{} {
	isVMSAdmin := IsVMSAdmin(vmsAdmin)
	if isVMSAdmin == true {
		return ErrorResponse(401, "User cannot have more that 1 team")
	}
	var role Roles
	role.Code = 44
	role.Role = "VMS Admin"
	role.User = vmsAdmin.FullName
	role.UserID = vmsAdmin.ID
	if createRole := Conn.Create(&role); createRole.Error != nil {
		return ErrorResponse(403, "Unable to add vms admin")
	}
	return ValidResponse(200, "Successfully added team lead", "success")
}

//AddTaskAdminOfficer adds a new task admin to the system
func AddTaskAdminOfficer(taskAdmin User) interface{} {
	isTaskAdmin := IsTaskAdmin(taskAdmin)
	if isTaskAdmin == true {
		return ValidResponse(401, "User already a task admin", "error")
	}
	var role Roles
	role.Code = UserRoles.TaskAdmin
	role.Role = "Task Admin"
	role.User = taskAdmin.FullName
	role.UserID = taskAdmin.ID
	if createRole := Conn.Create(&role); createRole.Error != nil {
		return ValidResponse(403, "Unable to add task admin", "error")
	}
	return ValidResponse(200, "Successfully added team lead", "success")
}

//GetAllVMSAdmin gets all vms admin on the system
func GetAllVMSAdmin() interface{} {
	var allVMSAdmin []Roles
	if getAll := Conn.Where("code = 44").Find(&allVMSAdmin); getAll.Error != nil {
		return ErrorResponse(401, "Unable to get all vms Admin")
	}
	var userArray []User
	var u User
	for _, role := range allVMSAdmin {
		u.ID = role.UserID
		if getUser := Conn.Where("id = ?", role.UserID).Find(&u); getUser.Error != nil {
			return ErrorResponse(401, "Error get user information from role email")
		}
		userArray = append(userArray, u)
	}
	return ValidResponse(200, userArray, "success")
}

//GetAllSystemTaskAdmin retrrieves a list of all task admin from the system
func GetAllSystemTaskAdmin() ([]User, error) {
	var allTaskAdmin []Roles
	if getAlladmin := Conn.Where("code = ?", UserRoles.TaskAdmin).Find(&allTaskAdmin); getAlladmin.Error != nil {
		return []User{}, getAlladmin.Error
	}
	var u User
	var userArray []User
	for _, role := range allTaskAdmin {
		u.ID = role.UserID
		if getUser := Conn.Where("id = ?", role.UserID).Find(&u); getUser.Error != nil {
			return userArray, getUser.Error
		}
		userArray = append(userArray, u)
	}
	return userArray, nil
}

//DeletevmsAdmin deletes a vms admin from the system
func DeletevmsAdmin(uid string) interface{} {
	var admin User
	admin, err := GetDataFromIDString(uid)
	if err != nil {
		return ErrorResponse(403, "User does not exist")
	}
	isVMSAdmin := IsVMSAdmin(admin)
	if isVMSAdmin == false {
		return ErrorResponse(403, "User is not a an Admin")
	}
	if deleteRole := Conn.Where("user_id = ?", uid).Delete(&Roles{}); deleteRole.Error != nil {
		return ErrorResponse(401, "Unable to delete Admin record")
	}
	return ValidResponse(200, "Delete Successful", "success")
}
