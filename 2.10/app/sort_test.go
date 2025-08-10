package main

import (
	"os"
	"testing"
)

func TestParseMonth(t *testing.T) {
	tests := []struct {
		input string
		want  int
		ok    bool
	}{
		{"Jan", 1, true},
		{"jan", 1, true},
		{"JAN", 1, true},
		{"Feb", 2, true},
		{"Mar", 3, true},
		{"Apr", 4, true},
		{"May", 5, true},
		{"Jun", 6, true},
		{"Jul", 7, true},
		{"Aug", 8, true},
		{"Sep", 9, true},
		{"Oct", 10, true},
		{"Nov", 11, true},
		{"Dec", 12, true},
		{"Invalid", 0, false},
		{"", 0, false},
		{" January ", 0, false},
		{"Jan ", 1, true},
		{" Jan", 1, true},
	}

	for _, tt := range tests {
		got, ok := parseMonth(tt.input)
		if got != tt.want || ok != tt.ok {
			t.Errorf("ParseMonth(%q) = (%d, %v), want (%d, %v)", tt.input, got, ok, tt.want, tt.ok)
		}
	}
}

func TestParseHumanNumber(t *testing.T) {
	tests := []struct {
		input string
		want  float64
		ok    bool
	}{
		{"1K", 1000, true},
		{"1.5K", 1500, true},
		{"1M", 1000000, true},
		{"1G", 1000000000, true},
		{"1T", 1000000000000, true},
		{"1P", 1000000000000000, true},
		{"1E", 1000000000000000000, true},
		{"1К", 1000, true},
		{"1М", 1000000, true},
		{"1Г", 1000000000, true},
		{"1Т", 1000000000000, true},
		{"123", 123, true},
		{"123.45", 123.45, true},
		{"-123", -123, true},
		{"-123.45", -123.45, true},
		{"1.5K", 1500, true},
		{"2.5M", 2500000, true},
		{"", 0, false},
		{"abc", 0, false},
		{"K", 0, false},
		{"1X", 1, true},
		{" 1K ", 1000, true},
	}

	for _, tt := range tests {
		got, ok := parseHumanNumber(tt.input)
		if got != tt.want || ok != tt.ok {
			t.Errorf("ParseHumanNumber(%q) = (%f, %v), want (%f, %v)", tt.input, got, ok, tt.want, tt.ok)
		}
	}
}

func TestParseNumber(t *testing.T) {
	tests := []struct {
		input string
		want  float64
		ok    bool
	}{
		{"123", 123, true},
		{"123.45", 123.45, true},
		{"-123", -123, true},
		{"-123.45", -123.45, true},
		{"0", 0, true},
		{"0.0", 0, true},
		{"", 0, false},
		{"abc", 0, false},
		{" 123 ", 123, true},
		{"123abc", 0, false},
	}

	for _, tt := range tests {
		got, ok := parseNumber(tt.input)
		if got != tt.want || ok != tt.ok {
			t.Errorf("ParseNumber(%q) = (%f, %v), want (%f, %v)", tt.input, got, ok, tt.want, tt.ok)
		}
	}
}

