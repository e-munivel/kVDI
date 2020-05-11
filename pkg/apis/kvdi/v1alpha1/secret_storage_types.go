package v1alpha1

import "sigs.k8s.io/controller-runtime/pkg/client"

// SecretsProvider provides an interface for an app instance to get and store
// any secrets it needs. Currenetly there is only a k8s secret provider, but
// this intreface could be implemented for things like vault.
type SecretsProvider interface {
	// Setup is called before the interface is used for any operations
	Setup(client.Client, *VDICluster) error
	// ReadSecret should return the contents of a secret by name.
	ReadSecret(name string) (contents []byte, err error)
	// WriteSecret should store a secret, replacing any existing one with the
	// same name.
	WriteSecret(name string, contents []byte) error
	// Close should handle any cleanup logic for the backend. This method is invoked
	// after temporary usages of the secret engine. This shouldn't be destructive,
	// but it should ensure any opened sockets are closed cleanly, spawned
	// goroutines are finished, and no other dangling references left behind.
	Close() error
}
