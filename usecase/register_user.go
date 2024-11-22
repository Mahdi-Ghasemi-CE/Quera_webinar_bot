package usecase

import (
	"Quera_webinar_bot/config"
	"Quera_webinar_bot/internal/enum"
	"Quera_webinar_bot/internal/models"
	"Quera_webinar_bot/internal/persistence/database"
	"Quera_webinar_bot/internal/telegram"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

type BotMainController struct {
	cfg      *config.Config
	localBot *tgbotapi.BotAPI
	db       *gorm.DB
}

func NewBotMainController(cfg *config.Config, localBot *tgbotapi.BotAPI) *BotMainController {
	return &BotMainController{cfg: cfg, db: database.GetDb(), localBot: localBot}
}

func (bmc BotMainController) SetupBot() {
	// Start receiving updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bmc.localBot.GetUpdatesChan(u)

	// State map to track user states
	var userStates = make(map[int64]string)

	for update := range updates {
		if update.Message != nil {
			chatID := update.Message.Chat.ID

			// Check the user's state and handle accordingly
			switch userStates[chatID] {
			case string(enum.AwaitingPassword):
				fmt.Println("1")
				// Handle login password
				if bmc.loginAdmin(update) {
					fmt.Println("2")
					bmc.sendReport(update)
					userStates[chatID] = string(enum.AdminLoggedIn)
				} else {
					fmt.Println("3")
					// Keep the user in "awaiting_password" state for retry
					userStates[chatID] = string(enum.AwaitingPassword)
				}

			case string(enum.AdminLoggedIn):
				fmt.Println("4")
				// Handle admin-specific commands after successful login
				if update.Message.Text == string(enum.AdminReport) {
					fmt.Println("5")
					bmc.sendReport(update)
				}

			default:
				// Handle general commands
				if update.Message.Text == string(enum.Start) {
					bmc.start(update)
				}
				if update.Message.Text == string(enum.Help) {
					bmc.help(update)
				}
				if update.Message.Contact != nil {
					bmc.saveContact(update)
				}
				if update.Message.Text == string(enum.AdminReport) {
					fmt.Println("6")

					// Start the login process
					userStates[chatID] = string(enum.AwaitingPassword)
					bmc.localBot.Send(tgbotapi.NewMessage(chatID, "لطفاً رمز عبور ادمین را وارد نمایید."))
					fmt.Println(userStates[chatID])
					fmt.Println("7")
				}
			}
		}
	}
}

func (bmc BotMainController) sendReport(update tgbotapi.Update) {
	fmt.Println("0")
	users, _ := getAllUsers(bmc.db)
	fmt.Println(users)
	// Create an Excel file
	f := excelize.NewFile()

	// Add headers to the first row
	headers := []string{"ID", "First Name", "Last Name", "Phone number"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string('A'+i))
		f.SetCellValue("Sheet1", cell, header)
	}
	fmt.Println("00")

	// Populate data rows
	rowNum := 2
	for i, user := range users {
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", rowNum), i)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", rowNum), user.FirstName)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", rowNum), user.LastName)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", rowNum), user.PhoneNumber)
		rowNum++
	}

	fmt.Println("000")

	// Save the Excel file to a temporary file
	tempFile := "user_report.xlsx"
	if err := f.SaveAs(tempFile); err != nil {
		log.Println("Error saving Excel file:", err)
		return
	}

	// Send the Excel file as a Telegram document
	file := tgbotapi.NewDocument(update.Message.Chat.ID, tgbotapi.FilePath(tempFile))
	if _, err := bmc.localBot.Send(file); err != nil {
		log.Println("Error sending file:", err)
		return
	}
	fmt.Println("0000")

	// Optionally, you can remove the temporary file after sending it
	if err := os.Remove(tempFile); err != nil {
		log.Println("Error deleting temporary file:", err)
	}
}

