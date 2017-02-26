package main

import (
	"github.com/kataras/iris/adaptors/sessions"
	"github.com/kataras/iris/middleware/logger"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/view"
)

// Todo struct of todo lisy
type Todo struct {
	Task string
	Done bool
}

// Contact struct for contacts
type Contact struct {
	Email   string
	Subject string
	Message string
}

// User struct
type User struct {
	FirstName string `json:"f_name"`
	LastName  string `json:"l_name"`
	Age       int    `json:"years"`
}

func main() {
	app := iris.New(iris.Configuration{Gzip: false, Charset: "UTF-8"})
	log := logger.New()

	app.Adapt(iris.DevLogger())
	app.Adapt(httprouter.New())
	app.Adapt(view.HTML("./app/views/", ".html"))
	app.StaticWeb("/public", "./app/assets")
	app.Favicon("./app/assets/images/favicon.ico", "/favicon.ico")
	app.Adapt(sessions.New(sessions.Config{
		Cookie:       "super_key",
		DecodeCookie: true,
	}))
	// error custom page
	app.OnError(iris.StatusInternalServerError, func(ctx *iris.Context) {
		log.Serve(ctx)
		ctx.RenderWithStatus(iris.StatusNotFound, "errors/500.html", nil)
	})
	app.OnError(iris.StatusNotFound, func(ctx *iris.Context) {
		log.Serve(ctx)
		ctx.RenderWithStatus(iris.StatusNotFound, "errors/404.html", nil)
	})

	// Root page
	app.Get("/", func(ctx *iris.Context) {
		log.Serve(ctx)
		ctx.Render("welcome.html", nil)
	})

	// Session
	app.Get("/login", func(ctx *iris.Context) {
		ctx.Session().Set("userID", 1001)
		userID, error := ctx.Session().GetInt("userID")
		if error != nil {
			ctx.Writef("Error has been occurred: %s", error)
			return
		}
		ctx.Writef("Added userID with %d", userID)
	})

	app.Get("/logout", func(ctx *iris.Context) {
		ctx.SessionDestroy()
		ctx.WriteString("Session has been destroyed")
	})

	app.Get("/try", func(ctx *iris.Context) {
		userID, error := ctx.Session().GetInt("userID")
		if error != nil {
			ctx.WriteString("User unregistred. Follow /login please.")
			return
		}
		ctx.Writef("USER_ID is %d", userID)
	})

	// todos
	todos := []Todo{
		{"Initialize application", true},
		{"Upload to github", true},
		{"Make nice template", false},
	}

	app.Get("/todo", func(ctx *iris.Context) {
		log.Serve(ctx)
		ctx.Render("todo.html", struct{ Todos []Todo }{todos})
	})

	// contact_details
	app.Get("/contact_detail", func(ctx *iris.Context) {
		log.Serve(ctx)
		ctx.Render("contact_detail.html", nil)
	})

	app.Post("/contact_detail", func(ctx *iris.Context) {
		log.Serve(ctx)
		// contact := Contact{
		// 	Email:   ctx.FormValue("email"),
		// 	Subject: ctx.FormValue("subject"),
		// 	Message: ctx.FormValue("message"),
		// }

		var contact Contact
		ctx.ReadForm(&contact)
		var emailPresent bool
		if contact.Email == "" {
			emailPresent = false
		} else {
			emailPresent = true
		}
		ctx.Render("contact_detail.html", struct{ Success bool }{emailPresent})
	})

	// JSON decode/encode
	app.Post("/decode_json", func(ctx *iris.Context) {
		log.Serve(ctx)
		var user User
		ctx.ReadJSON(&user)

		ctx.Writef("%s %s is %d ages!\n", user.FirstName, user.LastName, user.Age)
	})

	app.Get("/encode_json", func(ctx *iris.Context) {
		developer := User{"Oleg", "Sobchuk", 35}
		ctx.JSON(iris.StatusOK, developer)
	})

	app.Listen(":8080")
}
