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

Then, run the following command to execute your script in all of the repos listed in `config.yaml`:

```sh
rrr run # runs scr.sh script on all repos in repos subdir
rrr run --commit # commits changes from script
rrr run --push # same as above, but also pushes branches after commit (implies --commit)
rrr run --pr # same as above, but also opens PRs after commit (implies --push)

# use -i flag for interactive (as opposed the default -p for patch)
# use -A flag for adding all changes without prompts
```

Other available commands include:

```sh
rrr clone
rrr commit
rrr pr
rrr push
```

Use `rrr -h` or `rrr <command> -h` for more info and available flags

## To Do:

- [ ] Parallelize cloning of repos (go routine?)
- [ ] Improve logging during clones
- [ ] Unit tests
