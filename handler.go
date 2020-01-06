package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	//pares yaml
	var pathUrls []pathURL
	err := parseYaml(yml, &pathUrls)
	if err != nil {
		return nil, err
	}

	//convert yaml to map
	pathsToUrls := pathURLToMap(&pathUrls)

	//return map handler using MapHandler
	return MapHandler(pathsToUrls, fallback), nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
//    [ {"path": "/SomePath", "url": "https://..."},]
//
// The only errors that can be returned all related to having
// invalid JSON data.
func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathUrls []pathURL
	//parse json
	err := parseJSON(jsn, &pathUrls)
	if err != nil {
		return nil, err
	}
	//convert json to mapping
	pathsToUrls := pathURLToMap(&pathUrls)
	//return map handler
	return MapHandler(pathsToUrls, fallback), nil
}

func parseYaml(yml []byte, pathUrls *[]pathURL) error {
	err := yaml.Unmarshal(yml, pathUrls)
	if err != nil {
		return err
	}
	return nil
}

func parseJSON(jsn []byte, pathUrls *[]pathURL) error {
	err := json.Unmarshal(jsn, pathUrls)
	if err != nil {
		return err
	}
	return nil
}

func pathURLToMap(pathUrls *[]pathURL) map[string]string {
	ret := make(map[string]string)
	for _, pu := range *pathUrls {
		ret[pu.Path] = pu.URL
	}
	return ret
}

type pathURL struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}
