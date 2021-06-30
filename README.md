# User Service
The User Service is responsible for managing users, authentication, and user
permissions.

## Database Design
**`users` Table**
| Column Name | Type | Description |
|-------------|------|-------------|
| `Username`  | string | |
| `Password`  | string | |
| `Email`  | string | |

**`roles` Table**
| Column Name | Type | Description |
|-------------|------|-------------|
| `user_id`  | unsigned bigint | |
| `role_id`  | unsigned bigint | |

**`user_roles` Table**
| Column Name | Type | Description |
|-------------|------|-------------|
| `Username`  | string | |
| `Password`  | string | |
| `Email`  | string | |
