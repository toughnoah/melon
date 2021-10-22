package validationwebhooks

const (
	namingCheckError = "naming check failed: %s"

	decodeError = "decode error: %s"

	noContainerError = "deployment has no container specified"

	noResourcesLimitsError = "the resources limits of deployment can not be null"

	deniedErrorMessage = `not match the expr ^(?:noah|blackbean|melon)-(?:dev|qa|sa)-.+?-(?:test|prod)`
)
