package metadata

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GormAccess struct {
	db *gorm.DB
}

// NewGormAccess creates a new GormAccess, connects to the database.
func NewGormAccess(ctx context.Context, path string) (Access, error) {
	access := GormAccess{}
	err := access.Setup(ctx, path)
	if err != nil {
		log.Error(err, "failed to connect to the database")
		return nil, err
	}
	return &access, nil
}

// Setup connects to the database and creates the necessary tables if they don't exist
func (access *GormAccess) Setup(ctx context.Context, path string) error {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return errors.WithMessage(err, "failed to open db with gorm")
	}
	access.db = db
	return db.WithContext(ctx).AutoMigrate(&Dynakube{}, &Volume{}, &OsagentVolume{})
}

// InsertDynakube inserts a new Dynakube
func (access *GormAccess) InsertDynakube(ctx context.Context, dynakube *Dynakube) error {
	return access.db.WithContext(ctx).Create(dynakube).Error //could also be just .Save
}

// UpdateDynakube updates an existing Dynakube by matching the name
func (access *GormAccess) UpdateDynakube(ctx context.Context, dynakube *Dynakube) error {
	return access.db.WithContext(ctx).Save(dynakube).Error
}

// DeleteDynakube deletes an existing Dynakube using its name
func (access *GormAccess) DeleteDynakube(ctx context.Context, dynakubeName string) error {
	return access.db.WithContext(ctx).Where("Name = ?", dynakubeName).Delete(&Dynakube{}).Error
}

// GetDynakube gets Dynakube using its name
func (access *GormAccess) GetDynakube(ctx context.Context, dynakubeName string) (*Dynakube, error) {
	var dynakube Dynakube
	err := access.db.WithContext(ctx).Where("Name = ?", dynakubeName).Take(&dynakube).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &dynakube, err
}

// InsertVolume inserts a new Volume
func (access *GormAccess) InsertVolume(ctx context.Context, volume *Volume) error {
	return access.db.WithContext(ctx).Save(volume).Error
}

// GetVolume gets Volume by its ID
func (access *GormAccess) GetVolume(ctx context.Context, volumeID string) (*Volume, error) {
	var volume Volume
	err := access.db.WithContext(ctx).Where("ID = ?", volumeID).Take(&volume).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &volume, err
}

// DeleteVolume deletes a Volume by its ID
func (access *GormAccess) DeleteVolume(ctx context.Context, volumeID string) error {
	return access.db.WithContext(ctx).Where("ID = ?", volumeID).Delete(&Volume{}).Error
}

// InsertOsAgentVolume inserts a new OsAgentVolume
func (access *GormAccess) InsertOsAgentVolume(ctx context.Context, volume *OsagentVolume) error {
	return access.db.WithContext(ctx).Create(volume).Error //could also be just .Save
}

// UpdateOsAgentVolume updates an existing OsAgentVolume by matching the tenantUUID
func (access *GormAccess) UpdateOsAgentVolume(ctx context.Context, volume *OsagentVolume) error {
	return access.db.WithContext(ctx).Save(volume).Error
}

// GetOsAgentVolumeViaVolumeID gets an OsAgentVolume by its VolumeID
func (access *GormAccess) GetOsAgentVolumeViaVolumeID(ctx context.Context, volumeID string) (*OsagentVolume, error) {
	var volume OsagentVolume
	err := access.db.WithContext(ctx).Where("VolumeID = ?", volumeID).Take(&volume).Error
	return &volume, err
}

// GetOsAgentVolumeViaTenantUUID gets an OsAgentVolume by its tenantUUID
func (access *GormAccess) GetOsAgentVolumeViaTenantUUID(ctx context.Context, tenantUUID string) (*OsagentVolume, error) {
	var volume OsagentVolume
	err := access.db.WithContext(ctx).Where("TenantUUID = ?", tenantUUID).Take(&volume).Error
	return &volume, err
}

// GetAllVolumes gets all the Volumes from the database
func (access *GormAccess) GetAllVolumes(ctx context.Context) ([]*Volume, error) {
	var volumes []*Volume

	return volumes, access.db.WithContext(ctx).Find(&volumes).Error
}

// GetAllDynakubes gets all the Dynakubes from the database
func (access *GormAccess) GetAllDynakubes(ctx context.Context) ([]*Dynakube, error) {
	var dynakubes []*Dynakube

	return dynakubes, access.db.WithContext(ctx).Find(&dynakubes).Error
}

// GetAllOsAgentVolumes gets all the OsAgentVolume from the database
func (access *GormAccess) GetAllOsAgentVolumes(ctx context.Context) ([]*OsagentVolume, error) {
	var osVolumes []*OsagentVolume

	return osVolumes, access.db.WithContext(ctx).Find(&osVolumes).Error
}

