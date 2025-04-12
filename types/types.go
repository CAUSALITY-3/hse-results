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
	FullAplus  bool      `json:"fullAplus"`
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
	// Subject1      string
	// Subject2      string
	// Subject3      string
	// Subject4      string
	// Subject5      string
	// Subject6      string
	// Subject1Mark  *int16
	// Subject2Mark  *int16
	// Subject3Mark  *int16
	// Subject4Mark  *int16
	// Subject5Mark  *int16
	// Subject6Mark  *int16
	// Subject1Grade string
	// Subject2Grade string
	// Subject3Grade string
	// Subject4Grade string
	// Subject5Grade string
	// Subject6Grade string
	// Hide1         string
	// Hide2         string
	// Hide3         string
	// Hide4         string
	// Hide5         string
	// Hide6         string
	Subjects   template.HTML
	SchoolName string
	SchoolCode string
	FullAplus  bool
	TotalMarks int16
	Result     string
	Xvg        string
}

type SchoolResults []SchoolResult

type SchoolResult struct {
	RollNo   string    `json:"rollNo"`
	Regular  bool      `json:"regular"`
	Name     string    `json:"name"`
	Group    string    `json:"group"`
	Subjects []Subject `json:"subjects"`
	Result   string    `json:"result"`
}

type TemplateMappedSchoolResult struct {
	Name       string        `json:"name"`
	SchoolCode string        `json:"schoolCode"`
	Result     SchoolResults `json:"result"`
}
