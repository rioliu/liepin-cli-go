package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/rioliu/liepin-cli-go/internal/authstore"
)

// ─── helpers ────────────────────────────────────────────────────────

func resetAllFlags() {
	tokenFlag, baseURLFlag, timeoutFlag, outputFlag, inputFlag = "", "", 0, "", ""

	searchWorkExperience, searchEduLevel, searchCompNature, searchAddress = "", "", "", ""
	searchSalaryFloor, searchSalaryCap, searchSalaryKind = "", "", ""
	searchJobName, searchCompanyName, searchPage = "", "", ""
	applyJobID, applyJobKind = "", ""

	realName, sex, birthday, cityCode, startJob, startJobMonth = "", "", "", "", "", ""
	nowWorkStatus, nowSalarySecret, nameSecret, wechat = "", "", "", ""
	politicalStatusCode, nowIndusCode, nowJobTitleCode = "", "", ""
	nowJobName, nowComp, nowSalary, nowMonths = "", "", "", ""
	selfAssess = ""
	eduSchool, eduMajor, eduStart, eduEnd, eduDegree, eduTz, eduExperience, eduID = "", "", "", "", "", "", "", ""
	workCompName, workIndustry, workStart, workEnd, workRwTitle, workJobtitle = "", "", "", "", "", ""
	workDq, workDept, workReport, workDuty, workCompkind, workCompscale = "", "", "", "", "", ""
	workLabels, workSubordinate, workMonthsStr, workSalaryStr, workTypeStr = "", "", "", "", ""
	shieldComp, workID = false, ""
	projName, projStart, projEnd, projCompName, projPosition = "", "", "", "", ""
	projDescr, projDuty, projAchievement, projID = "", "", "", ""
	wantJobtitle, wantDq, wantWorkType = "", "", ""
	wantIndustries, wantOtherExpectDqs = nil, nil
	wantSalaryLowStr, wantSalaryHighStr, wantSalaryMonthsStr, wantWorkweekStr = "", "", "", ""
	wantPracticeMonthsStr, wantID = "", ""

	os.Unsetenv("LIEPIN_USER_TOKEN")
}

func startMockServer(handler http.HandlerFunc) *httptest.Server {
	wrapped := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer GinkgoRecover()
		handler(w, r)
	})
	srv := httptest.NewServer(wrapped)
	DeferCleanup(srv.Close)
	baseURLFlag = srv.URL
	tokenFlag = "test-token"
	outputFlag = "json"
	return srv
}

func captureStdout(fn func()) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func jsonOK() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	}
}

func decodeBody(r *http.Request) map[string]any {
	var body map[string]any
	json.NewDecoder(r.Body).Decode(&body)
	return body
}

// ─── root / help structure ──────────────────────────────────────────

var _ = Describe("Root command", func() {
	BeforeEach(func() {
		resetAllFlags()
	})

	It("shows all subcommands in help", func() {
		out := captureStdout(func() {
			rootCmd.SetArgs([]string{"--help"})
			rootCmd.Execute()
		})
		Expect(out).To(ContainSubstring("setup"))
		Expect(out).To(ContainSubstring("auth"))
		Expect(out).To(ContainSubstring("resume"))
		Expect(out).To(ContainSubstring("job"))
	})

	It("hides common flags from help", func() {
		out := captureStdout(func() {
			rootCmd.SetArgs([]string{"--help"})
			rootCmd.Execute()
		})
		Expect(out).NotTo(ContainSubstring("--token"))
		Expect(out).NotTo(ContainSubstring("--base-url"))
		Expect(out).NotTo(ContainSubstring("--timeout"))
		Expect(out).NotTo(ContainSubstring("--input"))
	})

	It("has no Chinese characters in help", func() {
		out := captureStdout(func() {
			rootCmd.SetArgs([]string{"--help"})
			rootCmd.Execute()
		})
		for _, ch := range []string{"设置", "求职", "简历", "搜索", "登录"} {
			Expect(out).NotTo(ContainSubstring(ch))
		}
	})
})

