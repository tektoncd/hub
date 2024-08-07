/*
Package codegen contains data structures and algorithms used by the Goa code
generation tool.

In particular package codegen defines the data structure that represents a
generated file (see File) which is composed of sections, each corresponding to a
Go text template and accompanying data used to render the final code.

The package also includes functions that generate code to transform a
given type into another (see GoTransform).
*/
package codegen
