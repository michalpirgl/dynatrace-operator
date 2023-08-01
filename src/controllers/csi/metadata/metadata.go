package metadata

import (
	"context"
	"time"
)

// Dynakube stores the necessary info from the Dynakube that is needed to be used during volume mount/unmount.
type Dynakube struct {
	Name                   string `json:"name" gorm:"column:Name;primaryKey;type:VARCHAR NOT NULL"`
	TenantUUID             string `json:"tenantUUID" gorm:"column:TenantUUID;type:VARCHAR NOT NULL"`
	LatestVersion          string `json:"latestVersion" gorm:"column:LatestVersion;type:VARCHAR NOT NULL"`
	ImageDigest            string `json:"imageDigest" gorm:"column:ImageDigest;type:VARCHAR NOT NULL DEFAULT ''"`
	MaxFailedMountAttempts int    `json:"maxFailedMountAttempts" gorm:"column:MaxFailedMountAttempts;type:VARCHAR NOT NULL;default:10"`
}

// NewDynakube returns a new metadata.Dynakube if all fields are set.
func NewDynakube(dynakubeName, tenantUUID, latestVersion, imageDigest string, maxFailedMountAttempts int) *Dynakube { //nolint:revive // argument-limit doesn't apply to constructors
	if tenantUUID == "" || dynakubeName == "" {
		return nil
	}

	return &Dynakube{
		Name:                   dynakubeName,
		TenantUUID:             tenantUUID,
		LatestVersion:          latestVersion,
		ImageDigest:            imageDigest,
		MaxFailedMountAttempts: maxFailedMountAttempts,
	}
}

type Volume struct {
	VolumeID      string `json:"volumeID" gorm:"column:ID;primaryKey;type:VARCHAR NOT NULL"`
	PodName       string `json:"podName" gorm:"column:PodName;type:VARCHAR NOT NULL"`
	Version       string `json:"version" gorm:"column:Version;type:VARCHAR NOT NULL"`
	TenantUUID    string `json:"tenantUUID" gorm:"column:TenantUUID;type:VARCHAR NOT NULL"`
	MountAttempts int    `json:"mountAttempts" gorm:"column:MountAttempts;type:VARCHAR NOT NULL;default:0"`
}

// NewVolume returns a new Volume if all fields (except version) are set.
func NewVolume(id, podName, version, tenantUUID string, mountAttempts int) *Volume { //nolint:revive // argument-limit doesn't apply to constructors
	if id == "" || podName == "" || tenantUUID == "" {
		return nil
	}

	if mountAttempts < 0 {
		mountAttempts = 0
	}

	return &Volume{
		VolumeID:      id,
		PodName:       podName,
		Version:       version,
		TenantUUID:    tenantUUID,
		MountAttempts: mountAttempts,
	}
}

type OsagentVolume struct {
	VolumeID     string     `json:"volumeID" gorm:"column:VolumeID;type:VARCHAR NOT NULL"`
	TenantUUID   string     `json:"tenantUUID" gorm:"column:TenantUUID;primaryKey;type:VARCHAR NOT NULL"`
	Mounted      bool       `json:"mounted" gorm:"column:Mounted;type:BOOLEAN NOT NULL"`
	LastModified *time.Time `json:"lastModified" gorm:"column:LastModified;type:DATETIME NOT NULL"`
}

// NewOsAgentVolume returns a new volume if all fields are set.
func NewOsAgentVolume(volumeID, tenantUUID string, mounted bool, timeStamp *time.Time) *OsagentVolume {
	if volumeID == "" || tenantUUID == "" || timeStamp == nil {
		return nil
	}
	return &OsagentVolume{
		VolumeID:     volumeID,
		TenantUUID:   tenantUUID,
		Mounted:      mounted,
		LastModified: timeStamp,
	}
}

type Access interface {
	Setup(ctx context.Context, path string) error

	InsertDynakube(ctx context.Context, dynakube *Dynakube) error
	UpdateDynakube(ctx context.Context, dynakube *Dynakube) error
	DeleteDynakube(ctx context.Context, dynakubeName string) error
	GetDynakube(ctx context.Context, dynakubeName string) (*Dynakube, error)
	GetTenantsToDynakubes(ctx context.Context) (map[string]string, error)
	GetAllDynakubes(ctx context.Context) ([]*Dynakube, error)

	InsertOsAgentVolume(ctx context.Context, volume *OsagentVolume) error
	GetOsAgentVolumeViaVolumeID(ctx context.Context, volumeID string) (*OsagentVolume, error)
	GetOsAgentVolumeViaTenantUUID(ctx context.Context, volumeID string) (*OsagentVolume, error)
	UpdateOsAgentVolume(ctx context.Context, volume *OsagentVolume) error
	GetAllOsAgentVolumes(ctx context.Context) ([]*OsagentVolume, error)

	InsertVolume(ctx context.Context, volume *Volume) error
	DeleteVolume(ctx context.Context, volumeID string) error
	GetVolume(ctx context.Context, volumeID string) (*Volume, error)
	GetAllVolumes(ctx context.Context) ([]*Volume, error)
	GetPodNames(ctx context.Context) (map[string]string, error)
	GetUsedVersions(ctx context.Context, tenantUUID string) (map[string]bool, error)
	GetAllUsedVersions(ctx context.Context) (map[string]bool, error)
	GetLatestVersions(ctx context.Context) (map[string]bool, error)
	GetUsedImageDigests(ctx context.Context) (map[string]bool, error)
	IsImageDigestUsed(ctx context.Context, imageDigest string) (bool, error)
}

type AccessOverview struct {
	Volumes        []*Volume        `json:"volumes"`
	Dynakubes      []*Dynakube      `json:"dynakubes"`
	OsAgentVolumes []*OsagentVolume `json:"osAgentVolumes"`
}

func NewAccessOverview(access Access) (*AccessOverview, error) {
	ctx := context.Background()
	volumes, err := access.GetAllVolumes(ctx)
	if err != nil {
		return nil, err
	}
	dynakubes, err := access.GetAllDynakubes(ctx)
	if err != nil {
		return nil, err
	}
	osVolumes, err := access.GetAllOsAgentVolumes(ctx)
	if err != nil {
		return nil, err
	}
	return &AccessOverview{
		Volumes:        volumes,
		Dynakubes:      dynakubes,
		OsAgentVolumes: osVolumes,
	}, nil
}

func LogAccessOverview(access Access) {
	overview, err := NewAccessOverview(access)
	if err != nil {
		log.Error(err, "Failed to get an overview of the stored csi metadata")
	}
	log.Info("The current overview of the csi metadata", "overview", overview)
}
