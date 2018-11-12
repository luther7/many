package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/alecthomas/kingpin.v2"
)

// A version of a service.
type Version struct {
	name        string
	description string
	date        time.Time
	author      string
}

// A collection of versions.
type Versions []Version

// A service.
type Service struct {
	name        string
	description string
	git         string
	docker      string
	candidate   Version
	versions    Versions
}

// A table of services. The key is the service's name.
type Services map[string]Service

// The Manyfile is the TOML config containing the versioning information.
type Manyfile struct {
	name       string
	remoteURL  string `toml:"remote_url"`
	remoteName string `toml:"remote_name"`
	versions   Versions
	services   Services
}

// A Many repository.
type Repo struct {
	path     string
	file     string
	manyFile Manyfile
}

// Part of the sort interface.
func (vs Versions) Len() int {
	return len(vs)
}

// Part of the sort interface.
func (vs Versions) Swap(i, j int) {
	vs[i], vs[j] = vs[j], vs[i]
}

// Part of the sort interface.
func (vs Versions) Less(i, j int) bool {
	return vs[i].name < vs[j].name
}

// Add a version to a collection of versions.
func (vs Versions) Add(v Version) {
	// Sort the versions and search for the version to be added.
	sort.Sort(vs)
	i := sort.Search(len(vs), func(i int) bool { return vs[i].name >= v.name })
	// The version already exists in the collection.
	if i < len(vs) && vs[i] == v {
		// Override the version.
		vs[i] = v
		// The version does not exist in the collection.
	} else {
		// Insert the version.
		vs = append(vs, Version{})
		copy(vs[i+1:], vs[i:])
		vs[i] = v
	}
}

// Merge services.
func (s1 *Service) Merge(s2 Service) error {
	if s2.name != "" {
		s1.name = s2.name
	}
	if s2.description != "" {
		s1.description = s2.description
	}
	if s2.git != "" {
		s1.git = s2.git
	}
	if s2.docker != "" {
		s1.docker = s2.docker
	}
	if s2.candidate != (Version{}) {
		s1.candidate = s2.candidate
	}
	if s2.versions != nil {
		for _, v := range s2.versions {
			s1.versions.Add(v)
		}
	}
	return nil
}

// Merge Manyfiles.
func (f1 *Manyfile) Merge(f2 Manyfile) error {
	if f2.name != "" {
		f1.name = f2.name
	}
	if f2.remoteURL != "" {
		f1.remoteURL = f2.remoteURL
	}
	if f2.remoteName != "" {
		f1.remoteName = f2.remoteName
	}
	if f2.versions != nil {
		for _, f := range f2.versions {
			f1.versions.Add(f)
		}
	}
	if f2.services != nil {
		for n, s2 := range f2.services {
			s1, ok := f1.services[n]
			if ok {
				s1.Merge(s2)
			}
		}
	}
	return nil
}

// Save the repo.
func (r *Repo) Save() error {
	// Check if the repo dir exists.
	_, err := os.Stat(r.path)
	if err != nil {
		// The repo dir exists and there was an error.
		if !os.IsNotExist(err) {
			return err
		}
		// Repo dir doesn't exist. Make the repo dir.
		err = os.MkdirAll(r.path, 700)
		if err != nil {
			return err
		}
	}
	// Check if the Manyfile exists.
	_, err = os.Stat(r.file)
	// The Manyfile exists and there was an error.
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	// Create the Manyfile. If it already exists it will be truncated.
	f, err := os.Create(r.file)
	defer f.Close()
	if err != nil {
		return err
	}
	// Write the Manyfile.
	e := toml.NewEncoder(f)
	err = e.Encode(r.manyFile)
	if err != nil {
		return err
	}
	return nil
}

// Load the repo.
func LoadRepo(repo string, file string) (*Repo, error) {
	// Clean the repo path.
	repo = filepath.Clean(repo)
	// Create a path from the repo path and the file path.
	file = filepath.Join(repo, file)
	// Check if the repo dir exists.
	_, err := os.Stat(repo)
	if err != nil {
		return nil, err
	}
	// Check if the repo file exists.
	_, err = os.Stat(file)
	if err != nil {
		return nil, err
	}
	// Decode the repo's Manyfile.
	var m *Manyfile
	_, err = toml.DecodeFile(file, m)
	if err != nil {
		return nil, err
	}
	// Return a new repo struct.
	return &Repo{
		path:     repo,
		file:     file,
		manyFile: *m,
	}, nil
}

