package main

import (
	"pasregistration/controllers"
	_ "pasregistration/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
)

func main() {
	logs.SetLogger(logs.AdapterFile, `{"filename": "./apilogs/apicalls.log", "level": 7, "maxlines": 0, "maxsize":0,"daily":true,"maxdays":10,"color":true}`)
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"https://*"},
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Access-control-allow-origin", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	beego.InsertFilter("/v1/*", beego.BeforeRouter, controllers.ValidateToken)
	// beego.InsertFilter("/v1/admin/*", beego.BeforeRouter, controllers.ValidateAdmin)
	beego.Run()
}
