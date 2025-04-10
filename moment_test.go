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
	customMoment := gomoment.NewMoment(customTime)

	if s, err := customMoment.Format("YYYY-MM-DD HH:mm:ss"); err == nil {
		t.Logf("Default Format: %s:", s)
	} else if s != "2023-05-15 14:30:45" {
		t.Fail()
	} else {
		t.Logf("Default Formatting YYYY-MM-DD HH:mm:ss Error: %v:", err)
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

	if s, err := customMoment.Format("YYYY-MM-DD HH:mm:ss Z"); err == nil {
		t.Logf("Format With Timezone: %s:", s)
	} else if s != "2023-05-15 14:30:45 +09:00" {
		t.Fail()
	} else {
		fmt.Println("Format With Timezone Error:", err)
	}

	// 잘못된 포맷 테스트
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

func TestMust(t *testing.T) {
	customTime := time.Date(2023, 5, 15, 14, 30, 45, 0, time.Local)
	customMoment := gomoment.NewMoment(customTime)

	if customMoment.Must("YYYY-MM-DD HH:mm:ss") != "2023-05-15 14:30:45" {
		t.Fail()
	}
}
