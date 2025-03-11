// -------------------------------------------------------------------------------------------
// Copyright (c) Team Sorghum. All rights reserved.
// Licensed under the GPL v3 License. See LICENSE in the project root for license information.
// -------------------------------------------------------------------------------------------

// nolint:lll
package encoding_test

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/teamsorghum/go-common/pkg/encoding"
)

func ExampleLoadConfig() {
	// Define config struct. Pointer fields are supported.
	// Remember the priority is: env > json/yaml/toml/xml > default
	type Config struct {
		Name  string          `json:"name" yaml:"name" toml:"name" xml:"name"`
		Host  string          `json:"host" yaml:"host" toml:"host" xml:"host"`
		Ports []int           `json:"ports" yaml:"ports" toml:"ports" xml:"ports" env:"PORTS" default:"[10001, 10002, 10003]"`
		Used  map[string]bool `json:"used" yaml:"used" toml:"used" xml:"used" env:"USED" default:"{\"foo\": true, \"bar\": false}"`
		Auth  struct {
			Username string `json:"username" yaml:"username" toml:"username" xml:"username" env:"USERNAME" default:"xxx"`
			Password string `json:"password" yaml:"password" toml:"password" xml:"password" env:"PASSWORD" default:"123"`
		} `json:"auth" yaml:"auth" toml:"auth" xml:"auth"`
	}

	// Define JSON, YAML, TOML and XML config content
	jsonContent := `
{
	"name": "mydb",
	"host": "localhost",
	"auth": {
		"username": "yyy",
		"password": "456"
	}
}
`
	yamlContent := `
name: mydb
host: localhost
auth:
  username: yyy
  password: 456
`

	tomlContent := `
name = "mydb"
host = "localhost"

[auth]
username = "yyy"
password = "456"
`

	xmlContent := `
<config>
  <name>mydb</name>
  <host>localhost</host>
  <auth>
    <username>yyy</username>
    <password>456</password>
  </auth>
</config>
`

	// Set environment variables
	if err := errors.Join(
		os.Setenv("USERNAME", "zzz"),
		os.Setenv("PORTS", "[10004, 10005, 10006]"),
	); err != nil {
		log.Fatalln(err.Error())
	}

	// Load JSON config
	jsonConfig, err := encoding.LoadConfig[Config]([]byte(jsonContent), encoding.TypeJSON)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Printf("%+v\n", jsonConfig)

	// Load YAML config
	yamlConfig, err := encoding.LoadConfig[Config]([]byte(yamlContent), encoding.TypeYAML)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Printf("%+v\n", yamlConfig)

	// Load TOML config
	tomlConfig, err := encoding.LoadConfig[Config]([]byte(tomlContent), encoding.TypeTOML)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Printf("%+v\n", tomlConfig)

	// Load XML config
	xmlConfig, err := encoding.LoadConfig[Config]([]byte(xmlContent), encoding.TypeXML)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Printf("%+v\n", xmlConfig)

	// Output:
	// &{Name:mydb Host:localhost Ports:[10004 10005 10006] Used:map[bar:false foo:true] Auth:{Username:zzz Password:456}}
	// &{Name:mydb Host:localhost Ports:[10004 10005 10006] Used:map[bar:false foo:true] Auth:{Username:zzz Password:456}}
	// &{Name:mydb Host:localhost Ports:[10004 10005 10006] Used:map[bar:false foo:true] Auth:{Username:zzz Password:456}}
	// &{Name:mydb Host:localhost Ports:[10004 10005 10006] Used:map[bar:false foo:true] Auth:{Username:zzz Password:456}}
}
