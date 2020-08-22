package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"html/template"
	"mime"
	"net/http"
	"os"
	"strings"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("src/views/*.html"))
}

func Index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index", nil)
}

func Connstr() (db *sql.DB) {
	dbDriver := "mssql"
	dbUser := "DESKTOP-1VA2HU8\\srave"
	//dbPass := "your_password"
	dbName := "portfolio"
	db, err := sql.Open(dbDriver, dbUser+":@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func DIndex(w http.ResponseWriter, r *http.Request) {
	db := Connstr()
	rows, err := db.Query("SELECT  * FROM fileupdown  ")
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		var filebytes []byte
		err = rows.Scan(&filebytes)
		if err != nil {
			panic(err.Error())
		}
		detectedFileType := http.DetectContentType(filebytes)
		switch detectedFileType {
		case "image/jpeg", "image/jpg":
		case "image/gif", "image/png":
		case "application/pdf":
			break
		default:
			renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
			return
		}
		fileEndings, err := mime.ExtensionsByType(detectedFileType)
		str3 := strings.Join(fileEndings, ", ")
		fileName := string(randToken(12)) + str3
		if err != nil {
			renderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
			return
		}
		contentType := http.DetectContentType(filebytes)
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Type", contentType)
		w.Write(filebytes)
	}
}
func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", Index)
	mux.HandleFunc("/download", DIndex)
	http.ListenAndServe(":"+port, mux)
}
