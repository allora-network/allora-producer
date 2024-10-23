package util

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/suite"
)

// RevisionTestSuite defines the test suite for Revision function.
type RevisionTestSuite struct {
	suite.Suite
	// originalReadBuildInfo holds the original readBuildInfo function to restore after tests.
	originalReadBuildInfo func() (*debug.BuildInfo, bool)
}

// SetupSuite runs before each test in the suite.
func (suite *RevisionTestSuite) SetupSuite() {
	// Save the original readBuildInfo function.
	suite.originalReadBuildInfo = readBuildInfo
}

// TearDownSuite runs after each test in the suite.
func (suite *RevisionTestSuite) TearDownSuite() {
	// Restore the original readBuildInfo function.
	readBuildInfo = suite.originalReadBuildInfo
}

// TestRevision tests the Revision function with various scenarios.
func (suite *RevisionTestSuite) TestRevision() {
	tests := []struct {
		name      string
		mockSetup func()
		want      string
	}{
		{
			name: "Valid Revision",
			mockSetup: func() {
				readBuildInfo = func() (*debug.BuildInfo, bool) {
					return &debug.BuildInfo{
						Settings: []debug.BuildSetting{
							{Key: "vcs.revision", Value: "abc123"},
						},
					}, true
				}
			},
			want: "abc123",
		},
		{
			name: "No Revision Key",
			mockSetup: func() {
				readBuildInfo = func() (*debug.BuildInfo, bool) {
					return &debug.BuildInfo{
						Settings: []debug.BuildSetting{
							{Key: "some.other.key", Value: "value"},
						},
					}, true
				}
			},
			want: "",
		},
		{
			name: "ReadBuildInfo Fails",
			mockSetup: func() {
				readBuildInfo = func() (*debug.BuildInfo, bool) {
					return nil, false
				}
			},
			want: "",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.mockSetup()
			got := Revision()
			suite.Equal(tt.want, got, "Revision() should return the expected value")
		})
	}
}

// TestRevisionTestSuite runs the RevisionTestSuite.
func TestRevisionTestSuite(t *testing.T) {
	suite.Run(t, new(RevisionTestSuite))
}
