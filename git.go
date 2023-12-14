package git_ver

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func GetLatestVersion(repo *git.Repository) (string, error) {
	refs, tagErr := repo.Tags()
	if tagErr != nil {
		return "", tagErr
	}
	tagMap := make(map[plumbing.Hash]string)
	if err := refs.ForEach(func(t *plumbing.Reference) error {
		tagMap[t.Hash()] = t.Name().Short()
		return nil
	}); err != nil {
		return "", err
	}
	log, logErr := repo.Log(&git.LogOptions{
		Order: git.LogOrderCommitterTime,
	})
	if logErr != nil {
		return "", logErr
	}
	for obj, err := log.Next(); err == nil; {
		if tag, exists := tagMap[obj.Hash]; exists {
			return tag, nil
		}
	}
	return "", nil
}
