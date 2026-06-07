package user

import (
	"context"
	"errors"
	"fmt"

	"todoshnik/internal/auth"
	apperrors "todoshnik/internal/errors"

	"todoshnik/internal/validation"

	"gorm.io/gorm"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Add(ctx context.Context, name string, login string, password string) (*User, error) {
	user, _ := s.repo.GetByLogin(ctx, login)
	if user != nil {
		return nil, apperrors.ErrConflict
	}

	passwordHash := auth.HashPassword(password)

	newUser := &User{
		Name:         name,
		Login:        login,
		PasswordHash: passwordHash,
	}

	ve := validateUser(newUser)
	if ve != nil {
		return nil, ve
	}

	user, err := s.repo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) AddFromTg(ctx context.Context, name string, telegramID int64) (*User, error) {
	user, _ := s.repo.GetByTgId(ctx, telegramID)
	if user != nil {
		return user, nil
	}

	newUser := &User{
		Name:       name,
		TelegramID: telegramID,
	}

	ve := validateUser(newUser)
	if ve != nil {
		return nil, ve
	}

	user, err := s.repo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) List(ctx context.Context) ([]*User, error) {
	return s.repo.List(ctx)
}

func (s *Service) Update(ctx context.Context, userID int, name string) (*User, error) {
	user, errNotFound := s.Get(ctx, userID, "")
	if errNotFound != nil {
		return nil, apperrors.ErrNotFound
	}

	updated := *user
	updated.Name = name

	if err := validateUser(&updated); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, &updated); err != nil {
		return nil, err
	}

	return &updated, nil
}

func (s *Service) Delete(ctx context.Context, userID int) error {
	user, err := s.Get(ctx, userID, "")
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, user)
}

// нарушение принципа единой ответственности, как будет вести себя функция если передать и юзера и логин? Это где-то
// зафиксирвоано? Если не передать ни то ни другое?
// Я бы предпочел две раздельные функции GetByID и GetByLogin явное лучше чем неявное, непонятно что мы пытаемся
// сэкономить? Бумагу на которой будем распечатывать исходники?
// Хотя для сервисного слоя это вожможно и допустимо, благо на уровне репозиториев два разных метода.
// Может быть стоит создать тип AuthRequest{id, login} и внутри показать что может быть заполнено одно из полей или ни
// одно.
func (s *Service) Get(ctx context.Context, userID int, login string) (*User, error) {
	var user *User
	var err error

	// TODO изучи библиотеки https://github.com/samber/mo и https://github.com/samber/lo
	// в Го часто возникаю тисуации когда надо отличать zero value от отсутствия значнеиия. Например есть структура
	// с персональной инфо о пользователе. Есть поле кол-во детей. В структуре стоит значение 0. Мы не знаем это у
	// человека нет детей или есть дети, просто у нас нет об этом информации. Конечно надо стараться так чтобы
	// программа одинаково работала с zero value и с отсутствием значения, но так не всегда получается. Соответственно
	// можно использовать или указатели или монады/структуры вида Value[T]{v T, isDefined bool}
	if login != "" {
		// у тебя уже есть Service.GetByID, почему бы не переиспользовать эти два метода? В целом магический метод
		// выглядит избыточно. Я бы предпочел чтобы клиент сам выбирал есть у него логин или id и вызывал
		// соответствующую ручку. В общем случае по строке нельзя понять id это или логин чтобы достоверно угадать-
		// никто не мешает сделать логин 12345 или засунуть в логин числовой номер теелфона 892542342342, ты решишь что
		// это число, и будешь искать по id, а в итоге это логин. Явное лучше неявного
		user, err = s.repo.GetByLogin(ctx, login)
	} else if userID != 0 {
		// можно ли передать отрицательное число?
		user, err = s.repo.GetByID(ctx, userID)
	}

	if err != nil {
		// Изучи как errors.Is работает с nil, можно не писать проверку err на nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}

		return nil, err
	}

	// я бы предложил отличать доменные ошибки: отсутствие пользователя/ресурса от транспортных ошибок и приводить
	// одни к другим на апликейшн слое.
	if user == nil {
		return nil, apperrors.ErrNotFound
	}

	return user, nil
}

func (s *Service) GetByTgId(ctx context.Context, userTgID int64) (*User, error) {
	user, err := s.repo.GetByTgId(ctx, userTgID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}

		return nil, err
	}

	return user, nil
}

func validateUser(user *User) error {
	ve := validation.Validate(user)
	if ve != nil {
		fmt.Println(ve)
		return apperrors.NewValidationErrorFromValidator(ve)
	}
	return nil
}
