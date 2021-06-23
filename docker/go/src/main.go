package main
// DB_HOST=localhost DB_USER=root DB_NAME=cdb_api_dev APP_URL=go.k99s.com go run main.go
import (
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
		"os"
		"io"
		"net/http"
		"log"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<h1>Welcome to %s</h1><br/> <p>Here is the list of tables in your database:</p> <br/>", os.Getenv("APP_URL"))
		writeData(w)
	})
	log.Print("Serving requests on port 80")
	log.Fatal(http.ListenAndServe(":80", nil))

}

func writeData(w io.Writer){
	connection := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
									os.Getenv("DB_HOST"), 3306, os.Getenv("DB_NAME"))
    db, err := sql.Open("mysql", connection)
    if err != nil {
			log.Print(err.Error())
			return
    }
    defer db.Close()
    results, err := db.Query("show tables;")
    if err != nil {
        log.Print(err.Error())
				return
    }

		for results.Next() {
			var row string
			err = results.Scan(&row)
			if err != nil {
        log.Print(err.Error())
				return
			}
			fmt.Fprintf(w, "%s<br/>", row)
	}
}