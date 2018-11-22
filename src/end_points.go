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
	storeClient.KVJSON.Delete("IplayerRecommend", "all")
	userAuthenticated = false
	if isRunning {
		close(StopChan)
	}
	http.Redirect(w, r, Host+"/ui", 302)
}

type Recommendations struct {
	Version             string `json:"version"`
	Schema              string `json:"schema"`
	UserRecommendations struct {
		Elements []struct {
			Type      string `json:"type"`
			Algorithm string `json:"algorithm"`
			Episode   struct {
				ID     string `json:"id"`
				Live   bool   `json:"live"`
				Type   string `json:"type"`
				Title  string `json:"title"`
				Images struct {
					Type     string `json:"type"`
					Standard string `json:"standard"`
				} `json:"images"`
			} `json:"episode"`
		} `json:"elements"`
	} `json:"user_recommendations"`
}

func info(w http.ResponseWriter, r *http.Request) {

	recommendations, err := storeClient.KVJSON.Read("IplayerRecommend", "all")
	libDatabox.ChkErr(err)

	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, recommendations, "", "    ")

	var recObj Recommendations
	err = json.Unmarshal(recommendations, &recObj)
	libDatabox.ChkErr(err)

	images := ""
	for _, el := range recObj.UserRecommendations.Elements {
		img := strings.Replace(el.Episode.Images.Standard, "{recipe}", "432x243", 1)
		images += `<a target="_blank" href="https://www.bbc.co.uk/iplayer/episode/` + el.Episode.ID + `"><img class="img-episode" src="` + img + `" /></a>`
	}

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
	<h1>iPlayer Recommendations</h1>
	<div style="float:right"><a href="` + BasePath + `/ui/logout">logout</a></div>
	<pre style="clear: both;">%s</pre>
	</body>
	</html>`

	fmt.Fprintf(w, body, images)

}

func index(w http.ResponseWriter, r *http.Request) {

	if userAuthenticated {
		http.Redirect(w, r, Host+"/ui/info", 302)
		return
	}

	callbackUrl := r.FormValue("post_auth_callback")
	PostAuthCallbackUrl := "/core-ui/ui/view/" + BasePath + "/info"
	if DataboxTestMode {
		PostAuthCallbackUrl = "/ui/info"
	}
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
