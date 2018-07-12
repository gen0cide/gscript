package testlib

import (
	"errors"
	"net/url"
	"strings"
)

// Aspy is a test type for me to work on the linker
type Aspy struct {
	Name string
}

// (s string, a2 map[*url.URL][]*Aspy, b []*Aspy, c url.URL, d *Aspy, e *url.URL) (*url.URL, error) {

// Test1  is a test function
func Test1(s string) (*url.URL, error) {
	u, err := url.Parse(s)
	_ = u
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Test2 is a test function
func Test2(p1, p2 string) string {
	return strings.Join([]string{p1, p2}, " ")
}

// Test3 is a test function
func Test3(s string) bool {
	if len(s) > 4 {
		return true
	}
	return false
}

// Test5 is a test function
func Test5() *Aspy {
	return &Aspy{
		Name: "simmons",
	}
}

// Test4 is a test function
func Test4(a *Aspy) string {
	return strings.ToUpper(a.Name)
}

// Test6 is a test function
func (a *Aspy) Test6(s string) (*Aspy, error) {
	if len(s) > 5 {
		return nil, errors.New("string should be 1-4 chars")
	}
	newName := strings.Join([]string{a.Name, s}, "-")
	return &Aspy{Name: newName}, nil
}

// Test7 is a test function
func Test7(s string) Aspy {
	return Aspy{
		Name: s,
	}
}

// Test8 is a test function
func (a Aspy) Test8() string {
	return a.Name
}
