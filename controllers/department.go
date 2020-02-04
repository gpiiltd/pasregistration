package controllers

import (
	"encoding/json"
	"pasregistration/models"

	"github.com/astaxie/beego"
)

//DepartmentController handles all data about departments
type DepartmentController struct {
	beego.Controller
}

//SetDepartmentHead sets a department head for a user.
// @Title SetDepartmentHead
// @Description set department head.
// @Param	body		body 	models.Departments	true		"body for user content"
// @Success 200 {int} models.ValidResponse
// @Failure 403 body is empty
// @router /head/ [post]
func (d *DepartmentController) SetDepartmentHead() {
	var departmentDetails models.Departments
	err := json.Unmarshal(d.Ctx.Input.RequestBody, &departmentDetails)
	if err != nil {
		d.Data["json"] = models.ErrorResponse(405, err.Error())
		d.ServeJSON()
		return
	}

	d.Data["json"] = models.MakeHeadDepartment(departmentDetails)
	d.ServeJSON()
}