var _ = Describe("Auth command", func() {
	BeforeEach(func() {
		resetAllFlags()
	})

	It("shows subcommands", func() {
		out := captureStdout(func() {
			rootCmd.SetArgs([]string{"auth", "--help"})
			rootCmd.Execute()
		})
		Expect(out).To(ContainSubstring("status"))
		Expect(out).To(ContainSubstring("clear"))
		Expect(out).To(ContainSubstring("open"))
		Expect(out).To(ContainSubstring("setup"))
	})
})

var _ = Describe("Resume command", func() {
	BeforeEach(func() {
		resetAllFlags()
	})

	It("shows all subcommands", func() {
		out := captureStdout(func() {
			rootCmd.SetArgs([]string{"resume", "--help"})
			rootCmd.Execute()
		})
		Expect(out).To(ContainSubstring("get"))
		Expect(out).To(ContainSubstring("update-base-info"))
		Expect(out).To(ContainSubstring("update-self-assess"))
		Expect(out).To(ContainSubstring("add-edu-exp"))
		Expect(out).To(ContainSubstring("update-edu-exp"))
		Expect(out).To(ContainSubstring("add-work-exp"))
		Expect(out).To(ContainSubstring("update-work-exp"))
		Expect(out).To(ContainSubstring("add-project-exp"))
		Expect(out).To(ContainSubstring("update-project-exp"))
		Expect(out).To(ContainSubstring("add-job-want"))
		Expect(out).To(ContainSubstring("update-job-want"))
	})
})

var _ = Describe("Job command", func() {
	BeforeEach(func() {
		resetAllFlags()
	})

	It("shows subcommands", func() {
		out := captureStdout(func() {
			rootCmd.SetArgs([]string{"job", "--help"})
			rootCmd.Execute()
		})
		Expect(out).To(ContainSubstring("search"))
		Expect(out).To(ContainSubstring("apply"))
	})
})

// ─── job search ─────────────────────────────────────────────────────

var _ = Describe("Job search", func() {
	BeforeEach(func() {
		resetAllFlags()
	})

	It("sends correct payload to /mcp/search-job", func() {
		startMockServer(func(w http.ResponseWriter, r *http.Request) {
			Expect(r.Method).To(Equal(http.MethodPost))
			Expect(r.URL.Path).To(Equal("/mcp/search-job"))
			Expect(r.Header.Get("x-user-token")).To(Equal("test-token"))

			body := decodeBody(r)
			Expect(body["jobName"]).To(Equal("Go Engineer"))
			Expect(body["address"]).To(Equal("Beijing"))
			Expect(body["page"]).To(Equal(float64(0)))

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"ok": true})
		})

		searchJobName = "Go Engineer"
		searchAddress = "Beijing"
		searchPage = "0"

		out := captureStdout(func() {
			jobSearchCmd.Run(jobSearchCmd, nil)
		})
		Expect(out).To(ContainSubstring(`"ok": true`))
	})

	It("sends salary filter fields", func() {
		startMockServer(func(w http.ResponseWriter, r *http.Request) {
			body := decodeBody(r)
			Expect(body["salaryFloor"]).To(Equal("10000"))
			Expect(body["salaryCap"]).To(Equal("30000"))
			Expect(body["salaryKind"]).To(Equal("1"))
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"ok": true})
		})

		searchSalaryFloor = "10000"
		searchSalaryCap = "30000"
		searchSalaryKind = "1"

		captureStdout(func() { jobSearchCmd.Run(jobSearchCmd, nil) })
	})

	It("omits empty flags from payload", func() {
		startMockServer(func(w http.ResponseWriter, r *http.Request) {
			body := decodeBody(r)
			for _, key := range []string{"workExperience", "eduLevel", "compNature", "address", "jobName", "companyName"} {
				Expect(body).NotTo(HaveKey(key))
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"ok": true})
		})

		captureStdout(func() { jobSearchCmd.Run(jobSearchCmd, nil) })
	})
})

// ─── job apply ──────────────────────────────────────────────────────

