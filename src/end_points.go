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
		callbackUrl := r.FormValue("post_auth_callback")
		PostAuthCallbackUrl := "/core-ui/ui/view/" + BasePath + "/info"
		if callbackUrl != "" {
			PostAuthCallbackUrl = callbackUrl
		}

		if DataboxTestMode {
			fmt.Fprintf(w, "<html><head><script>window.location = '%s';</script><head><body><body></html>", PostAuthCallbackUrl)
		} else {
			fmt.Fprintf(w, "<html><head><script>window.parent.location = '%s';</script><head><body><body></html>", PostAuthCallbackUrl)
		}

		if !isRunning {
			StopChan = make(chan struct{})
			go driverWork(token, StopChan)
		}
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
	if isRunning {
		close(StopChan)
	}
	http.Redirect(w, r, Host+"/ui", 302)
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

	<link rel="stylesheet" href="` + BasePath + `/ui/css/normalize.css">
	<link rel="stylesheet" href="` + BasePath + `/ui/css/main.css">
	<link rel="stylesheet" href="` + BasePath + `/ui/css/bbc.css">
	</head>

	<body>
	<h1>iPlayer driver loading recommendations!</h1>
	<div style="float:right"><a href="` + BasePath + `/ui/logout">logout</a></div>
	<pre style="clear: both;">%s</pre>
	</body>
	</html>`

	fmt.Fprintf(w, body, string(prettyJSON.Bytes()))

}

func index(w http.ResponseWriter, r *http.Request) {

	if userAuthenticated {
		http.Redirect(w, r, Host+"/ui/info", 302)
		return
	}

	callbackUrl := r.FormValue("post_auth_callback")
	PostAuthCallbackUrl := "/core-ui/ui/view/" + BasePath + "/info"
	if callbackUrl != "" {
		PostAuthCallbackUrl = callbackUrl
	}

	body := `<!doctype html>
	<head>
	  <meta charset="utf-8">
	  <meta http-equiv="x-ua-compatible" content="ie=edge">
	  <title></title>
	  <meta name="description" content="">
	  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

	  <link rel="stylesheet" href="` + BasePath + `/ui/css/normalize.css">
	  <link rel="stylesheet" href="` + BasePath + `/ui/css/main.css">
	  <link rel="stylesheet" href="` + BasePath + `/ui/css/bbc.css">
	</head>

	<body>
	  	<div class="form-login">
        <img class="logo" src="` + BasePath + `/ui/img/BBC_iPlayer_logo.svg" />
        <p>Sign in with your BBC account to download your iPlayer recommendations.</p>
			<form action="` + BasePath + `/ui/auth" method="post">
				<div class="row"> <label for="email">Email </label><input autocomplete="off" type="text" name="email" required></div>
				<div class="row"> <label for="password">Password </label><input autocomplete="off" type="password" name="password" required></div>
				<div class="row"> <input type="submit" class="btn-login" value="Sign in"></div>
				<input style="display: none" type="text" name="post_auth_callback" value="` + PostAuthCallbackUrl + `"/>
			</form>
		</div>
	</body>
	</html>`

	fmt.Fprintf(w, body)

}
