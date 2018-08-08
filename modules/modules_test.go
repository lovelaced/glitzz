package modules

import (
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
