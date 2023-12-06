package postgres

const initUUIDextension = `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp"
`

const createTableUsers = `
	CREATE TABLE 
	IF NOT EXISTS users (
	id UUID DEFAULT uuid_generate_v4(),
	login VARCHAR(255) NOT NULL UNIQUE,
	password VARCHAR(255) NOT NULL,
	current DECIMAL(15,2) NOT NULL DEFAULT 0.00,
	withdrawn DECIMAL(15,2) NOT NULL DEFAULT 0.00,
	PRIMARY KEY (id))`

const createTableOrders = `
	CREATE TABLE
	IF NOT EXISTS orders (
	user_id UUID REFERENCES users(id),
	number VARCHAR(20) UNIQUE NOT NULL,
	accrual DECIMAL(15,2) NOT NULL DEFAULT 0.00,
	status VARCHAR(15) DEFAULT 'NEW',
	uploaded_at TIMESTAMPTZ DEFAULT current_timestamp,
	PRIMARY KEY (number))`

const createTableWithdrawals = `
	CREATE TABLE
	IF NOT EXISTS withdrawals (
	id UUID DEFAULT uuid_generate_v4(),
	user_id UUID REFERENCES users(id),
	order_id VARCHAR(20) UNIQUE NOT NULL,
	sum DECIMAL(15,2) NOT NULL DEFAULT 0.00,
	processed_at TIMESTAMPTZ DEFAULT current_timestamp,
	PRIMARY KEY (id))`
