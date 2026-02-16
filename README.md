# treemux

`treemux` is a small tmux session helper for root/child session trees. It provides commands to manage
tmux sessions in a hierarchical manner, allowing you to organize your sessions based on their root
directories.

## Commands

### --help

```
treemux --help

Usage: treemux <command>

Commands:
  attach    Attach to a tmux session, creating it if it doesn't exist.
  show      Show a list of treemux sessions.

Run "treemux <command> --help" for more information.
```

### attach-root

Attach to a root tmux session, creating it if it doesn't exist. Works from both inside and outside
of tmux sessions.

ARGS

`--name <root-name>`:                     Root session name. If not specified, it is derived from
                                          the session directory.

`--cmd <command>`:                        Command to run in the session. If not specified, the
                                          session will start with the default shell.

`-d --dir <directory>`:                   Starting directory for the session. If not specified, it
                                          uses the current directory.

`-w --worktree <worktree/branch>`:        If specified, the session will be rooted at the given git
                                          worktree directory. If the worktree did not exist for the
                                          specified branch, it will be created. If `-d` is not
                                          specified, it will be rooted at .worktrees/<branch>.
                                          If `-d` is specified, the worktree will be created at that
                                          directory.

### attach-child

Attach to a child tmux session rooted at the specified root session, creating it if it doesn't
exist. Works from both inside and outside of tmux sessions. Child sessions always use the root
session directory.

ARGS

`--root <root-name>`:                     Root session name. If omitted inside tmux, treemux
                                          defaults to the current session's root.

`--name <child-name>`:                    Child session name.

`--cmd <command>`:                        Command to run in the session. If not specified, the
                                          session will start with the default shell.

### show

Prints a fzf-friendly list of tmux sessions, with root sessions and their child sessions grouped
together.

NOTE: treemux only shows sessions it manages (sessions with treemux metadata).

Example output:

```
root-a
root-a  child-1
root-a  child-2
root-b
root-b  child-1
```

Example usage:

```
child_of_a=$(treemux show | awk '/^root-a/ {print $2}' | fzf)
```
