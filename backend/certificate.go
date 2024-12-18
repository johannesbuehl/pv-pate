package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

type CertificateData struct {
	Reservation  ReservationData
	TemplateData SponsorshipTemplateData
	PDFFile      string
}

type SponsorshipTemplateData struct {
	Element string
	Article string
	Date    string
	Name    string
}

var months = [12]string{
	"Januar", "Februar", "März", "April", "Mai", "Juni", "Juli", "August", "September", "Oktober", "November", "Dezember",
}

func (data *SponsorshipTemplateData) populate(mid, name string) {
	*data = SponsorshipTemplateData{
		Name:    name,
		Element: fmt.Sprintf("%s %s", getElementType(mid), getElementID(mid)),
		Article: getElementArticle(mid),
		Date:    time.Now().Format(fmt.Sprintf("2. %s 2006", months[time.Now().Month()-1])),
	}
}

func (data *CertificateData) create() error {
	// populate the template-data
	data.TemplateData.populate(data.Reservation.Mid, data.Reservation.Name)

	// choose the svg-template wether a name is given or not
	var templateName string

	if data.Reservation.Name == "" {
		templateName = "template_without_name.svg"
	} else {
		templateName = "template_with_name.svg"
	}

	// open the svg-template
	// create temporary svg file
	if svgFile, err := os.CreateTemp("templates", "certificate.*.svg"); err != nil {
		return err
	} else {
		defer os.Remove(svgFile.Name())
		defer svgFile.Close()

		if svgString, err := parseTemplate(path.Join("templates", templateName), data.TemplateData); err != nil {
			return err
		} else {
			data.PDFFile = fmt.Sprintf("templates/certificate.%s.pdf", data.Reservation.Mid)

			// write the svg-template
			svgFile.WriteString(svgString)

			actionString := fmt.Sprintf(`--actions=export-filename:%s; export-area-page; export-do`, data.PDFFile)

			// create a pdf from the svg-file
			command := exec.Command("inkscape/AppRun", actionString, svgFile.Name())

			if err := command.Run(); err != nil {
				logger.Error().Msg(err.Error())

				return err
			}

			return nil
		}
	}
}

func (data CertificateData) send() error {
	email := mail.NewMSG()

	if subject, err := parseTemplate("templates/certificate_mail", data.TemplateData); err != nil {
		return err
	} else if bodyHTML, err := parseHTMLTemplate("templates/certificate_mail.html", data.TemplateData); err != nil {
		return err
	} else if bodyPlain, err := parseHTMLTemplate("templates/certificate_mail.txt", data.TemplateData); err != nil {
		return err
	} else {
		email.SetFrom(fmt.Sprintf("Klimaplus-Patenschaft <%s>", config.Mail.User)).AddTo(data.Reservation.Mail).SetSubject(subject)

		email.SetBody(mail.TextPlain, bodyPlain)

		email.AddAlternative(mail.TextHTML, bodyHTML)

		email.Attach(&mail.File{
			FilePath: data.PDFFile,
		})

		if mailClient, err := mailServer.Connect(); err != nil {
			logger.Fatal().Msgf("can't connect to to mail-server: %v", err)

			return err
		} else if err := email.Send(mailClient); err != nil {
			return err
		} else {
			return nil
		}
	}
}

func (data *CertificateData) cleanup() error {
	if data.PDFFile != "" {
		return os.Remove(data.PDFFile)
	} else {
		return nil
	}
}
