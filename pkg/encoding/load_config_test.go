// -------------------------------------------------------------------------------------------
// Copyright (c) Team Sorghum. All rights reserved.
// Licensed under the GPL v3 License. See LICENSE in the project root for license information.
// -------------------------------------------------------------------------------------------

package encoding_test

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/teamsorghum/go-common/pkg/encoding"
)

func TestLoadConfig_setVal(t *testing.T) {
	t.Parallel()

	type Embedded struct {
		InnerInt int `json:"inner_int"`
	}

	type Config struct {
		innerField string          `default:"nil"`
		String     string          `default:"string"`
		Bool       bool            `default:"true"`
		Int        int             `default:"-100"`
		Uint       uint            `default:"100"`
		Float64    float64         `default:"0.01"`
		Complex128 complex128      `default:"(1+1i)"`
		Slice      []int           `default:"[1, 2, 3]"`
		Array      [3]int          `default:"[1, 2, 3]"`
		Map        map[string]bool `default:"{\"foo\": true, \"bar\": false}"`
		Struct     Embedded        `default:"{\"inner_int\": -100}"`
	}

	want := Config{
		innerField: "",
		String:     "string",
		Bool:       true,
		Int:        -100,
		Uint:       100,
		Float64:    0.01,
		Complex128: complex(1, 1),
		Slice:      []int{1, 2, 3},
		Array:      [3]int{1, 2, 3},
		Map: map[string]bool{
			"foo": true,
			"bar": false,
		},
		Struct: Embedded{
			InnerInt: -100,
		},
	}

	got, err := encoding.LoadConfig[Config](nil, encoding.TypeNil)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, *got) {
		t.Fatalf("Want %+v\nGot %+v", want, *got)
	}
}

func TestLoadConfig_overrideWithEnvVars(t *testing.T) {
	t.Parallel()

	type Embedded struct {
		Num  int    `json:"num" env:"TEST_NUM"`
		Name string `json:"name"`
	}

	type Config struct {
		TestPriority        int       `json:"test_priority" env:"TEST_PRIORITY" default:"0"`
		TestEmbeddedPointer *Embedded `env:"TEST_EMBEDDED_POINTER"`
	}

	jsonContent := `{"test_priority": 1}`

	os.Setenv("TEST_PRIORITY", "2")
	os.Setenv("TEST_EMBEDDED_POINTER", `{"num": 1, "name": "foo"}`)
	os.Setenv("TEST_NUM", "2")

	want := Config{
		TestPriority: 2,
		TestEmbeddedPointer: &Embedded{
			Num:  2,
			Name: "foo",
		},
	}

	got, err := encoding.LoadConfig[Config]([]byte(jsonContent), encoding.TypeJSON)
	if err != nil {
		t.Fatal(err)
	}

	if want.TestPriority != got.TestPriority || !reflect.DeepEqual(want.TestEmbeddedPointer, got.TestEmbeddedPointer) {
		t.Fatalf("Want %+v\nGot %+v", want, *got)
	}
}

func TestLoadConfig_overrideWithDefaultValues(t *testing.T) {
	t.Parallel()

	type Embedded struct {
		Num  int    `json:"num" default:"2"`
		Name string `json:"name" default:"foo"`
	}

	type Config struct {
		TestPriority        int       `json:"test_priority" default:"1"`
		TestEmbeddedPointer *Embedded `json:"test_embedded_pointer" default:"{\"num\": 1}"`
	}

	jsonContent := `{"test_priority": 2}`

	want := Config{
		TestPriority: 2,
		TestEmbeddedPointer: &Embedded{
			Num:  2,
			Name: "foo",
		},
	}

	got, err := encoding.LoadConfig[Config]([]byte(jsonContent), encoding.TypeJSON)
	if err != nil {
		t.Fatal(err)
	}

	if want.TestPriority != got.TestPriority || !reflect.DeepEqual(want.TestEmbeddedPointer, got.TestEmbeddedPointer) {
		t.Fatalf("Want %+v\nGot %+v", want, *got)
	}
}

