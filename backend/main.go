package main

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// connection to database
var db *sql.DB

// cache for database
var dbCache *cache.Cache

// general message for REST-responses
type responseMessage struct {
	Status  int
	Message string
	Data    any
}

// query the database
func dbSelect[T any](table string, where string, args ...any) ([]T, error) {
	// validate columns against struct T
	tType := reflect.TypeOf(new(T)).Elem()
	columns := make([]string, tType.NumField())

	validColumns := make(map[string]any)
	for ii := 0; ii < tType.NumField(); ii++ {
		field := tType.Field(ii)
		validColumns[strings.ToLower(field.Name)] = struct{}{}
		columns[ii] = strings.ToLower(field.Name)
	}

	for _, col := range columns {
		if _, ok := validColumns[strings.ToLower(col)]; !ok {
			return nil, fmt.Errorf("invalid column: %s for struct type %T", col, new(T))
		}
	}

	// create the query
	completeQuery := fmt.Sprintf("SELECT %s FROM %s", strings.Join(columns, ", "), table)

	if where != "" && where != "*" {
		completeQuery = fmt.Sprintf("%s WHERE %s", completeQuery, where)
	}

	var rows *sql.Rows
	var err error

	if len(args) > 0 {
		db.Ping()

		rows, err = db.Query(completeQuery, args...)
	} else {
		db.Ping()

		rows, err = db.Query(completeQuery)
	}

	if err != nil {
		logger.Error().Msgf("database access failed with error %v", err)

		return nil, err
	}

	defer rows.Close()
	results := []T{}

	title := cases.Title(language.Und)

	for rows.Next() {
		var lineResult T

		scanArgs := make([]any, len(columns))
		v := reflect.ValueOf(&lineResult).Elem()

		for ii, col := range columns {
			colTitle := title.String(col)

			field := v.FieldByName(colTitle)

			if field.IsValid() && field.CanSet() {
				scanArgs[ii] = field.Addr().Interface()
			} else {
				logger.Warn().Msgf("Field %s not found in struct %T", col, lineResult)
				scanArgs[ii] = new(any) // save dummy value
			}
		}

		// scan the row into the struct
		if err := rows.Scan(scanArgs...); err != nil {
			logger.Warn().Msgf("Scan-error: %v", err)

			return nil, err
		}

		results = append(results, lineResult)
	}

	if err := rows.Err(); err != nil {
		logger.Error().Msgf("rows-error: %v", err)
		return nil, err
	} else {
		return results, nil
	}
}

// insert data intot the databse
func dbInsert(table string, vals any) error {
	// extract columns from vals
	v := reflect.ValueOf(vals)
	t := v.Type()

	columns := make([]string, t.NumField())
	values := make([]any, t.NumField())

	for ii := 0; ii < t.NumField(); ii++ {
		fieldValue := v.Field(ii)

		field := t.Field(ii)

		columns[ii] = strings.ToLower(field.Name)
		values[ii] = fieldValue.Interface()
	}

	placeholders := strings.Repeat(("?, "), len(columns))
	placeholders = strings.TrimSuffix(placeholders, ", ")

	completeQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(columns, ", "), placeholders)

	_, err := db.Exec(completeQuery, values...)

	return err
}

// update data in the database
func dbUpdate(table string, set, where any) error {
	setV := reflect.ValueOf(set)
	setT := setV.Type()

	setColumns := make([]string, setT.NumField())
	setValues := make([]any, setT.NumField())

	for ii := 0; ii < setT.NumField(); ii++ {
		fieldValue := setV.Field(ii)

		field := setT.Field(ii)

		setColumns[ii] = strings.ToLower(field.Name) + " = ?"
		setValues[ii] = fieldValue.Interface()
	}

	whereV := reflect.ValueOf(where)
	whereT := whereV.Type()

	whereColumns := make([]string, whereT.NumField())
	whereValues := make([]any, whereT.NumField())

	for ii := 0; ii < whereT.NumField(); ii++ {
		fieldValue := whereV.Field(ii)

		// skip empty (zero) values
		if !fieldValue.IsZero() {
			field := whereT.Field(ii)

			whereColumns[ii] = strings.ToLower(field.Name) + " = ?"
			whereValues[ii] = fmt.Sprint(fieldValue.Interface())
		}
	}

	sets := strings.Join(setColumns, ", ")
	wheres := strings.Join(whereColumns, " AND ")

	placeholderValues := append(setValues, whereValues...)

	completeQuery := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, sets, wheres)

	_, err := db.Exec(completeQuery, placeholderValues...)

	return err
}

