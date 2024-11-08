package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"text/template"
	"time"
)

type CertificateData struct {
	Mid          string
	Name         string
	PDFFile      string
	TemplateData CertificateTemplate
}

type CertificateTemplate struct {
	Name    string
	Element string
	Article string
	Date    string
}

var months = [12]string{
	"Januar", "Februar", "März", "April", "Mai", "Juni", "Juli", "August", "September", "Oktober", "November", "Dezember",
}

func (data *CertificateData) create() error {
	now := time.Now()

	// populate the template-data
	data.TemplateData.Name = data.Name
	data.TemplateData.Element = fmt.Sprintf("%s %s", getElementType(data.Mid), getElementID(data.Mid))
	data.TemplateData.Article = getElementArticle(data.Mid)
	data.TemplateData.Date = time.Now().Format(fmt.Sprintf("2. %s 2006", months[now.Month()-1]))

	// choose the svg-template wether a name is given or not
	var templateName string

	if data.Name == "" {
		templateName = "template_without_name.svg"
	} else {
		templateName = "template_with_name.svg"
	}

	// open the svg-template
	if buf, err := os.ReadFile(path.Join("certificates", templateName)); err != nil {
		return err
	} else {
		// open the template
		if svgTemplate, err := template.New("svgTemplate").Parse(string(buf)); err != nil {
			return err
		} else {
			// create temporary svg file
			if svgFile, err := os.CreateTemp("certificates", "certificate.*.svg"); err != nil {
				return err
			} else {
				defer os.Remove(svgFile.Name())
				defer svgFile.Close()

				data.PDFFile = fmt.Sprintf("certificates/certificate.%s.pdf", data.Mid)

				// write the svg-template
				svgTemplate.Execute(svgFile, data.TemplateData)

				actionString := fmt.Sprintf(`--actions=export-filename:%s; export-area-page; export-do`, data.PDFFile)

				// create a pdf from the svg-file
				command := exec.Command("certificates/inkscape/AppRun", actionString, svgFile.Name())
				command.Stderr = os.Stderr
				command.Stdout = os.Stdout

				if err := command.Run(); err != nil {
					logger.Error().Msg(err.Error())

					return err
				}

				return nil
			}
		}
	}
}

func (data *CertificateData) cleanup() error {
	return os.Remove(data.PDFFile)
}