var _ = Describe("Job apply", func() {
	BeforeEach(func() {
		resetAllFlags()
	})

	It("sends jobId and jobKind to /mcp/apply-job", func() {
		startMockServer(func(w http.ResponseWriter, r *http.Request) {
			Expect(r.URL.Path).To(Equal("/mcp/apply-job"))

			body := decodeBody(r)
			Expect(body["jobId"]).To(Equal(float64(42)))
			Expect(body["jobKind"]).To(Equal("1"))

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"applied": true})
		})

		applyJobID = "42"
		applyJobKind = "1"

		out := captureStdout(func() {
			jobApplyCmd.Run(jobApplyCmd, nil)
		})
		Expect(out).To(ContainSubstring(`"applied": true`))
	})
})

// ─── resume get ─────────────────────────────────────────────────────

var _ = Describe("Resume get", func() {
	BeforeEach(func() {
		resetAllFlags()
	})

	It("sends GET to /mcp/get-resume", func() {
		startMockServer(func(w http.ResponseWriter, r *http.Request) {
			Expect(r.Method).To(Equal(http.MethodGet))
			Expect(r.URL.Path).To(Equal("/mcp/get-resume"))

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"resumeId": float64(123)})
		})

		out := captureStdout(func() {
			resumeGetCmd.Run(resumeGetCmd, nil)
		})
		Expect(out).To(ContainSubstring(`"resumeId": 123`))
	})
})

// ─── resume update-base-info ────────────────────────────────────────

var _ = Describe("Resume update-base-info", func() {
	BeforeEach(func() {
		resetAllFlags()
	})

	It("sends string fields", func() {
		startMockServer(func(w http.ResponseWriter, r *http.Request) {
			Expect(r.URL.Path).To(Equal("/mcp/update-base-info"))
			body := decodeBody(r)
			Expect(body["realName"]).To(Equal("Test User"))
			Expect(body["sex"]).To(Equal("男"))
			Expect(body["birthday"]).To(Equal("19950101"))
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"ok": true})
		})

		realName = "Test User"
		sex = "男"
		birthday = "19950101"

		captureStdout(func() { resumeUpdateBaseInfoCmd.Run(resumeUpdateBaseInfoCmd, nil) })
	})

	It("sends startJob/startJobMonth pair", func() {
		startMockServer(func(w http.ResponseWriter, r *http.Request) {
			body := decodeBody(r)
			Expect(body["startJob"]).To(Equal("2020"))
			Expect(body["startJobMonth"]).To(Equal("06"))
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"ok": true})
		})

		startJob = "2020"
		startJobMonth = "06"

		captureStdout(func() { resumeUpdateBaseInfoCmd.Run(resumeUpdateBaseInfoCmd, nil) })
	})

	It("sends numeric fields as numbers", func() {
		startMockServer(func(w http.ResponseWriter, r *http.Request) {
			body := decodeBody(r)
			Expect(body["nowSalary"]).To(Equal(float64(30000)))
			Expect(body["nowMonths"]).To(Equal(float64(13)))
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"ok": true})
		})

		nowSalary = "30000"
		nowMonths = "13"

		captureStdout(func() { resumeUpdateBaseInfoCmd.Run(resumeUpdateBaseInfoCmd, nil) })
	})
})

// ─── resume update-self-assess ──────────────────────────────────────

var _ = Describe("Resume update-self-assess", func() {
	BeforeEach(func() {
		resetAllFlags()
	})

	It("sends selfAssess content", func() {
		startMockServer(func(w http.ResponseWriter, r *http.Request) {
			Expect(r.URL.Path).To(Equal("/mcp/update-self-assess"))
			body := decodeBody(r)
			Expect(body["selfAssess"]).To(Equal("Experienced backend engineer"))
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"ok": true})
		})

		selfAssess = "Experienced backend engineer"
		captureStdout(func() { resumeUpdateSelfAssessCmd.Run(resumeUpdateSelfAssessCmd, nil) })
	})
})

// ─── resume education experience ────────────────────────────────────