func (bmc BotMainController) saveContact(update tgbotapi.Update) {
	// Handle contact sharing
	phoneNumber := update.Message.Contact.PhoneNumber

	// Save the user to the database
	user := models.User{
		TelegramID:  update.Message.Chat.ID,
		PhoneNumber: phoneNumber,
		CreatedAt:   time.Now(),
		FirstName:   update.Message.Contact.FirstName,
		LastName:    update.Message.Contact.LastName,
	}

	existedUser, err := findUserByPhoneNumber(bmc.db, phoneNumber)
	if existedUser != nil {
		msg := "شماره تلفن شما قبلاً ثبت شده  \n ممنونیم که همچنان همراهمون هستی 😉	"
		bmc.removeKeyboard(msg, update.Message.Chat.ID)
	} else {
		if err = bmc.db.Create(&user).Error; err != nil {
			log.Printf("Failed to save user: %v", err)
			msg := "خطا در ذخیره اطلاعات لطفاً مجدد امتحان کن دوست من."
			bmc.removeKeyboard(msg, update.Message.Chat.ID)
		} else {
			msg := " ✅ ثبت نام شما با موفقیت انجام شد.\n\n به امید دیدار 👋"
			bmc.removeKeyboard(msg, update.Message.Chat.ID)
		}
	}
}

func (bmc BotMainController) removeKeyboard(msg string, chatID int64) {
	// Remove the keyboard
	removeKeyboard := tgbotapi.NewMessage(chatID, msg)
	removeKeyboard.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	bmc.localBot.Send(removeKeyboard)
}

func (bmc BotMainController) start(update tgbotapi.Update) {
	messageText := fmt.Sprintf("سلام رفیق به ربات کوئرا وبینار خوش اومدی 👋 \n\n" +
		"اگر دوست داری ثبت نام کنی لطفاً شماره تماست رو با ما به اشتراک بذار و روی دکمه ی 🔸اشتراک گذاری شماره همراه🔸 که این پایین 👇 برات نمایش داده میشه بزن و بعدش تایید کن که شماره ت برای ما ارسال بشه.")
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact("🔸اشتراک گذاری شماره همراه🔸"),
		),
	)
	_, err := bmc.localBot.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

func (bmc BotMainController) help(update tgbotapi.Update) {
	msg := fmt.Sprintf("")
	err := telegram.SendMessage(update.Message.Chat.ID, msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

func (bmc BotMainController) loginAdmin(update tgbotapi.Update) bool {
	// Predefined admin credentials
	const adminPassword = "admin@3041" // Replace with a strong, secure password

	// Validate that the update contains a message
	if update.Message == nil || update.Message.Chat == nil {
		log.Println("Invalid update: missing message or chat")
		return false
	}

	// Step 1: Prompt the user if no password is provided
	if update.Message.Text == "" {
		bmc.localBot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "لطفاً رمز عبور ادمین را وارد کنید :"))
		return false
	}

	// Step 2: Authenticate the password
	if update.Message.Text == adminPassword {
		// Step 3: Notify the admin of successful login
		successMessage := "✅ خوش اومدی ادمین گل الان گزارش ها رو برات میفرستم "
		bmc.localBot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, successMessage))

		return true
	}

	// Step 4: Notify of failed login attempt
	failMessage := "❌ رمز عبور غلطه ، اگر درستش رو میدونی بازم تلاش کن."
	bmc.localBot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, failMessage))

	return false
}

func findUserByPhoneNumber(db *gorm.DB, phoneNumber string) (*models.User, error) {
	var user models.User
	err := db.Where("phone_number = ?", phoneNumber).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No user found
		}
		return nil, err // Other errors
	}
	if user.ID == 0 {
		return nil, nil
	}
	return &user, nil
}

func getAllUsers(db *gorm.DB) ([]*models.User, error) {
	var users []models.User
	err := db.Find(&users).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No user found
		}
		return nil, err // Other errors
	}
	// Convert users to a slice of pointers
	var userPointers []*models.User
	for i := range users {
		userPointers = append(userPointers, &users[i])
	}
	return userPointers, nil
}
