package main

import (
	"fmt"
	"net/http"
	"os/exec"
)

func main() {

	http.HandleFunc("/", page)

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		add(w, r)
	})

	http.HandleFunc("/remove", func(w http.ResponseWriter, r *http.Request) {
		remove(w, r)
	})

	http.ListenAndServe(":9999", nil)

}

func page(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "<h1>Add container</h1>"+
		"<form action=\"/add\" method=\"POST\">"+
		"Id:"+
		"<textarea name=\"Post_id\"></textarea><br>"+
		"Port:"+
		"<textarea name=\"Post_port\"></textarea><br>"+
		"<input type=\"submit\" value=\"Submit\">"+
		"</form>"+

		"<h1>Delete container</h1>"+
		"<form action=\"/remove\" method=\"POST\">"+
		"Id:"+
		"<textarea name=\"Get_id\"></textarea><br>"+
		"<input type=\"submit\" value=\"Submit\">"+
		"</form>"+

		"</form>")

}

func add(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<p><a href=\"/\">go back</a></p>")
	fmt.Fprintf(w, "You have added a node to network")

	id := r.FormValue("Post_id")
	port := r.FormValue("Post_port")

	cmdStr := "sudo docker run -d -e ID=" + id + " -e PORT=" + port + " -p " + port + ":" + port + " -p " + port + ":" + port + "/udp master:ENV"

	fmt.Println(cmdStr)

	exec.Command("/bin/sh", "-c", cmdStr).Output()

}

func remove(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<p><a href=\"/\">go back</a></p>")
	fmt.Fprintf(w, "You have removed a node from network")

	id := r.FormValue("Get_id")

	cmdStr := "sudo docker rmi " + id

	fmt.Println(cmdStr)

	exec.Command("/bin/sh", "-c", cmdStr).Output()

}
