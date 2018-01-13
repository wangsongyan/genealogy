package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wangsongyan/genealogy/models"
)

func CoupleNew(context *gin.Context) {

}

func CoupleCreate(c *gin.Context) {
	var couple models.Couple
	if err := c.Bind(&couple); err == nil {
		if err := couple.Insert(); err == nil {
			c.JSON(http.StatusOK, gin.H{"msg": "succeed."})
		} else {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func GetFatherTree(peopleId uint) (uint, error) {
	people, err := models.GetPeopleById(peopleId)
	if err != nil {
		return 0, nil
	}
	if people.CoupleId == 0 {
		cur, err := models.GetCoupleByPeopleId(peopleId)
		if err != nil {
			return 0, err
		}
		return cur.ID, nil
	}
	couple, err := models.GetCoupleById(people.CoupleId)
	if err != nil {
		return 0, nil
	}
	return GetFatherTree(couple.HusbandId)
}

func GetMotherTree(peopleId uint) (uint, error) {
	people, err := models.GetPeopleById(peopleId)
	if err != nil {
		return 0, nil
	}
	if people.CoupleId == 0 {
		cur, err := models.GetCoupleByPeopleId(peopleId)
		if err != nil {
			return 0, err
		}
		return cur.ID, nil
	}
	couple, err := models.GetCoupleById(people.CoupleId)
	if err != nil {
		return 0, nil
	}
	return GetMotherTree(couple.WifeId)
}
