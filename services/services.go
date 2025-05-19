package services

import (
	"encoding/json"
	"fmt"
	"hse-results/types"
	"hse-results/utils"
	"html/template"
	"log"
	"math"
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

	student, err := utils.ReadFile[types.Student]("./public/" + requestType + "/students/" + rollNo + ".json")
	if err != nil {
		log.Println("Error Reading student data:", err)
		return c.Status(500).SendString("Error reading student data")
	}

	var studentTemplateMapping types.TemplateMappedStudentData

	// Safely map student data to avoid nil pointer dereferences
	studentTemplateMapping.Name = student.Name
	studentTemplateMapping.RollNo = student.RollNo
	studentTemplateMapping.Regular = student.Regular
	studentTemplateMapping.StudentGroup = student.Group
	studentTemplateMapping.SchoolName = student.SchoolName
	studentTemplateMapping.SchoolCode = student.SchoolCode
	studentTemplateMapping.FullAp = student.FullAp
	studentTemplateMapping.TotalMarks = student.TotalMarks
	studentTemplateMapping.SchoolRankCode = "school_" + student.SchoolCode
	studentTemplateMapping.Group = student.Group
	studentTemplateMapping.MainRoute = requestType

	if student.Result == "EHS" {
		studentTemplateMapping.Result = "Passed"
	} else {
		studentTemplateMapping.Result = "Failed"
	}

	if student.TotalMarks > 0 {
		studentTemplateMapping.Percentage = math.Round((float64(student.TotalMarks)/1200)*100*100) / 100.0
	} else {
		studentTemplateMapping.Percentage = 0
	}

	strTemp := ""

	for i := 0; i <= 5; i++ {
		if i < len(student.Subjects) && student.Subjects[i].Name != "" {
			marks := "N/A"
			if student.Subjects[i].Marks != nil {
				marks = fmt.Sprint(*student.Subjects[i].Marks)
			}
			strTemp += "<tr><td>" + student.Subjects[i].Name + "</td><td>" + marks + "</td><td>" + student.Subjects[i].Grade + "</td></tr>"
		}
	}

	studentTemplateMapping.Subjects = template.HTML(strTemp)

	tmpl, err := template.ParseFiles("./public/pages/students.html")
	if err != nil {
		log.Println("Error loading template:", err)
		return c.Status(500).SendString("Error loading template")
	}

	c.Set("Content-Type", "text/html")
	err = tmpl.Execute(c.Response().BodyWriter(), studentTemplateMapping)
	if err != nil {
		log.Println("Error executing template:", err)
		return c.Status(500).SendString("Error rendering template")
	}
	return nil
}

func GetSchoolResults(c *fiber.Ctx, requestType, schoolCode string) error {
	school, err := utils.ReadFile[types.SchoolResult]("./public/" + requestType + "/schools/" + schoolCode + ".json")
	if err != nil {
		log.Println("Error Reading school data:", err)
		return c.Status(500).SendString("Error reading school data")
	}

	// var schoolTemplateMapping types.SchoolResults
	tmpl, err := template.ParseFiles("./public/pages/school.html")
	if err != nil {
		log.Println("Error loading template:", err)
		return c.Status(500).SendString("Error loading template")
	}
	var temp types.TemplateMappedSchoolResultData

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
	temp.SchoolResult = string(SchoolResultJSON)
	temp.SchoolName = school.SchoolName
	temp.SchoolCode = school.SchoolCode
	temp.MainRoute = requestType
	temp.District = school.District
	temp.SchoolRankCode = "school_" + school.SchoolCode
	temp.TotalStudentsAppeared = len(school.RankList)
	temp.TotalStudentsPassed = school.TotalPass
	temp.TotalStudentsFailed = school.TotalFailure
	temp.PassPercentage = int(math.Round(school.PassPercentage*100) / 100.0)
	temp.TotalFullAPlus = school.TotalFullAp
	temp.TotalFullMarks = school.TotalFullMarks
	c.Set("Content-Type", "text/html")
	if err := tmpl.Execute(c.Response().BodyWriter(), temp); err != nil {
		log.Println("Error executing template:", err)
		return c.Status(500).SendString("Error rendering template")
	}
	return nil

}

