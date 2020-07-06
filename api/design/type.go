package design

import (
	. "goa.design/goa/v3/dsl"
)

var Tag = Type("Tag", func() {
	Attribute("id", UInt, "ID is the unique id of tag", func() {
		Example("id", 1)
	})
	Attribute("name", String, "Name of tag", func() {
		Example("name", "image-build")
	})

	Required("id", "name")
})

var Tags = ArrayOf(Tag)

var Category = Type("Category", func() {
	Attribute("id", UInt, "ID is the unique id of the category", func() {
		Example("id", 1)
	})
	Attribute("name", String, "Name of category", func() {
		Example("name", "Image Builder")
	})
	Attribute("tags", Tags, "List of tags associated with the category", func() {
		Example("tags", func() {
			Value([]Val{
				{"id": 1, "name": "image-build"},
				{"id": 2, "name": "kaniko"},
			})
		})
	})

	Required("id", "name", "tags")
})

var Catalog = Type("Catalog", func() {
	Attribute("id", UInt, "ID is the unique id of the catalog", func() {
		Example("id", 1)
	})
	Attribute("type", String, "Type of catalog", func() {
		Enum("official", "community")
		Example("type", "community")
	})

	Required("id", "type")
})

var Version = Type("Version", func() {
	Attribute("id", UInt, "ID is the unique id of resource's version", func() {
		Example("id", 1)
	})
	Attribute("version", String, "Version of resource", func() {
		Example("version", "0.1")
	})
	Attribute("displayName", String, "Display name of version", func() {
		Example("displayName", "Buildah")
	})
	Attribute("description", String, "Description of version", func() {
		Example("descripiton", "Buildah task builds source into a container image and then pushes it to a container registry.")
	})
	Attribute("minPipelinesVersion", String, "Minimum pipelines version the resource's version is compatible with", func() {
		Example("minPipelinesVersion", "0.12.1")
	})
	Attribute("rawURL", String, "Raw URL of resource's yaml file of the version", func() {
		Format(FormatURI)
		Example("rawURL", "https://raw.githubusercontent.com/tektoncd/catalog/master/task/buildah/0.1/buildah.yaml")
	})
	Attribute("webURL", String, "Web URL of resource's yaml file of the version", func() {
		Format(FormatURI)
		Example("webURL", "https://github.com/tektoncd/catalog/blob/master/task/buildah/0.1/buildah.yaml")
	})
	Attribute("updatedAt", String, "Timestamp when version was last updated", func() {
		Format(FormatDateTime)
		Example("updatedAt", "2020-01-01 12:00:00 +0000 UTC")
	})

	Required("id", "version", "description", "displayName", "minPipelinesVersion", "rawURL", "webURL", "updatedAt")
})

var Resource = ResultType("application/vnd.hub.resource", func() {
	Description("The resource type describes resource information.")
	TypeName("Resource")

	Attribute("id", UInt, "ID is the unique id of the resource", func() {
		Example("id", 1)
	})
	Attribute("name", String, "Name of resource", func() {
		Example("name", "buildah")
	})
	Attribute("catalog", Catalog, "Type of catalog to which resource belongs", func() {
		Example("catalog", func() {
			Value(Val{"id": 1, "type": "community"})
		})
	})
	Attribute("type", String, "Type of resource", func() {
		Example("type", "task")
	})
	Attribute("latestVersion", Version, "Latest version of resource", func() {
		Example("latestVersion", func() {
			Value(Val{
				"id":                  1,
				"version":             "0.1",
				"description":         "Buildah task builds source into a container image and then pushes it to a container registry.",
				"displayName":         "Buildah",
				"minPipelinesVersion": "0.12.1",
				"rawURL":              "https://raw.githubusercontent.com/tektoncd/catalog/master/task/buildah/0.1/buildah.yaml",
				"webURL":              "https://github.com/tektoncd/catalog/blob/master/task/buildah/0.1/buildah.yaml",
				"updatedAt":           "2020-01-01 12:00:00 +0000 UTC",
			})
		})
	})
	Attribute("tags", Tags, "Tags related to resource", func() {
		Example("tags", func() {
			Value([]Val{
				{"id": 1, "name": "image-build"},
			})
		})
	})
	Attribute("rating", Float64, "Rating of resource", func() {
		Example("rating", 4.3)
	})

	Required("id", "name", "catalog", "type", "latestVersion", "tags", "rating")
})
