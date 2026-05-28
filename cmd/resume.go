package cmd

import (
	"github.com/rioliu/liepin-cli-go/internal/models"
	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Query or update resume.",
	Long:  "Get resume, update base info and self-assessment, manage education/work/project experience and job preferences.",
}

// --- resume get ---

var resumeGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get current resume.",
	Run: func(_ *cobra.Command, _ []string) {
		executeGet("/mcp/get-resume")
	},
}

// --- update-base-info ---

var (
	realName, sex, birthday, cityCode, startJob, startJobMonth string
	nowWorkStatus, nowSalarySecret, nameSecret, wechat         string
	politicalStatusCode, nowIndusCode, nowJobTitleCode         string
	nowJobName, nowComp                                        string
	nowSalary, nowMonths                                       string
)

var resumeUpdateBaseInfoCmd = &cobra.Command{
	Use:   "update-base-info",
	Short: "Update base info.",
	Run: func(_ *cobra.Command, _ []string) {
		overrides := map[string]any{}
		if realName != "" {
			overrides["realName"] = realName
		}
		if sex != "" {
			overrides["sex"] = sex
		}
		if birthday != "" {
			overrides["birthday"] = birthday
		}
		if cityCode != "" {
			overrides["cityCode"] = cityCode
		}
		if startJob != "" {
			overrides["startJob"] = startJob
		}
		if startJobMonth != "" {
			overrides["startJobMonth"] = startJobMonth
		}
		if nowWorkStatus != "" {
			overrides["nowWorkStatus"] = nowWorkStatus
		}
		nowSalaryVal, err := models.ParseOptionalInt(nowSalary)
		if err != nil {
			handleError(err)
			return
		}
		if nowSalaryVal != nil {
			overrides["nowSalary"] = *nowSalaryVal
		}
		nowMonthsVal, err := models.ParseOptionalInt(nowMonths)
		if err != nil {
			handleError(err)
			return
		}
		if nowMonthsVal != nil {
			overrides["nowMonths"] = *nowMonthsVal
		}
		if nowSalarySecret != "" {
			overrides["nowSalarySecret"] = nowSalarySecret
		}
		if nowJobName != "" {
			overrides["jobName"] = nowJobName
		}
		if nowComp != "" {
			overrides["nowComp"] = nowComp
		}
		if nowIndusCode != "" {
			overrides["nowIndusCode"] = nowIndusCode
		}
		if nowJobTitleCode != "" {
			overrides["nowJobTitleCode"] = nowJobTitleCode
		}
		if nameSecret != "" {
			overrides["nameSecret"] = nameSecret
		}
		if wechat != "" {
			overrides["wechat"] = wechat
		}
		if politicalStatusCode != "" {
			overrides["politicalStatusCode"] = politicalStatusCode
		}

		m := models.UpdateBaseInfoInput{
			RealName:            models.StrPtr(realName),
			Sex:                 models.StrPtr(sex),
			Birthday:            models.StrPtr(birthday),
			CityCode:            models.StrPtr(cityCode),
			StartJob:            models.StrPtr(startJob),
			StartJobMonth:       models.StrPtr(startJobMonth),
			NowWorkStatus:       models.StrPtr(nowWorkStatus),
			NowSalarySecret:     models.StrPtr(nowSalarySecret),
			NameSecret:          models.StrPtr(nameSecret),
			Wechat:              models.StrPtr(wechat),
			PoliticalStatusCode: models.StrPtr(politicalStatusCode),
			NowIndusCode:        models.StrPtr(nowIndusCode),
			NowJobTitleCode:     models.StrPtr(nowJobTitleCode),
			JobName:             models.StrPtr(nowJobName),
			NowComp:             models.StrPtr(nowComp),
		}
		m.NowSalary = nowSalaryVal
		m.NowMonths = nowMonthsVal

		payload, err := buildPayload(overrides, m.Validate)
		if err != nil {
			handleError(err)
			return
		}
		executePost("/mcp/update-base-info", payload)
	},
}

// --- update-self-assess ---

var selfAssess string

var resumeUpdateSelfAssessCmd = &cobra.Command{
	Use:   "update-self-assess",
	Short: "Update self assessment.",
	Run: func(_ *cobra.Command, _ []string) {
		overrides := map[string]any{}
		if selfAssess != "" {
			overrides["selfAssess"] = selfAssess
		}

		m := models.UpdateSelfAssessInput{
			SelfAssess: selfAssess,
		}

		payload, err := buildPayload(overrides, m.Validate)
		if err != nil {
			handleError(err)
			return
		}
		executePost("/mcp/update-self-assess", payload)
	},
}

// --- add-edu-exp ---

