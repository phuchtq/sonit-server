package businesslogic

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	action_type "sonit_server/constant/action_type"
	"sonit_server/constant/env"
	mail_const "sonit_server/constant/mail_const"
	"sonit_server/constant/noti"
	page_url "sonit_server/constant/page_url"
	data_access "sonit_server/interface/data_access"
	"sonit_server/model/dto/request"
	entity "sonit_server/model/entity"
	"sonit_server/utils"
	"sync"
	"time"
)

const (
	sepChar        string = "|"
	mailSepChar    string = ":"
	id_validate    string = "ID_VALIDATE"
	email_validate string = "EMAIL_VALIDATE"
)

const (
	max_fail_attempts int = 5
)

const (
	min_length_verify_mail_component_combine int = 3
)

// Specific cases as a component of token action
const (
	activateType      string = "1"
	resetPassType     string = "2"
	updateProfileType string = "3"
	verifyType        string = "4"
)

const (
	Sale_action         string = "sale"
	Import_action       string = "import"
	Export_action       string = "export"
	Return_action       string = "return"
	Cancel_order_action string = "cancel"
)

// Supported filter property
const (
	date_filter  string = "DATE_FILTER"
	price_filter string = "PRICE_FILTER"
)

// Filter order
const (
	desc_order string = "DESC"
	asc_order  string = "ASC"
)

// -------------------- Generic usages --------------------
func isEntityExist(repo interface{}, prob, validateProb string, ctx context.Context) bool {
	var res bool
	switch r := repo.(type) {
	case data_access.IRoleRepo:
		if validateProb == id_type {
			data, _ := r.GetRoleById(prob, ctx)
			res = data != nil
		} else {
			data, _ := r.GetRolesByName(prob, ctx)
			res = data != nil
		}
	case data_access.IUserRepo:
		if validateProb == id_type {
			data, _ := r.GetUser(prob, ctx)
			res = data != nil
		} else {
			data, _ := r.GetUserByEmail(prob, ctx)
			res = data != nil
		}
	case data_access.ICollectionRepo:
		if validateProb == id_type {
			data, _ := r.GetCollectionById(prob, ctx)
			res = data != nil
		} else {
			data, _ := r.GetCollectionsByName(prob, ctx)
			res = data != nil
		}

	case data_access.IVoucherRepo:
		if validateProb == id_type {
			data, _ := r.GetVoucherByID(prob, ctx)
			res = data != nil
		} else {
			data, _ := r.GetVoucherByCode(prob, ctx)
			res = data != nil
		}
	case data_access.ICategoryRepo:
		if validateProb == id_type {
			data, _ := r.GetCategoryById(prob, ctx)
			res = data != nil
		} else {
			data, _ := r.GetCategoriesByName(prob, ctx)
			res = data != nil
		}
	case data_access.IProductRepo:
		if validateProb == id_type {
			data, _ := r.GetProductById(prob, ctx)
			res = data != nil
		} else {
			data, _, _ := r.GetProductsByName(1, prob, ctx)
			res = data != nil
		}
	}

	return res
}

func closeCnn(cnn *sql.DB) {
	cnn.Close()
}

// -------------------- ~~~~~ --------------------

// -------------------- USER SERVICE HELPER --------------------

