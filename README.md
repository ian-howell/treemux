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

`-d --dir <directory>`:                   Starting directory for the session. Mutually exclusive
                                          with `-w --worktree`. If not specified, it uses the
                                          current directory.

`-w --worktree <worktree/branch>`:        Mutually exclusive with `-d --dir`. If specified, the
                                          session will be rooted at the given git worktree
                                          directory. If the worktree did not exist for the
                                          specified branch, it will be created. The default root
                                          directory is `<repo>/.worktrees/<branch>`.

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

## Examples

From `~/src/tools`, attach a root (root dir `~/src/tools`; session name `tools`):

```
treemux attach-root
```

From `~/src/tools`, attach with explicit name (root dir `~/src/tools`; session name `work`):

```
treemux attach-root --name work
```

Attach with explicit directory (root dir `~/src/tools`; session name `tools`):

```
treemux attach-root -d ~/src/tools
```

Attach with explicit name and directory (root dir `~/src/tools`; session name `work`):

```
treemux attach-root --name work -d ~/src/tools
```

Attach for a worktree (root dir `<repo>/.worktrees/feature-x`; session name `feature-x`):

```
treemux attach-root -w feature-x
```

Attach for a worktree with explicit name (root dir `<repo>/.worktrees/feature-x`; session name `work`):

```
treemux attach-root --name work -w feature-x
```

Invalid: worktree and dir together (fails with `--dir and --worktree are mutually exclusive`):

```
treemux attach-root -d ~/src/tools -w feature-x
```

Attach a child inside tmux (root session and cwd derived from current session; child name `nvim`):

```
treemux attach-child --name nvim
```

Attach a child outside tmux (root session `tools`; child name `nvim`):

```
treemux attach-child --root tools --name nvim
```

Show roots with default sort (alphabetic by name):

```
treemux show-roots
```

Example output:

```
  alpha
* beta
  charlie
  delta
```

Show roots ordered by most recently used (`*` marks current root):

```
treemux show-roots --sort-by=most-recently-used
```

Example output:

```
* delta    # most recent
  alpha    # 2nd
  charlie  # 3rd
  beta     # 4th
```

Show children for a root while hiding the current session (`*` marks current session):

```
treemux show-children --root tools --hide-current
```

Example output:

```
  nvim
  logs
```

## Integration tests

Integration tests require a tmux session. Run them from inside tmux:

```
TREEMUX_TMUX=1 go test ./...
```
