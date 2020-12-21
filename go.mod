module github.com/GoogleContainerTools/container-diff

go 1.15

replace golang.org/x/sys => golang.org/x/sys v0.0.0-20190830141801-acfa387b8d69

require (
	code.cloudfoundry.org/bytefmt v0.0.0-20180906201452-2aa6f33b730c
	github.com/docker/distribution v0.0.0-20200319173657-742aab907b54 // indirect
	github.com/docker/docker v1.4.2-0.20190219180918-740349757396
	github.com/fsouza/go-dockerclient v1.3.6
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/google/go-containerregistry v0.0.0-20190214194807-bada66e31e55
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/nightlyone/lockfile v0.0.0-20180618180623-0ad87eef1443
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/pkg/profile v1.2.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	google.golang.org/grpc v1.28.1 // indirect
)
