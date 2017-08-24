package main

import (
  "github.com/gorilla/mux"
	"html/template"
  "log"
  "net/http"
  _ "github.com/go-sql-driver/mysql"
  "database/sql"
  "fmt"
  "golang.org/x/crypto/bcrypt"
	"io"
	"os"
	"encoding/base64"
	"strings"
)



func main() {

  router := mux.NewRouter()
  router.HandleFunc("/", home_page)
	router.HandleFunc("/companypage/account.html", cmpy_account_page)
  router.HandleFunc("/companypage/mainpage.html", cmpy_main_page)
	router.HandleFunc("/studentpage/account.html", std_account_page)
	router.HandleFunc("/studentpage/mainpage.html", std_main_page)
	router.HandleFunc("/logout", logout)
	router.HandleFunc("/forgotp", forgotp)
	router.HandleFunc("/companypage/edit-account.html", cmpy_edit_pos)
	router.HandleFunc("/delete_pos", delete_pos).Methods("Post")
	router.HandleFunc("/update_account", update_account).Methods("Post")
	router.HandleFunc("/upload_pic", upload_pic).Methods("Post")
	router.HandleFunc("/add_pos", add_pos).Methods("Post")
	router.HandleFunc("/register_user", register_user).Methods("Post")
  router.HandleFunc("/login_company", login).Methods("Post")
  router.HandleFunc("/register_company", register_company).Methods("Post")
  router.HandleFunc("/login_user", login).Methods("Post")

	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("/users/khalilbenayed/go/src/hello/projects/startconnectmaster"))))
	log.Println("Listening...")
	//port := ":" + os.Getenv("PORT")
	//fmt.Print(port)
	http.ListenAndServe(":3009", router)
}

