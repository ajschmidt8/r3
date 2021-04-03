package shared

const Script = `#!/bin/bash

# This script will be executed in each repo

# Update the changelog of each repo
echo "rrr is great!" > CHANGELOG.md
`

type ConfigInterface struct {
	PR struct {
		BaseBranch string   `yaml:"base_branch"`
		Draft      bool     `yaml:"draft"`
		Title      string   `yaml:"title"`
		Body       string   `yaml:"body"`
		Labels     []string `yaml:"labels"`
	} `yaml:"pr"`
	Repos      []string `yaml:"repos"`
	BranchName string   `yaml:"branch_name"`
	CommitMsg  string   `yaml:"commit_msg"`
}

const Config = `---
pr:
  base_branch: branch-0.19
  draft: false
  title: Update CHANGELOG.md
  body: |
    This PR updates the CHANGELOG.md file using the really great rrr tool.
  labels:
    - non-breaking
    - improvement
#    - breaking
#    - bug
#    - doc
#    - feature request

branch_name: my_new_branch
commit_msg: |
  updating CHANGELOG.md via some new automation

repos:
  - cusignal
#  - clx
#  - cudf
#  - cugraph
#  - cuml
#  - cusignal
#  - cuspatial
#  - cuxfilter
#  - dask-cuda
#  - raft
#  - rmm
#  - cumlprims_mg
`
