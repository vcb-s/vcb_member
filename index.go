package main

import (
	"fmt"
	"vcb_member/helper"
	"vcb_member/models"
)

func main() {
	helper.GenID()

	fmt.Print(models.Conf)

}