func TestLoadConfig_initNilPointers(t *testing.T) {
	t.Parallel()

	type Config struct {
		p1 *int
		P2 *struct {
			P3 *[3]int
			P4 *[]int
			P5 *map[int]int
		}
	}

	got, err := encoding.LoadConfig[Config](nil, encoding.TypeNil)
	if err != nil {
		t.Fatal(err)
	}

	if got.p1 != nil ||
		got.P2 == nil ||
		got.P2.P3 == nil ||
		got.P2.P4 == nil ||
		got.P2.P5 == nil {
		t.Fatalf("%+v", got)
	}

}

func TestLoadConfig_nonStructConfig(t *testing.T) {
	t.Parallel()

	_, err := encoding.LoadConfig[map[string]int]([]byte(`{"foo": 1}`), encoding.TypeJSON)
	if err == nil {
		t.Fatal("Expect non-nil error, got nil")
	}
}

func TestLoadConfig_invalidConfig(t *testing.T) {
	t.Parallel()

	type Config struct {
		Num int `json:"num" yaml:"num" toml:"num" xml:"num" default:"1"`
	}

	invalidContent := "invalid"

	_, err := encoding.LoadConfig[Config]([]byte(invalidContent), encoding.TypeJSON)
	if err == nil {
		t.Fatal("Expect err != nil, got nil")
	}
	_, err = encoding.LoadConfig[Config]([]byte(invalidContent), encoding.TypeYAML)
	if err == nil {
		t.Fatal("Expect err != nil, got nil")
	}
	_, err = encoding.LoadConfig[Config]([]byte(invalidContent), encoding.TypeTOML)
	if err == nil {
		t.Fatal("Expect err != nil, got nil")
	}
	_, err = encoding.LoadConfig[Config]([]byte(invalidContent), encoding.TypeXML)
	if err == nil {
		t.Fatal("Expect err != nil, got nil")
	}
}

func TestLoadConfig_types(t *testing.T) {
	t.Parallel()

	type Config struct {
		Num int `json:"num" yaml:"num" toml:"num" xml:"num" default:"1"`
	}

	jsonContent := `
{
	"num": 2
}
`

	yamlContent := `
num: 2
`

	tomlContent := `
num = 2
`

	xmlContent := `
<config>
  <num>
    2
  </num>
</config>
`

	jsonConfig, err := encoding.LoadConfig[Config]([]byte(jsonContent), encoding.TypeJSON)
	if err != nil {
		t.Fatal(err)
	}

	yamlConfig, err := encoding.LoadConfig[Config]([]byte(yamlContent), encoding.TypeYAML)
	if err != nil {
		t.Fatal(err)
	}

	tomlConfig, err := encoding.LoadConfig[Config]([]byte(tomlContent), encoding.TypeTOML)
	if err != nil {
		t.Fatal(err)
	}

	xmlConfig, err := encoding.LoadConfig[Config]([]byte(xmlContent), encoding.TypeXML)
	if err != nil {
		t.Fatal(err)
	}

	defaultConfig, err := encoding.LoadConfig[Config]([]byte(jsonContent), encoding.TypeNil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = encoding.LoadConfig[Config]([]byte(jsonContent), encoding.Type(-1))
	if !errors.Is(err, encoding.ErrLoadConfigUnsupportedType) {
		t.Fatalf("Want encoding.ErrLoadConfigUnsupportedType, got %+v", err)
	}

	wantConfig := Config{2}
	wantDefaultConfig := Config{1}

	if !reflect.DeepEqual(*jsonConfig, wantConfig) {
		t.Fatalf("Want %+v, got %+v", wantConfig, jsonConfig)
	}
	if !reflect.DeepEqual(*yamlConfig, wantConfig) {
		t.Fatalf("Want %+v, got %+v", wantConfig, yamlConfig)
	}
	if !reflect.DeepEqual(*tomlConfig, wantConfig) {
		t.Fatalf("Want %+v, got %+v", wantConfig, tomlConfig)
	}
	if !reflect.DeepEqual(*xmlConfig, wantConfig) {
		t.Fatalf("Want %+v, got %+v", wantConfig, xmlConfig)
	}
	if !reflect.DeepEqual(*defaultConfig, wantDefaultConfig) {
		t.Fatalf("Want %+v, got %+v", wantDefaultConfig, defaultConfig)
	}
}
