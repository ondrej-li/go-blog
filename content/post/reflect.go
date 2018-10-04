package main

import "reflect"
import "fmt"

func ref(v interface{}) {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)
	fmt.Println("Typ je:", typ.Name())
	fmt.Println("Je hodnota pointer:", val.Kind() == reflect.Ptr)
	fmt.Println("Je hodnota struct:", val.Elem().Kind() == reflect.Struct) // Elem() dereferencuje pointer
	fmt.Println("Hodnota prvního pole struktury:", val.Elem().Field(0).String())
	for i := 0; i < val.Elem().NumField(); i++ {
		fmt.Println("Nalezeno pole:", typ.Elem().Field(i).Name, "typu:", val.Elem().Field(i).Type().Name(), "s tagem:", typ.Elem().Field(i).Tag.Get("tag"))
	}
	// nastavíme novou hodnotu prvního pole 
	val.Elem().Field(0).SetString("newval")
}

func main() {
	s := &struct {
		V string `tag:"test"`
	}{V: "val"}
	ref(s)
	fmt.Printf("struktura po spuštění: %+v\n", s)
}
