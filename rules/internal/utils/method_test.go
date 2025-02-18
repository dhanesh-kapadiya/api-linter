// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package utils

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestIsCreateMethod(t *testing.T) {
	for _, test := range []struct {
		name string
		RPCs string
		want bool
	}{
		{"ValidBook", `
			rpc CreateBook(CreateBookRequest) returns (Book) {};
		`, true},
		{"InvalidNonCreate", `
			rpc GenerateBook(CreateBookRequest) returns (Book) {};
		`, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
				import "google/protobuf/field_mask.proto";
				service Foo {
					{{.RPCs}}
				}

				// This is at the top to make it retrievable
				// by the test code.
				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "books/{book}"
						singular: "book"
						plural: "books"
					};
				}

				message CreateBookRequest {
					// The parent resource where this book will be created.
					// Format: publishers/{publisher}
					string parent = 1;

					// The book to create.
					Book book = 2;
				}
			`, test)
			method := file.GetServices()[0].GetMethods()[0]
			got := IsCreateMethod(method)
			if got != test.want {
				t.Errorf("IsCreateMethod got %v, want %v", got, test.want)
			}
		})
	}
}

func TestIsUpdateMethod(t *testing.T) {
	for _, test := range []struct {
		name string
		RPCs string
		want bool
	}{
		{"ValidBook", `
			rpc UpdateBook(UpdateBookRequest) returns (Book) {};
		`, true},
		{"InvalidNonUpdate", `
			rpc UpsertBook(UpdateBookRequest) returns (Book) {};
		`, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
				import "google/protobuf/field_mask.proto";
				service Foo {
					{{.RPCs}}
				}

				// This is at the top to make it retrievable
				// by the test code.
				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "books/{book}"
						singular: "book"
						plural: "books"
					};
				}

				message UpdateBookRequest {
					Book book = 1;
					google.protobuf.FieldMask update_mask = 2;
				}
			`, test)
			method := file.GetServices()[0].GetMethods()[0]
			got := IsUpdateMethod(method)
			if got != test.want {
				t.Errorf("IsUpdateMethod got %v, want %v", got, test.want)
			}
		})
	}
}

func TestIsListMethod(t *testing.T) {
	for _, test := range []struct {
		name string
		RPCs string
		want bool
	}{
		{"ValidList", `
			rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {};
		`, true},
		{"InvalidListRevisionsMethod", `
			rpc ListBookRevisions(ListBooksRequest) returns (ListBooksResponse) {};
		`, false},
		{"InvalidNonList", `
			rpc EnumerateBooks(ListBooksRequest) returns (ListBooksResponse) {};
		`, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
				import "google/protobuf/field_mask.proto";
				service Foo {
					{{.RPCs}}
				}

				// This is at the top to make it retrievable
				// by the test code.
				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "books/{book}"
						singular: "book"
						plural: "books"
					};
				}

				message ListBooksRequest {
					string parent = 1;
					int32 page_size = 2;
					string page_token = 3;
				}

				message ListBooksResponse {
					repeated Book books = 1;
					string next_page_token = 2;
				}
			`, test)
			method := file.GetServices()[0].GetMethods()[0]
			got := IsListMethod(method)
			if got != test.want {
				t.Errorf("IsListMethod got %v, want %v", got, test.want)
			}
		})
	}
}

func TestIsListRevisionsMethod(t *testing.T) {
	for _, test := range []struct {
		name string
		RPCs string
		want bool
	}{
		{"ValidListRevisionsMethod", `
			rpc ListBookRevisions(ListBooksRequest) returns (ListBooksResponse) {};
		`, true},
		{"InvalidList", `
			rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {};
		`, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
				import "google/protobuf/field_mask.proto";
				service Foo {
					{{.RPCs}}
				}

				// This is at the top to make it retrievable
				// by the test code.
				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "books/{book}"
						singular: "book"
						plural: "books"
					};
				}

				message ListBooksRequest {
					string parent = 1;
					int32 page_size = 2;
					string page_token = 3;
				}

				message ListBooksResponse {
					repeated Book books = 1;
					string next_page_token = 2;
				}
			`, test)
			method := file.GetServices()[0].GetMethods()[0]
			got := IsListRevisionsMethod(method)
			if got != test.want {
				t.Errorf("IsListRevisionsMethod got %v, want %v", got, test.want)
			}
		})
	}
}

func TestGetListResourceMessage(t *testing.T) {
	for _, test := range []struct {
		name string
		RPCs string
		want string
	}{
		{"ValidBooks", `
			rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {};
		`, "Book"},
		{"InvalidNotListMethod", `
			rpc GetBook(ListBooksRequest) returns (Book) {};
		`, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
				import "google/protobuf/field_mask.proto";
				service Foo {
					{{.RPCs}}
				}

				// This is at the top to make it retrievable
				// by the test code.
				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "books/{book}"
						singular: "book"
						plural: "books"
					};
				}

				message ListBooksRequest {
					string parent = 1;
					int32 page_size = 2;
					string page_token = 3;
				}

				message ListBooksResponse {
					repeated Book books = 1;
					string next_page_token = 2;
				}
			`, test)
			method := file.GetServices()[0].GetMethods()[0]
			message := GetListResourceMessage(method)
			got := ""
			if message != nil {
				got = message.GetName()
			}
			if got != test.want {
				t.Errorf("GetListResourceMessage got %q, want %q", got, test.want)
			}
		})
	}
}
