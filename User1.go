package main

import (
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"text/template"
	//"os/user"

)

var store = sessions.NewCookieStore([]byte("login"))
var db *gorm.DB
var err error

type User struct {
	Id       int
	Fname    string
	Lname    string
	Email    string
	Password string
	Country  string
}

//================================User Register

func Index(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {

	http.ServeFile(w, r, "Login.html")
	r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//login
	//session.Values["foo"] = "hello"
	session, err := store.Get(r, "secret")
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["uname"] = " "
	session.Save(r, w)
	//fmt.Fprint(w, ServeRequest)
}

func Register(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
	if r.Method == "GET" {

		http.ServeFile(w, r, "Register.html")
		r.ParseForm()
	} else {

		fname := r.FormValue("fname")
		lname := r.FormValue("lname")
		email := r.FormValue("email")
		password := r.FormValue("password")
		country := r.FormValue("country")

		var pro User

		if db.HasTable("pro") == false {
			db.AutoMigrate(&User{})
			db.Find(&pro)
		}
		db.Create(&User{Fname: fname, Lname: lname, Email: email, Password: password, Country: country})
		if db.NewRecord(&User{}) == true {
			//fmt.Fprintln(w, "inserted<br>")
			//fmt.Fprint(w,tmpl)//http.Redirect(w,r,"/UserList",http.StatusAccepted)
			render(w, "Display.html")
		}
	}
}
func render(w http.ResponseWriter, Fname string) {
	t, _ := template.ParseFiles(Fname)
	t.Execute(w, nil)

}
var id int
func UpdateUser(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
	session, _ := store.Get(r, "secret")
	if r.Method == "GET" {

		var user User

		email := session.Values["uname"].(string)


		res := db.Select("id,email").Where("email = ?",email).First(&user)

		if res.RowsAffected > 0 {

			session.Values["uname"] = email
			r.ParseForm()

			type UserL struct {
				Fname   string
				Lname   string
				Email   string
				Country string
			}

			var userl UserL


			id = user.Id


			db.Where("id = ?", id).First(&user)

			userl.Fname = user.Fname
			userl.Lname = user.Lname
			userl.Email = user.Email
			userl.Country = user.Country
			t, err := template.ParseFiles("UpdateProfile.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			t.Execute(w, userl)
		}else {

			t, _ := template.ParseFiles("UpdateProfile.html")
			t.Execute(w, nil)
		}

	} else {

		fname := r.FormValue("fname")
		lname := r.FormValue("lname")
		email := r.FormValue("email")
		country := r.FormValue("country")

		//session := store.Get(r, "secret")


		var user User

		res := db.Model(&user).Where("id = ?",id).Updates(User{Fname: fname, Lname: lname, Email: email, Country: country})
		if res.RowsAffected > 0{

			session.Values["uname"]=""

			user = User{}

			db.Select("fname,lname,email,country").Where("id = ?",id).First(&user)
			 data := struct {
				 Fname string
				 Lname string
				 Email string
				 Country string
			 }{
				 Fname:user.Fname,
				 Lname:user.Lname,
				 Email:user.Email,
				 Country:user.Country,

			 }

			session.Values["uname"]=user.Email
			session.Save(r,w)
			t, _ := template.ParseFiles("UpdateProfile.html")
			t.Execute(w, data)

		}
	}

}

func ServeRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	session, err := store.Get(r, "secret")
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		r.ParseForm()

		t, _ := template.ParseFiles("Login.html")
		t.Execute(w, nil)
	} else {
		var user User

		email := r.FormValue("uname")
		password := r.FormValue("password")

		res := db.Select("id,email").Where("email = ? AND password = ?",email,password).First(&user)

		if res.RowsAffected > 0 {

			session.Values["uname"] = email
			session.Save(r,w)
			r.ParseForm()

			type UserL struct {
				Fname   string
				Lname   string
				Email   string
				Country string
			}

			var userl UserL


			id = user.Id


			db.Where("id = ?", id).First(&user)
			userl.Fname = user.Fname
			userl.Lname = user.Lname
			userl.Email = user.Email
			userl.Country = user.Country
			t, err := template.ParseFiles("UpdateProfile.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			t.Execute(w, userl)
		} else{
			render(w,"Login.html")
			return
		}

	}
}

func view(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, _ := template.ParseFiles("UserList.html")
	type UserL struct {
		Id       int
		Fname    string
		Lname    string
		Email    string
		Password string
		Country  string
	}

	//var prob []UserL
	var user []User
	res := db.Select("id, fname, lname, email, country").Find(&user)

	ans := res.RowsAffected

	prob := make([]UserL, ans)

	//	count := 1

	for i, val := range user {
		prob[i].Id = val.Id
		prob[i].Fname = val.Fname
		prob[i].Lname = val.Lname
		prob[i].Email = val.Email
		prob[i].Country = val.Country

	}

	t.Execute(w, prob)
}

/*
func Logout(w http.ResponseWriter,r *http.Response, _ httprouter.Param){

	http.Redirect(w, r, "/", 302)
}*/

func main() {

	db, err = gorm.Open("mysql", "root:password@/User?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/User", Register)

	router.POST("/User", Register)

	router.GET("/UpdatePro", UpdateUser)
	router.POST("/UpdatePro",UpdateUser)

	router.POST("/sr", ServeRequest)

	router.GET("/UserList", view)

	log.Fatal(http.ListenAndServe(":8000", context.ClearHandler(router)))
}
