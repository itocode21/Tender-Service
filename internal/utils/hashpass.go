package utils

/*
В ГОЛОВЕ ЭТО ДОЛЖНА БЫЛА БЫТЬ УТИЛКА ДЛЯ ХЕША ПАРОЛЯ ПРИ РЕГИСТРАЦИИ ПОЛЬЗОВАТЕЛЯ, НО НЕ ПРИГОДИТСЯ.
import "golang.org/x/crypto/bcrypt"

func HashPass(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "Error on hashPass func", err
	}

	return string(hashedPassword), nil
}

func CheckHash(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
*/
