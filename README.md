# User Service
The User Service is responsible for managing users, authentication, and user
permissions.

## Database Design
**`users` Table**
| Column Name | Type | Description |
|-------------|------|-------------|
| `id`        | unsigned bigint | GORM ID Column. |
| `created_at` | datetime | GORM Created At Column. |
| `updated_at` | datetime | GORM Updated At Column. |
| `deleted_at` | datetime | GORM Deleted At Column. |
| `Username`  | string | |
| `Password`  | string | |
| `Email`  | string | |
| `EmailVerifiedAt`  | datetime | |

**`roles` Table**
| Column Name | Type | Description |
|-------------|------|-------------|
| `user_id`  | unsigned bigint | |
| `role_id`  | unsigned bigint | |

**`user_roles` Table**
| Column Name | Type | Description |
|-------------|------|-------------|
| `id`        | unsigned bigint | GORM ID Column. |
| `created_at` | datetime | GORM Created At Column. |
| `updated_at` | datetime | GORM Updated At Column. |
| `deleted_at` | datetime | GORM Deleted At Column. |
| `Username`  | string | |
| `Password`  | string | |
| `Email`  | string | |

**`password_resets` Table**
| Column Name | Type | Description |
|-------------|------|-------------|
| `user_id`  | unsigned bigint | |
| `token`  | string | |
| `created_at`  | datetime | |
| `valid_until`  | datetime | |
