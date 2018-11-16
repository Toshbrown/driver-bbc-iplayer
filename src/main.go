package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	libDatabox "github.com/me-box/lib-go-databox"
)

//default addresses to be used in testing mode

const (
	testArbiterEndpoint    = "tcp://127.0.0.1:4444"
	testStoreEndpoint      = "tcp://127.0.0.1:5555"
	HostInsideDatabox      = "https://driver-bbc-iplayer:8080"
	HostOutsideDatabox     = "http://127.0.0.1:8080"
	BasePathInsideDatabox  = "/driver-bbc-iplayer"
	BasePathOutsideDatabox = ""
)

var (
	storeClient       *libDatabox.CoreStoreClient
	userAuthenticated bool
	stopChan          chan int
	Host              string
	BasePath          string
)

func main() {
	DataboxTestMode := os.Getenv("DATABOX_VERSION") == ""
	userAuthenticated = false
	stopChan = make(chan int)

	//Setup store client
	var DataboxStoreEndpoint string
	if DataboxTestMode {
		Host = HostOutsideDatabox
		BasePath = BasePathOutsideDatabox
		DataboxStoreEndpoint = testStoreEndpoint
		ac, _ := libDatabox.NewArbiterClient("./", "./", testArbiterEndpoint)
		storeClient = libDatabox.NewCoreStoreClient(ac, "./", DataboxStoreEndpoint, false)
		//turn on debug output for the databox library
		libDatabox.OutputDebug(true)
	} else {
		Host = HostInsideDatabox
		BasePath = BasePathInsideDatabox
		DataboxStoreEndpoint = os.Getenv("DATABOX_ZMQ_ENDPOINT")
		storeClient = libDatabox.NewDefaultCoreStoreClient(DataboxStoreEndpoint)
	}

	registerData()

	go func() {
		time.Sleep(time.Second * 5)
		token := authCheck()
		if token != "" {
			userAuthenticated = true
			libDatabox.Info("Email and password retrieved form DB starting do driver work")
			go driverWork(token, stopChan)
		}
	}()

	router := mux.NewRouter()
	router.HandleFunc("/status", statusEndpoint).Methods("GET")
	router.HandleFunc("/ui/auth", authUser)
	router.HandleFunc("/ui/logout", logout)
	router.HandleFunc("/ui/info", info)
	router.HandleFunc("/ui", index)
	router.PathPrefix("/ui/*").Handler(http.StripPrefix("/ui", http.FileServer(http.Dir("./static"))))
	setUpWebServer(DataboxTestMode, router, "8080")
}

func statusEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("active\n"))
}

func registerData() {

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

func authCheck() (token string) {
	tempUse, uErr := storeClient.KVText.Read("IplayerCred", "email")
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

func driverWork(token string, stop chan int) {
	for {
		recommendations, err := GetRecommendations(token)
		if err != nil {
			libDatabox.Err("Error Write Datasource " + err.Error())
		}

		aerr := storeClient.TSBlobJSON.Write("IplayerRecommend", []byte(recommendations))
		if aerr != nil {
			libDatabox.Err("Error Write Datasource " + aerr.Error())
		}

		libDatabox.Info("Storing data")

		select {
		case <-stop:
			libDatabox.Info("Stopping data updates stop message received")
			return
		case <-time.After(time.Hour * 1):
			libDatabox.Info("updating data after time out")
		}
	}
}