var (
	eduSchool, eduMajor, eduStart, eduEnd, eduDegree, eduTz, eduExperience string
)

var resumeAddEduExpCmd = &cobra.Command{
	Use:   "add-edu-exp",
	Short: "Add education experience.",
	Run: func(_ *cobra.Command, _ []string) {
		overrides := map[string]any{}
		if eduSchool != "" {
			overrides["school"] = eduSchool
		}
		if eduMajor != "" {
			overrides["major"] = eduMajor
		}
		if eduStart != "" {
			overrides["start"] = eduStart
		}
		if eduEnd != "" {
			overrides["end"] = eduEnd
		}
		if eduDegree != "" {
			overrides["degree"] = eduDegree
		}
		if eduTz != "" {
			overrides["tz"] = eduTz
		}
		if eduExperience != "" {
			overrides["experience"] = eduExperience
		}

		m := models.AddEduExpInput{
			School:     eduSchool,
			Major:      models.StrPtr(eduMajor),
			Start:      eduStart,
			End:        eduEnd,
			Degree:     eduDegree,
			Tz:         models.StrPtr(eduTz),
			Experience: models.StrPtr(eduExperience),
		}

		payload, err := buildPayload(overrides, m.Validate)
		if err != nil {
			handleError(err)
			return
		}
		executePost("/mcp/add-edu-exp", payload)
	},
}

// --- update-edu-exp ---

var eduID string

var resumeUpdateEduExpCmd = &cobra.Command{
	Use:   "update-edu-exp",
	Short: "Update education experience.",
	Run: func(_ *cobra.Command, _ []string) {
		overrides := map[string]any{}
		idVal, err := models.ParseOptionalInt(eduID)
		if err != nil {
			handleError(err)
			return
		}
		if idVal != nil {
			overrides["eduId"] = *idVal
		}
		if eduSchool != "" {
			overrides["school"] = eduSchool
		}
		if eduMajor != "" {
			overrides["major"] = eduMajor
		}
		if eduStart != "" {
			overrides["start"] = eduStart
		}
		if eduEnd != "" {
			overrides["end"] = eduEnd
		}
		if eduDegree != "" {
			overrides["degree"] = eduDegree
		}
		if eduTz != "" {
			overrides["tz"] = eduTz
		}
		if eduExperience != "" {
			overrides["experience"] = eduExperience
		}

		m := models.UpdateEduExpInput{
			EduID:      optionalIntValue(idVal),
			School:     models.StrPtr(eduSchool),
			Major:      models.StrPtr(eduMajor),
			Start:      models.StrPtr(eduStart),
			End:        models.StrPtr(eduEnd),
			Degree:     models.StrPtr(eduDegree),
			Tz:         models.StrPtr(eduTz),
			Experience: models.StrPtr(eduExperience),
		}

		payload, err := buildPayload(overrides, m.Validate)
		if err != nil {
			handleError(err)
			return
		}
		executePost("/mcp/update-edu-exp", payload)
	},
}

// --- add-work-exp ---

var (
	workCompName, workIndustry, workStart, workEnd, workRwTitle, workJobtitle string
	workDq, workDept, workReport, workDuty, workCompkind, workCompscale       string
	workLabels                                                                string
	workSubordinate, workMonthsStr, workSalaryStr, workTypeStr                string
	shieldComp                                                                bool
)

