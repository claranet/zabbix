package zabbix_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	zapi "github.com/extrading/go-zabbix-api"
)

func testCreateApplication(host *zapi.Host, t *testing.T) *zapi.Application {
	apps := zapi.Applications{{HostID: host.HostID, Name: fmt.Sprintf("App %d for %s", rand.Int(), host.Host)}}
	err := testGetAPI(t).ApplicationsCreate(apps)
	if err != nil {
		t.Fatal(err)
	}
	return &apps[0]
}

func testDeleteApplication(app *zapi.Application, t *testing.T) {
	err := testGetAPI(t).ApplicationsDelete(zapi.Applications{*app})
	if err != nil {
		t.Fatal(err)
	}
}

func TestApplications(t *testing.T) {
	api := testGetAPI(t)

	group := testCreateHostGroup(t)
	defer testDeleteHostGroup(group, t)

	host := testCreateHost(group, t)
	defer testDeleteHost(host, t)

	app := testCreateApplication(host, t)
	if app.ApplicationID == "" {
		t.Errorf("Id is empty: %#v", app)
	}

	app2 := testCreateApplication(host, t)
	if app2.ApplicationID == "" {
		t.Errorf("Id is empty: %#v", app2)
	}
	if reflect.DeepEqual(app, app2) {
		t.Errorf("Apps are equal:\n%#v\n%#v", app, app2)
	}

	apps, err := api.ApplicationsGet(zapi.Params{"hostids": host.HostID})
	if err != nil {
		t.Fatal(err)
	}
	if len(apps) != 2 {
		t.Errorf("Failed to create apps: %#v", apps)
	}

	app2, err = api.ApplicationGetByID(app.ApplicationID)
	if err != nil {
		t.Fatal(err)
	}
	app2.TemplateID = ""
	if !reflect.DeepEqual(app, app2) {
		t.Errorf("Apps are not equal:\n%#v\n%#v", app, app2)
	}

	app2, err = api.ApplicationGetByHostIDAndName(host.HostID, app.Name)
	if err != nil {
		t.Fatal(err)
	}
	app2.TemplateID = ""
	if !reflect.DeepEqual(app, app2) {
		t.Errorf("Apps are not equal:\n%#v\n%#v", app, app2)
	}

	testDeleteApplication(app, t)
}