func validateActorAndAffectedAccount(actorId, accountId string, userRepo data_access.IUserRepo, ctx context.Context) error {
	// Fetch account
	account, err := userRepo.GetUser(accountId, ctx)

	// Error while fetching account
	if err != nil {
		return err
	}

	// Account not existed
	if account == nil {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	// Actor and affected account is the samae
	if accountId == actorId {
		return nil
	}

	// Fetch actor
	actor, err := userRepo.GetUser(actorId, ctx)

	// Error while fetching actor
	if err != nil {
		return err
	}

	// Actor not existed
	if actor == nil {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	// Actor is not admin but edit other account's info
	if actor.RoleId != os.Getenv(env.ADMIN_ROLE) {
		return errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	return nil
}

// Verify if account is existed by fetching record from data sent by request coming supporting 2 fields: email or id
func verifyAccount(field, validateField string, user *entity.User, repo data_access.IUserRepo, ctx context.Context) error {
	if field == "" {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var res error
	var tmpUser *entity.User

	switch validateField {
	case id_validate:
		tmpUser, res = repo.GetUser(field, ctx)
	case email_validate:
		tmpUser, res = repo.GetUserByEmail(field, ctx)
	}

	// User ko tồn tại
	if tmpUser == nil && res == nil {
		// Phân chia lỗi trả về
		res = errors.New("User not found")
	}

	// Only set user if the pointer is not nil
	if tmpUser != nil {
		*user = *tmpUser
	}

	return res
}

// Processing with wrong-credentials login case
func processWrongCredentialsCase(accountId string, securityRepo data_access.IUserSecurityRepo, ctx context.Context) {
	info, err := securityRepo.GetUserSecurity(accountId, ctx)
	if err != nil {
		return
	}

	if info != nil {
		// Set fail access and last fail
		info.FailAccess += 1
		var curTime time.Time = time.Now()
		info.LastFail = &curTime

		// Update database
		securityRepo.EditUserSecurity(*info, ctx)
	}
}

// Processing with correct-credentials login case
func processCorrectCredentialsCase(account entity.User, logger *log.Logger, securityRepo data_access.IUserSecurityRepo, ctx context.Context) (string, string, error) {
	security, err := securityRepo.GetUserSecurity(account.UserId, ctx)
	if err != nil {
		return "", "", err
	}

	// Account need to be activated or verified as previous fail attempts
	if !account.IsActivated || security.FailAccess > max_fail_attempts {
		return processAccountVerifyCase(security, account.Email, logger, securityRepo, ctx)
	}

	var res1, res2 string

	// Need to reset password
	if account.IsHaveToResetPw != nil {
		// Token for verifying action
		token, err := utils.GenerateActionToken(account.Email, account.UserId, account.RoleId, logger)
		if err != nil {
			return "", "", err
		}

		res1 = action_type.REDIRECT             // Represent as redirect flag
		res2 = utils.ToCombinedString([]string{ // Represent url to be redirected
			page_url.RESET_PASSWORD_URL,
			token,
		}, mailSepChar)

		security.ActionToken = &token // Assigning created token
	} else {
		accessToken, refreshToken, err := utils.GenerateTokens(account.Email, account.UserId, account.RoleId, logger)

		if err != nil {
			return "", "", err
		}

		res1 = accessToken
		res2 = refreshToken
		security.AccessToken = &accessToken
		security.RefreshToken = &refreshToken
	}

	return res1, res2, securityRepo.EditUserSecurity(*security, ctx)
}

// Processing to prepare for sending mail in case of account needed to be activated or verified
func processAccountVerifyCase(security *entity.UserSecurity, email string, logger *log.Logger, secureRepo data_access.IUserSecurityRepo, ctx context.Context) (string, string, error) {
	var capturedErr error
	ctx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	var mu sync.Mutex

	var res1, res2 string

	wg.Add(2)

	go func() {
		defer wg.Done()

		var actionType, templatePath, subject string

		if security.FailAccess > max_fail_attempts {
			actionType = verifyType
			templatePath = mail_const.ACCOUNT_RECOVERY_MAIL_TEMPLATE
			subject = noti.VERIFY_ACCOUNT_MAIL_SUBJECT

			res1 = action_type.VERIFY_TYPE
			res2 = noti.VERIFY_ACCOUNT_MESSAGE

		} else {
			actionType = activateType
			templatePath = mail_const.ACCOUNT_REGISTRATION_MAIL_TEMPLATE
			subject = noti.REGISTRATION_ACCOUNT_MAIL_SUBJECT

			res1 = action_type.ACTIVATE_TYPE
			res2 = noti.ACTIVATE_ACCOUNT_MESSAGE
		}

		// Generate action token
		token, err := utils.GenerateActionToken(email, security.UserId, "", logger)
		if err != nil {
			// Lock other goroutines
			mu.Lock()
			if capturedErr == nil {
				// Assigning error
				capturedErr = err
				// Make a flag to cancel other goroutines
				cancel()
			}

			mu.Unlock()
		} else {
			// Assigning created token
			*security.ActionToken = token

			// Send mail
			if err := utils.SendMail(request.SendMailRequest{
				Body: request.MailBody{
					Email:   email,
					Subject: subject,
					Url: utils.ToCombinedString([]string{
						page_url.PROCESS_ENDPOINT,
						token,
						security.UserId,
						actionType,
					}, mailSepChar),
				},

				TemplatePath: templatePath,

				Logger: logger,
			}); err != nil {
				// Lock other goroutines
				mu.Lock()
				if capturedErr == nil {
					// Assigning error
					capturedErr = err
					// Make a flag to cancel other goroutines
					cancel()
				}

				mu.Unlock()
			}
		}
	}()

	go func() {
		defer wg.Done()

		security.FailAccess = 0
		// Check if LastFail is nil before dereferencing
		if security.LastFail == nil {
			currentTime := utils.GetPrimitiveTime()
			security.LastFail = &currentTime
		} else {
			*security.LastFail = utils.GetPrimitiveTime()
		}

		if err := secureRepo.EditUserSecurity(*security, ctx); err != nil {
			mu.Lock()

			if capturedErr == nil {
				capturedErr = err // Capture the first error
				cancel()          // Cancel the other goroutine
			}

			mu.Unlock()
		}
	}()

	wg.Wait()

	if capturedErr != nil {
		return "", "", capturedErr
	}

	return res1, res2, nil
}

// -------------------- ~~~~~ --------------------
// -------------------- VOUCHER SERVICE HELPER --------------------

func isVoucherAvailable(voucher entity.Voucher) bool {
	return !utils.IsActionExpired(voucher.ExpiredAt) && voucher.ActiveStatus && voucher.Amount > 0
}

// -------------------- ~~~~~ --------------------
// -------------------- PRODUCT INVENTORY TRANSACTION SERVICE HELPER --------------------

// Validate inventory transaction: export, import, ...
func isInventoryActionValid(action string) bool {
	var res bool = true

	switch action {
	case Sale_action:
	case Import_action:
	case Export_action:
	case Return_action:
	case Cancel_order_action:
	default:
		res = false
	}

	return res
}

// Execute inventory transaction when updated
func executeInventoryTransaction(action string, amount int64) int64 {
	// As sale and export are the actions which minus the product inventory -> turn the amount to negative
	if action == Sale_action || action == Export_action {
		amount = -amount
	}

	return amount
}

// Inverse inventory transaction when updated
func inverseInventoryTransaction(action string, amount int64) int64 {
	// As sale and export are the actions which minus the product inventory
	// When reverse these actions, we will add them back to the inventory
	// Otherwise, the rest actions as Adding quantity to inventory will be inverse to minus
	if action != Sale_action || action != Export_action {
		amount = -amount
	}

	return amount
}
