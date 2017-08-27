package lib1ssarp_php

import (
	"./../../lib1ssarp"
	"io/ioutil"
	"os"
	"fmt"
)


type ConfigurationBuilderPHP struct {

}



func (cb ConfigurationBuilderPHP)  Build(c lib1ssarp.Configuration){
	file := c.Output + "/index.php"
	fmt.Println("Output File: ", file)

	data := "<?php\necho \"Hello world!\";"
	e := ioutil.WriteFile(file, []byte(data), 0644);
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
}