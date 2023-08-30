package genson

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

func TestMain(m *testing.M) {
	keyring.MockInit() // don't mess with the actual keystore

	os.Exit(m.Run())
}

func TestAny(t *testing.T) {
	type testStruct struct {
		A Any[struct {
			Int int
			Str string
		}]
		B Any[struct {
			Int int
			Str string
		}]
	}
	type nonStruct struct {
		X Any[int]
	}
	type invalidField struct {
		X Any[struct {
			F func()
		}]
	}

	t.Run("encode", func(t *testing.T) {
		t.Run("int sniff", func(t *testing.T) {
			ts := testStruct{
				A: Any[struct {
					Int int
					Str string
				}]{
					Payload: struct {
						Int int
						Str string
					}{
						Int: 42,
					},
				},
			}
			expect := `{"A":42,"B":null}`

			p, err := json.Marshal(&ts)
			require.NoError(t, err)
			require.Equal(t, expect, string(p))
		})

		t.Run("string sniff", func(t *testing.T) {
			ts := testStruct{
				A: Any[struct {
					Int int
					Str string
				}]{
					Payload: struct {
						Int int
						Str string
					}{
						Str: "towel",
					},
				},
			}
			expect := `{"A":"towel","B":null}`

			p, err := json.Marshal(&ts)
			require.NoError(t, err)
			require.Equal(t, expect, string(p))
		})

		t.Run("int explicit", func(t *testing.T) {
			ts := testStruct{
				A: Any[struct {
					Int int
					Str string
				}]{
					FieldName: "Int",
					Payload: struct {
						Int int
						Str string
					}{
						Int: 42,
						Str: "towel",
					},
				},
			}
			expect := `{"A":42,"B":null}`

			p, err := json.Marshal(&ts)
			require.NoError(t, err)
			require.Equal(t, expect, string(p))
		})

		t.Run("string explicit", func(t *testing.T) {
			ts := testStruct{
				A: Any[struct {
					Int int
					Str string
				}]{
					FieldName: "Str",
					Payload: struct {
						Int int
						Str string
					}{
						Int: 42,
						Str: "towel",
					},
				},
			}
			expect := `{"A":"towel","B":null}`

			p, err := json.Marshal(&ts)
			require.NoError(t, err)
			require.Equal(t, expect, string(p))
		})

		t.Run("non struct", func(t *testing.T) {
			ts := nonStruct{}

			_, err := json.Marshal(&ts)
			require.Error(t, err)
		})

		t.Run("invalid field sniff", func(t *testing.T) {
			ts := invalidField{}
			ts.X.Payload.F = func() {}
			expect := `{"X":null}`

			p, err := json.Marshal(&ts)
			require.NoError(t, err)
			require.Equal(t, expect, string(p))
		})

		t.Run("invalid field explicit", func(t *testing.T) {
			ts := invalidField{}
			ts.X.FieldName = "F"
			ts.X.Payload.F = func() {}
			expect := `{"X":null}`

			p, err := json.Marshal(&ts)
			require.NoError(t, err)
			require.Equal(t, expect, string(p))
		})
	})

	t.Run("decode", func(t *testing.T) {
		t.Run("int", func(t *testing.T) {
			data := `{"A":42}`
			expect := testStruct{
				A: Any[struct {
					Int int
					Str string
				}]{
					FieldName: "Int",
					Payload: struct {
						Int int
						Str string
					}{
						Int: 42,
					},
				},
			}

			var got testStruct
			err := json.Unmarshal([]byte(data), &got)
			require.NoError(t, err)
			require.Equal(t, expect, got)
		})

		t.Run("string", func(t *testing.T) {
			data := `{"A":"towel"}`
			expect := testStruct{
				A: Any[struct {
					Int int
					Str string
				}]{
					FieldName: "Str",
					Payload: struct {
						Int int
						Str string
					}{
						Str: "towel",
					},
				},
			}

			var got testStruct
			err := json.Unmarshal([]byte(data), &got)
			require.NoError(t, err)
			require.Equal(t, expect, got)
		})

		t.Run("float", func(t *testing.T) {
			data := `{"A":42.42}`

			var got testStruct
			err := json.Unmarshal([]byte(data), &got)
			require.Error(t, err)
		})

		t.Run("non struct", func(t *testing.T) {
			data := `{"X":"towel"}`
			ts := nonStruct{}

			err := json.Unmarshal([]byte(data), &ts)
			require.Error(t, err)
		})
	})
}