// remove data from the database
func dbDelete(table string, vals any) error {
	// extract columns from vals
	v := reflect.ValueOf(vals)
	t := v.Type()

	columns := make([]string, t.NumField())
	values := make([]any, t.NumField())

	for ii := 0; ii < t.NumField(); ii++ {
		fieldValue := v.Field(ii)

		// skip empty (zero) values
		if !fieldValue.IsZero() {
			field := t.Field(ii)

			columns[ii] = strings.ToLower(field.Name) + " = ?"
			values[ii] = fmt.Sprint(fieldValue.Interface())
		}
	}

	completeQuery := fmt.Sprintf("DELETE FROM %s WHERE %s", table, strings.Join(columns, ", "))

	_, err := db.Exec(completeQuery, values...)

	return err
}

// answer the client request with the response-message
func (result responseMessage) send(c *fiber.Ctx) error {
	// if the status-code is in the error-region, return an error
	if result.Status >= 400 {
		// if available, include the message
		if result.Message != "" {
			return fiber.NewError(result.Status, result.Message)
		} else {
			return fiber.NewError(result.Status)
		}
	} else {
		// if there is data, send it as JSON
		if result.Data != nil {
			c.JSON(result.Data)

			// if there is a message, send it instead
		} else if result.Message != "" {
			c.SendString(result.Message)
		}

		return c.SendStatus(result.Status)
	}
}

// payload of the JSON webtoken
type JWTPayload struct {
	Uid int `json:"uid"`
	Tid int `json:"tid"`
}

// complete JSON webtoken
type JWT struct {
	Payload
	CustomClaims JWTPayload
}

// extracts the json webtoken from the request
//
// @returns (uID, tID, error)
func extractJWT(c *fiber.Ctx) (int, int, error) {
	// get the session-cookie
	cookie := c.Cookies("session")

	token, err := jwt.ParseWithClaims(cookie, &JWT{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected JWT signing method: %v", token.Header["alg"])
		}

		return []byte(config.ClientSession.JwtSignature), nil
	})

	if err != nil {
		return -1, -1, err
	}

	// extract the claims from the JWT
	if claims, ok := token.Claims.(*JWT); ok && token.Valid {
		return claims.CustomClaims.Uid, claims.CustomClaims.Tid, nil
	} else {
		return -1, -1, fmt.Errorf("invalid JWT")
	}
}

// checks wether the request is from a valid user
func checkUser(c *fiber.Ctx) (bool, error) {
	uid, tid, err := extractJWT(c)

	if err != nil {
		return false, nil
	}

	// retrieve the user from the database
	response, err := dbSelect[UserDB]("users", "uid = ? LIMIT 1", uid)

	if err != nil {
		return false, err
	}

	// if exactly one user came back and the tID is valid, the user is authorized
	return len(response) == 1 && response[0].Tid == tid, err
}

// checks wether the request is from the admin
func checkAdmin(c *fiber.Ctx) (bool, error) {
	uid, tid, err := extractJWT(c)

	if err != nil {
		return false, err
	}

	// retrieve the user from the database
	response, err := dbSelect[UserDB]("users", "uid = ? LIMIT 1", uid)

	if err != nil {
		return false, err
	}

	// if exactly one user came back and its name is "admin", the user is the admin
	if len(response) != 1 {
		return false, fmt.Errorf("user doesn't exist")
	} else {
		return response[0].Name == "admin" && response[0].Tid == tid, err
	}
}

