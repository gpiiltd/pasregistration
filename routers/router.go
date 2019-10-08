//Package routers handles routing
// @APIVersion 1.0.0
// @Title Century Group API Core
// @Description This handles all user and role management concerning century group.
// @Contact endy.apina@my-gpi.io
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"pasregistration/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/token",
			beego.NSInclude(
				&controllers.TokenController{},
			),
		),
		beego.NSNamespace("/validate",
			beego.NSInclude(
				&controllers.ValidateController{},
			),
		),
		beego.NSNamespace("/utility",
			beego.NSInclude(
				&controllers.UtilityController{},
			),
		),
		beego.NSNamespace("/upload",
			beego.NSInclude(
				&controllers.UploadController{},
			),
		),
		beego.NSNamespace("/admin",
			beego.NSInclude(
				&controllers.AdminController{},
			),
		),
	)
	beego.AddNamespace(ns)
	beego.SetStaticPath("/files", "files/")
}
