package validate_test

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/pansachin/employee-service/pkg/validate"
)

// Success and failure markers.
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// ErrorContains checks if the error message in out contains the text in
// want.
//
// This is safe when out is nil. Use an empty string for want if you want to
// test that err is nil.
func ErrorContains(out error, want string) bool {
	if out == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(out.Error(), want)
}

func NewUnit(t *testing.T) (*zap.SugaredLogger, func()) {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()

	var buf bytes.Buffer
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	writer := bufio.NewWriter(&buf)
	log := zap.New(zapcore.NewCore(encoder, zapcore.AddSync(writer), zapcore.ErrorLevel)).Sugar()

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		_ = log.Sync()
		_ = writer.Flush()
		log.Info("******************** LOGS ********************")
		log.Info(buf.String())
		log.Info("******************** LOGS ********************")
	}

	return log, teardown
}

func Test_CheckID(t *testing.T) {
	_, teardown := NewUnit(t)
	t.Cleanup(teardown)

	testID := 1
	t.Logf("Test:\tValidate Positive IDs (like database ids)")
	{

		cases := []struct {
			in       string
			expected string
		}{
			{
				in:       "1",
				expected: "",
			},
			{
				in:       "0",
				expected: "value cannot be zero",
			},
			{
				in:       "-1",
				expected: "value cannot be negative",
			},
			{
				in:       "!",
				expected: "! is not a valid number",
			},
			{
				in:       "a",
				expected: "a is not a valid number",
			},
			{
				// Biggest max unit 64 ... plus an additional two 0's
				// This number is TOO BIG
				in:       fmt.Sprintf("%d00", 1<<32-1),
				expected: "value cannot be greater than 4294967295",
			},
		}

		for _, tc := range cases {
			got := validate.CheckID(tc.in)
			if !ErrorContains(got, tc.expected) {
				t.Logf("%s\tTest %d:\tvalidate.CheckID(%q)", Failed, testID, tc.in)
				t.Fatalf("%s\t\tExpected: %q, Got: %q", Failed, tc.expected, got)
			} else {
				t.Logf("%s\tTest %d:\tvalidate.CheckID(%q)", Success, testID, tc.in)
			}
			testID++
		}
	}
}

func Test_CheckUUID(t *testing.T) {
	_, teardown := NewUnit(t)
	t.Cleanup(teardown)

	testID := 1
	t.Logf("Test:\tValidate UUIDs (v4)")
	{
		tu, _ := uuid.NewRandom()
		tus := tu.String()
		cases := []struct {
			in       string
			expected string
		}{
			{
				in:       tus,
				expected: "",
			},
			{
				in:       "",
				expected: "UUID is not in its proper form",
			},
			{
				in:       "  ",
				expected: "UUID is not in its proper form",
			},
			{
				in:       "123",
				expected: "UUID is not in its proper form",
			},
			{
				in:       fmt.Sprintf("%s0", tus),
				expected: "UUID is not in its proper form",
			},
		}

		for _, tc := range cases {
			got := validate.CheckUUID(tc.in)
			if !ErrorContains(got, tc.expected) {
				t.Logf("%s\tTest %d:\tvalidate.CheckUUID(%q)", Failed, testID, tc.in)
				t.Fatalf("%s\t\tExpected: %q, Got: %q", Failed, tc.expected, got)
			} else {
				t.Logf("%s\tTest %d:\tvalidate.CheckUUID(%q)", Success, testID, tc.in)
			}
			testID++
		}
	}
}

func Test_CheckString(t *testing.T) {
	_, teardown := NewUnit(t)
	t.Cleanup(teardown)

	testID := 1
	t.Logf("Test:\tValidate Strings")
	{
		cases := []struct {
			in       string
			expected string
		}{
			{
				in:       "x",
				expected: "",
			},
			{
				in:       "",
				expected: "string can not be blank",
			},
			{
				in:       "  ",
				expected: "string can not be blank",
			},
			{
				in:       "123",
				expected: "",
			},
		}

		for _, tc := range cases {
			got := validate.CheckString(tc.in)
			if !ErrorContains(got, tc.expected) {
				t.Logf("%s\tTest %d:\tvalidate.CheckString(%q)", Failed, testID, tc.in)
				t.Fatalf("%s\t\tExpected: %q, Got: %q", Failed, tc.expected, got)
			} else {
				t.Logf("%s\tTest %d:\tvalidate.CheckString(%q)", Success, testID, tc.in)
			}
			testID++
		}
	}
}

func Test_CheckSlug(t *testing.T) {
	_, teardown := NewUnit(t)
	t.Cleanup(teardown)

	testID := 1
	t.Logf("Test:\tValidate Slug")
	{
		cases := []struct {
			in       string
			expected string
		}{
			{
				in:       "sample",
				expected: "",
			},
			{
				in:       "-sample",
				expected: "invalid slug",
			},
		}

		for _, tc := range cases {
			got := validate.CheckSlug(tc.in)
			if !ErrorContains(got, tc.expected) {
				t.Logf("%s\tTest %d:\tvalidate.CheckSlug(%q)", Failed, testID, tc.in)
				t.Fatalf("%s\t\tExpected: %q, Got: %q", Failed, tc.expected, got)
			} else {
				t.Logf("%s\tTest %d:\tvalidate.CheckSlug(%q)", Success, testID, tc.in)
			}
			testID++
		}
	}
}

func Test_Check(t *testing.T) {
	type sample struct {
		UID  string `validate:"uuid"`
		ID   int    `validate:"required"`
		Str  string `validate:"omitempty,required,notblank"`
		Str2 string `validate:"omitempty"`
	}
	_, teardown := NewUnit(t)
	t.Cleanup(teardown)

	testID := 1
	t.Logf("Test:\tValidate Struct")
	{
		cases := []struct {
			in       sample
			expected string
		}{
			{
				in:       sample{},
				expected: "{\"FieldError\":[{\"field\":\"UID\",\"error\":\"UID must be a valid UUID\"},{\"field\":\"ID\",\"error\":\"ID is a required field\"}],\"omitempty\":\"\"}",
			},
			{
				in: sample{
					UID:  "f64d10bf-54b2-44e9-8de7-191fc75398be",
					ID:   1,
					Str:  "123",
					Str2: "234",
				},
				expected: "",
			},
		}

		for _, tc := range cases {
			got := validate.Check(tc.in)
			if !ErrorContains(got, tc.expected) {
				t.Logf("%s\tTest %d:\tvalidate.Check(%q)", Failed, testID, tc.in)
				t.Fatalf("%s\t\tExpected: %q, Got: %q", Failed, tc.expected, got)
			} else {
				t.Logf("%s\tTest %d:\tvalidate.Check(%q)", Success, testID, tc.in)
			}
			testID++
		}
	}
}
