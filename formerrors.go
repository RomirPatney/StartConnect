package main

import (
  "regexp"
  "strings"
  _ "github.com/go-sql-driver/mysql"
  "database/sql"
  "fmt"
  "log"
  "golang.org/x/crypto/bcrypt"
)
type Errmessage struct {
  Loginerr string
  Registererr string
}
type caccount struct {
  Cname string
  Cemail string
  Cwebsite string
  CImg string
}

type uaccount struct {
  Firstname string
  Lastname string
  Email string
}

type update_pword struct {
  what string
  email string
  op string
  ope string
  np string
  npe string
  cnp string
  cnpe string
}

type Description struct {
  D string
  DErr string
}
type Pos struct {
  Pos_name string
  Num_available string
  Desc string
  Pos_nameErr string
  Num_availableErr string
  DescErr string
}

type Login struct {
  What string
  Email string
  Password string
  EmailErr string
  PasswordErr string
  LoginErr string
}

type User struct {
  Firstname string
  Lastname string
  Email    string
  Password string
  Confirm_password string
  FirstnameErr string
  LastnameErr string
  EmailErr string
  PasswordErr string
  Confirm_passwordErr string
}

type Company struct {
  Website string
  WebsiteErr string
  Company_name string
  Email string
  Password string
  Confirm_password string
  Company_nameErr string
  EmailErr string
  PasswordErr string
  Confirm_passwordErr string
}

func (p *update_pword) Validate_update_pword() bool {
  if strings.TrimSpace(p.op) == "" {
    p.ope = "Please write your old password"
    return false
  }

  if strings.TrimSpace(p.np) == "" {
    p.npe = "Please write your new password"
    return false
  }

  if strings.TrimSpace(p.cnp) == "" {
    p.cnpe = "Please confirm your new password"
    return false
  }

  if p.cnp != p.np {
    p.cnpe = "Your passwords don't match"
    return false
  }

  db, err := db_connect()
  if err != nil {
    log.Fatal(err)
  } else {
    fmt.Println("Connection successful")
  }
  if p.what == "Companies" {
    rows, err := db.Query("select password from companies where email = ?", p.email)
    if err != nil {
      log.Fatal(err)
    }
    defer rows.Close()
    var hashfromdb string
    for rows.Next() {
      err := rows.Scan(&hashfromdb)
      if err != nil {
        log.Fatal(err)
      }
    }
    if hashfromdb == "" {
      p.ope = "Email invalid"
    } else {
      if err := bcrypt.CompareHashAndPassword([]byte(hashfromdb), []byte(p.op)); err != nil {
        p.ope = "Password invalid"
      }
    }
  } else {
    rows, err := db.Query("select password from users where email = ?", p.email)
    if err != nil {
      log.Fatal(err)
    }
    defer rows.Close()
    var hashfromdb string
    for rows.Next() {
      err := rows.Scan(&hashfromdb)
      if err != nil {
        log.Fatal(err)
      }
    }
    if hashfromdb == "" {
      p.ope = "Email invalid"
    } else {
      if err := bcrypt.CompareHashAndPassword([]byte(hashfromdb), []byte(p.op)); err != nil {
        p.ope = "Password invalid"
      }
    }
  }
  return p.ope == "" && p.npe == "" && p.cnpe == ""
}

func (p *Pos) Validate_pos() bool {
  if strings.TrimSpace(p.Pos_name) == "" {
    p.Pos_nameErr = "Please write a position name"
  }
  if strings.TrimSpace(p.Num_available) == "" {
    p.Num_availableErr = "Please write the number of positions available"
  }
  if strings.TrimSpace(p.Desc) == "" {
    p.DescErr = "Please write a description"
  }
  return p.Num_availableErr == "" && p.DescErr == "" && p.Pos_nameErr == ""
}

