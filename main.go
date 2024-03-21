package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"strings"

	"os"
	"runtime"
	"strconv"

	"templrjs/pkg/config"
	"templrjs/pkg/duckdb"
	"templrjs/pkg/router"

	"github.com/gin-contrib/cors" // Add this line to your existing imports
	"github.com/gin-gonic/gin"

	"github.com/natefinch/lumberjack"
	log "github.com/sirupsen/logrus"
)

var (
	stage             string
	flagPort          int
	flagSource        string
	flagEnv           string
	flagStaticDirPath string
	flagIndexFileName string
)

func init() {
	log.Info("Application init function start")

	log.Info("Initialize Logger")
	initLogger()
	log.Info("Initialized Logger")

	log.Info("Initialize Configuration")
	config.Setup()
	log.Info("Initialized Configuration")

	log.Info("Initialize command line args")
	flag.IntVar(&flagPort, "p", 8080, "port number for the api server")
	flag.StringVar(&flagStaticDirPath, "d", "./dist", "Website directory path")
	flag.StringVar(&flagIndexFileName, "f", "index.html", "Index document")
	flag.StringVar(&flagEnv, "e", "dev", "Development or Production stage")
	flag.StringVar(&flagSource, "s", "embed", "Host site from embedded source or local filesystem")
	log.Info("Initialized command line args")

	//database.Setup()
	log.Info("Initialize DuckDB")
	duckdb.Setup()
	log.Info("Initialized DuckDB")
	//broker.Setup()

	log.Info("Application init function end.")
}

//go:embed all:dist
var content embed.FS

func main() {

	flag.Parse()

	startServer()

}

func initLogger() {
	//logfilepath := AppExecutionPath() + "/logs/" + "templrjs.log"
	logfilepath := AppExecutionPath() + "/" + os.Args[0] + ".log"
	log.Info("logfilepath = " + logfilepath)
	// Set the Lumberjack logger
	ljack := &lumberjack.Logger{
		Filename:   logfilepath,
		MaxSize:    1,
		MaxBackups: 3,
		MaxAge:     3,
		LocalTime:  true,
	}

	//log := logrus.New()
	//
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

	mWriter := io.MultiWriter(os.Stdout, ljack)
	log.SetOutput(mWriter)
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})
	log.WithFields(log.Fields{
		"app":             os.Args[0],
		"Runtime Version": runtime.Version(),
		"Number of CPUs":  runtime.NumCPU(),
		"Arch":            runtime.GOARCH,
	}).Info("Application Initializing")

	/*if stage == "Dev" {
		log.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: true, TimestampFormat: time.RFC1123Z})
	} else {
		//log.SetFormatter(&log.JSONFormatter{})
		log.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: true, TimestampFormat: time.RFC1123Z})
	}*/
}

func Assets() (fs.FS, error) {
	return fs.Sub(content, "dist")
}

func startServer() {

	r := router.Setup()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*", "http://localhost:3000"} // Allow all origins
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	r.Use(cors.New(config))

	if strings.ToLower(flagSource) == "embed" {
		assets, _ := Assets()
		// Use the WrapH function to wrap the http.StripPrefix handler
		r.Use(func(c *gin.Context) {
			fs := http.StripPrefix("/", http.FileServer(http.FS(assets)))
			fs.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		})
	} else {
		log.Info("Serving from local filesystem")
		r.GET("/", func(c *gin.Context) { c.File("./dist/index.html") })
		r.GET("/resume", func(c *gin.Context) { c.File("./dist/resume/index.html") })
		r.GET("/cms", func(c *gin.Context) { c.File("./dist/cms/index.html") })
		r.GET("/aboutme", func(c *gin.Context) { c.File("./dist/cms/index.html") })
		r.GET("/sql-ide", func(c *gin.Context) { c.File("./dist/sql-ide/index.html") })
		r.Static("/_nuxt", "./dist/_nuxt")
		r.StaticFile("/favicon.ico", "./dist/favicon.ico")
		r.StaticFile("/logo.svg", "./dist/logo.svg")
		r.Static("/configs", "./dist/configs")
	}

	//log.Fatal(http.ListenAndServe(":"+port, a.Negroni))
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(flagPort))
	if err != nil {
		log.Error(err)
		log.Fatal(err)
	}
	log.Info("Starting the server " + strconv.Itoa(flagPort))
	done := make(chan bool)
	go http.Serve(listener, r)
	log.Info("Server started at port " + strconv.Itoa(flagPort))
	<-done
}

// AppExecutionPath returns the relative path where the application is executing
func AppExecutionPath() string {
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("App path = " + mydir)
	return mydir
}
