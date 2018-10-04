package roster

import (
	"fmt"
	"os"
)

type Connection interface {
	Query(filter string, attributes []string) ([]map[string][]string, error)
}

func GetGroups(ldapc Connection, group ...string) ([]map[string][]string, error) {
	prefix := os.Getenv("LDAP_GROUPS_PREFIX") // TODO: Find a better, generic place
	if group[0] != "" {
		prefix = group[0]
	}

	// objectClass rhatRoverGroup hardcoded due to app specific case, it's
	// unlikely to be used anywhere else
	filter := fmt.Sprintf("(&(objectClass=rhatRoverGroup)(cn=%s*))", prefix)
	attributes := []string{"cn", "description"}

	return ldapc.Query(filter, attributes)
}
