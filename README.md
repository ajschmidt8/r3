# rapids-repo-reviser

`rapids-repo-reviser` is a CLI tool built with Golang that automates the process of making changes to all of the RAPIDS repos.

## Usage

First, create and enter an empty directory and run:

```sh
rrr init
```

This command will generate the following files:

- `scr.sh` - The shell script to be run in each repo
- `config.yaml` - Some configuration settings (repo list, PR title, body, labels, etc.) to be used when committing the changes
- `README.md` - A short README file describing common usage

Then, run the following command to execute your script in all of the repos listed in `config.yaml`:

```sh
rrr run # runs scr.sh script on all repos in repos subdir
rrr run --push # same as above, but also pushes branches after each commit
rrr run --prs # same as above, but also opens PRs after each commit (implies --push)
# use -i flag for interactive (as opposed the default -p for patch) x
# use -a flag for adding all changes without prompts
```

Other commands:

```sh
rrr push # Pushes local changes to remote
rrr push --delete # (or -D) Deletes remote branch
```

```sh
rrr prs # Opens PRs (using config.yaml settings) for any repos that have outstanding changes in their directory
```

```sh
rrr auth # Sets GitHub API key used for opening PRs.
```

**To Do:**

- [x] Add setting for `git add` `-i` vs `-p` flags (`-i` is necessary for interactively staging files that aren't yet tracked)
- [ ] Make _base_branch_ dynamic - get default branch from cuDF repo
- [ ] Handle forks (create forks for user if they don't already exist)
- [ ] Setup global config file (`~/rrr.yaml`) on `rrr init`?
  - Contains GitHub username and GitHub API key
- [ ] Parallelize cloning of repos (go routine?)
- [ ] Update all help messages
- [ ] Add a global flag or config.yaml entry for repo owner (i.e. for testing in forks instead of `rapidsai` org)

  - Will be necessary to use for clones & PRs

- Handle case where: changes are committed/pushed & PRs are already opened & user needs to push updates without deleting existing branch/PR

- Clone
- Run script
- Stage
- Commit
- Push
- PR
