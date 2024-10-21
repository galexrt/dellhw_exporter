## Creating a new Release

1. Update the version.
    1. [`VERSION` file](VERSION)
    2. Helm chart: `charts/dellhw_exporter/Chart.yaml` `appVersion:` line and bump the Helm chart `version:` by a patch release version.
       1. Make sure to run `make helm-docs` in the root of the repo to update the helm chart docs.
2. Create an entry in the [`CHANGELOG.md` file](CHANGELOG.md).
    Example of a changelog entry:

    ```
    ## 1.12.0 / 2022-02-02

    * [ENHANCEMENT] Added Pdisk Remaining Rated Write Endurance Metric by @adityaborgaonkar
    * [BUGFIX] ci: fix build routine issues
    ```

    The following "kinds" of entries can be used:

    * `CHANGE`
    * `FEATURE`
    * `ENHANCEMENT`
    * `BUGFIX`
3. Commit the version increase with a commit messages in format: `version: bump to v1.12.0`
4. Create the `git` tag using `git tag v1.12.0`
5. Now push the changes and commit using `git push && git push --tags`
6. In a few minutes the new release should be available for "approval" under the [releases section](https://github.com/galexrt/dellhw_exporter/releases). Edit and save the release on GitHub and the release is complete.
