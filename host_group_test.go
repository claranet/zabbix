package zabbix_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	zapi "github.com/extrading/go-zabbix-api"
)

func testCreateHostGroup(t *testing.T) *zapi.HostGroup {
	hostGroups := zapi.HostGroups{{Name: fmt.Sprintf("zabbix-testing-%d", rand.Int())}}
	err := testGetAPI(t).HostGroupsCreate(hostGroups)
	if err != nil {
		t.Fatal(err)
	}
	return &hostGroups[0]
}

func testDeleteHostGroup(hostGroup *zapi.HostGroup, t *testing.T) {
	err := testGetAPI(t).HostGroupsDelete(zapi.HostGroups{*hostGroup})
	if err != nil {
		t.Fatal(err)
	}
}

func TestHostGroups(t *testing.T) {
	api := testGetAPI(t)

	groups, err := api.HostGroupsGet(zapi.Params{})
	if err != nil {
		t.Fatal(err)
	}

	hostGroup := testCreateHostGroup(t)
	if hostGroup.GroupID == "" || hostGroup.Name == "" {
		t.Errorf("Something is empty: %#v", hostGroup)
	}

	hostGroup2, err := api.HostGroupGetByID(hostGroup.GroupID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(hostGroup, hostGroup2) {
		t.Errorf("Error getting group.\nOld group: %#v\nNew group: %#v", hostGroup, hostGroup2)
	}

	groups2, err := api.HostGroupsGet(zapi.Params{})
	if err != nil {
		t.Fatal(err)
	}
	if len(groups2) != len(groups)+1 {
		t.Errorf("Error creating group.\nOld groups: %#v\nNew groups: %#v", groups, groups2)
	}

	testDeleteHostGroup(hostGroup, t)

	groups2, err = api.HostGroupsGet(zapi.Params{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(groups, groups2) {
		t.Errorf("Error deleting group.\nOld groups: %#v\nNew groups: %#v", groups, groups2)
	}
}
