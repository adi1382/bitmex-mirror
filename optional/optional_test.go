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
	Address Address `json:"address"`
}

func TestUnmarshal(t *testing.T) {
	jsonStr := `{"name":null,"age":30,"address":{"city":null,"state":"NY"}}`

	id := ID{}
	var err error
	err = json.Unmarshal([]byte(jsonStr), &id)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("Is Name of ID present: %v, value: %v\n", id.Name.IsPresent(), id.Name.ElseZero())
	fmt.Println("-------------------------------------------------------")
	fmt.Printf("Is Age of ID present: %v, value: %v\n", id.Age.IsPresent(), id.Age.ElseZero())
	fmt.Println("-------------------------------------------------------")
	fmt.Printf("Is Sex of ID present: %v, value: %v\n", id.Sex.IsPresent(), id.Sex.ElseZero())
	fmt.Println("-------------------------------------------------------")
	fmt.Printf("Is City of Address of ID present %v, value: %v\n: ", id.Address.City.IsPresent(), id.Address.City.ElseZero())
	fmt.Println("-------------------------------------------------------")
	fmt.Printf("Is State of Address of ID present %v, value: %v\n: ", id.Address.State.IsPresent(), id.Address.State.ElseZero())
	fmt.Println("-------------------------------------------------------")
	fmt.Printf("Is PostCode of Address of ID present %v, value: %v\n: ", id.Address.PostCode.IsPresent(), id.Address.PostCode.ElseZero())
}