var resumeAddWorkExpCmd = &cobra.Command{
	Use:   "add-work-exp",
	Short: "Add work experience.",
	Run: func(cmd *cobra.Command, _ []string) {
		overrides := map[string]any{}
		if workCompName != "" {
			overrides["compName"] = workCompName
		}
		if workIndustry != "" {
			overrides["industry"] = workIndustry
		}
		if workStart != "" {
			overrides["workStart"] = workStart
		}
		if workEnd != "" {
			overrides["workEnd"] = workEnd
		}
		if workRwTitle != "" {
			overrides["rwTitle"] = workRwTitle
		}
		if workJobtitle != "" {
			overrides["jobtitle"] = workJobtitle
		}
		if workDq != "" {
			overrides["dq"] = workDq
		}
		if workDept != "" {
			overrides["dept"] = workDept
		}
		if workReport != "" {
			overrides["report"] = workReport
		}
		if workDuty != "" {
			overrides["duty"] = workDuty
		}
		if workCompkind != "" {
			overrides["compkind"] = workCompkind
		}
		if workCompscale != "" {
			overrides["compscale"] = workCompscale
		}
		if workLabels != "" {
			overrides["labels"] = workLabels
		}

		sub, err := models.ParseOptionalInt(workSubordinate)
		if err != nil {
			handleError(err)
			return
		}
		if sub != nil {
			overrides["subordinate"] = *sub
		}
		months, err := models.ParseOptionalInt(workMonthsStr)
		if err != nil {
			handleError(err)
			return
		}
		if months != nil {
			overrides["months"] = *months
		}
		salary, err := models.ParseOptionalInt(workSalaryStr)
		if err != nil {
			handleError(err)
			return
		}
		if salary != nil {
			overrides["salary"] = *salary
		}
		wt, err := models.ParseOptionalInt(workTypeStr)
		if err != nil {
			handleError(err)
			return
		}
		if wt != nil {
			overrides["workType"] = *wt
		}
		if cmd.Flags().Changed("shield-comp") {
			overrides["shieldComp"] = shieldComp
		}

		m := models.AddWorkExpInput{
			CompName:    workCompName,
			Industry:    models.StrPtr(workIndustry),
			WorkStart:   workStart,
			WorkEnd:     workEnd,
			RwTitle:     workRwTitle,
			Jobtitle:    models.StrPtr(workJobtitle),
			Dq:          models.StrPtr(workDq),
			Dept:        models.StrPtr(workDept),
			Report:      models.StrPtr(workReport),
			Duty:        models.StrPtr(workDuty),
			Subordinate: sub,
			Months:      months,
			Salary:      salary,
			Compkind:    models.StrPtr(workCompkind),
			Compscale:   models.StrPtr(workCompscale),
			Labels:      models.StrPtr(workLabels),
			WorkType:    wt,
			ShieldComp:  models.ParseOptionalBool(cmd.Flags().Changed("shield-comp"), shieldComp),
		}

		payload, err := buildPayload(overrides, m.Validate)
		if err != nil {
			handleError(err)
			return
		}
		executePost("/mcp/add-work-exp", payload)
	},
}

// --- update-work-exp ---

var workID string

var resumeUpdateWorkExpCmd = &cobra.Command{
	Use:   "update-work-exp",
	Short: "Update work experience.",
	Run: func(cmd *cobra.Command, _ []string) {
		overrides := map[string]any{}
		idVal, err := models.ParseOptionalInt(workID)
		if err != nil {
			handleError(err)
			return
		}
		if idVal != nil {
			overrides["workId"] = *idVal
		}
		if workCompName != "" {
			overrides["compName"] = workCompName
		}
		if workIndustry != "" {
			overrides["industry"] = workIndustry
		}
		if workStart != "" {
			overrides["workStart"] = workStart
		}
		if workEnd != "" {
			overrides["workEnd"] = workEnd
		}
		if workRwTitle != "" {
			overrides["rwTitle"] = workRwTitle
		}
		if workJobtitle != "" {
			overrides["jobtitle"] = workJobtitle
		}
		if workDq != "" {
			overrides["dq"] = workDq
		}
		if workDept != "" {
			overrides["dept"] = workDept
		}
		if workReport != "" {
			overrides["report"] = workReport
		}
		if workDuty != "" {
			overrides["duty"] = workDuty
		}
		if workCompkind != "" {
			overrides["compkind"] = workCompkind
		}
		if workCompscale != "" {
			overrides["compscale"] = workCompscale
		}
		if workLabels != "" {
			overrides["labels"] = workLabels
		}

		sub, err := models.ParseOptionalInt(workSubordinate)
		if err != nil {
			handleError(err)
			return
		}
		if sub != nil {
			overrides["subordinate"] = *sub
		}
		months, err := models.ParseOptionalInt(workMonthsStr)
		if err != nil {
			handleError(err)
			return
		}
		if months != nil {
			overrides["months"] = *months
		}
		salary, err := models.ParseOptionalInt(workSalaryStr)
		if err != nil {
			handleError(err)
			return
		}
		if salary != nil {
			overrides["salary"] = *salary
		}
		wt, err := models.ParseOptionalInt(workTypeStr)
		if err != nil {
			handleError(err)
			return
		}
		if wt != nil {
			overrides["workType"] = *wt
		}
		if cmd.Flags().Changed("shield-comp") {
			overrides["shieldComp"] = shieldComp
		}

		m := models.UpdateWorkExpInput{
			WorkID:      optionalIntValue(idVal),
			CompName:    models.StrPtr(workCompName),
			Industry:    models.StrPtr(workIndustry),
			WorkStart:   models.StrPtr(workStart),
			WorkEnd:     models.StrPtr(workEnd),
			RwTitle:     models.StrPtr(workRwTitle),
			Jobtitle:    models.StrPtr(workJobtitle),
			Dq:          models.StrPtr(workDq),
			Dept:        models.StrPtr(workDept),
			Report:      models.StrPtr(workReport),
			Duty:        models.StrPtr(workDuty),
			Subordinate: sub,
			Months:      months,
			Salary:      salary,
			Compkind:    models.StrPtr(workCompkind),
			Compscale:   models.StrPtr(workCompscale),
			Labels:      models.StrPtr(workLabels),
			WorkType:    wt,
			ShieldComp:  models.ParseOptionalBool(cmd.Flags().Changed("shield-comp"), shieldComp),
		}

		payload, err := buildPayload(overrides, m.Validate)
		if err != nil {
			handleError(err)
			return
		}
		executePost("/mcp/update-work-exp", payload)
	},
}

