package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	rand.Seed(time.Now().Unix())
	r := mux.NewRouter()
	r.HandleFunc("/compile", compile)
	os.MkdirAll("/data", 7777)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), r)
	if err != nil {
		log.Fatal("listen: ", err)
	}
}

func compile(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength <= 0 {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Only fixed-size files are allowed")
		return
	}
	if r.ContentLength > 1*1024*1024 {
		w.WriteHeader(400)
		fmt.Fprintf(w, "filesize is > 1 MB")
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "body read: %v", err)
		return
	}
	inFile := fmt.Sprintf("/data/%v_%v.lua", rand.Int63(), time.Now().Unix())
	outFile := fmt.Sprintf("/data/%v_%v_out.luac", rand.Int63(), time.Now().Unix())
	err = ioutil.WriteFile(inFile, data, 7777)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "write file: %v", err)
		return
	}
	cmd := exec.Command("/app/luac.cross", "-f", "-s", "-o", outFile, inFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "luac.cross: %v %v", err, string(out))
		return
	}
	data, err = ioutil.ReadFile(outFile)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "read file: %v", err)
		return
	}
	w.Write(data)
}
