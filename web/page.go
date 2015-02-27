package main

import (
  "fmt"
  "html/template"
  //"io"
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
  web_port  int
  )

  const templResult = `{{define "T"}}<head>
    <meta charset="utf-8">
    <title>Swarm search engine</title>
    <meta name="description" content="description here">
    <meta name="author" content="content here">
    <link rel="stylesheet" href="static/css/5css.css">
  <style type="text/css">
  </style>
  </head>
  <body>
  <div id="header-container">
      <header class="wrapper">
          <nav>
           <div class="pul"><a href="http://localhost:8888">Home</a></div>
          </nav>
      </header>
  </div>
  <h1> Result of search </h1>
  <div id="main" class="wrapper">
  <table align="center" style="width:400px">
    {{range .}}
      <tr> <p> {{.}}</p></tr>
      {{end}}
  </table>
  </div>
    <div id="footer-container">
      <footer class="wrapper">
      </footer>
    </div>
  </body>
  </html>{{end}}`


func init(){
  server_port = resources.SERVER_PORT
  if os.Getenv(resources.PORT_FLAG) != "" {
      server_port, _ = strconv.Atoi(os.Getenv(resources.PORT_FLAG))
  }
  server_addr = resources.SERVER_ADDR
  if os.Getenv(resources.SERVER_ADDR_FLAG) != "" {
    server_addr = os.Getenv(resources.SERVER_ADDR_FLAG)
  }
  web_port = resources.WEB_SERVER_PORT
  if os.Getenv(resources.WEB_PORT_FLAG) != "" {
    web_port, _ = strconv.Atoi(os.Getenv(resources.WEB_PORT_FLAG))
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
    // Avoid to crash the server
    fmt.Fprintln(w, "ERROR: There has occurred an error when searching :( ")
  }

  t, err := template.New("result").Parse(templResult)
  text_to_render := arrayOfTuples(result)
  if len(text_to_render) == 0 {
    fmt.Sprintf("=> No results found with this criteria.")
  }

  err = t.ExecuteTemplate(w, "T", text_to_render)
}

func arrayOfTuples(result []map[string]interface{}) ([]string) {
  list := make([]string, 0)
  for _, user := range result {
    list = append(list, fmt.Sprintf("ID: %d USERNAME: %s ROLE: %s ORGANIZATION: %s",
                                user["id"].(int64),
                                user["username"].(string),
                                user["role"].(string),
                                user["organization"].(string)))
  }
  return list

}

func main() {
  http.HandleFunc("/", searchHandler)
  http.HandleFunc("/search", listHandler)

  fs := http.FileServer(http.Dir("static"))
  http.Handle("/static/", http.StripPrefix("/static/", fs))

  logger.Infof("Started web server for swarm search-engine at %d", web_port)

  err := http.ListenAndServe(fmt.Sprintf(":%d", web_port), nil)
  if err != nil {
    logger.Errorf("Error starting the web server at port %d", web_port)
    panic("Tracestack: " + err.Error())
  }
}