// --- add-project-exp ---

var (
	projName, projStart, projEnd, projCompName, projPosition string
	projDescr, projDuty, projAchievement                     string
)

var resumeAddProjectExpCmd = &cobra.Command{
	Use:   "add-project-exp",
	Short: "Add project experience.",
	Run: func(_ *cobra.Command, _ []string) {
		overrides := map[string]any{}
		if projName != "" {
			overrides["name"] = projName
		}
		if projStart != "" {
			overrides["start"] = projStart
		}
		if projEnd != "" {
			overrides["end"] = projEnd
		}
		if projCompName != "" {
			overrides["compName"] = projCompName
		}
		if projPosition != "" {
			overrides["position"] = projPosition
		}
		if projDescr != "" {
			overrides["descr"] = projDescr
		}
		if projDuty != "" {
			overrides["duty"] = projDuty
		}
		if projAchievement != "" {
			overrides["achievement"] = projAchievement
		}

		m := models.AddProjectExpInput{
			Name:        projName,
			Start:       projStart,
			End:         projEnd,
			CompName:    models.StrPtr(projCompName),
			Position:    models.StrPtr(projPosition),
			Descr:       models.StrPtr(projDescr),
			Duty:        models.StrPtr(projDuty),
			Achievement: models.StrPtr(projAchievement),
		}

		payload, err := buildPayload(overrides, m.Validate)
		if err != nil {
			handleError(err)
			return
		}
		executePost("/mcp/add-project-exp", payload)
	},
}

// --- update-project-exp ---

var projID string

var resumeUpdateProjectExpCmd = &cobra.Command{
	Use:   "update-project-exp",
	Short: "Update project experience.",
	Run: func(_ *cobra.Command, _ []string) {
		overrides := map[string]any{}
		idVal, err := models.ParseOptionalInt(projID)
		if err != nil {
			handleError(err)
			return
		}
		if idVal != nil {
			overrides["id"] = *idVal
		}
		if projName != "" {
			overrides["name"] = projName
		}
		if projStart != "" {
			overrides["start"] = projStart
		}
		if projEnd != "" {
			overrides["end"] = projEnd
		}
		if projCompName != "" {
			overrides["compName"] = projCompName
		}
		if projPosition != "" {
			overrides["position"] = projPosition
		}
		if projDescr != "" {
			overrides["descr"] = projDescr
		}
		if projDuty != "" {
			overrides["duty"] = projDuty
		}
		if projAchievement != "" {
			overrides["achievement"] = projAchievement
		}

		m := models.UpdateProjectExpInput{
			ID:          optionalIntValue(idVal),
			Name:        models.StrPtr(projName),
			Start:       models.StrPtr(projStart),
			End:         models.StrPtr(projEnd),
			CompName:    models.StrPtr(projCompName),
			Position:    models.StrPtr(projPosition),
			Descr:       models.StrPtr(projDescr),
			Duty:        models.StrPtr(projDuty),
			Achievement: models.StrPtr(projAchievement),
		}

		payload, err := buildPayload(overrides, m.Validate)
		if err != nil {
			handleError(err)
			return
		}
		executePost("/mcp/update-project-exp", payload)
	},
}

// --- add-job-want ---

var (
	wantJobtitle, wantDq, wantWorkType                                        string
	wantIndustries, wantOtherExpectDqs                                        []string
	wantSalaryLowStr, wantSalaryHighStr, wantSalaryMonthsStr, wantWorkweekStr string
	wantPracticeMonthsStr                                                     string
)

