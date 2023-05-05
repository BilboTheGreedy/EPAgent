package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Configuration struct {
	APIKey    string `json:"apiKey"`
	CertFile  string `json:"certFile"`
	KeyFile   string `json:"keyFile"`
	AgentPort int    `json:"agentPort"`
}

type PwshRequestBody struct {
	Path string
}

var (
	cfg Configuration
)

// Load config.json
func LoadConfiguration(file string) Configuration {
	var config Configuration
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func checkAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Request.Header.Get("X-API-Key")
		if apiKey != cfg.APIKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			return
		}
		c.Next()
	}
}

func Powershell(c *gin.Context) {
	var requestBody PwshRequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		//TODO: DO SOMETHING WITH THE ERROR
	}
	fmt.Println(requestBody.Path)
	if runtime.GOOS != "windows" {
		c.JSON(403, gin.H{
			"Message":  "Agent not running on Windows",
			"Finished": nil,
		})
	}
	cmd := exec.Command("powershell", "-file", requestBody.Path)
	stdout, err := cmd.Output()
	if err != nil {
		c.JSON(200, gin.H{
			"Message":  "Error occurs",
			"Finished": err,
		})
	} else {
		c.JSON(200, gin.H{
			"Message":  string(stdout),
			"Finished": nil,
		})
	}
}

func sh(c *gin.Context) {
	var requestBody PwshRequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		//TODO: DO SOMETHING WITH THE ERROR
	}
	fmt.Println(requestBody.Path)
	if runtime.GOOS != "linux" {
		c.JSON(403, gin.H{
			"Message":  "Agent not running on Linux",
			"Finished": nil,
		})
	}
	cmd := exec.Command("sh", requestBody.Path)

	stdout, err := cmd.Output()
	if err != nil {
		c.JSON(200, gin.H{
			"Message":  "Error occurs",
			"Finished": err,
		})
	} else {
		c.JSON(200, gin.H{
			"Message":  string(stdout),
			"Finished": nil,
		})
	}
}

func main() {
	cfg = LoadConfiguration("config.json")
	r := gin.Default()
	r.Use(checkAPIKey())
	r.POST("/pwsh", checkAPIKey(), func(c *gin.Context) {
		Powershell(c)
	})
	r.POST("/sh", checkAPIKey(), func(c *gin.Context) {
		sh(c)
	})
	r.RunTLS(":"+strconv.Itoa(cfg.AgentPort), cfg.CertFile, cfg.KeyFile)
}
