package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/labstack/echo"
)

type WebServer struct {
	path            string
	templates       *template.Template
	echo            *echo.Echo
	isRequireReload bool
	sync.Mutex
}

func NewWebServer(echo *echo.Echo, path string) *WebServer {
	web := &WebServer{
		echo: echo,
		path: path,
	}

	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		WebPath, err := filepath.Abs(path)
		if err != nil {
			log.Fatalln(err)
		}

		NewFileWatcher(WebPath, func(ev string, path string) {
			if strings.HasPrefix(filepath.Ext(path), ".htm") {
				web.isRequireReload = true
			}
		})
	}
	web.UpdateRender()

	return web
}

func (web *WebServer) CheckWatch() {
	if web.isRequireReload {
		web.Lock()
		if web.isRequireReload {
			err := web.UpdateRender()
			if err != nil {
				log.Println(err)
			} else {
				web.isRequireReload = false
			}
		}
		web.Unlock()
	}
}

func (web *WebServer) UpdateRender() error {
	tp := template.New("").Delims("<%", "%>").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	})
	filepath.Walk(web.path, func(path string, fi os.FileInfo, err error) error {
		if strings.HasPrefix(filepath.Ext(path), ".htm") {
			rel, err := filepath.Rel(web.path, path)
			if err != nil {
				return err
			}
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			rel = filepath.ToSlash(rel)
			template.Must(tp.New(rel).Parse(string(data)))
		}
		return nil
	})
	web.templates = tp

	return nil
}

func (web *WebServer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return web.templates.ExecuteTemplate(w, name, data)
}

func (web *WebServer) SetupStatic(e *echo.Echo, prefix string, root string) {
	h := func(c echo.Context) error {
		fname := c.Param("*")
		fpath := path.Join(root, fname)
		return c.File(fpath)
	}
	e.GET(prefix, h)
	if prefix == "/" {
		e.GET(prefix+"*", h)
	} else {
		e.GET(prefix+"/*", h)
	}
}
