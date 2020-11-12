package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// JSON suffix
const JSON = "json"

// JSONReader defines the contract for json file handler
type JSONReader struct{}

// NewJSONReader creates
func NewJSONReader() *JSONReader {
	return &JSONReader{}
}

// List lists json files in specified directory
func (JSONReader) List(directory string) (*[]string, error) {
	var files []string

	err := filepath.Walk(fmt.Sprint(directory), func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, JSON) {
			files = append(files, path)
		}
		return nil
	})

	return &files, err
}

// Read the specified file
func (reader JSONReader) Read(filepath string) *Activity {
	raw := RawActivity{}

	file, _ := ioutil.ReadFile(filepath)
	// Assumptions 9: JSON is always valid
	json.Unmarshal(file, &raw)

	days := make(map[int]time.Time)
	meals := make(map[time.Time]int)

	for date, day := range raw.Calendar.DateDayMap {
		if t, e := ParseTime(date); e == nil {
			days[day] = t
		}
	}

	for _, day := range raw.Calendar.MealDayMap {
		// Assumption 6: Meal ids appearing in the meal to day map are valid
		// Assumption 7: Meal can have 0 or more dishes
		if _, exists := meals[days[day]]; exists {
			meals[days[day]]++
		} else {
			meals[days[day]] = 1
		}
	}

	return &Activity{
		UserID: reader.FileName(filepath),
		Meals:  &meals,
	}
}

// FileName returns filename from path
func (*JSONReader) FileName(path string) string {
	base := filepath.Base(path)
	return base[:strings.LastIndex(base, JSON) - 1]
}
