package main

/**
Program ini ditujukan sebagai pemenuhan tugas Dumbways.Id Batch 46 dari day 11-16
---------------------------------------------------------------------------------
Program ini dibuat untuk menjalankan server localhost pada port 5000, sekaligus untuk keperluan autentikasi.
Menggunakan konsep struct interface,
Database PostgreSql,
Framework Echo Go-lang,
dan Bcrypt sebagai enkripsi password

"flash sebagai wadah pengiriman ke halaman web"
*/

/**
Fitur :
autentikasi login,
pengecekan email ketika didaftarkan,
penghitungan durasi pembuatan project dari hari, bulan, dan tahun,
menampilkan teknologi di home, dan detail blog sesuai yang di daftarkan user saat add project.
menampilkan project sesuai kepemilikan user yang login
*/

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"personal-web/connection"
	"personal-web/middleware"
	"reflect"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// TEMPLATE STRUCT
type Blog struct {
	Title, Description, SubDescription, Distance, Author, Image, StartDateConv, EndDateConv string
	StartDate, EndDate time.Time
	ID, Month, Day int
	Technology []string
}

type User struct {
	ID int
	Name, Email, Role, Password string
}

type Session struct {
	IsLogin bool
	NotLogin bool
	Name string
}

var dataBlog = []Blog{}
var sessionUser = Session{}

func main() {
	connection.DatabaseConnect()

	// new instance
	e := echo.New()

	// static files
	e.Static("/public", "public")
	e.Static("/upload", "upload")

	// initialitation to use session
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("session"))))

	// routing
	e.GET("/", home);
	e.GET("/home", home);
	e.GET("/contact", contact);
	e.GET("/blog", blog);
	e.GET("/blog-detail/:id", blogDetail);
	e.GET("/login", getLogin);
	e.GET("/register", getRegister);
	e.GET("/logout", logout)
	e.GET("/update-project/:id", updateProject);
	e.GET("/delete/:id", deleteBlog);

	e.POST("/addBlog", middleware.UploadFile(addBlog));
	e.POST("/update-project/:id", middleware.UploadFile(updateProjectPost));
	e.POST("/login", postLogin);
	e.POST("/register", postRegister);

	fmt.Println(sessionUser)

	e.Logger.Fatal(e.Start("Localhost:5000"))
}

func home(c echo.Context) error {
	session, _ := session.Get("session", c);

	flash := map[string]interface{} {
		"FlashId": session.Values["id"],
		"FlashIsLogin": session.Values["isLogin"],
		"FlashStatus": session.Values["status"],
		"FlashName": session.Values["name"],
		"FlashMessage": session.Values["message"],
	}

	fmt.Println(flash["FlashIsLogin"])
	fmt.Println(flash["FlashId"])

	// if session.Values["isLogin"] != true {
	// 	sessionUser.NotLogin = false
	// } else {
	// 	sessionUser.NotLogin = true
	// }
	
	delete(session.Values, "message");
	delete(session.Values, "status");
	session.Save(c.Request(), c.Response())

	template, err := template.ParseFiles("views/index.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	var result []Blog
	if flash["FlashIsLogin"] == true {
		// Query Get Database
		data, err := connection.Conn.Query(context.Background(), "SELECT tb_blog.id, title, description, tb_user.name AS author, image, start_date, end_date, technology FROM tb_blog INNER JOIN tb_user ON tb_blog.author = tb_user.id WHERE tb_blog.author = $1 ORDER BY tb_blog.id DESC", flash["FlashId"]);
	
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}
		
		for data.Next() {
			var temporary = Blog{}
			
			// Scanning 
			err := data.Scan(&temporary.ID, &temporary.Title, &temporary.Description, &temporary.Author, &temporary.Image, &temporary.StartDate, &temporary.EndDate, &temporary.Technology)
			
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
			}	
				
			temporary.Distance = getTime(temporary.StartDate, temporary.EndDate);
			
			// Substring use runes
			runes := []rune(temporary.Description);
			temporary.SubDescription  = string(runes[:130])
			
			/** Substring use ASCII
			temporary.SubDescription = temporary.Description[:147]
			*/
			
			result = append(result, temporary)
		}
	} else {
		data, err := connection.Conn.Query(context.Background(), "SELECT tb_blog.id, title, description, tb_user.name AS author, image, start_date, end_date, technology FROM tb_blog INNER JOIN tb_user ON tb_blog.author = tb_user.id ORDER BY tb_blog.id DESC");
	
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}
		
		// var result []Blog
		
		for data.Next() {
			var temporary = Blog{}
			
			// Scanning 
			err := data.Scan(&temporary.ID, &temporary.Title, &temporary.Description, &temporary.Author, &temporary.Image, &temporary.StartDate, &temporary.EndDate, &temporary.Technology)
			
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
			}	
				
			temporary.Distance = getTime(temporary.StartDate, temporary.EndDate);
			
			// Substring use runes
			runes := []rune(temporary.Description);
			temporary.SubDescription  = string(runes[:130])
			
			/** Substring use ASCII
			temporary.SubDescription = temporary.Description[:147]
			*/
			
			result = append(result, temporary)
		}
	}
	

	Blogs := map[string]interface{}{
		"Title":       "Hi, Welcome To My Hut",
		"Description": "Hello, my name is Ibnu. I completed high school with a focus on natural sciences and am currently studying to become a full-stack developer. I am particularly drawn to JavaScript due to its user-friendly syntax, which aligns well with my goal of becoming a proficient full-stack developer. Additionally you can find a photo of one of my favorite musical artists, Freddie Mercury, the iconic lead vocalist of Queen during the 70s and 80s, on this page site.",
		"Blog":        result,
		"Flash": 	flash,
		"Session": sessionUser,
	}

	return template.Execute(c.Response(), Blogs)
}

