package main

import (
	"flag"
	"fmt"
	"os"
	"io"
	"os/exec"
	"strings"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		i int
		b string
		c string
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.IntVar(&i, "i", 0, "Issue")
	flags.StringVar(&b, "b", "", "Compare branch name")
	flags.StringVar(&c, "c", "", "New create compare branch name")

	flVersion := flags.Bool("version", false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if *flVersion {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	var (
		cmd string
		out []byte
		err error
		message string
		issueNumber int
		issueComment string
		branch string
	)

	// issueNumber & issueComment
	if i == 0 {
		fmt.Fprintln(cli.errStream, "required issue number")
		return ExitCodeError
	}
	// issue comment
	pattern := "s@([0-9]+\\]) (.*) \\( https://github\\.com.+issues/[0-9]+ .*$@\\2@";
	cmd = fmt.Sprintf("hub issue | grep \"^\\s*%d\" | sed -E '%s'", i, pattern);
	out, err = exec.Command(os.Getenv("SHELL"), "-c", cmd).Output()
	if err != nil {
		return ExitCodeError
	}
	if len(out) == 0 {
		fmt.Fprintln(cli.errStream, "not found issue")
		return ExitCodeError
	}
	issueNumber = i
	issueComment = strings.TrimSpace(strings.TrimRight(fmt.Sprintf("%s", out), "\n"))

	// create branch
	if c != "" {
		// create branch
		cmd = fmt.Sprintf("git checkout -b  %s origin/master", c);
		_, err = exec.Command(os.Getenv("SHELL"), "-c", cmd).Output()
		if err != nil {
			return ExitCodeError
		}
		// create commit
		cmd = "git commit -m 'empty' --allow-empty";
		_, err = exec.Command(os.Getenv("SHELL"), "-c", cmd).Output()
		if err != nil {
			return ExitCodeError
		}
		// push
		cmd = fmt.Sprintf("git push origin %s", c);
		_, err = exec.Command(os.Getenv("SHELL"), "-c", cmd).Output()
		if err != nil {
			return ExitCodeError
		}
	} else if b != ""{
		// checkout
		cmd = fmt.Sprintf("git checkout %s", b);
		_, err = exec.Command(os.Getenv("SHELL"), "-c", cmd).Output()
		if err != nil {
			return ExitCodeError
		}
	} else {
		fmt.Fprintln(cli.errStream, "not define branch")
		return ExitCodeError
	}

	// branch
	cmd = "git rev-parse --abbrev-ref HEAD";
	out, err = exec.Command(os.Getenv("SHELL"), "-c", cmd).Output()
	if err != nil {
		return ExitCodeError
	}
	branch = strings.TrimRight(fmt.Sprintf("%s", out), "\n")
	if (branch == "master") {
		fmt.Fprintln(cli.errStream, "branch is master")
		return ExitCodeError
	}

	// message
	message = fmt.Sprintf("wip %s (closes #%d)\nref. #%d", issueComment, issueNumber, issueNumber)
	cmd = fmt.Sprintf("hub pull-request -m '%s' --browse -h %s", message, branch);
	_ , err = exec.Command(os.Getenv("SHELL"), "-c", cmd).Output()
	if err != nil {
		fmt.Fprintln(cli.errStream, err)
		return ExitCodeError
	}

	return ExitCodeOK
}
