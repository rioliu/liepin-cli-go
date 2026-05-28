package models_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/rioliu/liepin-cli-go/internal/models"
)

var _ = Describe("SearchJobInput", func() {
	It("rejects negative page", func() {
		page := -1
		m := models.SearchJobInput{Page: &page}
		Expect(m.Validate()).To(HaveOccurred())
	})

	It("allows zero page", func() {
		page := 0
		m := models.SearchJobInput{Page: &page}
		Expect(m.Validate()).To(Succeed())
	})

	It("allows nil page", func() {
		m := models.SearchJobInput{Page: nil}
		Expect(m.Validate()).To(Succeed())
	})
})

var _ = Describe("ApplyJobInput", func() {
	It("rejects missing jobID and kind", func() {
		m := models.ApplyJobInput{}
		Expect(m.Validate()).To(HaveOccurred())
	})

	It("rejects invalid jobKind", func() {
		m := models.ApplyJobInput{JobID: 1, JobKind: "social"}
		Expect(m.Validate()).To(HaveOccurred())

		m2 := models.ApplyJobInput{JobID: 1, JobKind: "3"}
		Expect(m2.Validate()).To(HaveOccurred())
	})

	It("accepts valid jobKind", func() {
		for _, k := range []string{"1", "2"} {
			m := models.ApplyJobInput{JobID: 1, JobKind: k}
			Expect(m.Validate()).To(Succeed())
		}
	})
})

var _ = Describe("UpdateBaseInfoInput", func() {
	It("requires startJob and startJobMonth as a pair", func() {
		sj := "2020"
		m := models.UpdateBaseInfoInput{StartJob: &sj}
		Expect(m.Validate()).To(HaveOccurred())

		sjm := "06"
		m2 := models.UpdateBaseInfoInput{StartJobMonth: &sjm}
		Expect(m2.Validate()).To(HaveOccurred())
	})

	It("accepts valid startJob pair", func() {
		sj, sjm := "2020", "06"
		m := models.UpdateBaseInfoInput{StartJob: &sj, StartJobMonth: &sjm}
		Expect(m.Validate()).To(Succeed())
	})

	It("rejects invalid sex", func() {
		s := "未知"
		m := models.UpdateBaseInfoInput{Sex: &s}
		Expect(m.Validate()).To(HaveOccurred())
	})

	It("accepts valid sex values", func() {
		for _, s := range []string{"男", "女"} {
			m := models.UpdateBaseInfoInput{Sex: &s}
			Expect(m.Validate()).To(Succeed())
		}
	})

	It("rejects invalid date formats", func() {
		bd := "1990-01-01"
		m := models.UpdateBaseInfoInput{Birthday: &bd}
		Expect(m.Validate()).To(HaveOccurred())

		bd2 := "19950101"
		m2 := models.UpdateBaseInfoInput{Birthday: &bd2}
		Expect(m2.Validate()).To(Succeed())
	})

	It("rejects invalid startJob format", func() {
		sj, sjm := "20", "06"
		m := models.UpdateBaseInfoInput{StartJob: &sj, StartJobMonth: &sjm}
		Expect(m.Validate()).To(HaveOccurred())

		sj2, sjm2 := "2020", "13"
		m2 := models.UpdateBaseInfoInput{StartJob: &sj2, StartJobMonth: &sjm2}
		Expect(m2.Validate()).To(HaveOccurred())
	})

	It("rejects invalid nowWorkStatus", func() {
		s := "9"
		m := models.UpdateBaseInfoInput{NowWorkStatus: &s}
		Expect(m.Validate()).To(HaveOccurred())
	})

	It("rejects invalid nameSecret", func() {
		s := "2"
		m := models.UpdateBaseInfoInput{NameSecret: &s}
		Expect(m.Validate()).To(HaveOccurred())
	})

	It("rejects negative nowSalary", func() {
		s := -1
		m := models.UpdateBaseInfoInput{NowSalary: &s}
		Expect(m.Validate()).To(HaveOccurred())
	})

	It("rejects negative nowMonths", func() {
		s := -1
		m := models.UpdateBaseInfoInput{NowMonths: &s}
		Expect(m.Validate()).To(HaveOccurred())
	})
})

var _ = Describe("AddEduExpInput", func() {
	It("rejects invalid degree", func() {
		m := models.AddEduExpInput{
			School: "Peking University", Start: "201909", End: "202306", Degree: "999",
		}
		Expect(m.Validate()).To(HaveOccurred())
	})

	It("rejects invalid tz", func() {
		tz := "2"
		m := models.AddEduExpInput{
			School: "Peking University", Start: "201909", End: "202306", Degree: "040", Tz: &tz,
		}
		Expect(m.Validate()).To(HaveOccurred())
	})

	It("rejects invalid date formats", func() {
		m := models.AddEduExpInput{
			School: "Peking University", Start: "2019-09", End: "202306", Degree: "040",
		}
		Expect(m.Validate()).To(HaveOccurred())

		m2 := models.AddEduExpInput{
			School: "Peking University", Start: "201909", End: "202313", Degree: "040",
		}
		Expect(m2.Validate()).To(HaveOccurred())
	})
})

