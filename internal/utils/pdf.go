package utils

import (
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf"
	"github.com/nikolai/ai-resume-builder/backend/internal/models"
)

func GeneratePDF(resume models.ResumeContent) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Personal Information
	pdf.Cell(0, 10, resume.PersonalInfo.Name)
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 5, resume.PersonalInfo.Email)
	pdf.Ln(5)
	pdf.Cell(0, 5, resume.PersonalInfo.Phone)
	pdf.Ln(10)

	// Summary
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, "Professional Summary")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(0, 5, resume.Summary, "", "", false)
	pdf.Ln(5)

	// Experience
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, "Experience")
	pdf.Ln(8)
	for _, exp := range resume.Experience {
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(0, 5, fmt.Sprintf("%s - %s", exp.Company, exp.Title))
		pdf.Ln(5)
		pdf.SetFont("Arial", "", 10)
		pdf.Cell(0, 5, fmt.Sprintf("%s - %s", exp.StartDate.Format("Jan 2006"), getEndDate(exp)))
		pdf.Ln(5)
		for _, desc := range exp.Description {
			pdf.Cell(0, 5, "• "+desc)
			pdf.Ln(5)
		}
		pdf.Ln(3)
	}

	// Education
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, "Education")
	pdf.Ln(8)
	for _, edu := range resume.Education {
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(0, 5, edu.School)
		pdf.Ln(5)
		pdf.SetFont("Arial", "", 10)
		pdf.Cell(0, 5, fmt.Sprintf("%s in %s", edu.Degree, edu.Field))
		pdf.Ln(5)
		pdf.Cell(0, 5, fmt.Sprintf("%s - %s", edu.StartDate.Format("Jan 2006"), getEducationEndDate(edu)))
		pdf.Ln(8)
	}

	// Skills
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, "Skills")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(0, 5, formatSkills(resume.Skills), "", "", false)

	// Generate PDF bytes
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %v", err)
	}

	return buf.Bytes(), nil
}

func getEndDate(exp models.Experience) string {
	if exp.Current {
		return "Present"
	}
	if exp.EndDate != nil {
		return exp.EndDate.Format("Jan 2006")
	}
	return "Present"
}

func getEducationEndDate(edu models.Education) string {
	if edu.EndDate != nil {
		return edu.EndDate.Format("Jan 2006")
	}
	return "Present"
}

func formatSkills(skills []string) string {
	return fmt.Sprintf("• %s", formatBulletList(skills))
}

func formatBulletList(items []string) string {
	var result string
	for i, item := range items {
		if i == len(items)-1 {
			result += item
		} else {
			result += item + " • "
		}
	}
	return result
}
