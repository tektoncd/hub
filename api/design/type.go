package design

import (
	. "goa.design/goa/v3/dsl"
)

var Category = Type("Category", func() {
	Attribute("id", UInt, "ID is the unique id of the category", func() {
		Example("id", 1)
	})
	Attribute("name", String, "Name of the Category", func() {
		Example("name", "Notification")
	})
	Attribute("tags", ArrayOf(ResourceTag), "Tags associated with the category")
	Required("id", "name", "tags")
})

var ResourceTag = Type("ResourceTag", func() {
	TypeName("Tag")
	Attribute("id", UInt, "ID is the unique id of the tag", func() {
		Example("id", 1)
	})
	Attribute("name", String, "Name of the tag", func() {
		Example("name", "notification")
	})
	Required("id", "name")
})
