// Package profile contains the profile value object attached to a person.
//
// Profile values are copied on change. Methods named With... return modified
// copies so the person aggregate remains responsible for committing profile
// changes to its state.
package profile
