package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	version = "0.1.0"
	app     = kingpin.New(
		"many",
		"Microservice versioning tool.",
	)
	repo = app.Flag(
		"repo",
		"Path to the Git repository containing the version file.",
	).Short('r').Default(".").File()
	manyFile = app.Flag(
		"file",
		"Path to the file containing the version information.",
	).Short('f').Default("./many.toml").String()

	initialise = app.Command(
		"init",
		"Initialize a new Many Git repository with an empty versioning file. "+
			"If a repository exists at the provided URL then it is cloned.",
	)
	initialiseUpdate = initialise.Flag(
		"update",
		"Update Many Git repository details if it is already initialised.",
	).Short('u').Default("false").Bool()
	initialiseName = initialise.Arg(
		"name",
		"Name of the Many Git repository.",
	).Required().String()
	initialiseRemoteUrl = initialise.Arg(
		"git-url",
		"URL of the Git remote.",
	).URL()
	initialiseRemoteName = initialise.Flag(
		"remote",
		"Name of the Git remote.",
	).Short('m').Default("origin").String()
	initialiseVersion = initialise.Flag(
		"inital-version",
		"The starting version.",
	).Short('i').Default("0.0.0").String()
	initialiseNoClone = initialise.Flag(
		"no-clone",
		"Do not clone the from an existing repository at the remote URL.",
	).Short('n').Default("false").Bool()

	pull = app.Command(
		"pull",
		"Pull changes from the remote Many Git repository.",
	)

	push = app.Command(
		"push",
		"Push changes to the remote Many Git repository.",
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
	).Short('g').URL()
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

func main() {
	app.HelpFlag.Short('h')
	app.Version(version)
	app.VersionFlag.Short('v')
	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	switch command {
	case "init":
		// TODO

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
}
