package util

import (
	"reflect"
	"regexp"

	"github.com/pkg/errors"
)

var (
	controlChars   = regexp.MustCompile("[[:cntrl:]]")
	controlEncoded = regexp.MustCompile("%[0-1][0-9,a-f,A-F]")
)

func ValidateURL(pathURL string) error {
	// Don't allow a URL containing control characters, standard or url-encoded
	if controlChars.FindStringIndex(pathURL) != nil || controlEncoded.FindStringIndex(pathURL) != nil {
		return errors.New("Invalid characters in url")
	}
	return nil
}

func StructToStrMap(i interface{}, len int) map[string]string {
	mapStr := make(map[string]string, len)
	iVal := reflect.ValueOf(i).Elem()
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		f := iVal.Field(i)
		mapStr[typ.Field(i).Name] = f.String()
	}
	return mapStr
}
