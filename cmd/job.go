package cmd

import (
	"fmt"
	"os"

	"github.com/rioliu/liepin-cli-go/internal/models"
	"github.com/rioliu/liepin-cli-go/internal/output"
	"github.com/rioliu/liepin-cli-go/internal/scraper"
	"github.com/spf13/cobra"
)

var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "Search and apply for jobs.",
	Long:  "Search and apply for jobs. Search first to obtain job-id and job-kind before applying.",
}

// --- job search ---

var (
	searchWorkExperience, searchEduLevel, searchCompNature, searchAddress string
	searchSalaryFloor, searchSalaryCap, searchSalaryKind                  string
	searchJobName, searchCompanyName                                      string
	searchPage                                                            string
)

var jobSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search jobs.",
	Run: func(_ *cobra.Command, _ []string) {
		overrides := map[string]any{}
		if searchWorkExperience != "" {
			overrides["workExperience"] = searchWorkExperience
		}
		if searchEduLevel != "" {
			overrides["eduLevel"] = searchEduLevel
		}
		if searchCompNature != "" {
			overrides["compNature"] = searchCompNature
		}
		if searchAddress != "" {
			overrides["address"] = searchAddress
		}
		if searchSalaryFloor != "" {
			overrides["salaryFloor"] = searchSalaryFloor
		}
		if searchSalaryCap != "" {
			overrides["salaryCap"] = searchSalaryCap
		}
		if searchSalaryKind != "" {
			overrides["salaryKind"] = searchSalaryKind
		}
		if searchJobName != "" {
			overrides["jobName"] = searchJobName
		}
		if searchCompanyName != "" {
			overrides["companyName"] = searchCompanyName
		}
		page, err := models.ParseOptionalInt(searchPage)
		if err != nil {
			handleError(err)
			return
		}
		if page != nil {
			overrides["page"] = page
		}

		m := models.SearchJobInput{
			WorkExperience: models.StrPtr(searchWorkExperience),
			EduLevel:       models.StrPtr(searchEduLevel),
			CompNature:     models.StrPtr(searchCompNature),
			Address:        models.StrPtr(searchAddress),
			SalaryFloor:    models.StrPtr(searchSalaryFloor),
			SalaryCap:      models.StrPtr(searchSalaryCap),
			SalaryKind:     models.StrPtr(searchSalaryKind),
			JobName:        models.StrPtr(searchJobName),
			CompanyName:    models.StrPtr(searchCompanyName),
			Page:           page,
		}

		payload, err := buildPayload(overrides, m.Validate)
		if err != nil {
			handleError(err)
			return
		}
		executePost("/mcp/search-job", payload)
	},
}

// --- job detail ---

var (
	detailJobID  string
	detailJobURL string

	// jobDetailBaseURL is the base URL for constructing job detail URLs
	// from --job-id. Overridable in tests.
	jobDetailBaseURL = "https://www.liepin.com"
)

func buildJobDetailURL() (string, error) {
	if detailJobURL != "" {
		return detailJobURL, nil
	}
	if detailJobID != "" {
		return fmt.Sprintf("%s/a/%s.shtml", jobDetailBaseURL, detailJobID), nil
	}
	return "", fmt.Errorf("either --job-id or --job-detail-url is required")
}

var jobDetailCmd = &cobra.Command{
	Use:   "detail",
	Short: "Get full job description.",
	Long:  "Get full job description by scraping the job detail page. Use --job-detail-url from search results, or --job-id for headhunter jobs.",
	Run: func(_ *cobra.Command, _ []string) {
		jobURL, err := buildJobDetailURL()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		data, err := scraper.FetchJobDetail(jobURL)
		if err != nil {
			handleError(err)
			return
		}
		output.Render(data, outputFlag)
	},
}

// --- job apply ---

var (
	applyJobID   string
	applyJobKind string
)

var jobApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply for a job.",
	Long:  "Apply for a job. --job-id and --job-kind are both required. Use the type value from search results for --job-kind.",
	Run: func(_ *cobra.Command, _ []string) {
		overrides := map[string]any{}
		idVal, err := models.ParseOptionalInt(applyJobID)
		if err != nil {
			handleError(err)
			return
		}
		if idVal != nil {
			overrides["jobId"] = *idVal
		}
		if applyJobKind != "" {
			overrides["jobKind"] = applyJobKind
		}

		m := models.ApplyJobInput{
			JobKind: applyJobKind,
		}
		if idVal != nil {
			m.JobID = *idVal
		}

		payload, err := buildPayload(overrides, m.Validate)
		if err != nil {
			handleError(err)
			return
		}
		executePost("/mcp/apply-job", payload)
	},
}

func init() {
	jobCmd.AddCommand(jobSearchCmd)
	jobCmd.AddCommand(jobDetailCmd)
	jobCmd.AddCommand(jobApplyCmd)

	addCommonFlags(jobSearchCmd)
	addCommonFlags(jobDetailCmd)
	addCommonFlags(jobApplyCmd)

	// search flags
	jobSearchCmd.Flags().StringVar(&searchWorkExperience, "work-experience", "", "Work experience requirement")
	jobSearchCmd.Flags().StringVar(&searchEduLevel, "edu-level", "", "Education level requirement")
	jobSearchCmd.Flags().StringVar(&searchCompNature, "comp-nature", "", "Company nature or type")
	jobSearchCmd.Flags().StringVar(&searchAddress, "address", "", "Work location")
	jobSearchCmd.Flags().StringVar(&searchSalaryFloor, "salary-floor", "", "Minimum salary")
	jobSearchCmd.Flags().StringVar(&searchSalaryCap, "salary-cap", "", "Maximum salary")
	jobSearchCmd.Flags().StringVar(&searchSalaryKind, "salary-kind", "", "Salary type")
	jobSearchCmd.Flags().StringVar(&searchJobName, "job-name", "", "Job keyword")
	jobSearchCmd.Flags().StringVar(&searchCompanyName, "company-name", "", "Company name")
	jobSearchCmd.Flags().StringVar(&searchPage, "page", "", "Result page number (first page is usually 0)")

	// detail flags
	jobDetailCmd.Flags().StringVar(&detailJobID, "job-id", "", "Job ID (for type 1 headhunter jobs)")
	jobDetailCmd.Flags().StringVar(&detailJobURL, "job-detail-url", "", "Full job detail URL from search results (overrides --job-id)")

	// apply flags
	jobApplyCmd.Flags().StringVar(&applyJobID, "job-id", "", "Job ID (required)")
	jobApplyCmd.Flags().StringVar(&applyJobKind, "job-kind", "", "Job type code, e.g. 1 or 2 (required; reuse from search results)")
}
