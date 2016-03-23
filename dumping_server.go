package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port, p",
			Value: 8080,
			Usage: "Which port to listen on for HTTP traffic",
		},
		cli.IntFlag{
			Name:  "https",
			Value: 9090,
			Usage: "Which port to listen on for HTTPS traffic",
		},
		cli.StringFlag{
			Name:  "key, k",
			Usage: "A key file to use",
		},
		cli.StringFlag{
			Name:  "cert, c",
			Usage: "A cert file to use",
		},
	}
	app.Action = start

	app.Run(os.Args)
}

// simply starts listening to http traffic and dumping the info to stdout
func start(c *cli.Context) {
	http.HandleFunc("/", dump)
	if c.GlobalString("cert") != "" || c.GlobalString("key") != "" {
		go startHTTPSServer(
			c.GlobalString("cert"),
			c.GlobalString("key"),
			c.GlobalInt("https"))
	}
	port := c.GlobalInt("port")
	fmt.Printf("Starting HTTP server on %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func dump(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Failed to read input: %v\n", err)
		return
	}

	fmt.Println("-- headers -- ")
	for k, v := range r.Header {
		fmt.Printf("%s: %v\n", k, v)
	}

	fmt.Println("-- body --")
	if r.Header.Get("Content-Type") == "application/json" {
		var out bytes.Buffer
		json.Indent(&out, body, "", "  ")
		out.WriteTo(os.Stdout)
		fmt.Println("")
	} else {
		fmt.Println(string(body))
	}
	fmt.Println("----------")
}

func startHTTPSServer(certPath, keyPath string, port int) {
	if certPath == "" {
		fmt.Println("No Cert file provided.")
		return
	}
	if keyPath == "" {
		fmt.Println("No Key file provided")
		return
	}
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		fmt.Println("Can't find key file at: " + keyPath)
		return
	}
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		fmt.Println("Can't find cert file at: " + keyPath)
		return
	}

	fmt.Printf("Starting HTTPS server on %d\n", port)
	http.ListenAndServeTLS(fmt.Sprintf(":%d", port), certPath, keyPath, nil)
}
