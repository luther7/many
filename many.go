package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	version = "0.1.0"
	app     = kingpin.New(
		"many",
		"Microservice versioning tool.",
	)
	argRepo = app.Flag(
		"repo",
		"Path to the Many repository.",
	).Short('r').Default(".").String()
	argFile = app.Flag(
		"file",
		"Name of the Many file.",
	).Short('f').Default("Many.toml").String()
	initial = app.Command(
		"init",
		"Initialize a new Many repository with an empty versioning file. "+
			"If a repository exists at the provided URL then it is cloned.",
	)
	initialUpdate = initial.Flag(
		"update",
		"Update Many repository details if it is already initialised.",
	).Short('u').Default("false").Bool()
	initialName = initial.Arg(
		"name",
		"Name of the Many repository.",
	).Required().String()
	initialRemoteURL = initial.Arg(
		"git-url",
		"URL of the Git remote.",
	).Required().String()
	initialRemoteName = initial.Flag(
		"remote",
		"Name of the Git remote.",
	).Short('m').Default("origin").String()
	initialVersion = initial.Flag(
		"inital-version",
		"The starting version.",
	).Short('i').Default("0.0.0").String()
	initialNoClone = initial.Flag(
		"no-clone",
		"Do not clone the from an existing repository at the remote URL.",
	).Short('n').Default("false").Bool()
	pull = app.Command(
		"pull",
		"Pull changes from the remote Many repository.",
	)
	push = app.Command(
		"push",
		"Push changes to the remote Many Grepository.",
	)
	create = app.Command(
		"create",
		"Register a new microservice with Many.",
	)
	createUpdate = create.Flag(
		"update",
		"Update microservice details if it already exists.",
	).Short('u').Default("false").Bool()
	createName = create.Arg(
		"service",
		"Name of microservice.",
	).Required().String()
	createDescription = create.Flag(
		"description",
		"Description of microservice.",
	).Short('s').String()
	createGit = create.Flag(
		"git",
		"URL of the Git repository for the microservice.",
	).Short('g').String()
	createDocker = create.Flag(
		"docker",
		"URL of the Docker repository for the microservice.",
	).Short('c').URL()
	view = app.Command(
		"view",
		"View details for microservices.",
	)
	viewName = view.Arg(
		"services",
		"CSV list of microservices.",
	).Required().String()
	delete = app.Command(
		"delete",
		"Delete a microservice.",
	)
	deleteName = delete.Arg(
		"service",
		"Name of microservice.",
	).Required().String()
	promote = app.Command(
		"promote",
		"Promote a candidate version of a microservice.",
	)
	promoteName = promote.Arg(
		"service",
		"Name of microservice.",
	).Required().String()
	promoteVersion = promote.Arg(
		"version",
		"Candidate version.",
	).Required().String()
	current = app.Command(
		"current",
		"View the current overall version.",
	)
	release = app.Command(
		"release",
		"Create a new overall version from the candidates.",
	)
	releaseCategory = release.Arg(
		"version",
		"Version to increment for this release.",
	).Required().Enum("patch", "minor", "major")
)

type ManyVersion struct {
	version     string
	date        time.Time
	description string
}

type ManyService struct {
	name        string
	description string
	git         string
	docker      string
	candidate   ManyVersion
	versions    []ManyVersion
}

type ManyFile struct {
	name       string
	remoteURL  string `toml:"remote_url"`
	remoteName string `toml:"remote_name"`
	versions   []ManyVersion
	services   map[string]ManyService
}

type ManyRepo struct {
	path     string
	file     string
	manyFile ManyFile
}

func (orig *ManyService) update(changes ManyService) error {
	orig.name = changes.name
	orig.description = changes.description
	orig.git = changes.git
	orig.docker = changes.docker
	orig.candidate = changes.candidate
	// TODO check merged versions
	// newest first
	orig.versions = append(changes.versions, orig.versions...)
	return nil
}

func (orig *ManyFile) update(changes ManyFile) error {
	orig.name = changes.name
	orig.remoteURL = changes.remoteURL
	orig.remoteName = changes.remoteName
	// TODO check merged versions
	// newest first
	orig.versions = append(changes.versions, orig.versions...)
	for name, service := range changes.services {
		origService, ok := orig.services[name]
		if ok {
			origService.update(service)
		}
	}
	return nil
}

func (orig *ManyRepo) update(changes ManyRepo) error {
	orig.path = changes.path
	orig.file = changes.file
	err := orig.manyFile.update(changes.manyFile)
	if err != nil {
		return err
	}
	return nil
}

func LoadManyRepo(repoPath string, filePath string) (*ManyRepo, error) {
	repoPath = filepath.Clean(repoPath)
	filePath = filepath.Join(repoPath, filePath)
	_, err := os.Stat(repoPath)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	var manyFile *ManyFile
	_, err = toml.DecodeFile(filePath, manyFile)
	if err != nil {
		return nil, err
	}
	return &ManyRepo{
		path:     repoPath,
		file:     filePath,
		manyFile: *manyFile,
	}, nil
}

func CreateManyRepo(
	repoPath string,
	filePath string,
	name string,
	remoteURL string,
	remoteName string,
	update bool,
	noClone bool,
) (*ManyRepo, error) {
	repo := ManyRepo{
		path: repoPath,
		file: filePath,
		manyFile: ManyFile{
		  name:       name,
		  remoteURL:  remoteURL,
		  remoteName: remoteName,
		  versions:   []ManyVersion{},
		  services:   map[string]ManyService{},
	  },
	}
	origRepo, err := LoadManyRepo(repoPath, filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if !update {
			return nil, errors.New("Repository already exists. Use --update to update it.")
		}
		err = origRepo.update(repo)
		if err != nil {
			return nil, err
		}
		return origRepo, nil
	}
	// TODO save
	return &repo, nil
}

func main() {
	app.HelpFlag.Short('h')
	app.Version(version)
	app.VersionFlag.Short('v')
	command := kingpin.MustParse(app.Parse(os.Args[1:]))
	fmt.Println("Started.")

	switch command {
	case "init":
		fmt.Println("Initialising Many repo.")

	case "pull":
		// TODO

	case "push":
		// TODO

	case "register":
		// TODO

	case "view":
		// TODO

	case "delete":
		// TODO

	case "promote":
		// TODO

	case "increment":
		// TODO
	}

	fmt.Println("Success.")
}
