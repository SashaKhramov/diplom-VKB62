package funcs

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	_ "github.com/lib/pq"
)

type Users struct {
	Id       int
	Username string
	Password string
	Realname string
	Role     string
}

type Roles struct {
	Id         int
	Rolename   string
	Main_ac    int
	Pay_ac     int
	Storage_ac int
	Setting_ac int
}

type Patients struct {
	Id          int
	Fio         string
	Birthdate   string
	Mestoprojiv string
	Number      string
	Email       string
}

type Pays struct {
	Id          int
	Paydate     string
	Summa       string
	Payform     string
	Patientname string
	PatientId   int
}

type Storage struct {
	Id          int
	Opertype    string
	Operdate    string
	Worker      string
	Productname string
	Count       int
}

var userRights [4]int
var currentUser Users

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func GetPatients() [][]string {
	psqlconn := "host=192.168.40.131 port=5432 user=postgres password=1 dbname=med sslmode=disable"
	db, err := sql.Open("postgres", psqlconn)
	querySQL := `SELECT Id,Fio, Birthdate, Mestoprojiv, Number, Email FROM patients`
	rows, err := db.Query(querySQL)
	CheckError(err)
	patientsSQL := []Patients{}
	patient := [][]string{}
	key := SetKey()
	for rows.Next() {
		p := Patients{}

		err := rows.Scan(&p.Id, &p.Fio, &p.Birthdate, &p.Mestoprojiv, &p.Number, &p.Email)
		CheckError(err)

		patientsSQL = append(patientsSQL, p)
	}

	for i := 0; i < len(patientsSQL); i++ {
		patientsSQL[i].Fio = Decrypt(patientsSQL[i].Fio, key)
		patientsSQL[i].Birthdate = Decrypt(patientsSQL[i].Birthdate, key)
		patientsSQL[i].Mestoprojiv = Decrypt(patientsSQL[i].Mestoprojiv, key)
		patientsSQL[i].Number = Decrypt(patientsSQL[i].Number, key)
		patientsSQL[i].Email = Decrypt(patientsSQL[i].Email, key)
		patient = append(patient, []string{strconv.Itoa(patientsSQL[i].Id), patientsSQL[i].Fio, patientsSQL[i].Birthdate, patientsSQL[i].Mestoprojiv, patientsSQL[i].Number, patientsSQL[i].Email})

	}
	return patient
}
func GetUsers() []Users {
	psqlconn := "host=192.168.40.131 port=5432 user=postgres password=1 dbname=med sslmode=disable"
	db, err := sql.Open("postgres", psqlconn)
	_ = db
	_ = err
	querySQL := `SELECT * FROM users`
	rows, err := db.Query(querySQL)
	CheckError(err)
	usersSQL := []Users{}

	for rows.Next() {
		u := Users{}
		err := rows.Scan(&u.Id, &u.Username, &u.Password, &u.Realname, &u.Role)
		CheckError(err)
		u.Username = Decrypt(u.Username, string(SetKey()))
		u.Role = Decrypt(u.Role, string(SetKey()))

		usersSQL = append(usersSQL, u)

	}

	return usersSQL
}
func GetOne(loginInput string) *sql.Row {
	psqlconn := "host=192.168.40.131 port=5432 user=postgres password=1 dbname=med sslmode=disable"
	db, err := sql.Open("postgres", psqlconn)
	_ = err
	querySQL := `SELECT * FROM users`

	rows, err := db.Query(querySQL)
	CheckError(err)
	usersSQL := []Users{}

	for rows.Next() {
		p := Users{}
		err := rows.Scan(&p.Id, &p.Username, &p.Password, &p.Realname, &p.Role)
		CheckError(err)

		usersSQL = append(usersSQL, p)
	}

	userId := 0
	_ = userId
	for i := 0; i < len(usersSQL); i++ {
		decryptedUser := Decrypt(usersSQL[i].Username, SetKey())

		if decryptedUser == loginInput {
			userId = usersSQL[i].Id

		}

	}
	newrows := db.QueryRow(`SELECT * FROM users WHERE id=$1`, userId)

	return newrows
}
func GetUser(loginInput string, myApp fyne.App) Users {
	psqlconn := "host=192.168.40.131 port=5432 user=postgres password=1 dbname=med sslmode=disable"
	db, err := sql.Open("postgres", psqlconn)
	_ = err
	rows := db.QueryRow(`SELECT * FROM users WHERE username=$1`, loginInput)
	user := Users{}
	err = rows.Scan(&user.Id, &user.Username, &user.Password, &user.Realname, &user.Role)
	if err != nil {
		myApp.SendNotification(fyne.NewNotification("Неверный логин или пароль", ""))
	}
	return user
}