var _ = Describe("Resume education experience", func() {
	BeforeEach(func() {
		resetAllFlags()
	})

	Describe("add-edu-exp", func() {
		It("sends school, degree, dates", func() {
			startMockServer(func(w http.ResponseWriter, r *http.Request) {
				Expect(r.URL.Path).To(Equal("/mcp/add-edu-exp"))
				body := decodeBody(r)
				Expect(body["school"]).To(Equal("Peking University"))
				Expect(body["degree"]).To(Equal("040"))
				Expect(body["start"]).To(Equal("201909"))
				Expect(body["end"]).To(Equal("202306"))
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]any{"ok": true})
			})

			eduSchool = "Peking University"
			eduDegree = "040"
			eduStart = "201909"
			eduEnd = "202306"

			captureStdout(func() { resumeAddEduExpCmd.Run(resumeAddEduExpCmd, nil) })
		})
	})

	Describe("update-edu-exp", func() {
		It("sends eduId with update fields", func() {
			startMockServer(func(w http.ResponseWriter, r *http.Request) {
				body := decodeBody(r)
				Expect(body["eduId"]).To(Equal(float64(100)))
				Expect(body["school"]).To(Equal("Tsinghua University"))
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]any{"ok": true})
			})

			eduID = "100"
			eduSchool = "Tsinghua University"
			eduDegree = "040"
			eduStart = "201909"
			eduEnd = "202306"

			captureStdout(func() { resumeUpdateEduExpCmd.Run(resumeUpdateEduExpCmd, nil) })
		})
	})
})

// ─── resume work experience ─────────────────────────────────────────

var _ = Describe("Resume work experience", func() {
	BeforeEach(func() {
		resetAllFlags()
	})

	Describe("add-work-exp", func() {
		It("omits shieldComp when flag not set", func() {
			startMockServer(func(w http.ResponseWriter, r *http.Request) {
				body := decodeBody(r)
				Expect(body).NotTo(HaveKey("shieldComp"))
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]any{"ok": true})
			})

			rootCmd.SetArgs([]string{
				"resume", "add-work-exp",
				"--comp-name", "Liepin",
				"--work-start", "202301",
				"--work-end", "202402",
				"--rw-title", "Backend Engineer",
			})
			captureStdout(func() { rootCmd.Execute() })
		})

		It("includes shieldComp:true when flag is set", func() {
			startMockServer(func(w http.ResponseWriter, r *http.Request) {
				body := decodeBody(r)
				Expect(body["shieldComp"]).To(BeTrue())
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]any{"ok": true})
			})

			rootCmd.SetArgs([]string{
				"resume", "add-work-exp",
				"--comp-name", "Liepin",
				"--work-start", "202301",
				"--work-end", "202402",
				"--rw-title", "Backend Engineer",
				"--shield-comp",
			})
			captureStdout(func() { rootCmd.Execute() })
		})

		It("sends numeric fields as numbers", func() {
			startMockServer(func(w http.ResponseWriter, r *http.Request) {
				body := decodeBody(r)
				Expect(body["subordinate"]).To(Equal(float64(5)))
				Expect(body["months"]).To(Equal(float64(14)))
				Expect(body["salary"]).To(Equal(float64(25000)))
				Expect(body["workType"]).To(Equal(float64(1)))
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]any{"ok": true})
			})

			workCompName = "Liepin"
			workRwTitle = "Engineer"
			workStart = "202301"
			workEnd = "202402"
			workSubordinate = "5"
			workMonthsStr = "14"
			workSalaryStr = "25000"
			workTypeStr = "1"

			captureStdout(func() { resumeAddWorkExpCmd.Run(resumeAddWorkExpCmd, nil) })
		})
	})
})

// ─── resume project experience ──────────────────────────────────────

var _ = Describe("Resume project experience", func() {
	BeforeEach(func() {
		resetAllFlags()
	})

	Describe("add-project-exp", func() {
		It("sends name, position, dates", func() {
			startMockServer(func(w http.ResponseWriter, r *http.Request) {
				body := decodeBody(r)
				Expect(body["name"]).To(Equal("Recommendation System"))
				Expect(body["position"]).To(Equal("Team Lead"))
				Expect(body["start"]).To(Equal("202301"))
				Expect(body["end"]).To(Equal("202312"))
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]any{"ok": true})
			})

			projName = "Recommendation System"
			projPosition = "Team Lead"
			projStart = "202301"
			projEnd = "202312"

			captureStdout(func() { resumeAddProjectExpCmd.Run(resumeAddProjectExpCmd, nil) })
		})
	})
})

