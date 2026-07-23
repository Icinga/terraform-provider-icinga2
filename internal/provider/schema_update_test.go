package provider

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestResourceSchemaUpdateSupport(t *testing.T) {
	ctx := context.Background()
	p := &icinga2Provider{}
	resources := p.Resources(ctx)

	for _, resourceFn := range resources {
		res := resourceFn()
		var resp resource.SchemaResponse
		res.Schema(ctx, resource.SchemaRequest{}, &resp)

		var metadataResp resource.MetadataResponse
		res.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "icinga2"}, &metadataResp)
		typeName := metadataResp.TypeName

		// Downtime does not support in-place updates (it deletes and recreates on changes).
		// Therefore, every user-configurable attribute in its schema must require replacement.
		supportsUpdate := typeName != "icinga2_downtime"

		if !supportsUpdate {
			for attrName, attr := range resp.Schema.Attributes {
				if attrName == "id" || attrName == "last_updated" || attrName == "names" {
					continue
				}

				if !isComputedAttr(attr) {
					if !hasRequiresReplaceModifier(attr) {
						t.Errorf("Resource %s does not support in-place updates, but attribute %s does not require replacement. This will cause runtime update errors!", typeName, attrName)
					}
				}
			}
		}
	}
}

func isComputedAttr(attr schema.Attribute) bool {
	switch a := attr.(type) {
	case schema.StringAttribute:
		return a.Computed
	case schema.BoolAttribute:
		return a.Computed
	case schema.Int64Attribute:
		return a.Computed
	case schema.ListAttribute:
		return a.Computed
	case schema.MapAttribute:
		return a.Computed
	}
	return false
}

func hasRequiresReplaceModifier(attr schema.Attribute) bool {
	switch a := attr.(type) {
	case schema.StringAttribute:
		for _, pm := range a.PlanModifiers {
			if strings.Contains(fmt.Sprintf("%T", pm), "requiresReplace") || strings.Contains(fmt.Sprintf("%T", pm), "RequiresReplace") {
				return true
			}
		}
	case schema.BoolAttribute:
		for _, pm := range a.PlanModifiers {
			if strings.Contains(fmt.Sprintf("%T", pm), "requiresReplace") || strings.Contains(fmt.Sprintf("%T", pm), "RequiresReplace") {
				return true
			}
		}
	case schema.Int64Attribute:
		for _, pm := range a.PlanModifiers {
			if strings.Contains(fmt.Sprintf("%T", pm), "requiresReplace") || strings.Contains(fmt.Sprintf("%T", pm), "RequiresReplace") {
				return true
			}
		}
	case schema.ListAttribute:
		for _, pm := range a.PlanModifiers {
			if strings.Contains(fmt.Sprintf("%T", pm), "requiresReplace") || strings.Contains(fmt.Sprintf("%T", pm), "RequiresReplace") {
				return true
			}
		}
	case schema.MapAttribute:
		for _, pm := range a.PlanModifiers {
			if strings.Contains(fmt.Sprintf("%T", pm), "requiresReplace") || strings.Contains(fmt.Sprintf("%T", pm), "RequiresReplace") {
				return true
			}
		}
	}
	return false
}
