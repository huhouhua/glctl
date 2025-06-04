### Setup your glctl Github Repository
Fork [glctl upstream](https://github.com/huhouhua/glctl/fork) source repository to your own personal repository.
```
$ mkdir -p $GOPATH/src/github.com/huhouhua
$ cd $GOPATH/src/github.com/huhouhua
$ git clone https://github.com/$USER_ID/glctl
$ cd glctl
$ make
$ ./_output/glctl-linux-amd64 --help
```

###  Developer Guidelines

#### Requirements
* [install](https://docs.docker.com/engine/install/) install docker cli
* Go 1.23.0 or later
* Make
* GitLab instance for testing
    - run `make run-gitlab`
* How to add test data
    - run `make testdata`
  
``glctl`` welcomes your contribution. To make the process as seamless as possible, we ask for the following:

* Go ahead and fork the project and make your changes. We encourage pull requests to discuss code changes.
    - Fork it
    - Create your feature branch (git checkout -b my-new-feature)
    - Commit your changes (git commit -am 'Add some feature')
    - Push to the branch (git push origin my-new-feature)
    - Create new Pull Request

* If you have additional dependencies for ``glctl``, ``glctl`` manages its dependencies using `go mod`
    - Run `go get foo/bar`
    - Edit your code to import foo/bar
    - Run `GO111MODULE=on go mod tidy` from top-level folder

* When you're ready to create a pull request, be sure to:
    - Have test cases for the new code. If you have questions about how to do it, please ask in your pull request.
    - Run `make format`
    - Squash your commits into a single commit. `git rebase -i`. It's okay to force update your pull request.
    - Make sure `make tools` completes.

* Read [Effective Go](https://github.com/golang/go/wiki/CodeReviewComments) article from Golang project
    - `glctl` project is conformant with Golang style
    - if you happen to observe offending code, please feel free to send a pull request
