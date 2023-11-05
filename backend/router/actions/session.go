package actions

import (
	"GrooveGuru/db"
	"GrooveGuru/ent"
	User "GrooveGuru/ent/user"
	"GrooveGuru/env"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Register creates a new user and session and sets Authorization cookie.
// returns a 400 if the username or email is already taken.
// returns a 201 if the user and session are created.
func Register(c *fiber.Ctx, password, username, email string) error {
	// check if the username is already taken.
	exists, err := client.User.
		Query().
		Where(User.UsernameEQ(username)).
		Exist(ctx)
	if err != nil {
		logError("Register", "username check", err)
		return internalServerError(c, "error checking username")
	} else if exists {
		return badRequest(c, "username already exists")
	}

	// check if an account with the email already exists.
	exists, err = client.User.
		Query().
		Where(User.EmailEQ(email)).
		Exist(ctx)
	if err != nil {
		logError("Register", "email check", err)
		return internalServerError(c, "error checking email")
	} else if exists {
		return badRequest(c, "email already exists")
	}

	// reassign password to the hashed and salted version.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		logError("Register", "encrypt password", err)
		return c.Status(500).SendString("A encryption error occurred")
	}
	password = string(hashedPassword)

	// create and save the user.
	user, err := client.User.Create().
		SetEmail(email).
		SetPassword(password).
		SetUsername(username).
		Save(ctx)
	if ent.IsValidationError(err) {
		// invalid username or email was provided.
		return badRequest(c, db.ValidateUser(username, email))
	} else if err != nil {
		logError("Register", "create user", err)
		return internalServerError(c, "error creating user")
	}

	// manage session creation and cookie.
	token := uuid.New().String()
	csrf := uuid.New().String()
	expiration := time.Now().Add(week)

	_, err = client.Session.Create().
		SetToken(token).
		SetUser(user).
		SetCsrf(csrf).
		SetExpiration(expiration).
		Save(ctx)
	if err != nil {
		logError("Register", "create session", err)
		return internalServerError(c, "error creating session")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    token,
		Expires:  expiration,
		HTTPOnly: true,
		SameSite: env.SameSite,
		Secure:   env.Secure,
	})

	c.Cookie(&fiber.Cookie{
		Name:    "Csrf",
		Value:   csrf,
		Expires: expiration,
		// HttpOnly is set to false because the frontend needs to access it.
		// This isn't a security risk because the cookie is for CSRF protection;
		// If XSS is present, the attacker can already do anything they want.
		HTTPOnly: false,
		SameSite: env.SameSite,
		Secure:   env.Secure,
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"acknowledged": true,
		"message":      "user " + username + " created",
		"user": fiber.Map{
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Login creates a new session and sets Authorization cookie.
// returns a 400 if the username does not exist or the password is incorrect.
// returns a 201 if the session is created.
func Login(c *fiber.Ctx, username, password string) error {

	// grab user from database.
	user, err := client.User.
		Query().
		Where(User.UsernameEQ(username)).
		First(ctx)
	if ent.IsNotFound(err) {
		return badRequest(c, "username does not exist")
	} else if err != nil {
		logError("Login", "check user", err)
		return internalServerError(c, "error getting account")
	}

	// check against stored password.
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return badRequest(c, "incorrect password")
	} else if err != nil {
		logError("Login", "check password", err)
		return badRequest(c, "error while authorizing")
	}

	// manage session creation and cookie.
	token := uuid.New().String()
	csrf := uuid.New().String()
	expiration := time.Now().Add(week)

	_, err = client.Session.Create().
		SetToken(token).
		SetUser(user).
		SetCsrf(csrf).
		SetExpiration(expiration).
		Save(ctx)
	if err != nil {
		logError("Login", "create session", err)
		return internalServerError(c, "error creating session")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    token,
		Expires:  expiration,
		HTTPOnly: true,
		SameSite: env.SameSite,
		Secure:   env.Secure,
	})

	c.Cookie(&fiber.Cookie{
		Name:    "Csrf",
		Value:   csrf,
		Expires: expiration,
		// HttpOnly is set to false because the frontend needs to access it.
		// This isn't a security risk because the cookie is for CSRF protection;
		// If XSS is present, the attacker can already do anything they want.
		HTTPOnly: false,
		SameSite: env.SameSite,
		Secure:   env.Secure,
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"acknowledged": true,
		"message":      "user " + username + " logged in",
	})
}

// Logout deletes the session and clears the Authorization cookie.
// returns a 204 if the session is deleted.
func Logout(c *fiber.Ctx) error {
	session := c.Locals("session").(*ent.Session)

	err := client.Session.DeleteOne(session).Exec(ctx)
	if err != nil {
		logError("Logout", "delete session", err)
		// we don't need to alert the user this failed. (it shouldn't fail anyway)
		// they will lose access to their account, and the session background worker will clean it up.
	}

	c.ClearCookie("Authorization")
	c.ClearCookie("Csrf")

	return c.SendStatus(fiber.StatusNoContent)
}
