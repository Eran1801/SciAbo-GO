package html


var EmailTemplates map[string]string = map[string]string{
  "reset_code": "utils/html/forget_password.html",
}

func GetEmailTemplate(templateName string) string {
  filePath := EmailTemplates[templateName]
  return filePath
}
