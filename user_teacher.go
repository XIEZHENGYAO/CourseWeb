package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	iconv "github.com/djimenez/iconv-go"
)

// 测试用例 0010633 张寅
type Teacher struct {
	ID                 string `bson:"_id"`
	Password           string
	Name               string
	Sex                string
	Introduction       string
	Email              string
	Phone              string
	College            string // 所在学院
	Department         string // 系
	AcademicBackground string // 学历
	AcademicTitle      string // 职称
	ResearchDirections string // 研究方向
	SecurityQuestions  []SecurityQuestion
}

func initTeachers(id int) {
	no := strconv.Itoa(id)
	for len(no) < 7 {
		no = "0" + no
	}
	url := "http://10.202.78.13/html_js/" + no + ".html"
	log.Println(id)
	r, err := http.Get(url)
	if err != nil {
		return
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	input := body
	out := make([]byte, len(input))
	out = out[:]
	iconv.Convert(input, out, "gb2312", "utf-8")
	str := string(out)

	reg, _ := regexp.Compile("[\\s\\S]*<span id=\"xm\">(.*)</span>[\\s\\S]*")
	if s := reg.FindString(str); len(s) == 0 {
		return
	}
	name := reg.ReplaceAllString(str, "$1")
	log.Println(name)

	reg, _ = regexp.Compile("[\\s\\S]*<span id=\"xb\">(.*)</span>[\\s\\S]*")
	sex := reg.ReplaceAllString(str, "$1")
	log.Println(sex)

	reg, _ = regexp.Compile("[\\s\\S]*<span id=\"lxdh\">(.*)</span>[\\s\\S]*")
	phone := reg.ReplaceAllString(str, "$1")
	log.Println(phone)

	reg, _ = regexp.Compile("[\\s\\S]*<span id=\"emldz\">(.*)</span>[\\s\\S]*")
	email := reg.ReplaceAllString(str, "$1")
	log.Println(email)

	reg, _ = regexp.Compile("[\\s\\S]*<span id=\"ks\">(.*)</span>[\\s\\S]*")
	college := reg.ReplaceAllString(str, "$1")
	log.Println(college)

	reg, _ = regexp.Compile("[\\s\\S]*<span id=\"bm\">(.*)</span>[\\s\\S]*")
	department := reg.ReplaceAllString(str, "$1")
	log.Println(department)

	reg, _ = regexp.Compile("[\\s\\S]*<span id=\"xl\">(.*)</span>[\\s\\S]*")
	AcademicBackground := reg.ReplaceAllString(str, "$1")
	log.Println(AcademicBackground)

	reg, _ = regexp.Compile("[\\s\\S]*<span id=\"zc\">(.*)</span>[\\s\\S]*")
	AcademicTitle := reg.ReplaceAllString(str, "$1")
	log.Println(AcademicTitle)

	reg, _ = regexp.Compile("[\\s\\S]*<span id=\"jxyjfx\">(.*)</span>[\\s\\S]*")
	ResearchDirections := reg.ReplaceAllString(str, "$1")
	log.Println(ResearchDirections)

	teacher := &Teacher{
		ID:                 no,
		Password:           no,
		Name:               name,
		Sex:                sex,
		Phone:              phone,
		Email:              email,
		College:            college,
		Department:         department,
		AcademicBackground: AcademicBackground,
		AcademicTitle:      AcademicTitle,
		ResearchDirections: ResearchDirections,
	}
	err = addTeacher(teacher)
	if err != nil {
		log.Println(err)
	}
	// ioutil.WriteFile("globeEarthquake_csn.html", out, 0644)
}

func addTeacher(teacher *Teacher) error {
	return mgoInsert("teacher", &teacher)
}

func removeTeachers(selector map[string]interface{}) (*mgo.ChangeInfo, error) {
	return mgoRemoveAll("teacher", selector)
}

func removeTeacherByID(id string) (*mgo.ChangeInfo, error) {
	return mgoRemove("teacher", bson.M{"_id": id})
}

func updateTeachers(selector map[string]interface{}, update map[string]interface{}) (*mgo.ChangeInfo, error) {
	return mgoUpdateAll("teacher", selector, update)
}

func updateTeacherByID(id string, update map[string]interface{}) error {
	return mgoUpdate("teacher", bson.M{"_id": id}, update)
}

func searchTeachers(selector map[string]interface{}) ([]map[string]interface{}, error) {
	return mgoFindAll("teacher", selector)
}

func getTeachers() ([]map[string]interface{}, error) {
	return mgoFindAll("teacher", nil)
}

func searchTeachersSelect(selector map[string]interface{}) ([]map[string]interface{}, error) {
	return mgoSearchSelect("teacher", selector)
}
func getTeachersByPage(page int) ([]map[string]interface{}, error) {
	return mgoFindByPage("teacher", page)
}

func getTeacherInfo(collection string, id string) (*map[string]interface{}, error) {
	person, err := mgoFind(collection,
		bson.M{"_id": id})
	personInfo := make(map[string]interface{})
	personInfo["id"] = person["_id"]
	personInfo["name"] = person["name"]
	personInfo["sex"] = person["sex"]
	personInfo["email"] = person["email"]
	personInfo["phone"] = person["phone"]
	personInfo["college"] = person["college"]
	personInfo["department"] = person["department"]
	personInfo["academicbackground"] = person["academicbackground"]
	personInfo["academictitle"] = person["academictitle"]
	personInfo["researchdirections"] = person["researchdirections"]
	personInfo["introduction"] = person["introduction"]
	return &personInfo, err
}
