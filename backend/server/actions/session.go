package actions

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"groove/pkgs/db"
	"groove/pkgs/ent"
	SpotifyLink "groove/pkgs/ent/spotifylink"
	User "groove/pkgs/ent/user"
	. "groove/pkgs/util"
	"time"
)

// Register creates a new user and session and sets Authorization cookie.
// returns a 400 if the username or email is already taken/invalid.
// returns a 201 if the user and session are created.
func (a *Actions) Register(c *fiber.Ctx, password, username, email string) error {
	// validate user input in order to prevent unnecessary database calls.
	if err := db.ValidateUser(username, email, password); err != nil {
		return BadRequest(c, err.Error())
	}

	ctx := c.Context()
	exists, err := a.client.User.
		Query().
		Where(User.UsernameEQ(username)).
		Exist(ctx)
	if err != nil {
		LogError("Register", "username check", err)
		return InternalServerError(c, "error checking username")
	} else if exists {
		return BadRequest(c, "username already exists")
	}

	// check if email is already taken.
	exists, err = a.client.User.
		Query().
		Where(User.EmailEQ(email)).
		Exist(ctx)
	if err != nil {
		LogError("Register", "email check", err)
		return InternalServerError(c, "error checking email")
	} else if exists {
		return BadRequest(c, "email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		LogError("Register", "hash password", err)
		return InternalServerError(c, "an error occurred")
	}

	// create and save the user.
	user, err := a.client.User.Create().
		SetEmail(email).
		SetPassword(string(hashedPassword)).
		SetUsername(username).
		Save(ctx)
	if err != nil {
		LogError("Register", "create user", err)
		return InternalServerError(c, "error creating account")
	}

	// manage session creation and cookie.
	token := uuid.New().String()
	csrf := uuid.New().String()
	expiration := time.Now().Add(TimeWeek)

	_, err = a.client.Session.Create().
		SetToken(token).
		SetUser(user).
		SetCsrf(csrf).
		SetExpiration(expiration).
		Save(ctx)
	if err != nil {
		LogError("Register", "create session", err)
		return InternalServerError(c, "error creating session")
	}
	SetSessionCookies(c, token, csrf, expiration, a.env.SameSite, a.env.Secure)

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
func (a *Actions) Login(c *fiber.Ctx, username, password string) error {
	ctx := c.Context()
	// check if username exists.
	user, err := a.client.User.
		Query().
		Where(User.UsernameEQ(username)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return BadRequest(c, "username does not exist")
		}
		LogError("Login", "check user", err)
		return InternalServerError(c, "error getting account")
	}

	// check against stored password.
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return BadRequest(c, "incorrect password")
		}
		LogError("Login", "check password", err)
		return BadRequest(c, "error while authorizing")
	}

	// manage session creation and cookie.
	token := uuid.New().String()
	csrf := uuid.New().String()
	expiration := time.Now().Add(TimeWeek)

	_, err = a.client.Session.Create().
		SetToken(token).
		SetUser(user).
		SetCsrf(csrf).
		SetExpiration(expiration).
		Save(ctx)
	if err != nil {
		LogError("Login", "create session", err)
		return InternalServerError(c, "error creating session")
	}
	SetSessionCookies(c, token, csrf, expiration, a.env.SameSite, a.env.Secure)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"acknowledged": true,
		"user": fiber.Map{
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Logout deletes the session and clears the Authorization cookie.
// returns a 204 if the session is deleted.
func (a *Actions) Logout(c *fiber.Ctx) error {
	session := c.Locals("session").(*ent.Session)
	ctx := c.Context()

	if err := a.client.Session.DeleteOne(session).Exec(ctx); err != nil {
		LogError("Logout", "delete session", err)
		// we don't need to alert the user this failed. (it shouldn't fail anyway)
		// they will lose access to their account, and the session background worker will clean it up.
	}

	ExpireSessionCookies(c, a.env.SameSite, a.env.Secure)
	return c.SendStatus(fiber.StatusNoContent)
}

// Authenticate resets the session expiration and returns the user's username, email, and status for spotify link.
// returns a 200 if the session is updated.
func (a *Actions) Authenticate(c *fiber.Ctx) error {
	ctx := c.Context()
	session := c.Locals("session").(*ent.Session)

	user, err := a.client.User.
		Query().
		Where(User.IDEQ(session.UserID)).
		First(ctx)
	if err != nil {
		LogError("Authenticate", "check user", err)
		return InternalServerError(c, "error getting account")
	}

	// refresh expirations.
	expiration := time.Now().Add(TimeWeek)
	if _, err = session.Update().SetExpiration(expiration).Save(ctx); err != nil {
		LogError("Authenticate", "update session", err)
		return InternalServerError(c, "error updating session")
	}
	// refresh cookie expiration with same values.
	SetSessionCookies(c, session.Token, session.Csrf, expiration, a.env.SameSite, a.env.Secure)

	// check for spotify link.
	exists, err := a.client.SpotifyLink.
		Query().
		Where(SpotifyLink.UserIDEQ(session.UserID)).
		Exist(ctx)
	if err != nil {
		LogError("Authenticate", "check spotify link", err)
		return InternalServerError(c, "error checking spotify account")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": fiber.Map{
			"username": user.Username,
			"email":    user.Email,
			"spotify":  exists,
		},
	})
}
