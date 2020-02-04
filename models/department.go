package models

//MakeHeadDepartment assigns a user as the head of a department
func MakeHeadDepartment(departmentInfo Departments) interface{} {
	var userDepartment Departments
	status, userDepartment := IsUserHead(departmentInfo.HeadID)
	if status == true {
		return ErrorResponse(403, "User is already the head of "+userDepartment.Department)
	}

	var departmentDetails Departments
	departmentDetails, err := GetDepartmentFromID(departmentInfo.ID)
	if err != nil {
		return ErrorResponse(403, "Department does not exist")
	}

	if departmentInfo.HeadID == 0 {
		return ErrorResponse(403, "Empty Head ID")
	}

	var departmentHead User
	Conn.Where("id = ?", departmentInfo.HeadID).Find(&departmentHead)

	// Conn.Create(&departmentDetails)
	Conn.Model(&departmentDetails).Where("id = ?", departmentDetails.ID).Updates(Departments{HeadID: departmentHead.ID, Head: departmentHead.FullName})

	return ValidResponse(200, "Success", "success")
}

//IsUserHead checks if a user is already a department head
func IsUserHead(userID uint64) (bool, Departments) {
	var departmentInfo Departments
	if findDepartmentHead := Conn.Where("head_id = ?", userID).Find(&departmentInfo); findDepartmentHead.Error != nil {
		return false, departmentInfo
	}

	return true, departmentInfo
}

//GetDepartmentFromID retrieves a department from department ID
func GetDepartmentFromID(departmentID uint64) (Departments, error) {
	var departmentInfo Departments
	if findDepartmentHead := Conn.Where("id = ?", departmentID).Find(&departmentInfo); findDepartmentHead.Error != nil {
		return departmentInfo, findDepartmentHead.Error
	}

	return departmentInfo, nil
}
