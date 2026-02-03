package auth

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"

	"github.com/dlsu-lscs/lscs-core-api/internal/database"
	"github.com/dlsu-lscs/lscs-core-api/internal/repository"
)

// Position hierarchy levels (higher number = more authority)
var positionHierarchy = map[string]int{
	"PRES": 7,
	"EVP":  6,
	"VP":   5,
	"AVP":  4,
	"CT":   3,
	"JO":   2,
	"MEM":  1,
}

// RBACService handles role-based access control
type RBACService struct {
	dbService database.Service
}

// NewRBACService creates a new RBAC service
func NewRBACService(dbService database.Service) *RBACService {
	return &RBACService{dbService: dbService}
}

// GetPositionLevel returns the hierarchy level for a position
func GetPositionLevel(positionID string) int {
	if level, ok := positionHierarchy[positionID]; ok {
		return level
	}
	return 0 // unknown position has no authority
}

// IsHigherPosition returns true if position1 has higher authority than position2
func IsHigherPosition(position1, position2 string) bool {
	return GetPositionLevel(position1) > GetPositionLevel(position2)
}

// IsHigherOrEqualPosition returns true if position1 has equal or higher authority than position2
func IsHigherOrEqualPosition(position1, position2 string) bool {
	return GetPositionLevel(position1) >= GetPositionLevel(position2)
}

// IsAdmin checks if a member has the ADMIN role
func (s *RBACService) IsAdmin(ctx context.Context, memberID int32) bool {
	q := repository.New(s.dbService.GetConnection())
	isAdmin, err := q.IsAdmin(ctx, memberID)
	if err != nil {
		log.Error().Err(err).Int32("member_id", memberID).Msg("failed to check admin status")
		return false
	}
	return isAdmin
}

// HasRole checks if a member has a specific role
func (s *RBACService) HasRole(ctx context.Context, memberID int32, roleID string) bool {
	q := repository.New(s.dbService.GetConnection())
	hasRole, err := q.HasRole(ctx, repository.HasRoleParams{
		MemberID: memberID,
		RoleID:   roleID,
	})
	if err != nil {
		log.Error().Err(err).Int32("member_id", memberID).Str("role", roleID).Msg("failed to check role")
		return false
	}
	return hasRole
}

// GetMemberRoles returns all roles assigned to a member
func (s *RBACService) GetMemberRoles(ctx context.Context, memberID int32) ([]repository.GetMemberRolesRow, error) {
	q := repository.New(s.dbService.GetConnection())
	return q.GetMemberRoles(ctx, memberID)
}

// GrantRole assigns a role to a member
func (s *RBACService) GrantRole(ctx context.Context, memberID int32, roleID string, grantedBy int32) error {
	q := repository.New(s.dbService.GetConnection())
	return q.GrantRole(ctx, repository.GrantRoleParams{
		MemberID:  memberID,
		RoleID:    roleID,
		GrantedBy: sql.NullInt32{Int32: grantedBy, Valid: true},
	})
}

// RevokeRole removes a role from a member
func (s *RBACService) RevokeRole(ctx context.Context, memberID int32, roleID string) error {
	q := repository.New(s.dbService.GetConnection())
	return q.RevokeRole(ctx, repository.RevokeRoleParams{
		MemberID: memberID,
		RoleID:   roleID,
	})
}

// CanEditMember checks if an actor can edit a target member based on:
// 1. Admin role (can edit anyone)
// 2. Same member (can edit own profile)
// 3. Position hierarchy within same committee (VP can edit lower positions in their committee)
// 4. Higher position hierarchy across committees (EVP/PRES can edit lower positions anywhere)
func (s *RBACService) CanEditMember(ctx context.Context, actorID, targetID int32) bool {
	// same member can always edit their own profile
	if actorID == targetID {
		return true
	}

	// admin can edit anyone
	if s.IsAdmin(ctx, actorID) {
		return true
	}

	q := repository.New(s.dbService.GetConnection())

	// get actor info
	actor, err := q.GetMemberInfoById(ctx, actorID)
	if err != nil {
		log.Error().Err(err).Int32("actor_id", actorID).Msg("failed to get actor info")
		return false
	}

	// get target info
	target, err := q.GetMemberInfoById(ctx, targetID)
	if err != nil {
		log.Error().Err(err).Int32("target_id", targetID).Msg("failed to get target info")
		return false
	}

	actorPosition := actor.PositionID.String
	targetPosition := target.PositionID.String
	actorCommittee := actor.CommitteeID.String
	targetCommittee := target.CommitteeID.String

	// check position hierarchy
	if !IsHigherPosition(actorPosition, targetPosition) {
		return false
	}

	// EVP and PRES can edit anyone with lower position (any committee)
	if actorPosition == "PRES" || actorPosition == "EVP" {
		return true
	}

	// VP can edit lower positions only in their own committee
	if actorPosition == "VP" && actorCommittee == targetCommittee {
		return true
	}

	return false
}

