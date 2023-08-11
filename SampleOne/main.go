package main

import (
	"fmt"
	"net/http"
	"text/template"
)

var temp *template.Template

var userData = make(map[string]Signupdata)

type Signupdata struct {
	ConfirmPassword string
	Email           string
	PhoneNumber     string
	Name            string
	Password        string
}

func init() {
	temp = template.Must(template.ParseGlob("template/*.html"))
}
func main() {
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/home", homePage)
	http.HandleFunc("/logined", postmethod)
	http.HandleFunc("/signed", signupmethod)
	http.HandleFunc("/logout", logout)
	http.Handle("/template/", http.StripPrefix("/template/", http.FileServer(http.Dir("template"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":9999", nil)
}
func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	cookie, err := r.Cookie("logincookie")
	if err == nil && cookie.Value != "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	temp.ExecuteTemplate(w, "homepage.html", nil)
}
func loginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	cookie, err := r.Cookie("logincookie")
	if err == nil && cookie.Value != "" {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	temp.ExecuteTemplate(w, "loginPage.html", nil)
}
func signupPage(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "signupPage.html", nil)

}
func postmethod(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("emailLogin")
	password := r.FormValue("passwordLogin")
	SignupData, ok := userData[email]
	if email == "" {
		temp.ExecuteTemplate(w, "loginPage.html", "email is invalid")
		fmt.Println("email is not given")
		return
	} else if password == "" {
		temp.ExecuteTemplate(w, "loginpage.html", " password invalid")
		fmt.Println("password is not given")
		return
	}
	if ok && password == SignupData.Password {
		CookieForLogin := &http.Cookie{}
		CookieForLogin.Name = "logincookie"
		CookieForLogin.Value = "cookievalue"
		CookieForLogin.MaxAge = 200
		CookieForLogin.Path = "/"
		http.SetCookie(w, CookieForLogin)

		http.Redirect(w, r, "/home", http.StatusSeeOther)

		fmt.Println(CookieForLogin)

	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

}
func logout(w http.ResponseWriter, r *http.Request) {
	Cookielogout := http.Cookie{}
	Cookielogout.Name = "logincookie"
	Cookielogout.Value = ""
	Cookielogout.MaxAge = -1
	Cookielogout.Path = "/"
	http.SetCookie(w, &Cookielogout)
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	// http.Redirect(w, r, "/login", http.StatusSeeOther)

	cookie, err := r.Cookie("logincookie")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}
func signupmethod(w http.ResponseWriter, r *http.Request) {

	firstname := r.FormValue("firstname")
	email := r.FormValue("email")
	phonenumber := r.FormValue("phonenumber")
	password := r.FormValue("password")
	confirmpassword := r.FormValue("confirmpassword")
	if firstname == "" {

		temp.ExecuteTemplate(w, "signupPage.html", "Name is required")
		return
	}

	if email == "" {

		temp.ExecuteTemplate(w, "signupPage.html", "Email is required")

		return
	}
	if password == "" {

		temp.ExecuteTemplate(w, "signupPage.html", "Password is required")

		return
	}
	if phonenumber == "" {

		temp.ExecuteTemplate(w, "signupPage.html", "password do not match")

		return
	}
	if confirmpassword != password {

		temp.ExecuteTemplate(w, "signupPage.html", "not matched")

		return
	}
	userData[email] = Signupdata{Email: email,
		Password:    password,
		Name:        firstname,
		PhoneNumber: phonenumber,
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	fmt.Print(userData)
}
