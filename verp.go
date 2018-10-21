package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	version = "0.1.0"
	app     = kingpin.New(
		"verp",
		"Microservice versioning tool.",
	)
	versionFile = app.Flag(
		"file",
		"File containing version info.",
	).Short('f').Default("./Versionfile").String()

	register = app.Command(
		"register",
		"Register a microservice.",
	)
	registerName = register.Arg(
		"name",
		"Name of microservice.",
	).Required().String()
	registerDescription = register.Flag(
		"description",
		"Description of microservice.",
	).Short('d').String()
	registerGit = register.Flag(
		"git",
		"URL of the Git repository for the microservice.",
	).Short('g').URL()
	registerDocker = register.Flag(
		"docker",
		"URL of the Docker repository for the microservice.",
	).Short('r').URL()
	registerUpdate = register.Flag(
		"update",
		"Update microservice details if it is already registered.",
	).Short('u').Default("false").Bool()

	delete = app.Command(
		"delete",
		"Delete a microservice.",
	)
	deleteName = delete.Arg(
		"name",
		"Name of microservice.",
	).Required().String()

	promote = app.Command(
		"promote",
		"Promote a candidate version.",
	)
	promoteName = promote.Arg(
		"name",
		"Name of microservice.",
	).Required().String()
	promoteVersion = promote.Arg(
		"version",
		"Candidate version.",
	).Required().String()

	increment = app.Command(
		"increment",
		"Increment the overall version.",
	)
	incrementCategory = increment.Arg(
		"category",
		"Category to increment.",
	).Required().Enum("patch", "minor", "major")
)

func main() {
	app.HelpFlag.Short('h')
	app.Version(version)
	app.VersionFlag.Short('v')
	command := kingpin.MustParse(app.Parse(os.Args[1:]))
	println(*versionFile)

	switch command {
	case "register":
		println(*registerName)
		if *registerDescription != "" {
			println(*registerDescription)
		}
		if *registerGit != nil {
			println((*registerGit).String())
		}
		if *registerDocker != nil {
			println((*registerDocker).String())
		}
		if *registerUpdate {
			println("update")
		}
		// TODO

	case "delete":
		println(*versionFile)
		println(*deleteName)
		// TODO

	case "promote":
		println(*versionFile)
		println(*promoteName)
		println(*promoteVersion)
		// TODO

	case "increment":
		println(*versionFile)
		println(*incrementCategory)
		// TODO
	}
}
