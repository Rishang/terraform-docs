/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package terraform

import (
	"fmt"
	"strings"

	terraformsdk "github.com/terraform-docs/plugin-sdk/terraform"
	"github.com/terraform-docs/terraform-docs/internal/types"
)

// Resource represents a managed or data type that is created by the module
type Resource struct {
	Type           string       `json:"type" toml:"type" xml:"type" yaml:"type"`
	Name           string       `json:"name" toml:"name" xml:"name" yaml:"name"`
	ProviderName   string       `json:"providerName" toml:"providerName" xml:"providerName" yaml:"providerName"`
	ProviderSource string       `json:"provicerSource" toml:"providerSource" xml:"providerSource" yaml:"providerSource"`
	Mode           string       `json:"mode" toml:"mode" xml:"mode" yaml:"mode"`
	Version        types.String `json:"version" toml:"version" xml:"version" yaml:"version"`
}

// Spec returns the resource spec addresses a specific resource in the config.
// It takes the form: resource_type.resource_name[resource index]
// For more details, see:
// https://www.terraform.io/docs/cli/state/resource-addressing.html#resource-spec
func (r *Resource) Spec() string {
	return r.ProviderName + "_" + r.Type + "." + r.Name
}

// GetMode returns normalized resource type as "resource" or "data source"
func (r *Resource) GetMode() string {
	switch r.Mode {
	case "managed":
		return "resource"
	case "data":
		return "data source"
	default:
		return "invalid"
	}
}

// URL returns a best guess at the URL for resource documentation
func (r *Resource) URL() string {
	kind := ""
	switch r.Mode {
	case "managed":
		kind = "resources"
	case "data":
		kind = "data-sources"
	default:
		return ""
	}

	if strings.Count(r.ProviderSource, "/") > 1 {
		return ""
	}
	return fmt.Sprintf("https://registry.terraform.io/providers/%s/%s/docs/%s/%s", r.ProviderSource, r.Version, kind, r.Type)
}

type resources []*Resource

func (rr resources) convert() []*terraformsdk.Resource {
	list := []*terraformsdk.Resource{}
	for _, r := range rr {
		list = append(list, &terraformsdk.Resource{
			Type:           r.Type,
			Name:           r.Name,
			ProviderName:   r.ProviderName,
			ProviderSource: r.ProviderSource,
			Version:        fmt.Sprintf("%v", r.Version.Raw()),
		})
	}

	return list
}

type resourcesSortedByType []*Resource

func (a resourcesSortedByType) Len() int      { return len(a) }
func (a resourcesSortedByType) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a resourcesSortedByType) Less(i, j int) bool {
	if a[i].Mode == a[j].Mode {
		if a[i].Spec() == a[j].Spec() {
			return a[i].Name <= a[j].Name
		}
		return a[i].Spec() < a[j].Spec()
	}
	return a[i].Mode > a[j].Mode
}
