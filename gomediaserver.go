package main

// import tells what libraries will be required
import (
  "runtime/debug" // used for playing with the gc
  "time"

  "html/template"
  "net/http"
  "strings"
  "io/ioutil"
  "encoding/json"
  "path/filepath"
  "fmt"
  "os"
  "image"                  // this
  "image/png"             // this
  "image/jpeg"             // this
  "github.com/nfnt/resize" // and this are used for thumbnail generation
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

// prioritise peformance over memory usage
func FolderScan (path string,extensions []string) []os.FileInfo {
  new_pathlist := make([]os.FileInfo,0,0); // the 0 is length, the 256 is capacity that can be filled before a resize is required
  // first we get a list of everything in the directory
  contents_array, err := ioutil.ReadDir(path);
  if err != nil {
      //panic(err)
      fmt.Print("Error:",err)
  } else {
    new_pathlist = make([]os.FileInfo,0, len(contents_array) ); // the 0 is length, the 256 is capacity that can be filled before a resize is required
    for contents_array_index := range contents_array {
      match_found := false;
      for extension_index := range extensions {
        if filepath.Ext(contents_array[contents_array_index].Name()) == extensions[extension_index] {
          match_found = true;
        }
      }
      if match_found {
        new_pathlist = append(new_pathlist, contents_array[contents_array_index]) // wow this is dumb syntax
      }
    }
  }
  return new_pathlist
}

// this converts pngs and jpgs to smaller jpgs, and writes them out to the http connection, this thing EATS ALL THE MEMORY
func generateThumb (w http.ResponseWriter, path string) {
  file, err := os.Open(path)
  if err != nil {
    fmt.Print("Error opening image:",err)
  }
  var img image.Image = nil;
  if filepath.Ext(path) == ".png" {
    img, err = png.Decode(file)
    if err != nil {
        fmt.Print("Error decoding image:",err)
        return
    }
    file.Close()
  } else if filepath.Ext(path) == ".jpg" {
    img, err = jpeg.Decode(file)
    if err != nil {
        fmt.Print("Error decoding image:",err)
        return
    }
    file.Close()
  } else {
    return
  }
  // resize to height 200
  err = jpeg.Encode(w, resize.Resize(0, 200, img, resize.NearestNeighbor), nil)
  if err != nil {
    fmt.Print("Error encoding thumb:",err)
  }
}

func ImageBrowseHandler (w http.ResponseWriter, r *http.Request) {
  full_path := r.URL.Path[1:];
  real_path := strings.TrimPrefix(full_path, "images/");
  final_path := config.ImageFolder + real_path;
  fmt.Printf("Full IMAGE path requested: "+final_path+"\n")
  if (strings.HasSuffix(final_path,"/")) {
    fmt.Printf("Requested image listing\n")
    // its actually a damn folder
    scanned_files := FolderScan(final_path, []string{ ".png",".jpg" } )
    fmt.Printf("Num images found: ", len(scanned_files) )
    t := template.Must(template.ParseFiles("templates/imagebrowse.html","templates/header.html","templates/footer.html") )  // Parse template file.
    t.Execute(w, scanned_files)
  } else { /* else if (strings.HasSuffix(final_path,"/thumb")) {} */
    fmt.Printf("Requested file\n")
    generateThumb(w,final_path);

    if time.Now().After(timeSinceMemFreed) { // run a custom little memory freeing loop, required because the thumb generator eats memory too fast for the gc
        debug.FreeOSMemory()
        fmt.Printf("<<<<<<< MEMORY FREED >>>>>>>\n")
        timeSinceMemFreed = time.Now().Add(time.Duration(5*time.Second))
    }

  }
}

// show a list of all songs found
func SongMusicView () {

}

// show list of albumns
func AlbumMusicView () {

}

// show list of playlists
func PlaylistMusicView () {

}

// this function offers various ways to enjoy music with ways to view and sort your music collection
func MusicBrowseHandler (w http.ResponseWriter, r *http.Request) {
  fmt.Printf("Music Browsing page requested.\n")
  scanned_files := FolderScan(config.FileFolder, []string{ ".png",".jpg" } )
  // TODO: should do some check on folder to make sure it can't break out of permitted folder

  t := template.Must(template.ParseFiles("templates/musicbrowse.html","templates/header.html","templates/footer.html") )  // Parse template file.
  t.Execute(w, scanned_files)
}

// this function is a bit shit, there has to be a better more idiomatic way
func FolderBrowseHandler (w http.ResponseWriter, r *http.Request) {
  fmt.Printf("File Browsing page requested.\n")
  full_path := r.URL.Path[1:];
  real_path := strings.TrimPrefix(full_path, "files/");
  final_path := config.FileFolder + real_path;
  // TODO: should do some check on folder to make sure it can't break out of permitted folder
  fmt.Printf("Full path requested: "+final_path+"\n")
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
    http.ServeFile(w, r, final_path) // consider using http.Dir to fix issues with browsing places that you shouldnt
  }
}

func HomeHandler (w http.ResponseWriter, r *http.Request) {
  fmt.Printf("Home page requested.\n")
  t := template.Must(template.ParseFiles("templates/main.html","templates/header.html","templates/footer.html") )  // Parse template file.
  t.Execute(w, nil)
}

var config Settings = LoadConfig();
var timeSinceMemFreed time.Time = time.Now().Add(time.Duration(5));

func main() {
  fmt.Printf("Go media server running on port 3000.\n")
  //http.Handle("/static/", http.FileServer(http.Dir("./static")) )
  http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
      http.ServeFile(w, r, r.URL.Path[1:])
  })
  http.HandleFunc("/images/", ImageBrowseHandler)
  http.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/files/", 301)} )
  http.HandleFunc("/files/", FolderBrowseHandler) // might be worth using stripprefix
  http.HandleFunc("/music",  MusicBrowseHandler)
  http.HandleFunc("/", HomeHandler)
  http.ListenAndServe(":3000", nil)
}
