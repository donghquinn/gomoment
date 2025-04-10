package gomoment_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/donghquinn/gomoment"
)

func TestDefault(t *testing.T) {
	// 예제 사용법
	now := gomoment.Now()
	formattedNow, err := now.Format("YYYY-MM-DD HH:mm:ss")

	if err != nil {
		t.Fatalf("Formatting Error: %v:", err)
	} else {
		t.Logf("Now: %s:", formattedNow)
	}
}

func TestCustom(t *testing.T) {
	customTime := time.Date(2023, 5, 15, 14, 30, 45, 0, time.Local)
	customMoment, momentErr := gomoment.NewMoment(customTime)

	if momentErr != nil {
		t.Fatalf("error creating moment: %v", momentErr)
	}

	if s, err := customMoment.Format("YY/MM/DD h:mm A"); err == nil {
		t.Logf("Custom Format: %s:", s)
	} else if s != "23/05/15 2:30 P5" {
		t.Fail()
	} else {
		t.Fatalf("Custom Format YY/MM/DD h:mm A Error: %v", err)
	}

	if s, err := customMoment.Format("M/D/YY"); err == nil {
		t.Logf("Custom Format: %s:", s)
	} else if s != "5/15/23" {
		t.Fail()
	} else {
		t.Fatalf("Custom Format M/D/YY Error: %v", err)
	}

	// Invalid Format test
	if _, err := customMoment.Format("INVALID-FORMAT"); err == nil {
		t.Fail()
	} else {
		t.Logf("Invalid Format Err: %v", err)
	}

	if _, err := customMoment.Format("YYYY-QQ-DD"); err == nil {
		t.Fail()
	} else {
		t.Logf("Invalid Format Err: %v", err)
	}
}

func TestDefaultFormat(t *testing.T) {
	customTime := time.Date(2023, 5, 15, 14, 30, 45, 0, time.Local)
	customMoment, momentErr := gomoment.NewMoment(customTime)

	if momentErr != nil {
		t.Fatalf("error creating moment: %v", momentErr)
	}

	if s, err := customMoment.Format("YYYY-MM-DD HH:mm:ss"); err == nil {
		t.Logf("Default Format: %s:", s)
	} else if s != "2023-05-15 14:30:45" {
		t.Fail()
	} else {
		t.Logf("Default Formatting YYYY-MM-DD HH:mm:ss Error: %v:", err)
	}
}

func TestDefaulFormatWithTimezone(t *testing.T) {
	customTime := time.Date(2023, 5, 15, 14, 30, 45, 0, time.Local)
	customMoment, momentErr := gomoment.NewMoment(customTime)

	if momentErr != nil {
		t.Fatalf("error creating moment: %v", momentErr)
	}

	if s, err := customMoment.Format("YYYY-MM-DD HH:mm:ss Z"); err == nil {
		t.Logf("Format With Timezone: %s:", s)
	} else if s != "2023-05-15 14:30:45 +09:00" {
		t.Fail()
	} else {
		fmt.Println("Format With Timezone Error:", err)
	}
}

func TestInvalidFormatHandling(t *testing.T) {
	customTime := time.Date(2023, 5, 15, 14, 30, 45, 0, time.Local)
	customMoment, momentErr := gomoment.NewMoment(customTime)

	if momentErr != nil {
		t.Fatalf("error creating moment: %v", momentErr)
	}

	// Invalid Format test
	if _, err := customMoment.Format("INVALID-FORMAT"); err == nil {
		t.Fail()
	} else {
		t.Logf("Invalid Format Err: %v", err)
	}

	if _, err := customMoment.Format("YYYY-QQ-DD"); err == nil {
		t.Fail()
	} else {
		t.Logf("Invalid Format Err: %v", err)
	}
}

func TestNow(t *testing.T) {
	nowMoment, _ := gomoment.NewMoment()

	expected := time.Now().Local()
	if !expected.Local().Equal(nowMoment.Time()) {
		t.Fatalf("Moment now is invalid: %v", nowMoment.Time())
	}
}

func TestMust(t *testing.T) {
	customTime := time.Date(2023, 5, 15, 14, 30, 45, 0, time.Local)
	customMoment, momentErr := gomoment.NewMoment(customTime)

	if momentErr != nil {
		t.Fatalf("error creating moment: %v", momentErr)
	}

	if customMoment.Must("YYYY-MM-DD HH:mm:ss") != "2023-05-15 14:30:45" {
		t.Fail()
	}
}

func TestMomentWithString(t *testing.T) {
	timeString, err := gomoment.NewMoment("2025-04-10 13:41:00")
	if err != nil {
		t.Fatalf("Create Moment With String Error: %v", err)
	}

	expected := time.Date(2025, 4, 10, 13, 41, 0, 0, time.UTC)

	if !timeString.Time().Local().Equal(expected) {
		t.Fatalf("Created Time value is invalid: got %v, want %v", timeString.Time().Local(), expected)
	}
}

// Create Current time
func TestMomentTime(t *testing.T) {
	now, _ := gomoment.NewMoment() // Current time moment
	timeNow := time.Now()

	if now.Time() == timeNow {
		t.Fatalf("not accurate moment: %v and now: %v", now, timeNow)
	}
}
