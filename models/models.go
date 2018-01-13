package models

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type BaseModel struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type People struct {
	BaseModel
	Name         string `form:"name"`                      //姓名
	Alias        string `form:"alias"`                     //别名
	Gender       uint   `form:"gender"`                    //性别
	Identitycard string `form:"identitycard"`              //身份证号
	Telephone    string `form:"telephone"`                 //电话
	Status       uint   `form:"status" gorm:"default:'0'"` //状态
	Photo        string `form:"photo"`                     //照片
	CoupleId     uint   `form:"coupleId"`                  //夫妻id
	//FatherId     uint   //父亲id
	//MotherId     uint   //母亲id
	OrderBy uint `form:"orderBy"` //排序
}

type Couple struct {
	BaseModel
	HusbandId  uint      `form:"husband_id"`                          //丈夫id
	WifeId     uint      `form:"wife_id"`                             //妻子id
	WeddingDay time.Time `form:"weddingDay" time_format:"2006-01-02"` //结婚日
	Status     uint      `form:"status" gorm:"default:'0'"`           //状态
}

type Node struct {
	RelationId    uint
	CoupleId      uint
	HusbandId     uint
	HusbandName   string
	HusbandAlias  string
	HusbandStatus uint
	WifeId        uint
	WifeName      string
	WifeAlias     string
	WifeStatus    uint
	Children      []*Node
	WeddingStatus uint
	Single        bool
	Sum           int
	Alive         int
}

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	//db, err := gorm.Open("sqlite3", "wblog.db?_loc=Asia/Shanghai")
	db, err := gorm.Open("mysql", "root:mysql@tcp(localhost:3306)/genealogy?charset=utf8&parseTime=True&loc=Local")
	if err == nil {
		DB = db
		db.LogMode(true)
		db.AutoMigrate(&People{}, &Couple{})
		db.Model(&People{}).AddUniqueIndex("uk_identitycard", "identitycard")
		db.Model(&Couple{}).AddUniqueIndex("uk_husband_wife", "husband_id", "wife_id")
		return db, err
	}
	return nil, err
}

func (p *People) Insert() error {
	return DB.FirstOrCreate(p, "identitycard = ?", p.Identitycard).Error
}

func (c *Couple) Insert() error {
	return DB.FirstOrCreate(c, "husband_id = ? and wife_id = ?", c.HusbandId, c.WifeId).Error
}

func GetCoupleById(coupleId uint) (*Couple, error) {
	var c Couple
	err := DB.Model(&Couple{}).Find(&c, "id = ?", coupleId).Error
	if err == nil {
		return &c, nil
	}
	return nil, err
}

func GetPeopleById(peopleId uint) (*People, error) {
	var p People
	err := DB.Model(&People{}).Find(&p, "id = ?", peopleId).Error
	if err == nil {
		return &p, nil
	}
	return nil, err
}

func GetCoupleByPeopleId(peopleId uint) (*Couple, error) {
	var c Couple
	err := DB.Model(&Couple{}).Find(&c, "husband_id = ? or wife_id = ?", peopleId, peopleId).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func GetPeopleByCoupId(coupleId uint) ([]*People, error) {
	var p []*People
	err := DB.Order("order_by asc").Where("couple_id = ?", coupleId).Find(&p).Error
	if err == nil {
		return p, nil
	}
	return nil, err
}

func GetCoupleNodeById(coupleId uint) (*Node, error) {
	sql := `
		select c.id couple_id,c.status wedding_status,
		h.id husband_id,h.name husband_name,h.alias husband_alias,h.status husband_status,
		w.id wife_id,w.name wife_name,w.gender,w.alias wife_alias,w.status wife_status
		from couples c inner join peoples h
		on c.husband_id = h.id
		inner join peoples w
		on c.wife_id = w.id where c.id = ?
	`
	rows, err := DB.Raw(sql, coupleId).Rows() // (*sql.Rows, error)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var node Node
	for rows.Next() {
		DB.ScanRows(rows, &node)
	}
	return &node, nil
}

func GetRootCoupleNode() ([]*Node, error) {
	sql := `select c.id couple_id ,h.name husband_name,w.name wife_name from couples c inner join
peoples h on c.husband_id = h.id
inner join peoples w
on c.wife_id = w.id
where h.couple_id = 0 and w.couple_id = 0`
	rows, err := DB.Raw(sql).Rows() // (*sql.Rows, error)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	nodes := make([]*Node, 0)
	for rows.Next() {
		var n Node
		DB.ScanRows(rows, &n)
		nodes = append(nodes, &n)
	}
	return nodes, nil
}
