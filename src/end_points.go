package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	libDatabox "github.com/me-box/lib-go-databox"
)

func authUser(w http.ResponseWriter, r *http.Request) {
	libDatabox.Info("Obtained auth")

	r.ParseForm()
	//Obtain user login details for their BCC account
	for k, v := range r.Form {
		err := storeClient.KVText.Write("IplayerCred", k, []byte(strings.Join(v, "")))
		libDatabox.ChkErr(err)
	}

	token := authCheck()

	if token != "" {
		http.Redirect(w, r, RedirectHostInsideDatabox+"/ui/info", 302)
		go driverWork(token, stopChan)
		userAuthenticated = true
	} else {
		fmt.Fprintf(w, "<h1>Auth Failed<h1>")
		fmt.Fprintf(w, " <button onclick='goBack()'>Go Back</button><script>function goBack() {	window.history.back();}</script> ")
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	storeClient.KVText.Delete("IplayerCred", "email")
	storeClient.KVText.Delete("IplayerCred", "password")
	userAuthenticated = false
	stopChan <- 1
	http.Redirect(w, r, RedirectHostInsideDatabox+"/ui", 302)
}

func info(w http.ResponseWriter, r *http.Request) {

	recommendations, err := storeClient.TSBlobJSON.Latest("IplayerRecommend")
	libDatabox.ChkErr(err)

	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, recommendations, "", "    ")

	body := `<!doctype html>
	<html class="no-js" lang="">

	<head>
	<meta charset="utf-8">
	<meta http-equiv="x-ua-compatible" content="ie=edge">
	<title></title>
	<meta name="description" content="">
	<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

	<link rel="stylesheet" href="./ui/css/normalize.css">
	<link rel="stylesheet" href="./ui/css/main.css">
	</head>

	<body>
	<h1>iPlayer driver loading recommendations!</h1>
	<div style="float:right"><a href="/driver-bbc-iplayer/ui/logout">logout</a></div>
	<pre style="clear: both;">%s</pre>
	</body>
	</html>`

	fmt.Fprintf(w, body, string(prettyJSON.Bytes()))

}

func index(w http.ResponseWriter, r *http.Request) {

	if userAuthenticated {
		http.Redirect(w, r, RedirectHostInsideDatabox+"/ui/info", 302)
		return
	}

	body := `<!doctype html>
	<head>
	  <meta charset="utf-8">
	  <meta http-equiv="x-ua-compatible" content="ie=edge">
	  <title></title>
	  <meta name="description" content="">
	  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

	  <link rel="stylesheet" href="./ui/css/normalize.css">
	  <link rel="stylesheet" href="./ui/css/main.css">
	</head>

	<body>
	  <h1>Authentication Form</h1>
		<form action="./ui/auth" method="post">
		  Username:<input type="text" name="email" required><br>
			Password: <input type="password" name="password" required><br>
				<input type="submit" value="Send">
	  </form>
	</body>
	</html>`

	fmt.Fprintf(w, body)

}