// information about an element in the database
type ElementDBNoReservation struct {
	Mid  string `json:"mid"`
	Name string `json:"name"`
}
type ElementDB struct {
	Mid         string  `json:"mid"`
	Name        string  `json:"name"`
	Reservation *string `json:"reservation"`
}

// client-data of the reserved elements
type ClientStatus struct {
	ReservedElements map[string]string `json:"reserved_elements"`
}

// caches the elements from the database
func cacheElements() error {
	if res, err := dbSelect[ElementDB]("elements", "*"); err != nil {
		return err
	} else {
		elementMap := make(map[string]string)

		for _, element := range res {
			elementMap[string(element.Mid[:])] = element.Name
		}

		dbCache.Set("elements", elementMap, cache.DefaultExpiration)

		return nil
	}
}

// gets the elements from the cache
func getElements(c *fiber.Ctx) responseMessage {
	response := responseMessage{}

	elements, found := dbCache.Get("elements")

	if !found {
		if err := cacheElements(); err != nil {
			response.Status = fiber.StatusInternalServerError
			response.Message = "can't get elements"

			logger.Error().Msgf("can't get elements from database: %v", err)
		} else if elements, found = dbCache.Get("elements"); !found {
			response.Status = fiber.StatusInternalServerError
			response.Message = "can't get elements"

			logger.Error().Msg("can't get 'elements' from cache")
		}
	}

	// if the reponse-status is still unset, there was no error
	if response.Status == 0 {
		response.Data = ClientStatus{
			ReservedElements: elements.(map[string]string),
		}

		logger.Debug().Msg("retrieved elements")
	}

	return response
}

// regex to match valid element-names
var midRegex = regexp.MustCompile(`^(pv-\w|(?:wr|bs)-)(\d{1,2})$`)

var validElements = map[string]struct {
	from int
	to   int
}{
	"bs-": {
		from: 1,
		to:   2,
	},
	"pv-a": {
		from: 1,
		to:   16,
	},
	"pv-b": {
		from: 2,
		to:   37,
	},
	"pv-c": {
		from: 3,
		to:   37,
	},
	"pv-d": {
		from: 3,
		to:   37,
	},
	"pv-e": {
		from: 1,
		to:   6,
	},
	"pv-f": {
		from: 1,
		to:   6,
	},
	"pv-g": {
		from: 1,
		to:   6,
	},
	"pv-h": {
		from: 1,
		to:   7,
	},
	"pv-i": {
		from: 1,
		to:   7,
	},
	"pv-j": {
		from: 1,
		to:   7,
	},
	"pv-k": {
		from: 1,
		to:   7,
	},
	"pv-l": {
		from: 1,
		to:   7,
	},
	"pv-m": {
		from: 1,
		to:   7,
	},
	"pv-n": {
		from: 1,
		to:   7,
	},
	"pv-o": {
		from: 1,
		to:   7,
	},
	"pv-p": {
		from: 1,
		to:   7,
	},
	"pv-q": {
		from: 1,
		to:   7,
	},
	"pv-r": {
		from: 1,
		to:   7,
	},
	"pv-s": {
		from: 1,
		to:   7,
	},
	"pv-t": {
		from: 1,
		to:   7,
	},
	"pv-u": {
		from: 1,
		to:   7,
	},
	"pv-v": {
		from: 1,
		to:   7,
	},
	"wr-": {
		from: 1,
		to:   4,
	},
}

func isValidMid(element string) (bool, error) {
	if results := midRegex.FindStringSubmatch(element); results == nil {
		return false, nil
	} else {
		// check wether the descriptor-part is valid
		if rng, ok := validElements[results[1]]; !ok {
			return false, nil

			// try to parse the mid-number
		} else if n, err := strconv.Atoi(results[2]); err != nil {
			return false, err
		} else {
			return rng.from <= n && n <= rng.to, nil
		}
	}
}

