package role

import "slices"

var permissionsByType = map[Type]map[Resource][]Action{
	TypeAdmin: {
		ResourceUser:        {ActionRead, ActionWrite, ActionCreate, ActionUpdate, ActionDelete},
		ResourceCourse:      {ActionRead, ActionWrite, ActionCreate, ActionUpdate, ActionDelete, ActionPublish},
		ResourceEnrollment:  {ActionRead, ActionWrite, ActionCreate, ActionUpdate, ActionDelete, ActionEnroll},
		ResourceQuiz:        {ActionRead, ActionWrite, ActionCreate, ActionUpdate, ActionDelete, ActionPublish},
		ResourceAttempt:     {ActionRead, ActionWrite, ActionUpdate, ActionDelete},
		ResourceSubmission:  {ActionRead, ActionWrite, ActionUpdate, ActionDelete, ActionGrade},
		ResourceGrade:       {ActionRead, ActionWrite, ActionUpdate, ActionDelete, ActionGrade},
		ResourceCertificate: {ActionRead, ActionWrite, ActionCreate, ActionUpdate, ActionDelete},
		ResourceFile:        {ActionRead, ActionWrite, ActionCreate, ActionUpdate, ActionDelete, ActionManage},
		ResourceAuditLog:    {ActionRead},
	},
	TypeCreator: {
		ResourceCourse:     {ActionRead, ActionWrite, ActionCreate, ActionUpdate, ActionPublish},
		ResourceEnrollment: {ActionRead, ActionEnroll},
		ResourceQuiz:       {ActionRead, ActionWrite, ActionCreate, ActionUpdate, ActionPublish},
		ResourceAttempt:    {ActionRead},
		ResourceSubmission: {ActionRead, ActionGrade},
		ResourceGrade:      {ActionRead, ActionWrite, ActionUpdate, ActionGrade},
		ResourceFile:       {ActionRead, ActionWrite, ActionCreate, ActionUpdate, ActionDelete, ActionManage},
	},
	TypeStudent: {
		ResourceCourse:     {ActionRead},
		ResourceEnrollment: {ActionRead},
		ResourceQuiz:       {ActionRead},
		ResourceAttempt:    {ActionRead, ActionCreate, ActionSubmit},
		ResourceSubmission: {ActionRead, ActionCreate, ActionSubmit},
		ResourceGrade:      {ActionRead},
		ResourceFile:       {ActionRead},
	},
	TypeOrganization: {
		ResourceUser:        {ActionRead},
		ResourceCourse:      {ActionRead},
		ResourceEnrollment:  {ActionRead},
		ResourceGrade:       {ActionRead},
		ResourceCertificate: {ActionRead},
	},
}

func (t Type) Allows(resource Resource, action Action) bool {
	if !t.IsValid() {
		return false
	}

	resources, ok := permissionsByType[t]
	if !ok {
		return false
	}

	actions, ok := resources[resource]
	if !ok {
		return false
	}

	return slices.Contains(actions, action)
}
