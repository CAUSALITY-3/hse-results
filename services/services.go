package services

import (
	"encoding/json"
	"fmt"
	"hse-results/types"
	"hse-results/utils"
	"html/template"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// type StudentData struct {
//     Thiruvananthapuram map[string][]string `json:"THIRUVANANTHAPURAM"`
//     Kollam    map[string][]string `json:"KOLLAM"`
// 	Pathanamthitta map[string][]string `json:"PATHANAMTHITTA"`
//     Alappuzha    map[string][]string `json:"ALAPPUZHA"`
// 	Kottayam map[string][]string `json:"KOTTAYAM"`
//     Idukky    map[string][]string `json:"IDUKKY"`
// 	Ernakulam map[string][]string `json:"ERNAKULAM"`
//     Thrissur    map[string][]string `json:"THRISSUR"`
// 	Palakkad map[string][]string `json:"PALAKKAD"`
// 	Kozhikode map[string][]string `json:"KOZHIKODE"`
//     Malapuram    map[string][]string `json:"MALAPPURAM"`
// 	Wayanad map[string][]string `json:"WAYANAD"`
//     Kannur    map[string][]string `json:"KANNUR"`
// 	Kasaragod map[string][]string `json:"KASARAGOD"`
//     MiddleEastern    map[string][]string `json:"MIDDLE-EASTERN"`
// 	Lakshadweep map[string][]string `json:"LAKSHADWEEP"`
//     Mahe    map[string][]string `json:"MAHE"`
// }

type StudentData map[string]map[string][]string

func SearchStudentByName(name string) ([]string, error) {
	formattedName := strings.ToUpper(name)

	// content := utils.SingletonInjector.Get[StudentData]("students")

	// if content == nil {
	// 	fmt.Println("from file")
	content, err := utils.ReadFile[StudentData]("districtWise-student-school_mapping.json")

	if err != nil {
		return nil, err
	}
	// 	utils.SingletonInjector.Bind("students", content)
	// }

	// districts := []string{
	// 	"Thiruvananthapuram",
	// 	"Kollam",
	// 	"Pathanamthitta",
	// 	"Alappuzha",
	// 	"Kottayam",
	// 	"Idukky",
	// 	"Ernakulam",
	// 	"Thrissur",
	// 	"Palakkad",
	// 	"Kozhikode",
	// 	"Malapuram",
	// 	"Wayanad",
	// 	"Kannur",
	// 	"Kasaragod",
	// 	"MiddleEastern",
	// 	"Lakshadweep",
	// 	"Mahe",
	// }
	var allStudentsStartingWithH []string
	for _, district := range *content {
		for _, students := range district {
			for _, student := range students {
				studentName := strings.Split(student, ",")[0]
				if strings.Contains(studentName, formattedName) {
					allStudentsStartingWithH = append(allStudentsStartingWithH, student)
				}
			}
		}

	}
	fmt.Println("total Count ", len(allStudentsStartingWithH))
	return allStudentsStartingWithH, nil

}

func GetStudentResults(c *fiber.Ctx, requestType, rollNo string) error {

	student, err := utils.ReadFile[types.Student]("./public/" + requestType + "/data/students/" + rollNo + ".json")
	if err != nil {
		log.Println("Error Reading student data:", err)
		return err
	}

	var studentTemplateMapping types.TemplateMappedStudentData

	studentTemplateMapping.Name = student.Name
	studentTemplateMapping.RollNo = student.RollNo
	studentTemplateMapping.Regular = student.Regular
	studentTemplateMapping.StudentGroup = student.Group
	studentTemplateMapping.SchoolName = student.SchoolName
	studentTemplateMapping.SchoolCode = student.SchoolCode
	studentTemplateMapping.FullAp = student.FullAp
	studentTemplateMapping.Result = student.Result
	studentTemplateMapping.TotalMarks = student.TotalMarks
	// studentTemplateMapping.Xvg = string(sub1)

	strTemp := ""

	for i := 0; i <= 5; i++ {
		if student.Subjects[i].Name != "" {
			strTemp += "<tr><td>" + student.Subjects[i].Name + "</td><td>" + fmt.Sprint(*student.Subjects[i].Marks) + "</td><td>" + student.Subjects[i].Grade + "</td></tr>"
		}
	}

	studentTemplateMapping.Subjects = template.HTML(strTemp)

	tmpl, err := template.ParseFiles("./public/" + requestType + "/students.html")
	if err != nil {
		log.Println("Error loading template:", err)
		return err
	}
	log.Println("FullAplus", tmpl)
	c.Set("Content-Type", "text/html")
	err = tmpl.Execute(c.Response().BodyWriter(), studentTemplateMapping)
	if err != nil {
		log.Println("Error executing template:", err)
		return c.Status(500).SendString("Error rendering template")
	}
	return nil
}

func GetSchoolResults(c *fiber.Ctx, requestType, schoolCode string) error {
	school, err := utils.ReadFile[types.SchoolResult]("./public/" + requestType + "/data/schools/" + schoolCode + ".json")
	if err != nil {
		log.Println("Error Reading student data:", err)
		return err
	}

	// var schoolTemplateMapping types.SchoolResults
	tmpl, err := template.ParseFiles("./public/" + requestType + "/school.html")
	if err != nil {
		log.Println("Error loading template:", err)
		return err
	}
	var temp struct {
		SchoolResult string
		Name         string
		SchoolCode   string
	}

	// SchoolResult, err := json.Marshal(school)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return err
	// }
	// fmt.Println("SchoolResult", string(SchoolResult))

	SchoolResultJSON, err := json.Marshal(school.Results)
	if err != nil {
		log.Println("Error marshaling school results:", err)
		return c.Status(500).SendString("Error processing school results")
	}
	fmt.Println("SchoolResult", string(SchoolResultJSON))
	temp.SchoolResult = string(SchoolResultJSON)
	temp.Name = school.SchoolName
	temp.SchoolCode = school.SchoolCode
	c.Set("Content-Type", "text/html")
	err = tmpl.Execute(c.Response().BodyWriter(), temp)
	if err != nil {
		log.Println("Error executing template:", err)
		return c.Status(500).SendString("Error rendering template")
	}
	return nil

}

// func GetSchoolResults(c *fiber.Ctx, requestType, schoolCode string) error {

// 	student, err := utils.ReadFile[types.SchoolResult]("./public/" + requestType + "/data/schools/" + schoolCode + ".json")
// 	if err != nil {
// 		log.Println("Error Reading student data:", err)
// 		return err
// 	}

// 	var studentTemplateMapping types.TemplateMappedStudentData

// 	tmpl, err := template.ParseFiles("./public/" + requestType + "/students.html")
// 	if err != nil {
// 		log.Println("Error loading template:", err)
// 		return err
// 	}
// 	log.Println("FullAplus", studentTemplateMapping.FullAplus)
// 	c.Set("Content-Type", "text/html")
// 	err = tmpl.Execute(c.Response().BodyWriter(), studentTemplateMapping)
// 	if err != nil {
// 		log.Println("Error executing template:", err)
// 		return c.Status(500).SendString("Error rendering template")
// 	}
// 	return nil
// }
