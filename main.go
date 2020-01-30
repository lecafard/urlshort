package main

import (
    "fmt"
    "net/http"
    "log"
    "encoding/json"
    "strings"
    "sort"
    "io/ioutil"
    "net/url"
    "math/rand"
    "time"

    "github.com/bmizerany/pat"
    "github.com/boltdb/bolt"
    "github.com/spf13/viper"
)

var URL_BUCKET = []byte("urls")
// Alphabetical list (sort later)
var BLACKLIST = []string{"api", "delete", "info", "list", "robots.txt", "shorten", "static", "stats"}
var ALPHANUM = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
const RAND_LENGTH = 6

// Does not have to be secure, can be shitty
func RandomString(n int, runes []rune) string {
    lr := len(runes)
    o := make([]rune, n)
    for i := 0; i < n; i++ {
        o[i] = runes[rand.Intn(lr)]
    }

    return string(o)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(404)
    fmt.Fprint(w, "URL Not Found\n")
}

func ShortenerUI(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/shorten.html")
}

func Robots(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w,"User-agent: *\nDisallow /\n")
}

func Lengthen(w http.ResponseWriter, r *http.Request) {
    params := r.URL.Query()
    short := params.Get(":short")
    long := db.Get(short)

    if long == "" {
        NotFound(w, r);
    } else {
        http.Redirect(w, r, long, 301);
    }
}

func List(w http.ResponseWriter, r *http.Request) {
    response, err := json.Marshal(db.List())
    if err != nil {
        log.Println("[ERROR]", err)
        ServerError(w, r)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    if string(response) == "null" {
        w.Write([]byte("[]"))
    } else {
        w.Write(response)
    }
}

func Delete(w http.ResponseWriter, r *http.Request) {
    params := r.URL.Query()
    short := params.Get(":short")
    err := db.Delete(short)
    if err != nil {
        log.Println("[ERROR]", err)
        ServerError(w, r)
    } else {
        w.WriteHeader(204)
    }
}

func Shorten(w http.ResponseWriter, r *http.Request) {
    var u URLRecord
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Println("[ERROR]", err)
        ServerError(w, r)
        return
    }
    err = json.Unmarshal(body, &u)
    // Check if no other short urls exist
    if err != nil {
        // Malformed Request
        w.WriteHeader(400);
        w.Header().Set("Content-Type", "application/json");
        w.Write([]byte("{\"status\":false,\"error\":\"MalformedRequest\"}"));
        return
    }
    checked := u.Overwrite || false
    // Check if url is valid
    _, err = url.ParseRequestURI(u.Long)
    if err != nil {
        log.Println("[ERROR]", err)
        w.WriteHeader(400)
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte("{\"status\":false,\"error\":\"InvalidURL\"}"))
        return
    }

    if u.Short != "" {
        ul := strings.ToLower(u.Short)
        idx := sort.SearchStrings(BLACKLIST, ul)

        // Check if short url is valid
        if strings.ContainsAny(u.Short, "~`!@#$%^&*()+={}[]|\\:;\"'<,>.?/") ||
            (idx < len(BLACKLIST) && BLACKLIST[idx] == ul) {
            w.WriteHeader(400)
            w.Header().Set("Content-Type", "application/json")
            w.Write([]byte("{\"status\":false,\"error\":\"InvalidShortURL\"}"))
            return
        }
    } else {
        u.Short = RandomString(RAND_LENGTH, ALPHANUM)
        checked = true
    }

    if checked == false {
        if db.Get(u.Short) != "" {
            w.WriteHeader(400)
            w.Header().Set("Content-Type", "application/json")
            w.Write([]byte("{\"status\":false,\"error\":\"ExistsShortURL\"}"))
            return
        }
    }

    err = db.Set(u.Short, u.Long)
    if err != nil {
        ServerError(w, r)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte("{\"status\":true,\"short\":\""+ u.Short +"\"}"))
}

func ServerError(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(500)
    fmt.Fprint(w, "Internal Server Error\n")
}

type URLRecord struct {
    Short string `json:"short"`
    Long string `json:"long"`
    Overwrite bool `json:"overwrite,omitempty"`
}

type URLStore struct {
    DB *bolt.DB
}

func (s *URLStore) Open(path string) error {
    db, err := bolt.Open(path, 0600, nil)
    // DB Startup create bucket
    db.Update(func (tx *bolt.Tx) error {
        _, err := tx.CreateBucketIfNotExists(URL_BUCKET)
        if err != nil {
            panic(err)
        }
        return err
    })

    s.DB = db
    return err
}

func (s *URLStore) Close() error {
    return s.DB.Close()
}

func (s *URLStore) Get(url string) (res string) {
    key := []byte(url)
    s.DB.View(func (tx *bolt.Tx) error {
        bucket := tx.Bucket(URL_BUCKET)
        res = string(bucket.Get(key))
        return nil
    })
    return
}

func (s *URLStore) Delete(url string) (err error) {
    key := []byte(url)
    s.DB.Update(func (tx *bolt.Tx) error {
        bucket := tx.Bucket(URL_BUCKET)
        return bucket.Delete(key)
    })
    return
}

func (s *URLStore) Set (short string, long string) (err error) {
    s.DB.Update(func (tx *bolt.Tx) error {
        bucket := tx.Bucket(URL_BUCKET)
        return bucket.Put([]byte(short), []byte(long))
    })
    return
}

func (s *URLStore) List() (list []URLRecord) {
    s.DB.View(func (tx *bolt.Tx) error {
        bucket := tx.Bucket(URL_BUCKET)
        cursor := bucket.Cursor()
        // Iterate over bucket
        for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
            list = append(list, URLRecord{string(key), string(value), false});
        }

        return nil
    })
    return
}

var db URLStore


func main() {
    viper.SetDefault("data", "data.db")
    viper.SetDefault("listen", ":8000")

    viper.SetEnvPrefix("urlshort")
    viper.BindEnv("data")
    viper.BindEnv("listen")


    dbpath := viper.GetString("data")
    host := viper.GetString("listen")

    // Generate random(ish) seed
    rand.Seed(time.Now().UTC().UnixNano())

    // Initialise database
    err := db.Open(dbpath)
    if err != nil {
        panic(err)
    }

    mux := pat.New()
    mux.Get("/api/list", http.HandlerFunc(List))
    mux.Post("/api/shorten", http.HandlerFunc(Shorten))
    mux.Del("/api/delete/:short", http.HandlerFunc(Delete))
    mux.Get("/static/*", http.FileServer(http.Dir("./static")))
    mux.Get("/", http.HandlerFunc(NotFound))
    mux.Get("/robots.txt", http.HandlerFunc(Robots))
    mux.Get("/shorten", http.HandlerFunc(ShortenerUI))
    mux.Get("/:short", http.HandlerFunc(Lengthen))

    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static", fs))
    http.Handle("/", mux)
    log.Fatal(http.ListenAndServe(host, nil))
}