func CreateUser(usernameEntry string, realName string, userRole string, userPassword string) {
	psqlconn := "host=192.168.40.131 port=5432 user=postgres password=1 dbname=med sslmode=disable"
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()
	username := Encrypt(usernameEntry, string(SetKey()))
	rolename := Encrypt(userRole, string(SetKey()))
	realname := Encrypt(realName, string(SetKey()))
	passhash := sha256.New()
	passhash.Write([]byte(userPassword))
	hashValue := passhash.Sum(nil)
	hashstring := hex.EncodeToString(hashValue)
	result, err := db.Exec(`insert into users (username,password,realname,role) values ($1,$2,$3,$4)`, username, hashstring, realname, rolename)
	fmt.Println()
	_ = result
	_ = err

}

func CreateRole() {
	psqlconn := "host=192.168.40.131 port=5432 user=postgres password=1 dbname=med sslmode=disable"
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()
	rolename := Encrypt("New_role", string(SetKey()))
	result, err := db.Exec(`insert into roles (rolename,main_ac,pay_ac,storage_ac,settings_ac) values ($1,1,0,0,0)`, rolename)
	_ = result
	_ = err
	fmt.Println(err)

}
func DeleteRole(rolename string) {
	psqlconn := "host=192.168.40.131 port=5432 user=postgres password=1 dbname=med sslmode=disable"
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()
	result, err := db.Exec(`delete from roles where rolename=$1;`, rolename)
	_ = result
	_ = err
}

func GetPays() [][]string {
	psqlconn := "host=192.168.40.131 port=5432 user=postgres password=1 dbname=med sslmode=disable"
	db, err := sql.Open("postgres", psqlconn)
	_ = err
	querySQL := `SELECT * FROM pays`
	rows, err := db.Query(querySQL)
	CheckError(err)
	paysSQL := []Pays{}
	pay := [][]string{}

	for rows.Next() {
		p := Pays{}
		err := rows.Scan(&p.Id, &p.Paydate, &p.Summa, &p.Payform, &p.Patientname, &p.PatientId)
		CheckError(err)

		paysSQL = append(paysSQL, p)
	}
	for i := 0; i < len(paysSQL); i++ {
		paysSQL[i].Paydate = Decrypt(paysSQL[i].Paydate, SetKey())
		paysSQL[i].Summa = Decrypt(paysSQL[i].Summa, SetKey())
		paysSQL[i].Payform = Decrypt(paysSQL[i].Payform, SetKey())
		paysSQL[i].Patientname = Decrypt(paysSQL[i].Patientname, SetKey())

		pay = append(pay, []string{strconv.Itoa(paysSQL[i].Id), paysSQL[i].Paydate, paysSQL[i].Summa, paysSQL[i].Payform, paysSQL[i].Patientname, strconv.Itoa(paysSQL[i].PatientId)})

	}
	return pay
}

func SelectPatient(searchEntry string) [][]string {
	psqlconn := "host=192.168.40.131 port=5432 user=postgres password=1 dbname=med sslmode=disable"
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()
	patientsSQL := []Patients{}
	patient := [][]string{}
	querySQL := `SELECT * FROM patients where id=$1`
	search, err := strconv.Atoi(searchEntry)
	CheckError(err)
	key := SetKey()
	rows, err := db.Query(querySQL, search)
	_ = err
	for rows.Next() {
		p := Patients{}
		err := rows.Scan(&p.Id, &p.Fio, &p.Birthdate, &p.Mestoprojiv, &p.Number, &p.Email)
		CheckError(err)

		patientsSQL = append(patientsSQL, p)
	}

	for i := 0; i < len(patientsSQL); i++ {
		patientsSQL[i].Fio = Decrypt(patientsSQL[i].Fio, key)
		patientsSQL[i].Birthdate = Decrypt(patientsSQL[i].Birthdate, key)
		patientsSQL[i].Mestoprojiv = Decrypt(patientsSQL[i].Mestoprojiv, key)
		patientsSQL[i].Number = Decrypt(patientsSQL[i].Number, key)
		patientsSQL[i].Email = Decrypt(patientsSQL[i].Email, key)
		patient = append(patient, []string{strconv.Itoa(patientsSQL[i].Id), patientsSQL[i].Fio, patientsSQL[i].Birthdate, patientsSQL[i].Mestoprojiv, patientsSQL[i].Number, patientsSQL[i].Email})

	}
	return patient
}