// handles post-requests for reserving new elements
func postElements(c *fiber.Ctx) responseMessage {
	response := responseMessage{}

	body := struct{ Name string }{}

	mid := c.Query("mid")

	if ok, err := isValidMid(mid); err != nil || !ok {
		response.Status = fiber.StatusBadRequest
		response.Message = "invalid mID"

		logger.Info().Msgf("can't reserve element: invalid element-name: %q", mid)
	} else if err := c.BodyParser(&body); err != nil {
		response.Status = fiber.StatusBadRequest
		response.Message = "invalid message-body"

		logger.Warn().Msg(`body can't be parsed as "struct{ name string }"`)
	} else {
		elements, found := dbCache.Get("elements")

		if !found {
			if err := cacheElements(); err != nil {
				response.Status = fiber.StatusInternalServerError
				response.Message = "can't get elements"

				logger.Error().Msgf("can't get elements from database: %v", err)
			} else if elements, found = dbCache.Get("elements"); !found {
				response.Status = fiber.StatusInternalServerError
				response.Message = "can't get elements"

				logger.Error().Msg("can't get 'elements' from cache")
			}
		}

		// if the status is still unset, there was no error
		if response.Status == 0 {
			// check wether the element already exists
			if _, ok := elements.(map[string]string)[mid]; ok {
				response.Status = fiber.StatusBadRequest
				response.Message = "element is already reserved"

				logger.Info().Msgf("element %q is already reserved", mid)

				return response
			}

			// clear the current cache
			dbCache.Delete("elements")

			// write the data to the database
			if err := dbInsert("elements", ElementDBNoReservation{Mid: mid, Name: body.Name}); err != nil {
				response.Status = fiber.StatusInternalServerError
				response.Message = "error while writing reservation to database"

				logger.Error().Msgf("can't write reservation to database: %v", err)
			} else {
				response = getElements(c)

				logger.Debug().Msgf("reserved element %q", mid)
			}
		}
	}

	return response
}

func sendReservationEmail() {

}

// handles patch-requests for modifying element reservations
func patchElements(c *fiber.Ctx) responseMessage {
	response := responseMessage{}

	if user, err := checkUser(c); err != nil {
		response.Status = fiber.StatusInternalServerError

		logger.Error().Msgf("can't check user: %v", err)
	} else if !user {
		response.Status = fiber.StatusUnauthorized

		logger.Info().Msg("request is not authorized as user")
	} else {
		body := struct{ Name string }{}

		mid := c.Query("mid")
		if ok, err := isValidMid(mid); err != nil || !ok {
			response.Status = fiber.StatusBadRequest
			response.Message = "invalid element name"

			logger.Info().Msgf("can't modify element: invalid element-name: %q", mid)
		} else if err := c.BodyParser(&body); err != nil {
			response.Status = fiber.StatusBadRequest
			response.Message = "invalid message-body"

			logger.Warn().Msg(`body can't be parsed as "struct{ name string }"`)
		} else {
			// check wether the element already exists
			if elements, found := dbCache.Get("elements"); found {
				if _, ok := elements.(map[string]string)[mid]; !ok {
					response.Status = fiber.StatusBadRequest
					response.Message = "element is already reserved"

					logger.Info().Msgf("element %q is already reserved", mid)

					return response
				}
			}

			// clear the current cache
			dbCache.Delete("elements")

			// write the data to the database
			if err := dbUpdate("elements", struct{ Name string }{Name: body.Name}, struct{ Mid string }{Mid: mid}); err != nil {
				response.Status = fiber.StatusInternalServerError
				response.Message = "error while writing reservation to database"

				logger.Error().Msgf("can't write reservation to database: %v", err)
			} else {
				response = getElements(c)

				logger.Debug().Msgf("modified reservation for element %q", mid)
			}
		}
	}

	return response
}

