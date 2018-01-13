package controllers

import (
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wangsongyan/genealogy/models"
)

func IndexGet(c *gin.Context) {
	couples, _ := models.GetRootCoupleNode()
	for _, couple := range couples {
		node, _ := GetCouple(couple.CoupleId)
		nodes := make([]*models.Node, 0)
		nodes = append(nodes, node)
		s, d := sumPeople(nodes)
		couple.Sum = s
		couple.Alive = s - d
	}
	c.HTML(http.StatusOK, "index", gin.H{
		"couples": couples,
	})
}

func CoupleGet(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"type":     "1",
		"coupleId": c.Param("id"),
	})
}

func CouplePost(c *gin.Context) {
	cid := c.Param("id")
	coupleId, err := strconv.ParseUint(cid, 10, 64)
	if err == nil {
		var node *models.Node
		node, err = GetCouple(uint(coupleId))
		if err == nil {
			nodes := make([]*models.Node, 0)
			nodes = append(nodes, node)
			sum, death := sumPeople(nodes)
			c.JSON(http.StatusOK, gin.H{
				"couple": nodes,
				"sum":    sum,
				"alive":  sum - death,
				"death":  death,
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"error": err.Error()})
}

func GetFamilyTreeByPeopleId(c *gin.Context) {
	pid := c.Param("id")
	peopleId, err := strconv.ParseUint(pid, 10, 64)
	if err == nil {
		fatherTreeRootId, _ := GetFatherTree(uint(peopleId))
		motherTreeRootId, _ := GetMotherTree(uint(peopleId))
		c.HTML(http.StatusOK, "index.html", gin.H{
			"type":             "2",
			"fatherTreeRootId": fatherTreeRootId,
			"motherTreeRootId": motherTreeRootId,
			"peopleId":         pid,
		})
	}

}

func GetCouple(coupleId uint) (*models.Node, error) {
	node, err := models.GetCoupleNodeById(coupleId)
	if err != nil {
		return nil, err
	}
	children, err := models.GetPeopleByCoupId(coupleId)
	if err != nil {
		return nil, err
	}
	node.Children = make([]*models.Node, 0)
	for _, c := range children {
		couple, err := models.GetCoupleByPeopleId(c.ID)
		if err != nil {
			n := new(models.Node)
			n.RelationId = c.ID
			n.Single = true
			if c.Gender == 0 {
				n.HusbandId = c.ID
				n.HusbandName = c.Name
				n.HusbandAlias = c.Alias
			} else {
				n.WifeId = c.ID
				n.WifeName = c.Name
				n.WifeAlias = c.Alias
			}
			node.Children = append(node.Children, n)
		} else {
			granChildren, err := GetCouple(couple.ID)
			if err == nil {
				granChildren.RelationId = c.ID
				node.Children = append(node.Children, granChildren)
			}
		}
	}
	return node, nil
}

func sumPeople(nodes []*models.Node) (int, int) {
	sum := 0
	death := 0
	for _, n := range nodes {
		if n.Single {
			sum++
		} else {
			sum += 2
		}
		if n.HusbandStatus == 1 {
			death++
		}
		if n.WifeStatus == 1 {
			death++
		}
		s, d := sumPeople(n.Children)
		sum += s
		death += d
	}
	return sum, death
}
