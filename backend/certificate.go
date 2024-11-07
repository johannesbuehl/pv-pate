package main

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

type CertificateData struct {
	Mid     string
	Element string
	Name    string
	PDFFile string
}

func (data *CertificateData) create() error {
	// open the svg-template
	if buf, err := os.ReadFile("certificates/template.svg"); err != nil {
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
				svgTemplate.Execute(svgFile, data)

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
