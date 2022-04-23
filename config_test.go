package config_test

import (
	"flag"
	"os"

	"github.com/portcullis/config"
)

func ExampleBind() {
	// create just a simple struct with some descriptive flags for the configuration
	myConfig := struct {
		Name     string `description:"This is a name" flag:"name"`
		Password string `description:"Super secret password" mask:"true"`
		HTTP     struct {
			Addr string `name:"Address" description:"Address to listen" flag:"address"`
			Port int16  `description:"What port to listen" flag:"port"`
		}
		Enabled bool `description:"Enable something"`
	}{
		Name: "Default User",
	}

	// set values like normal
	myConfig.HTTP.Addr = "0.0.0.0"
	myConfig.HTTP.Port = 8080

	// bind the configuration under MyApplication to the pointer of the config
	config.Subset("MyApplication").Bind(&myConfig)

	// parsing the flags, would normally be replaced with os.Args[1:]
	flag.CommandLine.Parse([]string{"-name=flagged", "-address=127.0.0.1", "-port=8090"})

	// manually update a setting by full path (the value being set can come from os.GetEnv())
	config.Update("MyApplication.Enabled", "true")

	// dump the output
	config.Dump(os.Stdout)

	// Output:
	// Path                        Type        Value           Default Value      Description
	// MyApplication.Enabled       *bool       "true"          "false"            Enable something
	// MyApplication.HTTP.Addr     *string     "127.0.0.1"     "0.0.0.0"          Address to listen
	// MyApplication.HTTP.Port     *int16      "8090"          "8080"             What port to listen
	// MyApplication.Name          *string     "flagged"       "Default User"     This is a name
	// MyApplication.Password      *string     "*****"         "*****"            Super secret password
}
