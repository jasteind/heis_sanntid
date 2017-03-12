package backup

import(
	"os"
	"encoding/json"
	"fmt"
	. "../types"
)



func SaveOrder(e ElevData_s){
	file, _ := os.Create("backup_local.json")
	buffer, _ := json.Marshal(e.Orders)
	file.WriteAt(buffer,0)
	fmt.Println("Saved orders: ")
	fmt.Println(string(buffer))
}


func LoadOrder(e *ElevData_s){
	file,_ := os.Open("backup_local.json")
	buffer := make([]byte, 1024)
	n, _ := file.ReadAt(buffer,0)
	json.Unmarshal(buffer[:n], &e.Orders)
	fmt.Println("Loaded orders: ")
	fmt.Println(string(buffer))
}
