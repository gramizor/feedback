package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	file, err := os.Open("resources/data/group.json")
	if err != nil {
		fmt.Println("Ошибка при открытии JSON файла:", err)
		return
	}
	defer file.Close()

	var groups []Group
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&groups); err != nil {
		fmt.Println("Ошибка при декодировании JSON данных:", err)
		return
	}

	r.LoadHTMLGlob("templates/*")

	r.Static("/images", "./resources/images")
	r.Static("/fonts", "./resources/fonts")
	r.Static("/data", "./resources/data")
	r.Static("/css", "./resources/css")

	r.GET("/", func(c *gin.Context) {
		data := gin.H{
			"groups": groups,
		}
		c.HTML(http.StatusOK, "mainPage.tmpl", data)
	})

	r.GET("/groups/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Print(err)
		}
		group := groups[id-1]
		c.HTML(http.StatusOK, "group.tmpl", group)
	})

	r.GET("/search", func(c *gin.Context) {
		searchQuery := c.DefaultQuery("groupSearch", "")
		var foundGroups []Group
		for _, group := range groups {
			if strings.HasPrefix(strings.ToLower(group.Name), strings.ToLower(searchQuery)) {
				foundGroups = append(foundGroups, group)
			}
		}
		data := gin.H{
			"groups":      foundGroups,
			"searchQuery": searchQuery,
		}
		c.HTML(http.StatusOK, "mainPage.tmpl", data)
	})

	r.Run()
	// go run cmd/main.go

	log.Println("Server down")
}
