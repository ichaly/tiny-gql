package core

import (
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {
	gql1 := `
	query {
		users {
			...userFields2
	
			created_at
			...userFields1
		}
	}

	fragment userFields1 on user {
		id
		email
	}
	
	fragment userFields2 on user {
		full_name
		phone
	}`
	lex1, err := newLexer([]byte(gql1))
	if err != nil {
		t.Fatalf("lex() error = %v", err)
	}
	gql2 := `query{users{...userFields2 created_at ...userFields1}}fragment userFields1 on user{id email}fragment userFields2 on user{full_name phone}`
	lex2, err := newLexer([]byte(gql2))
	if err != nil {
		t.Fatalf("lex() error = %v", err)
	}
	if len(lex1.items) != len(lex2.items) {
		t.Fatalf("lex() error = prettify and compress mode not equal")
	}
}

func TestSchema(t *testing.T) {
	schema := `type users {
	  phone: Text
	  category_counts: Json
	  avatar: Text
	  updated_at: TimestampWithTimeZone
	  stripe_id: Text
	  full_name: Text!
	  disabled: Boolean
	  created_at: TimestampWithTimeZone!
	  email: Text! @unique
	  id: Bigint! @id @unique
	}`
	lex, err := newLexer([]byte(schema))
	if err != nil {
		t.Fatalf("lex() error = %v", err)
	}
	got := string(lex.items[4].value)
	want := ":"
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("The values of %v is not %v\n", got, want)
	}
}
