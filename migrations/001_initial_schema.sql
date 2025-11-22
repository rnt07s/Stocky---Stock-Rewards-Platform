-- Initial database schema for Stocky assignment

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255),
    email VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_user_id ON users(user_id);

-- Stock symbols reference table
CREATE TABLE IF NOT EXISTS stocks (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(20) UNIQUE NOT NULL,
    company_name VARCHAR(255) NOT NULL,
    exchange VARCHAR(10) NOT NULL DEFAULT 'NSE', -- NSE or BSE
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_stocks_symbol ON stocks(symbol);

-- Reward events - immutable log of all rewards given
CREATE TABLE IF NOT EXISTS reward_events (
    id SERIAL PRIMARY KEY,
    idempotency_key VARCHAR(255) UNIQUE NOT NULL, -- Prevents duplicate rewards
    user_id VARCHAR(100) NOT NULL,
    stock_symbol VARCHAR(20) NOT NULL,
    shares_quantity NUMERIC(18, 6) NOT NULL CHECK (shares_quantity > 0),
    price_per_share NUMERIC(18, 4) NOT NULL CHECK (price_per_share > 0),
    total_value NUMERIC(18, 4) NOT NULL, -- shares_quantity * price_per_share
    
    -- Fee breakdown (what Stocky pays)
    brokerage_fee NUMERIC(18, 4) NOT NULL DEFAULT 0,
    stt_fee NUMERIC(18, 4) NOT NULL DEFAULT 0, -- Securities Transaction Tax
    gst_fee NUMERIC(18, 4) NOT NULL DEFAULT 0, -- GST on brokerage
    exchange_fee NUMERIC(18, 4) NOT NULL DEFAULT 0,
    sebi_fee NUMERIC(18, 4) NOT NULL DEFAULT 0, -- SEBI turnover charges
    total_fees NUMERIC(18, 4) NOT NULL DEFAULT 0,
    
    total_cost NUMERIC(18, 4) NOT NULL, -- total_value + total_fees
    
    reason VARCHAR(255), -- 'onboarding', 'referral', 'milestone', etc.
    metadata JSONB, -- Additional context
    
    rewarded_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_reward_events_user_id ON reward_events(user_id);
CREATE INDEX idx_reward_events_stock_symbol ON reward_events(stock_symbol);
CREATE INDEX idx_reward_events_rewarded_at ON reward_events(rewarded_at);
CREATE INDEX idx_reward_events_user_rewarded ON reward_events(user_id, rewarded_at);
CREATE INDEX idx_reward_events_idempotency ON reward_events(idempotency_key);

-- User holdings - aggregated current holdings per user per stock
CREATE TABLE IF NOT EXISTS user_holdings (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(100) NOT NULL,
    stock_symbol VARCHAR(20) NOT NULL,
    total_shares NUMERIC(18, 6) NOT NULL DEFAULT 0 CHECK (total_shares >= 0),
    average_price NUMERIC(18, 4) NOT NULL DEFAULT 0, -- Weighted average purchase price
    last_updated TIMESTAMP NOT NULL DEFAULT NOW(),
    
    UNIQUE(user_id, stock_symbol)
);

CREATE INDEX idx_user_holdings_user_id ON user_holdings(user_id);
CREATE INDEX idx_user_holdings_stock_symbol ON user_holdings(stock_symbol);

-- Stock prices - hourly snapshots
CREATE TABLE IF NOT EXISTS stock_prices (
    id SERIAL PRIMARY KEY,
    stock_symbol VARCHAR(20) NOT NULL,
    price NUMERIC(18, 4) NOT NULL CHECK (price > 0),
    timestamp TIMESTAMP NOT NULL,
    source VARCHAR(50) NOT NULL DEFAULT 'mock', -- 'mock', 'nse', 'bse', etc.
    
    UNIQUE(stock_symbol, timestamp)
);

CREATE INDEX idx_stock_prices_symbol ON stock_prices(stock_symbol);
CREATE INDEX idx_stock_prices_timestamp ON stock_prices(timestamp);
CREATE INDEX idx_stock_prices_symbol_timestamp ON stock_prices(stock_symbol, timestamp DESC);

-- Ledger entries - double-entry bookkeeping
CREATE TABLE IF NOT EXISTS ledger_entries (
    id SERIAL PRIMARY KEY,
    entry_group_id UUID NOT NULL, -- Groups related debit/credit entries
    reward_event_id INTEGER REFERENCES reward_events(id),
    
    account_type VARCHAR(50) NOT NULL, -- 'stock_inventory', 'cash_outflow', 'fees_expense'
    stock_symbol VARCHAR(20), -- NULL for cash/fee accounts
    
    debit_amount NUMERIC(18, 4) NOT NULL DEFAULT 0,
    credit_amount NUMERIC(18, 4) NOT NULL DEFAULT 0,
    
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CHECK (
        (debit_amount > 0 AND credit_amount = 0) OR
        (credit_amount > 0 AND debit_amount = 0)
    )
);

CREATE INDEX idx_ledger_entries_group_id ON ledger_entries(entry_group_id);
CREATE INDEX idx_ledger_entries_reward_event ON ledger_entries(reward_event_id);
CREATE INDEX idx_ledger_entries_account_type ON ledger_entries(account_type);

-- Stock events - track corporate actions (splits, mergers, delisting)
CREATE TABLE IF NOT EXISTS stock_events (
    id SERIAL PRIMARY KEY,
    stock_symbol VARCHAR(20) NOT NULL,
    event_type VARCHAR(50) NOT NULL, -- 'split', 'merger', 'delisting', 'bonus'
    event_date DATE NOT NULL,
    
    -- For splits: old_ratio:new_ratio (e.g., "1:10" means 1 share becomes 10)
    split_ratio_old INTEGER,
    split_ratio_new INTEGER,
    
    -- For mergers
    merged_into_symbol VARCHAR(20),
    conversion_ratio NUMERIC(18, 6),
    
    description TEXT,
    processed BOOLEAN NOT NULL DEFAULT FALSE,
    processed_at TIMESTAMP,
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_stock_events_symbol ON stock_events(stock_symbol);
CREATE INDEX idx_stock_events_date ON stock_events(event_date);
CREATE INDEX idx_stock_events_processed ON stock_events(processed);

-- Insert some default Indian stocks
INSERT INTO stocks (symbol, company_name, exchange) VALUES
    ('RELIANCE', 'Reliance Industries Limited', 'NSE'),
    ('TCS', 'Tata Consultancy Services Limited', 'NSE'),
    ('INFY', 'Infosys Limited', 'NSE'),
    ('HDFCBANK', 'HDFC Bank Limited', 'NSE'),
    ('ICICIBANK', 'ICICI Bank Limited', 'NSE'),
    ('HINDUNILVR', 'Hindustan Unilever Limited', 'NSE'),
    ('ITC', 'ITC Limited', 'NSE'),
    ('BHARTIARTL', 'Bharti Airtel Limited', 'NSE'),
    ('KOTAKBANK', 'Kotak Mahindra Bank Limited', 'NSE'),
    ('WIPRO', 'Wipro Limited', 'NSE')
ON CONFLICT (symbol) DO NOTHING;
