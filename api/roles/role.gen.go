// Code generated by "make api"; DO NOT EDIT.
package roles

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kr/pretty"

	"github.com/hashicorp/watchtower/api"
	"github.com/hashicorp/watchtower/api/scopes"
)

type Role struct {
	Id           string            `json:"id,omitempty"`
	Scope        *scopes.ScopeInfo `json:"scope,omitempty"`
	Name         string            `json:"name,omitempty"`
	Description  string            `json:"description,omitempty"`
	CreatedTime  time.Time         `json:"created_time,omitempty"`
	UpdatedTime  time.Time         `json:"updated_time,omitempty"`
	Disabled     bool              `json:"disabled,omitempty"`
	GrantScopeId string            `json:"grant_scope_id,omitempty"`
	Version      uint32            `json:"version,omitempty"`
	PrincipalIds []string          `json:"principal_ids,omitempty"`
	Principals   []*Principal      `json:"principals,omitempty"`
	GrantStrings []string          `json:"grant_strings,omitempty"`
	Grants       []*Grant          `json:"grants,omitempty"`
}

type roleClient struct {
	client *api.Client
}

func NewRoleClient(c *api.Client) *roleClient {
	return &roleClient{client: c}
}

func (c *roleClient) Create(ctx context.Context, opt ...Option) (*Role, *api.Error, error) {
	if c.client == nil {
		return nil, nil, fmt.Errorf("nil client")
	}

	opts, apiOpts := getOpts(opt...)

	req, err := c.client.NewRequest(ctx, "POST", fmt.Sprintf("roles"), opts.valueMap, apiOpts...)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating Create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error performing client request during Create call: %w", err)
	}

	target := new(Role)
	apiErr, err := resp.Decode(target)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding Create response: %w", err)
	}

	return target, apiErr, nil
}

func (c *roleClient) Read(ctx context.Context, roleId string, opt ...Option) (*Role, *api.Error, error) {
	if roleId == "" {
		return nil, nil, fmt.Errorf("empty roleId value passed into Read request")
	}

	if c.client == nil {
		return nil, nil, fmt.Errorf("nil client")
	}

	_, apiOpts := getOpts(opt...)

	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("roles/%s", roleId), nil, apiOpts...)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating Read request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error performing client request during Read call: %w", err)
	}

	target := new(Role)
	apiErr, err := resp.Decode(target)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding Read response: %w", err)
	}

	return target, apiErr, nil
}

func (c *roleClient) Update(ctx context.Context, roleId string, version uint32, opt ...Option) (*Role, *api.Error, error) {
	if roleId == "" {
		return nil, nil, fmt.Errorf("empty roleId value passed into Update request")
	}
	if c.client == nil {
		return nil, nil, fmt.Errorf("nil client")
	}

	opts, apiOpts := getOpts(opt...)

	if version == 0 {
		if !opts.withAutomaticVersioning {
			return nil, nil, errors.New("zero version number passed into Update request and automatic versioning not specified")
		}
		existingTarget, existingApiErr, existingErr := c.Read(ctx, roleId, opt...)
		if existingErr != nil {
			return nil, nil, fmt.Errorf("error performing initial check-and-set read: %w", existingErr)
		}
		if existingApiErr != nil {
			return nil, nil, fmt.Errorf("error from controller when performing initial check-and-set read: %s", pretty.Sprint(existingApiErr))
		}
		if existingTarget == nil {
			return nil, nil, errors.New("nil resource found when performing initial check-and-set read")
		}
		version = existingTarget.Version
	}
	opts.valueMap["version"] = version

	req, err := c.client.NewRequest(ctx, "PATCH", fmt.Sprintf("roles/%s", roleId), opts.valueMap, apiOpts...)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating Update request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error performing client request during Update call: %w", err)
	}

	target := new(Role)
	apiErr, err := resp.Decode(target)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding Update response: %w", err)
	}

	return target, apiErr, nil
}

