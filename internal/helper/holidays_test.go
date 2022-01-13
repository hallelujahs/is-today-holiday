package helper

import (
	"fmt"
	"path"
	"runtime"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	configPath := path.Join(path.Dir(filename), "..", "..", "configs")
	fmt.Println(configPath)

	h := Holidays{}
	h.Load(configPath)
	if h.holidays == nil || len(h.holidays) == 0 {
		t.Error("Load holidays failed")
	}
	if h.lieuDays == nil || len(h.lieuDays) == 0 {
		t.Error("Load lieu days failed")
	}
}

func TestIsHoliday(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	configPath := path.Join(path.Dir(filename), "..", "..", "configs")
	fmt.Println(configPath)

	h := Holidays{}
	h.Load(configPath)

	holiday, err := time.Parse("2006-01-02", "2022-10-01")
	if err != nil {
		t.Fatal("Parse date failed")
	}
	if !h.IsHoliday(holiday) {
		t.Error("Should be holiday")
	}
}
