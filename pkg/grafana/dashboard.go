package grafana

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

var (
	ErrGrafanaDashboardAlreadyExists = errors.New("grafana dashboard already exists")
	ErrGrafanaDashboardNotFound      = errors.New("grafana dashboard not found")
)

// DashboardId define a unique dashboard identifier.
type DashboardId uuid.UUID

// ParseDashboardId parses the given string and return a DashboardId if its valid.
// A valid DashboardId is a valid UUID v4.
func ParseDashboardId(dashboardID string) (DashboardId, error) {
	id, err := uuid.Parse(dashboardID)
	if err != nil {
		return DashboardId{}, err
	}

	return DashboardId(id), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (uid *DashboardId) UnmarshalJSON(rawJSON []byte) error {
	rawJSON = bytes.TrimPrefix(rawJSON, []byte(`"`))
	rawJSON = bytes.TrimSuffix(rawJSON, []byte(`"`))

	var err error
	*uid, err = ParseDashboardId(string(rawJSON))
	if err != nil {
		return err
	}

	return nil
}

// String implements fmt.Stringer.
func (uid DashboardId) String() string {
	return uuid.UUID(uid).String()
}

type DashboardMetadata struct {
	AnnotationsPermissions struct {
		dashboard struct {
			canAdd    bool
			canDelete bool
			canEdit   bool
		}
		organization struct {
			canAdd    bool
			canDelete bool
			canEdit   bool
		}
	}
	CanAdmin               bool
	CanDelete              bool
	CanEdit                bool
	CanSave                bool
	CanStar                bool
	Created                time.Time
	CreatedBy              string
	Expires                time.Time
	FolderId               int64
	FolderTitle            string
	FolderUid              FolderId
	FolderUrl              string
	HasAcl                 bool
	IsFolder               bool
	IsSnapshot             bool
	IsStarred              bool
	Provisioned            bool
	ProvisionedExternalId  string
	PublicDashboardEnabled bool
	PublicDashboardUid     string
	Slug                   string
	Type                   string
	Updated                time.Time
	UpdatedBy              string
	Url                    string
	Version                int
}

type Dashboard struct {
	Dashboard map[string]any    `json:"dashboard"`
	Metadata  DashboardMetadata `json:"meta"`
}

type SearchDashboardResult struct {
	Uid   DashboardId `json:"uid"`
	Title string      `json:"title"`
}

// CreateUpdateDashboard creates/updates a dashboard in the given organization and folder.
// dashboardJson map[string]any argument must contain a "uid" and "version" fields for updates.
// Version field must contains version BEFORE update, that is, the current version.
// If overwrite is sets to true, "version" field is optional.
//
// This method rely on user context and therefor, client mutex.
func (c Client) CreateUpdateDashboard(ctx context.Context, orgId OrgId, folder FolderId, dashboardJson map[string]any, overwrite bool) (DashboardId, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod("POST")
	req.SetRequestURI(fmt.Sprintf("%v/api/dashboards/db", c.cfg.Url))
	req.Header.Set(GrafanaOrgIdHeader, fmt.Sprint(orgId))
	c.addAuthorizationHeader(req)

	folderUid := folder.String()
	if folderUid == "00000000-0000-0000-0000-000000000000" {
		folderUid = ""
	}

	type requestBody struct {
		Dashboard map[string]any `json:"dashboard"`
		FolderUid string         `json:"folderUid"`
		Overwrite bool           `json:"overwrite"`
	}
	jsonBody, err := json.Marshal(requestBody{
		Dashboard: dashboardJson,
		FolderUid: folderUid,
		Overwrite: overwrite,
	})
	if err != nil {
		panic(err)
	}
	req.SetBody(jsonBody)
	req.Header.SetContentType("application/json")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err = c.do(ctx, req, resp)
	if err != nil {
		return DashboardId{}, fmt.Errorf("failed to query grafana to create dashboard: %w", err)
	}

	if resp.StatusCode() == 412 && strings.Contains(string(resp.Body()), "A dashboard with the same name in the folder already exists") {
		return DashboardId{}, ErrGrafanaDashboardAlreadyExists
	} else if resp.StatusCode() != 200 {
		return DashboardId{}, fmt.Errorf("failed to create/update grafana dashboard: %v %v", resp.StatusCode(), string(resp.Body()))
	}

	type responseBody struct {
		Uid DashboardId `json:"uid"`
	}
	respBody := responseBody{}
	err = json.Unmarshal(resp.Body(), &respBody)
	if err != nil {
		return DashboardId{}, fmt.Errorf("failed to parse grafana response: %w", err)
	}

	return respBody.Uid, nil
}

// GetDashboardByUid returns dashboard with the given id within the given organization.
// This method rely on user context and therefor, client mutex.
func (c Client) GetDashboardByUid(ctx context.Context, orgId OrgId, dashboardID DashboardId) (Dashboard, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod("GET")
	req.SetRequestURI(fmt.Sprintf("%v/api/dashboards/uid/%v", c.cfg.Url, dashboardID.String()))
	c.addAuthorizationHeader(req)
	req.Header.Set(GrafanaOrgIdHeader, fmt.Sprint(orgId))

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := c.do(ctx, req, resp)
	if err != nil {
		return Dashboard{}, fmt.Errorf("failed to query grafana to get dashboard by uid: %w", err)
	}

	if resp.StatusCode() == 404 {
		return Dashboard{}, ErrGrafanaDashboardNotFound
	} else if resp.StatusCode() != 200 {
		return Dashboard{}, fmt.Errorf("failed to get grafana dashboard by uid: %v %v", resp.StatusCode(), string(resp.Body()))
	}

	var respBody Dashboard
	err = json.Unmarshal(resp.Body(), &respBody)
	if err != nil {
		return Dashboard{}, fmt.Errorf("failed to parse grafana response: %w", err)
	}

	return respBody, nil
}

// DeleteDashboardByUid deletes a dashboard with the given ID within the given
// organization.
// This method rely on user context and therefor, client mutex.
func (c Client) DeleteDashboardByUid(ctx context.Context, orgId OrgId, dashboardID DashboardId) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod("DELETE")
	req.SetRequestURI(fmt.Sprintf("%v/api/dashboards/uid/%v", c.cfg.Url, dashboardID.String()))
	c.addAuthorizationHeader(req)
	req.Header.Set(GrafanaOrgIdHeader, fmt.Sprint(orgId))

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := c.do(ctx, req, resp)
	if err != nil {
		return fmt.Errorf("failed to query grafana to delete dashboard by uid: %w", err)
	}

	if resp.StatusCode() == 404 {
		return ErrGrafanaDashboardNotFound
	} else if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to delete grafana dashboard by uid: %v %v", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}

// SearchDashboards searches dashboard within the given organization.
// This method rely on user context and therefor, client mutex.
func (c Client) SearchDashboards(ctx context.Context, orgId OrgId, limit, page int) ([]SearchDashboardResult, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod("GET")
	req.SetRequestURI(fmt.Sprintf("%v/api/search?type=dash-db&limit=%v&page=%v", c.cfg.Url, limit, page))
	c.addAuthorizationHeader(req)
	req.Header.Set(GrafanaOrgIdHeader, fmt.Sprint(orgId))

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := c.do(ctx, req, resp)
	if err != nil {
		return nil, fmt.Errorf("failed to query grafana to delete dashboard by uid: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to search grafana dashboards: %v %v", resp.StatusCode(), string(resp.Body()))
	}

	var respBody []SearchDashboardResult
	err = json.Unmarshal(resp.Body(), &respBody)
	if err != nil {
		return nil, fmt.Errorf("failed to parse grafana response: %w", err)
	}

	return respBody, nil
}
