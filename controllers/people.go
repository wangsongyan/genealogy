package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wangsongyan/genealogy/models"
)

func PeopleNew(context *gin.Context) {

}

func PeopleCreate(c *gin.Context) {
	var p models.People
	if err := c.Bind(&p); err == nil {
		if len(p.Identitycard) == 0 {
			p.Identitycard = p.Name
		}
		if err := p.Insert(); err == nil {
			c.JSON(http.StatusOK, gin.H{"msg": "succeed."})
		} else {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
