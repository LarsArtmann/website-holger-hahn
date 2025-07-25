-- Initial database schema for Holger Hahn website
-- Contact form submissions and basic analytics

-- Contact form submissions
CREATE TABLE contacts (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    company TEXT,
    message TEXT NOT NULL,
    subject TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    status TEXT DEFAULT 'new' CHECK (status IN ('new', 'read', 'replied', 'archived')),
    source TEXT DEFAULT 'website' -- 'website', 'api', 'direct'
);

-- Analytics for page views and interactions
CREATE TABLE analytics_events (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    event_type TEXT NOT NULL, -- 'page_view', 'contact_form_submit', 'service_click', etc.
    page_path TEXT,
    user_agent TEXT,
    ip_address TEXT,
    session_id TEXT,
    referrer TEXT,
    metadata TEXT, -- JSON for additional event data
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Technologies showcase (can be managed dynamically)
CREATE TABLE technologies (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    name TEXT NOT NULL UNIQUE,
    category TEXT NOT NULL, -- 'language', 'infrastructure', 'blockchain'
    icon_class TEXT, -- CSS class for devicon or custom icon
    color_scheme TEXT, -- 'blue', 'green', 'orange', etc.
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Services offered (can be managed dynamically)
CREATE TABLE services (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    features TEXT, -- JSON array of feature strings
    icon_svg TEXT, -- SVG markup for service icon
    color_scheme TEXT, -- 'blue', 'green', 'purple', etc.
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Professional experience (can be managed dynamically)
CREATE TABLE experiences (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    company TEXT NOT NULL,
    position TEXT NOT NULL,
    description TEXT NOT NULL,
    achievements TEXT, -- JSON array of achievement strings
    technologies TEXT, -- JSON array of technology tags
    start_date DATE NOT NULL,
    end_date DATE, -- NULL for current position
    is_current BOOLEAN DEFAULT FALSE,
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_contacts_created_at ON contacts(created_at);
CREATE INDEX idx_contacts_status ON contacts(status);
CREATE INDEX idx_contacts_email ON contacts(email);

CREATE INDEX idx_analytics_events_created_at ON analytics_events(created_at);
CREATE INDEX idx_analytics_events_type ON analytics_events(event_type);
CREATE INDEX idx_analytics_events_page ON analytics_events(page_path);

CREATE INDEX idx_technologies_category ON technologies(category);
CREATE INDEX idx_technologies_active ON technologies(is_active);
CREATE INDEX idx_technologies_sort ON technologies(sort_order);

CREATE INDEX idx_services_active ON services(is_active);
CREATE INDEX idx_services_sort ON services(sort_order);

CREATE INDEX idx_experiences_dates ON experiences(start_date, end_date);
CREATE INDEX idx_experiences_current ON experiences(is_current);
CREATE INDEX idx_experiences_active ON experiences(is_active);
CREATE INDEX idx_experiences_sort ON experiences(sort_order);