var resumeAddJobWantCmd = &cobra.Command{
	Use:   "add-job-want",
	Short: "Add job preference.",
	Run: func(_ *cobra.Command, _ []string) {
		overrides := map[string]any{}
		if wantJobtitle != "" {
			overrides["jobtitle"] = wantJobtitle
		}
		if wantDq != "" {
			overrides["dq"] = wantDq
		}
		if len(wantIndustries) > 0 {
			overrides["industries"] = wantIndustries
		}
		if len(wantOtherExpectDqs) > 0 {
			overrides["otherExpectDqs"] = wantOtherExpectDqs
		}
		if wantWorkType != "" {
			overrides["workType"] = wantWorkType
		}

		low, err := models.ParseOptionalInt(wantSalaryLowStr)
		if err != nil {
			handleError(err)
			return
		}
		if low != nil {
			overrides["wantSalaryLow"] = low
		}
		high, err := models.ParseOptionalInt(wantSalaryHighStr)
		if err != nil {
			handleError(err)
			return
		}
		if high != nil {
			overrides["wantSalaryHigh"] = high
		}
		months, err := models.ParseOptionalInt(wantSalaryMonthsStr)
		if err != nil {
			handleError(err)
			return
		}
		if months != nil {
			overrides["wantSalaryMonths"] = months
		}
		week, err := models.ParseOptionalInt(wantWorkweekStr)
		if err != nil {
			handleError(err)
			return
		}
		if week != nil {
			overrides["workweek"] = week
		}
		practice, err := models.ParseOptionalInt(wantPracticeMonthsStr)
		if err != nil {
			handleError(err)
			return
		}
		if practice != nil {
			overrides["practiceMonths"] = practice
		}

		m := models.AddJobWantInput{
			Jobtitle:         wantJobtitle,
			Dq:               wantDq,
			Industries:       wantIndustries,
			WorkType:         models.StrPtr(wantWorkType),
			OtherExpectDqs:   wantOtherExpectDqs,
			WantSalaryLow:    low,
			WantSalaryHigh:   high,
			WantSalaryMonths: months,
			Workweek:         week,
			PracticeMonths:   practice,
		}

		payload, err := buildPayload(overrides, m.Validate)
		if err != nil {
			handleError(err)
			return
		}
		executePost("/mcp/add-job-want", payload)
	},
}

// --- update-job-want ---

var wantID string

var resumeUpdateJobWantCmd = &cobra.Command{
	Use:   "update-job-want",
	Short: "Update job preference.",
	Run: func(_ *cobra.Command, _ []string) {
		overrides := map[string]any{}
		idVal, err := models.ParseOptionalInt(wantID)
		if err != nil {
			handleError(err)
			return
		}
		if idVal != nil {
			overrides["id"] = *idVal
		}
		if wantJobtitle != "" {
			overrides["jobtitle"] = wantJobtitle
		}
		if wantDq != "" {
			overrides["dq"] = wantDq
		}
		if len(wantIndustries) > 0 {
			overrides["industries"] = wantIndustries
		}
		if len(wantOtherExpectDqs) > 0 {
			overrides["otherExpectDqs"] = wantOtherExpectDqs
		}
		if wantWorkType != "" {
			overrides["workType"] = wantWorkType
		}

		low, err := models.ParseOptionalInt(wantSalaryLowStr)
		if err != nil {
			handleError(err)
			return
		}
		if low != nil {
			overrides["wantSalaryLow"] = low
		}
		high, err := models.ParseOptionalInt(wantSalaryHighStr)
		if err != nil {
			handleError(err)
			return
		}
		if high != nil {
			overrides["wantSalaryHigh"] = high
		}
		months, err := models.ParseOptionalInt(wantSalaryMonthsStr)
		if err != nil {
			handleError(err)
			return
		}
		if months != nil {
			overrides["wantSalaryMonths"] = months
		}
		week, err := models.ParseOptionalInt(wantWorkweekStr)
		if err != nil {
			handleError(err)
			return
		}
		if week != nil {
			overrides["workweek"] = week
		}
		practice, err := models.ParseOptionalInt(wantPracticeMonthsStr)
		if err != nil {
			handleError(err)
			return
		}
		if practice != nil {
			overrides["practiceMonths"] = practice
		}

		m := models.UpdateJobWantInput{
			ID:               optionalIntValue(idVal),
			Jobtitle:         models.StrPtr(wantJobtitle),
			Dq:               models.StrPtr(wantDq),
			Industries:       wantIndustries,
			WorkType:         models.StrPtr(wantWorkType),
			OtherExpectDqs:   wantOtherExpectDqs,
			WantSalaryLow:    low,
			WantSalaryHigh:   high,
			WantSalaryMonths: months,
			Workweek:         week,
			PracticeMonths:   practice,
		}

		payload, err := buildPayload(overrides, m.Validate)
		if err != nil {
			handleError(err)
			return
		}
		executePost("/mcp/update-job-want", payload)
	},
}

