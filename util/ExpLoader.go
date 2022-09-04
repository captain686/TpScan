package util

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Exploit struct {
	Name       string           `yaml:"name"`
	Rules      map[string]Rules `yaml:"rules"`
	Expression string           `yaml:"expression"`
	Detail     detail           `yaml:"detail"`
}

type detail struct {
	Author string   `yaml:"author"`
	Links  []string `yaml:"links"`
}

type Rules struct {
	Request struct {
		Method          string            `yaml:"method"`
		Path            string            `yaml:"path"`
		Headers         map[string]string `yaml:"headers"`
		Body            string            `yaml:"body"`
		FollowRedirects string            `yaml:"follow_Redirects"`
		Proxy           string            `yaml:"proxy"`
	}
	Expression map[string]string `yaml:"expression"`
}

func ReadExp(file string) (*Exploit, error) {
	var exp Exploit

	File, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = yaml.Unmarshal(File, &exp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("[O] Load Exp:" + exp.Name)

	return &exp, nil
}
