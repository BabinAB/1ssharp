package main

import (
	"fmt"
	"flag"
	"os"
	"io/ioutil"
	"encoding/json"
	"./lib1ssarp"
	"./lib1ssarp/php"
)

func main()  {

	configureFile := flag.String("config", "configuration.test.json", "File configuration")

	flag.Parse()

	fmt.Println("Configure File: ", *configureFile)
	configuration := loadConfiguration(* configureFile)

	fmt.Println(configuration)


	var builder lib1ssarp.ConfigurationBuilder
	switch configuration.Language {
		case "php":
			builder = new(lib1ssarp_php.ConfigurationBuilderPHP)

		default:
			fmt.Printf("Configuration for language: %s not found\n", configuration.Language)
			os.Exit(1)
	}

	builder.Build(configuration)
}


func loadConfiguration(configureFile string)lib1ssarp.Configuration {
	file, e := ioutil.ReadFile(configureFile)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	fmt.Printf("Load Configure File: %s\n", configureFile)

	var configuration lib1ssarp.Configuration
	e = json.Unmarshal(file, &configuration)

	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	return configuration
}