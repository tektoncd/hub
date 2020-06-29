package design

import (
	. "goa.design/goa/v3/dsl"
	expr "goa.design/goa/v3/expr"
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
			Value([]expr.Val{
				{"id": 1, "name": "image-build"},
				{"id": 2, "name": "kaniko"},
			})
		})
	})

	Required("id", "name", "tags")
})
