package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"strings"
)

const (
	defaultServicePort  = 55555
	defaultServiceAddr  = "127.0.0.1"
	defaultPageSize     = 20
	initBasketCapacity  = 200
	maxBasketCapacity   = 2000
	defaultDatabaseType = DbTypeMemory
	serviceOldAPIPath   = "baskets"
	serviceAPIPath      = "api"
	serviceUIPath       = "web"
	serviceName         = "request-baskets"
	basketNamePattern   = `^[\w\d\-_\.]{1,250}$`
	sourceCodeURL       = "https://github.com/darklynx/request-baskets"
	defaultAllowForward = false
)

// ServerConfig describes server configuration.
type ServerConfig struct {
	ServerPort   int
	ServerAddr   string
	InitCapacity int
	MaxCapacity  int
	PageSize     int
	MasterToken  string
	DbType       string
	DbFile       string
	DbConnection string
	Baskets      []string
	PathPrefix   string
	Mode         string
	Theme        string
	ThemeCSS     template.HTML
	AllowForward bool
}

type arrayFlags []string

func (v *arrayFlags) String() string {
	return strings.Join(*v, ",")
}

func (v *arrayFlags) Set(value string) error {
	*v = append(*v, value)
	return nil
}

// CreateConfig creates server configuration base on application command line arguments
func CreateConfig() *ServerConfig {
	var port = flag.Int("p", defaultServicePort, "HTTP service port")
	var address = flag.String("l", defaultServiceAddr, "HTTP listen address")
	var initCapacity = flag.Int("size", initBasketCapacity, "Initial basket size (capacity)")
	var maxCapacity = flag.Int("maxsize", maxBasketCapacity, "Maximum allowed basket size (max capacity)")
	var pageSize = flag.Int("page", defaultPageSize, "Default page size")
	var masterToken = flag.String("token", "", "Master token, random token is generated if not provided")
	var dbType = flag.String("db", defaultDatabaseType, fmt.Sprintf(
		"Baskets storage type: \"%s\" - in-memory, \"%s\" - Bolt DB, \"%s\" - SQL database",
		DbTypeMemory, DbTypeBolt, DbTypeSQL))
	var dbFile = flag.String("file", "./baskets.db", "Database location, only applicable for file or SQL databases")
	var dbConnection = flag.String("conn", "", "Database connection string for SQL databases, if undefined \"file\" argument is considered")
	var prefix = flag.String("prefix", "", "Service URL path prefix")
	var mode = flag.String("mode", ModePublic, fmt.Sprintf(
		"Service mode: \"%s\" - any visitor can create a new basket, \"%s\" - baskets creation requires master token",
		ModePublic, ModeRestricted))
	var theme = flag.String("theme", ThemeStandard, fmt.Sprintf(
		"CSS theme for web UI, supported values: %s, %s, %s",
		ThemeStandard, ThemeAdaptive, ThemeFlatly))
	var allowForward = flag.Bool("allowforward", defaultAllowForward, "Allow forwards for basket configs (default false)")

	var baskets arrayFlags
	flag.Var(&baskets, "basket", "Name of a basket to auto-create during service startup (can be specified multiple times)")
	flag.Parse()

	var token = *masterToken
	if len(token) == 0 {
		token, _ = GenerateToken()
		log.Printf("[info] generated master token: %s", token)
	}

	return &ServerConfig{
		ServerPort:   *port,
		ServerAddr:   *address,
		InitCapacity: *initCapacity,
		MaxCapacity:  *maxCapacity,
		PageSize:     *pageSize,
		MasterToken:  token,
		DbType:       *dbType,
		DbFile:       *dbFile,
		DbConnection: *dbConnection,
		Baskets:      baskets,
		PathPrefix:   normalizePrefix(*prefix),
		Mode:         *mode,
		Theme:        *theme,
		ThemeCSS:     toThemeCSS(*theme),
		AllowForward: *allowForward,
	}
}

func normalizePrefix(prefix string) string {
	if (len(prefix) > 0) && (prefix[0] != '/') {
		return "/" + prefix
	} else {
		return prefix
	}
}
