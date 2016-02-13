package main

// import tells what libraries will be required
import (
  "html/template"
  "net/http"
  "strings"
  "io/ioutil"
  "encoding/json"
  "fmt"
)

type Settings struct {
    FileFolder string
    VideoFolder string
    MusicFolder string
}

type Page struct {
    Title string
    CurrentSection string // store name of the currently activated section of the website
}

func LoadConfig () Settings {
  content, err := ioutil.ReadFile("config.json")
  if err!=nil{
      fmt.Print("Error:",err)
      var TempConfig Settings = Settings{ FileFolder:"./",VideoFolder:"./",MusicFolder:"./" }
      content,err :=json.Marshal(&TempConfig)
      if err!=nil{
          fmt.Print("Error:",err)
          fmt.Print("Failed to serialise config file")
      }
      err = ioutil.WriteFile("config.json",content,0666)
      if err!=nil{
          fmt.Print("Error:",err)
          fmt.Print("Failed to write config file")
      }
  }
  var settings Settings
  err=json.Unmarshal(content, &settings)
  if err!=nil{
      fmt.Print("Error:",err)
  }
  return settings
}

func FolderScan (path string) {

}

func ImageBrowseHandler (w http.ResponseWriter, r *http.Request) {
  fmt.Printf("Image page requested.\n")
}

func FolderBrowseHandler (w http.ResponseWriter, r *http.Request) {
  fmt.Printf("File Browsing page requested.\n")
  full_path := r.URL.Path[1:];
  real_path := strings.TrimPrefix(full_path, "files/");

  //fmt.Fprintf(w, real_path)
  contents_array, err := ioutil.ReadDir(config.FileFolder + real_path);
  if err != nil {
      //panic(err)
      fmt.Fprintf(w, err.Error())
  } else {
    t := template.Must(template.ParseFiles("templates/filebrowse.html","templates/header.html","templates/footer.html") )  // Parse template file.
    t.Execute(w, contents_array) // note the limitation whereby only one object may be sent to the template
  }
}

func HomeHandler (w http.ResponseWriter, r *http.Request) {
  fmt.Printf("Home page requested.\n")
  //fmt.Fprintf(w, "Hi there, I love %s!")
  t := template.Must(template.ParseFiles("templates/main.html","templates/header.html","templates/footer.html") )  // Parse template file.
  t.Execute(w, nil)
}

var config Settings = LoadConfig();

func main() {

  http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
      http.ServeFile(w, r, r.URL.Path[1:])
  })
  http.HandleFunc("/images", ImageBrowseHandler)
  http.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/files/", 301)} )
  http.HandleFunc("/files/", FolderBrowseHandler)
  http.HandleFunc("/", HomeHandler)
  http.ListenAndServe(":3000", nil)
}