func (c *roleClient) Delete(ctx context.Context, roleId string, opt ...Option) (bool, *api.Error, error) {
	if roleId == "" {
		return false, nil, fmt.Errorf("empty roleId value passed into Delete request")
	}

	if c.client == nil {
		return false, nil, fmt.Errorf("nil client")
	}

	_, apiOpts := getOpts(opt...)

	req, err := c.client.NewRequest(ctx, "DELETE", fmt.Sprintf("roles/%s", roleId), nil, apiOpts...)
	if err != nil {
		return false, nil, fmt.Errorf("error creating Delete request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return false, nil, fmt.Errorf("error performing client request during Delete call: %w", err)
	}

	type deleteResponse struct {
		Existed bool
	}
	target := &deleteResponse{}
	apiErr, err := resp.Decode(target)
	if err != nil {
		return false, nil, fmt.Errorf("error decoding Delete response: %w", err)
	}

	return target.Existed, apiErr, nil
}

func (c *roleClient) List(ctx context.Context, opt ...Option) ([]*Role, *api.Error, error) {
	if c.client == nil {
		return nil, nil, fmt.Errorf("nil client")
	}

	_, apiOpts := getOpts(opt...)

	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("roles"), nil, apiOpts...)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating List request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error performing client request during List call: %w", err)
	}

	type listResponse struct {
		Items []*Role
	}
	target := &listResponse{}
	apiErr, err := resp.Decode(target)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding List response: %w", err)
	}

	return target.Items, apiErr, nil
}

func (c *roleClient) AddGrants(ctx context.Context, roleId string, version uint32, grantStrings []string, opt ...Option) (*Role, *api.Error, error) {
	if roleId == "" {
		return nil, nil, fmt.Errorf("empty roleId value passed into AddGrants request")
	}
	if c.client == nil {
		return nil, nil, fmt.Errorf("nil client")
	}

	opts, apiOpts := getOpts(opt...)

	if version == 0 {
		if !opts.withAutomaticVersioning {
			return nil, nil, errors.New("zero version number passed into AddGrants request")
		}
		existingTarget, existingApiErr, existingErr := c.Read(ctx, roleId, opt...)
		if existingErr != nil {
			return nil, nil, fmt.Errorf("error performing initial check-and-set read: %w", existingErr)
		}
		if existingApiErr != nil {
			return nil, nil, fmt.Errorf("error from controller when performing initial check-and-set read: %s", pretty.Sprint(existingApiErr))
		}
		if existingTarget == nil {
			return nil, nil, errors.New("nil resource found when performing initial check-and-set read")
		}
		version = existingTarget.Version
	}
	opts.valueMap["version"] = version

	if len(grantStrings) > 0 {
		opts.valueMap["grant_strings"] = grantStrings
	}

	req, err := c.client.NewRequest(ctx, "POST", fmt.Sprintf("roles/%s:add-grants", roleId), opts.valueMap, apiOpts...)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating AddGrants request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error performing client request during AddGrants call: %w", err)
	}

	target := new(Role)
	apiErr, err := resp.Decode(target)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding AddGrants response: %w", err)
	}

	return target, apiErr, nil
}

