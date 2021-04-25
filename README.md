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
rrr run --commit # commits changes from script
rrr run --push # same as above, but also pushes branches after commit (implies --commit)
rrr run --pr # same as above, but also opens PRs after commit (implies --push)
# use -i flag for interactive (as opposed the default -p for patch) x
# use -A flag for adding all changes without prompts
```

Other commands:

```sh
rrr push # Pushes local changes to remote x
rrr push --delete # (or -D) Deletes remote branch x
```

```sh
rrr clone --create-branch # (-b) Creates and checks out
```

```sh
rrr pr # Opens PRs (using config.yaml settings) for any repos that have outstanding changes in their directory
```

**To Do:**

- [ ] Parallelize cloning of repos (go routine?)
- [ ] Skip a repo during `run` if there are no changes to be staged

- Handle case where: changes are committed/pushed & PRs are already opened & user needs to push updates without deleting existing branch/PR

- Clone
- Run script
- Stage
- Commit
- Push
- PR
