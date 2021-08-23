module github.com/Hellcatlk/network-operator

go 1.16

require (
	github.com/go-logr/logr v0.4.0
	github.com/securego/gosec v0.0.0-20200401082031-e946c8c39989
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v0.21.2
	sigs.k8s.io/controller-runtime v0.9.2
	sigs.k8s.io/controller-tools v0.6.1
	sigs.k8s.io/kustomize/kustomize/v3 v3.10.0
)
