# access-refresh-tokens

Хочу подметить, что refresh token имеет структуру время_создания.Время_истечения.IDпользователя;
Для кодирования токена используется AES
В Dockerfile указаны переменные окружения SECRET_KEY для шифрования JWT (SHA512) и AES_KEY для refresh token