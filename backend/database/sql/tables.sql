-- Location table
CREATE TABLE IF NOT EXISTS Locations (
    ID   VARCHAR(36) PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    LegID VARCHAR(36) REFERENCES Legs(ID),
    FOREIGN KEY (LegID) REFERENCES Legs(ID)
);

-- Pricelist table
CREATE TABLE IF NOT EXISTS Pricelists (
    ID         VARCHAR(36) PRIMARY KEY,
    ValidUntil TIMESTAMP
);

-- Company table
CREATE TABLE IF NOT EXISTS Companies (
    ID   VARCHAR(36) PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    PriceListID VARCHAR(36) REFERENCES Pricelists(ID),
    FOREIGN KEY (PriceListID) REFERENCES Pricelists(ID)
);

-- RouteInfo table
CREATE TABLE IF NOT EXISTS RouteInfos (
    ID       VARCHAR(36) PRIMARY KEY,
    FromID   VARCHAR(36) REFERENCES Locations(ID),
    ToID     VARCHAR(36) REFERENCES Locations(ID),
    Distance INT,
    LegID    VARCHAR(36) REFERENCES Legs(ID),
    FOREIGN KEY (LegID) REFERENCES Legs(ID),
    FOREIGN KEY (FromID, ToID) REFERENCES Locations(ID, ID)
);

-- Provider table
CREATE TABLE IF NOT EXISTS Providers (
    ID          VARCHAR(36) PRIMARY KEY,
    CompanyID   VARCHAR(36) REFERENCES Companies(ID),
    Price       FLOAT,
    FlightStart TIMESTAMP,
    FlightEnd   TIMESTAMP,
    LegID       VARCHAR(36) REFERENCES Legs(ID), 
    FOREIGN KEY (CompanyID) REFERENCES Companies(ID),
    FOREIGN KEY (LegID) REFERENCES Legs(ID)
);

-- Leg table
CREATE TABLE IF NOT EXISTS Legs (
    ID            VARCHAR(36) PRIMARY KEY,
    RouteInfoID   VARCHAR(36) REFERENCES RouteInfos(ID),
    PriceListID   VARCHAR(36) REFERENCES Pricelists(ID),
    FOREIGN KEY (RouteInfoID) REFERENCES RouteInfos(ID),
    FOREIGN KEY (PriceListID) REFERENCES Pricelists(ID)
);

-- CachedRoutes table
CREATE TABLE IF NOT EXISTS CachedRoutes (
    ID              VARCHAR(36) PRIMARY KEY,
    PricelistID     VARCHAR(36) REFERENCES Pricelists(ID),
    FromLocation    VARCHAR(255) NOT NULL,
    ToLocation      VARCHAR(255) NOT NULL,
    ValidUntil      TIMESTAMP,
    Routes          TEXT
);

--Bookings table
CREATE TABLE IF NOT EXISTS Bookings (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    CompanyNames TEXT NOT NULL,
    StartTime TEXT NOT NULL,
    FirstName TEXT NOT NULL,
    LastName TEXT NOT NULL,
    TotalPrice REAL NOT NULL,
    TotalDuration TEXT NOT NULL,
    PricelistID INTEGER NOT NULL,
    FromCity TEXT NOT NULL,
    DestinationCity TEXT NOT NULL,
    ValidUntil TEXT
);