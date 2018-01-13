package main

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/wangsongyan/genealogy/controllers"
	"github.com/wangsongyan/genealogy/models"
)

func init() {
	//logger, err := seelog.LoggerFromConfigAsFile("conf/seelog.xml")
	//if err != nil {
	//	seelog.Critical("err parsing config log file", err)
	//	return
	//}
	//seelog.ReplaceLogger(logger)
	//defer seelog.Flush()
}

func main() {
	db, err := models.InitDB()
	if err != nil {
		seelog.Critical("err open databases: ", err)
		return
	}
	defer db.Close()

	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	setTemplate(router)

	router.Static("/static", "./static")

	router.GET("/", controllers.IndexGet)
	router.GET("/index", controllers.IndexGet)
	router.GET("/c/:id", controllers.CoupleGet)
	router.POST("/c/:id", controllers.CouplePost)
	router.GET("/p/:id", controllers.GetFamilyTreeByPeopleId)

	router.GET("/new_people", controllers.PeopleNew)
	router.POST("/new_people", controllers.PeopleCreate)
	//router.GET("/people/:id/edit", controllers.PeopleEdit)
	//router.POST("/people/:id/edit", controllers.PeopleUpdate)
	//router.POST("/people/:id/delete", controllers.PeopleDelete)

	router.GET("/new_couple", controllers.CoupleNew)
	router.POST("/new_couple", controllers.CoupleCreate)

	router.Run(":8080")
}

func setTemplate(engine *gin.Engine) {

	//funcMap := template.FuncMap{
	//	"dateFormat": helpers.DateFormat,
	//	"substring":  helpers.Substring,
	//	"isOdd":      helpers.IsOdd,
	//	"isEven":     helpers.IsEven,
	//	"truncate":   helpers.Truncate,
	//	"add":        helpers.Add,
	//}
	//
	//engine.SetFuncMap(funcMap)
	engine.LoadHTMLGlob("views/*")
}
