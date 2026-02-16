// Package cli defines the command-line interface for treemux.
package cli

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/ian-howell/treemux/internal/treemux"
)

// CLI defines the treemux command tree.
type CLI struct {
	// AttachRoot holds the attach-root subcommand configuration.
	AttachRoot AttachRootCmd `cmd:"" name:"attach-root" help:"Attach to a root session, creating it if it doesn't exist."`
	// AttachChild holds the attach-child subcommand configuration.
	AttachChild AttachChildCmd `cmd:"" name:"attach-child" help:"Attach to a child session rooted at the specified root."`
	// Show holds the show subcommand configuration.
	Show ShowCmd `cmd:"" help:"Show a list of treemux sessions."`
}

// AttachRootCmd configures the attach-root subcommand.
type AttachRootCmd struct {
	// Name is the name for the root session.
	Name string `name:"name" help:"Root session name."`
	// Command is the command to run in the session.
	Command string `name:"cmd" help:"Command to run in the session."`
	// Dir sets the session start directory.
	Dir string `name:"dir" short:"d" help:"Starting directory for the session."`
	// Worktree selects a git worktree branch.
	Worktree string `name:"worktree" short:"w" help:"Git worktree branch to use as the session root."`
}

// AttachChildCmd configures the attach-child subcommand.
type AttachChildCmd struct {
	// Root selects an existing root session.
	Root string `name:"root" required:"" help:"Root session name."`
	// Name is the name for the child session.
	Name string `name:"name" required:"" help:"Child session name."`
	// Command is the command to run in the session.
	Command string `name:"cmd" help:"Command to run in the session."`
}

// ShowCmd configures the show subcommand.
type ShowCmd struct{}

// Run parses CLI args and executes the requested command.
func Run() error {
	app := treemux.New()
	cli := CLI{}
	ctx := kong.Parse(&cli, kong.Name("treemux"))

	switch ctx.Command() {
	case "attach-root":
		return app.AttachRoot(treemux.AttachRootRequest{
			Name:     cli.AttachRoot.Name,
			Command:  cli.AttachRoot.Command,
			Dir:      cli.AttachRoot.Dir,
			Worktree: cli.AttachRoot.Worktree,
		})
	case "attach-child":
		return app.AttachChild(treemux.AttachChildRequest{
			Root:    cli.AttachChild.Root,
			Name:    cli.AttachChild.Name,
			Command: cli.AttachChild.Command,
		})
	case "show":
		lines, err := app.Show()
		if err != nil {
			return err
		}
		for _, line := range lines {
			fmt.Println(line)
		}
		return nil
	default:
		return fmt.Errorf("unknown command: %s", ctx.Command())
	}
}