func ApplyPatient(myApp fyne.App, win fyne.Window, fioEntry string, birthdateEntry string, mestoprojivEntry string, numberEntry string, emailEntry string, search int) {

	psqlconn := "host=192.168.40.131 port=5432 user=postgres password=1 dbname=med sslmode=disable"
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()
	result, err := db.Exec(`update patients set fio=$1,birthdate=$2,mestoprojiv=$3,number=$4,email=$5 where id=$6`, fioEntry, birthdateEntry, mestoprojivEntry, numberEntry, emailEntry, search)
	_ = result
	_ = err
	myApp.SendNotification(fyne.NewNotification("Изменения применены", ""))

}

func SelectPay(searchEntry string) [][]string {
	psqlconn := "host=192.168.40.131 port=5432 user=postgres password=1 dbname=med sslmode=disable"
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()
	paysSQL := []Pays{}
	pay := [][]string{}
	querySQL := `SELECT * FROM pays where id=$1`
	search, err := strconv.Atoi(searchEntry)
	CheckError(err)

	rows, err := db.Query(querySQL, search)
	_ = err
	for rows.Next() {
		p := Pays{}
		err := rows.Scan(&p.Id, &p.Paydate, &p.Summa, &p.Payform, &p.Patientname, &p.PatientId)
		CheckError(err)

		paysSQL = append(paysSQL, p)
	}

	for i := 0; i < len(paysSQL); i++ {
		paysSQL[i].Paydate = Decrypt(paysSQL[i].Paydate, SetKey())
		paysSQL[i].Summa = Decrypt(paysSQL[i].Summa, SetKey())
		paysSQL[i].Payform = Decrypt(paysSQL[i].Payform, SetKey())
		paysSQL[i].Patientname = Decrypt(paysSQL[i].Patientname, SetKey())
		pay = append(pay, []string{strconv.Itoa(paysSQL[i].Id), paysSQL[i].Paydate, paysSQL[i].Summa, paysSQL[i].Payform, paysSQL[i].Patientname, strconv.Itoa(paysSQL[i].PatientId)})

	}
	return pay

}
func CreatePay(paydateEntry string, summaEntry string, payformEntry string, patientnameEntry string, patientidEntry string) {
	psqlconn := "host=192.168.40.131 port=5432 user=postgres password=1 dbname=med sslmode=disable"
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()

	paydateEntry = Encrypt(paydateEntry, string(SetKey()))
	summaEntry = Encrypt(summaEntry, string(SetKey()))
	payformEntry = Encrypt(payformEntry, string(SetKey()))
	patientnameEntry = Encrypt(patientnameEntry, string(SetKey()))

	result, err := db.Exec(`insert into pays (paydate,summa,payform,patientname,patientid) values ($1,$2,$3,$4,$5)`, paydateEntry, summaEntry, payformEntry, patientnameEntry, patientidEntry)

	_ = result
	_ = err

}

func GetStorage() [][]string {
	psqlconn := "host=192.168.40.131 port=5432 user=postgres password=1 dbname=med sslmode=disable"
	db, err := sql.Open("postgres", psqlconn)
	querySQL := `SELECT * FROM storage`
	rows, err := db.Query(querySQL)
	CheckError(err)
	storageSQL := []Storage{}
	store := [][]string{}

	for rows.Next() {
		p := Storage{}
		err := rows.Scan(&p.Id, &p.Opertype, &p.Operdate, &p.Worker, &p.Productname, &p.Count)
		CheckError(err)

		storageSQL = append(storageSQL, p)
	}
	for i := 0; i < len(storageSQL); i++ {
		store = append(store, []string{strconv.Itoa(storageSQL[i].Id), storageSQL[i].Opertype, storageSQL[i].Operdate, storageSQL[i].Worker, storageSQL[i].Productname, strconv.Itoa(storageSQL[i].Count)})

	}
	return store
}

func CreateStorage(opertypeEntry string, operdateEntry string, workerEntry string, productnameEntry string, countEntry string) {
	psqlconn := "host=192.168.40.131 port=5432 user=postgres password=1 dbname=med sslmode=disable"
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()
	opertypeEntry = Encrypt(opertypeEntry, string(SetKey()))
	operdateEntry = Encrypt(operdateEntry, string(SetKey()))
	workerEntry = Encrypt(workerEntry, string(SetKey()))
	productnameEntry = Encrypt(productnameEntry, string(SetKey()))
	countEntry = Encrypt(countEntry, string(SetKey()))
	result, err := db.Exec(`insert into storage (opertype,operdate,worker,productname,count) values ($1,$2,$3,$4,$5)`, opertypeEntry, operdateEntry, workerEntry, productnameEntry, countEntry)
	_ = result
	_ = err

}

