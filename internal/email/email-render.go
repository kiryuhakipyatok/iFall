package email

import (
	"bytes"
	_ "embed"
	"html/template"
	"iFall/internal/domain/models"
)

//go:embed templates/email.html
var verifyEmailHTML string

func BuildEmailLetter(iphones []models.IPhone) (string, error) {

	funcMap := template.FuncMap{
		"abs": func(x float64) float64 {
			if x < 0 {
				return -x
			}
			return x
		},
	}

	tmpl, err := template.New("email").Funcs(funcMap).Parse(verifyEmailHTML)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, iphones); err != nil {
		return "", err
	}
	return buf.String(), nil
}
