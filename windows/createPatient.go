package windows

import (
	. "client_main/funcs"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

	"fyne.io/fyne/v2/widget"
	_ "github.com/lib/pq"
)

type patients struct {
	id          int
	fio         string
	birthdate   string
	mestoprojiv string
	number      string
	email       string
}

func PatientCreation(win fyne.Window) fyne.CanvasObject {

	fioLabel := widget.NewLabel("ФИО пациента:")
	fioEntry := widget.NewEntry()
	birthdateLabel := widget.NewLabel("Дата рождения:")
	birthdateEntry := widget.NewEntry()
	mestoprojivLabel := widget.NewLabel("Место проживания:")
	mestoprojivEntry := widget.NewEntry()
	numberLabel := widget.NewLabel("Телефон:")
	numberEntry := widget.NewEntry()
	emailLabel := widget.NewLabel("Почта:")
	emailEntry := widget.NewEntry()

	applyCreationButton := widget.NewButton("Создать пациента", func() {
		CreatePatient(fioEntry.Text, birthdateEntry.Text, mestoprojivEntry.Text, numberEntry.Text, emailEntry.Text, win)
		

	})
	cancelCreationButton := widget.NewButton("Отмена", func() {
		win.Hide()
	})
	return container.NewVBox(

		fioLabel,
		fioEntry,
		birthdateLabel,
		birthdateEntry,
		mestoprojivLabel,
		mestoprojivEntry,
		numberLabel,
		numberEntry,
		emailLabel,
		emailEntry,
		applyCreationButton,
		cancelCreationButton,
	)
}