// GetUsedVersions gets all UNIQUE versions present in the `volumes` for a given tenantUUID database in map.
// Map is used to make sure we don't return the same version multiple time,
// it's also easier to check if a version is in it or not. (a Set in style of Golang)
func (access *GormAccess) GetUsedVersions(ctx context.Context, tenantUUID string) (map[string]bool, error) {
	// TODO: Maybe less "brute force"?
	var volumes []Volume
	if err := access.db.WithContext(ctx).Where("TenantUUID = ?", tenantUUID).Find(&volumes).Error; err != nil {
		return nil, err
	}
	versions := map[string]bool{}
	for _, volume := range volumes {
		versions[volume.Version] = true
	}
	return versions, nil
}

// GetUsedVersions gets all UNIQUE versions present in the `volumes` database in map.
// Map is used to make sure we don't return the same version multiple time,
// it's also easier to check if a version is in it or not. (a Set in style of Golang)
func (access *GormAccess) GetAllUsedVersions(ctx context.Context) (map[string]bool, error) {
	// TODO: Maybe less "brute force" ?
	volumes, err := access.GetAllVolumes(ctx)
	if err != nil {
		return nil, err
	}
	versions := map[string]bool{}
	for _, volume := range volumes {
		versions[volume.Version] = true
	}
	return versions, nil
}

// GetLatestVersions gets all UNIQUE latestVersions present in the `dynakubes` database in map.
// Map is used to make sure we don't return the same version multiple time,
// it's also easier to check if a version is in it or not. (a Set in style of Golang)
func (access *GormAccess) GetLatestVersions(ctx context.Context) (map[string]bool, error) {
	// TODO: Maybe less "brute force" ?
	dynakubes, err := access.GetAllDynakubes(ctx)
	if err != nil {
		return nil, errors.WithStack(errors.WithMessage(err, "couldn't get all the latests version info for tenant uuid"))
	}
	versions := map[string]bool{}
	for _, dynakube := range dynakubes {
		versions[dynakube.LatestVersion] = true
	}
	return versions, nil
}

// GetUsedImageDigests gets all UNIQUE image digests present in the `dynakubes` database in a map.
// Map is used to make sure we don't return the same digest multiple time,
// it's also easier to check if a digest is in it or not. (a Set in style of Golang)
func (access *GormAccess) GetUsedImageDigests(ctx context.Context) (map[string]bool, error) {
	// TODO: Maybe less "brute force" ?
	dynakubes, err := access.GetAllDynakubes(ctx)
	if err != nil {
		return nil, errors.WithStack(errors.WithMessage(err, "couldn't get all the latests version info for tenant uuid"))
	}
	digests := map[string]bool{}
	for _, dynakube := range dynakubes {
		digests[dynakube.ImageDigest] = true
	}
	return digests, nil
}

// IsImageDigestUsed checks if the specified image digest is present in the database.
func (access *GormAccess) IsImageDigestUsed(ctx context.Context, imageDigest string) (bool, error) {
	var count int64
	err := access.db.WithContext(ctx).Where("ImageDigest = ?", imageDigest).Find(&Dynakube{}).Count(&count).Error
	return count > 0, err
}

// GetPodNames gets all PodNames present in the `volumes` database in map with their corresponding volumeIDs.
func (access *GormAccess) GetPodNames(ctx context.Context) (map[string]string, error) {
	// TODO: Maybe less "brute force" ?
	volumes, err := access.GetAllVolumes(ctx)
	if err != nil {
		return nil, errors.WithStack(errors.WithMessage(err, "couldn't get all pod names"))
	}
	podNames := map[string]string{}
	for _, volume := range volumes {
		if err != nil {
			return nil, errors.WithStack(errors.WithMessage(err, "couldn't scan pod name from database"))
		}
		podNames[volume.PodName] = volume.VolumeID
	}
	return nil, nil
}

// GetTenantsToDynakubes gets all Dynakubes and maps their name to the corresponding TenantUUID.
func (access *GormAccess) GetTenantsToDynakubes(ctx context.Context) (map[string]string, error) {
	// TODO: Maybe less "brute force" ?
	dynakubes, err := access.GetAllDynakubes(ctx)
	if err != nil {
		return nil, errors.WithStack(errors.WithMessage(err, "couldn't get all tenants to dynakube metadata"))
	}
	dynakubesMap := map[string]string{}
	for _, dynakube := range dynakubes {
		dynakubesMap[dynakube.Name] = dynakube.TenantUUID
	}
	return nil, nil
}
