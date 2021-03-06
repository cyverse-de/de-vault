package vaulter

import (
	"strings"

	vault "github.com/hashicorp/vault/api"
)

// Mounter is an interface for objects that can mount Vault backends.
type Mounter interface {
	Mount(path string, m *vault.MountInput) error
}

// MountLister is an interface for objects that can list mounted Vault backends.
type MountLister interface {
	ListMounts() (map[string]*vault.MountOutput, error)
}

// MountConfigGetter is an interface for objects that can get the configuration
// for a mount in Vault.
type MountConfigGetter interface {
	MountConfig(path string) (*vault.MountConfigOutput, error)
}

// MountTuner is an interface for objects that need to configure a mount in
// Vault.
type MountTuner interface {
	TuneMount(path string, input vault.MountConfigInput) error
}

// MountWriter is an interface for objects that can write to a path in a Vault
// backend.
type MountWriter interface {
	Write(c *vault.Client, path string, data map[string]interface{}) (*vault.Secret, error)
}

// MountReader is an interface for objects that can read data from a path in a
// Vault backend.
type MountReader interface {
	Read(c *vault.Client, path string) (*vault.Secret, error)
}

// PathDeleter is an interface for deleting information from a mount, not for
// deleting the mount itself.
type PathDeleter interface {
	Delete(c *vault.Client, path string) (*vault.Secret, error)
}

// Unmounter is an interface for objects that can unmount a Vault
// backend.
type Unmounter interface {
	Unmount(path string) error
}

// MountReaderWriter defines an interface for doing role related operations.
type MountReaderWriter interface {
	ClientGetter
	MountWriter
	MountReader
}

// MountDeleter defines and interface for deleting content from a path in a
// mounted backend.
type MountDeleter interface {
	ClientGetter
	PathDeleter
}

// MountConfiguration is a flattened representation of the configs that the Vault API
// supports for the backend mounts.
type MountConfiguration struct {
	Type            string
	Description     string
	DefaultLeaseTTL string
	MaxLeaseTTL     string
}

// Mount mounts a vault backend with the provided
func Mount(m Mounter, path string, c *MountConfiguration) error {
	return m.Mount(path, &vault.MountInput{
		Type:        c.Type,
		Description: c.Description,
		Config: vault.MountConfigInput{
			DefaultLeaseTTL: c.DefaultLeaseTTL,
			MaxLeaseTTL:     c.MaxLeaseTTL,
		},
	})
}

// Unmount unmounts a vault backend with the provided path.
func Unmount(u Unmounter, path string) error {
	return u.Unmount(path)
}

// MountConfig returns the config for the passed in mount rooted at the given
// path.
func MountConfig(m MountConfigGetter, path string) (*vault.MountConfigOutput, error) {
	return m.MountConfig(path)
}

// IsMounted returns true if the given path is mounted as a backend in Vault.
func IsMounted(l MountLister, path string) (bool, error) {
	var (
		hasPath bool
		err     error
	)
	mounts, err := l.ListMounts()
	if err != nil {
		return false, err
	}
	for m := range mounts {
		if strings.TrimSuffix(m, "/") == path {
			hasPath = true
		}
	}
	return hasPath, nil
}
