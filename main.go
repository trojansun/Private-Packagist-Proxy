package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/tidwall/gjson"
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
	Type       string `toml:"type"`
	JsonPath   string `toml:"jsonPath"`
	SourcePath string `toml:"sourcePath"`
}

// DomainConfig 对应 [domain] 部分
type DomainConfig struct {
	Original   string `toml:"original"`
	Bind       string `toml:"bind"`
	PrefixPath string `toml:"prefixPath"`
}

// Package represents the structure of the package entries in JSON
type Package struct {
	Dist struct {
		URL       string `json:"url"`
		Type      string `json:"type"`
		Reference string `json:"reference"`
	} `json:"dist"`
	Version string `json:"version"`
}

// ComposerJSON represents the top-level structure of the JSON with dynamic package keys
type ComposerJSON struct {
	Packages map[string][]Package `json:"packages"`
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
	router.GET("/*pathUri", func(c *gin.Context) {
		// 获取 URL 参数
		//pathUri := c.Param("pathUri")
		reqUrlPath := c.Request.URL.Path
		glog.Info("Request URL Path: ", reqUrlPath)

		// 忽略 favicon.ico 请求
		if reqUrlPath == "/favicon.ico" {
			return
		}

		// 如果是packages.json
		if reqUrlPath == "/packages.json" {
			servePackagesJSON(fmt.Sprintf("%s%s", config.Domain.Original, reqUrlPath), filepath.Join(config.Storage.JsonPath, reqUrlPath), true, c)
			return
		}

		// 如过是packages/list.json
		if reqUrlPath == "/packages/list.json" {
			servePackagesJSON(fmt.Sprintf("%s%s", config.Domain.Original, reqUrlPath), filepath.Join(config.Storage.JsonPath, reqUrlPath), true, c)
			return
		}

		//********************************处理资源相关 S****************************************
		// 不仅要保存json，还要替换资源文件的url
		// 步骤
		// 1. 先获取json文件
		// 2. 解析json文件，获取资源文件的url
		// 3. 根据文件类型，下载资源，并保存到本地
		// 4. 替换json文件中的url
		// 5. 保存json文件

		// 如果是/p/开头的 composer v1
		if len(reqUrlPath) > 3 && reqUrlPath[:3] == "/p/" {
			servePackagesJSON(fmt.Sprintf("%s%s", config.Domain.Original, reqUrlPath), filepath.Join(config.Storage.JsonPath, reqUrlPath), true, c)
			return
		}

		// 如果是/p2/开头的 composer v2
		if len(reqUrlPath) > 4 && reqUrlPath[:4] == "/p2/" {
			servePackagesJSON(fmt.Sprintf("%s%s", config.Domain.Original, reqUrlPath), filepath.Join(config.Storage.JsonPath, reqUrlPath), false, c)
			// 定位包名称,假设是/p2/laravel/laravel.json，替换掉里面的/p2/和.json
			_ = reqUrlPath[4 : len(reqUrlPath)-5]
			// 读取本地的json文件
			readFile, err := os.ReadFile(filepath.Join(config.Storage.JsonPath, reqUrlPath))
			if err != nil {
				return
			}
			jsonContent := string(readFile)
			gsonContent := gjson.Get(jsonContent, "packages")
			for packageName := range gsonContent.Map() {
				glog.Info(packageName)
				formattedKey := fmt.Sprintf(`%s`, packageName)
				for i := range gsonContent.Get(formattedKey).Array() {
					glog.Info(gsonContent.Get(formattedKey).Array()[i])
					// 要注意，URL不能为空或者三（四）项不能为空
					//glog.Info(jsonData.Packages[packageName][i].Dist.URL)
					////jsonData.Packages[packageName][i].Dist.URL = fmt.Sprintf("%s%s", config.Domain.Bind, jsonData.Packages[packageName][i].Dist.URL)
					//
					//glog.Info("version is ", jsonData.Packages[packageName][i].Version)
					//glog.Info("url is ", jsonData.Packages[packageName][i].Dist.URL)
					//glog.Info("type is ", jsonData.Packages[packageName][i].Dist.Type)
					//glog.Info("reference is ", jsonData.Packages[packageName][i].Dist.Reference)
					//
					//localSourceFileName := fmt.Sprintf("%s-%s-%s.%s", strings.ReplaceAll(packageName, "/", "-"),
					//	jsonData.Packages[packageName][i].Version,
					//	jsonData.Packages[packageName][i].Dist.Reference,
					//	jsonData.Packages[packageName][i].Dist.Type,
					//)
					//// 字符串packageName的"/"替换为"-"
					//// url + prefix + packageName + localSourceFileName
					//oldUrl := jsonData.Packages[packageName][i].Dist.URL
					//newUrl := fmt.Sprintf("%s%s/%s/%s",
					//	config.Domain.Bind,
					//	config.Domain.PrefixPath,
					//	packageName,
					//	localSourceFileName,
					//)
					//jsonData.Packages[packageName][i].Dist.URL = newUrl
					//// 判断文件是否存在
					//if _, err := os.Stat(fmt.Sprintf("%s/%s/%s", config.Storage.SourcePath, packageName, localSourceFileName)); err == nil {
					//	glog.Info("文件已经存在")
					//} else {
					//	servePackagesJSON(oldUrl, fmt.Sprintf("%s/%s/%s", config.Storage.SourcePath, packageName, localSourceFileName), false, c)
					//}
				}
			}
			// 保存json文件
			//jsonDataBytes, err := json.Marshal(jsonData)
			//if err != nil {
			//	return
			//}
			//err = os.WriteFile(filepath.Join(config.Storage.JsonPath, reqUrlPath), jsonDataBytes, 0644)
			//if err != nil {
			//	return
			//}
			//
			//// TODO 这个要把其他的json文件也要处理，否则无法运行的
			//
			//// 返回json文件
			//c.Data(http.StatusOK, "application/json", jsonDataBytes)
			return
		}

		// 如果是`config.Domain.prefixPath`开头的资源文件
		if len(reqUrlPath) > len(config.Domain.PrefixPath) && reqUrlPath[:len(config.Domain.PrefixPath)] == config.Domain.PrefixPath {
			// 获取资源文件的路径
			sourceFilePath := filepath.Join(config.Storage.SourcePath, reqUrlPath[len(config.Domain.PrefixPath)+1:])
			c.File(sourceFilePath)
		}

		//********************************处理资源相关 E****************************************

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	glog.Info("Application is running...")
	router.Run(":8000")
}

// fetchAndSavePackagesJSON 从 Packagist 获取 packages.json 并保存到本地
func fetchAndSavePackagesJSON(remoteUrl string, jsonPath string) error {
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
	if err := os.MkdirAll(filepath.Dir(jsonPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	// 写文件
	err = os.WriteFile(jsonPath, body, 0644)
	if err != nil {
		return err
	}

	return nil
}

// servePackagesJSON 服务 packages.json 到客户端
func servePackagesJSON(remoteUrl string, jsonPath string, directReturn bool, c *gin.Context) {
	body, err := os.ReadFile(jsonPath)
	if err != nil {
		// 如果读取失败，尝试从远程重新获取并保存
		err = fetchAndSavePackagesJSON(remoteUrl, jsonPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch packages.json"})
			return
		}

		// 再次尝试读取文件
		body, err = os.ReadFile(jsonPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read packages.json"})
			return
		}
	}
	if directReturn {
		c.Data(http.StatusOK, "application/json", body)
		return
	}
}