func (c *roleClient) AddPrincipals(ctx context.Context, roleId string, version uint32, principalIds []string, opt ...Option) (*Role, *api.Error, error) {
	if roleId == "" {
		return nil, nil, fmt.Errorf("empty roleId value passed into AddPrincipals request")
	}
	if c.client == nil {
		return nil, nil, fmt.Errorf("nil client")
	}

	opts, apiOpts := getOpts(opt...)

	if version == 0 {
		if !opts.withAutomaticVersioning {
			return nil, nil, errors.New("zero version number passed into AddPrincipals request")
		}
		existingTarget, existingApiErr, existingErr := c.Read(ctx, roleId, opt...)
		if existingErr != nil {
			return nil, nil, fmt.Errorf("error performing initial check-and-set read: %w", existingErr)
		}
		if existingApiErr != nil {
			return nil, nil, fmt.Errorf("error from controller when performing initial check-and-set read: %s", pretty.Sprint(existingApiErr))
		}
		if existingTarget == nil {
			return nil, nil, errors.New("nil resource found when performing initial check-and-set read")
		}
		version = existingTarget.Version
	}
	opts.valueMap["version"] = version

	if len(principalIds) > 0 {
		opts.valueMap["principal_ids"] = principalIds
	}

	req, err := c.client.NewRequest(ctx, "POST", fmt.Sprintf("roles/%s:add-principals", roleId), opts.valueMap, apiOpts...)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating AddPrincipals request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error performing client request during AddPrincipals call: %w", err)
	}

	target := new(Role)
	apiErr, err := resp.Decode(target)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding AddPrincipals response: %w", err)
	}

	return target, apiErr, nil
}

func (c *roleClient) SetGrants(ctx context.Context, roleId string, version uint32, grantStrings []string, opt ...Option) (*Role, *api.Error, error) {
	if roleId == "" {
		return nil, nil, fmt.Errorf("empty roleId value passed into SetGrants request")
	}
	if c.client == nil {
		return nil, nil, fmt.Errorf("nil client")
	}

	opts, apiOpts := getOpts(opt...)

	if version == 0 {
		if !opts.withAutomaticVersioning {
			return nil, nil, errors.New("zero version number passed into SetGrants request")
		}
		existingTarget, existingApiErr, existingErr := c.Read(ctx, roleId, opt...)
		if existingErr != nil {
			return nil, nil, fmt.Errorf("error performing initial check-and-set read: %w", existingErr)
		}
		if existingApiErr != nil {
			return nil, nil, fmt.Errorf("error from controller when performing initial check-and-set read: %s", pretty.Sprint(existingApiErr))
		}
		if existingTarget == nil {
			return nil, nil, errors.New("nil resource found when performing initial check-and-set read")
		}
		version = existingTarget.Version
	}
	opts.valueMap["version"] = version

	if len(grantStrings) > 0 {
		opts.valueMap["grant_strings"] = grantStrings
	} else if grantStrings != nil {
		// In this function, a non-nil but empty list means clear out
		opts.valueMap["grant_strings"] = nil
	}

	req, err := c.client.NewRequest(ctx, "POST", fmt.Sprintf("roles/%s:set-grants", roleId), opts.valueMap, apiOpts...)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating SetGrants request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error performing client request during SetGrants call: %w", err)
	}

	target := new(Role)
	apiErr, err := resp.Decode(target)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding SetGrants response: %w", err)
	}

	return target, apiErr, nil
}

func (c *roleClient) SetPrincipals(ctx context.Context, roleId string, version uint32, principalIds []string, opt ...Option) (*Role, *api.Error, error) {
	if roleId == "" {
		return nil, nil, fmt.Errorf("empty roleId value passed into SetPrincipals request")
	}
	if c.client == nil {
		return nil, nil, fmt.Errorf("nil client")
	}

	opts, apiOpts := getOpts(opt...)

	if version == 0 {
		if !opts.withAutomaticVersioning {
			return nil, nil, errors.New("zero version number passed into SetPrincipals request")
		}
		existingTarget, existingApiErr, existingErr := c.Read(ctx, roleId, opt...)
		if existingErr != nil {
			return nil, nil, fmt.Errorf("error performing initial check-and-set read: %w", existingErr)
		}
		if existingApiErr != nil {
			return nil, nil, fmt.Errorf("error from controller when performing initial check-and-set read: %s", pretty.Sprint(existingApiErr))
		}
		if existingTarget == nil {
			return nil, nil, errors.New("nil resource found when performing initial check-and-set read")
		}
		version = existingTarget.Version
	}
	opts.valueMap["version"] = version

	if len(principalIds) > 0 {
		opts.valueMap["principal_ids"] = principalIds
	} else if principalIds != nil {
		// In this function, a non-nil but empty list means clear out
		opts.valueMap["principal_ids"] = nil
	}

	req, err := c.client.NewRequest(ctx, "POST", fmt.Sprintf("roles/%s:set-principals", roleId), opts.valueMap, apiOpts...)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating SetPrincipals request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error performing client request during SetPrincipals call: %w", err)
	}

	target := new(Role)
	apiErr, err := resp.Decode(target)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding SetPrincipals response: %w", err)
	}

	return target, apiErr, nil
}