// handle delete-requets for deleting an element reservation
func deleteElements(c *fiber.Ctx) responseMessage {
	response := responseMessage{}

	if user, err := checkUser(c); err != nil {
		response.Status = fiber.StatusInternalServerError

		logger.Error().Msgf("can't check user: %v", err)
	} else if !user {
		response.Status = fiber.StatusUnauthorized

		logger.Info().Msg("request is not authorized as user")
	} else {
		mid := c.Query("mid")

		if ok, err := isValidMid(mid); !ok || err != nil {
			response.Status = fiber.StatusBadRequest
			response.Message = "invalid element name"

			logger.Info().Msgf("can't delete element: invalid element-name: %q", mid)
		} else {
			dbCache.Delete("elements")

			if err := dbDelete("elements", struct{ Mid string }{Mid: mid}); err != nil {
				response.Status = fiber.StatusInternalServerError
				response.Message = "error while deleting reservation from database"

				logger.Error().Msgf("can't delete reservation from database: %v", err)
			} else {
				response = getElements(c)

				logger.Debug().Msgf("deleted reservation for %q", mid)
			}
		}
	}

	return response
}

// request for adding a user
type AddUserBody struct {
	Name     string
	Password string
}

// user-entry in the database
type UserDB struct {
	Uid      int    `json:"uid"`
	Name     string `json:"name"`
	Password []byte `json:"password"`
	Tid      int    `json:"tid"`
}

// hashes a password
func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// handles get-request for the users
func getUsers(c *fiber.Ctx) responseMessage {
	var response responseMessage

	if isAdmin, err := checkAdmin(c); err != nil {
		response.Status = fiber.StatusInternalServerError

		logger.Error().Msgf("can't check for admin-user: %v", err)
	} else if !isAdmin {
		response.Status = fiber.StatusUnauthorized

		logger.Info().Msg("request is not authorized as admin")
	} else {
		// retrieve all users
		if users, err := dbSelect[struct {
			Uid  int    `json:"uid"`
			Name string `json:"name"`
		}]("users", ""); err != nil {
			response.Status = fiber.StatusInternalServerError
			response.Message = "can't get users from database"

			logger.Error().Msgf("can't get users from database: %v", err)
		} else {
			response.Data = users

			logger.Debug().Msg("retrieved users from database")
		}
	}

	return response
}

// validates a password against the password-rules
func validatePassword(password string) bool {
	return len(password) >= 12 && len(password) <= 64
}

// handles post-request to add a new user to the database
func postUsers(c *fiber.Ctx) responseMessage {
	response := responseMessage{}
	body := AddUserBody{}

	if admin, err := checkAdmin(c); err != nil {
		response.Status = fiber.StatusInternalServerError

		logger.Error().Msgf("can't check for admin-user: %v", err)
	} else if !admin {
		response.Status = fiber.StatusUnauthorized

		logger.Info().Msg("request is not authorized as user")
	} else if err := c.BodyParser(&body); err != nil {
		response.Status = fiber.StatusBadRequest
		response.Message = "invalid message-body"

		logger.Warn().Msg(`body can't be parsed as "struct{ name string; Password string }"`)
	} else {
		if dbUsers, err := dbSelect[UserDB]("users", "name = ? LIMIT 1", body.Name); err != nil {
			response.Status = fiber.StatusInternalServerError

			logger.Error().Msgf("can't read users from database: %v", err)
		} else if len(dbUsers) != 0 {
			response.Status = fiber.StatusBadRequest
			response.Message = "user already exists"

			logger.Info().Msgf("can't add user: user with name %q already exists", body.Name)
		} else {
			// everything is valid
			if hashedPassword, err := hashPassword(body.Password); err != nil {
				response.Status = fiber.StatusInternalServerError

				logger.Error().Msgf("can't hash password: %v", err)
			} else {
				if err := dbInsert("users", struct {
					Name     string
					Password []byte
				}{Name: body.Name, Password: hashedPassword}); err != nil {
					response.Status = fiber.StatusInternalServerError
					response.Message = "can't add user to database"

					logger.Error().Msgf("can't add user to database: %v", err)
				} else {
					response = getUsers(c)

					logger.Debug().Msgf("added user %q", body.Name)
				}
			}
		}
	}

	return response
}