// ─── resume job want ────────────────────────────────────────────────

var _ = Describe("Resume job want", func() {
	BeforeEach(func() {
		resetAllFlags()
	})

	Describe("add-job-want", func() {
		It("sends industries as repeatable array", func() {
			startMockServer(func(w http.ResponseWriter, r *http.Request) {
				body := decodeBody(r)
				industries, ok := body["industries"].([]any)
				Expect(ok).To(BeTrue())
				Expect(industries).To(HaveLen(2))
				Expect(industries[0]).To(Equal("IT"))
				Expect(industries[1]).To(Equal("Finance"))
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]any{"ok": true})
			})

			wantJobtitle = "Java"
			wantDq = "Shanghai"
			wantIndustries = []string{"IT", "Finance"}

			captureStdout(func() { resumeAddJobWantCmd.Run(resumeAddJobWantCmd, nil) })
		})

		It("sends salary range as numbers", func() {
			startMockServer(func(w http.ResponseWriter, r *http.Request) {
				body := decodeBody(r)
				Expect(body["wantSalaryLow"]).To(Equal(float64(20000)))
				Expect(body["wantSalaryHigh"]).To(Equal(float64(40000)))
				Expect(body["wantSalaryMonths"]).To(Equal(float64(15)))
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]any{"ok": true})
			})

			wantJobtitle = "Java"
			wantDq = "Shanghai"
			wantSalaryLowStr = "20000"
			wantSalaryHighStr = "40000"
			wantSalaryMonthsStr = "15"

			captureStdout(func() { resumeAddJobWantCmd.Run(resumeAddJobWantCmd, nil) })
		})

		It("sends otherExpectDqs as array", func() {
			startMockServer(func(w http.ResponseWriter, r *http.Request) {
				body := decodeBody(r)
				dqs := body["otherExpectDqs"].([]any)
				Expect(dqs).To(ConsistOf("Beijing", "Shenzhen"))
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]any{"ok": true})
			})

			wantJobtitle = "Java"
			wantDq = "Shanghai"
			wantOtherExpectDqs = []string{"Beijing", "Shenzhen"}

			captureStdout(func() { resumeAddJobWantCmd.Run(resumeAddJobWantCmd, nil) })
		})
	})

	Describe("update-job-want", func() {
		It("sends id with update fields", func() {
			startMockServer(func(w http.ResponseWriter, r *http.Request) {
				body := decodeBody(r)
				Expect(body["id"]).To(Equal(float64(99)))
				Expect(body["jobtitle"]).To(Equal("Go Developer"))
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]any{"ok": true})
			})

			wantID = "99"
			wantJobtitle = "Go Developer"
			wantDq = "Shanghai"

			captureStdout(func() { resumeUpdateJobWantCmd.Run(resumeUpdateJobWantCmd, nil) })
		})
	})
})

// ─── input file merging ─────────────────────────────────────────────

var _ = Describe("Payload file merging", func() {
	BeforeEach(func() {
		resetAllFlags()
	})

	It("merges input file with overrides, preferring explicit values", func() {
		dir, _ := os.MkdirTemp("", "payload-test")
		defer os.RemoveAll(dir)
		path := filepath.Join(dir, "base.json")
		os.WriteFile(path, []byte(`{"jobId":12,"jobKind":"2"}`), 0644)
		inputFlag = path

		applyJobID = "42"
		overrides := map[string]any{}
		idVal, _ := parseOptionalIntForTest(applyJobID)
		if idVal != nil {
			overrides["jobId"] = *idVal
		}
		overrides["jobKind"] = "1"

		payload, err := buildPayload(overrides, func() error { return nil })
		Expect(err).NotTo(HaveOccurred())
		Expect(payload["jobId"]).To(Equal(42))
		Expect(payload["jobKind"]).To(Equal("1"))
	})

	It("returns error from validation", func() {
		overrides := map[string]any{"jobId": 1}
		_, err := buildPayload(overrides, func() error {
			return &testValidationError{}
		})
		Expect(err).To(HaveOccurred())
	})
})

