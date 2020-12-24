package repository

import "fmt"

type RepoType = string

type RepoFactory interface {
	Build(cfg interface{}) (*Repository, error)
	RepoType() RepoType
}

func GetFactory(repoType RepoType, factories ...RepoFactory) (RepoFactory, error) {
	for _, factory := range factories {
		if factory.RepoType() == repoType {
			return factory, nil
		}
	}

	return nil, fmt.Errorf("unknown repo type: " + repoType)
}
