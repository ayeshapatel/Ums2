package main

import (

//"net/http"
//_ "github.com/go-sql-driver/mysql"
//  "github.com/julienschmidt/httprouter"
//"html/template"
)

/*
func view(w http.ResponseWriter,r *http.Request,_ httprouter.Params) {
	t, _ := template.ParseFiles("UserList.html")
	type UserL struct {
		Id      int
		Fname   string
		Lname   string
		Email   string
		Password    string
		Country string
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


}*/
