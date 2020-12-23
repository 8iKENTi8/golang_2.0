package main

import(
  "fmt"
  "net/http"
  "github.com/gorilla/mux"
  "html/template"
  "database/sql"
  _"github.com/go-sql-driver/mysql"
)

type Article struct{
  Id uint16
  Title, Anons, Full_txt string
}

var posts = []Article{}
var showPost = Article{}

func select_InfoDb(){

  db,err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
  if err != nil{
    panic(err)
  }

  defer db.Close()

  //выборка данных из базы данных
  res, err := db.Query("SELECT * FROM `articles`")
  if err != nil{
    panic(err)
  }
  posts = []Article{}
  for res.Next(){
    var post Article
    err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Full_txt)

    if err != nil{
      panic(err)
    }

    posts = append(posts, post)

  }
}

func index(w http.ResponseWriter, r *http.Request){
  t, err := template.ParseFiles("templates/index.html",
    "templates/header.html","templates/footer.html")

  if err != nil{
    fmt.Fprintf(w, err.Error())
  }

  select_InfoDb()

  t.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request){
  t, err := template.ParseFiles("templates/create.html",
    "templates/header.html","templates/footer.html")

  if err != nil{
    fmt.Fprintf(w, err.Error())
  }

  t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request){

  title := r.FormValue("title")
  anons := r.FormValue("anons")
  full_txt := r.FormValue("full_txt")

  if title == "" || anons == "" || full_txt == ""{
    fmt.Fprintf(w, "Не все данные заполнены")
    } else{

      db,err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
      if err != nil{
        panic(err)
      }

      defer db.Close()

      //Занесение данныз в табилцу

      insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles`" +
        "(`id`, `title`, `anons`, `full_txt`)"+
        "VALUES (NULL, '%s', '%s', '%s')",title, anons, full_txt))
        if err != nil {
          panic(err)
        }

        defer insert.Close()

        http.Redirect(w, r, "/", http.StatusSeeOther)
      }
    }

    

func show_post(w http.ResponseWriter, r *http.Request){
vars := mux.Vars(r)

t, err := template.ParseFiles("templates/show.html",
  "templates/header.html","templates/footer.html")

db,err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
if err != nil{
  panic(err)
}

defer db.Close()

//выборка данных из базы данных
res, err := db.Query(fmt.Sprintf("SELECT * FROM"+
  "`articles` WHERE `id` = '%s'",vars["id"]))
if err != nil{
  panic(err)
}
showPost = Article{}
for res.Next(){
  var post Article
  err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Full_txt)

  if err != nil{
    panic(err)
  }

showPost = post
}

  t.ExecuteTemplate(w, "show", showPost)
}

func handleFunc()  {

      rtr := mux.NewRouter()

      rtr.HandleFunc("/" , index).Methods("GET")
      rtr.HandleFunc("/create" , create).Methods("GET")
      rtr.HandleFunc("/save_article" , save_article).Methods("POST")
      rtr.HandleFunc("/post/{id:[0-9]+}" , show_post).Methods("GET")


      http.Handle("/", rtr)
      http.Handle("/static/",
        http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
      http.ListenAndServe(":8080", nil)
    }

func main(){
      handleFunc()
    }
