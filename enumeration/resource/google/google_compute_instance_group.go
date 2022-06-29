package google

import "github.com/snyk/driftctl/enumeration/resource"

const GoogleComputeInstanceGroupResourceType = "google_compute_instance_group"

func initGoogleComputeInstanceGroupMetadata(resourceSchemaRepository resource.SchemaRepositoryInterface) {
	resourceSchemaRepository.SetResolveReadAttributesFunc(GoogleComputeInstanceGroupResourceType, func(res *resource.Resource) map[string]string {
		return map[string]string{
			"name":    *res.Attributes().GetString("name"),
			"project": *res.Attributes().GetString("project"),
			"zone":    *res.Attributes().GetString("location"),
		}
	})
	resourceSchemaRepository.SetHumanReadableAttributesFunc(GoogleComputeInstanceGroupResourceType, func(res *resource.Resource) map[string]string {
		attrs := make(map[string]string)
		if v := res.Attributes().GetString("name"); v != nil && *v != "" {
			attrs["Name"] = *v
		}
		return attrs
	})
	resourceSchemaRepository.SetFlags(GoogleComputeInstanceGroupResourceType, resource.FlagDeepMode)
}
