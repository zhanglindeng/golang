package main

import (
    "net/http"
    "os"
    "log"
    "io/ioutil"
)
// <video src="http://127.0.0.1:8000/video" width="520" controls="controls">error</video>
func video(w http.ResponseWriter, r *http.Request) {

    path := "go.mp4"

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

    w.Header().Add("Content-Type", "video/mp4")

    n, err := w.Write(fd)
    if err != nil {
        log.Println(err, "Write")
        http.Error(w, "Internal Server Error: " + err.Error(), http.StatusInternalServerError)
    } else {
        log.Println(n, "Ok")
    }
}
