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
	// Remember the priority is: env > json/yaml > default
	type Config struct {
		Name  string          `json:"name" yaml:"name"`
		Host  string          `json:"host" yaml:"host"`
		Ports []int           `json:"ports" yaml:"ports" env:"PORTS" default:"[10001, 10002, 10003]"`
		Used  map[string]bool `json:"used" yaml:"used" env:"USED" default:"{\"foo\": true, \"bar\": false}"`
		Auth  struct {
			Username string `json:"username" yaml:"username" env:"USERNAME" default:"xxx"`
			Password string `json:"password" yaml:"password" env:"PASSWORD" default:"123"`
		} `json:"auth" yaml:"auth"`
	}

	// Define JSON and YAML config content
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

	// Output:
	// &{Name:mydb Host:localhost Ports:[10004 10005 10006] Used:map[bar:false foo:true] Auth:{Username:zzz Password:456}}
	// &{Name:mydb Host:localhost Ports:[10004 10005 10006] Used:map[bar:false foo:true] Auth:{Username:zzz Password:456}}
}