func (c *roleClient) RemoveGrants(ctx context.Context, roleId string, version uint32, grantStrings []string, opt ...Option) (*Role, *api.Error, error) {
	if roleId == "" {
		return nil, nil, fmt.Errorf("empty roleId value passed into RemoveGrants request")
	}
	if c.client == nil {
		return nil, nil, fmt.Errorf("nil client")
	}

	opts, apiOpts := getOpts(opt...)

	if version == 0 {
		if !opts.withAutomaticVersioning {
			return nil, nil, errors.New("zero version number passed into RemoveGrants request")
		}
		existingTarget, existingApiErr, existingErr := c.Read(ctx, roleId, opt...)
		if existingErr != nil {
			return nil, nil, fmt.Errorf("error performing initial check-and-set read: %w", existingErr)
		}
		if existingApiErr != nil {
			return nil, nil, fmt.Errorf("error from controller when performing initial check-and-set read: %s", pretty.Sprint(existingApiErr))
		}
		if existingTarget == nil {
			return nil, nil, errors.New("nil resource found when performing initial check-and-set read")
		}
		version = existingTarget.Version
	}
	opts.valueMap["version"] = version

	if len(grantStrings) > 0 {
		opts.valueMap["grant_strings"] = grantStrings
	}

	req, err := c.client.NewRequest(ctx, "POST", fmt.Sprintf("roles/%s:remove-grants", roleId), opts.valueMap, apiOpts...)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating RemoveGrants request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error performing client request during RemoveGrants call: %w", err)
	}

	target := new(Role)
	apiErr, err := resp.Decode(target)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding RemoveGrants response: %w", err)
	}

	return target, apiErr, nil
}

func (c *roleClient) RemovePrincipals(ctx context.Context, roleId string, version uint32, principalIds []string, opt ...Option) (*Role, *api.Error, error) {
	if roleId == "" {
		return nil, nil, fmt.Errorf("empty roleId value passed into RemovePrincipals request")
	}
	if c.client == nil {
		return nil, nil, fmt.Errorf("nil client")
	}

	opts, apiOpts := getOpts(opt...)

	if version == 0 {
		if !opts.withAutomaticVersioning {
			return nil, nil, errors.New("zero version number passed into RemovePrincipals request")
		}
		existingTarget, existingApiErr, existingErr := c.Read(ctx, roleId, opt...)
		if existingErr != nil {
			return nil, nil, fmt.Errorf("error performing initial check-and-set read: %w", existingErr)
		}
		if existingApiErr != nil {
			return nil, nil, fmt.Errorf("error from controller when performing initial check-and-set read: %s", pretty.Sprint(existingApiErr))
		}
		if existingTarget == nil {
			return nil, nil, errors.New("nil resource found when performing initial check-and-set read")
		}
		version = existingTarget.Version
	}
	opts.valueMap["version"] = version

	if len(principalIds) > 0 {
		opts.valueMap["principal_ids"] = principalIds
	}

	req, err := c.client.NewRequest(ctx, "POST", fmt.Sprintf("roles/%s:remove-principals", roleId), opts.valueMap, apiOpts...)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating RemovePrincipals request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error performing client request during RemovePrincipals call: %w", err)
	}

	target := new(Role)
	apiErr, err := resp.Decode(target)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding RemovePrincipals response: %w", err)
	}

	return target, apiErr, nil
}
