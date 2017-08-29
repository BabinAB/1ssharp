package lib1ssarp_php

import (
	"./../../lib1ssarp"
	"io/ioutil"
	"os"
	"fmt"
    "text/template"
	"bytes"
)


type ConfigurationBuilderPHP struct {
	Config lib1ssarp.Configuration
}


func (cb ConfigurationBuilderPHP)  Build(c lib1ssarp.Configuration){

	cb.Config = c

	cb.BuildDataBase()
	cb.BuildFiles()
}




func (cb ConfigurationBuilderPHP)  BuildDataBase(){
	//TODO create data base

	fmt.Println("Create data base: not yet")
}


func (cb ConfigurationBuilderPHP)  BuildFiles(){

	file := cb.Config.Output + "/index.php"

	ProcessorTemplate("./lib1ssarp/php/templates/index.php", file, cb)
}


func (cb ConfigurationBuilderPHP) RunServer(c lib1ssarp.Configuration) {
	fmt.Println("Start Server ...")
	s := Server{c}
	s.Launch()
}


/**
Паристся шаблон и сохраняется в файл
TODO move to own package
 */
func ProcessorTemplate(sourceFile string, distFile string, context interface{}) {

	fmt.Println("Read File: ", sourceFile)

	t, e := template.ParseFiles(sourceFile)
	if e != nil {
		fmt.Printf("Error Parse Template : %v\n", e)
		os.Exit(1)
	}
	/*t.Funcs(template.FuncMap{
		"html": func(value interface{}) template.HTML {
			return template.HTML(fmt.Sprint(value))
		},
	})*/
	var tpl bytes.Buffer

	e = t.Execute(&tpl, context)
	if e != nil {
		fmt.Printf("Error Execute Template : %v\n", e)
		os.Exit(1)
	}

	fmt.Println("Output File: ", distFile)

	e = ioutil.WriteFile(distFile, tpl.Bytes(), 0644);
	if e != nil {
		fmt.Printf("Error Save File: %v\n", e)
		os.Exit(1)
	}
}