package hcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

// Image represents an Image in the Hetzner Cloud.
type Image struct {
	ID          int64
	Name        string
	Type        ImageType
	Status      ImageStatus
	Description string
	ImageSize   float32
	DiskSize    float32
	Created     time.Time
	CreatedFrom *Server
	BoundTo     *Server
	RapidDeploy bool

	OSFlavor     string
	OSVersion    string
	Architecture Architecture

	Protection ImageProtection
	Deprecated time.Time // The zero value denotes the image is not deprecated.
	Labels     map[string]string
	Deleted    time.Time
}

// IsDeprecated returns whether the image is deprecated.
func (image *Image) IsDeprecated() bool {
	return !image.Deprecated.IsZero()
}

// IsDeleted returns whether the image is deleted.
func (image *Image) IsDeleted() bool {
	return !image.Deleted.IsZero()
}

// ImageProtection represents the protection level of an image.
type ImageProtection struct {
	Delete bool
}

// ImageType specifies the type of an image.
type ImageType string

const (
	// ImageTypeSnapshot represents a snapshot image.
	ImageTypeSnapshot ImageType = "snapshot"
	// ImageTypeBackup represents a backup image.
	ImageTypeBackup ImageType = "backup"
	// ImageTypeSystem represents a system image.
	ImageTypeSystem ImageType = "system"
	// ImageTypeApp represents a one click app image.
	ImageTypeApp ImageType = "app"
)

// ImageStatus specifies the status of an image.
type ImageStatus string

const (
	// ImageStatusCreating is the status when an image is being created.
	ImageStatusCreating ImageStatus = "creating"
	// ImageStatusAvailable is the status when an image is available.
	ImageStatusAvailable ImageStatus = "available"
)

// ImageClient is a client for the image API.
type ImageClient struct {
	client *Client
	Action *ResourceActionClient
}

// GetByID retrieves an image by its ID. If the image does not exist, nil is returned.
func (c *ImageClient) GetByID(ctx context.Context, id int64) (*Image, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("/images/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.ImageGetResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		if IsError(err, ErrorCodeNotFound) {
			return nil, resp, nil
		}
		return nil, nil, err
	}
	return ImageFromSchema(body.Image), resp, nil
}

// GetByName retrieves an image by its name. If the image does not exist, nil is returned.
//
// Deprecated: Use [ImageClient.GetByNameAndArchitecture] instead.
func (c *ImageClient) GetByName(ctx context.Context, name string) (*Image, *Response, error) {
	if name == "" {
		return nil, nil, nil
	}
	images, response, err := c.List(ctx, ImageListOpts{Name: name})
	if len(images) == 0 {
		return nil, response, err
	}
	return images[0], response, err
}

// GetByNameAndArchitecture retrieves an image by its name and architecture. If the image does not exist,
// nil is returned.
// In contrast to [ImageClient.Get], this method also returns deprecated images. Depending on your needs you should
// check for this in your calling method.
func (c *ImageClient) GetByNameAndArchitecture(ctx context.Context, name string, architecture Architecture) (*Image, *Response, error) {
	if name == "" {
		return nil, nil, nil
	}
	images, response, err := c.List(ctx, ImageListOpts{Name: name, Architecture: []Architecture{architecture}, IncludeDeprecated: true})
	if len(images) == 0 {
		return nil, response, err
	}
	return images[0], response, err
}

// Get retrieves an image by its ID if the input can be parsed as an integer, otherwise it
// retrieves an image by its name. If the image does not exist, nil is returned.
//
// Deprecated: Use [ImageClient.GetForArchitecture] instead.
func (c *ImageClient) Get(ctx context.Context, idOrName string) (*Image, *Response, error) {
	if id, err := strconv.ParseInt(idOrName, 10, 64); err == nil {
		img, res, err := c.GetByID(ctx, id)
		if img != nil {
			return img, res, err
		}
	}
	return c.GetByName(ctx, idOrName)
}

