package roster

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type Connection interface {
	Query(filter string, attributes []string) ([]map[string][]string, error)
}

func GetGroups(ldapc Connection, groups string, filter *string) ([]map[string][]string, error) {
	// TODO: Find a better, generic place
	groupPrefix := os.Getenv("LDAP_GROUPS_PREFIX")
	subGroupPrefix := os.Getenv("LDAP_SUBGROUPS_PREFIX")
	rolesPrefix := os.Getenv("LDAP_GROUPS_ROLES_PREFIX")

	if groupPrefix == "" {
		return nil, errors.New("LDAP_GROUP_PREFIX variable is not defined.")
	}

	// Get all custom roles
	var ldapRoles []map[string][]string
	var ldapMembersRoles []map[string][]string
	if rolesPrefix != "" {
		*filter = fmt.Sprintf("(&(objectClass=rhatRoverGroup)(cn=%s*))", rolesPrefix)
		ldapRoles, _ = ldapc.Query(*filter, []string{"description"})

		// Get members data that belong to custom roles
		*filter = fmt.Sprintf("(&(objectClass=rhatPerson)(memberOf=*%s*))", rolesPrefix)
		ldapMembersRoles, _ = ldapc.Query(*filter, []string{"uid", "cn", "memberOf"})
	}

	// Building up the map with members and its roles
	// e.g. mapMemberRole["uid=?,ou=?,ou=?"] = [...]["jdoe", "John Doe", "Fancy Role Name"]
	mapMemberRole := make(map[string][][]string)
	var roles []string // used later to fill in the groups info
	for _, ldapRole := range ldapRoles {
		roledn := ldapRole["dn"][0]            // e.g. "cn=?,ou=?,ou=?"
		rolename := ldapRole["description"][0] // e.g. "Fancy Role Name"

		roles = append(roles, rolename)

		for _, member := range ldapMembersRoles {
			if contains(member["memberOf"], roledn) {
				memberdn := member["dn"][0] // e.g. "uid=?,ou=?,ou=?"
				memberuid := member["uid"][0]
				membername := member["cn"][0]

				info := []string{memberuid, membername, rolename}

				mapMemberRole[memberdn] = append(mapMemberRole[memberdn], info)
			}
		}
	}

	// Get groups
	*filter = fmt.Sprintf("(&(objectClass=rhatRoverGroup)(cn=%s*))", groupPrefix)
	if subGroupPrefix != "" {
		*filter = fmt.Sprintf("(&(objectClass=rhatRoverGroup)(!cn=*%s*)(cn=%s*))", subGroupPrefix, groupPrefix)
	}
	// Amend the default query in case we have groups defined
	if groups != "" {
		groupPrefix = ""
		for _, group := range strings.Split(groups, ",") {
			groupPrefix += fmt.Sprintf("(cn=%s)", group)
		}
		*filter = fmt.Sprintf("(&(objectClass=rhatRoverGroup)(|%s))", groupPrefix)
		if subGroupPrefix != "" {
			*filter = fmt.Sprintf("(&(objectClass=rhatRoverGroup)(!cn=*%s*)(|%s))", subGroupPrefix, groupPrefix)
		}
	}
	ldapGroups, err := ldapc.Query(*filter, []string{"cn", "description", "rhatGroupNotes", "uniqueMember"})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Building up a Group data
	for _, ldapGroup := range ldapGroups {
		var uniqueMembers []string
		uniqueMembers = append(uniqueMembers, ldapGroup["uniqueMember"]...)

		// Each group must have consistent structure, even if the its content is null
		ldapGroup["subGroup"] = []string{}
		for _, role := range roles {
			ldapGroup[role] = []string{}
		}
		ldapGroup["links"] = []string{}
		ldapGroup["roles"] = roles

		// Get possible sub-groups, and merge sub-group info and its members into parent group
		if subGroupPrefix != "" {
			filter := fmt.Sprintf("(&(objectClass=rhatRoverGroup)(cn=%s%s*))", ldapGroup["cn"][0], subGroupPrefix)
			ldapSubGroups, _ := ldapc.Query(filter, []string{"cn", "description", "rhatGroupNotes", "uniqueMember"})
			for _, ldapSubGroup := range ldapSubGroups {
				// Extending Group data with the information about subGroups
				subGroup := fmt.Sprintf("%s,%s", ldapSubGroup["cn"][0], ldapSubGroup["description"][0])
				ldapGroup["subGroup"] = append(ldapGroup["subGroup"], subGroup)

				ldapGroup["rhatGroupNotes"] = append(ldapGroup["rhatGroupNotes"], ldapSubGroup["rhatGroupNotes"]...)
				uniqueMembers = append(uniqueMembers, ldapSubGroup["uniqueMember"]...)
			}
			removeDuplicates(&uniqueMembers)
			ldapGroup["uniqueMember"] = uniqueMembers
		}

		for _, note := range ldapGroup["rhatGroupNotes"] {
			ldapGroup["links"] = append(ldapGroup["links"], decodeNote(note)...)
		}
		delete(ldapGroup, "rhatGroupNotes")

		// Extend Group with the special roles and members
		for _, uniqueMember := range uniqueMembers {
			if _, ok := mapMemberRole[uniqueMember]; !ok {
				continue
			}

			for _, member := range mapMemberRole[uniqueMember] {
				name := fmt.Sprintf("%s,%s", member[0], member[1])
				role := member[2]

				ldapGroup[role] = append(ldapGroup[role], name)
			}
		}
	}

	return ldapGroups, err
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func decodeNote(note string) []string {
	var result []string

	// accepts:
	// pile:key=value or pile:key="value value"
	// can be separated by , or space
	re, _ := regexp.Compile(`pile:(\w*=[\w:/@.-]+|\w*="[\w\s!,:/@.-]+")`)
	// TODO: take care of error here
	pile := re.FindAllStringSubmatch(note, -1)
	// TODO: code below is fragile, very fragile
	for i := range pile {
		kv := strings.Split(pile[i][1], "=")
		note := fmt.Sprintf("%s,%s", strings.Title(kv[0]), strings.Trim(kv[1], "\""))
		result = append(result, note)
	}

	return result
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
