package menu

import (
	"email-app/src/decorative"
	"email-app/src/entity"
	"email-app/src/features/authentication"
	"email-app/src/features/emails"
	"email-app/src/util"
	"encoding/json"
	"fmt"
	"os"
)

var routes entity.USER_TYPE_LIST
var userTypeIndex, routeIndex, choiceIndex int

var USERS entity.USER_LIST
var ADMINS entity.USER_ADMIN_LIST
var EMAILS entity.EMAIL_LIST
var CurrentLogged entity.LoggedUser

var printStatus string
var printText string

var selectedEmailIdx entity.Email

// Reply email purpose
var replyEmail string

func InitRoutes() {

	// Inisialisasi Route
	userTypeIndex = -1 // 0 for admin 1 for user
	routeIndex = -1    // menu
	choiceIndex = -1   // sub menu

	// Init Seed
	seedData()

	CurrentLogged = entity.LoggedUser{Id: -1, Name: "", Email: "", Role: -1}

	routes[0] = entity.UserType{
		UserType: util.ADMIN_TYPE,
		RouteList: [10]entity.Route{
			{
				RouteName: util.ADMIN_AUTH_MENU,
				RouteFunc: func(printStatus *string, printText *string, choiceIndex *int) {
					decorative.HeaderTemplate()
					// Menambahkan menu user dan admin
					decorative.PrintLine()
					decorative.PrintTitle(" Admin Auth Menu ")
					decorative.PrintDecorativeLine()
					decorative.PrintMenu(1, "Register")
					decorative.PrintMenu(2, "Login")
					decorative.PrintMenu(3, "Back")
					decorative.PrintDecorativeLine()
					decorative.PrintInstruction(" Choose the number of the menu to continue ")
					decorative.PrintBottomLine()

					navigateInputIndex(3, choiceIndex)

				},
				ChoiceList: [4]entity.Choice{
					{
						ChoiceText: util.ADMIN_REGISTER_MENU,
						ChoiceFunc: func(userTypeIndex *int, routeIndex *int, choiceIndex *int) {
							decorative.HeaderTemplate()
							decorative.HeaderPage[string]("Admin Register Page")

							// Function to register user
							loggedIn := CurrentLogged.Id != -1 && CurrentLogged.Role == 1 // check is user loggedin

							for !loggedIn {
								name, email, password := authentication.InputUserRegister(func() {
									decorative.ResetPrintStatus(&printStatus, &printText)
									navigateRoute(util.ADMIN_AUTH_MENU, userTypeIndex, routeIndex, choiceIndex)
								})
								err, status := authentication.RegisterAdmin(name, email, password, &ADMINS)

								if err {
									util.ClearScreen()
									printStatus = util.PRINT_STATUS_ERROR
									printText = status
									navigateRoute(util.ADMIN_AUTH_MENU, userTypeIndex, routeIndex, choiceIndex)
								} else {
									printStatus = util.PRINT_STATUS_SUCCESS
									printText = status
									loggedIn = true
								}

								decorative.PrintStatus(printStatus, printText)
							}

							decorative.PrintText("Press any key and enter to navigate to login page: ")
							var key int
							fmt.Scan(&key)
							navigateRoute(util.ADMIN_LOGIN_MENU, userTypeIndex, routeIndex, choiceIndex)
						},
					},
					{
						ChoiceText: util.ADMIN_LOGIN_MENU,
						ChoiceFunc: func(userTypeIndex *int, routeIndex *int, choiceIndex *int) {

							// Function to login user
							loggedIn := CurrentLogged.Id != -1
							for !loggedIn {
								decorative.HeaderTemplate()
								decorative.HeaderPage[string]("Admin Login Page")
								if printStatus == util.PRINT_STATUS_ERROR {
									decorative.PrintStatus(printStatus, printText)
								} else if printStatus == util.PRINT_STATUS_NOTHING {
									decorative.PrintNothing()
								}

								email, password := authentication.InputUserLogin(func() {
									decorative.ResetPrintStatus(&printStatus, &printText)
									navigateRoute(util.ADMIN_AUTH_MENU, userTypeIndex, routeIndex, choiceIndex)
								})
								err, message := authentication.LoginAsAdmin(email, password, ADMINS, &CurrentLogged)

								if err {
									printStatus = util.PRINT_STATUS_ERROR
									printText = message
									util.ClearScreen()
								} else {
									loggedIn = true
								}
							}

							navigateRoute(util.ADMIN_APPROVAL_MENU, userTypeIndex, routeIndex, choiceIndex)
						},
					},
					{
						ChoiceText: util.ADMIN_AUTH_MENU_BACK,
						ChoiceFunc: func(userTypeIndex *int, routeIndex *int, choiceIndex *int) {
							decorative.ResetPrintStatus(&printStatus, &printText)
							navigateRoute(util.ADMIN_TYPE, userTypeIndex, routeIndex, choiceIndex)
						},
					},
				},
			},
			{
				RouteName: util.ADMIN_APPROVAL_MENU,
				RouteFunc: func(printStatus *string, printText *string, choiceIndex *int) {
					decorative.HeaderTemplate()
					// Menambahkan menu user dan admin
					decorative.PrintLine()
					decorative.PrintTitle(" Admin Approval Menu ")
					decorative.PrintDecorativeLine()
					decorative.PrintMenu(1, "Approve user")
					decorative.PrintMenu(2, "Logout")
					decorative.PrintDecorativeLine()
					decorative.PrintInstruction(" Choose the number of the menu to continue ")
					decorative.PrintBottomLine()

					navigateInputIndex(2, choiceIndex)

				},
				ChoiceList: [4]entity.Choice{
					{
						ChoiceText: util.ADMIN_APPROVE_MENU,
						ChoiceFunc: func(userTypeIndex *int, routeIndex *int, choiceIndex *int) {
							decorative.HeaderAdminMenu(CurrentLogged)
							decorative.HeaderPage[int]("Approve User Page")

							decorative.PrintInfo(" List Unverified User")

							id := 0
							for *userTypeIndex == 0 && *routeIndex == 1 && *choiceIndex == 0 {
								authentication.RetrieveUnverifiedUser(USERS)
								fmt.Println("Enter user id to approve: ")
								fmt.Scan(&id)

								util.CheckForExitInput[int](id, func() {
									navigateRoute(util.ADMIN_APPROVAL_MENU, userTypeIndex, routeIndex, choiceIndex)
								})

								err, message := authentication.VerifyUser(id, &USERS)

								if err {
									decorative.PrintAlert(message)
								} else {
									decorative.PrintInfo(message)
								}

							}
						},
					},
					{
						ChoiceText: util.LOGOUT,
						ChoiceFunc: func(userTypeIndex *int, routeIndex *int, choiceIndex *int) {
							authentication.LogoutUser(&CurrentLogged)
							decorative.ResetPrintStatus(&printStatus, &printText)
							navigateRoute(util.ADMIN_AUTH_MENU, userTypeIndex, routeIndex, choiceIndex)
						},
					},
				},
			},
		},
	}
	routes[1] = entity.UserType{
		UserType: util.USER_TYPE,
		RouteList: [10]entity.Route{
			{
				RouteName: util.USER_AUTH_MENU,
				RouteFunc: func(printStatus *string, printText *string, hoiceIndex *int) {
					decorative.HeaderTemplate()
					// Menambahkan menu user dan admin
					if *printStatus == util.PRINT_STATUS_ERROR {
						decorative.PrintStatus(*printStatus, *printText)
					} else if *printStatus == util.PRINT_STATUS_NOTHING {
						decorative.PrintNothing()
					}
					decorative.PrintLine()
					decorative.PrintTitle(" User Auth Menu ")
					decorative.PrintDecorativeLine()
					decorative.PrintMenu(1, "Register")
					decorative.PrintMenu(2, "Login")
					decorative.PrintMenu(3, "Back")
					decorative.PrintDecorativeLine()
					decorative.PrintInstruction(" Choose the number of the menu to continue ")
					decorative.PrintBottomLine()

					navigateInputIndex(3, &choiceIndex)
				},
				ChoiceList: [4]entity.Choice{
					{
						ChoiceText: util.USER_AUTH_REGISTER_MENU,
						ChoiceFunc: func(userTypeIndex *int, routeIndex *int, choiceIndex *int) {
							decorative.HeaderTemplate()
							decorative.HeaderPage[string]("User Register Page")

							finishRegister := false

							for !finishRegister {
								name, email, password := authentication.InputUserRegister(func() {
									decorative.ResetPrintStatus(&printStatus, &printText)
									navigateRoute(util.USER_AUTH_MENU, userTypeIndex, routeIndex, choiceIndex)
								})
								err, status := authentication.RegisterUser(name, email, password, &USERS)

								if err {
									printStatus = util.PRINT_STATUS_ERROR
									printText = status
									navigateRoute(util.USER_AUTH_MENU, userTypeIndex, routeIndex, choiceIndex)
								} else {
									finishRegister = true
									printStatus = util.PRINT_STATUS_SUCCESS
									printText = status
								}

								decorative.PrintStatus(printStatus, printText)
							}

							decorative.PrintText("Press any key and enter to navigate to login page: ")
							var key int
							fmt.Scan(&key)

							navigateRoute(util.USER_AUTH_LOGIN_MENU, userTypeIndex, routeIndex, choiceIndex)
						},
					},
					{
						ChoiceText: util.USER_AUTH_LOGIN_MENU,
						ChoiceFunc: func(userTypeIndex *int, routeIndex *int, choiceIndex *int) {
							isLoggedIn := false

							for !isLoggedIn {
								decorative.HeaderTemplate()
								decorative.HeaderPage[string]("User Login Page")

								if printStatus == util.PRINT_STATUS_ERROR {
									decorative.PrintStatus(printStatus, printText)
								} else {
									decorative.PrintNothing()
								}

								email, password := authentication.InputUserLogin(func() {
									decorative.ResetPrintStatus(&printStatus, &printText)
									navigateRoute(util.USER_AUTH_MENU, userTypeIndex, routeIndex, choiceIndex)
								})
								err, status := authentication.LoginUser(email, password, USERS, &CurrentLogged)

								printText = status

								if err {
									util.ClearScreen()
									printStatus = util.PRINT_STATUS_ERROR
								} else {
									isLoggedIn = true
									printStatus = util.PRINT_STATUS_SUCCESS
									decorative.ResetPrintStatus(&printStatus, &printText)
									navigateRoute(util.USER_SUB_MENU, userTypeIndex, routeIndex, choiceIndex)
								}
							}
						},
					},
					{
						ChoiceText: util.USER_AUTH_BACK_MENU,
						ChoiceFunc: func(userTypeIndex *int, routeIndex *int, choiceIndex *int) {
							decorative.ResetPrintStatus(&printStatus, &printText)
							navigateRoute(util.ADMIN_TYPE, userTypeIndex, routeIndex, choiceIndex)
						},
					},
				},
			},
			{
				RouteName: util.USER_SUB_MENU,
				RouteFunc: func(printStatus *string, printText *string, choiceIndex *int) {
					decorative.HeaderUserMenu(CurrentLogged)
					// Menambahkan menu user dan admin
					if *printStatus == util.PRINT_STATUS_SUCCESS {
						decorative.PrintStatus(*printStatus, *printText)
					}
					decorative.PrintLine()
					decorative.PrintTitle(" User Menu ")
					decorative.PrintDecorativeLine()
					decorative.PrintMenu(1, "Send Email")
					decorative.PrintMenu(2, "Inbox")
					decorative.PrintMenu(3, "Logout")
					decorative.PrintDecorativeLine()
					decorative.PrintInstruction(" Choose the number of the menu to continue ")
					decorative.PrintBottomLine()

					navigateInputIndex(3, choiceIndex)

					*choiceIndex -= 1
					// Pilihan input nomor dari user dikurang 1 dan kita memanggil Menu untuk mengubah dan memanggil fungsi yang ada pada struct Choice
					util.ClearScreen()
					Menu()
				},
				ChoiceList: [4]entity.Choice{
					{
						ChoiceText: util.USER_SUB_MENU_SEND_EMAIL,
						ChoiceFunc: func(userTypeIndex *int, routeIndex *int, choiceIndex *int) {

							sent := false
							for !sent {
								decorative.HeaderUserMenu(CurrentLogged)
								decorative.HeaderPage[string]("Send Email Page")

								if printStatus == util.PRINT_STATUS_SUCCESS || printStatus == util.PRINT_STATUS_ERROR {
									decorative.PrintStatus(printStatus, printText)
								} else {
									decorative.PrintNothing()
								}

								decorative.ResetPrintStatus(&printStatus, &printText)
								to, subject, body := emails.WriteEmail(replyEmail, &CurrentLogged, func() {
									replyEmail = ""
									decorative.ResetPrintStatus(&printStatus, &printText)
									navigateRoute(util.USER_SUB_MENU, userTypeIndex, routeIndex, choiceIndex)
								})
								err, message := emails.SendEmail(CurrentLogged.Email, to, subject, body, &EMAILS)

								if err {
									util.ClearScreen()
									printStatus = util.PRINT_STATUS_ERROR
									printText = message
								} else {
									printStatus = util.PRINT_STATUS_SUCCESS
									printText = message
									sent = true
								}

								decorative.PrintStatus(printStatus, printText)
							}

							navigateRoute(util.USER_SUB_MENU_SEND_EMAIL, userTypeIndex, routeIndex, choiceIndex)
						},
					},
					{
						ChoiceText: util.USER_SUB_MENU_INBOX,
						ChoiceFunc: func(userTypeIndex *int, routeIndex *int, choiceIndex *int) {
							decorative.ResetPrintStatus(&printStatus, &printText)
							decorative.HeaderUserMenu(CurrentLogged)
							decorative.HeaderPage[int]("Inbox/Outbox Page")

							mail := emails.RetrieveEmails(EMAILS, CurrentLogged.Email)
							totalIdx := emails.ShowEmailListInbox(mail)

							idx := 0
							decorative.PrintInfo(" Input email number: ")
							inputsMenus(totalIdx, &idx)
							util.CheckForExitInput[int](idx, func() {
								navigateRoute(util.USER_SUB_MENU, userTypeIndex, routeIndex, choiceIndex)
							})

							selectedEmailIdx = mail[idx-1]

							*choiceIndex = 3
							Menu()
						},
					},
					{
						ChoiceText: util.LOGOUT,
						ChoiceFunc: func(userTypeIndex *int, routeIndex *int, choiceIndex *int) {
							authentication.LogoutUser(&CurrentLogged)
							navigateRoute(util.USER_AUTH_MENU, userTypeIndex, routeIndex, choiceIndex)
						},
					},
					{
						ChoiceText: util.USER_SUB_MENU_EMAIL_LIST,
						ChoiceFunc: func(userTypeIndex *int, routeIndex *int, choiceIndex *int) {
							decorative.HeaderUserMenu(CurrentLogged)
							decorative.HeaderPage[int]("Email Page")
							decorative.PrintStatus(printStatus, printText)
							decorative.ResetPrintStatus(&printStatus, &printText)

							emails.ReadEmail(selectedEmailIdx.From, selectedEmailIdx.To, &EMAILS, CurrentLogged)

							// fmt.Println("Selected Email:  ", selectedEmailIdx) // debug,.. need deleted
							list := emails.EmailList(selectedEmailIdx.From, selectedEmailIdx.To, EMAILS)
							totalIdx := emails.ShowEmailList(list)

							for *userTypeIndex == 1 && *routeIndex == 1 && *choiceIndex == 3 {
								decorative.PrintInfo(" Input action number: ")
								decorative.PrintInfo(" 1. Reply")
								decorative.PrintInfo(" 2. Delete")

								var key int
								inputsMenus(2, &key)
								util.CheckForExitInput[int](key, func() {
									navigateRoute(util.USER_SUB_MENU_INBOX, userTypeIndex, routeIndex, choiceIndex)
								})

								if key == 1 {
									if CurrentLogged.Email == selectedEmailIdx.From {
										replyEmail = selectedEmailIdx.To
									} else if CurrentLogged.Email == selectedEmailIdx.To {
										replyEmail = selectedEmailIdx.From
									}
									navigateRoute(util.USER_SUB_MENU_SEND_EMAIL, userTypeIndex, routeIndex, choiceIndex)
								} else if key == 2 {
									decorative.PrintInfo(" Input email number for delete: ")
									var idx int
									inputsMenus(totalIdx, &idx)
									util.CheckForExitInput[int](idx, func() {
										navigateRoute(util.USER_SUB_MENU_EMAIL_LIST, userTypeIndex, routeIndex, choiceIndex)
									})

									errDelete, message := emails.DeleteEmail(&EMAILS, list[idx-1].Id)

									if errDelete {
										printStatus = util.PRINT_STATUS_ERROR
										printText = message
									} else {
										printStatus = util.PRINT_STATUS_SUCCESS
										printText = message
									}

									navigateRoute(util.USER_SUB_MENU_EMAIL_LIST, userTypeIndex, routeIndex, choiceIndex)
								}
							}
						},
					},
				},
			},
		},
	}
}

