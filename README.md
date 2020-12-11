# Envosaurus

It helps you do a bunch of stuff

The aim of envosaurus is to make it easier to work with multiple projects in
different directories (e.g., git repositories).  We want to bring the benefits
of monolithic repositories to many-repository organizations.

This is still very much a work in progress but we'd love to hear your ideas!

## Usage

### Help

```
envosaurus help
```

### Clone

To clone repos from a repo spec file:

```
envosaurus clone --repo-config path/to/projects.json
```

where `projects.json` has the following structure (see `samples/projects.json`):

```
{
    "rootDirectory": "${HOME}/projects",
    "projects": [{
        "name": "project 1",
        "git": {
            "clone": git@github.com:myorg/project1
        }
    }]
}
```

**NOTE** Make sure you have `ssh-agent` running and all relevant identities
loaded.

### Add

To add a repo to a spec file:

```
envosaurus add --repo-config path/to/projects.json
```

This will add the git repository at the current path to the projects repository.

## Contributing

Standard procedure should apply here.  Fork the repo, make a PR, all tests
should pass.