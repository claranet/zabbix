package zabbix

// LLDRulesFilterCondition represent zabbix low-level discovery rules filter condition(LLD rule file condition) object
// https://www.zabbix.com/documentation/3.2/manual/api/reference/discoveryrule/object#lld_rule_filter_condition
type LLDRulesFilterCondition struct {
	LLDMacro  string `json:"macro"` // Required
	Value     string `json:"value"` // Required
	FormulaID string `json:"formulaid,omitempty"`
	Operator  int    `json:"operator,omitempty,string"`
}

// LLDRulesFilterConditions is an array of LLDRulesFilterCondition
type LLDRulesFilterConditions []LLDRulesFilterCondition

// LLDRuleFilter represent zabbix low-level discovery rules filter(LLD rule filter) object
// https://www.zabbix.com/documentation/3.2/manual/api/reference/discoveryrule/object#lld_rule_filter
type LLDRuleFilter struct {
	Conditions  LLDRulesFilterConditions `json:"conditions"`      // Required
	EvalType    int                      `json:"evaltype,string"` // Required
	EvalFormula string                   `json:"eval_formula,omitempty"`
	Formula     string                   `json:"formula,omitempty"`
}

// LLDRule represent Zabbix low-level discovery rule(LLD rule) object
// https://www.zabbix.com/documentation/3.2/manual/api/reference/discoveryrule/object#lld_rule
type LLDRule struct {
	ItemID               string        `json:"itemid,omitempty"` // Readonly
	Delay                string        `json:"delay"`            // Required
	HostID               string        `json:"hostid"`           // Required
	InterfaceID          string        `json:"interfaceid"`      // Required
	Key                  string        `json:"key_"`             // Required
	Name                 string        `json:"name"`             // Required
	Type                 ItemType      `json:"type,string"`      // Required
	AuthType             string        `json:"authtype,omitempty"`
	DelayFlex            string        `json:"delay_flex,omitempty"`
	Description          string        `json:"description,omitempty"`
	Error                string        `json:"error,omitempty"` // Readonly
	IpmiSensor           string        `json:"ipmi_sensor,omitempty"`
	LifeTime             string        `json:"lifetime,omitempty"`
	Params               string        `json:"params,omitempty"`
	Password             string        `json:"password,omitempty"`
	Port                 string        `json:"port,omitempty"`
	PrivateKey           string        `json:"privatekey,omitempty"`
	PublicKey            string        `json:"publickey,omitempty"`
	SnmpCommunity        string        `json:"snmp_community,omitempty"`
	SnmpOid              string        `json:"snmp_oid,omitempty"`
	Snmpv3Authpassphrase string        `json:"snmpv3_authpassphrase,omitempty"`
	Snmpv3Authprotocol   int           `json:"snmpv3_authprotocol,omitempty,string"`
	Snmpv3Contextname    string        `json:"snmpv3_contextname,omitempty"`
	Snmpv3Privpassphrase string        `json:"snmpv3_privpassphrase,omitempty"`
	Snmpv3Privprotocol   int           `json:"snmpv3_privprotocol,omitempty,string"`
	Snmpv3Securitylevel  int           `json:"snmpv3_securitylevel,omitempty,string"`
	Snmpv3Securityname   string        `json:"snmpv3_securityname,omitempty"`
	State                int           `json:"state,omitempty,string"`
	Status               int           `json:"status,omitempty,string"`
	Templateid           string        `json:"templateid,omitempty"`
	TrapperHosts         string        `json:"trapper_hosts,omitempty"`
	Username             string        `json:"username,omitempty"`
	Filter               LLDRuleFilter `json:"filter"`
}

// LLDRules is an array of LLDRule
type LLDRules []LLDRule

// DiscoveryRulesGet Wrapper for discoveryrule.get
// https://www.zabbix.com/documentation/3.2/manual/api/reference/discoveryrule/get
func (api *API) DiscoveryRulesGet(params Params) (res LLDRules, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	err = api.CallWithErrorParse("discoveryrule.get", params, &res)
	return
}

// DiscoveryRulesGetByID Gets discovery rule by id only if there is exactly 1 matching discovery rule.
func (api *API) DiscoveryRulesGetByID(id string) (res *LLDRule, err error) {
	result, err := api.DiscoveryRulesGet(Params{"itemids": id})
	if err != nil {
		return
	}

	if len(result) == 1 {
		res = &result[0]
	} else {
		e := ExpectedOneResult(len(result))
		err = &e
	}
	return
}

// DiscoveryRulesCreate Wrapper for discoveryrule.create
// https://www.zabbix.com/documentation/3.2/manual/api/reference/discoveryrule/create
func (api *API) DiscoveryRulesCreate(rules LLDRules) error {
	result, err := api.CallWithError("discoveryrule.create", rules)
	if err != nil {
		return err
	}

	res := result.Result.(map[string]interface{})
	ids := res["itemids"].([]interface{})
	for i, id := range ids {
		rules[i].ItemID = id.(string)
	}
	return nil
}

// DiscoveryRulesUpdate Wrapper for discoveryrule.update
// https://www.zabbix.com/documentation/3.2/manual/api/reference/discoveryrule/update
func (api *API) DiscoveryRulesUpdate(rules LLDRules) error {
	_, err := api.CallWithError("discoveryrule.update", rules)
	return err
}

// DiscoveryRulesDelete Wrapper for discoveryrule.delete
// Cleans ItemID in all discovery rule element if call succeed.
// https://www.zabbix.com/documentation/3.2/manual/api/reference/discoveryrule/delete
func (api *API) DiscoveryRulesDelete(rules LLDRules) (err error) {
	var ids []string
	for _, rule := range rules {
		ids = append(ids, rule.ItemID)
	}

	err = api.DiscoveryRulesDeletesByIDs(ids)
	if err == nil {
		for i := range rules {
			rules[i].ItemID = ""
		}
	}
	return
}

// DiscoveryRulesDeletesByIDs  Wrapper for discorveryrule.delete
// Delete the discovery rule with the given ids
// https://www.zabbix.com/documentation/3.2/manual/api/reference/discoveryrule/delete
func (api *API) DiscoveryRulesDeletesByIDs(ids []string) (err error) {
	drulsids, err := api.DiscoveryRulesDeletesIDs(ids)
	if err != nil {
		return
	}

	if len(drulsids) != len(ids) {
		err = &ExpectedMore{len(ids), len(drulsids)}
	}
	return
}

// DiscoveryRulesDeletesIDs  Wrapper for discorveryrule.delete
// Delete the item and return the id of the deleted item
// https://www.zabbix.com/documentation/3.2/manual/api/reference/discoveryrule/delete
func (api *API) DiscoveryRulesDeletesIDs(ids []string) (drulsids []interface{}, err error) {
	res, err := api.CallWithError("discoveryrule.delete", ids)
	if err != nil {
		return
	}

	result := res.Result.(map[string]interface{})
	drulsids = result["ruleids"].([]interface{})
	return
}
