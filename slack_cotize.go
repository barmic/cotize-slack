package main

import (
	"bufio"
	"fmt"
	"os"
	"io"
	"log"
	"encoding/json"
	"net/http"
	"strings"
	"path/filepath"
)

var runtimePath string
var tokenSlack string

type SlackResponse struct {
	ResponseType string `json:"response_type"`
	Text string `json:"text"`
}

/*
 * Commands
 */
func getVersion(env string) string {
	log.Print("env : " + env)
	realPath, _ := filepath.EvalSymlinks(runtimePath + env)
	_, file := filepath.Split(realPath)
	return file
}

/*
 * Web functions
 */
func parseInput(input io.ReadCloser) map[string]string {
	rd := bufio.NewReader(input)
	str, err := rd.ReadString('\n')
	params := make(map[string]string)
	for ; err == nil; str, err = rd.ReadString('\n') {
		s := strings.Split(str, "=")
		params[s[0]] = strings.Trim(s[1], "\n")
	}
	return params
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func mepHandler(w http.ResponseWriter, r *http.Request) {
	params := parseInput(r.Body)
	if params["token"] != tokenSlack {
		log.Print("The token " + params["token"] + " isn't good (" + tokenSlack + ")")
		return
	}
	//for key, value := range params {
	//	fmt.Fprintf(w, "Key: %s; Value %s\n", key, value)
	//}

	//fmt.Fprintf(w, "The command is: %s !\n", params["command"])
	var response SlackResponse
	response.ResponseType = "in_channel"

	switch params["command"] {
	case "/versions":
		test := getVersion("test")
		prod := getVersion("production")
		response.Text = "Version de test -> " + test + " | version de production -> " + prod
	case "/mep":
		fmt.Fprintf(w, "Je ne sais pas encore faire de mise en production\n")
	default:
		fmt.Fprintf(w, "Je ne connais pas la commande %s\n", params["command"])
	}
	responseB, _ := json.MarshalIndent(response, "", "    ")
	responseStr := string(responseB)
	log.Print(responseStr)
	fmt.Fprintf(w, responseStr)
	// chose compute by command
}

func main() {
	// configuration
	tokenSlack = os.Args[1]
	runtimePath = os.Args[2]

	// run web server
	http.HandleFunc("/mep", mepHandler)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
