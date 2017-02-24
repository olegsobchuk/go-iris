package main

import (
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

	app.Adapt(iris.DevLogger())
	app.Adapt(httprouter.New())
	app.Adapt(view.HTML("./app/views/", ".html"))

	// todos
	todos := []Todo{
		{"Initialize application", true},
		{"Upload to github", true},
		{"Make nice template", false},
	}

	app.Get("/todo", func(ctx *iris.Context) {
		ctx.Render("todo.html", struct{ Todos []Todo }{todos})
	})

	// contact_details
	app.Get("/contact_detail", func(ctx *iris.Context) {
		ctx.Render("contact_detail.html", nil)
	})

	app.Post("/contact_detail", func(ctx *iris.Context) {
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
