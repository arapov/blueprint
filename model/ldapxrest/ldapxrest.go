package ldapxrest

import (
	"fmt"
	"log"
	"strings"
)

// Connection interface
type Connection interface {
	Query(filter string, attributes []string) ([]map[string][]string, error)
	Ping() error
}

func Query(ldapc Connection, formValues map[string][]string, filter *string) ([]map[string][]string, error) {
	var attributes []string

	var andKeyList = make(map[string][]string)
	var orKeyList = make(map[string][]string)
	for key, formValue := range formValues {
		var values []string
		for _, value := range formValue {
			values = append(values, strings.Split(value, ",")...)
		}

		if key == "attributes" {
			attributes = values
			continue
		}

		for _, value := range values {
			if key[0] == byte('|') {
				cleanKey := strings.TrimLeft(key, "|")
				orKeyList[cleanKey] = append(orKeyList[cleanKey], value)
			} else {
				andKeyList[key] = append(andKeyList[key], value)
			}
		}
	}

	var andFilter string
	for k, values := range andKeyList {
		for _, value := range values {
			andFilter += fmt.Sprintf("(%s=%s)", k, value)
		}
	}
	if andFilter != "" {
		*filter = fmt.Sprintf("(&%s)", andFilter)
	}

	var orFilter string
	for k, values := range orKeyList {
		for _, value := range values {
			if andFilter == "" {
				orFilter += fmt.Sprintf("(%s=%s)", k, value)
			} else if len(orKeyList)+len(values) == 2 {
				orFilter += fmt.Sprintf("(&%s)(%s=%s)", andFilter, k, value)
			} else {
				orFilter += fmt.Sprintf("(&%s(%s=%s))", andFilter, k, value)
			}
		}
	}
	if orFilter != "" {
		*filter = fmt.Sprintf("(|%s)", orFilter)
	}

	res, err := ldapc.Query(*filter, attributes)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return res, err
}

func Ping(ldapc Connection) error {
	return ldapc.Ping()
}
