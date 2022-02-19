package optional

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Address struct {
	City     String `json:"city"`
	State    String `json:"state"`
	PostCode Int64  `json:"postCode"`
}

type ID struct {
	Name    String  `json:"name"`
	Age     Int64   `json:"age"`
	Sex     String  `json:"sex"`
	Weight  Float64 `json:"weight"`
	Address Address `json:"address"`
}

func TestUnmarshal(t *testing.T) {
	jsonStr := `{"name":null,"age":30,"weight":10.12,"address":{"city":null,"state":"NY"}}`

	id := ID{}
	var err error
	err = json.Unmarshal([]byte(jsonStr), &id)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(id.Weight.Present())
	fmt.Println(id.Weight.ElseZero())

	fmt.Printf("Is Name of ID present: %v, value: %v\n", id.Name.Present(), id.Name.ElseZero())
	fmt.Println("-------------------------------------------------------")
	fmt.Printf("Is Age of ID present: %v, value: %v\n", id.Age.Present(), id.Age.ElseZero())
	fmt.Println("-------------------------------------------------------")
	fmt.Printf("Is Sex of ID present: %v, value: %v\n", id.Sex.Present(), id.Sex.ElseZero())
	fmt.Println("-------------------------------------------------------")
	fmt.Printf("Is City of Address of ID present %v, value: %v\n: ", id.Address.City.Present(), id.Address.City.ElseZero())
	fmt.Println("-------------------------------------------------------")
	fmt.Printf("Is State of Address of ID present %v, value: %v\n: ", id.Address.State.Present(), id.Address.State.ElseZero())
	fmt.Println("-------------------------------------------------------")
	fmt.Printf("Is PostCode of Address of ID present %v, value: %v\n: ", id.Address.PostCode.Present(), id.Address.PostCode.ElseZero())
}

type o struct {
	Name String `json:"name"`
	Age  Int64  `json:"age"`
}

func TestCopy(t *testing.T) {
	orig := make([]o, 0, 10)
	orig = append(orig, o{Name: OfString("a"), Age: OfInt64(1)})

	newest := make([]o, len(orig))
	copy(newest, orig)

	fmt.Println("New: ", newest[0].Name.ElseZero())
	fmt.Println("Orig: ", orig[0].Name.ElseZero())

	newStr := "b"
	orig[0].Name = OfString(newStr)
	//copy(newest, orig)

	fmt.Println("New: ", newest[0].Name.ElseZero())
	fmt.Println("Orig: ", orig[0].Name.ElseZero())
}
