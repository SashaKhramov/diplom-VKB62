package windows

import (
	. "client_main/funcs"
	"database/sql"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "github.com/lib/pq"
)

func DeletePatient(win fyne.Window, searchEntry string) fyne.CanvasObject {
	confirmLabel := widget.NewLabel("Вы действительно хотите удалить карту пациента?")
	fmt.Println(searchEntry)
	confirmButton := widget.NewButton("Да", func() {
		psqlconn := "host=192.168.40.131 port=5432 user=postgres password=1 dbname=med sslmode=disable"
		db, err := sql.Open("postgres", psqlconn)
		CheckError(err)
		result, err := db.Exec(`delete from patients where id=$1`, searchEntry)
		_ = result
		CheckError(err)
		defer db.Close()

		win.Hide()
	})
	denyButton := widget.NewButton("Нет", func() {
		win.Hide()
	})
	buttons := container.NewGridWithColumns(2, confirmButton, denyButton)
	return container.NewCenter(container.NewVBox(confirmLabel, buttons))
}