func forgotp(w http.ResponseWriter, r *http.Request) {
	where := r.FormValue("where")
	db, err := sql.Open("mysql", "root:@/Velocity_Connect")
	if err != nil {
		log.Fatal(err)
	}
	if where == "companies" {
		email := r.FormValue("cemail")
		cname := r.FormValue("cname")
		pwd := r.FormValue("pwd")
		cpwd := r.FormValue("cpwd")
		rows, err := db.Query("select company_name from companies where email = ? and company_name = ?", email, cname)
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
		if name != "" && pwd == cpwd {
			hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
			if err != nil {
				log.Fatal(err)
			}
			stmt, err := db.Prepare("UPDATE companies SET password = ? WHERE email = ?")
			if err != nil {
				log.Fatal(err)
			}
			_, err1 := stmt.Exec(hash, email)
			if err1 != nil {
				log.Fatal(err1)
			}
			http.Redirect(w, r, "publicfile/one-page.html", 301)
		} else {
			http.Redirect(w, r, "companypage/forgotpwd.html", 301)
		}
	} else {
		email := r.FormValue("email")
		firstname := r.FormValue("firstname")
		lastname := r.FormValue("lastname")
		pwd := r.FormValue("pwd")
		cpwd := r.FormValue("cpwd")
		rows, err := db.Query("select firstname from users where email = ? and firstname = ? and lastname = ?", email, firstname, lastname)
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
		if name != "" && pwd == cpwd {
			hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
			if err != nil {
				log.Fatal(err)
			}
			stmt, err := db.Prepare("UPDATE users SET password = ? WHERE email = ?")
			if err != nil {
				log.Fatal(err)
			}
			_, err1 := stmt.Exec(hash, email)
			if err1 != nil {
				log.Fatal(err1)
			}
			http.Redirect(w, r, "publicfile/one-page.html", 301)
		} else {
			http.Redirect(w, r, "studentpage/forgotpwd.html", 301)
		}
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	email := http.Cookie{
		Name: "cemail",
		MaxAge: -1}
	http.SetCookie(w, &email)
	http.Redirect(w, r, "publicfile/one-page.html", 301)
}

func delete_pos(w http.ResponseWriter, r *http.Request) {
	email, err := r.Cookie("cemail")
	if err == nil && email.Value != "" {
		pos := r.FormValue("pos")
		db, err := sql.Open("mysql", "root:@/Velocity_Connect")
		if err != nil {
			log.Fatal(err)
		}
		_, err1 := db.Query("DELETE FROM Positions WHERE CompanyEmail=? AND PositionName=?", email.Value, pos)
		if err1 != nil {
			log.Fatal(err1)
		}
		http.Redirect(w, r, "companypage/edit-account.html", 301)
	} else {
		fmt.Print(err)
		http.Redirect(w, r, "publicfile/one-page.html", 301)
	}
}

func upload_pic(w http.ResponseWriter, r *http.Request) {
	email, err := r.Cookie("cemail")
	if err == nil && email.Value != "" {
		db, err := sql.Open("mysql", "root:@/Velocity_Connect")
		if err != nil {
			log.Fatal(err)
		}

		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		f, err := os.OpenFile("./"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		fil, err := os.Open("./" + handler.Filename)

		defer fil.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		buff := make([]byte, 1024*1024) // why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
		_, err = fil.Read(buff)

		filetype := http.DetectContentType(buff)

		stat, err1 := fil.Stat()
		if err1 != nil {
			log.Fatal(err1)
		}
		filesize := stat.Size()

		err2 := os.Remove("./" + handler.Filename)

		if err2 != nil {
			fmt.Println(err2)
			return
		}

		rows, err := db.Query("select imagename from companyprofilepic where companyemail = ?", email.Value)
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
			stmt, err := db.Prepare("UPDATE companyprofilepic SET ImageType = ?, image = ?, imagesize = ?, imagename = ? WHERE companyemail = ?")
			if err != nil {
				log.Fatal(err)
			}
			_, err1 := stmt.Exec(filetype, buff, filesize, handler.Filename, email.Value)
			if err1 != nil {
				log.Fatal(err1)
			}
		} else {

			stmt, err := db.Prepare("INSERT INTO CompanyProfilePic (CompanyEmail, ImageType, Image, ImageSize, ImageName) VALUES (?, ?, ?, ?, ?)")
			if err != nil {
				log.Fatal(err)
			}
			res, err := stmt.Exec(email.Value, filetype, buff, filesize, handler.Filename)
			if err != nil {
				log.Fatal(err)
			}
			lastId, err := res.LastInsertId()
			if err != nil {
				log.Fatal(err)
			}
			rowCnt, err := res.RowsAffected()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
			defer stmt.Close()
		}
		http.Redirect(w, r, "companypage/account.html", 301)
	} else {
		fmt.Print(err)
		http.Redirect(w, r, "publicfile/onepage.html", 301)
	}
}

func update_account(w http.ResponseWriter, r *http.Request) {
	email, err := r.Cookie("cemail")
	if err == nil && email.Value != "" {
		what := r.FormValue("what")
		if what == "Company_Name" {
			with := r.FormValue(what)
			if with != "" {
				db, err := sql.Open("mysql", "root:@/Velocity_Connect")
				if err != nil {
					log.Fatal(err)
				}
				stmt, err := db.Prepare("UPDATE Companies SET Company_Name = ? WHERE email=?")
				if err != nil {
					fmt.Print(err)
				}
				_, err1 := stmt.Exec(with, email.Value)
				if err1 != nil {
					fmt.Print(err1)
				}
			}
			http.Redirect(w, r, "companypage/account.html", 301)
		} else if what == "Firstname" {
			with := r.FormValue(what)
			if with != "" {
				db, err := sql.Open("mysql", "root:@/Velocity_Connect")
				if err != nil {
					log.Fatal(err)
				}
				stmt, err := db.Prepare("UPDATE Users SET firstname = ? WHERE email=?")
				if err != nil {
					fmt.Print(err)
				}
				_, err1 := stmt.Exec(with, email.Value)
				if err1 != nil {
					fmt.Print(err1)
				}
			}
			http.Redirect(w, r, "studentpage/account.html", 301)
		} else if what == "Lastname" {
			with := r.FormValue(what)
			if with != "" {
				db, err := sql.Open("mysql", "root:@/Velocity_Connect")
				if err != nil {
					log.Fatal(err)
				}
				stmt, err := db.Prepare("UPDATE Users SET lastname = ? WHERE email=?")
				if err != nil {
					fmt.Print(err)
				}
				_, err1 := stmt.Exec(with, email.Value)
				if err1 != nil {
					fmt.Print(err1)
				}
			}
			http.Redirect(w, r, "studentpage/account.html", 301)
		} else if what == "Cemail" {
			with := r.FormValue(what)
			if with != "" {
				db, err := sql.Open("mysql", "root:@/Velocity_Connect")
				if err != nil {
					log.Fatal(err)
				}
				rows, err := db.Query("select company_name from companies where email = ?", with)
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
				emailerr := ""
				if name != "" {
					emailerr = "This email already exists in our database"
				}
				if emailerr == "" {
					stmt, err := db.Prepare("UPDATE Companies SET email = ? WHERE email=?")
					if err != nil {
						fmt.Print(err)
					}
					_, err1 := stmt.Exec(with, email.Value)
					if err1 != nil {
						fmt.Print(err1)
					}
					stmt1, err := db.Prepare("UPDATE Positions SET CompanyEmail = ? WHERE CompanyEmail=?")
					if err != nil {
						fmt.Print(err)
					}
					_, err2 := stmt1.Exec(with, email.Value)
					if err2 != nil {
						fmt.Print(err1)
					}
					stmt2, err := db.Prepare("UPDATE CompanyProfilePic SET CompanyEmail = ? WHERE CompanyEmail=?")
					if err != nil {
						fmt.Print(err)
					}
					_, err3 := stmt2.Exec(with, email.Value)
					if err3 != nil {
						fmt.Print(err1)
					}
					email := http.Cookie{
						Name:   "cemail",
						Value:  with,
						MaxAge: 3600}
					http.SetCookie(w, &email)
				} else {
					fmt.Print(emailerr, name)

				}
			}
			http.Redirect(w, r, "companypage/account.html", 301)
		} else if what == "Uemail" {
			with := r.FormValue(what)
			if with != "" {
				db, err := sql.Open("mysql", "root:@/Velocity_Connect")
				if err != nil {
					log.Fatal(err)
				}
				rows, err := db.Query("SELECT Firstname from Users where email =?", with)
				var name string
				for rows.Next() {
					err = rows.Scan(&name)
				}
				var emailerr string
				if name != "" {
					emailerr = "This email already exists in our database"
				}
				if emailerr == "" {
					stmt, err := db.Prepare("UPDATE Users SET email = ? WHERE email=?")
					if err != nil {
						fmt.Print(err)
					}
					_, err1 := stmt.Exec(with, email.Value)
					if err1 != nil {
						fmt.Print(err1)
					}
					email := http.Cookie{
						Name:   "cemail",
						Value:  with,
						MaxAge: 3600}
					http.SetCookie(w, &email)
				} else {
					fmt.Print(emailerr, name)
				}
			}
			http.Redirect(w, r, "studentpage/account.html", 301)
		} else if what == "Website" {
			with := r.FormValue(what)
			if with != "" {
				db, err := sql.Open("mysql", "root:@/Velocity_Connect")
				if err != nil {
					log.Fatal(err)
				}
				stmt, err := db.Prepare("UPDATE Companies SET website = ? WHERE email=?")
				if err != nil {
					fmt.Print(err)
				}
				_, err1 := stmt.Exec(with, email.Value)
				if err1 != nil {
					fmt.Print(err1)
				}
			}
			http.Redirect(w, r, "companypage/account.html", 301)
		} else if what == "Cpassword" {
		p := &update_pword{
				what: "Companies",
				email: email.Value,
				op: r.FormValue("op"),
				np: r.FormValue("np"),
				cnp: r.FormValue("cnp"),
			}
			if p.Validate_update_pword() {
				db, err := sql.Open("mysql", "root:@/Velocity_Connect")
				if err != nil {
					log.Fatal(err)
				}
				hash, err := bcrypt.GenerateFromPassword([]byte(p.np), bcrypt.DefaultCost)
				if err != nil {
					log.Fatal(err)
				}
				stmt, err := db.Prepare("UPDATE Companies SET Password = ? WHERE email=?")
				if err != nil {
					fmt.Print(err)
				}
				_, err1 := stmt.Exec(hash, email.Value)
				if err1 != nil {
					fmt.Print(err1)
				}
			}
			http.Redirect(w, r, "companypage/account.html", 301)
		} else {
			p := &update_pword{
			what: "Users",
			email: email.Value,
			op: r.FormValue("op"),
			np: r.FormValue("np"),
			cnp: r.FormValue("cnp"),
		}
		if p.Validate_update_pword() {
			db, err := sql.Open("mysql", "root:@/Velocity_Connect")
			if err != nil {
				log.Fatal(err)
			}
			hash, err := bcrypt.GenerateFromPassword([]byte(p.np), bcrypt.DefaultCost)
			if err != nil {
				log.Fatal(err)
			}
			stmt, err := db.Prepare("UPDATE Users SET Password = ? WHERE email=?")
			if err != nil {
				fmt.Print(err)
			}
			_, err1 := stmt.Exec(hash, email.Value)
			if err1 != nil {
				fmt.Print(err1)
			}
		}
			http.Redirect(w, r, "studentpage/account.html", 301)
		}
	} else {
		fmt.Print(err)
		http.Redirect(w, r, "publicfile/onepage.html", 301)
	}
}

func add_pos(w http.ResponseWriter, r *http.Request) {
	email, err := r.Cookie("cemail")
	if err == nil && email.Value != "" {
		pos := &Pos{
			Pos_name:      r.FormValue("pos_name"),
			Num_available: r.FormValue("num_available"),
			Desc:          r.FormValue("desc"),
		}

		if pos.Validate_pos() == false {
			return
		} else {
			db, err := sql.Open("mysql", "root:@/Velocity_Connect")
			if err != nil {
				log.Fatal(err)
			}
			// insert
			stmt, err := db.Prepare("INSERT INTO Positions (CompanyEmail, PositionName, NumberAvailable, Description) VALUES (?, ?, ?, ?)")
			if err != nil {
				log.Fatal(err)
			}
			res, err := stmt.Exec(email.Value, pos.Pos_name, pos.Num_available, pos.Desc)
			if err != nil {
				log.Fatal(err)
			}
			lastId, err := res.LastInsertId()
			if err != nil {
				log.Fatal(err)
			}
			rowCnt, err := res.RowsAffected()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
			defer stmt.Close()
			http.Redirect(w, r, "companypage/edit-account.html", 301)
		}
	} else {
		fmt.Print(err)
		http.Redirect(w, r, "publicfile/onepage.html", 301)
	}
}

func cmpy_edit_pos(w http.ResponseWriter, r *http.Request) {
	email, err := r.Cookie("cemail")
	if err == nil && email.Value != "" {
		db, err := sql.Open("mysql", "root:@/Velocity_Connect")
		if err != nil {
			log.Fatal(err)
		}
		rows, err := db.Query("SELECT PositionName, NumberAvailable, Description FROM Positions WHERE CompanyEmail=?", email.Value)

		cols, _ := rows.Columns()

		final_result := make(map[int]map[string]interface{})
		final_map := map[int]map[string]map[string][]string{}

		count := 0

		for rows.Next() {
			// Create a slice of interface{}'s to represent each column,
			// and a second slice to contain pointers to each item in the columns slice.
			columns := make([]interface{}, len(cols))
			columnPointers := make([]interface{}, len(cols))
			for i, _ := range columns {
				columnPointers[i] = &columns[i]
			}

			// Scan the result into the column pointers...
			if err := rows.Scan(columnPointers...); err != nil {
				fmt.Print(err)
			}

			// Create our map, and retrieve the value for each column from the pointers slice,
			// storing it in the map with the name of the column as the key.

			tmp := make(map[string]interface{})

			for i, col := range cols {
				//var v interface{}
				val := columns[i]
				/*b, ok := val.([]byte)
			if (ok) {
				v = string(b)
			} else {
				v = val
			}*/
				c := fmt.Sprintf("%s", col)
				v := fmt.Sprintf("%s", val)
				tmp[c] = v
			}
			final_result[count] = tmp
			count++
		}

		i := 0
		for i < count {
			desc_map := make(map[string][]string)
			positions_map := make(map[string]map[string][]string)

			pos := fmt.Sprintf("%s", final_result[i]["PositionName"])
			num := fmt.Sprintf("%s", final_result[i]["NumberAvailable"])
			desc := fmt.Sprintf("%s", final_result[i]["Description"])

			body := strings.Split(desc,"\n")
			desc_map[num] = body
			positions_map[pos] = desc_map

			final_map[i] = positions_map
			i++
		}
		render(w, "companypage/edit-account.html", final_map)
	} else {
		fmt.Print(err)
		http.Redirect(w, r, "publicfile/onepage.html", 301)
	}
}

func cmpy_account_page(w http.ResponseWriter, r *http.Request) {
	email, err := r.Cookie("cemail")
	if err == nil && email.Value != "" {
		db, err := sql.Open("mysql", "root:@/Velocity_Connect")
		if err != nil {
			log.Fatal(err)
		}
		var imagetype, imagesize, imagename, img string

		rows, err := db.Query("SELECT Companies.Company_name, Companies.email, Companies.website, Companyprofilepic.ImageType, Companyprofilepic.Image, Companyprofilepic.ImageSize, Companyprofilepic.ImageName  FROM Companies LEFT JOIN Companyprofilepic ON Companies.Email = Companyprofilepic.companyemail WHERE  email =?", email.Value)
		account := &caccount{}

		for rows.Next() {
			err = rows.Scan(&account.Cname, &account.Cemail, &account.Cwebsite, &imagetype, &img, &imagesize, &imagename)
		}

		if imagesize != "" {
			sEnc := base64.StdEncoding.EncodeToString([]byte(img))
			account.CImg = sEnc
		} else {
			account.CImg = ""
		}
		render(w, "companypage/account.html", account)
	} else {
		fmt.Print(err)
		http.Redirect(w, r, "publicfile/onepage.html", 301)
	}
}

func std_account_page(w http.ResponseWriter, r *http.Request) {
	email, err := r.Cookie("cemail")
	if err == nil && email.Value != "" {
		db, err := sql.Open("mysql", "root:@/Velocity_Connect")
		if err != nil {
			log.Fatal(err)
		}

		rows, err := db.Query("SELECT firstname, lastname, email FROM Users WHERE  email =?", email.Value)
		account := &uaccount{}

		for rows.Next() {
			err = rows.Scan(&account.Firstname, &account.Lastname, &account.Email)
		}
		fmt.Print("  hi  ")
		render(w, "studentpage/account.html", account)
	} else {
		fmt.Print("  hi  ", err)
		http.Redirect(w, r, "publicfile/onepage.html", 301)
	}
}


func cmpy_main_page(w http.ResponseWriter, r *http.Request) {
	email, err := r.Cookie("cemail")
	if err == nil && email.Value != "" {
		db, err := sql.Open("mysql", "root:@/Velocity_Connect")
		if err != nil {
			log.Fatal(err)
		}

		rows, err := db.Query("SELECT Companies.Company_name, Companies.Website, Companyprofilepic.Image, Companyprofilepic.Imagename, Positions.PositionName, Positions.NumberAvailable, Positions.Description FROM Companies LEFT JOIN Companyprofilepic ON Companies.Email=Companyprofilepic.Companyemail INNER JOIN Positions ON Companies.Email=Positions.CompanyEmail")
		cols, _ := rows.Columns()

		final_result := make(map[int]map[string]interface{})

		count := 0

		for rows.Next() {
			// Create a slice of interface{}'s to represent each column,
			// and a second slice to contain pointers to each item in the columns slice.
			columns := make([]interface{}, len(cols))
			columnPointers := make([]interface{}, len(cols))
			for i, _ := range columns {
				columnPointers[i] = &columns[i]
			}

			// Scan the result into the column pointers...
			if err := rows.Scan(columnPointers...); err != nil {
				fmt.Print(err)
			}

			// Create our map, and retrieve the value for each column from the pointers slice,
			// storing it in the map with the name of the column as the key.

			tmp := make(map[string]interface{})

			for i, col := range cols {
				val := columns[i]
				c := fmt.Sprintf("%s", col)
				v := fmt.Sprintf("%s", val)
				tmp[c] = v
			}
			final_result[count] = tmp
			count++
		}

		final_map := map[int]map[string]map[string]map[string]map[int]map[string]map[string][]string{}

		i := 0
		k := 0

		for i < count {
			id_map := make(map[int]map[string]map[string][]string)
			website_map := make(map[string]map[int]map[string]map[string][]string)
			cmpy_map := make(map[string]map[string]map[int]map[string]map[string][]string)
			img_map := make(map[string]map[string]map[string]map[int]map[string]map[string][]string)
			var c string
			if final_result[i]["Company_name"] != nil {
				c = fmt.Sprintf("%s", final_result[i]["Company_name"])
			} else {
				c = ""
			}

			w := fmt.Sprintf("%s", final_result[i]["Website"])
			var sEnc string
			img := fmt.Sprintf("%s", final_result[i]["Image"])

			imgname := fmt.Sprintf("%s", final_result[i]["Imagename"])

			if imgname != "%!s(<nil>)" {
				sEnc = base64.StdEncoding.EncodeToString([]byte(img))
			} else {
				sEnc = ""
			}

			j := i

			for j < count {
				desc_map := make(map[string][]string)
				positions_map := make(map[string]map[string][]string)
				var compare string
				if final_result[j]["Company_name"] != nil {
					compare = fmt.Sprintf("%s", final_result[j]["Company_name"])
				} else {
					compare = ""
				}
				if compare == c && c != "" {
					p := fmt.Sprintf("%s", final_result[j]["PositionName"])
					n := fmt.Sprintf("%s", final_result[j]["NumberAvailable"])
					d := fmt.Sprintf("%s", final_result[j]["Description"])
					body := strings.Split(d,"\n")
					desc_map[n] = body
					positions_map[p] = desc_map
					id_map[k] = positions_map
					final_result[j] = nil
					k++
				}
				j++
			}

			if c != "" {
				website_map[w] = id_map

				cmpy_map[c] = website_map

				img_map[sEnc] = cmpy_map

				final_map[i] = img_map
			}
			i++
		}

		render(w, "companypage/mainpage.html", final_map)
	} else {
		http.Redirect(w, r, "../publicfile/one-page.html", 301)
	}
}

func std_main_page(w http.ResponseWriter, r *http.Request) {
	email, err := r.Cookie("cemail")
	if err == nil && email.Value != "" {
		db, err := sql.Open("mysql", "root:@/Velocity_Connect")
		if err != nil {
			log.Fatal(err)
		}

		rows, err := db.Query("SELECT Companies.Email, Companies.Company_name, Companies.Website, Companyprofilepic.Image, Companyprofilepic.Imagename, Positions.PositionName, Positions.NumberAvailable, Positions.Description FROM Companies LEFT JOIN Companyprofilepic ON Companies.Email=Companyprofilepic.Companyemail INNER JOIN Positions ON Companies.Email=Positions.CompanyEmail")
		cols, _ := rows.Columns()

		final_result := make(map[int]map[string]interface{})

		count := 0

		for rows.Next() {
			// Create a slice of interface{}'s to represent each column,
			// and a second slice to contain pointers to each item in the columns slice.
			columns := make([]interface{}, len(cols))
			columnPointers := make([]interface{}, len(cols))
			for i, _ := range columns {
				columnPointers[i] = &columns[i]
			}

			// Scan the result into the column pointers...
			if err := rows.Scan(columnPointers...); err != nil {
				fmt.Print(err)
			}

			// Create our map, and retrieve the value for each column from the pointers slice,
			// storing it in the map with the name of the column as the key.

			tmp := make(map[string]interface{})

			for i, col := range cols {
				val := columns[i]
				c := fmt.Sprintf("%s", col)
				v := fmt.Sprintf("%s", val)
				tmp[c] = v
			}
			final_result[count] = tmp
			count++
		}

		final_map := map[int]map[string]map[string]map[string]map[string]map[int]map[string]map[string][]string{}

		i := 0
		k := 0

		for i < count {

			id_map := make(map[int]map[string]map[string][]string)
			website_map := make(map[string]map[int]map[string]map[string][]string)
			cmpy_map := make(map[string]map[string]map[int]map[string]map[string][]string)
			img_map := make(map[string]map[string]map[string]map[int]map[string]map[string][]string)
			email_map := make(map[string]map[string]map[string]map[string]map[int]map[string]map[string][]string)
			var c string
			if final_result[i]["Company_name"] != nil {
				c = fmt.Sprintf("%s", final_result[i]["Company_name"])
			} else {
				c = ""
			}

			e := fmt.Sprintf("%s", final_result[i]["Email"])
			w := fmt.Sprintf("%s", final_result[i]["Website"])
			var sEnc string
			img := fmt.Sprintf("%s", final_result[i]["Image"])

			imgname := fmt.Sprintf("%s", final_result[i]["Imagename"])

			if imgname != "%!s(<nil>)" {
				sEnc = base64.StdEncoding.EncodeToString([]byte(img))
			} else {
				sEnc = ""
			}

			j := i

			for j < count {
				desc_map := make(map[string][]string)
				positions_map := make(map[string]map[string][]string)
				var compare string
				if final_result[j]["Company_name"] != nil {
					compare = fmt.Sprintf("%s", final_result[j]["Company_name"])
				} else {
					compare = ""
				}
				if compare == c && c != "" {
					p := fmt.Sprintf("%s", final_result[j]["PositionName"])
					n := fmt.Sprintf("%s", final_result[j]["NumberAvailable"])
					d := fmt.Sprintf("%s", final_result[j]["Description"])
					body := strings.Split(d,"\n")
					desc_map[n] = body
					positions_map[p] = desc_map
					id_map[k] = positions_map
					final_result[j] = nil
					k++
				}
				j++
			}

			if c != "" {
				website_map[w] = id_map

				cmpy_map[c] = website_map

				img_map[sEnc] = cmpy_map

				email_map[e] = img_map

				final_map[i] = email_map
			}
			i++
		}

		render(w, "studentpage/mainpage.html", final_map)
	} else {
		http.Redirect(w, r, "publicfile/one-page.html", 301)
	}
}

func home_page(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "publicfile/one-page.html", 301)
}