// CanViewMember checks if an actor can view a target member's full info.
// All authenticated members can view basic info, but some fields may require higher permissions.
func (s *RBACService) CanViewMember(ctx context.Context, actorID, targetID int32) bool {
	// for now, all authenticated users can view any member
	// this can be made more restrictive if needed
	return true
}

// CanManageRoles checks if an actor can grant/revoke roles
// Only ADMIN role holders can manage roles
func (s *RBACService) CanManageRoles(ctx context.Context, actorID int32) bool {
	return s.IsAdmin(ctx, actorID)
}

// CanAccessAPIKeyManagement checks if a member can access API key management
// Requirements: Member of RND committee OR AVP+ position
func (s *RBACService) CanAccessAPIKeyManagement(ctx context.Context, memberID int32) bool {
	if s.IsAdmin(ctx, memberID) {
		return true
	}

	q := repository.New(s.dbService.GetConnection())
	member, err := q.GetMemberInfoById(ctx, memberID)
	if err != nil {
		log.Error().Err(err).Int32("member_id", memberID).Msg("failed to get member info for API key access check")
		return false
	}

	// RND members can access
	if member.CommitteeID.String == "RND" {
		return true
	}

	// AVP and above can access
	return GetPositionLevel(member.PositionID.String) >= GetPositionLevel("AVP")
}

// EditableField represents which fields can be edited by whom
type EditableField string

const (
	FieldNickname      EditableField = "nickname"
	FieldTelegram      EditableField = "telegram"
	FieldDiscord       EditableField = "discord"
	FieldInterests     EditableField = "interests"
	FieldContactNumber EditableField = "contact_number"
	FieldFbLink        EditableField = "fb_link"
	FieldFullName      EditableField = "full_name"
	FieldEmail         EditableField = "email"
	FieldPositionID    EditableField = "position_id"
	FieldCommitteeID   EditableField = "committee_id"
	FieldCollege       EditableField = "college"
	FieldProgram       EditableField = "program"
	FieldHouseID       EditableField = "house_id"
)

// selfEditableFields are fields that members can edit on their own profile
var selfEditableFields = map[EditableField]bool{
	FieldNickname:      true,
	FieldTelegram:      true,
	FieldDiscord:       true,
	FieldInterests:     true,
	FieldContactNumber: true,
	FieldFbLink:        true,
}

// authorizedEditableFields are fields that authorized users can edit on other profiles
// (in addition to self-editable fields)
var authorizedEditableFields = map[EditableField]bool{
	FieldFullName:    true,
	FieldEmail:       true,
	FieldPositionID:  true,
	FieldCommitteeID: true,
	FieldCollege:     true,
	FieldProgram:     true,
	FieldHouseID:     true,
}

// CanEditField checks if an actor can edit a specific field on a target member
func (s *RBACService) CanEditField(ctx context.Context, actorID, targetID int32, field EditableField) bool {
	// admin can edit any field
	if s.IsAdmin(ctx, actorID) {
		return true
	}

	// self-editing: only self-editable fields
	if actorID == targetID {
		return selfEditableFields[field]
	}

	// authorized user editing: check if they can edit the member first
	if !s.CanEditMember(ctx, actorID, targetID) {
		return false
	}

	// authorized users can edit both self-editable and authorized-editable fields
	return selfEditableFields[field] || authorizedEditableFields[field]
}

// GetEditableFields returns the list of fields an actor can edit on a target member
func (s *RBACService) GetEditableFields(ctx context.Context, actorID, targetID int32) []EditableField {
	var fields []EditableField

	// admin can edit everything
	if s.IsAdmin(ctx, actorID) {
		for field := range selfEditableFields {
			fields = append(fields, field)
		}
		for field := range authorizedEditableFields {
			fields = append(fields, field)
		}
		return fields
	}

	// self-editing
	if actorID == targetID {
		for field := range selfEditableFields {
			fields = append(fields, field)
		}
		return fields
	}

	// authorized user editing
	if s.CanEditMember(ctx, actorID, targetID) {
		for field := range selfEditableFields {
			fields = append(fields, field)
		}
		for field := range authorizedEditableFields {
			fields = append(fields, field)
		}
	}

	return fields
}
