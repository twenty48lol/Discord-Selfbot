package commands

import (
	"fmt"
	"strings"
)

type CommandList struct {
	Commands []Command
}

func InitCommands() CommandList {
	list := emptyList()
	return list
}

func emptyList() CommandList {
	return CommandList{}
}

func (list *CommandList) AddCommands(commands ...Command) {
	for _, command := range commands {
		for _, c := range list.Commands {
			found := false
			var conflicting Command
			for one := range c.Names {
				if found {
					break
				}
				for two := range command.Names {
					if strings.EqualFold(c.Names[one], command.Names[two]) {
						found = true
						conflicting = c
					}
				}
			}
			if found {
				panic(fmt.Sprintf("error adding %s cuz it conflicts with existing command %s", command.Names[0], conflicting.Names[0]))
			}
		}
	}
}
