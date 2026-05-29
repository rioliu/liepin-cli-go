// Package models defines the input payload structs sent to the Liepin API
// and the shared validation helpers used to check user-supplied data
// before issuing a request.
package models

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	reYear4     = regexp.MustCompile(`^\d{4}$`)
	reMonth2    = regexp.MustCompile(`^(0[1-9]|1[0-2])$`)
	reYearMonth = regexp.MustCompile(`^\d{4}(0[1-9]|1[0-2])$`)
	reDateYmd   = regexp.MustCompile(`^\d{4}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])$`)
)

// ValidationError describes a single input-validation failure on a named
// field. It implements the error interface so individual failures can also
// be returned directly when convenient.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors is a collection of ValidationError values gathered while
// validating a single input struct. Its Error method joins all individual
// messages with "; " for display.
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	var msgs []string
	for _, err := range e {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

func validateRequired(field, value string, errors *[]ValidationError) {
	if value == "" {
		*errors = append(*errors, ValidationError{field, "is required"})
	}
}

func validateRequiredInt(field string, _ int, isSet bool, errors *[]ValidationError) {
	if !isSet {
		*errors = append(*errors, ValidationError{field, "is required"})
	}
}

func validateEnum(field, value string, allowed []string, errors *[]ValidationError) {
	if value == "" {
		return
	}
	for _, a := range allowed {
		if value == a {
			return
		}
	}
	*errors = append(*errors, ValidationError{field, fmt.Sprintf("must be one of: %s", strings.Join(allowed, ", "))})
}

func validatePattern(field, value string, re *regexp.Regexp, desc string, errors *[]ValidationError) {
	if value == "" {
		return
	}
	if !re.MatchString(value) {
		*errors = append(*errors, ValidationError{field, fmt.Sprintf("must match format %s", desc)})
	}
}

func validateNonNegative(field string, value int, isSet bool, errors *[]ValidationError) {
	if isSet && value < 0 {
		*errors = append(*errors, ValidationError{field, "must be non-negative"})
	}
}

// StrPtr returns a pointer to s, or nil if s is empty. It is useful for
// building optional fields in request payloads from CLI flag values.
func StrPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// SafeString gets string from pointer or empty.
func SafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// --- SearchJobInput ---

// SearchJobInput is the request body for the job-search endpoint. All
// fields are optional; nil pointers are omitted from the marshalled JSON.
type SearchJobInput struct {
	WorkExperience *string `json:"workExperience,omitempty"`
	EduLevel       *string `json:"eduLevel,omitempty"`
	CompNature     *string `json:"compNature,omitempty"`
	Address        *string `json:"address,omitempty"`
	SalaryFloor    *string `json:"salaryFloor,omitempty"`
	SalaryCap      *string `json:"salaryCap,omitempty"`
	SalaryKind     *string `json:"salaryKind,omitempty"`
	JobName        *string `json:"jobName,omitempty"`
	CompanyName    *string `json:"companyName,omitempty"`
	Page           *int    `json:"page,omitempty"`
}

// Validate checks SearchJobInput for invalid field values and returns a
// ValidationErrors collection when one or more checks fail.
func (m *SearchJobInput) Validate() error {
	var errs []ValidationError
	if m.Page != nil {
		validateNonNegative("page", *m.Page, true, &errs)
	}
	if len(errs) > 0 {
		return ValidationErrors(errs)
	}
	return nil
}

// --- ApplyJobInput ---

// ApplyJobInput is the request body for applying to a job. Both fields are
// required; JobKind must be either "1" or "2" per the API contract.
type ApplyJobInput struct {
	JobID   int    `json:"jobId"`
	JobKind string `json:"jobKind"`
}

// Validate checks ApplyJobInput for missing or out-of-range fields.
func (m *ApplyJobInput) Validate() error {
	var errs []ValidationError
	validateRequiredInt("jobId", m.JobID, true, &errs)
	validateRequired("jobKind", m.JobKind, &errs)
	validateEnum("jobKind", m.JobKind, []string{"1", "2"}, &errs)
	if len(errs) > 0 {
		return ValidationErrors(errs)
	}
	return nil
}

// --- UpdateBaseInfoInput ---

// UpdateBaseInfoInput is the request body for updating a candidate's basic
// resume information. All fields are optional patches; only non-nil fields
// are sent to the API.
type UpdateBaseInfoInput struct {
	RealName            *string `json:"realName,omitempty"`
	Sex                 *string `json:"sex,omitempty"`
	Birthday            *string `json:"birthday,omitempty"`
	CityCode            *string `json:"cityCode,omitempty"`
	StartJob            *string `json:"startJob,omitempty"`
	StartJobMonth       *string `json:"startJobMonth,omitempty"`
	NowWorkStatus       *string `json:"nowWorkStatus,omitempty"`
	NowSalary           *int    `json:"nowSalary,omitempty"`
	NowMonths           *int    `json:"nowMonths,omitempty"`
	NowSalarySecret     *string `json:"nowSalarySecret,omitempty"`
	JobName             *string `json:"jobName,omitempty"`
	NowComp             *string `json:"nowComp,omitempty"`
	NowIndusCode        *string `json:"nowIndusCode,omitempty"`
	NowJobTitleCode     *string `json:"nowJobTitleCode,omitempty"`
	NameSecret          *string `json:"nameSecret,omitempty"`
	Wechat              *string `json:"wechat,omitempty"`
	PoliticalStatusCode *string `json:"politicalStatusCode,omitempty"`
}

// Validate enforces enum, date-format, and pairing constraints on
// UpdateBaseInfoInput (for example startJob and startJobMonth must be set
// together).
func (m *UpdateBaseInfoInput) Validate() error {
	var errs []ValidationError
	hasStartJob := m.StartJob != nil
	hasStartJobMonth := m.StartJobMonth != nil
	if hasStartJob != hasStartJobMonth {
		errs = append(errs, ValidationError{"startJob/startJobMonth", "must be provided together"})
	}
	if m.Sex != nil {
		validateEnum("sex", *m.Sex, []string{"男", "女"}, &errs)
	}
	if m.Birthday != nil {
		validatePattern("birthday", *m.Birthday, reDateYmd, "yyyyMMdd (e.g. 19950101)", &errs)
	}
	if m.StartJob != nil {
		validatePattern("startJob", *m.StartJob, reYear4, "yyyy (e.g. 2020)", &errs)
	}
	if m.StartJobMonth != nil {
		validatePattern("startJobMonth", *m.StartJobMonth, reMonth2, "MM (e.g. 06)", &errs)
	}
	if m.NowWorkStatus != nil {
		validateEnum("nowWorkStatus", *m.NowWorkStatus, []string{"0", "1", "2", "3", "4", "5", "6", "7"}, &errs)
	}
	if m.NowSalary != nil {
		validateNonNegative("nowSalary", *m.NowSalary, true, &errs)
	}
	if m.NowMonths != nil {
		validateNonNegative("nowMonths", *m.NowMonths, true, &errs)
	}
	if m.NowSalarySecret != nil {
		validateEnum("nowSalarySecret", *m.NowSalarySecret, []string{"0", "1"}, &errs)
	}
	if m.NameSecret != nil {
		validateEnum("nameSecret", *m.NameSecret, []string{"0", "1"}, &errs)
	}
	if m.PoliticalStatusCode != nil {
		validateEnum("politicalStatusCode", *m.PoliticalStatusCode, []string{"1", "2", "3", "4", "5", "6"}, &errs)
	}
	if len(errs) > 0 {
		return ValidationErrors(errs)
	}
	return nil
}

// --- UpdateSelfAssessInput ---

// UpdateSelfAssessInput is the request body for updating the self-assessment
// section of a resume. SelfAssess is required.
type UpdateSelfAssessInput struct {
	SelfAssess string `json:"selfAssess"`
}

// Validate ensures SelfAssess is non-empty.
func (m *UpdateSelfAssessInput) Validate() error {
	var errs []ValidationError
	validateRequired("selfAssess", m.SelfAssess, &errs)
	if len(errs) > 0 {
		return ValidationErrors(errs)
	}
	return nil
}

// --- AddEduExpInput ---

// AddEduExpInput is the request body for adding an education experience.
// School, Start, End, and Degree are required; date fields use yyyyMM.
type AddEduExpInput struct {
	School     string  `json:"school"`
	Major      *string `json:"major,omitempty"`
	Start      string  `json:"start"`
	End        string  `json:"end"`
	Degree     string  `json:"degree"`
	Tz         *string `json:"tz,omitempty"`
	Experience *string `json:"experience,omitempty"`
}

// Validate enforces required fields, allowed degree codes, and the yyyyMM
// date format on AddEduExpInput.
func (m *AddEduExpInput) Validate() error {
	var errs []ValidationError
	validateRequired("school", m.School, &errs)
	validateRequired("start", m.Start, &errs)
	validateRequired("end", m.End, &errs)
	validateRequired("degree", m.Degree, &errs)
	validatePattern("start", m.Start, reYearMonth, "yyyyMM (e.g. 201909)", &errs)
	validatePattern("end", m.End, reYearMonth, "yyyyMM (e.g. 202306)", &errs)
	validateEnum("degree", m.Degree, []string{"090", "080", "060", "050", "040", "030", "020", "010"}, &errs)
	if m.Tz != nil {
		validateEnum("tz", *m.Tz, []string{"0", "1"}, &errs)
	}
	if len(errs) > 0 {
		return ValidationErrors(errs)
	}
	return nil
}

// --- UpdateEduExpInput ---

// UpdateEduExpInput is the request body for patching an existing education
// experience identified by EduID. All other fields are optional patches.
type UpdateEduExpInput struct {
	EduID      int     `json:"eduId"`
	School     *string `json:"school,omitempty"`
	Major      *string `json:"major,omitempty"`
	Start      *string `json:"start,omitempty"`
	End        *string `json:"end,omitempty"`
	Degree     *string `json:"degree,omitempty"`
	Tz         *string `json:"tz,omitempty"`
	Experience *string `json:"experience,omitempty"`
}

// Validate requires EduID and validates any supplied optional fields.
func (m *UpdateEduExpInput) Validate() error {
	var errs []ValidationError
	validateRequiredInt("eduId", m.EduID, true, &errs)
	if m.Start != nil {
		validatePattern("start", *m.Start, reYearMonth, "yyyyMM", &errs)
	}
	if m.End != nil {
		validatePattern("end", *m.End, reYearMonth, "yyyyMM", &errs)
	}
	if m.Degree != nil {
		validateEnum("degree", *m.Degree, []string{"090", "080", "060", "050", "040", "030", "020", "010"}, &errs)
	}
	if m.Tz != nil {
		validateEnum("tz", *m.Tz, []string{"0", "1"}, &errs)
	}
	if len(errs) > 0 {
		return ValidationErrors(errs)
	}
	return nil
}

// --- AddWorkExpInput ---

// AddWorkExpInput is the request body for adding a work experience entry.
// CompName, WorkStart, WorkEnd, and RwTitle are required.
type AddWorkExpInput struct {
	CompName    string  `json:"compName"`
	Industry    *string `json:"industry,omitempty"`
	WorkStart   string  `json:"workStart"`
	WorkEnd     string  `json:"workEnd"`
	RwTitle     string  `json:"rwTitle"`
	Jobtitle    *string `json:"jobtitle,omitempty"`
	Dq          *string `json:"dq,omitempty"`
	Dept        *string `json:"dept,omitempty"`
	Report      *string `json:"report,omitempty"`
	Subordinate *int    `json:"subordinate,omitempty"`
	Duty        *string `json:"duty,omitempty"`
	Months      *int    `json:"months,omitempty"`
	Salary      *int    `json:"salary,omitempty"`
	Compkind    *string `json:"compkind,omitempty"`
	Compscale   *string `json:"compscale,omitempty"`
	ShieldComp  *bool   `json:"shieldComp,omitempty"`
	Labels      *string `json:"labels,omitempty"`
	WorkType    *int    `json:"workType,omitempty"`
}

// Validate enforces required fields, yyyyMM date format, non-negative
// numeric fields, and an allowed WorkType (1 or 2) on AddWorkExpInput.
func (m *AddWorkExpInput) Validate() error {
	var errs []ValidationError
	validateRequired("compName", m.CompName, &errs)
	validateRequired("workStart", m.WorkStart, &errs)
	validateRequired("workEnd", m.WorkEnd, &errs)
	validateRequired("rwTitle", m.RwTitle, &errs)
	validatePattern("workStart", m.WorkStart, reYearMonth, "yyyyMM (e.g. 202301)", &errs)
	validatePattern("workEnd", m.WorkEnd, reYearMonth, "yyyyMM (e.g. 202402)", &errs)
	if m.Subordinate != nil {
		validateNonNegative("subordinate", *m.Subordinate, true, &errs)
	}
	if m.Months != nil {
		validateNonNegative("months", *m.Months, true, &errs)
	}
	if m.Salary != nil {
		validateNonNegative("salary", *m.Salary, true, &errs)
	}
	if m.WorkType != nil {
		if *m.WorkType != 1 && *m.WorkType != 2 {
			errs = append(errs, ValidationError{"workType", "must be 1 or 2"})
		}
	}
	if len(errs) > 0 {
		return ValidationErrors(errs)
	}
	return nil
}

// --- UpdateWorkExpInput ---

// UpdateWorkExpInput is the request body for patching an existing work
// experience identified by WorkID. All other fields are optional patches.
type UpdateWorkExpInput struct {
	WorkID      int     `json:"workId"`
	CompName    *string `json:"compName,omitempty"`
	Industry    *string `json:"industry,omitempty"`
	WorkStart   *string `json:"workStart,omitempty"`
	WorkEnd     *string `json:"workEnd,omitempty"`
	RwTitle     *string `json:"rwTitle,omitempty"`
	Jobtitle    *string `json:"jobtitle,omitempty"`
	Dq          *string `json:"dq,omitempty"`
	Dept        *string `json:"dept,omitempty"`
	Report      *string `json:"report,omitempty"`
	Subordinate *int    `json:"subordinate,omitempty"`
	Duty        *string `json:"duty,omitempty"`
	Months      *int    `json:"months,omitempty"`
	Salary      *int    `json:"salary,omitempty"`
	Compkind    *string `json:"compkind,omitempty"`
	Compscale   *string `json:"compscale,omitempty"`
	ShieldComp  *bool   `json:"shieldComp,omitempty"`
	Labels      *string `json:"labels,omitempty"`
	WorkType    *int    `json:"workType,omitempty"`
}

// Validate requires WorkID and validates any supplied optional fields,
// including yyyyMM date format, non-negative numerics, and WorkType (1 or 2).
func (m *UpdateWorkExpInput) Validate() error {
	var errs []ValidationError
	validateRequiredInt("workId", m.WorkID, true, &errs)
	if m.WorkStart != nil {
		validatePattern("workStart", *m.WorkStart, reYearMonth, "yyyyMM", &errs)
	}
	if m.WorkEnd != nil {
		validatePattern("workEnd", *m.WorkEnd, reYearMonth, "yyyyMM", &errs)
	}
	if m.Subordinate != nil {
		validateNonNegative("subordinate", *m.Subordinate, true, &errs)
	}
	if m.Months != nil {
		validateNonNegative("months", *m.Months, true, &errs)
	}
	if m.Salary != nil {
		validateNonNegative("salary", *m.Salary, true, &errs)
	}
	if m.WorkType != nil {
		if *m.WorkType != 1 && *m.WorkType != 2 {
			errs = append(errs, ValidationError{"workType", "must be 1 or 2"})
		}
	}
	if len(errs) > 0 {
		return ValidationErrors(errs)
	}
	return nil
}

// --- AddProjectExpInput ---

// AddProjectExpInput is the request body for adding a project experience
// entry. Name, Start, and End are required.
type AddProjectExpInput struct {
	Name        string  `json:"name"`
	Start       string  `json:"start"`
	End         string  `json:"end"`
	CompName    *string `json:"compName,omitempty"`
	Position    *string `json:"position,omitempty"`
	Descr       *string `json:"descr,omitempty"`
	Duty        *string `json:"duty,omitempty"`
	Achievement *string `json:"achievement,omitempty"`
}

// Validate enforces required fields and yyyyMM date format on
// AddProjectExpInput.
func (m *AddProjectExpInput) Validate() error {
	var errs []ValidationError
	validateRequired("name", m.Name, &errs)
	validateRequired("start", m.Start, &errs)
	validateRequired("end", m.End, &errs)
	validatePattern("start", m.Start, reYearMonth, "yyyyMM", &errs)
	validatePattern("end", m.End, reYearMonth, "yyyyMM", &errs)
	if len(errs) > 0 {
		return ValidationErrors(errs)
	}
	return nil
}

// --- UpdateProjectExpInput ---

// UpdateProjectExpInput is the request body for patching an existing
// project experience identified by ID. Other fields are optional patches.
type UpdateProjectExpInput struct {
	ID          int     `json:"id"`
	Name        *string `json:"name,omitempty"`
	Start       *string `json:"start,omitempty"`
	End         *string `json:"end,omitempty"`
	CompName    *string `json:"compName,omitempty"`
	Position    *string `json:"position,omitempty"`
	Descr       *string `json:"descr,omitempty"`
	Duty        *string `json:"duty,omitempty"`
	Achievement *string `json:"achievement,omitempty"`
}

// Validate requires ID and validates any supplied yyyyMM date fields.
func (m *UpdateProjectExpInput) Validate() error {
	var errs []ValidationError
	validateRequiredInt("id", m.ID, true, &errs)
	if m.Start != nil {
		validatePattern("start", *m.Start, reYearMonth, "yyyyMM", &errs)
	}
	if m.End != nil {
		validatePattern("end", *m.End, reYearMonth, "yyyyMM", &errs)
	}
	if len(errs) > 0 {
		return ValidationErrors(errs)
	}
	return nil
}

// --- AddJobWantInput ---

// AddJobWantInput is the request body for adding a desired-job preference.
// Jobtitle and Dq are required; salary and WorkType fields are validated
// when present.
type AddJobWantInput struct {
	Industries       []string `json:"industries,omitempty"`
	Jobtitle         string   `json:"jobtitle"`
	Dq               string   `json:"dq"`
	WantSalaryLow    *int     `json:"wantSalaryLow,omitempty"`
	WantSalaryHigh   *int     `json:"wantSalaryHigh,omitempty"`
	WantSalaryMonths *int     `json:"wantSalaryMonths,omitempty"`
	WorkType         *string  `json:"workType,omitempty"`
	OtherExpectDqs   []string `json:"otherExpectDqs,omitempty"`
	Workweek         *int     `json:"workweek,omitempty"`
	PracticeMonths   *int     `json:"practiceMonths,omitempty"`
}

// Validate enforces required fields and salary-range consistency
// (wantSalaryLow must not exceed wantSalaryHigh) on AddJobWantInput.
func (m *AddJobWantInput) Validate() error {
	var errs []ValidationError
	validateRequired("jobtitle", m.Jobtitle, &errs)
	validateRequired("dq", m.Dq, &errs)
	if m.WantSalaryLow != nil {
		validateNonNegative("wantSalaryLow", *m.WantSalaryLow, true, &errs)
	}
	if m.WantSalaryHigh != nil {
		validateNonNegative("wantSalaryHigh", *m.WantSalaryHigh, true, &errs)
	}
	if m.WantSalaryLow != nil && m.WantSalaryHigh != nil && *m.WantSalaryLow > *m.WantSalaryHigh {
		errs = append(errs, ValidationError{"wantSalaryLow/wantSalaryHigh", "wantSalaryLow must not exceed wantSalaryHigh"})
	}
	if m.WantSalaryMonths != nil {
		validateNonNegative("wantSalaryMonths", *m.WantSalaryMonths, true, &errs)
	}
	if m.WorkType != nil {
		validateEnum("workType", *m.WorkType, []string{"0", "1"}, &errs)
	}
	if m.PracticeMonths != nil {
		validateNonNegative("practiceMonths", *m.PracticeMonths, true, &errs)
	}
	if len(errs) > 0 {
		return ValidationErrors(errs)
	}
	return nil
}

// --- UpdateJobWantInput ---

// UpdateJobWantInput is the request body for patching an existing
// desired-job entry identified by ID. Remaining fields are optional patches.
type UpdateJobWantInput struct {
	ID               int      `json:"id"`
	Jobtitle         *string  `json:"jobtitle,omitempty"`
	Dq               *string  `json:"dq,omitempty"`
	Industries       []string `json:"industries,omitempty"`
	WantSalaryLow    *int     `json:"wantSalaryLow,omitempty"`
	WantSalaryHigh   *int     `json:"wantSalaryHigh,omitempty"`
	WantSalaryMonths *int     `json:"wantSalaryMonths,omitempty"`
	WorkType         *string  `json:"workType,omitempty"`
	OtherExpectDqs   []string `json:"otherExpectDqs,omitempty"`
	Workweek         *int     `json:"workweek,omitempty"`
	PracticeMonths   *int     `json:"practiceMonths,omitempty"`
}

// Validate requires ID and enforces the same salary-range and enum
// constraints as AddJobWantInput for any supplied optional fields.
func (m *UpdateJobWantInput) Validate() error {
	var errs []ValidationError
	validateRequiredInt("id", m.ID, true, &errs)
	if m.WantSalaryLow != nil {
		validateNonNegative("wantSalaryLow", *m.WantSalaryLow, true, &errs)
	}
	if m.WantSalaryHigh != nil {
		validateNonNegative("wantSalaryHigh", *m.WantSalaryHigh, true, &errs)
	}
	if m.WantSalaryLow != nil && m.WantSalaryHigh != nil && *m.WantSalaryLow > *m.WantSalaryHigh {
		errs = append(errs, ValidationError{"wantSalaryLow/wantSalaryHigh", "wantSalaryLow must not exceed wantSalaryHigh"})
	}
	if m.WantSalaryMonths != nil {
		validateNonNegative("wantSalaryMonths", *m.WantSalaryMonths, true, &errs)
	}
	if m.WorkType != nil {
		validateEnum("workType", *m.WorkType, []string{"0", "1"}, &errs)
	}
	if m.PracticeMonths != nil {
		validateNonNegative("practiceMonths", *m.PracticeMonths, true, &errs)
	}
	if len(errs) > 0 {
		return ValidationErrors(errs)
	}
	return nil
}

// --- Helpers for CLI parsing ---

// ParseOptionalInt converts a CLI string flag into an optional integer.
// An empty input returns (nil, nil) to mean "not provided"; a non-empty
// input that fails to parse returns a descriptive error.
func ParseOptionalInt(raw string) (*int, error) {
	if raw == "" {
		return nil, nil
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid number: %s", raw)
	}
	return &v, nil
}

// ParseOptionalBool returns a pointer to v when isSet is true (meaning the
// caller explicitly provided the flag), or nil otherwise. This lets request
// payloads distinguish "false" from "unset".
func ParseOptionalBool(isSet bool, v bool) *bool {
	if !isSet {
		return nil
	}
	return &v
}
