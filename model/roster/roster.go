package roster

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Connection interface {
	Query(filter string, attributes []string) ([]map[string][]string, error)
}

func GetGroups(ldapc Connection, groups string, filter *string) ([]map[string][]string, error) {

	attributes := []string{"cn", "description", "uniqueMember"}

	// TODO: Find a better, generic place
	subGroupPrefix := os.Getenv("LDAP_SUBGROUPS_PREFIX")
	groupPrefix := os.Getenv("LDAP_GROUPS_PREFIX")
	// objectClass rhatRoverGroup hardcoded due to app specific case, it's
	// unlikely and not meant to be used anywhere else
	*filter = fmt.Sprintf("(&(objectClass=rhatRoverGroup)(!cn=*%s*)(cn=%s*))", subGroupPrefix, groupPrefix)

	// Amend the default query in case we have groups defined
	if groups != "" {
		groupPrefix = ""
		for _, group := range strings.Split(groups, ",") {
			groupPrefix += fmt.Sprintf("(cn=%s)", group)
		}
		*filter = fmt.Sprintf("(&(objectClass=rhatRoverGroup)(!cn=*%s*)(|%s))", subGroupPrefix, groupPrefix)
	}

	ldapGroups, err := ldapc.Query(*filter, attributes)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Look for subgroups and merge stats
	var res []map[string][]string
	for _, ldapGroup := range ldapGroups {
		var uniqueMember []string

		uniqueMember = append(uniqueMember, ldapGroup["uniqueMember"]...)
		ldapGroup["subGroup"] = []string{}

		filter := fmt.Sprintf("(&(objectClass=rhatRoverGroup)(cn=%s%s*))", ldapGroup["cn"][0], subGroupPrefix)
		ldapSubGroups, _ := ldapc.Query(filter, attributes)
		for _, ldapSubGroup := range ldapSubGroups {
			// Extending Group data with the information about subGroups
			ldapGroup["subGroup"] = append(ldapGroup["subGroup"], ldapSubGroup["cn"][0])

			// Merging subGroup members with group members
			uniqueMember = append(uniqueMember, ldapSubGroup["uniqueMember"]...)
		}

		removeDuplicates(&uniqueMember)
		ldapGroup["uniqueMember"] = uniqueMember

		res = append(res, ldapGroup)
	}

	return res, err
}

func removeDuplicates(xs *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *xs {
		if !found[x] {
			found[x] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}
