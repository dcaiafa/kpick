package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
)

type Config struct {
	Clusters       []*Cluster `json:"clusters"`
	Users          []*User    `json:"users"`
	Contexts       []*Context
	CurrentContext string `json:"current-context"`
}

type Context struct {
	Name    string `json:"name"`
	Context struct {
		Cluster string `json:"cluster"`
		User    string `json:"user"`
	}
}

type Cluster struct {
	Name string `json:"name"`
}

type User struct {
	Name string `json:"name"`
}

func getConfig() *Config {
	cmd := exec.Command("kubectl", "config", "view", "-o", "json")
	res, err := cmd.CombinedOutput()
	if err != nil {
		die(err)
	}
	config := &Config{}
	err = json.Unmarshal(res, config)
	if err != nil {
		die(err)
	}
	return config
}

func useContext(c string) {
	cmd := exec.Command("kubectl", "config", "use-context", c)
	err := cmd.Run()
	if err != nil {
		die(err)
	}
}

func deleteContext(contextName string) {
	config := getConfig()

	var context *Context
	for i := 0; i < len(config.Contexts); i++ {
		if config.Contexts[i].Name == contextName {
			context = config.Contexts[i]
			break
		}
	}

	if context == nil {
		die(errors.New("context not found"))
	}

	err := exec.Command("kubectl", "config", "delete-context", contextName).Run()
	if err != nil {
		die(err)
	}

	userName := context.Context.User
	clusterName := context.Context.Cluster

	deleteUser := true
	for _, context = range config.Contexts {
		if context.Context.User == userName && context.Name != contextName {
			// User is used in another context.
			deleteUser = false
			break
		}
	}

	if deleteUser {
		exec.Command(
			"kubectl", "config", "unset", fmt.Sprintf("users.%s", userName)).Run()
	}

	deleteCluster := true
	for _, context = range config.Contexts {
		if context.Context.Cluster == clusterName && context.Name != contextName {
			// Cluster is used in another context.
			deleteCluster = false
			break
		}
	}

	if deleteCluster {
		exec.Command(
			"kubectl", "config", "delete-cluster", clusterName).Run()
	}
}

func renameContext(originalName, newName string) {
	exec.Command(
		"kubectl", "config", "rename-context", originalName, newName).Run()
}