// GetForArchitecture retrieves an image by its ID if the input can be parsed as an integer, otherwise it
// retrieves an image by its name and architecture. If the image does not exist, nil is returned.
//
// In contrast to [ImageClient.Get], this method also returns deprecated images. Depending on your needs you should
// check for this in your calling method.
func (c *ImageClient) GetForArchitecture(ctx context.Context, idOrName string, architecture Architecture) (*Image, *Response, error) {
	if id, err := strconv.ParseInt(idOrName, 10, 64); err == nil {
		img, res, err := c.GetByID(ctx, id)
		if img != nil || err != nil {
			return img, res, err
		}
	}
	return c.GetByNameAndArchitecture(ctx, idOrName, architecture)
}

// ImageListOpts specifies options for listing images.
type ImageListOpts struct {
	ListOpts
	Type              []ImageType
	BoundTo           *Server
	Name              string
	Sort              []string
	Status            []ImageStatus
	IncludeDeprecated bool
	Architecture      []Architecture
}

func (l ImageListOpts) values() url.Values {
	vals := l.ListOpts.Values()
	for _, typ := range l.Type {
		vals.Add("type", string(typ))
	}
	if l.BoundTo != nil {
		vals.Add("bound_to", strconv.FormatInt(l.BoundTo.ID, 10))
	}
	if l.Name != "" {
		vals.Add("name", l.Name)
	}
	if l.IncludeDeprecated {
		vals.Add("include_deprecated", strconv.FormatBool(l.IncludeDeprecated))
	}
	for _, sort := range l.Sort {
		vals.Add("sort", sort)
	}
	for _, status := range l.Status {
		vals.Add("status", string(status))
	}
	for _, arch := range l.Architecture {
		vals.Add("architecture", string(arch))
	}
	return vals
}

// List returns a list of images for a specific page.
//
// Please note that filters specified in opts are not taken into account
// when their value corresponds to their zero value or when they are empty.
func (c *ImageClient) List(ctx context.Context, opts ImageListOpts) ([]*Image, *Response, error) {
	path := "/images?" + opts.values().Encode()
	req, err := c.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.ImageListResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	images := make([]*Image, 0, len(body.Images))
	for _, i := range body.Images {
		images = append(images, ImageFromSchema(i))
	}
	return images, resp, nil
}

// All returns all images.
func (c *ImageClient) All(ctx context.Context) ([]*Image, error) {
	return c.AllWithOpts(ctx, ImageListOpts{ListOpts: ListOpts{PerPage: 50}})
}

// AllWithOpts returns all images for the given options.
func (c *ImageClient) AllWithOpts(ctx context.Context, opts ImageListOpts) ([]*Image, error) {
	allImages := []*Image{}

	err := c.client.all(func(page int) (*Response, error) {
		opts.Page = page
		images, resp, err := c.List(ctx, opts)
		if err != nil {
			return resp, err
		}
		allImages = append(allImages, images...)
		return resp, nil
	})
	if err != nil {
		return nil, err
	}

	return allImages, nil
}

// Delete deletes an image.
func (c *ImageClient) Delete(ctx context.Context, image *Image) (*Response, error) {
	req, err := c.client.NewRequest(ctx, "DELETE", fmt.Sprintf("/images/%d", image.ID), nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req, nil)
}

// ImageUpdateOpts specifies options for updating an image.
type ImageUpdateOpts struct {
	Description *string
	Type        ImageType
	Labels      map[string]string
}

// Update updates an image.
func (c *ImageClient) Update(ctx context.Context, image *Image, opts ImageUpdateOpts) (*Image, *Response, error) {
	reqBody := schema.ImageUpdateRequest{
		Description: opts.Description,
	}
	if opts.Type != "" {
		reqBody.Type = Ptr(string(opts.Type))
	}
	if opts.Labels != nil {
		reqBody.Labels = &opts.Labels
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/images/%d", image.ID)
	req, err := c.client.NewRequest(ctx, "PUT", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ImageUpdateResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ImageFromSchema(respBody.Image), resp, nil
}

// ImageChangeProtectionOpts specifies options for changing the resource protection level of an image.
type ImageChangeProtectionOpts struct {
	Delete *bool
}

// ChangeProtection changes the resource protection level of an image.
func (c *ImageClient) ChangeProtection(ctx context.Context, image *Image, opts ImageChangeProtectionOpts) (*Action, *Response, error) {
	reqBody := schema.ImageActionChangeProtectionRequest{
		Delete: opts.Delete,
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/images/%d/actions/change_protection", image.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ImageActionChangeProtectionResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, err
}
