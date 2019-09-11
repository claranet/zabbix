package zabbix_test

import (
	"testing"

	. "."
)

func CreateTemplate(hostGroup *HostGroup, t *testing.T) *Template {
	template := Templates{Template{
		Host:   "template name",
		Groups: HostGroups{*hostGroup},
	}}
	err := getAPI(t).TemplatesCreate(template)
	if err != nil {
		t.Fatal(err)
	}
	return &template[0]
}

func DeleteTemplate(template *Template, t *testing.T) {
	err := getAPI(t).TemplatesDelete(Templates{*template})
	if err != nil {
		t.Fatal(err)
	}
}

func TestTemplates(t *testing.T) {
	api := getAPI(t)

	hostGroup := CreateHostGroup(t)
	defer DeleteHostGroup(hostGroup, t)

	templates, err := api.TemplatesGet(Params{})
	if err != nil {
		t.Fatal(err)
	}

	if len(templates) == 0 {
		t.Fatal("No templates were obtained")
	}

	template := CreateTemplate(hostGroup, t)
	if template.TemplateID == "" {
		t.Errorf("Template id is empty %#v", template)
	}

	template.Name = "new template name"
	err = api.TemplatesUpdate(Templates{*template})
	if err != nil {
		t.Error(err)
	}

	DeleteTemplate(template, t)
}
