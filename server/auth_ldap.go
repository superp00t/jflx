package server

import (
	"crypto/subtle"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-ldap/ldap/v3"
)

type ldap_cached_credential struct {
	Time     time.Time
	Password string
	Address  string
}

type ldap_cached_auth_provider struct {
	server             *Server
	mutex              sync.Mutex
	cached_credentials map[string]*ldap_cached_credential
}

func (p *ldap_cached_auth_provider) fetch_cached_credential(address, usergroup, username, password string) (ok bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var cached_credential *ldap_cached_credential
	cached_credential, ok = p.cached_credentials[username]
	if !ok {
		return
	}

	if subtle.ConstantTimeCompare([]byte(password), []byte(cached_credential.Password)) != 1 {
		ok = false
	}

	age := time.Since(cached_credential.Time)
	if age > p.server.config.LDAP.CacheExpiry {
		delete(p.cached_credentials, username)
		ok = false
	}

	return
}

func (p *ldap_cached_auth_provider) cache_credentials(address, usergroup, username, password string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	cached_credential, ok := p.cached_credentials[username]
	if !ok {
		cached_credential = new(ldap_cached_credential)
		p.cached_credentials[username] = cached_credential
	}
	cached_credential.Address = address
	cached_credential.Password = password
}

func (p *ldap_cached_auth_provider) AuthenticateCredentials(address, usergroup, username, password string) (ok bool) {
	ok = p.fetch_cached_credential(address, usergroup, username, password)
	if ok {
		return
	}

	conn, err := ldap.DialURL(p.server.config.LDAP.URL)
	if err != nil {
		log.Println("error contacting LDAP host:", err)
		return
	}

	if err := conn.Bind(p.server.config.LDAP.Username, p.server.config.LDAP.Password); err != nil {
		log.Println("error binding to LDAP host:", err)
		return
	}

	search_request := ldap.NewSearchRequest(
		p.server.config.LDAP.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=person)(uid=%s)(&(memberOf=cn=jflx_user,ou=groups,%s)))", ldap.EscapeFilter(username), p.server.config.LDAP.BaseDN),
		[]string{"dn"},
		nil,
	)

	search_result, err := conn.Search(search_request)
	if err != nil {
		log.Println("error searching for user:", err)
		return
	}

	if len(search_result.Entries) != 1 {
		log.Println("user not found:", username)
		return
	}

	user_entry := search_result.Entries[0]

	if err := conn.Bind(user_entry.DN, password); err != nil {
		log.Println("failed to authenticate user:", user_entry.DN)
	} else {
		ok = true
	}

	conn.Close()

	if ok {
		p.cache_credentials(address, usergroup, username, password)
	}

	return
}

func new_ldap_cached_auth_provider(s *Server) auth_provider {
	p := new(ldap_cached_auth_provider)
	p.cached_credentials = make(map[string]*ldap_cached_credential)
	p.server = s
	return p
}
