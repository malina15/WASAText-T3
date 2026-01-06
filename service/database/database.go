/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.
To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.
For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// Errors section
var ErrPhotoDoesntExist = errors.New("photo doesn't exist")
var ErrUserBanned = errors.New("user is banned")
var ErrNicknameAlreadyTaken = errors.New("nickname already taken")
var ErrGroupNotFound = errors.New("group not found")
var ErrUserAlreadyInGroup = errors.New("user already in group")
var ErrUserNotInGroup = errors.New("user not in group")
var ErrUserPhotoNotFound = errors.New("user photo not found")
var ErrMessageNotFound = errors.New("message not found")
var ErrForbiddenMessageAction = errors.New("forbidden message action")

/*
var ErrUserAutoLike = errors.New("users can't like their own photos")
var ErrUserAutoFollow = errors.New("users can't follow themselfes")
*/

// Constants
const PhotosPerUserHome = 3

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	// Creates a new user in the database. It returns an error
	CreateUser(User) error

	// Modifies the nickname of a user in the database. It returns an error
	ModifyNickname(User, Nickname) error

	// Searches all the users that match the given name (both identifier and nickname). Returns the list of matching users and an error
	SearchUser(searcher User, userToSearch User) ([]CompleteUser, error)

	// Creates a new photo in the database. It returns the photo identifier and an error
	CreatePhoto(Photo) (int64, error)

	// Inserts a like of a user for a specified photo in the database. It returns an error
	LikePhoto(PhotoId, User) error

	// Removes a like of a user for a specified photo from the database. It returns an error
	UnlikePhoto(PhotoId, User) error

	// Adds a comment from a user to a specified photo in the database. It returns the unique comment id and an error
	CommentPhoto(PhotoId, User, Comment) (int64, error)

	// Deletes a comment from a user from a specified photo in the database. It returns an error
	UncommentPhoto(PhotoId, User, CommentId) error

	// Adds a follower (a) to the user that is being followed (b). It returns an error
	FollowUser(a User, b User) error

	// Removes a follower (a) from the user that is being unfollowed (b). It returns an error
	UnfollowUser(a User, b User) error

	// Adds a user (b) to the banned list of another (a). It returns an error
	BanUser(a User, b User) error

	// Removes a user (b) from the banned list of another (a). It returns an error
	UnbanUser(a User, b User) error

	// Get the a user's stream (photos of people who are followed by the user in reversed chronological order). It returns the photos and an error
	GetStream(User) ([]Photo, error)

	// Removes a photo from the database. The removal includes likes and comments.  It returns an error
	RemovePhoto(User, PhotoId) error

	// ____________________________________  Util Methods ____________________________________

	// Gets the followers list for the specified user. Returns the followers list and an error
	GetFollowers(User) ([]User, error)

	// Gets the following list for the specified user. Returns the following list and an error
	GetFollowing(User) ([]User, error)

	// Gets all users
	GetAllUsers() ([]User, error)

	// Gets the photos list of user b for the user a. Returns the photo list and an error
	GetPhotosList(a User, b User) ([]Photo, error)

	// Allows the author of a photo to remove a comment from another user on his/her photo. Returns an error
	UncommentPhotoAuthor(PhotoId, CommentId) error

	// Gets the nickname of a user. Returns the nickname and an error
	GetNickname(User) (string, error)

	// Finds a user by nickname. Returns the user, whether it was found, and an error
	FindUserByNickname(nickname string) (User, bool, error)

	// Creates a group and returns its identifier
	CreateGroup(creator User, name string, members []User) (int64, error)

	// Adds an existing user to an existing group
	AddUserToGroup(groupId int64, user User) error

	// Removes a user from a group
	RemoveUserFromGroup(groupId int64, user User) error

	// Checks if a user is member of a group
	IsUserInGroup(groupId int64, user User) (bool, error)

	// Updates group name
	SetGroupName(groupId int64, name string) error

	// Updates group photo path (local server path)
	SetGroupPhotoPath(groupId int64, photoPath string) error

	// Retrieves group info
	GetGroup(groupId int64) (Group, error)

	// Lists groups of a user
	ListGroupsForUser(user User) ([]Group, error)

	// Sets user's profile photo path (local server path)
	SetUserPhotoPath(user User, photoPath string) error

	// Gets user's profile photo path (local server path)
	GetUserPhotoPath(user User) (string, error)

	// Lists conversations (direct peers and groups), sorted by last activity (reverse chronological)
	ListConversations(user User) ([]Conversation, error)

	// Checks if a user (a) is banned by another (b). Returns a boolean
	BannedUserCheck(a User, b User) (bool, error)

	// Checks if a user (a) exists
	CheckUser(a User) (bool, error)

	// Checks if a photo (via its id) exists. Returns an error
	CheckPhotoExistence(p PhotoId) (bool, error)

	// Ping checks whether the database is available or not (in that case, an error will be returned)
	Ping() error

    // Chat methods
    CreateMessage(from User, to User, body string) (int64, error)
    ListMessages(a User, b User, limit int, offset int) ([]Message, error)

	// Group chat methods
	CreateGroupMessage(groupId int64, from User, body string) (int64, error)
	ListGroupMessages(groupId int64, limit int, offset int) ([]GroupMessage, error)

	// Message operations (delete / reactions)
	GetDirectMessageInConversation(a User, b User, messageId int64) (Message, error)
	GetGroupMessageInGroup(groupId int64, messageId int64) (GroupMessage, error)
	DeleteDirectMessage(messageId int64, deletedBy User) error
	DeleteGroupMessage(groupId int64, messageId int64, deletedBy User) error

	SetDirectMessageReaction(messageId int64, user User, reaction string) error
	RemoveDirectMessageReaction(messageId int64, user User) error
	ListDirectMessageReactions(messageId int64) ([]MessageReaction, error)

	SetGroupMessageReaction(groupId int64, messageId int64, user User, reaction string) error
	RemoveGroupMessageReaction(groupId int64, messageId int64, user User) error
	ListGroupMessageReactions(groupId int64, messageId int64) ([]MessageReaction, error)

	// Read receipts (checkmarks)
	MarkDirectConversationRead(reader User, peer User) error
	MarkGroupConversationRead(groupId int64, reader User) error
	GetDirectMessageCheckmarks(messageId int64) (int, error)
	GetGroupMessageCheckmarks(groupId int64, messageId int64) (int, error)
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Activate foreign keys for db

	_, errPramga := db.Exec(`PRAGMA foreign_keys= ON`)
	if errPramga != nil {
		return nil, fmt.Errorf("error setting pragmas: %w", errPramga)
	}

	// Ensure all required tables exist (CREATE TABLE IF NOT EXISTS).
	// This also works as a minimal migration step when new tables are added.
	err := createDatabase(db)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

