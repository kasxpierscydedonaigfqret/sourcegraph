package graphqlbackend

import (
	"context"

	"github.com/sourcegraph/sourcegraph/cmd/frontend/internal/backend"
)

func (r *schemaResolver) ExtensionRegistry(ctx context.Context) (*extensionRegistryResolver, error) {
	if err := backend.CheckActorHasPlatformEnabled(ctx); err != nil {
		return nil, err
	}
	return &extensionRegistryResolver{}, nil
}

type extensionRegistryResolver struct{}

func (*extensionRegistryResolver) Extension(ctx context.Context, args *struct {
	ExtensionID string
}) (*registryExtensionMultiResolver, error) {
	return getExtensionByExtensionID(ctx, args.ExtensionID)
}

func getExtensionByExtensionID(ctx context.Context, extensionID string) (*registryExtensionMultiResolver, error) {
	local, remote, err := backend.GetExtensionByExtensionID(ctx, extensionID)
	if err != nil {
		return nil, err
	}
	if local != nil {
		return &registryExtensionMultiResolver{local: &registryExtensionDBResolver{local}}, nil
	}
	return &registryExtensionMultiResolver{remote: &registryExtensionRemoteResolver{remote}}, nil
}

func (*extensionRegistryResolver) LocalExtensionIDPrefix() (*string, error) {
	return backend.GetLocalRegistryExtensionIDPrefix()
}
