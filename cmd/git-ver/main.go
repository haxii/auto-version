package main

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/haxii/git-ver"
	"os"
	"strings"
)

func CheckIfError(err error) {
	if err == nil {
		return
	}
	fmt.Printf("error: %s\n", err)
	os.Exit(1)
}

func helpMessage() {
	fmt.Printf("Usage: %s [patch|minor|major] to bump a version\n", os.Args[0])
	os.Exit(1)
}

func isIncValid(inc string) bool {
	return inc == "patch" || inc == "minor" || inc == "major"
}

func main() {
	path, pathErr := os.Getwd()
	CheckIfError(pathErr)
	inc := ""
	if len(os.Args) > 2 {
		helpMessage()
	} else if len(os.Args) == 2 {
		if inc = os.Args[1]; !isIncValid(inc) {
			helpMessage()
		}
	}
	//path = "/Users/zichao/workspace/sif"
	repo, openErr := git.PlainOpen(path)
	CheckIfError(openErr)
	head, headErr := repo.Head()
	CheckIfError(headErr)
	latestVer, verErr := git_ver.GetLatestVersion(repo)
	CheckIfError(verErr)
	if latestVer == nil {
		fmt.Println("no latest version found")
		return
	}
	ver, err := semver.NewVersion(latestVer.Name)
	if err != nil {
		fmt.Printf("latest version %s is not a valid semver: %s\n", latestVer.Name, err)
		os.Exit(1)
	}
	if head.Hash() == latestVer.Hash {
		fmt.Printf("head already tagged with version %s\n", latestVer.Name)
		os.Exit(1)
	}
	if len(inc) == 0 {
		fmt.Printf("latest version %s, type patch, minor or major to bump, default is patch: ", ver)
		_, err = fmt.Scanln(&inc)
		if err != nil && strings.Contains(err.Error(), "unexpected newline") {
			err = nil
		}
		CheckIfError(err)
		if len(inc) > 0 && !isIncValid(inc) {
			CheckIfError(fmt.Errorf("invalid command %s, only patch, minor or major acceptted", inc))
		}
	}
	newVer := ""
	switch inc {
	case "minor":
		newVer = ver.IncMinor().String()
	case "major":
		newVer = ver.IncMajor().String()
	default:
		newVer = ver.IncPatch().String()
	}
	if _, err = repo.CreateTag("v"+newVer, head.Hash(), nil); err != nil {
		fmt.Printf("new version %s tag failed with error %s\n", newVer, err)
	} else {
		fmt.Printf("new version %s tagged\n", newVer)
	}
}
