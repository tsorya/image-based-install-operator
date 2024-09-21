package installer_types

import (
	lca_api "github.com/openshift-kni/lifecycle-agent/api/seedreconfig"
	apicfgv1 "github.com/openshift/api/config/v1"
)

type InstallConfig struct {
	APIVersion string         `json:"apiVersion"`
	BaseDomain string         `json:"baseDomain"`
	Proxy      *lca_api.Proxy `json:"proxy,omitempty"`
	Networking struct {
		NetworkType    string           `json:"networkType"`
		MachineNetwork []MachineNetwork `json:"machineNetwork,omitempty"`
	} `json:"networking"`
	Metadata              Metadata            `json:"metadata"`
	Compute               []Compute           `json:"compute"`
	ControlPlane          Compute             `json:"controlPlane"`
	Platform              Platform            `json:"platform"`
	FIPS                  bool                `json:"fips"`
	PullSecret            string              `json:"pullSecret"`
	SSHKey                string              `json:"sshKey"`
	AdditionalTrustBundle string              `json:"additionalTrustBundle,omitempty"`
	ImageDigestSources    []ImageDigestSource `json:"imageDigestSources,omitempty"`
}

type Networking struct {
	NetworkType    string           `json:"networkType"`
	MachineNetwork []MachineNetwork `json:"machineNetwork,omitempty"`
}

type MachineNetwork struct {
	Cidr string `json:"cidr"`
}

type Platform struct {
	None PlatformNone `json:"none,omitempty"`
}

type PlatformNone struct {
}

type Metadata struct {
	Name string `json:"name"`
}

type Compute struct {
	Name     string `json:"name"`
	Replicas int    `json:"replicas"`
}

type ImageDigestSource struct {
	// Source is the repository that users refer to, e.g. in image pull specifications.
	Source string `json:"source"`

	// Mirrors is one or more repositories that may also contain the same images.
	Mirrors []apicfgv1.ImageMirror `json:"mirrors,omitempty"`
}
