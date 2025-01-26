package email

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed templates/*.gotmpl
var emailTemplates embed.FS

func getTemplate(templateName string) ([]byte, error) {
	return emailTemplates.ReadFile("templates/" + templateName + ".gotmpl")
}

func RenderTemplate(templateName string, data interface{}) (string, error) {
	templateBytes, err := getTemplate(templateName)
	if err != nil {
		return "", err
	}
	t := template.Must(template.New(templateName).Parse(string(templateBytes)))
	var renderedTemplate bytes.Buffer
	err = t.Execute(&renderedTemplate, data)
	if err != nil {
		return "", err
	}
	return renderedTemplate.String(), nil
}
