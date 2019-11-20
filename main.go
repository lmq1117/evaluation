// file: main.go

package main

import (
	"evaluation/datamodels"
	//"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	//"github.com/kataras/iris/v12/middleware/basicauth"
	//"time"

	"evaluation/services"
	"evaluation/web/controllers"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

var (
	cookieNameForSessionID = "sessioncookiename1"
	sessManager            = sessions.New(sessions.Config{
		Cookie: cookieNameForSessionID,
		//Expires: 30 * 24 * time.Hour,
	})
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	tmpl := iris.HTML("./web/views", ".html").
		Layout("shared/layout.html").
		Reload(true)
	app.RegisterView(tmpl)
	app.HandleDir("/public", "./web/public")
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().
			GetStringDefault("message", "The page you're looking for doesn't exist"))
		ctx.View("shared/error.html")
	})

	db, err := NewDataBaseEngine()
	if err != nil {
		app.Logger().Fatalf("连接数据库失败: %v", err)
	}
	userService := services.NewUserService(db)

	//authConfig := basicauth.Config{
	//	Users:   map[string]string{"admin": "123456", "admin888": "678910"},
	//	Realm:   "Authorization Required",
	//	Expires: time.Duration(30) * time.Minute,
	//	OnAsk:   nil,
	//}
	//user := mvc.New(app.Party("/user", basicauth.New(authConfig)))
	user := mvc.New(app.Party("/user"))
	user.Register(
		userService,
		sessManager.Start,
	)
	user.Handle(new(controllers.UserController))

	// http://localhost:8080/noexist
	// and all controller's methods like
	// http://localhost:8080/users/1
	// http://localhost:8080/user/register
	// http://localhost:8080/user/login
	// http://localhost:8080/user/me
	// http://localhost:8080/user/logout
	// basic auth: "admin", "password", see "./middleware/basicauth.go" source file.
	app.Run(
		// Starts the web server at localhost:8080
		iris.Addr("localhost:8080"),
		// Ignores err server closed log when CTRL/CMD+C pressed.
		iris.WithoutServerError(iris.ErrServerClosed),
		// Enables faster json serialization and more.
		iris.WithOptimizations,
	)

	defer db.Close()
}

func NewDataBaseEngine() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "d:\\sqlite\\evaluation.db")
	//fmt.Println(err)
	db.AutoMigrate(&datamodels.User{}, &datamodels.Device{})
	//defer db.Close()
	return db, err
}