func contact(c echo.Context) error {
	session, err := session.Get("session", c);

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	flash := map[string]interface{} {
		"FlashIsLogin": session.Values["isLogin"],
		"FlashName": session.Values["name"],
	}

	template, err := template.ParseFiles("views/contact-me.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return template.Execute(c.Response(), flash)
}

func blog(c echo.Context) error {
	session, _ := session.Get("session", c);

	flash := map[string]interface{} {
		"FlashIsLogin": session.Values["isLogin"],
		"FlashName": session.Values["name"],
	}

	template, err := template.ParseFiles("views/blog.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return template.Execute(c.Response(), flash)
}

func blogDetail(c echo.Context) error {
	session, _ := session.Get("session", c);
	flash := map[string]interface{} {
		"FlashIsLogin": session.Values["isLogin"],
		"FlashName": session.Values["name"],
	}

	id, _ := strconv.Atoi(c.Param("id"))

	template, err := template.ParseFiles("views/blog-detail.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	var BlogData = Blog{}

	// id, title, description, end_date - start_date, start_date, end_date, technology, image, tb_user.name AS author;

	dataErr := connection.Conn.QueryRow((context.Background()), "SELECT tb_blog.id, title, description, end_date - start_date, start_date, end_date, technology, image, tb_user.name AS author FROM tb_blog INNER JOIN tb_user ON tb_blog.author = tb_user.id WHERE tb_blog.id = $1;", id).Scan(&BlogData.ID, &BlogData.Title, &BlogData.Description, &BlogData.Day, &BlogData.StartDate, &BlogData.EndDate, &BlogData.Technology, &BlogData.Image, &BlogData.Author)

	if dataErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": dataErr.Error()})
	}
	// Format date
	var layout = "2 January 2006"
	var startDate =  BlogData.StartDate.Format(layout)
	var endDate =  BlogData.EndDate.Format(layout)

	BlogData.StartDateConv = startDate;
	BlogData.EndDateConv = endDate;

	// Distance
	BlogData.Distance = getTime(BlogData.StartDate, BlogData.EndDate)
	fmt.Println(BlogData.Distance)
	fmt.Println(reflect.TypeOf(BlogData.Distance))
	fmt.Println(BlogData)


	data := map[string]interface{}{
		"Blog": BlogData,
		"Flash": flash,
	}

	return template.Execute(c.Response(), data)
}

func addBlog(c echo.Context) error {
	session, _ := session.Get("session", c);
	author := session.Values["id"];
	
	title := c.FormValue("title");
	startDate := c.FormValue("start-date");
	endDate := c.FormValue("end-date");
	description	:= c.FormValue("description");
	image := c.Get("dataFile").(string)
	nodeJs	:= c.FormValue("node-js");
	reactJs	:= c.FormValue("react-js");
	golang	:= c.FormValue("golang");
	python	:= c.FormValue("python");

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO public.tb_blog(title, description, author, start_date, end_date, technology, image) VALUES($1, $2, $3, $4, $5, Array[$6, $7, $8, $9], $10 )", title, description, author, startDate, endDate, nodeJs, reactJs, golang, python, image)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/home")
}

func updateProject(c echo.Context) error {
	session, _ := session.Get("session", c);
	flash := map[string]interface{} {
		"FlashIsLogin": session.Values["isLogin"],
		"FlashName": session.Values["name"],
	}

	id, _ := strconv.Atoi(c.Param("id"))

	template, err := template.ParseFiles("views/update-project.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": err.Error()})
	}

	var BlogData = Blog{}

	err = connection.Conn.QueryRow((context.Background()), "SELECT id, title, description FROM tb_blog WHERE id=$1;", id).Scan(&BlogData.ID, &BlogData.Title, &BlogData.Description);

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	data := map[string]interface{}{
		"Blog": BlogData,
		"Flash": flash,
	}

	return template.Execute(c.Response(), data)
}