// change the password in the database
func changePassword(uid int, password string) responseMessage {
	response := responseMessage{}

	// hash the new password
	if hashedPassword, err := hashPassword(password); err != nil {
		response.Status = fiber.StatusInternalServerError

		logger.Error().Msgf("can't hash password: %v", err)
	} else {
		// increase the token-id of the user to make the current-token invalid
		if err := incTokenId(uid); err != nil {
			response.Status = fiber.StatusInternalServerError

			logger.Error().Msgf("can't increase the tid: %v", err)
		} else {
			// update the databse with the new password
			if err := dbUpdate("users", struct{ Password []byte }{Password: hashedPassword}, struct{ Uid int }{Uid: uid}); err != nil {
				response.Status = fiber.StatusInternalServerError
				response.Message = "can't update password"

				logger.Error().Msgf("can't update password: %v", err)
			} else {
				logger.Debug().Msgf("updated password for user %q", uid)

				response.Status = fiber.StatusOK
			}
		}
	}

	return response
}

// handles patch-request to change a useres password
func patchUsers(c *fiber.Ctx) responseMessage {
	response := responseMessage{}

	if admin, err := checkAdmin(c); err != nil {
		response.Status = fiber.StatusInternalServerError
		response.Message = "error while checking the authorization"

		logger.Error().Msgf("can't check for admin: %v", err)
	} else if !admin {
		response.Status = fiber.StatusUnauthorized

		logger.Info().Msg("user is no admin")
	} else {
		body := struct {
			Password string `json:"password"`
		}{}

		// check wether a valid uid is present
		if uid := c.QueryInt("uid", -1); uid < 0 {
			response.Status = fiber.StatusBadRequest
			response.Message = "query doesn't include valid uid"

			logger.Info().Msg("query doesn't include valid uid")
		} else {
			// try to parse the body
			if err := c.BodyParser(&body); err != nil {
				response.Status = fiber.StatusBadRequest
				response.Message = "invalid message-body"

				logger.Warn().Msg(`body can't be parsed as "struct{ password string }"`)
			} else {
				// check, wether the user exists
				if dbUsers, err := dbSelect[UserDB]("users", "uid = ? LIMIT 1", uid); err != nil {
					response.Status = fiber.StatusInternalServerError

					logger.Error().Msgf("can't read users from database: %v", err)
				} else if len(dbUsers) != 1 {
					response.Status = fiber.StatusBadRequest
					response.Message = "user doesn't exist"

					logger.Info().Msgf("can't modify user: user with uid %q doesn't exist", uid)
				} else {
					// everything is valid

					if response = changePassword(uid, body.Password); response.Status == fiber.StatusOK {
						response = getUsers(c)
					}
				}
			}
		}
	}

	return response
}

// handle delete-request for removing a user
func deleteUsers(c *fiber.Ctx) responseMessage {
	response := responseMessage{}

	if admin, err := checkAdmin(c); err != nil {
		response.Status = fiber.StatusInternalServerError

		logger.Error().Msgf("can't check for admin-user: %v", err)
	} else if !admin {
		response.Status = fiber.StatusUnauthorized

		// check wether there is a valid uid
	} else if uid := c.QueryInt("uid", -1); uid < 0 {
		response.Status = fiber.StatusBadRequest
		response.Message = "query doesn't include valid uid"

		logger.Info().Msg("query doesn't include valid uid")
	} else {
		// delete the user from the database
		if err := dbDelete("users", struct{ Uid int }{Uid: uid}); err != nil {
			response.Status = fiber.StatusInternalServerError
			response.Message = "can't delete user"

			logger.Error().Msgf("can't delete user with uid = %q: %v", uid, err)
		} else {
			logger.Debug().Msgf("deleted user with uid = %q", uid)

			response = getUsers(c)
		}
	}

	return response
}

