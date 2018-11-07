package roster

import (
	"fmt"
	"os"
	"strings"
)

type Connection interface {
	Query(filter string, attributes []string) ([]map[string][]string, error)
}

func GetGroups(ldapc Connection, groups string, filter *string) ([]map[string][]string, error) {
	// TODO: Find a better, generic place:
	subGroupPrefix := os.Getenv("LDAP_SUBGROUPS_PREFIX")
	groupPrefix := os.Getenv("LDAP_GROUPS_PREFIX")

	// objectClass rhatRoverGroup hardcoded due to app specific case, it's
	// unlikely and not meant to be used anywhere else
	*filter = fmt.Sprintf("(&(objectClass=rhatRoverGroup)(!cn=*%s*)(cn=%s*))", subGroupPrefix, groupPrefix)
	attributes := []string{"cn", "description", "uniqueMember"}

	if groups != "" {
		groupPrefix = ""
		for _, group := range strings.Split(groups, ",") {
			groupPrefix += fmt.Sprintf("(cn=%s)", group)
		}
		*filter = fmt.Sprintf("(&(objectClass=rhatRoverGroup)(!cn=*%s*)(|%s))", subGroupPrefix, groupPrefix)
	}

	return ldapc.Query(*filter, attributes)
}