func updateProjectPost(c echo.Context) error {
	err := c.Request().ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	session, _ := session.Get("session", c);
	author := session.Values["id"]

	id, _ := strconv.Atoi(c.Param("id"));
	
	title := c.FormValue("title");
	description := c.FormValue("description");
	startDate :=   c.FormValue("start-date");
	endDate :=     c.FormValue("end-date");
	image := c.Get("dataFile").(string);
	nodeJs :=      c.FormValue("node-js");
	reactJs :=     c.FormValue("react-js");
	golang :=      c.FormValue("golang");
	python :=      c.FormValue("python");

	data, err := connection.Conn.Exec(context.Background(), "UPDATE tb_blog SET title=$1, description=$2, author=$5, start_date=$3, end_date=$4, technology=ARRAY[$6, $7, $8, $9], image=$11 WHERE id=$10", title, description, startDate, endDate, author, nodeJs, reactJs, golang, python, id, image)

	fmt.Println(data)
	fmt.Println(err)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message" : err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/home")
}

func deleteBlog(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id")) // id = 0 string => 0 int

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM public.tb_blog where id=$1", id);

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message" : err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/home")
}

func getLogin(c echo.Context) error {
	session, _ := session.Get("session", c)

	flash := map[string]interface{}{
		"FlashStatus":  session.Values["status"], 
		"FlashMessage": session.Values["message"],    
	}

	delete(session.Values, "message")
	delete(session.Values, "status")
	session.Save(c.Request(), c.Response())

	template, err := template.ParseFiles("views/login.html");

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message" : err.Error()})
	}

	return template.Execute(c.Response(), flash)
}

func getRegister(c echo.Context) error {
	session, _ := session.Get("session", c)

	flash := map[string]interface{}{
		"FlashStatus":  session.Values["status"], 
		"FlashMessage": session.Values["message"],    
	}

	delete(session.Values, "status");
	delete(session.Values, "message");
	session.Save(c.Request(), c.Response())

	template, err := template.ParseFiles("views/register.html");

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message" : err.Error()})
	}

	return template.Execute(c.Response(), flash)
}

func postLogin(c echo.Context) error {
	err := c.Request().ParseForm();

	if err != nil {
		log.Fatal(err)
	}

	email := c.FormValue("email");
	password := c.FormValue("password");

	user := User{};

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, email, password FROM tb_user WHERE email=$1", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password);

	if err != nil {
		// return c.JSON(http.StatusInternalServerError, map[string]string{"message" : err.Error()})
		return redirectWithMessage(c, "Invalid Email!", false, "/login")
	}
	fmt.Println(user);

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		// return c.JSON(http.StatusInternalServerError, map[string]string{"message" : err.Error()})
		return redirectWithMessage(c, "Invalid Password", false, "/login")
	}

	session, _ := session.Get("session", c);
	session.Options.MaxAge = 10800; // 3 jam aja
	session.Values["message"] = "Login Success";
	session.Values["status"] = true;
	session.Values["name"] = user.Name;
	session.Values["id"] = user.ID;
	session.Values["isLogin"] = true;
	session.Values["email"] = user.Email;
	session.Save(c.Request(), c .Response())

	return c.Redirect(http.StatusMovedPermanently, "/home")
}

func postRegister(c echo.Context) error {
	err := c.Request().ParseForm();
	if err != nil {
		log.Fatal(err);
	};

	name := c.FormValue("name");
	email := c.FormValue("email");
	password := c.FormValue("password");

	// email checking
	var count int;
	err = connection.Conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM tb_user WHERE email=$1", email).Scan(&count);
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	fmt.Println(count)

	if count > 0 {
		return redirectWithMessage(c, "Email Already Exists", false, "/register")
	}	
	
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10);

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(name, email, password) VALUES($1, $2, $3)", name, email, passwordHash);

	if err != nil {
		// return c.JSON(http.StatusInternalServerError, map[string]string {"message" : err.Error()})
		return redirectWithMessage(c, "Registration failed, please try again!", false, "/register");
	};
	
	return redirectWithMessage(c, "Registration Success, welcome to our family", true, "/login")
}

func logout(c echo.Context) error {
	session, err :=  session.Get("session", c);

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	session.Options.MaxAge = -1;
	session.Values["isLogin"] = false;
	session.Save(c.Request(), c.Response());
	
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func redirectWithMessage(c echo.Context, message string, status bool, path string) error {
	session, err := session.Get("session", c);

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message" : err.Error()})
	}
	session.Values["message"] = message;
	session.Values["status"] = status;
	session.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, path)
}

func getTime(start time.Time, end time.Time) string {

	distance := int(end.Sub(start).Hours());
	day := distance / 24;
	month := day / 30;
	year := month / 12;
	
	if day > 0 && day <= 29 {
		return fmt.Sprintf("%d Day", day);
	} else if day >= 30 && month < 12 {
		return fmt.Sprintf("%d Month", month);
	} else if month >= 12 {
		return fmt.Sprintf("%d Year", year);
	} else if day >= 0 && distance <= 24 {
		return "1 Day"
	}

	return ""
}