// handles patch-requests to change the users password
func patchUserPassword(c *fiber.Ctx) responseMessage {
	response := responseMessage{}

	if user, err := checkUser(c); err != nil {
		response.Status = fiber.StatusInternalServerError

		logger.Error().Msgf("can't check for user: %v", err)
	} else if !user {
		response.Status = fiber.StatusUnauthorized

		logger.Info().Msg("request is not authorized as user")
	} else {
		// parse the body
		var body struct {
			Password string `json:"password"`
		}

		if uid, _, err := extractJWT(c); err != nil {
			response.Status = fiber.StatusBadRequest
			response.Message = "query doesn't include valid uid"

			logger.Warn().Msg("can't extract uid from query")
		} else if err := c.BodyParser(&body); err != nil {
			response.Status = fiber.StatusBadRequest

			logger.Warn().Msg(`body can't be parsed as "struct{ password string }"`)
		} else if !validatePassword(body.Password) {
			response.Status = fiber.StatusBadRequest
			response.Message = "invalid password"

			logger.Info().Msg("invalid password")
		} else {
			// everything is valid

			return changePassword(uid, body.Password)
		}
	}

	return response
}

// handle welcome-messages from clients
func handleWelcome(c *fiber.Ctx) error {
	logger.Debug().Msgf("HTTP %s request: %q", c.Method(), c.OriginalURL())

	response := responseMessage{}
	response.Data = UserLogin{
		LoggedIn: false,
	}

	if ok, err := checkUser(c); err != nil {
		response.Status = fiber.StatusInternalServerError

		logger.Warn().Msgf("can't check user: %v", err)
	} else if !ok {
		response.Status = fiber.StatusNoContent
	} else {
		if uid, _, err := extractJWT(c); err != nil {
			response.Status = fiber.StatusBadRequest

			logger.Error().Msgf("can't extract JWT: %v", err)
		} else {
			if users, err := dbSelect[UserDB]("users", "uid = ? LIMIT 1", strconv.Itoa(uid)); err != nil {
				response.Status = fiber.StatusInternalServerError

				logger.Error().Msgf("can't get users from database: %v", err)
			} else {
				if len(users) != 1 {
					response.Status = fiber.StatusForbidden
					response.Message = "unknown user"

					removeSessionCookie(c)
				} else {
					user := users[0]

					response.Data = UserLogin{
						Uid:      user.Uid,
						Name:     user.Name,
						LoggedIn: true,
					}
				}

				logger.Debug().Msgf("welcomed user with uid = %v", uid)
			}
		}
	}

	return response.send(c)
}

// body from a login-request
type LoginBody struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

// data of the logged-in-status
type UserLogin struct {
	Uid      int    `json:"uid"`
	Name     string `json:"name"`
	LoggedIn bool   `json:"logged_in"`
}

// retrieves the current tid for a specific user from the database
func getTokenId(uid int) (int, error) {
	if response, err := dbSelect[UserDB]("users", "uid = ? LIMIT 1", uid); err != nil {
		return -1, err
	} else if len(response) != 1 {
		return -1, fmt.Errorf("can't get user with uid = %q from database", uid)
	} else {
		return response[0].Tid, nil
	}
}

// increases the tid of a user
func incTokenId(uid int) error {
	_, err := db.Exec("UPDATE users SET tid = tid + 1 WHERE uid = ?", uid)

	return err
}

var messageWrongLogin = "Unkown user or wrong password"