func init() {
	resumeCmd.AddCommand(resumeGetCmd)
	resumeCmd.AddCommand(resumeUpdateBaseInfoCmd)
	resumeCmd.AddCommand(resumeUpdateSelfAssessCmd)
	resumeCmd.AddCommand(resumeAddEduExpCmd)
	resumeCmd.AddCommand(resumeUpdateEduExpCmd)
	resumeCmd.AddCommand(resumeAddWorkExpCmd)
	resumeCmd.AddCommand(resumeUpdateWorkExpCmd)
	resumeCmd.AddCommand(resumeAddProjectExpCmd)
	resumeCmd.AddCommand(resumeUpdateProjectExpCmd)
	resumeCmd.AddCommand(resumeAddJobWantCmd)
	resumeCmd.AddCommand(resumeUpdateJobWantCmd)

	for _, cmd := range []*cobra.Command{
		resumeGetCmd,
		resumeUpdateBaseInfoCmd, resumeUpdateSelfAssessCmd,
		resumeAddEduExpCmd, resumeUpdateEduExpCmd,
		resumeAddWorkExpCmd, resumeUpdateWorkExpCmd,
		resumeAddProjectExpCmd, resumeUpdateProjectExpCmd,
		resumeAddJobWantCmd, resumeUpdateJobWantCmd,
	} {
		addCommonFlags(cmd)
	}

	// update-base-info flags
	resumeUpdateBaseInfoCmd.Flags().StringVar(&realName, "real-name", "", "Name")
	resumeUpdateBaseInfoCmd.Flags().StringVar(&sex, "sex", "", "Gender")
	resumeUpdateBaseInfoCmd.Flags().StringVar(&birthday, "birthday", "", "Birthday, e.g. 19950101")
	resumeUpdateBaseInfoCmd.Flags().StringVar(&cityCode, "city-code", "", "City code")
	resumeUpdateBaseInfoCmd.Flags().StringVar(&startJob, "start-job", "", "First job year, e.g. 2020")
	resumeUpdateBaseInfoCmd.Flags().StringVar(&startJobMonth, "start-job-month", "", "First job month, e.g. 06")
	resumeUpdateBaseInfoCmd.Flags().StringVar(&nowWorkStatus, "now-work-status", "", "Current work status code")
	resumeUpdateBaseInfoCmd.Flags().StringVar(&nowSalary, "now-salary", "", "Current monthly salary")
	resumeUpdateBaseInfoCmd.Flags().StringVar(&nowMonths, "now-months", "", "Annual salary months")
	resumeUpdateBaseInfoCmd.Flags().StringVar(&nowSalarySecret, "now-salary-secret", "", "Hide current salary")
	resumeUpdateBaseInfoCmd.Flags().StringVar(&nowJobName, "job-name", "", "Current job title")
	resumeUpdateBaseInfoCmd.Flags().StringVar(&nowComp, "now-comp", "", "Current company name")
	resumeUpdateBaseInfoCmd.Flags().StringVar(&nowIndusCode, "now-indus-code", "", "Current industry code")
	resumeUpdateBaseInfoCmd.Flags().StringVar(&nowJobTitleCode, "now-job-title-code", "", "Current function code")
	resumeUpdateBaseInfoCmd.Flags().StringVar(&nameSecret, "name-secret", "", "Hide name")
	resumeUpdateBaseInfoCmd.Flags().StringVar(&wechat, "wechat", "", "WeChat ID")
	resumeUpdateBaseInfoCmd.Flags().StringVar(&politicalStatusCode, "political-status-code", "", "Political status code")

	// update-self-assess flags
	resumeUpdateSelfAssessCmd.Flags().StringVar(&selfAssess, "self-assess", "", "Self assessment content")

	// add-edu-exp flags
	resumeAddEduExpCmd.Flags().StringVar(&eduSchool, "school", "", "School name")
	resumeAddEduExpCmd.Flags().StringVar(&eduMajor, "major", "", "Major name")
	resumeAddEduExpCmd.Flags().StringVar(&eduStart, "start", "", "Start date, e.g. 201909")
	resumeAddEduExpCmd.Flags().StringVar(&eduEnd, "end", "", "End date, e.g. 202306")
	resumeAddEduExpCmd.Flags().StringVar(&eduDegree, "degree", "", "Degree code")
	resumeAddEduExpCmd.Flags().StringVar(&eduTz, "tz", "", "Unified admission")
	resumeAddEduExpCmd.Flags().StringVar(&eduExperience, "experience", "", "Campus experience notes")

	// update-edu-exp flags
	resumeUpdateEduExpCmd.Flags().StringVar(&eduID, "edu-id", "", "Education experience ID")
	resumeUpdateEduExpCmd.Flags().StringVar(&eduSchool, "school", "", "School name")
	resumeUpdateEduExpCmd.Flags().StringVar(&eduMajor, "major", "", "Major name")
	resumeUpdateEduExpCmd.Flags().StringVar(&eduStart, "start", "", "Start date")
	resumeUpdateEduExpCmd.Flags().StringVar(&eduEnd, "end", "", "End date")
	resumeUpdateEduExpCmd.Flags().StringVar(&eduDegree, "degree", "", "Degree code")
	resumeUpdateEduExpCmd.Flags().StringVar(&eduTz, "tz", "", "Unified admission")
	resumeUpdateEduExpCmd.Flags().StringVar(&eduExperience, "experience", "", "Campus experience notes")

	// add-work-exp flags
	resumeAddWorkExpCmd.Flags().StringVar(&workCompName, "comp-name", "", "Company name")
	resumeAddWorkExpCmd.Flags().StringVar(&workIndustry, "industry", "", "Industry")
	resumeAddWorkExpCmd.Flags().StringVar(&workStart, "work-start", "", "Start date, e.g. 202301")
	resumeAddWorkExpCmd.Flags().StringVar(&workEnd, "work-end", "", "End date, e.g. 202402")
	resumeAddWorkExpCmd.Flags().StringVar(&workRwTitle, "rw-title", "", "Job title")
	resumeAddWorkExpCmd.Flags().StringVar(&workJobtitle, "jobtitle", "", "Function name")
	resumeAddWorkExpCmd.Flags().StringVar(&workDq, "dq", "", "Work location")
	resumeAddWorkExpCmd.Flags().StringVar(&workDept, "dept", "", "Department")
	resumeAddWorkExpCmd.Flags().StringVar(&workReport, "report", "", "Reports to")
	resumeAddWorkExpCmd.Flags().StringVar(&workSubordinate, "subordinate", "", "Number of subordinates")
	resumeAddWorkExpCmd.Flags().StringVar(&workDuty, "duty", "", "Job responsibilities")
	resumeAddWorkExpCmd.Flags().StringVar(&workMonthsStr, "months", "", "Salary months")
	resumeAddWorkExpCmd.Flags().StringVar(&workSalaryStr, "salary", "", "Monthly salary")
	resumeAddWorkExpCmd.Flags().StringVar(&workCompkind, "compkind", "", "Company nature")
	resumeAddWorkExpCmd.Flags().StringVar(&workCompscale, "compscale", "", "Company size")
	resumeAddWorkExpCmd.Flags().StringVar(&workLabels, "labels", "", "Company labels")
	resumeAddWorkExpCmd.Flags().StringVar(&workTypeStr, "work-type", "", "Work type")
	resumeAddWorkExpCmd.Flags().BoolVar(&shieldComp, "shield-comp", false, "Shield company name")

	// update-work-exp flags
	resumeUpdateWorkExpCmd.Flags().StringVar(&workID, "work-id", "", "Work experience ID")
	resumeUpdateWorkExpCmd.Flags().StringVar(&workCompName, "comp-name", "", "Company name")
	resumeUpdateWorkExpCmd.Flags().StringVar(&workIndustry, "industry", "", "Industry")
	resumeUpdateWorkExpCmd.Flags().StringVar(&workStart, "work-start", "", "Start date")
	resumeUpdateWorkExpCmd.Flags().StringVar(&workEnd, "work-end", "", "End date")
	resumeUpdateWorkExpCmd.Flags().StringVar(&workRwTitle, "rw-title", "", "Job title")
	resumeUpdateWorkExpCmd.Flags().StringVar(&workJobtitle, "jobtitle", "", "Function name")
	resumeUpdateWorkExpCmd.Flags().StringVar(&workDq, "dq", "", "Work location")
	resumeUpdateWorkExpCmd.Flags().StringVar(&workDept, "dept", "", "Department")
	resumeUpdateWorkExpCmd.Flags().StringVar(&workReport, "report", "", "Reports to")
	resumeUpdateWorkExpCmd.Flags().StringVar(&workSubordinate, "subordinate", "", "Number of subordinates")
	resumeUpdateWorkExpCmd.Flags().StringVar(&workDuty, "duty", "", "Job responsibilities")
	resumeUpdateWorkExpCmd.Flags().StringVar(&workMonthsStr, "months", "", "Salary months")
	resumeUpdateWorkExpCmd.Flags().StringVar(&workSalaryStr, "salary", "", "Monthly salary")
	resumeUpdateWorkExpCmd.Flags().StringVar(&workCompkind, "compkind", "", "Company nature")
	resumeUpdateWorkExpCmd.Flags().StringVar(&workCompscale, "compscale", "", "Company size")
	resumeUpdateWorkExpCmd.Flags().StringVar(&workLabels, "labels", "", "Company labels")
	resumeUpdateWorkExpCmd.Flags().StringVar(&workTypeStr, "work-type", "", "Work type")
	resumeUpdateWorkExpCmd.Flags().BoolVar(&shieldComp, "shield-comp", false, "Shield company name")

	// add-project-exp flags
	resumeAddProjectExpCmd.Flags().StringVar(&projName, "name", "", "Project name")
	resumeAddProjectExpCmd.Flags().StringVar(&projStart, "start", "", "Start date, e.g. 202301")
	resumeAddProjectExpCmd.Flags().StringVar(&projEnd, "end", "", "End date, e.g. 202312")
	resumeAddProjectExpCmd.Flags().StringVar(&projCompName, "comp-name", "", "Company name")
	resumeAddProjectExpCmd.Flags().StringVar(&projPosition, "position", "", "Role")
	resumeAddProjectExpCmd.Flags().StringVar(&projDescr, "descr", "", "Project description")
	resumeAddProjectExpCmd.Flags().StringVar(&projDuty, "duty", "", "Responsibilities")
	resumeAddProjectExpCmd.Flags().StringVar(&projAchievement, "achievement", "", "Achievements")

	// update-project-exp flags
	resumeUpdateProjectExpCmd.Flags().StringVar(&projID, "id", "", "Project experience ID")
	resumeUpdateProjectExpCmd.Flags().StringVar(&projName, "name", "", "Project name")
	resumeUpdateProjectExpCmd.Flags().StringVar(&projStart, "start", "", "Start date")
	resumeUpdateProjectExpCmd.Flags().StringVar(&projEnd, "end", "", "End date")
	resumeUpdateProjectExpCmd.Flags().StringVar(&projCompName, "comp-name", "", "Company name")
	resumeUpdateProjectExpCmd.Flags().StringVar(&projPosition, "position", "", "Role")
	resumeUpdateProjectExpCmd.Flags().StringVar(&projDescr, "descr", "", "Project description")
	resumeUpdateProjectExpCmd.Flags().StringVar(&projDuty, "duty", "", "Responsibilities")
	resumeUpdateProjectExpCmd.Flags().StringVar(&projAchievement, "achievement", "", "Achievements")

	// add-job-want flags
	resumeAddJobWantCmd.Flags().StringVar(&wantJobtitle, "jobtitle", "", "Desired job title")
	resumeAddJobWantCmd.Flags().StringVar(&wantDq, "dq", "", "Desired work location")
	resumeAddJobWantCmd.Flags().StringArrayVar(&wantIndustries, "industries", nil, "Desired industries (repeatable)")
	resumeAddJobWantCmd.Flags().StringArrayVar(&wantOtherExpectDqs, "other-expect-dqs", nil, "Other desired locations (repeatable)")
	resumeAddJobWantCmd.Flags().StringVar(&wantSalaryLowStr, "want-salary-low", "", "Minimum desired salary")
	resumeAddJobWantCmd.Flags().StringVar(&wantSalaryHighStr, "want-salary-high", "", "Maximum desired salary")
	resumeAddJobWantCmd.Flags().StringVar(&wantSalaryMonthsStr, "want-salary-months", "", "Desired annual salary months")
	resumeAddJobWantCmd.Flags().StringVar(&wantWorkType, "work-type", "", "Work type")
	resumeAddJobWantCmd.Flags().StringVar(&wantWorkweekStr, "workweek", "", "Work days per week")
	resumeAddJobWantCmd.Flags().StringVar(&wantPracticeMonthsStr, "practice-months", "", "Internship months")

	// update-job-want flags
	resumeUpdateJobWantCmd.Flags().StringVar(&wantID, "id", "", "Job preference ID")
	resumeUpdateJobWantCmd.Flags().StringVar(&wantJobtitle, "jobtitle", "", "Desired job title")
	resumeUpdateJobWantCmd.Flags().StringVar(&wantDq, "dq", "", "Desired work location")
	resumeUpdateJobWantCmd.Flags().StringArrayVar(&wantIndustries, "industries", nil, "Desired industries (repeatable)")
	resumeUpdateJobWantCmd.Flags().StringArrayVar(&wantOtherExpectDqs, "other-expect-dqs", nil, "Other desired locations (repeatable)")
	resumeUpdateJobWantCmd.Flags().StringVar(&wantSalaryLowStr, "want-salary-low", "", "Minimum desired salary")
	resumeUpdateJobWantCmd.Flags().StringVar(&wantSalaryHighStr, "want-salary-high", "", "Maximum desired salary")
	resumeUpdateJobWantCmd.Flags().StringVar(&wantSalaryMonthsStr, "want-salary-months", "", "Desired annual salary months")
	resumeUpdateJobWantCmd.Flags().StringVar(&wantWorkType, "work-type", "", "Work type")
	resumeUpdateJobWantCmd.Flags().StringVar(&wantWorkweekStr, "workweek", "", "Work days per week")
	resumeUpdateJobWantCmd.Flags().StringVar(&wantPracticeMonthsStr, "practice-months", "", "Internship months")
}
