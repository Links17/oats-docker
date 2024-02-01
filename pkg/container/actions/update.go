package actions

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"oats-docker/pkg/container"
	"oats-docker/pkg/container/lifecycle"
	"oats-docker/pkg/container/session"
	"oats-docker/pkg/container/sorter"
	"oats-docker/pkg/types"
	"oats-docker/pkg/util"
	"strings"
)

// Update looks at the running Docker containers to see if any of the images
// used to start those containers have been updated. If a change is detected in
// any of the images, the associated containers are stopped and restarted with
// the new image.
func Update(client container.Client, params types.UpdateParams) (types.Report, error) {
	log.Debug("Checking containers for updated images")
	progress := &session.Progress{}
	staleCount := 0

	if params.LifecycleHooks {
		lifecycle.ExecutePreChecks(client, params)
	}

	containers, err := client.ListContainers()
	if err != nil {
		return nil, err
	}

	staleCheckFailed := 0

	for i, targetContainer := range containers {
		stale, newestImage, err := client.IsContainerStale(targetContainer, params)
		shouldUpdate := stale && !params.NoRestart && !targetContainer.IsMonitorOnly(params)
		if err == nil && shouldUpdate {
			// Check to make sure we have all the necessary information for recreating the container
			err = targetContainer.VerifyConfiguration()
			// If the image information is incomplete and trace logging is enabled, log it for further diagnosis
			if err != nil && log.IsLevelEnabled(log.TraceLevel) {
				imageInfo := targetContainer.ImageInfo()
				log.Tracef("Image info: %#v", imageInfo)
				log.Tracef("Container info: %#v", targetContainer.ContainerInfo())
				if imageInfo != nil {
					log.Tracef("Image config: %#v", imageInfo.Config)
				}
			}
		}

		if err != nil {
			log.Infof("Unable to update container %q: %v. Proceeding to next.", targetContainer.Name(), err)
			stale = false
			staleCheckFailed++
			progress.AddSkipped(targetContainer, err)
		} else {
			progress.AddScanned(targetContainer, newestImage)
		}
		containers[i].SetStale(stale)

		if stale {
			staleCount++
		}
	}

	containers, err = sorter.SortByDependencies(containers)
	if err != nil {
		return nil, err
	}

	UpdateImplicitRestart(containers)

	var containersToUpdate []types.Container
	for _, c := range containers {
		if !c.IsMonitorOnly(params) {
			name := c.Name()
			tag := util.RepoTag(name, params)
			c.ContainerInfo().Config.Image = tag
			if tag == "" {
				continue
			}
			c.ImageInfo().RepoTags = []string{tag}
			id, err := client.ImageID(tag)
			if err != nil {
				continue
			}
			c.ContainerInfo().Image = id
			containersToUpdate = append(containersToUpdate, c)
			progress.MarkForUpdate(c.ID())
		}
	}

	if params.RollingRestart {
		progress.UpdateFailed(PerformRollingRestart(containersToUpdate, client, params))
	} else {
		failedStop, stoppedImages := stopContainersInReversedOrder(containersToUpdate, client, params)
		progress.UpdateFailed(failedStop)
		failedStart := restartContainersInSortedOrder(containersToUpdate, client, params, stoppedImages)
		progress.UpdateFailed(failedStart)
	}

	if params.LifecycleHooks {
		lifecycle.ExecutePostChecks(client, params)
	}
	return progress.Report(), nil
}

func PerformRollingRestart(containers []types.Container, client container.Client, params types.UpdateParams) map[types.ContainerID]error {
	cleanupImageIDs := make(map[types.ImageID]bool, len(containers))
	failed := make(map[types.ContainerID]error, len(containers))

	for i := len(containers) - 1; i >= 0; i-- {
		if containers[i].ToRestart() {
			err := stopStaleContainer(containers[i], client, params)
			if err != nil {
				failed[containers[i].ID()] = err
			} else {
				if err := restartStaleContainer(containers[i], client, params); err != nil {
					failed[containers[i].ID()] = err
				} else if containers[i].IsStale() {
					// Only add (previously) stale containers' images to clean up
					cleanupImageIDs[containers[i].ImageID()] = true
				}
			}
		}
	}

	if params.Cleanup {
		cleanupImages(client, cleanupImageIDs)
	}
	return failed
}