func (login *Login) Validate_user_login() bool {
  if strings.TrimSpace(login.Email) == "" {
    login.EmailErr = "Please write an email"
  }

  if strings.TrimSpace(login.Password) == "" {
    login.PasswordErr = "Please write a password"
  }

  db, err := db_connect()
  if err != nil {
    log.Fatal(err)
  } else {
    fmt.Println("Connection successful")
  }
  var rows *sql.Rows
  if login.What == "users" {
    rows, err = db.Query("select password from users where email = ?", login.Email)
    if err != nil {
      log.Fatal(err)
    }
  } else {
    rows, err = db.Query("select password from companies where email = ?", login.Email)
    if err != nil {
      log.Fatal(err)
    }
  }
  defer rows.Close()
  var hashfromdb string
  for rows.Next() {
    err := rows.Scan(&hashfromdb)
    if err != nil {
      log.Fatal(err)
    }
  }
  if hashfromdb == "" {
    login.LoginErr = "Email invalid"
  } else {
    if err := bcrypt.CompareHashAndPassword([]byte(hashfromdb), []byte(login.Password)); err != nil {
        login.LoginErr = "Password invalid"
    }
  }
  return login.EmailErr == "" && login.PasswordErr == "" && login.LoginErr == ""
}


func (usr *User) Validate_user_registration() bool {

  if strings.TrimSpace(usr.Firstname) == "" {
    usr.FirstnameErr = "Please write a first name"
  }

  if strings.TrimSpace(usr.Lastname) == "" {
    usr.LastnameErr = "Please write a last name"
  }

  re := regexp.MustCompile(".+@.+\\..+")
  matched := re.Match([]byte(usr.Email))
  if strings.TrimSpace(usr.Email) == "" {
    usr.EmailErr = "Please write an email"
  } else if matched == false {
    usr.EmailErr = "Please enter a valid email address"
  } else {
    db, err := db_connect()
    if err != nil {
      log.Fatal(err)
    } else {
      fmt.Println("Connection successful")
    }
    rows, err := db.Query("select firstname from users where email = ?", usr.Email)
    if err != nil {
      log.Fatal(err)
    }
    defer rows.Close()
    var name string
    for rows.Next() {
      err := rows.Scan(&name)
      if err != nil {
        log.Fatal(err)
      }
    }
    if name != "" {
      usr.EmailErr = "This email already exists in our database"
    }
  }

  if strings.TrimSpace(usr.Password) == "" {
    usr.PasswordErr = "Please write a password"
  }

  if strings.TrimSpace(usr.Confirm_password) == "" {
    usr.Confirm_passwordErr = "Please confirm your password"
  } else if usr.Password != usr.Confirm_password {
    usr.Confirm_passwordErr = "You wrote different passwords"
  }
  return usr.FirstnameErr == ""  && usr.LastnameErr == "" &&
  usr.EmailErr == "" && usr.PasswordErr ==  "" && usr.Confirm_passwordErr == ""
}

func (cmpy *Company) Validate_company_registration() bool {

  if strings.TrimSpace(cmpy.Company_name) == "" {
    cmpy.Company_nameErr = "Please write a company name"
  }

  re := regexp.MustCompile(".+@.+\\..+")
  matched := re.Match([]byte(cmpy.Email))
  if strings.TrimSpace(cmpy.Email) == "" {
    cmpy.EmailErr = "Please write an email"
  } else if matched == false {
    cmpy.EmailErr = "Please enter a valid email address"
  } else {
    db, err := db_connect()
    if err != nil {
      log.Fatal(err)
    }

    fmt.Print("hi   ")
    rows, err := db.Query("select company_name from companies where email = ?", cmpy.Email)
    if err != nil {
      log.Fatal(err)
    }

    fmt.Print("hi1   ")
    defer rows.Close()
    var name string
    for rows.Next() {
      err := rows.Scan(&name)
      if err != nil {
        log.Fatal(err)
      }
    }
    if name != "" {
      cmpy.EmailErr = "This email already exists in our database"
    }
  }

  if strings.TrimSpace(cmpy.Password) == "" {
    cmpy.PasswordErr = "Please write a password"
  }

  if strings.TrimSpace(cmpy.Confirm_password) == "" {
    cmpy.Confirm_passwordErr = "Please confirm your password"
  } else if cmpy.Password != cmpy.Confirm_password {
    cmpy.Confirm_passwordErr = "You wrote different passwords"
  }
  return cmpy.Company_nameErr == "" && cmpy.EmailErr == "" && cmpy.PasswordErr ==  "" && cmpy.Confirm_passwordErr == ""
}
