package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["pasregistration/controllers:AdminController"] = append(beego.GlobalControllerRouter["pasregistration/controllers:AdminController"],
        beego.ControllerComments{
            Method: "GetAllFrontDeskOfficer",
            Router: `/frontdesk/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pasregistration/controllers:AdminController"] = append(beego.GlobalControllerRouter["pasregistration/controllers:AdminController"],
        beego.ControllerComments{
            Method: "AddFrontDeskOfficer",
            Router: `/frontdesk/:id`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pasregistration/controllers:AdminController"] = append(beego.GlobalControllerRouter["pasregistration/controllers:AdminController"],
        beego.ControllerComments{
            Method: "DeleteFrontDeskOfficer",
            Router: `/frontdesk/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pasregistration/controllers:AdminController"] = append(beego.GlobalControllerRouter["pasregistration/controllers:AdminController"],
        beego.ControllerComments{
            Method: "GetAllTeamLead",
            Router: `/teamlead/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pasregistration/controllers:AdminController"] = append(beego.GlobalControllerRouter["pasregistration/controllers:AdminController"],
        beego.ControllerComments{
            Method: "AddNewTeamLead",
            Router: `/teamlead/:id`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pasregistration/controllers:AdminController"] = append(beego.GlobalControllerRouter["pasregistration/controllers:AdminController"],
        beego.ControllerComments{
            Method: "DeleteTeamLead",
            Router: `/teamlead/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pasregistration/controllers:TokenController"] = append(beego.GlobalControllerRouter["pasregistration/controllers:TokenController"],
        beego.ControllerComments{
            Method: "ValidateAttachedToken",
            Router: `/validate`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pasregistration/controllers:UploadController"] = append(beego.GlobalControllerRouter["pasregistration/controllers:UploadController"],
        beego.ControllerComments{
            Method: "UploadImage",
            Router: `/image`,
            AllowHTTPMethods: []string{"POST"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pasregistration/controllers:UserController"] = append(beego.GlobalControllerRouter["pasregistration/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetAllUsers",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pasregistration/controllers:UserController"] = append(beego.GlobalControllerRouter["pasregistration/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetAllFromDepartment",
            Router: `/:departmentid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pasregistration/controllers:UserController"] = append(beego.GlobalControllerRouter["pasregistration/controllers:UserController"],
        beego.ControllerComments{
            Method: "Login",
            Router: `/login`,
            AllowHTTPMethods: []string{"POST"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pasregistration/controllers:UserController"] = append(beego.GlobalControllerRouter["pasregistration/controllers:UserController"],
        beego.ControllerComments{
            Method: "RecoverPassword",
            Router: `/password/recover/:email`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pasregistration/controllers:UserController"] = append(beego.GlobalControllerRouter["pasregistration/controllers:UserController"],
        beego.ControllerComments{
            Method: "PasswordReset",
            Router: `/password/reset`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pasregistration/controllers:UserController"] = append(beego.GlobalControllerRouter["pasregistration/controllers:UserController"],
        beego.ControllerComments{
            Method: "PasswordUpdate",
            Router: `/password/update`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pasregistration/controllers:UserController"] = append(beego.GlobalControllerRouter["pasregistration/controllers:UserController"],
        beego.ControllerComments{
            Method: "RegisterUsers",
            Router: `/register`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pasregistration/controllers:UserController"] = append(beego.GlobalControllerRouter["pasregistration/controllers:UserController"],
        beego.ControllerComments{
            Method: "UpdateProfile",
            Router: `/update/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pasregistration/controllers:UserController"] = append(beego.GlobalControllerRouter["pasregistration/controllers:UserController"],
        beego.ControllerComments{
            Method: "UpdateProfilePicture",
            Router: `/update/avatar/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pasregistration/controllers:UtilityController"] = append(beego.GlobalControllerRouter["pasregistration/controllers:UtilityController"],
        beego.ControllerComments{
            Method: "GetDepartmentList",
            Router: `/departments/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pasregistration/controllers:UtilityController"] = append(beego.GlobalControllerRouter["pasregistration/controllers:UtilityController"],
        beego.ControllerComments{
            Method: "GetSubDepartmentList",
            Router: `/departments/:subsidiaryid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pasregistration/controllers:UtilityController"] = append(beego.GlobalControllerRouter["pasregistration/controllers:UtilityController"],
        beego.ControllerComments{
            Method: "GetSubsidiaryList",
            Router: `/subsidiary/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pasregistration/controllers:ValidateController"] = append(beego.GlobalControllerRouter["pasregistration/controllers:ValidateController"],
        beego.ControllerComments{
            Method: "ValidateResetPasswordCode",
            Router: `/password/code`,
            AllowHTTPMethods: []string{"POST"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pasregistration/controllers:ValidateController"] = append(beego.GlobalControllerRouter["pasregistration/controllers:ValidateController"],
        beego.ControllerComments{
            Method: "ValidateTeamLead",
            Router: `/teamlead/`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
