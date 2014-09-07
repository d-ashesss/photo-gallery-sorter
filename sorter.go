package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"encoding/json"
	"flag"
	"os"
	"os/exec"
	"net/http"
	"regexp"
	"./tpl"
)

var (
	availableExtensions = []string{"jpg", "jpeg", "png", "gif", "apng", "agif", "swf"}
)

const (
	noTagDir = "_"
)

func init() {
	var extensions string
	var target string

	flag.StringVar(&extensions, "e", "", "Show only specified file types")
	flag.StringVar(&target, "t", "", "Directory to work in")
	flag.Parse()

	if len(extensions) > 0 {
		availableExtensions = strings.Split(extensions, ",")
	}

	if len(target) > 0 {
		if err := os.Chdir(target); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/image/", imageHandler)

	http.HandleFunc("/images", imagesHandler)
	http.HandleFunc("/tags", tagsHandler)

	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/close", closeHandler)

	cmd := exec.Command("C:\\Program Files (x86)\\Opera\\opera.exe", "localhost:8080")
	cmd.Run()

	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	template := tpl.GetTemplate()
	template = strings.Replace(template, "%%NO_TAG_DIR%%", noTagDir, 1)
	io.WriteString(w, template)
}

func imageHandler(w http.ResponseWriter, req *http.Request) {
	pathInfo := strings.SplitN(req.URL.Path, "/", 3)

	fd, err := os.Open(pathInfo[2])

	if err != nil {
		fmt.Println("images handler:", err)
		return
	}

	io.Copy(w, fd)
	fd.Close()
}

func imagesHandler(w http.ResponseWriter, req *http.Request) {
	images, err := getImages()

	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.Write(images)
}

func tagsHandler(w http.ResponseWriter, req *http.Request) {
	tags, err := getTags()

	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.Write(tags)
}

func saveHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	for file, tag := range req.Form {
		dir := strings.Join(tag, "")

		_, err := os.Stat(file)

		if err != nil {
			fmt.Println(err)
			continue
		}

		err = mkdir(dir)

		if err != nil {
			fmt.Println(err)
			continue
		}

		os.Rename(file, dir + "/" + file)
	}
}

func closeHandler(w http.ResponseWriter, req *http.Request) {
	os.Exit(0)
}

func getImages() ([]byte, error) {
	files, err := ioutil.ReadDir(".")

	if err != nil {
		return nil, err
	}

	images := []string{}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		nameParts := strings.Split(file.Name(), ".")
		ext := nameParts[len(nameParts) - 1]

		if !inArray(ext, availableExtensions) {
			continue
		}

		images = append(images, file.Name())
	}

	if len(images) == 0 {
		return []byte("[]"), nil
	}

	return json.Marshal(images)
}

func getTags() ([]byte, error) {
	preTags := []string{}

	for _, tag := range flag.Args() {
		preTags = append(preTags, tag)
	}

	wd, _ := os.Getwd()
	fi, _ := os.Stat(wd)
	dirTags := strings.Split(fi.Name(), " ")
	preTags = append(preTags, dirTags...)

	dirs, err := ioutil.ReadDir(".")
	if err == nil {
		for _, dir := range dirs {
			if dir.IsDir() {
				preTags = append(preTags, dir.Name())
			}
		}
	}

	tags := []string{}
	tagExp, _ := regexp.Compile("^[a-zA-Z0-9]+$")

	for _, tag := range preTags {
		if tagExp.MatchString(tag) {
			tags = append(tags, tag)
		}
	}

	return json.Marshal(uniq(tags))
}

func inArray(needle string, haystack []string) bool {
	for _, val := range haystack {
		if strings.ToLower(val) == strings.ToLower(needle) {
			return true
		}
	}

	return false
}

func mkdir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if err := os.Mkdir(path, 0755); err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func uniq(strings []string) []string {
	result := []string{}

	for _, str := range strings {
		if !inArray(str, result) {
			result = append(result, str)
		}
	}

	return result
}
