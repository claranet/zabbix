package zabbix_test

import (
	"testing"

	dd "github.com/extrading/go-zabbix-api"
)

func testCreateLLDRule(template *dd.Template, t *testing.T) *dd.LLDRule {
	lldrulefiltercondition := dd.LLDRulesFilterCondition{
		LLDMacro: "{#MY_MACRO_NAME_FD}",
		Value:    "^lo$",
		Operator: 8,
	}
	lldrulefilter := dd.LLDRuleFilter{
		Conditions: dd.LLDRulesFilterConditions{lldrulefiltercondition},
		EvalType:   0,
	}
	lldrules := dd.LLDRules{{
		Delay:       "34",
		HostID:      template.TemplateID,
		InterfaceID: "0",
		Key:         "lld_rule_key",
		Name:        "My low level discovery rule",
		Type:        dd.ZabbixAgent,
		Filter:      lldrulefilter,
	}}
	err := testGetAPI(t).DiscoveryRulesCreate(lldrules)
	if err != nil {
		t.Fatal(err)
	}
	return &lldrules[0]
}

func testDeleteLLDRule(rule *dd.LLDRule, t *testing.T) {
	err := testGetAPI(t).DiscoveryRulesDelete(dd.LLDRules{*rule})
	if err != nil {
		t.Fatal(err)
	}
}

func TestLLDRule(t *testing.T) {
	api := testGetAPI(t)

	hostGroup := testCreateHostGroup(t)
	defer testDeleteHostGroup(hostGroup, t)

	template := testCreateTemplate(hostGroup, t)
	defer testDeleteTemplate(template, t)

	lldRule := testCreateLLDRule(template, t)

	getlldRule, err := api.DiscoveryRulesGet(dd.Params{"itemids": lldRule.ItemID})
	if err != nil {
		t.Error(err)
	}
	if len(getlldRule) != 1 {
		t.Errorf("Expecting one result and got : %d", len(getlldRule))
	}

	lldRule.Name = "update_lld_name"
	err = api.DiscoveryRulesUpdate(dd.LLDRules{*lldRule})
	if err != nil {
		t.Error(err)
	}

	updateRule, err := api.DiscoveryRulesGetByID(lldRule.ItemID)
	if err != nil {
		t.Error(err)
	}
	if updateRule.Name != lldRule.Name {
		t.Errorf("LLD rule name is %q and should be %q", updateRule.Name, lldRule.Name)
	}

	testDeleteLLDRule(lldRule, t)
}
