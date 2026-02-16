# treemux

`treemux` is a small tmux session helper for root/child session trees. It provides commands to manage
tmux sessions in a hierarchical manner, allowing you to organize your sessions based on their root
directories.

## Commands

### attach-root

Attach to a root tmux session, creating it if it doesn't exist. Works from both inside and outside
of tmux sessions.

ARGS

`--name <root-name>`:                     Root session name. If not specified, it is derived from
                                          the session directory.

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

### show-roots

Prints a list of treemux root sessions. The current root is prefixed with `*`. Use `--hide-current`
to omit the current root. Use `--sort-by=most-recently-used` to order by most recently attached.

NOTE: treemux only shows sessions it manages (sessions with treemux metadata).

Example output:

```
* root-a
  root-b
```

### show-children

Prints child sessions for a root. The current session is prefixed with `*`. If `--root` is omitted
inside tmux, treemux defaults to the current session's root. Use `--hide-current` to omit the
current session. Use `--sort-by=most-recently-used` to order by most recently attached.

Example output:

```
  child-1
  child-2
* child-3
```

## Integration tests

Integration tests require a tmux session. Run them from inside tmux:

```
TREEMUX_TMUX=1 go test ./...
```
