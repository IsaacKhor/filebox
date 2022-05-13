package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func main() {
	log.SetLevel(log.DEBUG)

	ech := echo.New()
	ech.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(10)))
	ech.Use(middleware.Logger())
	ech.Use(middleware.Recover())

	t := template.New("")
	t.Funcs(template.FuncMap{"ToBinarySuffix": ToBinarySuffix})
	template.Must(t.ParseGlob("templates/*.html"))
	ech.Renderer = &Template{templates: t}

	ech.GET("/files/:fileid", downloadFile)
	ech.DELETE("/files/:fileid", deleteFile)
	ech.POST("/files", uploadFile)
	ech.GET("/files", homeView)

	ech.GET("/views/:viewid", getView)
	ech.DELETE("/views/:viewid", deleteView)
	ech.POST("/views", createView)
	ech.GET("/views", allViews)

	ech.GET("/", homeView)

	ech.File("/js/main.js", "static/main.js")
	ech.File("/favicon.ico", "static/favicon.ico")

	// Start server
	go func() {
		if err := ech.Start(":8080"); err != nil {
			ech.Logger.Fatal(err)
		}
	}()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)

	// Gracefully shutdown
	<-sigc
	writeFileEntries(DbPath)
	fmt.Println("SIGINT: cleaning up and shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := ech.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(
	w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func downloadFile(c echo.Context) error {
	ids := c.Param("fileid")
	id := dieOnErr2(strconv.Atoi(ids))
	entry := GetEntryById(id)
	return c.File(entry.GetFilepath())
}

func uploadFile(c echo.Context) error {
	return c.String(http.StatusNotImplemented, "Not implemented")
}

func deleteFile(c echo.Context) error {
	return c.String(http.StatusNotImplemented, "Not implemented")
}

func allViews(c echo.Context) error {
	return c.String(http.StatusNotImplemented, "Not implemented")
}

func homeView(c echo.Context) error {
	return c.Render(http.StatusOK, "filesview.html", GetFileEntries())
}

func getView(c echo.Context) error {
	return c.String(http.StatusNotImplemented, "Not implemented")
}

func createView(c echo.Context) error {
	return c.String(http.StatusNotImplemented, "Not implemented")
}

func deleteView(c echo.Context) error {
	return c.String(http.StatusNotImplemented, "Not implemented")
}
