package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/oasis14/GophercisesUrlShorten"
)

func main() {
	yamlFile := flag.String("yamlFile", "", "Yaml file name used to set up routes")
	jsonFile := flag.String("jsonFile", "", "JSON file name used to set up routes")
	flag.Parse()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	yamlHandler, err := processYaml(*yamlFile, mapHandler)
	if err != nil {
		panic(err)
	}

	//Json File handler
	JSONHandler, err := processJSON(*jsonFile, yamlHandler)
	//we can technically ignore these error
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", JSONHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func readFile(fileName string) ([]byte, error) {
	//open file
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func processYaml(fileName string, mapHandler http.HandlerFunc) (http.HandlerFunc, error) {
	//no file name specified return the passed in mapHandler
	if fileName == "" {
		return mapHandler, nil
	}

	//file does exist open file
	yaml, err := readFile(fileName)
	if err != nil {
		fmt.Printf("Error reading file: %s. Yaml extensions wont be avaliable\n", fileName)
		return mapHandler, err
	}

	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}
	return yamlHandler, err
}

func processJSON(fileName string, mapHandler http.HandlerFunc) (http.HandlerFunc, error) {
	//no file name specified return the passed in mapHandler
	if fileName == "" {
		return mapHandler, nil
	}

	//file does exist open file
	json, err := readFile(fileName)
	if err != nil {
		fmt.Printf("Error reading file: %s. JSON extensions wont be avaliable\n", fileName)
		return mapHandler, err
	}

	jsonHandler, err := urlshort.JSONHandler([]byte(json), mapHandler)
	if err != nil {
		panic(err)
	}
	return jsonHandler, err
}
