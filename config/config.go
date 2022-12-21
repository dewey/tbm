package config

import (
	"strings"
	"text/template"
	"text/template/parse"
)

// Service is a single configuration option for a service we want to run
type Service struct {
	Command     string `yaml:"command"`
	Environment string `yaml:"environment"`
	Enable      bool   `yaml:"enable"`
	// Variables are string mappings, the key can be used as $KEY in the "Command" string. It will be interpolated when
	// it is used to spawn the proc
	Variables []map[string]string `yaml:"variables"`
}

// Configuration holds a configuration, the key of the map is the name of the configuration. This is a string defined by
// the user to differentiate the various services started.
type Configuration map[string]Service

// Valid validates a full configuration. This is mainly aiming at making sure we have unique port configurations.
func (s Configuration) Valid() bool {
	m := make(map[string]struct{})
	for _, service := range s {
		for _, variable := range service.Variables {
			_, ok := m[variable["port"]]
			if ok {
				return false
			} else {
				m[variable["port"]] = struct{}{}
			}
		}
	}
	return true
}

// InterpolatedCommand is replacing the variable placeholders in a string with the variable value
func (s Service) InterpolatedCommand() (string, error) {
	var finalCommand string
	tmpl, err := template.New("command").Parse(s.Command)
	if err != nil {
		return "", err
	}

	// Replace variables in command string if variables exist, otherwise we just return the original command
	variables := ListTemplateFields(tmpl)
	if len(s.Variables) > 0 {
		for _, val := range s.Variables {
			for _, val := range val {
				for _, variable := range variables {
					if strings.Contains(s.Command, variable) {
						finalCommand = strings.Replace(s.Command, variable, val, -1)
					}
				}
			}
		}
		return finalCommand, nil
	}
	return s.Command, nil
}

// Valid returns true if a service is enabled and has all the required values set
func (s Service) Valid() bool {
	// Fail early if the service is not enabled
	if !s.Enable {
		return false
	}

	vars, err := extractVariables(s.Command)
	if err != nil {
		return false
	}

	// Fail early if different counts
	if len(vars) != len(s.Variables) {
		return false
	}

	vm := make(map[string]struct{})
	for _, v := range vars {
		if _, ok := vm[v]; !ok {
			vm[v] = struct{}{}
		}
	}
	for _, variable := range s.Variables {
		for key, _ := range variable {
			if val, ok := vm[key]; ok {
				if vm[key] != val {
					return false
				}
			} else {
				return false
			}
		}
	}

	return true
}

// extractVariables parses a command template and returns the Go template variables that were used
func extractVariables(command string) ([]string, error) {
	tmpl, err := template.New("command").Parse(command)
	if err != nil {
		return nil, err
	}
	variables := ListTemplateFields(tmpl)
	for i, _ := range variables {
		variables[i] = strings.Replace(variables[i], "{{", "", -1)
		variables[i] = strings.Replace(variables[i], "}}", "", -1)
		variables[i] = strings.Replace(variables[i], ".", "", -1)
		variables[i] = strings.ToLower(variables[i])
	}
	return variables, nil
}

// ListTemplateFields lists the fields used in a template. Sourced and adapted from: https://stackoverflow.com/a/40584967
func ListTemplateFields(t *template.Template) []string {
	return listNodeFields(t.Tree.Root, nil)
}

// listNodeFields iterates over the parsed tree and extracts fields
func listNodeFields(node parse.Node, res []string) []string {
	//fmt.Println("p", node.String())
	//fmt.Println("p", node.Type())
	// Only looking at fields, needs to be adapted if further template entities should be supported
	//if node.Type() == parse.NodeField {
	//	res = append(res, node.String())
	//}

	if node.Type() == parse.NodeAction {
		res = append(res, node.String())
	}

	if ln, ok := node.(*parse.ListNode); ok {
		for _, n := range ln.Nodes {
			res = listNodeFields(n, res)
		}
	}
	return res
}
