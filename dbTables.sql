/* This also needs auth */
CREATE TABLE `users` (
	uid INTEGER PRIMARY KEY AUTOINCREMENT,
	uname TEXT UNIQUE NOT NULL
);

CREATE TABLE `links` (
	uid INTEGER,
	link BLOB,
	PRIMARY KEY (uid, link)
	FOREIGN KEY(uid) REFERENCES users(uid)
);
