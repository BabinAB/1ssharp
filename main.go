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
	configuration := loadConfuguration(* configureFile)

	fmt.Println(configuration)

	//TODO add check type by configuration mode
	//configuration.Language == 'php'
	var builder lib1ssarp_php.ConfigurationBuilderPHP
	builder.Build(configuration)
}


func loadConfuguration(configureFile string)lib1ssarp.Configuration {
	file, e := ioutil.ReadFile(configureFile)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	//fmt.Printf("%s\n", string(file))

	var configuration lib1ssarp.Configuration
	e = json.Unmarshal(file, &configuration)

	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	return configuration
}