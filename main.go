package main

import(
  "fmt"
  "net/http"
  "html/template"
  "database/sql"
  _"github.com/go-sql-driver/mysql"
)

type Article struct{
  Id uint16
  Title, Anons, Full_txt string
}

var posts = []Article{}

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

func handleFunc()  {
  http.Handle("/static/",
  http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
  http.HandleFunc("/" , index)
  http.HandleFunc("/create" , create)
  http.HandleFunc("/save_article" , save_article)
  http.ListenAndServe(":8080", nil)
}

func main(){
  handleFunc()
}
