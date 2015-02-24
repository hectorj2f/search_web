package main

import (
  "fmt"
  "html/template"
  "net/http"
  "os"
  "strconv"

  "github.com/hectorj2f/search_networking/networking"
  "github.com/hectorj2f/search_web/resources"

  logger "github.com/Sirupsen/logrus"
)

var (
  server_addr string
  server_port int
  )

func init(){
  server_port = resources.SERVER_PORT
  if os.Getenv(resources.PORT_FLAG) != "" {
      server_port, _ = strconv.Atoi(os.Getenv(resources.PORT_FLAG))
  }
  server_addr = resources.SERVER_ADDR
  if os.Getenv(resources.SERVER_ADDR_FLAG) != "" {
    server_addr = os.Getenv(resources.SERVER_ADDR_FLAG)
  }
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("index.html")
    t.Execute(w, nil)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
  search_query := make(map[string]interface{})
  if r.FormValue("organization") != "" {
    search_query["organization"] = r.FormValue("organization")
  }
  if r.FormValue("role") != "" {
    search_query["role"] = r.FormValue("role")
  }
  if r.FormValue("id") != "" {
    search_query["id"] = r.FormValue("id")
  }
  if r.FormValue("username") != "" {
    search_query["username"] = r.FormValue("username")
  }

  result, err := networking.Query(search_query, server_addr, true, server_port)
  if err != nil {
    logger.Errorf("Error while searching users %s", err)
    fmt.Fprintln(w, "ERROR: There has occurred an error when searching :( ")
    return
  }

  for _, user := range result {
    fmt.Fprintf(w,"ID: %d CREATED: %s USERNAME: %s ROLE: %s ORGANIZATION: %s\n",
                 user["id"].(int64),
                 user["created"].(string),
                 user["username"].(string),
                 user["role"].(string),
                 user["organization"].(string))
  }
  if len(result) == 0 {
    fmt.Fprintln(w,"=> No results found with this criteria.")
  }
}

func main() {
  http.HandleFunc("/", searchHandler)
  http.HandleFunc("/search", listHandler)
  http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, r.URL.Path[1:])
  })

  logger.Infof("Started web server for swarm search engine at %s", resources.WEB_SERVER_PORT)

  err := http.ListenAndServe(fmt.Sprintf(":%s",resources.WEB_SERVER_PORT), nil)
  if err != nil {
    logger.Errorf("Error starting the web server at port %s", resources.WEB_SERVER_PORT)
    panic("Tracestack: " + err.Error())
  }
}