func TestKeyLess(t *testing.T) {
	tests := []struct {
		a    *key
		b    *key
		want bool
	}{
		{
			&key{strVal: "a", isNum: false},
			&key{strVal: "b", isNum: false},
			true,
		},
		{
			&key{strVal: "b", isNum: false},
			&key{strVal: "a", isNum: false},
			false,
		},
		{
			&key{strVal: "a", isNum: false},
			&key{strVal: "a", isNum: false},
			false,
		},
		{
			&key{numVal: 1, isNum: true},
			&key{numVal: 2, isNum: true},
			true,
		},
		{
			&key{numVal: 2, isNum: true},
			&key{numVal: 1, isNum: true},
			false,
		},
		{
			&key{numVal: 1, isNum: true},
			&key{numVal: 1, isNum: true},
			false,
		},
		{
			&key{strVal: "a", isNum: false},
			&key{numVal: 1, isNum: true},
			true,
		},
		{
			&key{numVal: 1, isNum: true},
			&key{strVal: "a", isNum: false},
			false,
		},
	}

	for _, tt := range tests {
		got := keyLess(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("KeyLess(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestApplicationKey(t *testing.T) {
	tests := []struct {
		name      string
		app       *application
		input     string
		wantStr   string
		wantNum   float64
		wantIsNum bool
	}{
		{
			name:      "basic string",
			app:       &application{},
			input:     "hello",
			wantStr:   "hello",
			wantNum:   0,
			wantIsNum: false,
		},
		{
			name:      "numeric sort",
			app:       &application{numeric: true},
			input:     "123",
			wantStr:   "123",
			wantNum:   123,
			wantIsNum: true,
		},
		{
			name:      "month sort",
			app:       &application{sortByMonth: true},
			input:     "Jan",
			wantStr:   "Jan",
			wantNum:   1,
			wantIsNum: true,
		},
		{
			name:      "human numeric sort",
			app:       &application{humanNumeric: true},
			input:     "1K",
			wantStr:   "1K",
			wantNum:   1000,
			wantIsNum: true,
		},
		{
			name:      "column extraction",
			app:       &application{column: 2},
			input:     "a\tb\tc",
			wantStr:   "b",
			wantNum:   0,
			wantIsNum: false,
		},
		{
			name:      "column out of range",
			app:       &application{column: 5},
			input:     "a\tb\tc",
			wantStr:   "",
			wantNum:   0,
			wantIsNum: false,
		},
		{
			name:      "ignore trailing blanks",
			app:       &application{ignoreTrailingBlanks: true},
			input:     "hello   ",
			wantStr:   "hello",
			wantNum:   0,
			wantIsNum: false,
		},
		{
			name:      "priority: month over numeric",
			app:       &application{sortByMonth: true, numeric: true},
			input:     "Jan",
			wantStr:   "Jan",
			wantNum:   1,
			wantIsNum: true,
		},
		{
			name:      "priority: human numeric over numeric",
			app:       &application{humanNumeric: true, numeric: true},
			input:     "1K",
			wantStr:   "1K",
			wantNum:   1000,
			wantIsNum: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := tt.app.key(tt.input)
			if key.strVal != tt.wantStr {
				t.Errorf("key.strVal = %q, want %q", key.strVal, tt.wantStr)
			}
			if key.numVal != tt.wantNum {
				t.Errorf("key.numVal = %f, want %f", key.numVal, tt.wantNum)
			}
			if key.isNum != tt.wantIsNum {
				t.Errorf("key.isNum = %v, want %v", key.isNum, tt.wantIsNum)
			}
			if key.original != tt.input {
				t.Errorf("key.original = %q, want %q", key.original, tt.input)
			}
		})
	}
}

func TestApplicationFullLess(t *testing.T) {
	app := &application{}

	tests := []struct {
		name string
		a    *key
		b    *key
		want bool
	}{
		{
			name: "different keys",
			a:    &key{strVal: "a", original: "a"},
			b:    &key{strVal: "b", original: "b"},
			want: true,
		},
		{
			name: "same keys, different originals",
			a:    &key{strVal: "a", original: "a"},
			b:    &key{strVal: "a", original: "b"},
			want: true,
		},
		{
			name: "same keys, same originals",
			a:    &key{strVal: "a", original: "a"},
			b:    &key{strVal: "a", original: "a"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := app.fullLess(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("FullLess(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestCheckSorted(t *testing.T) {
	tests := []struct {
		name  string
		app   *application
		lines []string
		want  bool
	}{
		{
			name:  "sorted ascending",
			app:   &application{},
			lines: []string{"a", "b", "c"},
			want:  true,
		},
		{
			name:  "sorted descending",
			app:   &application{reverse: true},
			lines: []string{"c", "b", "a"},
			want:  true,
		},
		{
			name:  "unsorted ascending",
			app:   &application{},
			lines: []string{"b", "a", "c"},
			want:  false,
		},
		{
			name:  "unsorted descending",
			app:   &application{reverse: true},
			lines: []string{"a", "b", "c"},
			want:  false,
		},
		{
			name:  "empty list",
			app:   &application{},
			lines: []string{},
			want:  true,
		},
		{
			name:  "single element",
			app:   &application{},
			lines: []string{"a"},
			want:  true,
		},
		{
			name:  "numeric sorted",
			app:   &application{numeric: true},
			lines: []string{"1", "2", "10"},
			want:  true,
		},
		{
			name:  "numeric unsorted",
			app:   &application{numeric: true},
			lines: []string{"2", "1", "10"},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.app.checkSorted(tt.lines)
			if got != tt.want {
				t.Errorf("CheckSorted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name      string
		keys      []*key
		lines     []string
		wantKeys  []*key
		wantLines []string
	}{
		{
			name:      "no duplicates",
			keys:      []*key{{strVal: "a"}, {strVal: "b"}, {strVal: "c"}},
			lines:     []string{"a", "b", "c"},
			wantKeys:  []*key{{strVal: "a"}, {strVal: "b"}, {strVal: "c"}},
			wantLines: []string{"a", "b", "c"},
		},
		{
			name:      "with duplicates",
			keys:      []*key{{strVal: "a"}, {strVal: "a"}, {strVal: "b"}, {strVal: "b"}, {strVal: "c"}},
			lines:     []string{"a1", "a2", "b1", "b2", "c"},
			wantKeys:  []*key{{strVal: "a"}, {strVal: "b"}, {strVal: "c"}},
			wantLines: []string{"a1", "b1", "c"},
		},
		{
			name:      "numeric duplicates",
			keys:      []*key{{numVal: 1, isNum: true}, {numVal: 1, isNum: true}, {numVal: 2, isNum: true}},
			lines:     []string{"1a", "1b", "2"},
			wantKeys:  []*key{{numVal: 1, isNum: true}, {numVal: 2, isNum: true}},
			wantLines: []string{"1a", "2"},
		},
		{
			name:      "mixed types",
			keys:      []*key{{strVal: "a", isNum: false}, {numVal: 1, isNum: true}, {strVal: "a", isNum: false}},
			lines:     []string{"a1", "1", "a2"},
			wantKeys:  []*key{{strVal: "a", isNum: false}, {numVal: 1, isNum: true}, {strVal: "a", isNum: false}},
			wantLines: []string{"a1", "1", "a2"},
		},
		{
			name:      "empty input",
			keys:      []*key{},
			lines:     []string{},
			wantKeys:  []*key{},
			wantLines: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKeys, gotLines := removeDuplicates(tt.keys, tt.lines)
			if len(gotKeys) != len(tt.wantKeys) {
				t.Errorf("RemoveDuplicates() returned %d keys, want %d", len(gotKeys), len(tt.wantKeys))
			}
			if len(gotLines) != len(tt.wantLines) {
				t.Errorf("RemoveDuplicates() returned %d lines, want %d", len(gotLines), len(tt.wantLines))
			}
		})
	}
}

func TestReadFromStdin(t *testing.T) {
	// Создаем временный файл для stdin
	content := "line1\nline2\nline3"
	tmpfile, err := os.CreateTemp("", "stdin")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Сохраняем оригинальный stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Открываем временный файл как stdin
	file, err := os.Open(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	os.Stdin = file
	defer file.Close()

	lines, err := readFromStdin()
	if err != nil {
		t.Errorf("ReadFromStdin() error = %v", err)
		return
	}

	expected := []string{"line1", "line2", "line3"}
	if len(lines) != len(expected) {
		t.Errorf("ReadFromStdin() returned %d lines, want %d", len(lines), len(expected))
		return
	}

	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("ReadFromStdin()[%d] = %q, want %q", i, line, expected[i])
		}
	}
}

func TestReadFromFile(t *testing.T) {
	// Создаем временный файл
	content := "line1\nline2\nline3"
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	lines, err := readFromFile(tmpfile.Name())
	if err != nil {
		t.Errorf("ReadFromFile() error = %v", err)
		return
	}

	expected := []string{"line1", "line2", "line3"}
	if len(lines) != len(expected) {
		t.Errorf("ReadFromFile() returned %d lines, want %d", len(lines), len(expected))
		return
	}

	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("ReadFromFile()[%d] = %q, want %q", i, line, expected[i])
		}
	}
}

func TestReadFromFileError(t *testing.T) {
	_, err := readFromFile("nonexistent_file")
	if err == nil {
		t.Error("ReadFromFile() should return error for nonexistent file")
	}
}

func TestReadLines(t *testing.T) {
	// Создаем временные файлы
	content1 := "file1_line1\nfile1_line2"
	content2 := "file2_line1\nfile2_line2"

	tmpfile1, err := os.CreateTemp("", "test1")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile1.Name())

	tmpfile2, err := os.CreateTemp("", "test2")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile2.Name())

	if _, err := tmpfile1.WriteString(content1); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile1.Close(); err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile2.WriteString(content2); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile2.Close(); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		files   []string
		want    []string
		wantErr bool
	}{
		{
			name:    "multiple files",
			files:   []string{tmpfile1.Name(), tmpfile2.Name()},
			want:    []string{"file1_line1", "file1_line2", "file2_line1", "file2_line2"},
			wantErr: false,
		},
		{
			name:    "single file",
			files:   []string{tmpfile1.Name()},
			want:    []string{"file1_line1", "file1_line2"},
			wantErr: false,
		},
		{
			name:    "nonexistent file",
			files:   []string{"nonexistent"},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readLines(tt.files)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("ReadLines() returned %d lines, want %d", len(got), len(tt.want))
					return
				}
				for i, line := range got {
					if line != tt.want[i] {
						t.Errorf("ReadLines()[%d] = %q, want %q", i, line, tt.want[i])
					}
				}
			}
		})
	}
}

func TestReadLinesNoFiles(t *testing.T) {
	_, err := readLines([]string{})
	if err != nil {
		t.Errorf("ReadLines() with no files should not return error, got %v", err)
	}
}
