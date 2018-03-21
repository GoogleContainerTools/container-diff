# Copyright 2018 Google, Inc. All rights reserved.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

def _impl(ctx):
  """Implementation of docker_diff"""

  container_diff_loction = ctx.executable._container_diff.short_path
  image_location = ctx.file.image.short_path
  
  # Shell script to execute container-diff with appropriate flags
  content = """\
#!/bin/bash
set -e
./%s diff %s %s""" % (container_diff_loction, image_location, ctx.attr.diff_base)
  content += " ".join(["--type=%s" % t for t in ctx.attr.diff_types])

  ctx.file_action(
      output = ctx.outputs.executable,
      content = content,
  )

  return struct(runfiles=ctx.runfiles(
    files = [
      ctx.executable._container_diff,
      ctx.file.image
    ]),
  )

#   Diff a bazel image against an image in production with bazel run
#   Runs container-diff on the two images and prints the output.
#   Args:
#     name: name of the rule
#     image: bazel target to an image you have bazel built (must be a tar)
#        ex: image = "@//target/to:image.tar"
#     diff_base: Tag or digest in a remote registry you want to diff against
#        ex: diff_base = "gcr.io/google-appengine/debian8:latest"
#     diff_types: Types flag to pass to container diff
#        ex: ["pip", "file"]

docker_diff = rule(
    attrs = {
        "image": attr.label(
            allow_files = [".tar"],
            single_file = True,
            mandatory = True,
        ),
        "diff_base": attr.string(mandatory = True),
        "diff_types": attr.string_list(
            allow_empty = True,
        ),
        "_container_diff": attr.label(
            default = Label("//:container-diff"),
            executable = True,
            cfg = "host",
        ),
    },
    executable = True,
    implementation = _impl,
)
