package types

import "html/template"

type Student struct {
	RollNo     string    `json:"rollNo"`
	Regular    bool      `json:"regular"`
	Name       string    `json:"name"`
	Group      string    `json:"group"`
	Subjects   []Subject `json:"subjects"`
	Result     string    `json:"result"`
	SchoolName string    `json:"schoolName"`
	SchoolCode string    `json:"schoolCode"`
	TotalMarks int16     `json:"totalMarks"`
	FullAp     bool      `json:"fullAplus"`
}
type Subject struct {
	Name  string `json:"name"`
	Grade string `json:"grade"`
	Marks *int16 `json:"marks"`
}

type TemplateMappedStudentData struct {
	RollNo       string
	Regular      bool
	Name         string
	StudentGroup string
	Subjects     template.HTML
	SchoolName   string
	SchoolCode   string
	FullAp       bool
	TotalMarks   int16
	Result       string
	Xvg          string
}

type SchoolResults []SchoolResult

type SchoolResult struct {
	SchoolName           string     `json:"schoolName"`
	SchoolCode           string     `json:"schoolCode"`
	PhoneNo              string     `json:"phoneNo"`
	EmailId              string     `json:"emailId"`
	SchoolType           string     `json:"schoolType"`
	LocalBody            string     `json:"localBody"`
	EduDistrict          string     `json:"eduDistrict"`
	AssemblyConstituency string     `json:"assemblyConstituency"`
	SchoolGender         string     `json:"schoolGender"`
	SecondLanguage       string     `json:"secondLanguage"`
	CourseCode           string     `json:"courseCode"`
	District             string     `json:"district"`
	Results              []Student  `json:"results"`
	PassPercentage       float64    `json:"passPercentage"`
	TotalFullAp          int        `json:"totalFullAp"`
	TotalFullMarks       int        `json:"totalFullMarks"`
	FullApStudents       []string   `json:"fullApStudents"`
	FullMarkStudents     []string   `json:"fullMarkStudents"`
	RankList             []RankList `json:"rankList"`
}

type RankList struct {
	RollNo     string `json:"rollNo"`
	Name       string `json:"name"`
	Group      string `json:"group"`
	TotalMarks int16  `json:"totalMarks"`
	Rank       int16  `json:"rank"`
}

type TemplateMappedSchoolResult struct {
	SchoolName           string        `json:"schoolName"`
	SchoolCode           string        `json:"schoolCode"`
	PhoneNo              string        `json:"phoneNo"`
	EmailId              string        `json:"emailId"`
	SchoolType           string        `json:"schoolType"`
	LocalBody            string        `json:"localBody"`
	EduDistrict          string        `json:"eduDistrict"`
	AssemblyConstituency string        `json:"assemblyConstituency"`
	SchoolGender         string        `json:"schoolGender"`
	SecondLanguage       string        `json:"secondLanguage"`
	CourseCode           string        `json:"courseCode"`
	District             string        `json:"district"`
	Results              SchoolResults `json:"results"`
	PassPercentage       float64       `json:"passPercentage"`
	TotalFullAp          int           `json:"totalFullAp"`
	TotalFullMarks       int           `json:"totalFullMarks"`
	FullApStudents       []string      `json:"fullApStudents"`
	FullMarkStudents     []string      `json:"fullMarkStudents"`
	RankList             []RankList    `json:"rankList"`
}