func UpdateEnv(containers []types.Container, client container.Client, params types.UpdateParams) map[types.ContainerID]error {
	cleanupImageIDs := make(map[types.ImageID]bool, len(containers))
	failed := make(map[types.ContainerID]error, len(containers))

	for i := len(containers) - 1; i >= 0; i-- {
		if err := updateContainer(containers[i], client, params); err != nil {
			failed[containers[i].ID()] = err
		} else if containers[i].IsStale() {
			// Only add (previously) stale containers' images to clean up
			cleanupImageIDs[containers[i].ImageID()] = true
		}
	}

	if params.Cleanup {
		cleanupImages(client, cleanupImageIDs)
	}
	return failed
}

func stopContainersInReversedOrder(containers []types.Container, client container.Client, params types.UpdateParams) (failed map[types.ContainerID]error, stopped map[types.ImageID]bool) {
	failed = make(map[types.ContainerID]error, len(containers))
	stopped = make(map[types.ImageID]bool, len(containers))
	for i := len(containers) - 1; i >= 0; i-- {
		if err := stopStaleContainer(containers[i], client, params); err != nil {
			failed[containers[i].ID()] = err
		} else {
			// NOTE: If a container is restarted due to a dependency this might be empty
			stopped[containers[i].SafeImageID()] = true
		}

	}
	return
}

