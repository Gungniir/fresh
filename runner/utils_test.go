package runner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsWatchedFile(t *testing.T) {
	tests := []struct {
		file     string
		expected bool
	}{
		{"test.go", true},
		{"test.tpl", true},
		{"test.tmpl", true},
		{"test.html", true},
		{"test.css", false},
		{"test-executable", false},
		{"./tmp/test.go", false},
	}

	for _, test := range tests {
		actual := isWatchedFile(test.file)

		if actual != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, actual)
		}
	}
}

func TestShouldRebuild(t *testing.T) {
	tests := []struct {
		eventName string
		expected  bool
	}{
		{`"test.go": MODIFIED`, true},
		{`"test.tpl": MODIFIED`, false},
		{`"test.tmpl": DELETED`, false},
		{`"unknown.extension": DELETED`, true},
		{`"no_extension": ADDED`, true},
		{`"./a/path/test.go": MODIFIED`, true},
	}

	for _, test := range tests {
		actual := shouldRebuild(test.eventName)

		if actual != test.expected {
			t.Errorf("Expected %v, got %v (event was '%s')", test.expected, actual, test.eventName)
		}
	}
}

func TestIsIgnoredFolder(t *testing.T) {
	settings["ignored"] = "assets, tmp, **/build"

	tests := []struct {
		name     string
		dir      string
		expected bool
	}{
		{
			"assets",
			"assets",
			true,
		},
		{
			"assets node_modules",
			"assets/node_modules",
			true,
		},
		{
			"tmp",
			"tmp",
			true,
		},
		{
			"tmp pid",
			"tmp/pid",
			true,
		},
		{
			"app controllers",
			"app/controllers",
			false,
		},
		{
			"... build",
			"authorization-service/build",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := isIgnoredFolder(tt.dir)

			assert.Equal(t, tt.expected, actual)
		})
	}
}
