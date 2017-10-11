load("@io_bazel_rules_go//go:def.bzl", "gazelle", "go_binary", "go_library", "go_prefix")

gazelle(
    name = "gazelle",
    build_tags = [
        "container_image_ostree_stub",
        "containers_image_openpgp",
    ],
    external = "vendored",
    prefix = "github.com/GoogleCloudPlatform/container-diff",
)

go_prefix("github.com/GoogleCloudPlatform/container-diff")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    visibility = ["//visibility:private"],
    deps = ["//cmd:go_default_library"],
)

go_binary(
    name = "container-diff",
    library = ":go_default_library",
    visibility = ["//visibility:public"],
)
