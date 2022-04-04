package gents

import (
	"fmt"
	"go/token"
	"go/types"
	"net/http"
	"testing"

	tstypes "github.com/benoitkugler/structgen/ts-types"
)

func TestGenerate(t *testing.T) {
	apis := Service{
		{
			Url: "/samlskm/", Method: http.MethodPost, Contrat: Contrat{
				Input:       TypeNoId{Type: types.NewArray(types.Typ[types.Byte], 5), NoId: true},
				HandlerName: "M1",
				Return:      types.NewSlice(types.Typ[types.Int]),
			},
		},
		{
			Url: "/samlskm/:param1", Method: http.MethodGet, Contrat: Contrat{
				HandlerName: "M2",
				QueryParams: []TypedParam{
					{Name: "arg1", Type: tstypes.TsString},
					{Name: "arg2", Type: tstypes.TsNumber},
					{Name: "arg3", Type: tstypes.TsBoolean},
				},
				Return: types.NewSlice(types.Typ[types.Int]),
			},
		},
	}
	fmt.Println(apis.Render(nil, types.NewScope(nil, token.NoPos, token.NoPos, "")))
}

func TestGenerateMaps(t *testing.T) {
	api := API{
		Url:    "/samlskm/",
		Method: http.MethodPost,
		Contrat: Contrat{
			Input:       TypeNoId{Type: types.NewArray(types.Typ[types.Byte], 5), NoId: true},
			HandlerName: "M1",
			Return:      types.NewMap(types.Typ[types.String], types.NewSlice(types.Typ[types.Int])),
		},
	}
	fmt.Println(api.generateMethod(types.NewScope(nil, token.NoPos, token.NoPos, "")))
}
