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
    ImageFolder string
    VideoFolder string
    MusicFolder string
}

type Page struct {
    Title string
    CurrentSection string // store name of the currently activated section of the website
}

func SaveConfig (settings Settings) {
  content,err :=json.Marshal(&settings)
  if err!=nil{
      fmt.Print("Error:",err)
  } else {
    err = ioutil.WriteFile("config.json",content,0666)
    if err!=nil{
        fmt.Print("Error:",err)
    }
  }
}

func LoadConfig () Settings {
  var settings Settings = Settings{ FileFolder:"./",ImageFolder:"./",VideoFolder:"./",MusicFolder:"./" }
  content, err := ioutil.ReadFile("config.json")
  if err!=nil{
      fmt.Print("Error:",err)
      fmt.Print("\nSo a new one is being created")
      SaveConfig(settings)
  } else {
    err=json.Unmarshal(content, &settings)
    if err!=nil{
        fmt.Print("Error:",err)
    }
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
  final_path := config.FileFolder + real_path;
  // should do some check on folder to make sure it can't break out of permitted folder
  fmt.Printf("Full path requested:"+final_path)
  // do some check to see if it points to a file or a folder
  if (strings.HasSuffix(final_path,"/")) {
    fmt.Printf("Requested folder listing\n")
    contents_array, err := ioutil.ReadDir(final_path);
    if err != nil {
        //panic(err)
        fmt.Fprintf(w, err.Error())
    } else {
      t := template.Must(template.ParseFiles("templates/filebrowse.html","templates/header.html","templates/footer.html") )  // Parse template file.
      t.Execute(w, contents_array) // note the limitation whereby only one object may be sent to the template
    }
  } else {
    fmt.Printf("Requested file\n")
    http.ServeFile(w, r, final_path)
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
