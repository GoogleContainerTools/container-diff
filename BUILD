load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_prefix")
load("@bazel_gazelle//:def.bzl", "gazelle")

gazelle(
    name = "gazelle",
    external = "vendored",
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
    embed = [":go_default_library"],
    importpath = "github.com/GoogleCloudPlatform/container-diff",
    visibility = ["//visibility:public"],
)