func login(w http.ResponseWriter, r *http.Request){
  login := &Login{
	  What: r.FormValue("what"),
    Email: r.FormValue("email"),
    Password: r.FormValue("password"),
  }
  if login.Validate_user_login() == false {
    fmt.Print(login.LoginErr, login.EmailErr, login.PasswordErr)
	  if login.What == "companies" {
		  http.Redirect(w, r, "sign-up-login-form/index.html", 301)
	  } else {
		  http.Redirect(w, r, "sign-up-login-form/index1.html", 301)
	  }
  } else {
	  email := http.Cookie{
		  Name: "cemail",
		  Value: login.Email,
		  MaxAge: 3600}
	  http.SetCookie(w, &email)
	  if login.What == "companies" {
		  http.Redirect(w, r, "companypage/mainpage.html", 301)
	  } else {
		  http.Redirect(w, r, "studentpage/mainpage.html", 301)
	  }
  }
}

func register_user(w http.ResponseWriter, r *http.Request) {
  usr := &User{
    Firstname: r.FormValue("firstname"),
    Lastname: r.FormValue("lastname"),
    Email: r.FormValue("email"),
    Password: r.FormValue("password"),
    Confirm_password: r.FormValue("confirm_password"),
  }
  if usr.Validate_user_registration() == false {
	  fmt.Printf("%s %s %s %s %s", usr.FirstnameErr, usr.LastnameErr, usr.EmailErr, usr.PasswordErr, usr.Confirm_passwordErr)
	  http.Redirect(w, r, "sign-up-login-form/index1.html", 301)
  } else {
      db, err := sql.Open("mysql", "root:@/Velocity_Connect")
      if err != nil {
        log.Fatal(err)
      }
      // insert
      hash, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
      if err != nil {
        log.Fatal(err)
      }
      stmt, err := db.Prepare("INSERT INTO Users (Firstname, Lastname, Email, Password) VALUES (?, ?, ?, ?)")
      if err != nil {
        log.Fatal(err)
      }
      res, err := stmt.Exec(usr.Firstname, usr.Lastname, usr.Email, hash)//, usr.H, 0)
      if err != nil {
        log.Fatal(err)
      }
      lastId, err := res.LastInsertId()
      if err != nil {
        log.Fatal(err)
      }
      rowCnt, err := res.RowsAffected()
      if err != nil {
        log.Fatal(err)
      }
      log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
      defer stmt.Close()
	  email := http.Cookie{
		  Name: "cemail",
		  Value: usr.Email,
		  MaxAge: 3600}
	  http.SetCookie(w, &email)
	  http.Redirect(w, r, "studentpage/mainpage.html", 301)
    }
  }
  func register_company(w http.ResponseWriter, r *http.Request) {
    cmpy := &Company{
		Website: r.FormValue("website"),
      Company_name: r.FormValue("companyname"),
      Email: r.FormValue("email"),
      Password: r.FormValue("password"),
      Confirm_password: r.FormValue("confirm_password"),
    }
    if cmpy.Validate_company_registration() == false {
		fmt.Printf("%s %s %s %s %s", cmpy.Company_nameErr, cmpy.Confirm_passwordErr, cmpy.EmailErr, cmpy.PasswordErr, cmpy.Confirm_passwordErr)
		http.Redirect(w, r, "sign-up-login-form/index.html", 301)
	} else {
        db, err := sql.Open("mysql", "root:@/Velocity_Connect")
        if err != nil {
          log.Fatal(err)
        }
        // insert
        hash, err := bcrypt.GenerateFromPassword([]byte(cmpy.Password), bcrypt.DefaultCost)
        if err != nil {
          log.Fatal(err)
        }
        stmt, err := db.Prepare("INSERT INTO Companies (Company_name, Email, Password, Website) VALUES (?, ?, ?, ?)")
        if err != nil {
          log.Fatal(err)
        }
        res, err := stmt.Exec(cmpy.Company_name, cmpy.Email, hash, cmpy.Website)
        if err != nil {
          log.Fatal(err)
        }
        lastId, err := res.LastInsertId()
        if err != nil {
          log.Fatal(err)
        }
        rowCnt, err := res.RowsAffected()
        if err != nil {
          log.Fatal(err)
        }
        log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
        defer stmt.Close()
		email := http.Cookie{
			Name: "cemail",
			Value: cmpy.Email,
			MaxAge: 3600}
		http.SetCookie(w, &email)
		http.Redirect(w, r, "companypage/mainpage.html", 301)
      }
    }

  func render(w http.ResponseWriter, filename string, data interface{}) {
    tmpl, err := template.ParseFiles(filename)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    if err := tmpl.Execute(w, data); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }
  }