// handles login-requests
func handleLogin(c *fiber.Ctx) error {
	logger.Debug().Msgf("HTTP %s request: %q", c.Method(), c.OriginalURL())

	var response responseMessage

	body := LoginBody{}

	if err := c.BodyParser(&body); err != nil {
		response.Status = fiber.StatusBadRequest
		response.Message = "can't parse message-body"

		logger.Warn().Msgf("can't parse login-body: %v", err)
	} else {
		// try to get the hashed password from the database
		dbResult, err := dbSelect[UserDB]("users", "name = ? LIMIT 1", body.User)

		if err != nil {
			response.Status = fiber.StatusInternalServerError

			logger.Error().Msgf("can't get users from the database: %v", err)
		} else if len(dbResult) != 1 {
			response.Status = fiber.StatusForbidden
			response.Message = messageWrongLogin

			logger.Info().Msgf("user with name = %q doesn't exist", body.User)
		} else {
			response.Data = UserLogin{
				LoggedIn: false,
			}

			user := dbResult[0]

			if len(dbResult) != 1 || bcrypt.CompareHashAndPassword(user.Password, []byte(body.Password)) != nil {
				response.Status = fiber.StatusUnauthorized
				response.Message = messageWrongLogin

				logger.Debug().Msgf("can't login: wrong username or password")
			} else {
				// get the token-id
				if tid, err := getTokenId(user.Uid); err != nil {
					response.Status = fiber.StatusInternalServerError

					logger.Error().Msgf("can't get tid for user with uid = %q", user.Uid)
				} else {
					// create the jwt
					jwt, err := config.signJWT(JWTPayload{
						Uid: user.Uid,
						Tid: tid,
					})

					if err != nil {
						response.Status = fiber.StatusInternalServerError

						logger.Error().Msgf("json-webtoken creation failed: %v", err)
					} else {
						c.Cookie(&fiber.Cookie{
							Name:     "session",
							Value:    jwt,
							HTTPOnly: true,
							SameSite: "strict",
							MaxAge:   int(config.SessionExpire.Seconds()),
						})

						response.Data = UserLogin{
							Uid:      user.Uid,
							Name:     user.Name,
							LoggedIn: true,
						}

						logger.Info().Msgf("user with uid = %q logged in", user.Uid)
					}
				}
			}
		}
	}

	return response.send(c)
}

// removes the session-coockie from a request
func removeSessionCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    "",
		HTTPOnly: true,
		SameSite: "strict",
		Expires:  time.Unix(0, 0),
	})
}

// handles logout-requests
func handleLogout(c *fiber.Ctx) error {
	logger.Debug().Msgf("HTTP %s request: %q", c.Method(), c.OriginalURL())

	removeSessionCookie(c)

	return responseMessage{
		Data: UserLogin{
			LoggedIn: false,
		},
	}.send(c)
}

func main() {
	// setup the database-connection
	sqlConfig := mysql.Config{
		AllowNativePasswords: true,
		Net:                  "tcp",
		User:                 config.Database.User,
		Passwd:               config.Database.Password,
		Addr:                 config.Database.Host,
		DBName:               config.Database.Database,
	}

	// connect to the database
	db, _ = sql.Open("mysql", sqlConfig.FormatDSN())
	db.SetMaxIdleConns(10)
	db.SetMaxIdleConns(100)
	db.SetConnMaxLifetime(time.Minute)

	// setup the cache
	dbCache = cache.New(config.Cache.Expiration, config.Cache.Purge)

	// setup fiber
	app := fiber.New(fiber.Config{
		AppName:               "johannes-pv",
		DisableStartupMessage: true,
	})

	// map with the individual methods
	handleMethods := map[string]func(path string, handlers ...func(*fiber.Ctx) error) fiber.Router{
		"GET":    app.Get,
		"POST":   app.Post,
		"PATCH":  app.Patch,
		"DELETE": app.Delete,
	}

	// map with the individual registered endpoints
	endpoints := map[string]map[string]func(*fiber.Ctx) responseMessage{
		"GET": {
			"elements": getElements,
			"users":    getUsers,
		},
		"POST": {
			"elements": postElements,
			"users":    postUsers,
		},
		"PATCH": {
			"elements":      patchElements,
			"users":         patchUsers,
			"user/password": patchUserPassword,
		},
		"DELETE": {
			"elements": deleteElements,
			"users":    deleteUsers,
		},
	}

	// handle specific requests special
	app.Get("/pv/api/welcome", handleWelcome)
	app.Post("/pv/api/login", handleLogin)
	app.Get("/pv/api/logout", handleLogout)

	// register the registered endpoints
	for method, handlers := range endpoints {
		for address, handler := range handlers {
			handleMethods[method]("/pv/api/"+address, func(c *fiber.Ctx) error {
				logger.Debug().Msgf("HTTP %s request: %q", c.Method(), c.OriginalURL())

				return handler(c).send(c)
			})
		}
	}

	// start the server
	app.Listen(fmt.Sprintf("localhost:%d", config.Server.Port))
}
