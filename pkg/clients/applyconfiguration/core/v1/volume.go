// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus", DO NOT EDIT.

package v1

// VolumeApplyConfiguration represents an declarative configuration of the Volume type for use
// with apply.
type VolumeApplyConfiguration struct {
	Name                           *string `json:"name,omitempty"`
	VolumeSourceApplyConfiguration `json:",inline"`
}

// VolumeApplyConfiguration constructs an declarative configuration of the Volume type for use with
// apply.
func Volume() *VolumeApplyConfiguration {
	return &VolumeApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithName(value string) *VolumeApplyConfiguration {
	b.Name = &value
	return b
}

// WithHostPath sets the HostPath field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the HostPath field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithHostPath(value *HostPathVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.HostPath = value
	return b
}

// WithEmptyDir sets the EmptyDir field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the EmptyDir field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithEmptyDir(value *EmptyDirVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.EmptyDir = value
	return b
}

// WithGCEPersistentDisk sets the GCEPersistentDisk field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the GCEPersistentDisk field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithGCEPersistentDisk(value *GCEPersistentDiskVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.GCEPersistentDisk = value
	return b
}

// WithAWSElasticBlockStore sets the AWSElasticBlockStore field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the AWSElasticBlockStore field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithAWSElasticBlockStore(value *AWSElasticBlockStoreVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.AWSElasticBlockStore = value
	return b
}

// WithGitRepo sets the GitRepo field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the GitRepo field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithGitRepo(value *GitRepoVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.GitRepo = value
	return b
}

// WithSecret sets the Secret field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Secret field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithSecret(value *SecretVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.Secret = value
	return b
}

// WithNFS sets the NFS field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the NFS field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithNFS(value *NFSVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.NFS = value
	return b
}

// WithISCSI sets the ISCSI field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ISCSI field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithISCSI(value *ISCSIVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.ISCSI = value
	return b
}

// WithGlusterfs sets the Glusterfs field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Glusterfs field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithGlusterfs(value *GlusterfsVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.Glusterfs = value
	return b
}

// WithPersistentVolumeClaim sets the PersistentVolumeClaim field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PersistentVolumeClaim field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithPersistentVolumeClaim(value *PersistentVolumeClaimVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.PersistentVolumeClaim = value
	return b
}

// WithRBD sets the RBD field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RBD field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithRBD(value *RBDVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.RBD = value
	return b
}

// WithFlexVolume sets the FlexVolume field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the FlexVolume field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithFlexVolume(value *FlexVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.FlexVolume = value
	return b
}

// WithCinder sets the Cinder field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Cinder field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithCinder(value *CinderVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.Cinder = value
	return b
}

// WithCephFS sets the CephFS field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CephFS field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithCephFS(value *CephFSVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.CephFS = value
	return b
}

// WithFlocker sets the Flocker field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Flocker field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithFlocker(value *FlockerVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.Flocker = value
	return b
}

// WithDownwardAPI sets the DownwardAPI field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DownwardAPI field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithDownwardAPI(value *DownwardAPIVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.DownwardAPI = value
	return b
}

// WithFC sets the FC field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the FC field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithFC(value *FCVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.FC = value
	return b
}

// WithAzureFile sets the AzureFile field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the AzureFile field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithAzureFile(value *AzureFileVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.AzureFile = value
	return b
}

// WithConfigMap sets the ConfigMap field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ConfigMap field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithConfigMap(value *ConfigMapVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.ConfigMap = value
	return b
}

// WithVsphereVolume sets the VsphereVolume field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the VsphereVolume field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithVsphereVolume(value *VsphereVirtualDiskVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.VsphereVolume = value
	return b
}

// WithQuobyte sets the Quobyte field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Quobyte field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithQuobyte(value *QuobyteVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.Quobyte = value
	return b
}

// WithAzureDisk sets the AzureDisk field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the AzureDisk field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithAzureDisk(value *AzureDiskVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.AzureDisk = value
	return b
}

// WithPhotonPersistentDisk sets the PhotonPersistentDisk field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PhotonPersistentDisk field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithPhotonPersistentDisk(value *PhotonPersistentDiskVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.PhotonPersistentDisk = value
	return b
}

// WithProjected sets the Projected field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Projected field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithProjected(value *ProjectedVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.Projected = value
	return b
}

// WithPortworxVolume sets the PortworxVolume field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PortworxVolume field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithPortworxVolume(value *PortworxVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.PortworxVolume = value
	return b
}

// WithScaleIO sets the ScaleIO field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ScaleIO field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithScaleIO(value *ScaleIOVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.ScaleIO = value
	return b
}

// WithStorageOS sets the StorageOS field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the StorageOS field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithStorageOS(value *StorageOSVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.StorageOS = value
	return b
}

// WithCSI sets the CSI field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CSI field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithCSI(value *CSIVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.CSI = value
	return b
}

// WithEphemeral sets the Ephemeral field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Ephemeral field is set to the value of the last call.
func (b *VolumeApplyConfiguration) WithEphemeral(value *EphemeralVolumeSourceApplyConfiguration) *VolumeApplyConfiguration {
	b.Ephemeral = value
	return b
}
