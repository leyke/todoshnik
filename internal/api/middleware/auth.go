package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	authapi "todoshnik/internal/auth/api"
	authcontext "todoshnik/internal/auth/context"
)

func Auth(ah *authapi.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, ok := extractTokenFromHeader(r)
			if !ok {
				http.Error(w, "Unauthorized: не передан токен", http.StatusUnauthorized)
				return
			}

			user, err := ah.GetAuthorizedUser(r.Context(), tokenString)
			if err != nil {
				// ты тут где-то логируешь, где-то не логируешь что немного смотрится неоднородно, но основная проблема
				// что написать юнит тест на такую функцию сложно: у нее много сайд эфектов: это и логирование (на что
				// можно подзабить и не тестировать), а еще вызов http.Error. Крч написать тест на такое тяжело,
				// нужно будет мокать.
				// 1. Нужно стараться писать чисые фукции и вытаскивать сайд-эфекты наружу
				// 2. Тут у меня нет сильно мнения, все таки это мидлваря, но даже тут можно выделить чистую функцию и
				//    возвращать результат и ошибку
				// 3. В целом функция которая возвращает замыкание внутри другого замыкания выглядит странновато, тем
				//    более я не вижу особо какого-то захвата переменных и скоупа
				fmt.Printf("Ошибка валидации токена: %v\n", err)
				http.Error(w, "Unauthorized: не найден пользователь", http.StatusUnauthorized)
				return
			}

			ctx := authcontext.SetUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractTokenFromHeader(r *http.Request) (string, bool) {
	// я бы советова выносить все строковые литералы в неимпортируемые константы/переменные
	authHeader := r.Header.Get("Authorization")

	// опять же про чистые функции и тестирование, можно все что ниже вынести в отдельную неимпортированную функцию
	// пакета, покрыть ее табличным тестом и будет хорошо и без лишних щависимостей и ты разделить логику извлечения
	// заголовка от санитации какой-то

	// - n > 0: at most n substrings; the last substring will be the unsplit remainder;
	// если будет передано bearer hello awesome world, то все равно вернется 2 части, это точно тебя устраивает?
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return "", false
	}

	// заголовки в http регистронезависимы, могут быть проблемы
	if !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", false
	}

	return token, true
}

// аргумент не используется, это для совместимости с интерфейсом?
func BotAuth(ah *authapi.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// аналогично– вынести заголовки в константы
			authHeader := r.Header.Get("X-Bot-Service-Token")

			token := strings.TrimSpace(authHeader)
			if token == "" {
				// вынести ошибки в переменные пакета
				// ErrUnauthorized = errors.New("Unauthorized: не передан токен: %w", http.StatusUnauthorized)
				// ...
				// return ErrUnauthorized
				// или
				// http.Error(w, ErrUnauthorized)
				http.Error(w, "Unauthorized: не передан токен", http.StatusUnauthorized)
				return
			}

			if token != os.Getenv("BOT_SERVICE_TOKEN") {
				fmt.Printf("Неверный сервисный токен бота")
				http.Error(w, "Unauthorized: сервисный токен не верен", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
