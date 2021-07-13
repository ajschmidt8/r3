# r3

`r3` (_RAPIDS repository reviser_) is a CLI tool built with Golang that automates the process of making changes to all of the RAPIDS repository.

## Installation

Download the latest binary from GitHub and add it to one of your `PATH` directories:

```sh
# wget https://github.com/ajschmidt8/r3/releases/latest/download/r3_macos -O r3
wget https://github.com/ajschmidt8/r3/releases/latest/download/r3_linux -O r3
chmod +x ./r3
sudo mv ./r3 /usr/local/bin
```

## Usage

First, create and enter an empty directory and run:

```sh
r3 init
```

The first time any command is run, it will prompt you to authenticate with GitHub:

![](/doc/r3_init.png)

The `init` command will then generate the following files:

- `scr.sh` - The shell script to be run in each repository
- `config.yaml` - Some configuration settings (repository list, PR title, body, labels, etc.) to be used when committing the changes

Now run `r3 run` in order to execute `scr.sh` in each repository that's listed in `config.yaml`. By default, `r3 run` will allow you to review all of your changes interactively using `git add -p`.

![](/doc/r3_run.png)

If the changes look correct, you can re-run the command again with the `--pr` flag (i.e. `r3 run --pr`) in order to open pull-requests in each repository. You can also use `r3 run -A --pr` to avoid being prompted to review the changes again.

![](/doc/r3_run_pr.png)

Use `r3 -h` or `r3 <command> -h` for more info and available flags
