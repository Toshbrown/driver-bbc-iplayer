package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	libDatabox "github.com/me-box/lib-go-databox"
)

//default addresses to be used in testing mode
const testArbiterEndpoint = "tcp://127.0.0.1:4444"
const testStoreEndpoint = "tcp://127.0.0.1:5555"

var (
	storeClient *libDatabox.CoreStoreClient
)

func main() {
	DataboxTestMode := os.Getenv("DATABOX_VERSION") == ""
	registerData(DataboxTestMode)

	router := mux.NewRouter()
	router.HandleFunc("/status", statusEndpoint).Methods("GET")
	router.HandleFunc("/ui/saved", infoUser)
	router.HandleFunc("/ui/info", infoUser)
	router.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(http.Dir("./static"))))
	setUpWebServer(DataboxTestMode, router, "8080")
}

func statusEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("active\n"))
}

func registerData(testMode bool) {
	//Setup store client
	var DataboxStoreEndpoint string
	if testMode {
		DataboxStoreEndpoint = testStoreEndpoint
		ac, _ := libDatabox.NewArbiterClient("./", "./", testArbiterEndpoint)
		storeClient = libDatabox.NewCoreStoreClient(ac, "./", DataboxStoreEndpoint, false)
		//turn on debug output for the databox library
		libDatabox.OutputDebug(true)
	} else {
		DataboxStoreEndpoint = os.Getenv("DATABOX_ZMQ_ENDPOINT")
		storeClient = libDatabox.NewDefaultCoreStoreClient(DataboxStoreEndpoint)
	}
	//Setup datastore for main data
	recomDatasource := libDatabox.DataSourceMetadata{
		Description:    "IPlayer Recommendation data", //required
		ContentType:    libDatabox.ContentTypeJSON,    //required
		Vendor:         "databox-test",                //required
		DataSourceType: "recommendData",               //required
		DataSourceID:   "IplayerRecommend",            //required
		StoreType:      libDatabox.StoreTypeTSBlob,    //required
		IsActuator:     false,
		IsFunc:         false,
	}
	dErr := storeClient.RegisterDatasource(recomDatasource)
	if dErr != nil {
		libDatabox.Err("Error Registering Datasource " + dErr.Error())
		return
	}
	libDatabox.Info("Registered Datasource")
	//Setup authentication datastore
	authDatasource := libDatabox.DataSourceMetadata{
		Description:    "IPlayer Login Data",       //required
		ContentType:    libDatabox.ContentTypeTEXT, //required
		Vendor:         "databox-test",             //required
		DataSourceType: "loginData",                //required
		DataSourceID:   "IplayerCred",              //required
		StoreType:      libDatabox.StoreTypeKV,     //required
		IsActuator:     false,
		IsFunc:         false,
	}
	cErr := storeClient.RegisterDatasource(authDatasource)
	if cErr != nil {
		libDatabox.Err("Error Registering Credential Datasource " + cErr.Error())
		return
	}
	libDatabox.Info("Registered Credential Datasource")
}

func setUpWebServer(testMode bool, r *mux.Router, port string) {

	//Start up a well behaved HTTP/S server for displying the UI

	srv := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
		Handler:      r,
	}
	if testMode {
		//set up an http server for testing
		libDatabox.Info("Waiting for http requests on port http://127.0.0.1" + srv.Addr)
		log.Fatal(srv.ListenAndServe())
	} else {
		//configure tls
		tlsConfig := &tls.Config{
			PreferServerCipherSuites: true,
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
			},
		}
		srv.TLSConfig = tlsConfig

		libDatabox.Info("Waiting for https requests on port " + srv.Addr)
		log.Fatal(srv.ListenAndServeTLS(libDatabox.GetHttpsCredentials(), libDatabox.GetHttpsCredentials()))
	}
}

func infoUser(w http.ResponseWriter, r *http.Request) {
	libDatabox.Info("Obtained auth")

	r.ParseForm()
	//Obtain user login details for their BCC account
	for k, v := range r.Form {
		if k == "email" {
			err := storeClient.KVText.Write("IplayerCred", "username", []byte(strings.Join(v, "")))
			if err != nil {
				libDatabox.Err("Error Write Datasource " + err.Error())
			}

		} else {
			err := storeClient.KVText.Write("IplayerCred", "password", []byte(strings.Join(v, "")))
			if err != nil {
				libDatabox.Err("Error Write Datasource " + err.Error())
			}
		}

	}

	token := authCheck()

	if token != "" {
		fmt.Fprintf(w, "<h1>Auth success<h1>")
		go driverWork(token)
	} else {
		fmt.Fprintf(w, "<h1>Auth Failed<h1>")
		fmt.Fprintf(w, " <button onclick='goBack()'>Go Back</button><script>function goBack() {	window.history.back();}</script> ")
	}
}

func infoSaved(w http.ResponseWriter, r *http.Request) {
	//Check to see if any password is saved inside the auth datastore
	tempPas, pErr := storeClient.KVText.Read("YoutubeHistoryCred", "password")
	if pErr != nil {
		fmt.Println(pErr.Error())
		return
	}
	//If there is no saved password, warn the user, otherwise run the driver
	if tempPas != nil {
		libDatabox.Info("Saved auth detected")
		fmt.Fprintf(w, "<h1>Saved authentication detected<h1>")
		token := authCheck()

		if token != "" {
			fmt.Fprintf(w, "<h1>Auth success<h1>")
			go driverWork(token)
		} else {
			fmt.Fprintf(w, "<h1>Auth Failed<h1>")
			fmt.Fprintf(w, " <button onclick='goBack()'>Go Back</button><script>function goBack() {	window.history.back();}</script> ")
		}

	} else {
		fmt.Fprintf(w, "<h1>No saved auth detected<h1>")
		fmt.Fprintf(w, " <button onclick='goBack()'>Go Back</button><script>function goBack() {	window.history.back();}</script> ")
	}
}

func authCheck() (token string) {
	tempUse, uErr := storeClient.KVText.Read("IplayerCred", "username")
	if uErr != nil {
		fmt.Println(uErr.Error())
		return
	}

	tempPas, pErr := storeClient.KVText.Read("IplayerCred", "password")

	if pErr != nil {
		fmt.Println(pErr.Error())
		return
	}
	token, err := Auth(string(tempUse), string(tempPas))
	if err != nil {
		fmt.Println("Error ", err)
		return ""
	}
	return token

}

func driverWork(token string) {
	for {
		recommendations, err := GetRecommendations(token)
		if err != nil {
			libDatabox.Err("Error Write Datasource " + err.Error())
		}

		aerr := storeClient.TSBlobJSON.Write("IplayerRecommend", []byte(recommendations))
		if aerr != nil {
			libDatabox.Err("Error Write Datasource " + aerr.Error())
		}
		//libDatabox.Info("Data written to store: " + recommendations)
		libDatabox.Info("Storing data")

		//time.Sleep(time.Hour * 24)
		time.Sleep(time.Second * 30)
	}
}
