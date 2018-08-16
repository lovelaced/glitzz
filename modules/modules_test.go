package modules

import (
	"github.com/lovelaced/glitzz/config"
	"testing"
)

func TestGetModuleConstructorsShouldReturnNonEmptyMap(t *testing.T) {
	modules := getModuleConstructors()
	if modules == nil {
		t.Fatal("modules was nil")
	}
	if len(modules) == 0 {
		t.Fatal("modules was empty")
	}
}

func TestGetModuleConstructorsShouldNotContainNulls(t *testing.T) {
	for key, value := range getModuleConstructors() {
		if value == nil {
			t.Errorf("nil value for key %s", key)
		}
	}
}

func TestCreateModulesModuleDoesNotExist(t *testing.T) {
	_, err := CreateModules(nil, config.Config{EnabledModules: []string{"invalid_module"}})
	if err == nil {
		t.Fatal("err was nil")
	}
}
