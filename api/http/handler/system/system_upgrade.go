package system

import (
	"net/http"
	"regexp"
	"slices"

	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/internal/upgrade"
	"github.com/portainer/portainer/api/platform"
	plf "github.com/portainer/portainer/api/platform"
	httperror "github.com/portainer/portainer/pkg/libhttp/error"
	"github.com/portainer/portainer/pkg/libhttp/request"
	"github.com/portainer/portainer/pkg/libhttp/response"

	"github.com/pkg/errors"
)

type systemUpgradePayload struct {
	License string
}

var re = regexp.MustCompile(`^\d-.+`)

func (payload *systemUpgradePayload) Validate(r *http.Request) error {
	if payload.License == "" {
		return errors.New("license is missing")
	}

	if !re.MatchString(payload.License) {
		return errors.New("license is invalid")
	}

	return nil
}

var platformToEndpointType = map[platform.ContainerPlatform][]portainer.EndpointType{
	platform.PlatformDocker:     {portainer.AgentOnDockerEnvironment, portainer.DockerEnvironment},
	platform.PlatformKubernetes: {portainer.KubernetesLocalEnvironment},
}

// @id systemUpgrade
// @summary Upgrade Portainer to BE
// @description Upgrade Portainer to BE
// @description **Access policy**: administrator
// @tags system
// @produce json
// @success 204 {object} status "Success"
// @router /system/upgrade [post]
func (handler *Handler) systemUpgrade(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	payload, err := request.GetPayload[systemUpgradePayload](r)
	if err != nil {
		return httperror.BadRequest("Invalid request payload", err)
	}

	environment, err := handler.guessLocalEndpoint()
	if err != nil {
		return httperror.InternalServerError("Failed to guess local endpoint", err)
	}

	err = handler.upgradeService.Upgrade(environment, payload.License)
	if err != nil {
		return httperror.InternalServerError("Failed to upgrade Portainer", err)
	}

	return response.Empty(w)
}

func (handler *Handler) guessLocalEndpoint() (*portainer.Endpoint, error) {
	platform, err := plf.DetermineContainerPlatform()
	if err != nil {
		return nil, errors.Wrap(err, "failed to determine container platform")
	}

	endpoints, err := handler.dataStore.Endpoint().Endpoints()
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve endpoints")
	}

	endpointTypes, ok := platformToEndpointType[platform]
	if !ok {
		return nil, errors.New("failed to determine endpoint type")
	}

	for _, endpoint := range endpoints {
		if slices.Contains(endpointTypes, endpoint.Type) {
			if platform != plf.PlatformDocker || upgrade.CheckDockerEnvTypeForUpgrade(&endpoint) != "" {
				return &endpoint, nil
			}
		}
	}

	return nil, errors.New("failed to find local endpoint")
}
