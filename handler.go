package urlshort

import (
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(rw, r, dest, http.StatusFound)
		}

		fallback.ServeHTTP(rw, r)
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
	y, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}

	pathMap := buildMap(y)

	return MapHandler(pathMap, fallback), nil
}

func buildMap(parsed []yamlStruct) map[string]string {
	mappings := make(map[string]string)
	for _, v := range parsed {
		mappings[v.Path] = v.Url
	}
	return mappings
}

func parseYaml(yml []byte) ([]yamlStruct, error) {
	y := []yamlStruct{}
	err := yaml.Unmarshal(yml, &y)
	if err != nil {
		return []yamlStruct{}, err
	}
	return y, nil
}

type yamlStruct struct {
	Path string
	Url  string
}
