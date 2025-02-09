package services

import (
	"hse-results/utils"
	"strings"
	"fmt"
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
	formattedName:=strings.ToUpper(name)

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
