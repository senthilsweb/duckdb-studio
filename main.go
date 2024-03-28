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

	"templrjs/pkg/router"

	"github.com/gin-contrib/cors"
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

	log.Info("Initialize command line args")
	flag.IntVar(&flagPort, "p", 8080, "port number for the api server")
	flag.StringVar(&flagStaticDirPath, "d", "./.dist", "Website directory path") // Updated to reflect .dist directory
	flag.StringVar(&flagIndexFileName, "f", "index.html", "Index document")
	flag.StringVar(&flagEnv, "e", "dev", "Development or Production stage")
	flag.StringVar(&flagSource, "s", "embed", "Host site from embedded source or local filesystem")
	log.Info("Initialized command line args")

	log.Info("Application init function end.")
}

//go:embed all:.dist
var content embed.FS

func main() {

	flag.Parse()

	startServer()

}

func initLogger() {
	logfilepath := AppExecutionPath() + "/" + os.Args[0] + ".log"
	log.Info("logfilepath = " + logfilepath)

	ljack := &lumberjack.Logger{
		Filename:   logfilepath,
		MaxSize:    1,
		MaxBackups: 3,
		MaxAge:     3,
		LocalTime:  true,
	}

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

}

func Assets() (fs.FS, error) {
	return fs.Sub(content, ".dist") // Updated to match .dist directory
}

func startServer() {

	r := router.Setup()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*", "http://localhost:3000"}
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	r.Use(cors.New(config))

	if strings.ToLower(flagSource) == "embed" {
		assets, _ := Assets()

		r.Use(func(c *gin.Context) {
			fs := http.StripPrefix("/", http.FileServer(http.FS(assets)))
			fs.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		})
	} else {
		log.Info("Serving from local filesystem")
		r.GET("/", func(c *gin.Context) { c.File("./.dist/index.html") })                // Updated paths to match .dist directory
		r.GET("/sql-ide", func(c *gin.Context) { c.File("./.dist/sql-ide/index.html") }) // Updated paths
		r.Static("/_nuxt", "./.dist/_nuxt")                                              // Updated paths
		r.StaticFile("/favicon.ico", "./.dist/favicon.ico")                              // Updated paths
		r.StaticFile("/logo.svg", "./.dist/logo.svg")                                    // Updated paths
	}

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

func AppExecutionPath() string {
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("App path = " + mydir)
	return mydir
}
