package helper

import (
	"encoding/json"
	"go.uber.org/zap"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"time"
)

type Holidays struct {
	lieuDays map[string]struct{}
	holidays map[string]struct{}
}

func init() {
	time.Local = time.FixedZone("CST", 8*60*60)
}

func parseDays(data []string) []string {
	switch len(data) {
	case 1:
		return data
	case 2:
		s, err := time.Parse("2006-01-02", data[0])
		if err != nil {
			return nil
		}
		e, err := time.Parse("2006-01-02", data[1])
		if err != nil {
			return nil
		}
		var days []string
		for s.Before(e) || s.Equal(e) {
			days = append(days, s.Format("2006-01-02"))
			s = s.Add(time.Hour * 24)
		}
		return days
	default:
		zap.L().Error("Parse days failed", zap.Int("days", len(data)))
		return nil
	}
}

func (h *Holidays) loadFile(configFile string) {
	zapFileField := zap.String("path", configFile)
	zap.L().Info("Load config file", zapFileField)
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		zap.L().Error("Load config file failed", zapFileField, zap.Error(err))
		return
	}
	data := struct {
		LieuDays [][]string `json:"LieuDays"`
		Holidays [][]string `json:"Holidays"`
	}{}
	if err := json.Unmarshal(b, &data); err != nil {
		zap.L().Error("Unmarshal config file failed", zapFileField, zap.Error(err))
		return
	}
	for _, d := range data.LieuDays {
		days := parseDays(d)
		if days == nil {
			continue
		}
		for _, day := range days {
			h.lieuDays[day] = struct{}{}
		}
	}
	for _, d := range data.Holidays {
		days := parseDays(d)
		if days == nil {
			continue
		}
		for _, day := range days {
			h.holidays[day] = struct{}{}
		}
	}
}

func (h *Holidays) Load(configPath string) {
	zap.L().Info("Load config files", zap.String("path", configPath))
	var configFiles []string
	err := filepath.Walk(configPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		configFiles = append(configFiles, path)
		return nil
	})
	if err != nil {
		zap.L().Error("Load config files failed", zap.Error(err))
		return
	}

	newHolidays := &Holidays{
		make(map[string]struct{}),
		make(map[string]struct{}),
	}
	for _, configFile := range configFiles {
		newHolidays.loadFile(configFile)
	}
	h.lieuDays = newHolidays.lieuDays
	h.holidays = newHolidays.holidays
}

func (h *Holidays) IsHoliday(date time.Time) bool {
	day := date.Format("2006-01-02")
	if _, ok := h.holidays[day]; ok {
		return true
	}
	if _, ok := h.lieuDays[day]; ok {
		return false
	}
	return date.Weekday() == time.Saturday || date.Weekday() == time.Sunday
}

func (h *Holidays) IsTodayHoliday() bool {
	return h.IsHoliday(time.Now().Local())
}
