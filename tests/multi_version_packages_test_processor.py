import json
import sys


def _process_test_diff(file_path):
    with open(file_path) as f:
        diffs = json.load(f)

    for diff in diffs:
        if diff["DiffType"] == "NodeDiffer":
            diff_result = diff["Diff"]
            package1_dict = diff_result["Packages1"]
            package2_dict = diff_result["Packages2"]
            diff_result["Packages1"] = _trim_layer_paths(package1_dict)
            diff_result["Packages2"] = _trim_layer_paths(package2_dict)
            diff["Diff"] = diff_result

    with open(file_path, 'w') as f:
        json.dump(diffs, f, indent=4)


def _trim_layer_paths(packages):
    new_packages = {}
    for package, versions in packages.items():
        versions_to_size = {}
        for path in versions.keys():
            version = versions[path]["Version"]
            size = versions[path]["Size"]
            versions_to_size[version] = size
        new_packages[package] = versions_to_size
    return new_packages


if __name__ == '__main__':
    sys.exit(_process_test_diff(sys.argv[1]))
