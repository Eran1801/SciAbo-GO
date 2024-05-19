package html


var email_templates map[string]string = map[string]string{
  "reset_code": "utils/html/forget_password.html",
  "new_massage": "utils/html/new_message.html",
}

func GetEmailTemplate(template_name string) string {
  filePath := email_templates[template_name]
  return filePath
}
