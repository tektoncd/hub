package design

import (
	. "goa.design/goa/v3/dsl"
	expr "goa.design/goa/v3/expr"
)

var Category = Type("Category", func() {
	Attribute("id", UInt, "ID is the unique id of the category", func() {
		Example("id", 1)
	})
	Attribute("name", String, "Name of category", func() {
		Example("name", "Image-build")
	})
	Attribute("tags", ArrayOf(ResourceTag), "List of tag associated with category", func() {
		Example("tags", func() {
			Value([]expr.Val{
				{"id": 1, "name": "image-build"},
			})
		})
	})

	Required("id", "name", "tags")
})

var ResourceTag = Type("Tag", func() {
	Attribute("id", UInt, "Id is the unique id of tag", func() {
		Example("id", 1)
	})
	Attribute("name", String, "Name of tag", func() {
		Example("name", "image-build")
	})

	Required("id", "name")
})
