package installer_types

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	// ImageBasedConfigVersion is the version supported by this package.
	ImageBasedConfigVersion = "v1beta1"
)

// Config is the API for specifying configuration for the image-based configuration ISO.
// image-based-config.yaml
type ImageBasedConfig struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        `json:"metadata,omitempty"`

	// AdditionalNTPSources is a list of NTP sources (hostname or IP) to be added to all cluster
	// hosts. They are added to any NTP sources that were configured through other means.
	// +optional
	AdditionalNTPSources []string `json:"additionalNTPSources,omitempty"`

	// Hostname is the desired hostname of the SNO node.
	Hostname string `json:"hostname,omitempty"`

	// NetworkConfig is a YAML manifest that can be processed by nmstate, using custom
	// marshaling/unmarshaling that will allow to populate nmstate config as plain yaml.
	// +optional
	NetworkConfig string `json:"networkConfig,omitempty"`

	// ReleaseRegistry is the container registry used to host the release image of the seed cluster.
	// +optional
	ReleaseRegistry string `json:"releaseRegistry,omitempty"`
}