func GetSchoolResultsServerSide(c *fiber.Ctx, requestType, schoolCode string) error {
	school, err := utils.ReadFile[types.SchoolResult]("./public/" + requestType + "/schools/" + schoolCode + ".json")
	if err != nil {
		log.Println("Error Reading school data:", err)
		return c.Status(500).SendString("Error reading school data")
	}

	// var schoolTemplateMapping types.SchoolResults
	tmpl, err := template.ParseFiles("./public/pages/school.html")
	if err != nil {
		log.Println("Error loading template:", err)
		return c.Status(500).SendString("Error loading template")
	}
	var temp types.TemplateMappedSchoolResultData

	SchoolResultJSON, err := json.Marshal(school.Results)
	if err != nil {
		log.Println("Error marshaling school results:", err)
		return c.Status(500).SendString("Error processing school results")
	}
	temp.SchoolResult = string(SchoolResultJSON)
	temp.SchoolName = school.SchoolName
	temp.SchoolCode = school.SchoolCode
	temp.MainRoute = requestType
	temp.SchoolRankCode = "school_" + school.SchoolCode
	temp.District = school.District
	temp.TotalStudentsAppeared = len(school.RankList)
	temp.TotalStudentsPassed = school.TotalPass
	temp.TotalStudentsFailed = school.TotalFailure
	temp.PassPercentage = int(math.Round(school.PassPercentage*100) / 100.0)
	temp.TotalFullAPlus = school.TotalFullAp
	temp.TotalFullMarks = school.TotalFullMarks
	temp.DisplayRank = "true"
	temp.DisplayToppers = "true"
	temp.DisplayFullAplus = "display-none"
	temp.DisplayFullMark = "display-none"
	var tableRows strings.Builder
	if len(school.Results) == 0 {
		temp.DisplayRank = "display-none"
		tableRows.WriteString("<tr><td colspan='10'>No data available</td></tr>")
	} else {
		for i, student := range school.Results {
			rowClass := "even-row"
			if i%2 != 0 {
				rowClass = "odd-row"
			}
			tableRows.WriteString(fmt.Sprintf(
				`<tr class="%s">
                    <td class="school-student-redirect" hx-get="/%s/search/student/%s" hx-swap="outerHTML" hx-target="#school-page" hx-select="#students-page" hx-push-url="true">%s</td>
                    <td class="school-student-redirect" hx-get="/%s/search/student/%s" hx-swap="outerHTML" hx-target="#school-page" hx-select="#students-page" hx-push-url="true">%s</td>
                    <td>%s</td>`,
				rowClass,
				requestType, student.RollNo, student.RollNo,
				requestType, student.RollNo, student.Name,
				student.Group,
			))
			for _, subject := range student.Subjects {
				subjectCell := ""
				if subject.Name != "" {
					marks := ""
					if subject.Marks != nil {
						marks = fmt.Sprintf("%v", *subject.Marks)
					}
					subjectCell = fmt.Sprintf(
						`%s: <span class='marks'>%s</span> (<span class='grade'>%s</span>)`,
						subject.Name, marks, subject.Grade,
					)
				}
				tableRows.WriteString(fmt.Sprintf("<td>%s</td>", subjectCell))
			}
			resultClass := "fail"
			if student.Result == "EHS" {
				resultClass = "pass"
			}
			tableRows.WriteString(fmt.Sprintf(`<td class="%s">%s</td></tr>`, resultClass, student.Result))
		}
	}
	temp.TableRows = template.HTML(tableRows.String())

	var fullMarkStudentsRows strings.Builder
	var fullApStudentsRows strings.Builder
	if len(school.FullMarkStudents) == 0 && len(school.FullApStudents) == 0 {
		temp.DisplayToppers = "display-none"
	} else {
		if len(school.FullMarkStudents) > 0 {
			temp.DisplayFullMark = "true"
			for _, student := range school.FullMarkStudents {
				studentRollNo := utils.ExtractRollNo(student)
				fullMarkStudentsRows.WriteString(fmt.Sprintf(
					`<div class="school-student-redirect school-topper-student" hx-get="/%s/search/student/%s" hx-swap="outerHTML" hx-target="#school-page" hx-select="#students-page" hx-push-url="true">%s</div>`,
					requestType, studentRollNo, student,
				))
			}
		}
		if len(school.FullApStudents) > 0 {
			temp.DisplayFullAplus = "true"
			for _, student := range school.FullApStudents {
				studentRollNo := utils.ExtractRollNo(student)
				fullApStudentsRows.WriteString(fmt.Sprintf(
					`<div class="school-student-redirect school-topper-student" hx-get="/%s/search/student/%s" hx-swap="outerHTML" hx-target="#school-page" hx-select="#students-page" hx-push-url="true">%s</div>`,
					requestType, studentRollNo, student,
				))
			}
		}
	}
	temp.FullMarkStudents = template.HTML(fullMarkStudentsRows.String())
	temp.FullApStudents = template.HTML(fullApStudentsRows.String())

	var rankTableRows strings.Builder
	if len(school.RankList) > 0 {
		for i, student := range school.RankList {
			rowClass := "even-row"
			if i%2 != 0 {
				rowClass = "odd-row"
			}
			rankTableRows.WriteString(fmt.Sprintf(
				`<tr class="%s">
                    <td class="school-student-redirect" hx-get="/%s/search/student/%s" hx-swap="outerHTML" hx-target="#school-page" hx-select="#students-page" hx-push-url="true">%s</td>
                    <td class="school-student-redirect" hx-get="/%s/search/student/%s" hx-swap="outerHTML" hx-target="#school-page" hx-select="#students-page" hx-push-url="true">%s</td>
                    <td>%s</td>
					<td>%d</td>
					<td class="school-student-rank">%d</td>`,
				rowClass,
				requestType, student.RollNo, student.RollNo,
				requestType, student.RollNo, student.Name,
				student.Group,
				student.TotalMarks,
				student.Rank,
			))
			resultClass := "fail"
			resultText := "Fail"
			if student.Pass {
				resultClass = "pass"
				resultText = "Pass"
			}
			rankTableRows.WriteString(fmt.Sprintf(`<td class="%s">%s</td></tr>`, resultClass, resultText))
		}
	}
	temp.RankTableRows = template.HTML(rankTableRows.String())

	c.Set("Content-Type", "text/html")
	if err := tmpl.Execute(c.Response().BodyWriter(), temp); err != nil {
		log.Println("Error executing template:", err)
		return c.Status(500).SendString("Error rendering template")
	}
	return nil
}

// func GetSchoolResults(c *fiber.Ctx, requestType, schoolCode string) error {

// 	student, err := utils.ReadFile[types.SchoolResult]("./public/" + requestType + "/schools/" + schoolCode + ".json")
// 	if err != nil {
// 		log.Println("Error Reading student data:", err)
// 		return err
// 	}

// 	var studentTemplateMapping types.TemplateMappedStudentData

// 	tmpl, err := template.ParseFiles("./public/pages/students.html")
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