type testValidationError struct{}

func (e *testValidationError) Error() string { return "mock validation failure" }

// ─── output format ──────────────────────────────────────────────────

var _ = Describe("Output format", func() {
	BeforeEach(func() {
		resetAllFlags()
	})

	It("renders valid JSON in json mode", func() {
		startMockServer(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"resumeId": float64(1)})
		})
		outputFlag = "json"

		out := captureStdout(func() {
			resumeGetCmd.Run(resumeGetCmd, nil)
		})
		var result map[string]any
		Expect(json.Unmarshal([]byte(out), &result)).To(Succeed())
		Expect(result["resumeId"]).To(Equal(float64(1)))
	})

	It("renders plain text as-is in pretty mode", func() {
		startMockServer(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte("hello world"))
		})
		outputFlag = "pretty"

		out := captureStdout(func() {
			resumeGetCmd.Run(resumeGetCmd, nil)
		})
		Expect(out).To(ContainSubstring("hello world"))
	})
})

// ─── auth status / clear ────────────────────────────────────────────

var _ = Describe("Auth status", func() {
	BeforeEach(func() {
		resetAllFlags()
	})

	It("shows no-token message when config is empty", func() {
		dir, _ := os.MkdirTemp("", "auth-test")
		defer os.RemoveAll(dir)
		os.Setenv("XDG_CONFIG_HOME", dir)
		defer os.Unsetenv("XDG_CONFIG_HOME")

		out := captureStdout(func() {
			authStatusCmd.Run(authStatusCmd, nil)
		})
		Expect(out).To(ContainSubstring("No saved token"))
	})

	It("shows masked token when token is set", func() {
		dir, _ := os.MkdirTemp("", "auth-test")
		defer os.RemoveAll(dir)
		os.Setenv("XDG_CONFIG_HOME", dir)
		defer os.Unsetenv("XDG_CONFIG_HOME")
		authstore.SaveToken("my-secret-token-value")

		out := captureStdout(func() {
			authStatusCmd.Run(authStatusCmd, nil)
		})
		Expect(out).To(ContainSubstring("my-secr"))
		Expect(out).NotTo(ContainSubstring("my-secret-token-value"))
	})
})

var _ = Describe("Auth clear", func() {
	BeforeEach(func() {
		resetAllFlags()
	})

	It("clears token and prints message", func() {
		dir, _ := os.MkdirTemp("", "auth-test")
		defer os.RemoveAll(dir)
		os.Setenv("XDG_CONFIG_HOME", dir)
		defer os.Unsetenv("XDG_CONFIG_HOME")
		authstore.SaveToken("some-token")

		out := captureStdout(func() {
			authClearCmd.Run(authClearCmd, nil)
		})
		Expect(out).To(ContainSubstring("cleared"))

		token, _ := authstore.LoadToken()
		Expect(token).To(BeEmpty())
	})

	It("prints message when nothing to clear", func() {
		dir, _ := os.MkdirTemp("", "auth-test")
		defer os.RemoveAll(dir)
		os.Setenv("XDG_CONFIG_HOME", dir)
		defer os.Unsetenv("XDG_CONFIG_HOME")

		out := captureStdout(func() {
			authClearCmd.Run(authClearCmd, nil)
		})
		Expect(out).To(ContainSubstring("No saved token"))
	})
})

// ─── parseOptionalInt test helper ───────────────────────────────────

type strconvError struct{ s string }

func (e *strconvError) Error() string { return "invalid integer: " + e.s }

func parseOptionalIntForTest(s string) (*int, error) {
	if s == "" {
		return nil, nil
	}
	v := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return nil, &strconvError{s}
		}
		v = v*10 + int(c-'0')
	}
	return &v, nil
}
