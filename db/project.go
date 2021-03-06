package db

import "errors"

func ProjectExists(projectID uint) bool {
	doesNotExist := DB.Where("id = ?", projectID).Find(&Project{}).RecordNotFound()
	return !doesNotExist
}

func ProjectHasEntries(projectID uint) bool {
	if !ProjectExists(projectID) {
		return false
	}
	project, _ := GetProjectByID(projectID)
	return len(project.Entries) > 0
}

func ProjectIsOwnedByUser(projectID uint, userID uint) (bool, error) {
	var project Project
	err := DB.Where("id = ?", projectID).Find(&project).Error
	if err != nil {
		return false, err
	}
	return project.UserID == userID, nil
}

func GetProjectByID(projectID uint) (Project, error) {
	var project Project
	err := DB.Where("id = ?", projectID).Preload("Entries").Preload("EntryTypes").Find(&project).Error
	if err != nil {
		return project, err
	}
	return project, nil
}

func ProjectSupportsEntryType(projectID uint, entryType string) (bool, error) {
	project, err := GetProjectByID(uint(projectID))
	if err != nil {
		return false, err
	}

	// Ensure project supports given entry type
	entryTypeSupported := false
	for _, enabledEntryType := range project.EntryTypes {
		if entryType == enabledEntryType.Title {
			entryTypeSupported = true
		}
	}
	return entryTypeSupported, nil
}

func CreateProject(title string, ownerID uint) error {
	if len(title) == 0 {
		return errors.New("Project title must not be empty")
	}
	newProject := Project{Title: title, UserID: ownerID}
	err := DB.Create(&newProject).Error
	return err
}

func ProjectAlreadyExistsByTitleForUser(title string, ownerID uint) bool {
	doesNotExist := DB.Where("title = ? AND user_id = ?", title, ownerID).Find(&Project{}).RecordNotFound()
	return !doesNotExist
}
