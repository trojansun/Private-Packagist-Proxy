package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Config 对应 TOML 配置文件的顶级结构
type Config struct {
	Storage StorageConfig `toml:"storage"`
	Domain  DomainConfig  `toml:"domain"`
}

// StorageConfig 对应 [storage] 部分
type StorageConfig struct {
	Type string `toml:"type"`
	Path string `toml:"path"`
}

// DomainConfig 对应 [domain] 部分
type DomainConfig struct {
	Original string `toml:"original"`
	Bind     string `toml:"bind"`
}

func main() {

	var config Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println("Error reading TOML config:", err)
		return
	}

	flag.Parse()
	glog.MaxSize = 1024 * 1024 * 10

	glog.Info("Starting the application...")
	router := gin.Default()
	router.GET("/:pathUri", func(c *gin.Context) {
		// 获取 URL 参数
		pathUri := c.Param("pathUri")

		// 如果是packages.json
		if pathUri == "packages.json" {
			servePackagesJSON(fmt.Sprintf("%s/%s", config.Domain.Original, pathUri), filepath.Join(config.Storage.Path, pathUri), c)
			return
		}

		glog.Info("pathUri: ", pathUri)
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	glog.Info("Application is running...")
	router.Run(":8000")
}

// fetchAndSavePackagesJSON 从 Packagist 获取 packages.json 并保存到本地
func fetchAndSavePackagesJSON(remoteUrl string, localPath string) error {
	resp, err := http.Get(remoteUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to fetch packages.json: status code %d", resp.StatusCode)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	// 写文件
	err = os.WriteFile(localPath, body, 0644)
	if err != nil {
		return err
	}

	return nil
}

// servePackagesJSON 服务 packages.json 到客户端
func servePackagesJSON(remoteUrl string, localPath string, c *gin.Context) {
	body, err := os.ReadFile(localPath)
	if err != nil {
		// 如果读取失败，尝试从远程重新获取并保存
		err = fetchAndSavePackagesJSON(remoteUrl, localPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch packages.json"})
			return
		}

		// 再次尝试读取文件
		body, err = os.ReadFile(localPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read packages.json"})
			return
		}
	}

	c.Data(http.StatusOK, "application/json", body)
}
