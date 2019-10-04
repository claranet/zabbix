package zabbix

import "github.com/AlekSi/reflector"

type Macro struct {
	MacroID   string `json:"hostmacroids,omitempty"`
	HostID    string `json:"hostid"`
	MacroName string `json:"macro"`
	Value     string `json:"value"`
}

type Macros []Macro

func (api *API) MacroGet(params Params) (res Macros, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("usermacro.get", params)
	if err != nil {
		return
	}
	reflector.MapsToStructs2(response.Result.([]interface{}), &res, reflector.Strconv, "json")
	return
}

func (api *API) MacroGetByID(id string) (res *Macro, err error) {
	triggers, err := api.MacroGet(Params{"hostmacroids": id})
	if err != nil {
		return
	}

	if len(triggers) == 1 {
		res = &triggers[0]
	} else {
		e := ExpectedOneResult(len(triggers))
		err = &e
	}
	return
}

func (api *API) MacroCreate(macros Macros) error {
	response, err := api.CallWithError("usermacro.create", macros)
	if err != nil {
		return err
	}

	result := response.Result.(map[string]interface{})
	macroids := result["hostmacroids"].([]interface{})
	for i, id := range macroids {
		macros[i].HostID = id.(string)
	}
	return nil
}

func (api *API) MacroUpdate(macros Macros) (err error) {
	_, err = api.CallWithError("usermacro.create", macros)
	return
}

func (api *API) MacroDeleteByID(ids []string) (err error) {
	response, err := api.CallWithError("usermacro.delete", ids)

	result := response.Result.(map[string]interface{})
	hostmacroids := result["hostmacroids"].([]interface{})
	if len(ids) != len(hostmacroids) {
		err = &ExpectedMore{len(ids), len(hostmacroids)}
	}
	return
}

func (api *API) MacroDelete(macros Macros) (err error) {
	ids := make([]string, len(macros))
	for i, macro := range macros {
		ids[i] = macro.MacroID
	}

	err = api.MacroDeleteByID(ids)
	if err == nil {
		for i := range macros {
			macros[i].MacroID = ""
		}
	}
	return
}
