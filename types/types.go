package types

type Student struct {
	RollNo     string  `json:"rollNo"`
	Regular    bool    `json:"regular"`
	Name       string  `json:"name"`
	Group      string  `json:"group"`
	Subject1   Subject `json:"subject1"`
	Subject2   Subject `json:"subject2"`
	Subject3   Subject `json:"subject3"`
	Subject4   Subject `json:"subject4"`
	Subject5   Subject `json:"subject5"`
	Subject6   Subject `json:"subject6"`
	Result     string  `json:"result"`
	SchoolName string  `json:"schoolName"`
	SchoolCode string  `json:"schoolCode"`
	TotalMarks int16   `json:"totalMarks"`
}
type Subject struct {
	Name  string `json:"name"`
	Grade string `json:"grade"`
	Marks *int16 `json:"marks"`
}

type TemplateMappedStudentData struct {
	RollNo        string
	Regular       bool
	Name          string
	StudentGroup  string
	Subject1      string
	Subject2      string
	Subject3      string
	Subject4      string
	Subject5      string
	Subject6      string
	Subject1Mark  *int16
	Subject2Mark  *int16
	Subject3Mark  *int16
	Subject4Mark  *int16
	Subject5Mark  *int16
	Subject6Mark  *int16
	Subject1Grade string
	Subject2Grade string
	Subject3Grade string
	Subject4Grade string
	Subject5Grade string
	Subject6Grade string
	Hide1         string
	Hide2         string
	Hide3         string
	Hide4         string
	Hide5         string
	Hide6         string
	SchoolName    string
	SchoolCode    string
	FullAplus     bool
	TotalMarks    int16
	Result        string
	Xvg           string
}

type SchoolResults []SchoolResult

type SchoolResult struct {
	RollNo   string  `json:"rollNo"`
	Regular  bool    `json:"regular"`
	Name     string  `json:"name"`
	Group    string  `json:"group"`
	Subject1 Subject `json:"subject1"`
	Subject2 Subject `json:"subject2"`
	Subject3 Subject `json:"subject3"`
	Subject4 Subject `json:"subject4"`
	Subject5 Subject `json:"subject5"`
	Subject6 Subject `json:"subject6"`
	Result   string  `json:"result"`
}

type TemplateMappedSchoolResult struct {
	Name       string        `json:"name"`
	SchoolCode string        `json:"schoolCode"`
	Result     SchoolResults `json:"result"`
}
