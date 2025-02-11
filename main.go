package main

import (
	"client_main/funcs"
	. "client_main/funcs"
	"client_main/windows"

	"database/sql"

	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	_ "github.com/lib/pq"
)

func UpdateList(rolelist *widget.List) {
	roledata := GetRoles()
	rolelist = widget.NewList(
		func() int {
			return len(roledata)
		},
		func() fyne.CanvasObject {

			return widget.NewLabel("Template")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(roledata[id].Rolename)
		},
	)

}
func main() {

	var userRights [4]int

	psqlconn := "host=192.168.40.131 port=5432 user=postgres password=1 dbname=med sslmode=disable"
	db, err := sql.Open("postgres", psqlconn)
	_ = err
	usersSQL := funcs.GetUsers()
	patientsSQL := funcs.GetPatients()
	paysSQL := funcs.GetPays()
	storageSQL := funcs.GetStorage()
	myApp := app.New()
	myApp.Settings().SetTheme(theme.LightTheme())

	loginWindow := myApp.NewWindow("МедАссист логин")
	loginWindow.Resize(fyne.NewSize(300, 500))

	startWindow := myApp.NewWindow("МедАссист")
	startWindow.Resize(fyne.NewSize(1900, 900))
	startWindow.CenterOnScreen()

	mainWindow := myApp.NewWindow("МедАссист")
	mainWindow.Resize(fyne.NewSize(1900, 900))
	mainWindow.CenterOnScreen()

	payWindow := myApp.NewWindow("МедАссист")
	payWindow.Resize(fyne.NewSize(1900, 900))
	payWindow.CenterOnScreen()

	storageWindow := myApp.NewWindow("МедАссист")
	storageWindow.Resize(fyne.NewSize(1900, 900))
	storageWindow.CenterOnScreen()

	settingsWindow := myApp.NewWindow("МедАссист")
	settingsWindow.Resize(fyne.NewSize(1900, 900))
	settingsWindow.CenterOnScreen()

	createPatientWindow := myApp.NewWindow("МедАссист")
	createPatientWindow.Resize(fyne.NewSize(400, 600))
	createPatientWindow.CenterOnScreen()

	editPatientWindow := myApp.NewWindow("МедАссист")
	editPatientWindow.Resize(fyne.NewSize(400, 600))
	editPatientWindow.CenterOnScreen()

	deletePatientWindow := myApp.NewWindow("МедАссист")
	deletePatientWindow.Resize(fyne.NewSize(200, 200))
	deletePatientWindow.CenterOnScreen()

	createPayWindow := myApp.NewWindow("МедАссист")
	createPayWindow.Resize(fyne.NewSize(400, 600))
	createPayWindow.CenterOnScreen()

	editPayWindow := myApp.NewWindow("МедАссист")
	editPayWindow.Resize(fyne.NewSize(400, 600))
	editPayWindow.CenterOnScreen()

	createStorageWindow := myApp.NewWindow("МедАссист")
	createStorageWindow.Resize(fyne.NewSize(400, 600))
	createStorageWindow.CenterOnScreen()

	toolbar := widget.NewToolbar(widget.NewToolbarAction(theme.AccountIcon(), func() {
		if userRights[0] == 0 {
			myApp.SendNotification(fyne.NewNotification("Отказано в доступе", ""))
		} else {
			mainWindow.Show()
			startWindow.Hide()
			payWindow.Hide()
			storageWindow.Hide()
			settingsWindow.Hide()
		}

	}),
		widget.NewToolbarAction(theme.DocumentIcon(), func() {
			if userRights[1] == 0 {
				myApp.SendNotification(fyne.NewNotification("Отказано в доступе", ""))
			} else {
				payWindow.Show()
				startWindow.Hide()
				mainWindow.Hide()
				storageWindow.Hide()
				settingsWindow.Hide()
			}

		}),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() {
			if userRights[2] == 0 {
				myApp.SendNotification(fyne.NewNotification("Отказано в доступе", ""))
			} else {
				storageWindow.Show()
				startWindow.Hide()
				payWindow.Hide()
				mainWindow.Hide()
				settingsWindow.Hide()
			}

		}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			if userRights[3] == 0 {
				myApp.SendNotification(fyne.NewNotification("Отказано в доступе", ""))
			} else {
				settingsWindow.Show()
				startWindow.Hide()
				payWindow.Hide()
				mainWindow.Hide()
				storageWindow.Hide()
			}

		}),
	)
	content := container.NewBorder(toolbar, nil, nil, nil, nil)

	//Окно входа
	loginInput := widget.NewEntry()
	passwordInput := widget.NewPasswordEntry()
	loginWindowLayout := (fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),

		widget.NewLabelWithStyle("Вход в систему", fyne.TextAlignCenter, fyne.TextStyle{Symbol: true}),
		widget.NewLabelWithStyle("Имя пользователя:", fyne.TextAlignCenter, fyne.TextStyle{Symbol: true}),
		loginInput,
		widget.NewLabelWithStyle("Пароль:", fyne.TextAlignCenter, fyne.TextStyle{Symbol: true}),
		passwordInput,
		layout.NewSpacer(),
		widget.NewButton("Войти", func() {
			userRights = SysLogin(loginWindow, startWindow, myApp, loginInput.Text, passwordInput.Text)

		}),
		widget.NewButton("Выход", func() {
			myApp.Quit()
		}),
	))

	loginWindow.SetContent(loginWindowLayout)
	loginWindow.CenterOnScreen()
	loginWindow.Show()

	//Стартовое окно
	startLabel := widget.NewLabel("Добро пожаловать в МедАссист!")

	startWindow.SetContent(container.NewBorder(
		toolbar,
		nil,
		nil,
		nil,
		startLabel,
	))

	//Окно с пациентами

	search_data := widget.NewEntry()
	search_data.SetText("1")

	buttonBar := container.NewGridWithColumns(4,
		widget.NewButton("Поиск", func() {}),
		widget.NewButton("Изменить пациента", func() {
			editPatientWindow.Show()
		}),
		widget.NewButton("Удалить пациента", func() {}))
	searchbar := (fyne.NewContainerWithLayout(
		layout.NewGridLayout(2),
		search_data,
		buttonBar,
	))
	searchbar.Move(fyne.Position{0, 50})
	searchbar.Resize(fyne.NewSize(500, 30))

	//var data = [][]string{
	//[]string{"1", "Иванов Иван Иванович", "01.01.1999", "Ростов-на-Дону, ул.Краснойармейская", "+79182013432", "ivanov@mail.ru"}}
	list := widget.NewTableWithHeaders(
		func() (int, int) {
			return len(patientsSQL), len(patientsSQL[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(patientsSQL[i.Row][i.Col])

		},
	)
	list.SetColumnWidth(0, 50)
	list.SetColumnWidth(1, 200)
	list.SetColumnWidth(2, 120)
	list.SetColumnWidth(3, 400)
	list.SetColumnWidth(4, 150)
	list.SetColumnWidth(5, 200)
	list.CreateHeader = func() fyne.CanvasObject {

		return widget.NewButton("000", func() {})
	}
	list.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		b := o.(*widget.Button)
		if id.Col == -1 {
			b.SetText(strconv.Itoa(id.Row))
			b.Importance = widget.LowImportance
			b.Disable()
		} else {
			switch id.Col {
			case 0:
				b.SetText("ID")
			case 1:
				b.SetText("ФИО")
			case 2:
				b.SetText("Дата рождения")
			case 3:
				b.SetText("Место проживания")
			case 4:
				b.SetText("Телефон")
			case 5:
				b.SetText("E-mail")
				b.Enable()
				b.Refresh()
			}
		}
	}

	//окно создания пациента

	createPatientButton := widget.NewButton("Создать пациента", func() {
		createPatientWindow.SetContent(windows.PatientCreation(createPatientWindow))
		createPatientWindow.Show()
	})

	MW_Vbox := (fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		content,
		searchbar,
		createPatientButton,
	))
	refreshPatientsButton := widget.NewButton("Обновить", func() {
		patientsSQL := GetPatients()
		list := widget.NewTableWithHeaders(
			func() (int, int) {
				return len(patientsSQL), len(patientsSQL[0])
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("wide content")
			},
			func(i widget.TableCellID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(patientsSQL[i.Row][i.Col])

			},
		)

		list.SetColumnWidth(0, 50)
		list.SetColumnWidth(1, 200)
		list.SetColumnWidth(2, 120)
		list.SetColumnWidth(3, 400)
		list.SetColumnWidth(4, 150)
		list.SetColumnWidth(5, 200)
		list.CreateHeader = func() fyne.CanvasObject {

			return widget.NewButton("000", func() {})
		}
		list.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
			b := o.(*widget.Button)
			if id.Col == -1 {
				b.SetText(strconv.Itoa(id.Row))
				b.Importance = widget.LowImportance
				b.Disable()
			} else {
				switch id.Col {
				case 0:
					b.SetText("ID")
				case 1:
					b.SetText("ФИО")
				case 2:
					b.SetText("Дата рождения")
				case 3:
					b.SetText("Место проживания")
				case 4:
					b.SetText("Телефон")
				case 5:
					b.SetText("E-mail")
					b.Enable()
					b.Refresh()
				}
			}
		}
		mainWindow.SetContent(container.NewBorder(
			MW_Vbox,
			nil,
			nil,
			nil,
			list,
		),
		)
	})
	buttonBar = container.NewGridWithColumns(4,

		widget.NewButton("Изменить пациента", func() {
			editPatientWindow.SetContent(windows.EditPatient(editPatientWindow, myApp, search_data.Text))
			editPatientWindow.Show()
		}),
		widget.NewButton("Удалить пациента", func() {
			deletePatientWindow.SetContent(windows.DeletePatient(deletePatientWindow, search_data.Text))
			deletePatientWindow.Show()

		}),
		refreshPatientsButton,
	)
	searchbar = (fyne.NewContainerWithLayout(
		layout.NewGridLayout(2),
		search_data,
		buttonBar,
	))
	MW_Vbox = (fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		content,
		searchbar,
		createPatientButton,
	))

	mainWindow.SetContent(container.NewBorder(
		MW_Vbox,
		nil,
		nil,
		nil,
		list,
	),
	)

	//Окно с счетами

	paysearch_data := widget.NewEntry()

	paybuttonBar := container.NewGridWithColumns(3,

		widget.NewButton("Изменить cчет", func() {
			editPatientWindow.Show()
		}),
		widget.NewButton("Удалить счет", func() {}))

	paysearchbar := (fyne.NewContainerWithLayout(
		layout.NewGridLayoutWithColumns(2),
		paysearch_data,
		paybuttonBar,
	))
	paysearchbar.Resize(fyne.NewSize(600, 30))
	//var paydata = [][]string{
	//[]string{"1", "10.11.2024", "1500", "Наличные", "Иванов Иван Иванович", "1"}}
	paylist := widget.NewTableWithHeaders(
		func() (int, int) {
			return len(paysSQL), len(paysSQL[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(pay_i widget.TableCellID, pay_o fyne.CanvasObject) {
			pay_o.(*widget.Label).SetText(paysSQL[pay_i.Row][pay_i.Col])

		},
	)
	paylist.SetColumnWidth(0, 50)
	paylist.SetColumnWidth(1, 180)
	paylist.SetColumnWidth(2, 120)
	paylist.SetColumnWidth(3, 120)
	paylist.SetColumnWidth(4, 250)
	paylist.SetColumnWidth(5, 100)

	paylist.CreateHeader = func() fyne.CanvasObject {

		return widget.NewButton("000", func() {})
	}
	paylist.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		b := o.(*widget.Button)
		if id.Col == -1 {
			b.SetText(strconv.Itoa(id.Row))
			b.Importance = widget.LowImportance
			b.Disable()
		} else {
			switch id.Col {
			case 0:
				b.SetText("ID")
			case 1:
				b.SetText("Дата оплаты")
			case 2:
				b.SetText("Сумма")
			case 3:
				b.SetText("Форма оплаты")
			case 4:
				b.SetText("Пациент")
			case 5:
				b.SetText("ID Пациентa")
			}
		}
	}
	createPayButton := widget.NewButton("Создать счёт", func() {
		createPayWindow.SetContent(windows.PayCreation(createPayWindow, myApp, db))
		createPayWindow.Show()
	})

	PW_Vbox := (fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		content,
		paybuttonBar,
		createPayButton,
	))
	refreshPaysButton := widget.NewButton("Обновить", func() {
		paysSQL := GetPays()
		paylist := widget.NewTableWithHeaders(
			func() (int, int) {
				return len(paysSQL), len(paysSQL[0])
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("wide content")
			},
			func(pay_i widget.TableCellID, pay_o fyne.CanvasObject) {
				pay_o.(*widget.Label).SetText(paysSQL[pay_i.Row][pay_i.Col])

			},
		)
		paylist.SetColumnWidth(0, 50)
		paylist.SetColumnWidth(1, 180)
		paylist.SetColumnWidth(2, 120)
		paylist.SetColumnWidth(3, 120)
		paylist.SetColumnWidth(4, 250)
		paylist.SetColumnWidth(5, 100)

		paylist.CreateHeader = func() fyne.CanvasObject {

			return widget.NewButton("000", func() {})
		}
		paylist.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
			b := o.(*widget.Button)
			if id.Col == -1 {
				b.SetText(strconv.Itoa(id.Row))
				b.Importance = widget.LowImportance
				b.Disable()
			} else {
				switch id.Col {
				case 0:
					b.SetText("ID")
				case 1:
					b.SetText("Дата оплаты")
				case 2:
					b.SetText("Сумма")
				case 3:
					b.SetText("Форма оплаты")
				case 4:
					b.SetText("Пациент")
				case 5:
					b.SetText("ID Пациентa")
				}
			}
		}
		payWindow.SetContent(container.NewBorder(
			PW_Vbox,
			nil,
			nil,
			nil,
			paylist,
		),
		)

	})

	PW_Vbox = (fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		content,
		paysearchbar,
		refreshPaysButton,
		createPayButton,
	))
	payWindow.SetContent(container.NewBorder(
		PW_Vbox,
		nil,
		nil,
		nil,
		paylist,
	),
	)

	//Окно склада

	storagesearch_data := widget.NewEntry()
	storagesearchbar := (fyne.NewContainerWithLayout(
		layout.NewGridLayout(2),
		storagesearch_data,

		widget.NewButton("Изменить счёт", func() {

		}),
	))

	storagesearchbar.Resize(fyne.NewSize(600, 30))
	//var storageData = [][]string{
	//[]string{"1", "Прием", "09.11.2024", "Петров Петр Петрович", "Салфетки спиртовые антисептические(50 шт.)", "10"}}
	storageList := widget.NewTableWithHeaders(
		func() (int, int) {
			return len(storageSQL), len(storageSQL[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(store_i widget.TableCellID, store_o fyne.CanvasObject) {
			store_o.(*widget.Label).SetText(storageSQL[store_i.Row][store_i.Col])

		},
	)
	storageList.SetColumnWidth(0, 50)
	storageList.SetColumnWidth(1, 180)
	storageList.SetColumnWidth(2, 120)
	storageList.SetColumnWidth(3, 300)
	storageList.SetColumnWidth(4, 350)
	storageList.SetColumnWidth(5, 100)

	storageList.CreateHeader = func() fyne.CanvasObject {

		return widget.NewButton("000", func() {})
	}
	storageList.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		b := o.(*widget.Button)
		if id.Col == -1 {
			b.SetText(strconv.Itoa(id.Row))
			b.Importance = widget.LowImportance
			b.Disable()
		} else {
			switch id.Col {
			case 0:
				b.SetText("ID")
			case 1:
				b.SetText("Тип операции")
			case 2:
				b.SetText("Дата")
			case 3:
				b.SetText("Ответственное лицо")
			case 4:
				b.SetText("Название продукта")
			case 5:
				b.SetText("Количество")
			}
		}
	}
	createStorageButton := widget.NewButton("Создать операцию", func() {
		createStorageWindow.SetContent(windows.StorageCreation(createStorageWindow, myApp))
		createStorageWindow.Show()
	})
	STW_Vbox := (fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		content,
		storagesearchbar,
		createStorageButton,
	))
	refreshStorageButton := widget.NewButton("Обновить", func() {
		storageSQL := GetStorage()
		storageList := widget.NewTableWithHeaders(
			func() (int, int) {
				return len(storageSQL), len(storageSQL[0])
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("wide content")
			},
			func(store_i widget.TableCellID, store_o fyne.CanvasObject) {
				store_o.(*widget.Label).SetText(storageSQL[store_i.Row][store_i.Col])

			},
		)
		storageList.SetColumnWidth(0, 50)
		storageList.SetColumnWidth(1, 180)
		storageList.SetColumnWidth(2, 120)
		storageList.SetColumnWidth(3, 300)
		storageList.SetColumnWidth(4, 350)
		storageList.SetColumnWidth(5, 100)

		storageList.CreateHeader = func() fyne.CanvasObject {

			return widget.NewButton("000", func() {})
		}
		storageList.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
			b := o.(*widget.Button)
			if id.Col == -1 {
				b.SetText(strconv.Itoa(id.Row))
				b.Importance = widget.LowImportance
				b.Disable()
			} else {
				switch id.Col {
				case 0:
					b.SetText("ID")
				case 1:
					b.SetText("Тип операции")
				case 2:
					b.SetText("Дата")
				case 3:
					b.SetText("Ответственное лицо")
				case 4:
					b.SetText("Название продукта")
				case 5:
					b.SetText("Количество")
				}
			}
		}
		storageWindow.SetContent(container.NewBorder(
			STW_Vbox,
			nil,
			nil,
			nil,
			storageList,
		))
	})
	STW_Vbox = (fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		content,
		storagesearchbar,
		createStorageButton,
		refreshStorageButton,
	))
	storageWindow.SetContent(container.NewBorder(
		STW_Vbox,
		nil,
		nil,
		nil,
		storageList,
	),
	)

	//Окно с настройками

	//roledata := make([]string, 2)
	//roledata[0] = "admin"
	//roledata[1] = "Врач"
	roledata := GetRoles()
	rolelabel := widget.NewLabel("Название роли:")
	roleEntry := widget.NewEntry()
	
	rolePrivMain := widget.NewLabel("Доступ к Пациентам")
	rolePrivPay := widget.NewLabel("Доступ к Счетам")
	rolePrivStorage := widget.NewLabel("Доступ к Складу")
	rolePrivSettings := widget.NewLabel("Доступ к Настройкам")
	radioMain := widget.NewRadioGroup([]string{"Разрешено", "Запрещено"}, func(s string) {})
	radioPay := widget.NewRadioGroup([]string{"Разрешено", "Запрещено"}, func(s string) {})
	radioStorage := widget.NewRadioGroup([]string{"Разрешено", "Запрещено"}, func(s string) {})
	radioSettings := widget.NewRadioGroup([]string{"Разрешено", "Запрещено"}, func(s string) {})
	
	rolesArray := []string{}
	for i := 0; i < len(roledata); i++ {
		rolesArray = append(rolesArray, roledata[i].Rolename)
	}
	rolelist := widget.NewList(
		func() int {
			return len(roledata)
		},
		func() fyne.CanvasObject {

			return widget.NewLabel("Template")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(roledata[id].Rolename)
		},
	)
	rolelist.OnSelected = func(id widget.ListItemID) {
		roleEntry.SetText(roledata[id].Rolename)
		if roledata[id].Main_ac == 1 {
			radioMain.SetSelected("Разрешено")
		} else {
			radioMain.SetSelected("Запрещено")
		}
		if roledata[id].Pay_ac == 1 {
			radioPay.SetSelected("Разрешено")
		} else {
			radioPay.SetSelected("Запрещено")
		}
		if roledata[id].Storage_ac == 1 {
			radioStorage.SetSelected("Разрешено")
		} else {
			radioStorage.SetSelected("Запрещено")
		}
		if roledata[id].Setting_ac == 1 {
			radioSettings.SetSelected("Разрешено")
		} else {
			radioSettings.SetSelected("Запрещено")
		}

	}
	
	addroleButton := widget.NewButton("Добавить роль", func() {
		CreateRole()

	})
	deleteroleButton := widget.NewButton("Удалить роль", func() {
		DeleteRole(Encrypt(roleEntry.Text, SetKey()))
	})
	roleButtons := container.NewGridWithColumns(2, addroleButton, deleteroleButton)
	role_hbox := fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(2), rolelabel, roleEntry, rolePrivMain, radioMain, rolePrivPay, radioPay, rolePrivStorage, radioStorage, rolePrivSettings, radioSettings)
	SR_Hbox := container.NewHSplit(rolelist, container.NewBorder(roleButtons, nil, nil, nil, role_hbox))

	userlabel := widget.NewLabel("Имя пользователя:")
	usernameEntry := widget.NewEntry()
	usernameEntry.SetText("admin")
	realNameLabel := widget.NewLabel("ФИО:")
	realName := widget.NewEntry()
	realName.Resize(fyne.NewSize(300, 20))
	realName.SetText("Храмов Александр Максимович")
	userRoleLabel := widget.NewLabel("Роль пользователя:")
	userRole := widget.NewSelectEntry(rolesArray)
	userPasswordLabel := widget.NewLabel("Пароль пользователя:")
	userPassword := widget.NewPasswordEntry()
	userButton := widget.NewButton("Применить", func() {})
	adduserButton := widget.NewButton("Добавить пользователя", func() { CreateUser(usernameEntry.Text, realName.Text, userRole.Text, userPassword.Text) })

	userButtons := container.NewGridWithColumns(2, adduserButton)
	userlist := widget.NewList(
		func() int {
			return len(usersSQL)
		},
		func() fyne.CanvasObject {

			return widget.NewLabel("Template")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(usersSQL[id].Username)
		},
	)
	userlist.OnSelected = func(id widget.ListItemID) {
		usernameEntry.SetText(usersSQL[id].Username)
		roleEntry.SetText(usersSQL[id].Role)
		realName.SetText(usersSQL[id].Realname)
		userPassword.SetText(usersSQL[id].Password)
	}

	user_hbox := fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(2), realNameLabel, realName, userRoleLabel, userRole, userlabel, usernameEntry, userPasswordLabel, userPassword, userButton)
	SU_Hbox := container.NewHSplit(userlist, container.NewBorder(userButtons, nil, nil, nil, user_hbox))

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Роли", theme.ComputerIcon(), SR_Hbox),
		container.NewTabItemWithIcon("Пользователи", theme.AccountIcon(), SU_Hbox),
	)
	settingsWindow.SetContent(container.NewBorder(
		content,
		nil,
		nil,
		nil,
		tabs,
	))

	myApp.Run()

}
