package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type (
	rootCommand struct {
		Name     string     `json:"name"`
		Commands []*command `json:"commands"`
		Flags    []*cmdFlag `json:"flags"`
	}

	cmdFlag struct {
		Long        string `json:"long,omitempty"`
		Short       string `json:"short,omitempty"`
		Description string `json:"description,omitempty"`
		Argument    string `json:"argument,omitempty"`
		Required    bool   `json:"required"`
	}

	command struct {
		Name  string     `json:"name"`
		Flags []*cmdFlag `json:"flags,omitempty"`
	}
)

func main() {
	var (
		output string
		rc     rootCommand
	)
	{
		flag.StringVar(&output, "output", "goagen.fish", "Output `file`")
		flag.Parse()

		cmd := exec.Command("goagen", "commands")
		out, err := cmd.Output()
		if err != nil {
			fail("could not run goagen, check it's installed and in the PATH", err)
		}
		err = json.NewDecoder(bytes.NewBuffer(out)).Decode(&rc)
		if err != nil {
			fail("could not parse 'goagen commands' output", err)
		}
	}
	f, err := os.Create(output)
	if err != nil {
		fail("could not open output file", err)
	}
	generateScript(f, &rc)
	fmt.Println(output)
}

func generateScript(file *os.File, rc *rootCommand) {
	_, err := file.WriteString(strings.Replace(scriptPrefix, "__$$COMMAND$$__", rc.Name, 1))
	if err != nil {
		fail("could not write script", err)
	}
	for _, f := range rc.Flags {
		file.WriteString(flagComplete(rc.Name, "", f) + "\n")
	}
	file.WriteString("\n")
	for _, c := range rc.Commands {
		file.WriteString(cmdComplete(rc.Name, c))
	}
	file.Close()
}

func fail(msg string, err error) {
	if err != nil {
		msg += ": "
	}
	fmt.Fprintf(os.Stderr, "%s%s", msg, err)
	os.Exit(1)
}

func cmdComplete(cmdName string, cmd *command) string {
	entry := "# " + cmd.Name + "\ncomplete -f -c " + cmdName + " -n '__fish_prog_needs_command' -a " + cmd.Name + "\n"
	fls := make([]string, len(cmd.Flags))
	for i, fl := range cmd.Flags {
		fls[i] = flagComplete(cmdName, cmd.Name, fl)
	}
	return entry + strings.Join(fls, "\n") + "\n\n"
}

func flagComplete(rootCmdName, parCmdName string, fl *cmdFlag) string {
	line := "complete -f -c " + rootCmdName
	if parCmdName != "" {
		line += " -n '__fish_prog_using_command " + parCmdName + "'"
	}
	if fl.Long != "" {
		line += " -l " + fl.Long
	}
	if fl.Short != "" {
		line += " -s " + fl.Short
	}
	if fl.Required {
		line += " -r"
	}
	switch fl.Argument {
	case "$DIR":
		line += " -a '(__fish_complete_directories)'"
	case "$PKG":
		line += " -a '(__list_pkg)'"
	case "$DESIGN_PKG":
		line += " -a '(__list_design_pkg)'"
	}
	if fl.Description != "" {
		line += " -d '" + fl.Description + "'"
	}
	return line
}

const scriptPrefix = `#################################################################
# This script was auto-generated with goagenfish. DO NOT MODIFY #
#################################################################

# This is a fish shell completion script for the goagen tool
# See https://github.com/goadesign/goagenfish

function __escaped_go_path_src
    echo -n $GOPATH/src/ | sed -e 's/[\/&]/\\\&/g'
end

function __list_pkg
    set -l esc (__escaped_go_path_src)
    find $GOPATH/src -type d -not -path '*/.git/*' | sed -e s/^"$esc"//
end

function __list_design_pkg
    set -l esc (__escaped_go_path_src)
    find $GOPATH/src -type d -not -path '*/.git/*' | egrep "/design\$" | sed -e s/^"$esc"//
end

function __fish_prog_needs_command
  set cmd (commandline -opc)
  if [ (count $cmd) -eq 1 -a $cmd[1] = '__$$COMMAND$$__' ]
    return 0
  end
  return 1
end

function __fish_prog_using_command
  set cmd (commandline -opc)
  if [ (count $cmd) -gt 1 ]
    if [ $argv[1] = $cmd[2] ]
      return 0
    end
  end
  return 1
end

`
