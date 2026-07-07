package service

import (
	"testing"

	"go-backend/internal/models"

	"gorm.io/gorm"
)

func TestCanAccessUserForTeacherWithAssignedGroups(t *testing.T) {
	svc := NewStudentService(nil)
	currentUser := models.User{Role: "teacher", Group: "VM,SM"}
	targetUser := models.User{Role: "student", Group: "VM"}

	if !svc.CanAccessUser(currentUser, targetUser) {
		t.Fatalf("teacher should access student from assigned group")
	}
}

func TestCanAccessUserForTeacherWithoutAssignedGroup(t *testing.T) {
	svc := NewStudentService(nil)
	currentUser := models.User{Role: "teacher", Group: "VM"}
	targetUser := models.User{Role: "student", Group: "KT"}

	if svc.CanAccessUser(currentUser, targetUser) {
		t.Fatalf("teacher should not access student from non-assigned group")
	}
}

func TestCanAccessUserForStudentOwnProfile(t *testing.T) {
	svc := NewStudentService(nil)
	currentUser := models.User{Model: gorm.Model{ID: 7}, Role: "student", Group: "VM"}
	targetUser := models.User{Model: gorm.Model{ID: 7}, Role: "student", Group: "VM"}

	if !svc.CanAccessUser(currentUser, targetUser) {
		t.Fatalf("student should access own profile")
	}
}
