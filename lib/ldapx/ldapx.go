// Package ldapx - Lorem ipsum...
package ldapx

import (
	"errors"
	"fmt"
	"log"

	ldap "gopkg.in/ldap.v2"
)

// Info keeps infromation about LDAP server
type Info struct {
	Hostname string
	Port     string
	BaseDN   string
}

// Conn represents an LDAP Connection
type Conn struct {
	*ldap.Conn
	BaseDN string
}

// Dial connects to the LDAP server and then returns a new Conn for the connection.
func (c Info) Dial() (*Conn, error) {
	parentConn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", c.Hostname, c.Port))
	if err != nil {
		return nil, err
	}

	return &Conn{Conn: parentConn, BaseDN: c.BaseDN}, err
}

// Query performs a LDAP search request, where filter is the search string and
// must obey LDAP filter syntax and attributes are the fields we want to get
// from LDAP. It returns all the fields if attributes is not set and nil.
func (c *Conn) Query(filter string, attributes []string) ([]map[string][]string, error) {
	if c == nil {
		return nil, errors.New("roster was started without LDAP setup, LDAP configuration must be provided")
	}

	request := ldap.NewSearchRequest(
		c.BaseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter, attributes, nil,
	)
	ldapRes, err := c.Search(request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var res []map[string][]string
	for _, entry := range ldapRes.Entries {
		ldapEntry := make(map[string][]string)
		for _, attr := range entry.Attributes {
			ldapEntry[attr.Name] = attr.Values
		}

		ldapEntry["dn"] = append(ldapEntry["dn"], entry.DN)

		res = append(res, ldapEntry)
	}

	return res, err
}

// Ping is ensuring we are connected to LDAP server and
// able to query. It is also used to keep connection alive.
func (c *Conn) Ping() error {
	if c == nil {
		return errors.New("LDAP connection is dead")
	}

	request := ldap.NewSearchRequest(
		c.BaseDN, ldap.ScopeSingleLevel, ldap.NeverDerefAliases, 0, 0, false,
		"(ou=*)", nil, nil,
	)
	_, err := c.Search(request)
	if err != nil {
		log.Println(err)
		return err
	}

	return err
}