func GetRoles() []Roles {
	psqlconn := "host=192.168.40.131 port=5432 user=postgres password=1 dbname=med sslmode=disable"
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()
	querySQL := `SELECT * FROM roles`
	rows, err := db.Query(querySQL)
	CheckError(err)
	rolesSQL := []Roles{}

	for rows.Next() {
		u := Roles{}
		err := rows.Scan(&u.Id, &u.Rolename, &u.Main_ac, &u.Pay_ac, &u.Storage_ac, &u.Setting_ac)
		CheckError(err)

		u.Rolename = Decrypt(u.Rolename, SetKey())
		rolesSQL = append(rolesSQL, u)

	}

	return rolesSQL
}

func SetRights(user Users) [4]int {
	psqlconn := "host=192.168.40.131 port=5432 user=postgres password=1 dbname=med sslmode=disable"
	db, err := sql.Open("postgres", psqlconn)

	var rightsArray [4]int
	rows, err := db.Query(`SELECT * FROM roles`)
	rolesSQL := []Roles{}
	for rows.Next() {
		p := Roles{}
		err := rows.Scan(&p.Id, &p.Rolename, &p.Main_ac, &p.Pay_ac, &p.Storage_ac, &p.Setting_ac)
		CheckError(err)

		rolesSQL = append(rolesSQL, p)
	}

	roleId := 0
	_ = roleId
	decryptedUser := Decrypt(user.Role, SetKey())
	for i := 0; i < len(rolesSQL); i++ {
		decryptedRole := Decrypt(rolesSQL[i].Rolename, SetKey())

		if decryptedUser == decryptedRole {
			roleId = rolesSQL[i].Id

		}

	}
	newrows := db.QueryRow(`SELECT * FROM roles WHERE id=$1`, roleId)
	roles := Roles{}
	err = newrows.Scan(&roles.Id, &roles.Rolename, &roles.Main_ac, &roles.Pay_ac, &roles.Storage_ac, &roles.Setting_ac)
	_ = err

	rightsArray[0] = roles.Main_ac
	rightsArray[1] = roles.Pay_ac
	rightsArray[2] = roles.Storage_ac
	rightsArray[3] = roles.Setting_ac

	return rightsArray
}

func CreatePatient(fioEntry string, birthdateEntry string, mestoprojivEntry string, numberEntry string, emailEntry string, win fyne.Window) {
	psqlconn := "host=192.168.40.131 port=5432 user=postgres password=1 dbname=med sslmode=disable"
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()
	var id int
	fioEntry = Encrypt(fioEntry, string(SetKey()))
	mestoprojivEntry = Encrypt(mestoprojivEntry, string(SetKey()))
	emailEntry = Encrypt(emailEntry, string(SetKey()))
	numberEntry = Encrypt(numberEntry, string(SetKey()))
	birthdateEntry = Encrypt(birthdateEntry, string(SetKey()))
	db.QueryRow(`insert into patients (fio,birthdate,mestoprojiv,number,email) values ($1,$2,$3,$4,$5) returning id`, fioEntry, birthdateEntry, mestoprojivEntry, numberEntry, emailEntry).Scan(&id)

	win.Hide()
}

func SysLogin(loginWindow fyne.Window, startWindow fyne.Window, myApp fyne.App, loginInput string, passwordInput string) [4]int {
	rows := GetOne(loginInput)
	user := Users{}
	err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Realname, &user.Role)
	if err != nil {
		myApp.SendNotification(fyne.NewNotification("Неверный логин или пароль", ""))
	}
	if passwordInput == "" || loginInput == "" {
		myApp.SendNotification(fyne.NewNotification("Неверный логин или пароль", ""))
	} else {
		passhash := sha256.New()
		passhash.Write([]byte(passwordInput))
		hashValue := passhash.Sum(nil)
		if hex.EncodeToString(hashValue) != user.Password {
			myApp.SendNotification(fyne.NewNotification("Неверный логин или пароль", ""))
		} else {
			userAc := [4]int{0, 0, 0, 0}
			_ = userAc
			currentUser = user
			userRights = SetRights(currentUser)

			loginWindow.Hide()
			startWindow.Show()
		}

	}
	return userRights
}