func navigateRoute(name string, userTypeIndex *int, routeIndex *int, choiceIndex *int) {
	if name == util.ADMIN_TYPE || name == util.USER_TYPE {
		if name == util.ADMIN_TYPE {
			*userTypeIndex = -1
			*routeIndex = -1
			*choiceIndex = -1
			Menu()
		}
		if name == util.USER_TYPE {
			*userTypeIndex = -1
			*routeIndex = -0
			*choiceIndex = -1
			Menu()
		}
	} else {
		menuConfig := util.GetMenuConfig(name)
		*userTypeIndex = menuConfig.UserTypeIdx
		*routeIndex = menuConfig.RouteIdx
		*choiceIndex = menuConfig.ChoiceIdx
		Menu()
	}
}

func Menu() {
	util.ClearScreen()
	if userTypeIndex != -1 && routeIndex != -1 {
		if choiceIndex != -1 {
			/* Jika user sudah memilih nomor input, maka akan menampilkan fungsi yang ada pada
			struct Choice yaitu berupa ChoiceFunc*/
			routes[userTypeIndex].RouteList[routeIndex].ChoiceList[choiceIndex].ChoiceFunc(&userTypeIndex, &routeIndex, &choiceIndex)
		} else {
			/* Jika user belum memilih nomor input, maka akan menampilkan fungsi yang ada pada
			struct Route yaitu RouteFunc untuk memberikan pilihan kepada user*/
			routes[userTypeIndex].RouteList[routeIndex].RouteFunc(&printStatus, &printText, &choiceIndex)
		}
	} else {
		PrintStartMenu(&userTypeIndex, &routeIndex)
	}
}

