package windows

import (
	. "client_main/funcs"

	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	_ "github.com/lib/pq"
)

func EditPatient(win fyne.Window, myApp fyne.App, searchEntry string) fyne.CanvasObject {
	patient := SelectPatient(searchEntry)
	search, err := strconv.Atoi(searchEntry)
	_ = err
	fioLabel := widget.NewLabel("ФИО пациента:")
	fioEntry := widget.NewEntry()
	fioEntry.SetText(patient[0][1])
	birthdateLabel := widget.NewLabel("Дата рождения:")
	birthdateEntry := widget.NewEntry()
	birthdateEntry.SetText(patient[0][2])
	mestoprojivLabel := widget.NewLabel("Место проживания:")
	mestoprojivEntry := widget.NewEntry()
	mestoprojivEntry.SetText(patient[0][3])
	numberLabel := widget.NewLabel("Телефон:")
	numberEntry := widget.NewEntry()
	numberEntry.SetText(patient[0][4])
	emailLabel := widget.NewLabel("Дата рождения")
	emailEntry := widget.NewEntry()
	emailEntry.SetText(patient[0][5])

	applyEditPatient := widget.NewButton("Применить изменения", func() {
		ApplyPatient(myApp, win, fioEntry.Text, birthdateEntry.Text, mestoprojivEntry.Text, numberEntry.Text, emailEntry.Text, search)

		win.Hide()
	})
	cancelEditButton := widget.NewButton("Отмена", func() {
		win.Hide()
	})

	return container.NewVBox(fioLabel,
		fioEntry,
		birthdateLabel,
		birthdateEntry,
		mestoprojivLabel,
		mestoprojivEntry,
		numberLabel,
		numberEntry,
		emailLabel,
		emailEntry,
		applyEditPatient,
		cancelEditButton,
	)

}
