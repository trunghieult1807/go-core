package registry

import (
	"sync"

	"github.com/sarulabs/di"
)

// DIBuilder public method to generate  definition for building DI
type DIBuilder func() []di.Def

var (
	buildOnce sync.Once
	builder   *di.Builder
	container di.Container
	// ConfigsBuilder builder for config
	ConfigsBuilder DIBuilder
	// HelpersBuilder builder for all helpers
	HelpersBuilder DIBuilder
	// RepositoriesBuilder builder for repositories
	RepositoriesBuilder DIBuilder
	// AdaptersBuilder builder for adapters
	AdaptersBuilder DIBuilder
	// UsecasesBuilder builder for usecase
	UsecasesBuilder DIBuilder
	// FacadiesBuilder builder for usecase
	FacadiesBuilder DIBuilder
	// APIsBuilder builder for apis
	APIsBuilder DIBuilder
)

// BuildDIContainer build DI container
func BuildDIContainer() {
	buildOnce.Do(func() {
		builder, _ = di.NewBuilder()
		doBuild()
		container = builder.Build()
	})
}

func doBuild() {
	if err := buildConfigs(); err != nil {
		panic(err)
	}
	if err := buildHelpers(); err != nil {
		panic(err)
	}
	if err := buildRepositories(); err != nil {
		panic(err)
	}
	if err := buildAdapters(); err != nil {
		panic(err)
	}
	if err := buildUsecases(); err != nil {
		panic(err)
	}
	if err := buildFacadies(); err != nil {
		panic(err)
	}
	if err := buildAPIs(); err != nil {
		panic(err)
	}
}

// GetDependency gets dependency from DI container
func GetDependency(dependencyName string) interface{} {
	return container.Get(dependencyName)
}

// CleanDependency cleans dependency
func CleanDependency() error {
	return container.Clean()
}

func buildConfigs() error {
	defs := []di.Def{}
	if ConfigsBuilder == nil {
		ConfigsBuilder = defaultBuilder
	}
	defs = ConfigsBuilder()
	err := builder.Add(defs...)
	if err != nil {
		return err
	}
	return nil
}

func buildAPIs() error {
	defs := []di.Def{}
	if APIsBuilder == nil {
		APIsBuilder = defaultBuilder
	}
	defs = APIsBuilder()
	err := builder.Add(defs...)
	if err != nil {
		return err
	}
	return nil
}

func buildUsecases() error {
	defs := []di.Def{}
	if UsecasesBuilder == nil {
		UsecasesBuilder = defaultBuilder
	}
	defs = UsecasesBuilder()
	err := builder.Add(defs...)
	if err != nil {
		return err
	}
	return nil
}

func buildFacadies() error {
	defs := []di.Def{}
	if FacadiesBuilder == nil {
		FacadiesBuilder = defaultBuilder
	}
	defs = FacadiesBuilder()
	err := builder.Add(defs...)
	if err != nil {
		return err
	}
	return nil
}

func buildRepositories() error {
	defs := []di.Def{}
	if RepositoriesBuilder == nil {
		RepositoriesBuilder = defaultBuilder
	}
	defs = RepositoriesBuilder()
	err := builder.Add(defs...)
	if err != nil {
		return err
	}
	return nil
}

func buildAdapters() error {
	defs := []di.Def{}
	if AdaptersBuilder == nil {
		AdaptersBuilder = defaultBuilder
	}
	defs = AdaptersBuilder()
	err := builder.Add(defs...)
	if err != nil {
		return err
	}
	return nil
}

func buildHelpers() error {
	defs := []di.Def{}
	if HelpersBuilder == nil {
		HelpersBuilder = defaultBuilder
	}
	defs = HelpersBuilder()
	err := builder.Add(defs...)
	if err != nil {
		return err
	}
	return nil
}

func defaultBuilder() []di.Def {
	return []di.Def{}
}
