package zabbix

import (
	"github.com/AlekSi/reflector"
)

// https://www.zabbix.com/documentation/2.2/manual/api/reference/template/object
type Template struct {
	TemplateId  string     `json:"templateid,omitempty"`
	Host        string     `json:"host"`
	Description string     `json:"description,omitempty"`
	Name        string     `json:"name,omitempty"`
	Groups      HostGroups `json:"groups"`
}

type Templates []Template

type TemplateId struct {
	TemplateId string `json:"templateid"`
}

type TemplateIds []TemplateId

// Wrapper for template.get: https://www.zabbix.com/documentation/2.2/manual/api/reference/template/get
func (api *API) TemplatesGet(params Params) (res Templates, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("template.get", params)
	if err != nil {
		return
	}

	reflector.MapsToStructs2(response.Result.([]interface{}), &res, reflector.Strconv, "json")
	return
}

func (api *API) TemplatesCreate(templates Templates) (err error) {
	response, err := api.CallWithError("template.create", templates)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	templateids := result["templateids"].([]interface{})
	for i, id := range templateids {
		templates[i].TemplateId = id.(string)
	}
	return
}

func (api *API) TemplatesDelete(templates Templates) (err error) {
	templatesIds := make([]string, len(templates))
	for i, template := range templates {
		templatesIds[i] = template.TemplateId
	}

	err = api.TemplatesDeleteByIds(templatesIds)
	if err == nil {
		for i := range templates {
			templates[i].TemplateId = ""
		}
	}
	return
}

func (api *API) TemplatesDeleteByIds(ids []string) (err error) {
	response, err := api.CallWithError("template.delete", ids)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	templateids := result["templateids"].([]interface{})
	if len(ids) != len(templateids) {
		err = &ExpectedMore{len(ids), len(templateids)}
	}
	return
}
