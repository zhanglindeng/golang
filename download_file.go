package main

import (
    "net/http"
    "log"
    "io"
    "os"
    "io/ioutil"
    "strconv"
)

func main() {
    http.HandleFunc("/", index)
    http.HandleFunc("/download", download)

    err := http.ListenAndServe(":8000", nil)
    if err != nil {
        log.Fatalln(err)
    }
}

func index(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "hello")
}

func download(w http.ResponseWriter, r *http.Request) {
    path := "go.zip"

    file, err := os.Open(path)
    if err != nil {
        log.Println(err, "Open")
        http.Error(w, "Internal Server Error: " + err.Error(), http.StatusInternalServerError)
        return
    }
    defer file.Close()

    fd, err := ioutil.ReadAll(file)
    if err != nil {
        log.Println(err, "ReadAll")
        http.Error(w, "Internal Server Error: " + err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Add("Content-Description", "File Transfer")
    w.Header().Add("Content-Type", "application/zip")
    w.Header().Add("Content-Disposition", "attachment;filename=" + path)
    w.Header().Add("Expires", "0")
    w.Header().Add("Cache-Control", "must-revalidate")
    w.Header().Add("Pragma", "public")
    w.Header().Add("Content-Length", strconv.Itoa(len(fd)))

    n, err := w.Write(fd)
    if err != nil {
        log.Println(err, "Write")
        http.Error(w, "Internal Server Error: " + err.Error(), http.StatusInternalServerError)
    } else {
        log.Println(n, "Ok")
    }

}
