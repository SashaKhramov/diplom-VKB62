package windows

import (
	. "client_main/funcs"
	"database/sql"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "github.com/lib/pq"
)



func PayCreation(win fyne.Window, myApp fyne.App, db *sql.DB) fyne.CanvasObject {
	patientSelectWindow := myApp.NewWindow("МедАссист")
	patientSelectWindow.Resize(fyne.NewSize(400, 600))
	patientSelectWindow.CenterOnScreen()
	patientsSQL := GetPatients()
	patientlist := widget.NewList(
		func() int {
			return len(patientsSQL)
		},
		func() fyne.CanvasObject {

			return widget.NewLabel("Template")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(patientsSQL[id][0] + " " + patientsSQL[id][1] + " " + patientsSQL[id][2])
		},
	)
	patientlist.OnSelected = func(id widget.ListItemID) {

	}
	payformdata := make([]string, 2)
	payformdata[0] = "Наличные"
	payformdata[1] = "Безналичные"

	paydateLabel := widget.NewLabel("Дата оплаты")
	paydateEntry := widget.NewEntry()
	summaLabel := widget.NewLabel("Сумма:")
	summaEntry := widget.NewEntry()
	payformLabel := widget.NewLabel("Форма оплаты:")
	payformEntry := widget.NewSelectEntry(payformdata)
	patientnameLabel := widget.NewLabel("Имя пациента")
	patientnameEntry := widget.NewEntry()
	patientidLabel := widget.NewLabel("Id пациента")
	patientidEntry := widget.NewEntry()
	patientlist.OnSelected = func(id widget.ListItemID) {
		patientnameEntry.SetText(patientsSQL[id][1])
		patientidEntry.SetText(patientsSQL[id][0])
		patientSelectWindow.Hide()
	}
	patientidEntry.Disable()
	applyCreationButton := widget.NewButton("Создать счёт", func() {
		CreatePay(paydateEntry.Text, summaEntry.Text, payformEntry.Text, patientnameEntry.Text, patientidEntry.Text)

		win.Hide()
		//insertstm := "insert into patients (fio,birthdate,mestoprojiv,number,email) values ($1)"
		//input := [5]string{fioEntry.Text, birthdateEntry.Text, mestoprojivEntry.Text, numberEntry.Text, emailEntry.Text}

	})
	selectPatientButton := widget.NewButton("Выбрать пациента", func() {
		cancelButton := widget.NewButton("Отмена", func() {
			patientSelectWindow.Hide()
		})
		patientSelectWindow.SetContent(container.NewBorder(nil, cancelButton, nil, nil,
			patientlist,
		))
		patientSelectWindow.Show()
	})
	cancelCreationButton := widget.NewButton("Отмена", func() {
		win.Hide()
	})
	emptyBox := container.NewVBox()
	return container.NewVBox(

		paydateLabel,
		paydateEntry,
		summaLabel,
		summaEntry,
		payformLabel,
		payformEntry,
		patientnameLabel,
		patientnameEntry,
		selectPatientButton,
		emptyBox,
		patientidLabel,
		patientidEntry,
		applyCreationButton,
		cancelCreationButton,
	)
}