func stopStaleContainer(container types.Container, client container.Client, params types.UpdateParams) error {
	if container.IsWatchtower() {
		log.Debugf("This is the watchtower container %s", container.Name())
		return nil
	}

	if !container.ToRestart() {
		return nil
	}

	// Perform an additional check here to prevent us from stopping a linked container we cannot restart
	if container.IsLinkedToRestarting() {
		if err := container.VerifyConfiguration(); err != nil {
			return err
		}
	}

	if params.LifecycleHooks {
		skipUpdate, err := lifecycle.ExecutePreUpdateCommand(client, container)
		if err != nil {
			log.Error(err)
			log.Info("Skipping container as the pre-update command failed")
			return err
		}
		if skipUpdate {
			log.Debug("Skipping container as the pre-update command returned exit code 75 (EX_TEMPFAIL)")
			return errors.New("skipping container as the pre-update command returned exit code 75 (EX_TEMPFAIL)")
		}
	}

	if err := client.StopContainer(container, params.Timeout); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func restartContainersInSortedOrder(containers []types.Container, client container.Client, params types.UpdateParams, stoppedImages map[types.ImageID]bool) map[types.ContainerID]error {
	cleanupImageIDs := make(map[types.ImageID]bool, len(containers))
	failed := make(map[types.ContainerID]error, len(containers))

	for _, c := range containers {
		if !c.ToRestart() {
			continue
		}
		if stoppedImages[c.SafeImageID()] {
			if err := restartStaleContainer(c, client, params); err != nil {
				failed[c.ID()] = err
			} else if c.IsStale() {
				// Only add (previously) stale containers' images to cleanup
				cleanupImageIDs[c.ImageID()] = true
			}
		}
	}

	if params.Cleanup {
		cleanupImages(client, cleanupImageIDs)
	}

	return failed
}

func cleanupImages(client container.Client, imageIDs map[types.ImageID]bool) {
	for imageID := range imageIDs {
		if imageID == "" {
			continue
		}
		if err := client.RemoveImageByID(imageID); err != nil {
			log.Error(err)
		}
	}
}

func restartStaleContainer(container types.Container, client container.Client, params types.UpdateParams) error {
	// Since we can't shutdown a watchtower container immediately, we need to
	// start the new one while the old one is still running. This prevents us
	// from re-using the same container name so we first rename the current
	// instance so that the new one can adopt the old name.
	if container.IsWatchtower() {
		if err := client.RenameContainer(container, util.RandName()); err != nil {
			log.Error(err)
			return nil
		}
	}
	if !params.NoRestart {
		if newContainerID, err := client.StartContainer(container); err != nil {
			log.Error(err)
			return err
		} else if container.ToRestart() && params.LifecycleHooks {
			lifecycle.ExecutePostUpdateCommand(client, newContainerID)
		}
	}
	return nil
}

func updateContainer(container types.Container, client container.Client, params types.UpdateParams) error {
	// Since we can't shutdown a watchtower container immediately, we need to
	// start the new one while the old one is still running. This prevents us
	// from re-using the same container name so we first rename the current
	// instance so that the new one can adopt the old name.
	if container.IsWatchtower() {
		if err := client.RenameContainer(container, util.RandName()); err != nil {
			log.Error(err)
			return nil
		}
	}
	check := CheckEnv(container, params)
	if check && !params.NoRestart {
		if err := client.RemoveContainer(container, params.Timeout); err != nil {
			log.Error(err)
			return err
		}
		if newContainerID, err := client.StartContainer(container); err != nil {
			log.Error(err)
			return err
		} else if container.ToRestart() && params.LifecycleHooks {
			lifecycle.ExecutePostUpdateCommand(client, newContainerID)
		}
	}
	return nil
}

// CheckEnv checkEnv update value and update status
func CheckEnv(container types.Container, params types.UpdateParams) bool {
	existingEnv := make(map[string]bool)
	changed := false

	for i, env := range container.ContainerInfo().Config.Env {
		existingEnv[env] = true
		for _, updateEnv := range params.UpdateEnv {
			u := strings.Split(updateEnv, "=")
			fmt.Println(u)
			e := strings.Split(env, "=")
			fmt.Println(e)
			if len(u) > 1 && len(e) > 1 && strings.HasPrefix(u[0], e[0]) && !strings.HasPrefix(u[1], e[1]) {
				container.ContainerInfo().Config.Env[i] = updateEnv
				existingEnv[updateEnv] = true
				changed = true
				break
			}
			if len(u) == 1 && len(e) > 1 && strings.HasPrefix(u[0], e[0]) {
				container.ContainerInfo().Config.Env = append(container.ContainerInfo().Config.Env[:i], container.ContainerInfo().Config.Env[i+1:]...)
				changed = true
				break
			}
		}
	}

	for _, updateEnv := range params.UpdateEnv {
		i := len(strings.Split(updateEnv, "="))
		if !existingEnv[updateEnv] && i != 1 {
			container.ContainerInfo().Config.Env = append(container.ContainerInfo().Config.Env, updateEnv)
			existingEnv[updateEnv] = true
			changed = true
		}
	}
	return changed
}

// UpdateImplicitRestart iterates through the passed containers, setting the
// `LinkedToRestarting` flag if any of it's linked containers are marked for restart
func UpdateImplicitRestart(containers []types.Container) {

	for ci, c := range containers {
		if c.ToRestart() {
			// The container is already marked for restart, no need to check
			continue
		}

		if link := linkedContainerMarkedForRestart(c.Links(), containers); link != "" {
			log.WithFields(log.Fields{
				"restarting": link,
				"linked":     c.Name(),
			}).Debug("container is linked to restarting")
			// NOTE: To mutate the array, the `c` variable cannot be used as it's a copy
			containers[ci].SetLinkedToRestarting(true)
		}

	}
}

// linkedContainerMarkedForRestart returns the name of the first link that matches a
// container marked for restart
func linkedContainerMarkedForRestart(links []string, containers []types.Container) string {
	for _, linkName := range links {
		for _, candidate := range containers {
			if candidate.Name() == linkName && candidate.ToRestart() {
				return linkName
			}
		}
	}
	return ""
}
