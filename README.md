# r3

`r3` (_RAPIDS repo reviser_) is a CLI tool built with Golang that automates the process of making changes to all of the RAPIDS repos.

## Usage

First, create and enter an empty directory and run:

```sh
r3 init
```

This command will generate the following files:

- `scr.sh` - The shell script to be run in each repo
- `config.yaml` - Some configuration settings (repo list, PR title, body, labels, etc.) to be used when committing the changes

Then, run the following command to execute your script in all of the repos listed in `config.yaml`:

```sh
r3 run # runs scr.sh script on all repos in repos subdir
r3 run --commit # commits changes from script
r3 run --push # same as above, but also pushes branches after commit (implies --commit)
r3 run --pr # same as above, but also opens PRs after commit (implies --push)

# use -i flag for interactive (as opposed the default -p for patch)
# use -A flag for adding all changes without prompts
```

Other available commands include:

```sh
r3 clone
r3 commit
r3 pr
r3 push
```

Use `r3 -h` or `r3 <command> -h` for more info and available flags

## To Do:

- [ ] Improve logging during clones
- [ ] Unit tests