func PrintStartMenu(userTypeIndex *int, routeIndex *int) {

	decorative.HeaderTemplate()

	// Menambahkan menu user dan admin
	decorative.PrintLine()
	decorative.PrintTitle(" Primary Menu ")
	decorative.PrintDecorativeLine()
	decorative.PrintMenu(1, "Admin Menu")
	decorative.PrintMenu(2, "User Menu")
	decorative.PrintDecorativeLine()
	decorative.PrintInstruction(" Choose the number of the menu to continue ")
	decorative.PrintBottomLine()

	*routeIndex = 0
	navigateInputIndex(2, userTypeIndex)

	// Menambahkan pesan penutup dengan warna yang berbeda
	// color.New(color.FgHiYellow, color.Bold).Println("\n🌟 Thanks for using this app! 🌟")
}

func SimpleHeaderTemplate() {
	decorative.PrintLine()
	decorative.PrintSubtitle(" Welcome to Email App")
	decorative.PrintBottomLine()
}

func HeaderAuthMenu(userType int) {
	SimpleHeaderTemplate()

	decorative.PrintLine()
	if userType == 0 {
		decorative.PrintTitle(" Admin Auth Menu ")
	} else {
		decorative.PrintTitle(" User Auth Menu ")
	}
	decorative.PrintBottomLine()

}

