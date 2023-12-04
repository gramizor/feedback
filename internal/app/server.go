package app

import (
	"feedback/internal/app/ds"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (a *Application) StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.Static("/images", "./resources/images")
	r.Static("/fonts", "./resources/fonts")
	r.Static("/data", "./resources/data")
	r.Static("/css", "./resources/css")

	r.GET("/", func(c *gin.Context) {
		groups, err := a.repo.GetActiveGroups()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		data := gin.H{
			"groups": groups,
		}
		c.HTML(http.StatusOK, "mainPage.tmpl", data)
	})

	r.GET("/groups/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		group, err := a.repo.GetActiveGroupById(id)
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.HTML(http.StatusOK, "group.tmpl", group)
	})

	r.GET("/search", func(c *gin.Context) {
		groups, err := a.repo.GetActiveGroups()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		groupSlice := *groups
		searchQuery := c.DefaultQuery("groupSearch", "")
		var foundGroups []ds.Group
		for _, group := range groupSlice {
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

	r.POST("/delete", func(c *gin.Context) {
		id, err := strconv.Atoi(c.DefaultQuery("delete", ""))
		log.Print(c.DefaultQuery("delete", ""))
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = a.repo.DeactivateGroupByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		groups, err := a.repo.GetActiveGroups()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		data := gin.H{
			"groups": groups,
		}
		c.HTML(http.StatusOK, "mainPage.tmpl", data)
	})

	r.Run()
	// go run cmd/feedback/main.go

	log.Println("Server down")
}
