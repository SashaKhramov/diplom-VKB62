package windows

import (
	. "client_main/funcs"
	"database/sql"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "github.com/lib/pq"
)

type storage struct {
	id          int
	opertype    string
	operdate    string
	worker      string
	productname string
	count       int
}

type users struct {
	id       int
	username string
	password string
	realname string
	role     string
}

func getUsers() []users {
	psqlconn := "host=192.168.40.131 port=5432 user=postgres password=1 dbname=med sslmode=disable"
	db, err := sql.Open("postgres", psqlconn)
	_ = err
	querySQL := `SELECT * FROM users`
	rows, err := db.Query(querySQL)
	CheckError(err)
	usersSQL := []users{}

	for rows.Next() {
		u := users{}
		err := rows.Scan(&u.id, &u.username, &u.password, &u.realname, &u.role)
		CheckError(err)
		usersSQL = append(usersSQL, u)

	}

	return usersSQL
}

func StorageCreation(win fyne.Window, myApp fyne.App) fyne.CanvasObject {

	userSelectWindow := myApp.NewWindow("МедАссист")
	userSelectWindow.Resize(fyne.NewSize(400, 600))
	userSelectWindow.CenterOnScreen()
	usersSQL := getUsers()
	userlist := widget.NewList(
		func() int {
			return len(usersSQL)
		},
		func() fyne.CanvasObject {

			return widget.NewLabel("Template")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(usersSQL[id].realname)
		},
	)

	opertypedata := make([]string, 2)
	opertypedata[0] = "Прием"
	opertypedata[1] = "Выдача"

	opertypeLabel := widget.NewLabel("Тип операции:")
	opertypeEntry := widget.NewSelectEntry(opertypedata)
	operdateLabel := widget.NewLabel("Дата операции:")
	operdateEntry := widget.NewEntry()
	workerLabel := widget.NewLabel("Сотрудник:")
	workerEntry := widget.NewEntry()
	productnameLabel := widget.NewLabel("Название продукта:")
	productnameEntry := widget.NewEntry()
	countLabel := widget.NewLabel("Количество:")
	countEntry := widget.NewEntry()
	userlist.OnSelected = func(id widget.ListItemID) {
		workerEntry.SetText(usersSQL[id].realname)
		userSelectWindow.Hide()
	}
	applyCreationButton := widget.NewButton("Создать операцию", func() {

		win.Hide()
		//insertstm := "insert into patients (fio,birthdate,mestoprojiv,number,email) values ($1)"
		//input := [5]string{fioEntry.Text, birthdateEntry.Text, mestoprojivEntry.Text, numberEntry.Text, emailEntry.Text}

	})
	selectWorkerButton := widget.NewButton("Выбрать сотрудника", func() {
		cancelButton := widget.NewButton("Отмена", func() {
			userSelectWindow.Hide()
		})
		userSelectWindow.SetContent(container.NewBorder(nil, cancelButton, nil, nil,
			userlist,
		))
		userSelectWindow.Show()
	})
	cancelCreationButton := widget.NewButton("Отмена", func() {
		win.Hide()
	})
	emptyBox := container.NewVBox()
	return container.NewVBox(
		opertypeLabel,
		opertypeEntry,
		operdateEntry,
		operdateLabel,
		operdateEntry,
		workerLabel,
		workerEntry,
		selectWorkerButton,
		emptyBox,
		productnameLabel,
		productnameEntry,
		countLabel,
		countEntry,
		applyCreationButton,
		cancelCreationButton,
	)
}
