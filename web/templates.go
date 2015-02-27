package main


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
         <div class="pul"><a href="http://localhost:8080">Home</a></div>
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
