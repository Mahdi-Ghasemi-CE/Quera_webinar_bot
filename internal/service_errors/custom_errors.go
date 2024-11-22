package service_errors

import (
	"errors"
	"fmt"
)

func DuplicateDataError(dataName string) error {
	text := fmt.Sprintf("%s وارد شده تکراری می باشد ، لطفاً با %s جدید دوباره تلاش کنید. ", dataName, dataName)
	return errors.New(text)
}

func InvalidDataError() error {
	text := fmt.Sprintf("اطلاعات وارد شده نامعتبر است.")
	return errors.New(text)
}

func UsedOtpError() error {
	text := fmt.Sprintf("کد فعالسازی قبلاً استفاده شده است.")
	return errors.New(text)
}

func ValidOldOtpError() error {
	text := fmt.Sprintf("کد فعالسازی قبلی شما هنوز معتبر است.")
	return errors.New(text)
}

func InvalidCustomDataError(dataName string) error {
	text := fmt.Sprintf("%s وارد شده نامعتبر است.", dataName)
	return errors.New(text)
}

func DataNotFoundError() error {
	text := fmt.Sprintf("داده ای پیدا نشد.")
	return errors.New(text)
}

func InternalServerError() error {
	text := fmt.Sprintf("خطای سمت سرور رخ داد.")
	return errors.New(text)
}