// Creates all the necessary sql tables for the WASAPhoto app.
func createDatabase(db *sql.DB) error {
    tables := [17]string{
		`CREATE TABLE IF NOT EXISTS users (
			id_user VARCHAR(16) NOT NULL PRIMARY KEY,
			nickname VARCHAR(16) NOT NULL
			);`,
		`CREATE TABLE IF NOT EXISTS user_photos (
			id_user VARCHAR(16) NOT NULL PRIMARY KEY,
			photo_path TEXT NOT NULL,
			FOREIGN KEY(id_user) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS groups (
			id_group INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			photo_path TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL
			);`,
		`CREATE TABLE IF NOT EXISTS group_members (
			id_group INTEGER NOT NULL,
			id_user VARCHAR(16) NOT NULL,
			PRIMARY KEY (id_group, id_user),
			FOREIGN KEY(id_group) REFERENCES groups (id_group) ON DELETE CASCADE,
			FOREIGN KEY(id_user) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS photos (
			id_photo INTEGER PRIMARY KEY AUTOINCREMENT,
			id_user VARCHAR(16) NOT NULL,
			date DATETIME NOT NULL,
			FOREIGN KEY(id_user) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS  likes (
			id_photo INTEGER NOT NULL,
			id_user VARCHAR(16) NOT NULL,
			PRIMARY KEY (id_photo,id_user),
			FOREIGN KEY(id_photo) REFERENCES photos (id_photo) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS comments (
			id_comment INTEGER PRIMARY KEY AUTOINCREMENT,
			id_photo INTEGER NOT NULL,
			id_user VARCHAR(16) NOT NULL,
			comment VARCHAR(30) NOT NULL,
			FOREIGN KEY(id_photo) REFERENCES photos (id_photo) ON DELETE CASCADE,
			FOREIGN KEY(id_user) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS banned_users (
			banner VARCHAR(16) NOT NULL,
			banned VARCHAR(16) NOT NULL,
			PRIMARY KEY (banner,banned),
			FOREIGN KEY(banner) REFERENCES users (id_user) ON DELETE CASCADE,
			FOREIGN KEY(banned) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
        `CREATE TABLE IF NOT EXISTS followers(
			follower VARCHAR(16) NOT NULL,
			followed VARCHAR(16) NOT NULL,
			PRIMARY KEY (follower,followed),
			FOREIGN KEY(follower) REFERENCES users (id_user) ON DELETE CASCADE,
			FOREIGN KEY(followed) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
        `CREATE TABLE IF NOT EXISTS messages (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            sender VARCHAR(16) NOT NULL,
            receiver VARCHAR(16) NOT NULL,
            body TEXT NOT NULL,
            date DATETIME NOT NULL,
            FOREIGN KEY(sender) REFERENCES users (id_user) ON DELETE CASCADE,
            FOREIGN KEY(receiver) REFERENCES users (id_user) ON DELETE CASCADE
            );`,
		`CREATE TABLE IF NOT EXISTS group_messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			id_group INTEGER NOT NULL,
			sender VARCHAR(16) NOT NULL,
			body TEXT NOT NULL,
			date DATETIME NOT NULL,
			FOREIGN KEY(id_group) REFERENCES groups (id_group) ON DELETE CASCADE,
			FOREIGN KEY(sender) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS direct_message_deletions (
			message_id INTEGER NOT NULL PRIMARY KEY,
			deleted_at DATETIME NOT NULL,
			deleted_by VARCHAR(16) NOT NULL,
			FOREIGN KEY(message_id) REFERENCES messages (id) ON DELETE CASCADE,
			FOREIGN KEY(deleted_by) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS group_message_deletions (
			message_id INTEGER NOT NULL PRIMARY KEY,
			id_group INTEGER NOT NULL,
			deleted_at DATETIME NOT NULL,
			deleted_by VARCHAR(16) NOT NULL,
			FOREIGN KEY(message_id) REFERENCES group_messages (id) ON DELETE CASCADE,
			FOREIGN KEY(id_group) REFERENCES groups (id_group) ON DELETE CASCADE,
			FOREIGN KEY(deleted_by) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS direct_message_reactions (
			message_id INTEGER NOT NULL,
			id_user VARCHAR(16) NOT NULL,
			reaction TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			PRIMARY KEY (message_id, id_user),
			FOREIGN KEY(message_id) REFERENCES messages (id) ON DELETE CASCADE,
			FOREIGN KEY(id_user) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS group_message_reactions (
			message_id INTEGER NOT NULL,
			id_user VARCHAR(16) NOT NULL,
			reaction TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			PRIMARY KEY (message_id, id_user),
			FOREIGN KEY(message_id) REFERENCES group_messages (id) ON DELETE CASCADE,
			FOREIGN KEY(id_user) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS direct_message_receipts (
			message_id INTEGER NOT NULL PRIMARY KEY,
			receiver_id VARCHAR(16) NOT NULL,
			received_at DATETIME NOT NULL,
			read_at DATETIME,
			FOREIGN KEY(message_id) REFERENCES messages (id) ON DELETE CASCADE,
			FOREIGN KEY(receiver_id) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS group_message_receipts (
			message_id INTEGER NOT NULL,
			id_group INTEGER NOT NULL,
			id_user VARCHAR(16) NOT NULL,
			received_at DATETIME NOT NULL,
			read_at DATETIME,
			PRIMARY KEY (message_id, id_user),
			FOREIGN KEY(message_id) REFERENCES group_messages (id) ON DELETE CASCADE,
			FOREIGN KEY(id_group) REFERENCES groups (id_group) ON DELETE CASCADE,
			FOREIGN KEY(id_user) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
	}

	// Iteration to create all the needed sql schemas
	for i := 0; i < len(tables); i++ {

		sqlStmt := tables[i]
		_, err := db.Exec(sqlStmt)

		if err != nil {
			return err
		}
	}
	return nil
}
