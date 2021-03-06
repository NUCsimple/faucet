/*
 * Kubernetes
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: v1.10.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package client

// Represents storage that is managed by an external CSI volume driver (Beta feature)
type V1CsiPersistentVolumeSource struct {

	// ControllerPublishSecretRef is a reference to the secret object containing sensitive information to pass to the CSI driver to complete the CSI ControllerPublishVolume and ControllerUnpublishVolume calls. This field is optional, and  may be empty if no secret is required. If the secret object contains more than one secret, all secrets are passed.
	ControllerPublishSecretRef *V1SecretReference `json:"controllerPublishSecretRef,omitempty"`

	// Driver is the name of the driver to use for this volume. Required.
	Driver string `json:"driver"`

	// Filesystem type to mount. Must be a filesystem type supported by the host operating system. Ex. \"ext4\", \"xfs\", \"ntfs\". Implicitly inferred to be \"ext4\" if unspecified.
	FsType string `json:"fsType,omitempty"`

	// NodePublishSecretRef is a reference to the secret object containing sensitive information to pass to the CSI driver to complete the CSI NodePublishVolume and NodeUnpublishVolume calls. This field is optional, and  may be empty if no secret is required. If the secret object contains more than one secret, all secrets are passed.
	NodePublishSecretRef *V1SecretReference `json:"nodePublishSecretRef,omitempty"`

	// NodeStageSecretRef is a reference to the secret object containing sensitive information to pass to the CSI driver to complete the CSI NodeStageVolume and NodeStageVolume and NodeUnstageVolume calls. This field is optional, and  may be empty if no secret is required. If the secret object contains more than one secret, all secrets are passed.
	NodeStageSecretRef *V1SecretReference `json:"nodeStageSecretRef,omitempty"`

	// Optional: The value to pass to ControllerPublishVolumeRequest. Defaults to false (read/write).
	ReadOnly bool `json:"readOnly,omitempty"`

	// Attributes of the volume to publish.
	VolumeAttributes map[string]string `json:"volumeAttributes,omitempty"`

	// VolumeHandle is the unique volume name returned by the CSI volume plugin???s CreateVolume to refer to the volume on all subsequent calls. Required.
	VolumeHandle string `json:"volumeHandle"`
}
