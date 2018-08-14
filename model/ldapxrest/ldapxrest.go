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

func Query(ldapc Connection, formValues map[string][]string) ([]map[string][]string, error) {
	var filter string
	var attributes []string

	attributes = nil
	filter = "(&"
	for key, values := range formValues {
		if key == "attributes" {
			attributes = strings.Split(values[0], ",")
			continue
		}

		value := strings.Split(values[0], ",")
		if len(value) > 1 {
			filter += "(|"
			for _, v := range value {
				filter += fmt.Sprintf("(%s=%s)", key, v)
			}
			filter += ")"
		} else {
			filter += fmt.Sprintf("(%s=%s)", key, values[0])
		}
	}
	filter += ")"

	res, err := ldapc.Query(filter, attributes)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return res, err
}

func Ping(ldapc Connection) error {
	return ldapc.Ping()
}
