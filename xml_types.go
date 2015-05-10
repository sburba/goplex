package goplex

import (
	"encoding/xml"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type CommaSeperatedSlice []string

func (s *CommaSeperatedSlice) UnmarshalXMLAttr(attr xml.Attr) error {
	*s = CommaSeperatedSlice(strings.Split(attr.Value, ","))
	return nil
}

type HTTPURL struct {
	url.URL
}

func (u *HTTPURL) UnmarshalXMLAttr(attr xml.Attr) error {
	if !strings.HasPrefix(attr.Value, "http://") {
		attr.Value = "http://" + attr.Value
	}

	url, err := url.Parse(attr.Value)
	*u = HTTPURL{*url}
	return err
}

type URLPath struct {
	url.URL
}

func (p *URLPath) UnmarshalXMLAttr(attr xml.Attr) error {
	url, err := url.Parse(attr.Value)
	*p = URLPath{*url}
	return err
}

type UnixTime struct {
	time.Time
}

func (t *UnixTime) UnmarshalXMLAttr(attr xml.Attr) error {
	sec, err := strconv.ParseInt(attr.Value, 10, 64)
	if err != nil {
		return err
	}

	*t = UnixTime{time.Unix(sec, 0)}
	return nil
}

type MillisDuration time.Duration

func (dur *MillisDuration) UnmarshalXMLAttr(attr xml.Attr) error {
	millis, err := strconv.ParseInt(attr.Value, 10, 64)
	if err != nil {
		return err
	}

	*dur = MillisDuration(time.Duration(millis) * time.Millisecond)
	return nil
}

type IntAsBool bool

func (v *IntAsBool) UnmarshalXMLAttr(attr xml.Attr) error {
	*v = attr.Value == "1"
	return nil
}
