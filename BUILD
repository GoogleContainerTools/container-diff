load("@io_bazel_rules_go//go:def.bzl", "gazelle", "go_binary", "go_library", "go_prefix")

gazelle(
    name = "gazelle",
    build_tags = [
        "container_image_ostree_stub",
        "containers_image_openpgp",
    ],
    external = "vendored",
    mode = "fix",
    prefix = "github.com/GoogleCloudPlatform/container-diff",
)

go_prefix("github.com/GoogleCloudPlatform/container-diff")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/GoogleCloudPlatform/container-diff",
    visibility = ["//visibility:private"],
    deps = [
        "//cmd:go_default_library",
        "//vendor/github.com/pkg/profile:go_default_library",
    ],
)

go_binary(
    name = "container-diff",
    importpath = "github.com/GoogleCloudPlatform/container-diff",
    library = ":go_default_library",
    visibility = ["//visibility:public"],
)
