package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type Rsvp struct {
	Name, Email, Phone string
	WillAttend         bool
}

var responses = make([]*Rsvp, 0, 10)
var templates = make(map[string]*template.Template, 3)

func loadTemplates() {
	templateNames := [5]string{"welcome", "form", "thanks", "sorry", "list"}

	for index, name := range templateNames {
		t, err := template.ParseFiles("templates/layout.html", "templates/"+name+".html")
		if err == nil {
			templates[name] = t
			fmt.Println("Loaded template", index, name)
		} else {
			panic(err)
		}
	}
}

func main() {
	loadTemplates()

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/form", formHandler)

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	templates["welcome"].Execute(writer, nil)
}

func listHandler(writer http.ResponseWriter, request *http.Request) {
	templates["list"].Execute(writer, responses)
}

type formData struct {
	*Rsvp
	Errors []string
}

func formHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		templates["form"].Execute(writer, formData{
			Rsvp: &Rsvp{}, Errors: []string{},
		})
	} else if request.Method == http.MethodPost {

		request.ParseForm()
		responseData := Rsvp{
			Name:       request.FormValue("name"),
			Email:      request.FormValue("email"),
			Phone:      request.FormValue("phone"),
			WillAttend: request.FormValue("willAttend") == "true",
		}

		errors := []string{}
		if strings.TrimSpace(responseData.Name) == "" {
			errors = append(errors, "Please enter your name.")
		}
		if strings.TrimSpace(responseData.Phone) == "" {
			errors = append(errors, "Please enter your phone number.")
		}
		if strings.TrimSpace(responseData.Email) == "" {
			errors = append(errors, "Please enter your email.")
		}
		if len(errors) != 0 {
			templates["form"].Execute(writer, formData{
				Rsvp: &responseData, Errors: errors,
			})
		} else {
			responses = append(responses, &responseData)

			if responseData.WillAttend {
				templates["thanks"].Execute(writer, responseData.Name)
			} else {
				templates["sorry"].Execute(writer, responseData.Name)
			}
		}
	}
}
