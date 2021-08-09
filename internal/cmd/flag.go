package cmd

import (
	"strconv"
	"strings"
)

// PIN flag type
type pinValue int

func newPinValue(val int, p *int) *pinValue {
	*p = val
	return (*pinValue)(p)
}

func (r *pinValue) String() string {
	if *r < 0 {
		return "none"
	}
	return strconv.Itoa(int(*r))
}

func (r *pinValue) Set(s string) error {
	v, err := strconv.Atoi(s)
	*r = pinValue(v)
	return err
}

func (r *pinValue) Type() string {
	return "int"
}

// case insensitive domain name
type domainValue string

func newDomainValue(val string, p *string) *domainValue {
	*p = val
	return (*domainValue)(p)
}

func (r *domainValue) String() string {
	return string(*r)
}

func (r *domainValue) Set(v string) error {
	s := strings.ToLower(strings.TrimSpace(v))
	*r = domainValue(s)
	return nil
}

func (r *domainValue) Type() string {
	return "string"
}
