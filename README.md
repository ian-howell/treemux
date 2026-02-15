# treemux

`treemux` is a small tmux session helper for root/child session trees. It provides commands to manage
tmux sessions in a hierarchical manner, allowing you to organize your sessions based on their root
directories.

## Commands

### attach

Attach to a tmux session, creating it if it doesn't exist.

ARGS

`--rooted-at <root-name>`:                Attach to a child session rooted at the specified root
                                          session.

`--cmd <command>`:                        Command to run in the session. If not specified, the
                                          session will start with the default shell.

`-d --dir <directory>`:                   Starting directory for the session. If not specified, it
                                          will inherit from the root session (for child sessions) or
                                          use the current directory (for root sessions).

`-w --worktree <worktree/branch>`:        If specified, the session will be rooted at the given git
                                          worktree directory. If the worktree did not exist for the
                                          specified branch, it will be created. If `-d` is not
                                          specified, it will be rooted at .worktrees/<branch>.

### show

Prints a fzf-friendly list of tmux sessions, with root sessions and their child sessions grouped
together.
