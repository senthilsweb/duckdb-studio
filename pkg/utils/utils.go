package utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aymerick/raymond"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// EncodeAuthToken creates authentication token
func EncodeAuthToken(uid uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["userID"] = uid
	claims["IssuedAt"] = time.Now().Unix()
	claims["ExpiresAt"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	return token.SignedString([]byte(os.Getenv("SECRET")))
}

func GetValElseSetEnvFallback(input string, key string) string {
	value := gjson.Get(input, key)
	if len(value.String()) == 0 {
		return os.Getenv(key)
	}
	return value.String()
}

func GetFileNameWithoutExt(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

func GetFileExt(fileName string) string {
	return filepath.Ext(fileName)
}

// contains
func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

// RenderFromTemplate returns  content
func RenderFromTemplate(filename string, data interface{}) string {
	var err error
	fp := AppExecutionPath() + "/data/templates/" + filename
	fmt.Println("\nTemplate Path " + fp)
	//tpl := raymond.MustParse("{{> myPartial name=hero }}")
	tpl, err := raymond.ParseFile(fp)
	result, err := tpl.Exec(data)
	if err != nil {
		log.Error(err)
		panic("Please report a bug :)")
	}
	log.Info(result)
	return result
}

// GinContextFromContext  from context
func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value("GinContextKey")
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}

func GetStringFromGinRequestBody(c *gin.Context) string {
	jsonData, err := c.GetRawData()
	res := bytes.NewBuffer(jsonData).String()
	if err != nil {
		return ""
	}
	return res
}

// AppExecutionPath returns the relative path where the application is executing
func AppExecutionPath() string {
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mydir)
	return mydir
}

func DecodeAsStringFromBase64EncodedString(msg string) (output string) {
	decodedstg, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		log.Fatal("Error during decode")
	}
	body := []byte(decodedstg)
	var strToConvert string
	strToConvert = bytes.NewBuffer(body).String()
	return strToConvert
}
func GetContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

func ValInSlice(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func KeyInSlice(slice []string, val string) int {
	for i, item := range slice {
		if item == val {
			return i
		}
	}
	return -1
}

func KeyInMap(data map[string]interface{}, key string) bool {
	for item, _ := range data {
		if item == key {
			return true
		}
	}
	return false
}

func KeyInStringMap(data map[string]string, key string) bool {
	for item, _ := range data {
		if item == key {
			return true
		}
	}
	return false
}

func GetFormateTime() string {
	t := time.Now()
	day := strconv.Itoa(t.Day())
	mon := string(int(t.Month()))
	hour := strconv.Itoa(t.Hour())
	min := strconv.Itoa(t.Minute())
	//week := t.Weekday().String()[0:3]

	return mon + "-" + day + "-" + hour + ":" + min
}

func FmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}

func FormatedDuration(dur time.Duration) string {
	m := math.Round(dur.Minutes())
	h := math.Floor(m / 60)
	if h == 0 {
		return fmt.Sprintf("%02dm", int(m))
	}
	d := int(math.Floor(h / 24))
	if d == 0 {
		return fmt.Sprintf("%02dh", int(h))
	}
	return fmt.Sprintf("%02dd", int(d))
}
