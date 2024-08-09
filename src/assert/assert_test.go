package assert

import (
	"testing"
	"time"
)

func TestStringifyInt(t *testing.T) {
	s := stringify(1)
	if s != "1" {
		t.Fatalf(`Expected "1", received %q`, s)
	}
}

func TestStringifyString(t *testing.T) {
	s := stringify("a")
	if s != "a" {
		t.Fatalf(`Expected "a", received %q`, s)
	}
}

func TestStringifyBool(t *testing.T) {
	s := stringify(true)
	if s != "true" {
		t.Fatalf(`Expected "true", received %q`, s)
	}
}

func TestStringifyNil(t *testing.T) {
	s := stringify(nil)
	if s != "nil" {
		t.Fatalf(`Expected "nil", received %q`, s)
	}
}

func TestStringifyBytes(t *testing.T) {
	s := stringify([]byte("a"))
	if s != "a" {
		t.Fatalf(`Expected "a", received %q`, s)
	}
}

func TestStringifyMap(t *testing.T) {
	s := stringify(map[string]string{"a": "b", "c": "d"})
	expected := "{\"a\":\"b\",\"c\":\"d\"}"
	if s != expected {
		t.Fatalf(`Expected "%s", received %s`, expected, s)
	}
}

func TestStringifyStruct(t *testing.T) {
	type TestStruct struct {
		A string
		B int
	}
	s := stringify(TestStruct{"a", 1})
	expected := "{\"A\":\"a\",\"B\":1}"
	if s != expected {
		t.Fatalf(`Expected "%s", received %s`, expected, s)
	}
}

func TestStringifySlice(t *testing.T) {
	s := stringify([]string{"a", "b", "c"})
	expected := "[\"a\",\"b\",\"c\"]"
	if s != expected {
		t.Fatalf(`Expected "%s", received %s`, expected, s)
	}
}

func TestStringifyTime(t *testing.T) {
	date := time.Date(2023, 3, 1, 15, 0, 0, 0, time.UTC)
	s := stringify(date)
	expected := "2023-03-01T15:00:00Z"
	if s != expected {
		t.Fatalf(`Expected "%s", received %s`, expected, s)
	}
}

func TestTrue(t *testing.T) {
	True(true, "This should pass")
	//True(false, "This should fail")
}
