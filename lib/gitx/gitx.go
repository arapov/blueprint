package gitx

import (
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type Info struct {
	Hostname string
	Username string
	SshKey   string
	Repo     string
}

type Tree struct {
	*git.Worktree
}

func (c *Tree) Root() string {
	return c.Filesystem.Root()
}

func (c Info) sshAuth() (*ssh.PublicKeys, error) {
	sshAuth, err := ssh.NewPublicKeysFromFile(c.Username, c.SshKey, "")
	if err != nil {
		return nil, err
	}

	return sshAuth, err
}

func (c Info) Connect() (*Tree, error) {
	sshAuth, err := c.sshAuth()
	if err != nil {
		return nil, err
	}

	repo, err := git.PlainOpen(c.Repo)
	if err != nil {
		repo, err = git.PlainClone(c.Repo, false, &git.CloneOptions{
			URL:               c.Hostname,
			ReferenceName:     "refs/heads/master",
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			Auth:              sshAuth,
		})
		if err != nil {
			return nil, err
		}
	}

	wt, _ := repo.Worktree()
	wt.Pull(&git.PullOptions{RemoteName: "origin", Auth: sshAuth})

	return &Tree{wt}, err
}