var _ = Describe("AddWorkExpInput", func() {
	It("rejects invalid workType", func() {
		wt := 3
		m := models.AddWorkExpInput{
			CompName: "Liepin", RwTitle: "Backend Engineer",
			WorkStart: "202301", WorkEnd: "202402", WorkType: &wt,
		}
		Expect(m.Validate()).To(HaveOccurred())
	})

	It("rejects negative numeric fields", func() {
		neg := -1
		m := models.AddWorkExpInput{
			CompName: "Liepin", RwTitle: "Backend Engineer",
			WorkStart: "202301", WorkEnd: "202402", Subordinate: &neg,
		}
		Expect(m.Validate()).To(HaveOccurred())

		m2 := models.AddWorkExpInput{
			CompName: "Liepin", RwTitle: "Backend Engineer",
			WorkStart: "202301", WorkEnd: "202402", Months: &neg,
		}
		Expect(m2.Validate()).To(HaveOccurred())

		m3 := models.AddWorkExpInput{
			CompName: "Liepin", RwTitle: "Backend Engineer",
			WorkStart: "202301", WorkEnd: "202402", Salary: &neg,
		}
		Expect(m3.Validate()).To(HaveOccurred())
	})

	It("rejects invalid YYYYMM dates", func() {
		m := models.AddWorkExpInput{
			CompName: "Liepin", RwTitle: "Backend Engineer",
			WorkStart: "20231", WorkEnd: "202402",
		}
		Expect(m.Validate()).To(HaveOccurred())

		m2 := models.AddWorkExpInput{
			CompName: "Liepin", RwTitle: "Backend Engineer",
			WorkStart: "202301", WorkEnd: "202400",
		}
		Expect(m2.Validate()).To(HaveOccurred())
	})
})

var _ = Describe("AddProjectExpInput", func() {
	It("rejects invalid YYYYMM dates", func() {
		m := models.AddProjectExpInput{
			Name: "Recommendation System", Start: "2023/01", End: "202312",
		}
		Expect(m.Validate()).To(HaveOccurred())

		m2 := models.AddProjectExpInput{
			Name: "Recommendation System", Start: "202301", End: "202313",
		}
		Expect(m2.Validate()).To(HaveOccurred())
	})
})

var _ = Describe("AddJobWantInput", func() {
	It("rejects invalid workType", func() {
		wt := "2"
		m := models.AddJobWantInput{Jobtitle: "Java", Dq: "Shanghai", WorkType: &wt}
		Expect(m.Validate()).To(HaveOccurred())
	})

	It("rejects negative numeric fields", func() {
		neg := -1
		m1 := models.AddJobWantInput{Jobtitle: "Java", Dq: "Shanghai", WantSalaryLow: &neg}
		Expect(m1.Validate()).To(HaveOccurred())

		m2 := models.AddJobWantInput{Jobtitle: "Java", Dq: "Shanghai", WantSalaryHigh: &neg}
		Expect(m2.Validate()).To(HaveOccurred())

		m3 := models.AddJobWantInput{Jobtitle: "Java", Dq: "Shanghai", WantSalaryMonths: &neg}
		Expect(m3.Validate()).To(HaveOccurred())

		m4 := models.AddJobWantInput{Jobtitle: "Java", Dq: "Shanghai", PracticeMonths: &neg}
		Expect(m4.Validate()).To(HaveOccurred())
	})

	It("rejects inverted salary range", func() {
		low, high := 30000, 20000
		m := models.AddJobWantInput{
			Jobtitle: "Java", Dq: "Shanghai", WantSalaryLow: &low, WantSalaryHigh: &high,
		}
		Expect(m.Validate()).To(HaveOccurred())
	})

	It("accepts valid salary range", func() {
		low, high := 20000, 30000
		m := models.AddJobWantInput{
			Jobtitle: "Java", Dq: "Shanghai", WantSalaryLow: &low, WantSalaryHigh: &high,
		}
		Expect(m.Validate()).To(Succeed())
	})
})

var _ = Describe("Helper functions", func() {
	Describe("StrPtr", func() {
		It("returns nil for empty string", func() {
			Expect(models.StrPtr("")).To(BeNil())
		})
		It("returns pointer for non-empty string", func() {
			p := models.StrPtr("abc")
			Expect(p).NotTo(BeNil())
			Expect(*p).To(Equal("abc"))
		})
	})

	Describe("ParseOptionalInt", func() {
		It("returns nil,nil for empty string", func() {
			v, err := models.ParseOptionalInt("")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(BeNil())
		})
		It("returns pointer for valid int", func() {
			v, err := models.ParseOptionalInt("42")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).NotTo(BeNil())
			Expect(*v).To(Equal(42))
		})
		It("returns error for invalid input", func() {
			_, err := models.ParseOptionalInt("abc")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("ParseOptionalBool", func() {
		It("returns nil when isSet=false", func() {
			Expect(models.ParseOptionalBool(false, false)).To(BeNil())
		})
		It("returns pointer to true", func() {
			p := models.ParseOptionalBool(true, true)
			Expect(p).NotTo(BeNil())
			Expect(*p).To(BeTrue())
		})
		It("returns pointer to false", func() {
			p := models.ParseOptionalBool(true, false)
			Expect(p).NotTo(BeNil())
			Expect(*p).To(BeFalse())
		})
	})

	Describe("ValidationErrors", func() {
		It("formats errors as a string", func() {
			errs := models.ValidationErrors{
				{Field: "a", Message: "is required"},
				{Field: "b", Message: "must be non-negative"},
			}
			Expect(errs.Error()).NotTo(BeEmpty())
		})
	})
})
