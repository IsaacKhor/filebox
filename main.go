package main

import (
	"context"
	"crypto/subtle"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	loadConfig()
	fileEntries = loadFileEntries(Config.DbPath)
	// Create directory if it doesn't exist
	err := os.MkdirAll(Config.FilesPath, 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	ech := echo.New()
	ech.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(10)))
	ech.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "REQ ${time_rfc3339}: ${remote_ip} ${method} ${uri}. " +
			"I/O ${bytes_in}/${bytes_out}. " +
			"${status} in ${latency_human}.\n",
	}))
	ech.Use(middleware.BasicAuth(checkToken))
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
	go startServer(ech, Config)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)

	// Gracefully shutdown
	<-sigc
	writeFileEntries(Config.DbPath)
	log.Println("SIGINT: cleaning up and shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := ech.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}

func startServer(ech *echo.Echo, config FileboxConfig) {
	// Test environment, run at port 8080 over http
	if !config.Production {
		if err := ech.Start(":8080"); err != nil {
			ech.Logger.Error(err)
		}
	}

	addr := fmt.Sprintf(":%d", config.ProductionPort)
	err := ech.StartTLS(addr, config.CertificatePath, config.PrivKeyPath)
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func checkToken(username string, pwd string, ctx echo.Context) (bool, error) {
	// Password is ignored, use username is auth token
	if subtle.ConstantTimeCompare(
		[]byte(username), []byte(Config.AuthToken)) == 1 {
		return true, nil
	}
	return false, nil
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
	id, err := strToI64(ids)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid fileid")
	}

	if !HasFileEntry(id) {
		return c.String(http.StatusNotFound, "Fileid not found")
	}

	entry := GetEntryById(id)
	return c.File(entry.GetFilepath())
}

func uploadFile(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid upload POST request")
	}
	files := form.File["files"]
	log.Println(form, files)

	for _, file := range files {
		log.Println("Uploading: ", file)
		saveFile(file)
	}

	return c.Render(http.StatusCreated, "postUpload.html", nil)
}

func saveFile(file *multipart.FileHeader) {
	entry := CreateFileEntry(file.Filename, file.Size)

	src := panicOnErr(file.Open())
	defer closeOrPanic(src)

	dst := panicOnErr(os.Create(entry.GetFilepath()))
	defer closeOrPanic(dst)

	_ = panicOnErr(io.Copy(dst, src))
}

func deleteFile(c echo.Context) error {
	ids := c.Param("fileid")
	id, err := strToI64(ids)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid fileid")
	}
	if !HasFileEntry(id) {
		return c.String(http.StatusNotFound, "Fileid not found")
	}

	entry := GetEntryById(id)
	checkErr(os.Remove(entry.GetFilepath()))
	RemoveFileEntry(id)
	return c.String(http.StatusAccepted, "File deleted")
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
