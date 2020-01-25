package main

import (
    "database/sql"
    "log"
    "net/http"
    "text/template"
    "fmt"
    
    _ "github.com/go-sql-driver/mysql"
    "github.com/kataras/go-sessions"
)


func dbConn() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := ""
    dbName := "website_crud"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    if err != nil {
        panic(err.Error())
    }
    return db
}

func HandlerIndex(w http.ResponseWriter, r *http.Request) {
    session := sessions.Start(w, r)
    var data = map[string]string{
        "name": session.GetString("name"),
        "message":  "Welcome Login APP",
    }
    var tmp = template.Must(template.ParseFiles(
        "views/Header.html",
        "views/Menu.html",
        "views/Index.html",
        "views/Footer.html",
    ))

    var error = tmp.ExecuteTemplate(w,"Index",data)
    if error != nil {
        http.Error(w, error.Error(), http.StatusInternalServerError)
    }
}

func HandlerRegister(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        name := r.FormValue("name")
        email := r.FormValue("email")
        password := r.FormValue("password")
        insForm, err := db.Prepare("INSERT INTO users (name, email, password) VALUES(?,?,?)")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(name,email,password)
        session := sessions.Start(w, r)
        session.Set("email", email)
        session.Set("name", name)
        defer db.Close()
        http.Redirect(w, r, "/user", 301)

    }else{
        var tmp = template.Must(template.ParseFiles(
            "views/Header.html",
            "views/Menu.html",
            "views/Register.html",
            "views/Footer.html",
        ))
        data:=""
        var error = tmp.ExecuteTemplate(w,"Register",data)
        if error != nil {
            http.Error(w, error.Error(), http.StatusInternalServerError)
        }
    }

    
}


func HandlerLogin(w http.ResponseWriter, r *http.Request) {
    session := sessions.Start(w, r)
	if len(session.GetString("name")) != 0 {
		http.Redirect(w, r, "/user", 302)
    }else{
        if r.Method == "POST" {
            emailLogin:=r.FormValue("email")
            passwordLogin:=r.FormValue("password")
    
            db := dbConn()
            selDB, err := db.Query("SELECT * FROM users WHERE email=?", emailLogin)
            if err != nil {
                panic(err.Error())
            }
                
            for selDB.Next() {
                var id int
                var name,email,password string
                err = selDB.Scan(&id, &name, &email, &password)
                if err != nil {
                    panic(err.Error())
                
                }
                if password==passwordLogin {
                    session := sessions.Start(w, r)
                    session.Set("email", email)
                    session.Set("name", name)
                    http.Redirect(w, r, "/user", 302)
                }else{//wrong password
                    fmt.Fprintf(w,"wrong password")
                }
            }
            fmt.Fprintf(w,"user not found")
            
            defer db.Close()
    
    
        }else{
            var tmp = template.Must(template.ParseFiles(
                "views/Header.html",
                "views/Menu.html",
                "views/Login.html",
                "views/Footer.html",
            ))
            data:=""
            var error = tmp.ExecuteTemplate(w,"Login",data)
            if error != nil {
                http.Error(w, error.Error(), http.StatusInternalServerError)
            }
        }

    }

}


func HandlerUser(w http.ResponseWriter, r *http.Request) {
    session := sessions.Start(w, r)
	if len(session.GetString("name")) == 0 {
		http.Redirect(w, r, "/user/login", 301)
	}else{
        var data = map[string]string{
            "name": session.GetString("name"),
            "message":  "Welcome User Home !",
        }
    
        var tmp = template.Must(template.ParseFiles(
            "views/Header.html",
            "views/Menu.html",
            "views/User.html",
            "views/Footer.html",
        ))
    
        var error = tmp.ExecuteTemplate(w,"User",data)
        if error != nil {
            http.Error(w, error.Error(), http.StatusInternalServerError)
        }
    }

}


func HandlerLogout(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	session.Clear()
	sessions.Destroy(w, r)
	http.Redirect(w, r, "/", 302)
}



func main() {
    log.Println("Server started on: http://localhost:8000")
    http.HandleFunc("/", HandlerIndex)
    http.HandleFunc("/user", HandlerUser)
    http.HandleFunc("/user/register", HandlerRegister)
    http.HandleFunc("/user/login", HandlerLogin)
    http.HandleFunc("/user/logout", HandlerLogout)

    http.ListenAndServe(":8000", nil)
}