//go:build windows && !ordax_owner_prototype

package main

func ownerReadOnlyBackendExecutable(string) (string, bool) { return "", false }

func ownerPrototypePhysicalBackend() (string, bool) {
	return "", false
}

func ownerPrototypeBuildInfo() (string, string, bool) {
	return "", "", false
}
