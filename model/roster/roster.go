package roster

import (
	"log"
)

type Connection interface {
	Query(filter string, attributes []string) ([]map[string][]string, error)
}

func GetGroups(ldapc Connection) ([]map[string][]string, error) {
	log.Println("getting groups...")

	// TODO: c.Query()

	return nil, nil
}
