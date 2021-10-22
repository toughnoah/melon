package validationwebhooks

const (
	namingCheckError = "%s naming check failed: %s"

	decodeError = "decode error: %s"

	noContainerError = "deployment has no container specified"

	noResourcesLimitsError = "the resources limits of deployment can not be null"

	testDeniedErrorMessage = `not match the expr ^(?:noah|blackbean|melon)-(?:dev|qa|sa)-.+?-(?:test|prod)`

	testImageNamingFailedMessage = `not match the expr ^(?:docker.io)/(?:toughnoah|test)/.+?:v1.0`

	DeploymentNamingKind = "deployment.naming"

	ImageKind = "image"

	DeploymentLimitsKind = "deployment.limits"

	NamespaceNamingKind = "namespace.naming"

	ServiceNamingKind = "service.naming"

	ConfigmapNamingKind = "configmap.naming"

	DeploymentLimitRangeKind = "deployment.resources.limitRangeSpec"

	DeploymentResourceQuatoKind = "deployment.resources.resourceQuato"
)
