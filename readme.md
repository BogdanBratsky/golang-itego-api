# Документация API для блога

###### Основные возможности

***

## Эндпоинты для работы с пользователями

***

#### POST api/register
>Создаёт нового пользователя.

**METHOD**: POST
**URL**: api/register

**request:** 

>POST api/register

>Content-type: application/json
>Authorization: Bearer {token}

```
{
    "userName": "User_Name",
    "userEmail": "user@email.com",
    "userPassword": "ExampleOfPassword2048"
}
```

**response (201):**
>Content-type: application/json

```
{
    "success": true,
    "message": "Пользователь успешно создан",
    "user": {
        "userId": 1,
        "userName": "@User_Name",
        "userEmail": "user@email.com",
        "createdAt": "2024-02-25T12:00:00Z"
    }
}
```

**response (400):**
>Content-type: application/json

```
{
    "success": false,
    "message": "Некорректный запрос"
}
```
**response (401):**

>Content-type: application/json
```

{
    "success": false,
    "message": "Необходимо авторизоваться"
}
```
**response (403):**

>Content-type: application/json
```

{
    "success": false,
    "message": "Нет прав на это действие"
}
```
**response (404):**
>Content-type: application/json

```
{
    "success": false,
    "message": "Ресурс не найден"
}
```

**response (500):**
>Content-type: application/json

```
{
    "success": false,
    "message": "Внутренняя ошибка сервера"
}
```
---
#### POST api/login
>Эндпоинт для авторизации.

**METHOD**: POST
**URL**: api/login



**request:** 
```
POST api/login
```
>Content-type: application/json
```
{
    "userName": "User_Name",
    "userPassword": "ExampleOfPassword2048"
}
```

**response (200):**
>Content-type: application/json

```
{
    "success": true,
    "message": "Пользователь успешно авторизован",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
    "user": {
        "userId": 1,
        "userName": "@User_Name",
        "userEmail": "user@email.com",
        "createdAt": "2024-02-25T12:00:00Z"
    }
}
```

**response (400):**
>Content-type: application/json

```
{
    "success": false,
    "message": "Некорректный запрос"
}
```
**response (404):**
>Content-type: application/json

```
{
    "success": false,
    "message": "Ресурс не найден"
}
```
**response (500):**
>Content-type: application/json

```
{
    "success": false,
    "message": "Внутренняя ошибка сервера"
}
```
---

#### GET api/users
>Возвращает список пользователей. Если не указать в квери-параметрах 'page' и 'per_page', то в ответе будет возвращено 50 последних зарегистрировавшихся пользователей. В 'per_page' можно передать от 1 до 100. Если в этот параметр передано больше 100, то в ответе вернётся ошибка 404.

**METHOD**: GET
**URL**: api/users



**request:**
```
GET api/users
```
**or**
```
GET api/users?page=1&per_page=100
```
**response (200):**

>Content-type: application/json
```

{
    "success": true,
    "items": [
        {
            "userId": 1,
            "userName": "@User_Name",
            "createdAt": "2024-02-25T12:00:00Z",
        },
        ...
        // прочие пользователи
    ],
    "pagination": {
        "page": 1,
        "perPage": 100,
        "pagesCount": 50,
    },
    "totalCount": 5000
}
```

**response (400):**

>Content-type: application/json
```

{
    "success": false,
    "message": "Некорректный запрос"
}
```
**response (404):**

>Content-type: application/json
```

{
    "success": false,
    "message": "Ресурс не найден"
}
```
**response (500):**

>Content-type: application/json
```

{
    "success": false,
    "message": "Внутренняя ошибка сервера"
}
```

---
#### GET api/users/:id
>Возвращает конкретного пользователя по указанному в запросе идентификатору. 

**METHOD**: GET
**URL**: api/users/:id



**request:**
```
GET api/users/1
```
**response (200):**

>Content-type: application/json
```

{
    "success": true,
    "user": {
        "userId": 1,
        "userName": "@User_Name",
        "createdAt": "2024-02-25T12:00:00Z",
    }
}
```

**response (400):**

>Content-type: application/json
```

{
    "success": false,
    "message": "Некорректный запрос"
}
```
**response (404):**

>Content-type: application/json
```

{
    "success": false,
    "message": "Ресурс не найден"
}
```
**response (500):**

>Content-type: application/json
```

{
    "success": false,
    "message": "Внутренняя ошибка сервера"
}
```

---
