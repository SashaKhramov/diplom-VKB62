package windows

import (
	. "client_main/funcs"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "github.com/lib/pq"
)

func PayEdit(win fyne.Window, myApp fyne.App, searchEntry string) fyne.CanvasObject {
	patientSelectWindow := myApp.NewWindow("МедАссист")
	patientSelectWindow.Resize(fyne.NewSize(400, 600))
	patientSelectWindow.CenterOnScreen()
	patientsSQL := GetPatients()
	pay := SelectPatient(searchEntry)
	search, err := strconv.Atoi(searchEntry)
	_ = search
	_ = err
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
	paydateEntry.SetText(pay[0][1])
	summaLabel := widget.NewLabel("Сумма:")
	summaEntry := widget.NewEntry()
	summaEntry.SetText(pay[0][2])
	payformLabel := widget.NewLabel("Форма оплаты:")
	payformEntry := widget.NewSelectEntry(payformdata)
	payformEntry.SetText(pay[0][3])
	patientnameLabel := widget.NewLabel("Имя пациента")
	patientnameEntry := widget.NewEntry()
	patientnameEntry.SetText(pay[0][4])
	patientidLabel := widget.NewLabel("Id пациента")
	patientidEntry := widget.NewEntry()
	patientidEntry.SetText(pay[0][5])
	patientlist.OnSelected = func(id widget.ListItemID) {
		patientnameEntry.SetText(patientsSQL[id][1])
		patientidEntry.SetText(patientsSQL[id][0])
		patientSelectWindow.Hide()
	}
	patientidEntry.Disable()

	applyCreationButton := widget.NewButton("Изменить счёт", func() {
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
