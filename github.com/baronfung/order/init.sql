CREATE DATABASE IF NOT EXISTS orders;
USE orders;

CREATE TABLE IF NOT EXISTS Orders(
	ORDER_ID int NOT NULL AUTO_INCREMENT,
	START_LATITUDE VARCHAR(50),
	START_LONGTITUDE VARCHAR(50),
	END_LATITUDE VARCHAR(50),
	END_LONGTITUDE VARCHAR(50),
	TOTAL_DISTANCE VARCHAR(50),
	STATUS VARCHAR(50),
	ERROR_DESCRIPTION VARCHAR(255),
	PRIMARY KEY (ORDER_ID)
);

INSERT INTO Orders VALUES(
	NULL,
	22.277627,
	114.173463,
	22.2783034,
	114.1796477,
	1.1,
	"TAKEN",
	null
);
INSERT INTO Orders VALUES(
	NULL,
	22.2783034,
	114.1796477,
	22.2783034,
	114.1796477,
	1.1,
	"TAKEN",
	null
);