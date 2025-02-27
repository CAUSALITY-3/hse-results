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

	sub1, _ := json.Marshal(student.Subject1)
	studentTemplateMapping.Name = student.Name
	studentTemplateMapping.RollNo = student.RollNo
	studentTemplateMapping.Regular = student.Regular
	studentTemplateMapping.StudentGroup = student.Group
	studentTemplateMapping.SchoolName = student.SchoolName
	studentTemplateMapping.SchoolCode = student.SchoolCode
	studentTemplateMapping.FullAplus = true
	studentTemplateMapping.Xvg = string(sub1)
	if student.Subject1.Name != "" {
		studentTemplateMapping.Subject1 = student.Subject1.Name
		studentTemplateMapping.Subject1Mark = student.Subject1.Marks
		studentTemplateMapping.Subject1Grade = student.Subject1.Grade
		if student.Subject1.Grade != "A+" {
			studentTemplateMapping.FullAplus = false
		}
	} else {
		studentTemplateMapping.Hide1 = "hide-subject"
		studentTemplateMapping.FullAplus = false
	}
	if student.Subject2.Name != "" {
		studentTemplateMapping.Subject2 = student.Subject2.Name
		studentTemplateMapping.Subject2Mark = student.Subject2.Marks
		studentTemplateMapping.Subject2Grade = student.Subject2.Grade
		if student.Subject2.Grade != "A+" {
			studentTemplateMapping.FullAplus = false
		}
	} else {
		studentTemplateMapping.Hide2 = "hide-subject"
		studentTemplateMapping.FullAplus = false
	}
	if student.Subject3.Name != "" {
		studentTemplateMapping.Subject3 = student.Subject3.Name
		studentTemplateMapping.Subject3Mark = student.Subject3.Marks
		studentTemplateMapping.Subject3Grade = student.Subject3.Grade
		if student.Subject3.Grade != "A+" {
			studentTemplateMapping.FullAplus = false
		}
	} else {
		studentTemplateMapping.Hide3 = "hide-subject"
		studentTemplateMapping.FullAplus = false
	}
	if student.Subject4.Name != "" {
		studentTemplateMapping.Subject4 = student.Subject4.Name
		studentTemplateMapping.Subject4Mark = student.Subject4.Marks
		studentTemplateMapping.Subject4Grade = student.Subject4.Grade
		if student.Subject4.Grade != "A+" {
			studentTemplateMapping.FullAplus = false
		}
	} else {
		studentTemplateMapping.Hide4 = "hide-subject"
		studentTemplateMapping.FullAplus = false
	}
	if student.Subject5.Name != "" {
		studentTemplateMapping.Subject5 = student.Subject5.Name
		studentTemplateMapping.Subject5Mark = student.Subject5.Marks
		studentTemplateMapping.Subject5Grade = student.Subject5.Grade
		if student.Subject5.Grade != "A+" {
			studentTemplateMapping.FullAplus = false
		}
	} else {
		studentTemplateMapping.Hide5 = "hide-subject"
		studentTemplateMapping.FullAplus = false
	}
	if student.Subject6.Name != "" {
		studentTemplateMapping.Subject6 = student.Subject6.Name
		studentTemplateMapping.Subject6Mark = student.Subject6.Marks
		studentTemplateMapping.Subject6Grade = student.Subject6.Grade
		if student.Subject6.Grade != "A+" {
			studentTemplateMapping.FullAplus = false
		}
	} else {
		studentTemplateMapping.Hide6 = "hide-subject"
		studentTemplateMapping.FullAplus = false
	}

	tmpl, err := template.ParseFiles("./public/" + requestType + "/students.html")
	if err != nil {
		log.Println("Error loading template:", err)
		return err
	}
	log.Println("FullAplus", studentTemplateMapping.FullAplus)
	c.Set("Content-Type", "text/html")
	err = tmpl.Execute(c.Response().BodyWriter(), studentTemplateMapping)
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
