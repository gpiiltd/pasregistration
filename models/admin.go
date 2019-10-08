package models

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

//GetAllTeamLead gets all front team leads in the system
func GetAllTeamLead() interface{} {
	var teamLeads []Roles
	if getAll := Conn.Where("code = 66").Find(&teamLeads); getAll.Error != nil {
		return ErrorResponse(401, "Unable to get front desk officers")
	}
	var userArray []User
	Conn.Find(&userArray)
	var allTeamLeadArray []User

	for _, user := range userArray {
		for _, teamLead := range teamLeads {
			if teamLead.UserID == user.ID {
				allTeamLeadArray = append(allTeamLeadArray, user)
			}
		}
	}
	return ValidResponse(200, allTeamLeadArray, "success")
}

//AddFrontDeskOfficer adds a front desk user to the system
func AddFrontDeskOfficer(frontDesk User) interface{} {
	isFrontDeskOfficer := IsFrontDesk(frontDesk)
	if isFrontDeskOfficer == true {
		return ErrorResponse(401, "User already a front desk officer")
	}
	var role Roles
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
		return ErrorResponse(401, "User cannot have more that 1 team")
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

//DeleteTeamLeadOfficer deletes team lead from the system
func DeleteTeamLeadOfficer(uid string) interface{} {
	var teamLead User
	teamLead, err := GetDataFromIDString(uid)
	if err != nil {
		return ErrorResponse(403, "User does not exist")
	}
	IsTeamLead := IsTeamLead(teamLead)
	if IsTeamLead == false {
		return ErrorResponse(403, "User is not a Team Lead")
	}
	if deleteRole := Conn.Where("user_id = ?", uid).Delete(&Roles{}); deleteRole.Error != nil {
		return ErrorResponse(401, "Unable to delete Team Lead record")
	}
	return ValidResponse(200, "Delete Successful", "success")
}
