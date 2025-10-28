### План по API для Dvorfs Repository Manager

#### 1. Аутентификация (`/api/v1/auth`)
Эндпоинты для управления сессиями пользователей.

- `POST /login`
  - **Описание:** Аутентификация пользователя по логину и паролю.
  - **Ответ:** JWT (JSON Web Token) для доступа к защищенным ресурсам.
- `POST /logout`
  - **Описание:** Завершение сессии пользователя.
- `GET /me`
  - **Описание:** Получение информации о текущем аутентифицированном пользователе.

#### 2. Управление репозиториями (`/api/v1/repositories`)
CRUD-операции для репозиториев (`hosted`, `proxy`, `group`).

- `GET /`
  - **Описание:** Получить список всех репозиториев.
- `POST /`
  - **Описание:** Создать новый репозиторий.
  - **Тело запроса:** `{ "name": "npm-hosted", "format": "npm", "type": "hosted", "attributes": { ... } }`
- `GET /{name}`
  - **Описание:** Получить детальную информацию о репозитории по его имени.
- `PUT /{name}`
  - **Описание:** Обновить конфигурацию репозитория.
- `DELETE /{name}`
  - **Описание:** Удалить репозиторий.

#### 3. Работа с артефактами (`/repository/{repository-name}/...`)
Эти эндпоинты будут использоваться клиентами (Maven, NPM, Docker и т.д.) для загрузки и скачивания артефактов. Структура URL соответствует стандартам каждого формата.

- **Generic (RAW):**
  - `PUT /repository/raw-hosted/{path/to/artifact}`: Загрузить артефакт.
  - `GET /repository/raw-hosted/{path/to/artifact}`: Скачать артефакт.
- **Maven:**
  - `PUT /repository/maven-hosted/{groupId}/{artifactId}/{version}/{file}`: Загрузить Maven-артефакт.
  - `GET /repository/maven-proxy/{...}`: Скачать артефакт через прокси.
- **NPM:**
  - `PUT /repository/npm-hosted/{package}`: Опубликовать NPM-пакет.
  - `GET /repository/npm-group/{package}`: Скачать NPM-пакет через группу.
- **Docker:**
  - Эндпоинты будут соответствовать Docker Registry API v2 (например, `GET /v2/`, `GET /v2/{name}/manifests/{reference}`).

#### 4. Поиск (`/api/v1/search`)
Поиск артефактов по метаданным.

- `GET /artifacts`
  - **Описание:** Поиск артефактов.
  - **Параметры:** `?q=query`, `&repository=repo-name`, `&format=maven`.

#### 5. Управление пользователями и ролями (RBAC) (`/api/v1/security`)

- `GET /users`: Список пользователей.
- `POST /users`: Создать пользователя.
- `PUT /users/{username}`: Обновить пользователя.
- `PUT /users/{username}/password`: Сменить пароль пользователя.
- `DELETE /users/{username}`: Удалить пользователя.
- `GET /roles`: Список ролей.
- `POST /roles`: Создать роль с набором привилегий.
- `PUT /roles/{roleId}`: Обновить роль.
- `DELETE /roles/{roleId}`: Удалить роль.

#### 6. Политики очистки (`/api/v1/cleanup-policies`)
Управление правилами автоматического удаления старых артефактов.

- `GET /`: Получить список всех политик.
- `POST /`: Создать новую политику.
- `PUT /{policyId}`: Обновить политику.
- `DELETE /{policyId}`: Удалить политику.
