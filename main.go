package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"time"
)

// FileEntry represents a single uploaded file
type FileEntry struct {
	Id         int
	Name       string
	Size       int // in bytes
	UploadDate time.Time
}

type View struct {
	viewid int
	files  []FileEntry
	token  string
}

var (
	files    = []FileEntry{}
	basePath = "/tmp/filebox"
)

func main() {
	files = append(files, FileEntry{1, "foo", 100, time.Now()})
	files = append(files, FileEntry{2, "bar", 200, time.Now()})

	e := echo.New()
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(10)))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Renderer = &Template{templates: template.Must(template.ParseGlob("templates/*.html"))}

	e.GET("/files/:fileid", downloadFile)
	e.DELETE("/files/:fileid", deleteFile)
	e.POST("/files", uploadFile)
	e.GET("/files", homeView)

	e.GET("/views/:viewid", getView)
	e.DELETE("/views/:viewid", deleteView)
	e.POST("/views", createView)
	e.GET("/views", allViews)

	e.GET("/", homeView)

	e.File("/js/main.js", "static/main.js")
	e.File("/favicon.ico", "static/favicon.ico")

	e.Logger.Fatal(e.Start(":8080"))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(
	w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func downloadFile(c echo.Context) error {
	id := c.Param("fileid")
	return c.File(filepath.Join(basePath, id))
}

func uploadFile(c echo.Context) error {
	log.Error("Not implemented")
	return c.Render(http.StatusCreated, "uploadSuccess.html", files)
}

func deleteFile(c echo.Context) error {
	return c.String(http.StatusAccepted, "hi")
}

func allViews(c echo.Context) error {
	return c.String(http.StatusAccepted, "hi")
}

func homeView(c echo.Context) error {
	return c.Render(http.StatusOK, "filesview.html", files)
}

func getView(c echo.Context) error {
	return c.String(http.StatusAccepted, "hi")
}

func createView(c echo.Context) error {
	return c.String(http.StatusAccepted, "hi")
}

func deleteView(c echo.Context) error {
	return c.String(http.StatusAccepted, "hi")
}
