# Database Migrations

This directory contains SQL migrations for the User Management API database.

## Migration Files

- `001_initial_schema.up.sql` - Initial database schema
- `001_initial_schema.down.sql` - Rollback initial schema

## Database Schema

### Tables

1. **users** - User accounts
   - `id` (BIGSERIAL) - Primary key
   - `username` (VARCHAR) - Unique username
   - `referrer_id` (BIGINT) - Reference to user who invited this user
   - `created_at` (TIMESTAMP) - Account creation time

2. **tasks** - Available tasks for users to complete
   - `id` (BIGSERIAL) - Primary key
   - `code` (VARCHAR) - Unique task code
   - `title` (VARCHAR) - Task title
   - `reward_points` (BIGINT) - Points awarded for completion
   - `is_active` (BOOLEAN) - Whether task is currently available
   - `created_at` (TIMESTAMP) - Task creation time

3. **user_tasks** - Completed tasks by users
   - `user_id` (BIGINT) - User who completed the task
   - `task_id` (BIGINT) - Completed task
   - `completed_at` (TIMESTAMP) - Completion time
   - Primary key: (user_id, task_id)

4. **balances** - User point balances
   - `user_id` (BIGINT) - Primary key, references users
   - `points` (BIGINT) - Current point balance
   - `updated_at` (TIMESTAMP) - Last update time

5. **transactions** - Transaction history (audit log)
   - `id` (BIGSERIAL) - Primary key
   - `user_id` (BIGINT) - User involved in transaction
   - `delta` (BIGINT) - Point change (positive or negative)
   - `reason` (VARCHAR) - Transaction reason
   - `reference_type` (VARCHAR) - Type of reference (e.g., "task", "referral")
   - `reference_id` (BIGINT) - ID of related entity
   - `created_at` (TIMESTAMP) - Transaction time

## Running Migrations

### Using psql directly:

Apply migration

```
psql -U postgres -d user_management -f migrations/001_initial_schema.up.sql
```

Rollback migration

```
psql -U postgres -d user_management -f migrations/001_initial_schema.down.sql
```


### Using Docker:

Apply migration

```
docker exec -i postgres_container psql -U postgres -d user_management < migrations/001_initial_schema.up.sql
```

Rollback migration
```
docker exec -i postgres_container psql -U postgres -d user_management < migrations/001_initial_schema.down.sql
```

## Features

- Automatic balance creation for new users via trigger
- Foreign key constraints for data integrity
- Indexes for optimized queries
- Sample tasks pre-populated for testing
