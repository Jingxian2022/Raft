package tapestry

import (
	"fmt"
	"testing"
)

func TestRoutingTableAdd(t *testing.T) {
	localId := MakeID(1)

	newId1, err := ParseID("2300000000000000000000000000000000000000")
	if err != nil {
		t.Errorf("Parse ID error!")
	}
	table := NewRoutingTable(localId)
	table.Add(newId1)
	if len(table.Rows[0][2]) != 1 && table.Rows[0][2][0] != newId1 {
		t.Errorf("Wrong answer when adding node to an empty slot!")
	}

	newId2, err := ParseID("2100000000000000000000000000000000000000")
	if err != nil {
		t.Errorf("Parse ID error!")
	}
	table.Add(newId2)
	if len(table.Rows[0][2]) != 2 && table.Rows[0][2][0] != newId2 {
		t.Errorf("Wrong answer when adding node to a slot with one element!")
	}

	newId3, err := ParseID("2200000000000000000000000000000000000000")
	if err != nil {
		t.Errorf("Parse ID error!")
	}
	table.Add(newId3)
	if len(table.Rows[0][2]) != 3 && table.Rows[0][2][1] != newId3 {
		t.Errorf("Wrong answer when adding node to a slot with two elements!")
	}

	newId4, err := ParseID("2000000000000000000000000000000000000000")
	if err != nil {
		t.Errorf("Parse ID error!")
	}
	table.Add(newId4)
	if len(table.Rows[0][2]) != 3 && table.Rows[0][2][0] != newId4 {
		t.Errorf("Wrong answer when adding node to a full slot!")
	}

	for i := 0; i < DIGITS; i++ {
		for j := 0; j < BASE; j++ {
			for k := 0; k < len(table.Rows[i][j]); k++ {
				fmt.Printf("Row %v slot %v element %v: %v\n", i, j, k, table.Rows[i][j][k])
			}
		}
	}
	//
	//newId5, err := ParseID("1100000000000000000000000000000000000004")
	//if err != nil {
	//	t.Errorf("Parse ID error!")
	//}
	//table.Add(newId5)
	//if len(table.Rows[1][1]) != 3 && table.Rows[1][1][2] != newId5 {
	//	t.Errorf("Wrong answer when adding node to a full slot!")
	//}
}

func TestRoutingTableGetLevel(t *testing.T) {
	localId := MakeID(1)

	newId1, err := ParseID("1100000000000000000000000000000000000000")
	newId2, err := ParseID("1100000000000000000000000000000000000001")
	newId3, err := ParseID("1100000000000000000000000000000000000005")
	newId4, err := ParseID("1100000000000000000000000000000000000007")
	newId5, err := ParseID("1100000000000000000000000000000000000001")
	if err != nil {
		t.Errorf("Parse ID error!")
	}

	table := NewRoutingTable(localId)
	table.Add(newId1)
	table.Add(newId2)
	table.Add(newId3)
	table.Add(newId4)
	table.Add(newId5)

	nodeIds := table.GetLevel(1)
	for _, id := range nodeIds {
		fmt.Printf("%v\n", id)
		if id == table.localId {
			t.Errorf("Not excluding local node id!")
		}
	}
}
