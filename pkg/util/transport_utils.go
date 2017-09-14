/*
Copyright 2017 Google, Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"github.com/containers/image/docker"
	"github.com/containers/image/docker/archive"
	"github.com/containers/image/docker/daemon"
	"github.com/containers/image/types"
)

// assume that this artifact is one of three things:
//   1) a path to an image in the local docker daemon,
//   2) a path to an image in a remote registry,
//   3) a path to a tarball on the local fs
func GetImageReference(artifact string) (types.ImageReference, error) {
	if IsTar(artifact) {
		transport := archive.Transport
		return transport.ParseReference(artifact)
	}

	localTransport := daemon.Transport
	ref, err := localTransport.ParseReference(artifact)
	if err == nil {
		// this worked, so image was in local daemon. return the ref
		return ref, err
	}
	// otherwise, assume image is in remote registry
	ref, err = docker.ParseReference(artifact)
	if err != nil {
		// nothing worked: return an error
		return nil, err
	}

	return ref, nil
}
