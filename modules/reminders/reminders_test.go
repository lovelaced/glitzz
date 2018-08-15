package reminders

import (
	"fmt"
	"testing"
	"time"
)

type parseMessageTest struct {
	UnitValue         string
	UnitTexts         []string
	ResultingDuration time.Duration
}

var parseMessageTests = []parseMessageTest{
	{
		UnitValue:         "7",
		UnitTexts:         []string{"s", "sec", "second", "seconds"},
		ResultingDuration: 7 * time.Second,
	},
	{
		UnitValue:         "7.00",
		UnitTexts:         []string{"s", "sec", "second", "seconds"},
		ResultingDuration: 7 * time.Second,
	},
	{
		UnitValue:         "7.31",
		UnitTexts:         []string{"s", "sec", "second", "seconds"},
		ResultingDuration: 7310 * time.Millisecond,
	},
}

func validateParseMessage(t *testing.T, p parseMessageTest) func(time.Duration, string, error) {
	return func(duration time.Duration, message string, err error) {
		t.Logf("testing %+v", p)
		if err != nil {
			t.Fatalf("error is not nil: %s", err)
		}
		if duration != p.ResultingDuration {
			t.Errorf("duration is invalid: %s", duration.String())
		}
		if message != "reminder message" {
			t.Errorf("messsage is invalid: %s", message)
		}
	}
}

func TestParseMessage(t *testing.T) {
	for _, test := range parseMessageTests {
		for _, unitText := range test.UnitTexts {
			u := test.UnitValue + unitText
			arguments := []string{u, "reminder", "message"}
			validateParseMessage(t, test)(parseMessage(arguments))
		}
	}
}

func TestParseMessageSeparateUnit(t *testing.T) {
	for _, test := range parseMessageTests {
		for _, unitText := range test.UnitTexts {
			arguments := []string{test.UnitValue, unitText, "reminder", "message"}
			validateParseMessage(t, test)(parseMessage(arguments))
		}
	}
}

func TestParseMessageInvalidNoUnit(t *testing.T) {
	u := fmt.Sprintf("%f", 7.0)
	arguments := []string{u, "reminder", "message"}
	_, _, err := parseMessage(arguments)
	if err == nil {
		t.Fatal("error is nil")
	}
}

func TestParseMessageInvalidNoAmountNoUnit(t *testing.T) {
	arguments := []string{"something", "reminder", "message"}
	_, _, err := parseMessage(arguments)
	if err == nil {
		t.Fatal("error is nil")
	}
}