// Initialise the repo.
func InitRepo(
	repo string,
	file string,
	name string,
	remoteURL string,
	remoteName string,
	update bool,
	noClone bool,
) error {
	// Attempt to load an existing repo.
	r, err := LoadRepo(repo, file)
	if err != nil {
		// Repo exists and there was an error loading it.
		if !os.IsNotExist(err) {
			return err
		}
		// TODO clone.
		// Repo does not exist. Create it.
		err = r.Save()
		if err != nil {
			return err
		}
		return nil
	}
	// Repo exists, update it if flagged.
	if !update {
		return errors.New("Repository already exists. Use --update to update it.")
	}
	// Update the repo. Merge in the new repo details.
	err = r.manyFile.Merge(
		Manyfile{
			name:       name,
			remoteURL:  remoteURL,
			remoteName: remoteName,
		},
	)
	if err != nil {
		return err
	}
	// Save the updated repo.
	err = r.Save()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var (
		// The application's version.
		version = "0.1.0"
		// Kingpin vars.
		a = kingpin.New(
			"many",
			"Service versioning tool.",
		)
		argRepo = a.Flag(
			"repo",
			"Path to the Many repository.",
		).Short('r').Default(".").String()
		argFile = a.Flag(
			"file",
			"Name of the Many file.",
		).Short('f').Default("Many.toml").String()
		argInit = a.Command(
			"init",
			"Initialize a new Many repository with an empty versioning file. "+
				"If a repository exists at the provided URL then it is cloned.",
		)
		argInitName = argInit.Arg(
			"name",
			"Name of the Many repository.",
		).Required().String()
		argInitRemoteURL = argInit.Arg(
			"git-url",
			"URL of the Git remote.",
		).Required().String()
		argInitRemoteName = argInit.Flag(
			"remote",
			"Name of the Git remote.",
		).Short('m').Default("origin").String()
		argInitUpdate = argInit.Flag(
			"update",
			"Update Many repository details if it is already initialised.",
		).Short('u').Default("false").Bool()
		argInitNoClone = argInit.Flag(
			"no-clone",
			"Do not clone the from an existing repository at the remote URL.",
		).Short('n').Default("false").Bool()
		argPull = a.Command(
			"pull",
			"Pull changes from the remote Many repository.",
		)
		argPush = a.Command(
			"push",
			"Push changes to the remote Many Grepository.",
		)
		argCreate = a.Command(
			"create",
			"Register a new microservice with Many.",
		)
		argCreateUpdate = argCreate.Flag(
			"update",
			"Update microservice details if it already exists.",
		).Short('u').Default("false").Bool()
		argCreateName = argCreate.Arg(
			"service",
			"Name of microservice.",
		).Required().String()
		argCreateDescription = argCreate.Flag(
			"description",
			"Description of microservice.",
		).Short('s').String()
		argCreateGit = argCreate.Flag(
			"git",
			"URL of the Git repository for the microservice.",
		).Short('g').String()
		argCreateDocker = argCreate.Flag(
			"docker",
			"URL of the Docker repository for the microservice.",
		).Short('c').URL()
		argView = a.Command(
			"view",
			"View details for microservices.",
		)
		argViewName = argView.Arg(
			"services",
			"CSV list of microservices.",
		).Required().String()
		argDelete = a.Command(
			"delete",
			"Delete a microservice.",
		)
		argDeleteName = argDelete.Arg(
			"service",
			"Name of microservice.",
		).Required().String()
		argPromote = a.Command(
			"promote",
			"Promote a candidate version of a microservice.",
		)
		argPromoteName = argPromote.Arg(
			"service",
			"Name of microservice.",
		).Required().String()
		argPromoteVersion = argPromote.Arg(
			"version",
			"Candidate version.",
		).Required().String()
		argCurrent = a.Command(
			"current",
			"View the current overall version.",
		)
		argRelease = a.Command(
			"release",
			"Create a new overall version from the candidates.",
		)
		argReleaseCategory = argRelease.Arg(
			"version",
			"Version to increment for this release.",
		).Required().Enum("patch", "minor", "major")
	)
	// Kingpin.
	a.HelpFlag.Short('h')
	a.Version(version)
	a.VersionFlag.Short('v')
	c := kingpin.MustParse(a.Parse(os.Args[1:]))
	// Loggers. No prefix. No timestamps.
	lstd := log.New(os.Stout, "", 0)
	lerr := log.New(os.Sterr, "", 0)
	// Switch on command.
	switch c {
	case "init":
		err := InitRepo(
			argRepo,
			argFile,
			argInitName,
			argInitRemoteURL,
			argInitRemoteName,
			argInitUpdate,
			argInitNoClone,
		)
		if err != nil {
			lstderr.Fatal(err)
		}
		lstdout.Println("Initialised Many repo.")
		// case "pull":
		// 	lstdout.Println("Pulling Many repo.")
		// 	// TODO
		// 	if err != nil {
		// 		lstderr.Fatal(err)
		// 	}
		// 	lstdout.Println("Success.")
		// case "push":
		// 	lstdout.Println("Pushing Many repo.")
		// 	// TODO
		// 	if err != nil {
		// 		lstderr.Fatal(err)
		// 	}
		// 	lstdout.Println("Success.")
		// case "register":
		// 	// TODO
		// 	if err != nil {
		// 		lstderr.Fatal(err)
		// 	}
		// 	lstdout.Println("Registered service.")
		// case "view":
		// case "delete":
		// 	// TODO
		// 	if err != nil {
		// 		lstderr.Fatal(err)
		// 	}
		// 	lstdout.Println("Deleted service.")
		// case "promote":
		// 	// TODO
		// 	if err != nil {
		// 		lstderr.Fatal(err)
		// 	}
		// 	lstdout.Println("Promoted service.")
		// case "increment":
		// 	// TODO
		// 	if err != nil {
		// 		lstderr.Fatal(err)
		// 	}
		// 	lstdout.Println("Incremented service.")
	}
}
