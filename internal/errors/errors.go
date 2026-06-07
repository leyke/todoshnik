package errors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	// под стандартные статус коды есть стандартые ошибки, кажется в пакете http, кажется что никакой доп. нагрузки эти
	// ошибки не несут
	ErrNotFound        = errors.New("not found")
	ErrNotValidate     = errors.New("not validate")
	ErrConflict        = errors.New("conflict")
	ErrUnAuth          = errors.New("unauthorized")
	ErrEmptyBody       = errors.New("empty request body")
	ErrInvalidJSON     = errors.New("invalid json")
	ErrServerException = errors.New("server exception")
	ErrUnknownMethod   = errors.New("unknown method")
	// не смотрел еще, но кажется что эта доменная ошибка и ее лучше определить в пекете где она используется, это
	// слишком общий пакет чтобы тут это хранить.
	// тут я бы придерживался той же логики что и с константами, интерфейсами, метриками, моками... все это должно быть
	// объявлено по месту использования, тем более если это какая-то конкретная ошибка. Ошибки это такой же контракт
	// пакета как и объявления функций. Даже если 2 пакета возвращают "invalid task id" то скорее всего сценарии
	// того что это за task id и что с этим делать разные, в одно случае надо вернуть invalid argument, в другом
	// откатить финансовую транзакцию и стригерить алерт для ИБшников
	ErrInvalidTaskID = errors.New("invalid task id")
)

type ValidationError struct {
	error
	Errors map[string]string `json:"errors"`
}

// Странный метод павлик-морозов чот мне мешает воспользоваться стандратной error.Is тем что тут не будет поиска
// вглубину? Ну это скорее минус, чем плюс. Потом кто-то заврапает ошибку и все сломается. Это место где не нужно
// бороться с Го, привыкай это удобно.
func (e ValidationError) Is(target error) bool {
	return target == ErrNotValidate
}

func (e ValidationError) Error() string {
	if len(e.Errors) == 0 {
		return "Ошибка валидации"
	}

	var errs []string
	// поздравляю, ты придумал erros.Join, изучи стантартные возможности
	for field, message := range e.Errors {
		errs = append(errs, fmt.Sprintf("%s: %s", field, message))
	}
	return "Ошибка валидации: " + strings.Join(errs, "; ")
}

func NewValidationErrorFromValidator(err error) ValidationError {
	// Я глубоко не разбираляс, но это тоже выглядит странно. Зачем-то приведение типа и без проверки на ошибку–
	// очень смело
	ve := err.(validator.ValidationErrors)
	errors := make(map[string]string)

	// не знаю целенаправленно ты это делал или нет, но переменная цикла shadow аргумент функции, это может быть
	// неожиданно.
	// Какое значение будет у аргумента err после прохода цикла?
	for _, err := range ve {
		field := err.Field()
		tag := err.Tag()
		param := err.Param()

		// я бы постарался не аллоцировать message на каждой итерации, хотя я предполагаю что тут ботл-нека не будет
		// и массив ve короткий. Тем не менее лучше сразу аллоцировать errors нужного размера, это позволит избежать
		// переаллокаций map что довольно трудоемкая операция. И объявил бы промежуточную переменну вне цикла.
		// можно ли написать errors[field] = getValidatorMessage(field, tag, param) ?
		// может ли быть несколько ошибок в одном field?
		// можно ли обойтись массивом errors[]{field, message} и вставлять по индексу
		// errors[i] = {field, message}?
		message := getValidatorMessage(field, tag, param)
		errors[field] = message
	}

	return ValidationError{
		error:  ErrNotValidate,
		Errors: errors,
	}
}

// getValidatorMessage возвращает человеко-читаемое сообщение
func getValidatorMessage(field, tag, param string) string {
	switch tag {
	// Может быть стоит объвить свой тип и перегрузить метод String. Проблема таких селектов что надо не забыть
	// добавить в него правило и потом это превратится в портянку
	case "required":
		return fmt.Sprintf("%s обязательно", field)
	case "min":
		return fmt.Sprintf("%s минимум символов: %s", field, param)
	case "max":
		return fmt.Sprintf("%s максимум символов: %s", field, param)
	case "number":
		return fmt.Sprintf("%s должно быть числом", field)
	default:
		return fmt.Sprintf("%s неверно (%s)", field, tag)
	}
}