func navigateInputIndex(menuMax int, inputIndex *int) {
	valid := false

	for !valid {
		fmt.Print("Enter input number: ")
		fmt.Scan(inputIndex)
		status := validateChoiceInput(menuMax, inputIndex)

		if status {
			valid = true
		}
	}
	fmt.Println(*inputIndex)
	*inputIndex -= 1
	Menu()
}
func validateChoiceInput(max int, input *int) (status bool) {
	if *input >= 1 && *input <= max {
		return true
	}

	return false
}
func inputsMenus(menuMax int, input *int) {
	fmt.Scan(input)

	check := true
	for check {
		err := validateMenuInputs(menuMax, input)

		if err {
			fmt.Println("Enter correct input number: ")
			fmt.Scan(input)
		} else {
			check = false
		}
	}
}
func validateMenuInputs(max int, input *int) (err bool) {
	if *input == -1 {
		return false
	}

	if *input < 1 || *input > max {
		return true
	}

	return false
}
func seedData() {
	// Read data from JSON file
	path, _ := os.Getwd()
	fmt.Println(path)
	file, err := os.Open(path + "/src/seed/data.json")
	if err != nil {
		fmt.Println("Error opening data file:", err)
		return
	}
	defer file.Close()

	// Decode JSON data into structs
	var data struct {
		Users  []entity.User      `json:"users"`
		Admins []entity.UserAdmin `json:"admins"`
		Emails []entity.Email     `json:"emails"`
	}
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON data:", err)
		return
	}

	// Populate USERS, ADMINS, and EMAILS with data from JSON
	copy(USERS[:], data.Users)
	copy(ADMINS[:], data.Admins)
	copy(EMAILS[:], data.Emails)
}
