package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"

	"github.com/rotisserie/eris"

	"github.com/solo-io/go-utils/changelogutils"
)

func main() {
	ctx := context.Background()
	app := rootApp(ctx)
	if err := app.Execute(); err != nil {
		fmt.Printf("unable to run: %v\n", err)
		os.Exit(1)
	}
}

type options struct {
	ctx              context.Context
	HugoDataSoloOpts HugoDataSoloOpts
}
type HugoDataSoloOpts struct {
	product string
	version string
	// if set, will override the version when rendering the
	callLatest bool
	noScope    bool
}

func rootApp(ctx context.Context) *cobra.Command {
	opts := &options{
		ctx: ctx,
	}
	app := &cobra.Command{
		Use: "docs-util",
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}
	app.AddCommand(writeVersionScopeDataForHugo(opts))
	app.AddCommand(changelogMdCmd(opts))

	app.PersistentFlags().StringVar(&opts.HugoDataSoloOpts.version, "version", "", "version of docs and code")
	app.PersistentFlags().StringVar(&opts.HugoDataSoloOpts.product, "product", "gloo", "product to which the docs refer (defaults to gloo)")
	app.PersistentFlags().BoolVar(&opts.HugoDataSoloOpts.noScope, "no-scope", false, "if set, will not nest the served docs by product or version")
	app.PersistentFlags().BoolVar(&opts.HugoDataSoloOpts.callLatest, "call-latest", false, "if set, will use the string 'latest' in the scope, rather than the particular release version")

	return app
}
func changelogMdCmd(opts *options) *cobra.Command {
	app := &cobra.Command{
		Use:   "gen-changelog-md",
		Short: "generate a markdown file from changelogs",
		RunE: func(cmd *cobra.Command, args []string) error {

			return generateChangelogMd(opts, args)
			return nil
		},
	}
	return app
}

type HugoDataSoloYaml struct {
	// this is the string that is used to scope the hosted docs' urls
	// it can be either "latest" or a full semver version
	DocsVersion string `yaml:"DocsVersion"`
	// this is the string that tells the reader which version of code was used to publish the docs
	// it can only be a full semver version, unless user passes the --no-scope flag, which indicates "dev" mode
	CodeVersion string `yaml:"CodeVersion"`
}

const hugoDataSoloFilename = "data/Solo.yaml"

func writeVersionScopeDataForHugo(opts *options) *cobra.Command {
	app := &cobra.Command{
		Use:   "gen-version-scope-data",
		Short: "generate a data file for Hugo that indicates the docs version",
		RunE: func(cmd *cobra.Command, args []string) error {
			data := &HugoDataSoloYaml{}
			if err := getDocsVersionFromOpts(data, opts.HugoDataSoloOpts); err != nil {
				return err
			}
			marshalled, err := yaml.Marshal(data)
			if err != nil {
				return err
			}
			return ioutil.WriteFile(hugoDataSoloFilename, marshalled, 0644)
		},
	}
	return app
}

const (
	latestVersionPath = "latest"
)

func getDocsVersionFromOpts(soloData *HugoDataSoloYaml, hugoOpts HugoDataSoloOpts) error {
	if hugoOpts.noScope {
		soloData.CodeVersion = "dev"
		return nil
	}
	if hugoOpts.version == "" {
		return eris.New("must provide a version for scoped docs generation")
	}
	if hugoOpts.product == "" {
		return eris.New("must provide a product for scoped docs generation")
	}
	soloData.CodeVersion = hugoOpts.version
	version := hugoOpts.version
	if hugoOpts.callLatest {
		version = latestVersionPath
	}
	soloData.DocsVersion = fmt.Sprintf("/%v/%v", hugoOpts.product, version)
	return nil
}

const (
	glooDocGen  = "gloo"
	glooEDocGen = "glooe"
)

var (
	InvalidInputError = func(arg string) error {
		return eris.Errorf("invalid input, must provide exactly one argument, either '%v' or '%v', (provided %v)",
			glooDocGen,
			glooEDocGen,
			arg)
	}
)

func generateChangelogMd(opts *options, args []string) error {
	if len(args) != 1 {
		return InvalidInputError(fmt.Sprintf("%v", len(args)-1))
	}
	target := args[0]

	changelogDirPath := changelogutils.ChangelogDirectory
	var repoRootPath, repo string
	switch target {
	case glooDocGen:
		repoRootPath = ".."
		repo = "gloo"
	case glooEDocGen:
		// files should already be there because we download them in CI, via `download-glooe-changelog` make target
		repoRootPath = "../../solo-projects"
		repo = "solo-projects"
	default:
		return InvalidInputError(target)
	}

	// consider writing to stdout to enhance makefile/io readability `go run cmd/main.go > changelogSummary.md`
	owner := "solo-io"
	w := os.Stdout
	err := changelogutils.GenerateChangelogFromLocalDirectory(opts.ctx, repoRootPath, owner, repo, changelogDirPath, w)
	return err
}