package permission

type Permission struct {
	resource Resource
	actions  []Action
}

func New(resource Resource, actions []Action) (*Permission, error) {
	if !resource.IsValid() {
		return nil, ErrInvalidResource
	}

	if err := reviewActions(actions); err != nil {
		return nil, err
	}

	copyActions := make([]Action, len(actions))
	copy(copyActions, actions)

	return &Permission{
		resource: resource,
		actions:  copyActions,
	}, nil
}

func (p *Permission) Resource() Resource {
	return p.resource
}

func (p *Permission) Actions() []Action {
	copyActions := make([]Action, len(p.actions))
	copy(copyActions, p.actions)

	return copyActions
}
