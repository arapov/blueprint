package roster

import (
	"fmt"
	"os"
)

type Connection interface {
	Query(filter string, attributes []string) ([]map[string][]string, error)
}

func GetGroups(ldapc Connection, group ...string) ([]map[string][]string, error) {
	// TODO: Find a better, generic place:
	groupPrefix := os.Getenv("LDAP_GROUPS_PREFIX")
	subGroupPrefix := os.Getenv("LDAP_SUBGROUPS_PREFIX")

	// TODO: Generalize. Handles just one specific group atm.
	if group[0] != "" {
		groupPrefix = group[0]
	}

	// objectClass rhatRoverGroup hardcoded due to app specific case, it's
	// unlikely and not meant to be used anywhere else
	filter := fmt.Sprintf("(&(objectClass=rhatRoverGroup)(cn=%s*)(!cn=*%s*))", groupPrefix, subGroupPrefix)
	attributes := []string{"cn", "description"}

	return ldapc.Query(filter, attributes)
}
