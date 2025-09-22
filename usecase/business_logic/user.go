package businesslogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sonit_server/constant/env"
	"sonit_server/constant/env/auth"
	mail_const "sonit_server/constant/mail_const"
	"sonit_server/constant/noti"
	page_url "sonit_server/constant/page_url"
	repo "sonit_server/data_access" // Role data access ~~ Role repository
	"sonit_server/data_access/db"
	db_server "sonit_server/data_access/db_server"
	business_logic "sonit_server/interface/business_logic"
	"sonit_server/model/dto/request"
	"sonit_server/model/dto/response"
	entity "sonit_server/model/entity"
	"sonit_server/utils"
	"strings"
	"sync"
	"time"

	data_access "sonit_server/interface/data_access"
)

type userService struct {
	logger           *log.Logger
	roleRepo         data_access.IRoleRepo
	userSecurityRepo data_access.IUserSecurityRepo
	userRepo         data_access.IUserRepo
}

func InitializeUserService(db *sql.DB, logger *log.Logger) business_logic.IUserService {
	return &userService{
		logger:           logger,
		roleRepo:         repo.InitializeRoleRepo(db, logger),
		userSecurityRepo: repo.InitializeUserSecurityRepo(db, logger),
		userRepo:         repo.InitializeUserRepo(db, logger),
	}
}

func GenerateUserService() (business_logic.IUserService, error) {
	var logger = utils.GetLogConfig()

	cnn, err := db.ConnectDB(logger, db_server.InitializePostgreSQL())

	if err != nil {
		return nil, err
	}

	user_cnn = cnn

	return InitializeUserService(cnn, logger), nil
}

var user_cnn *sql.DB

// CreateVipCode implements businesslogic.IUserService.
func (u *userService) CreateVipCode(req request.CreateVipCodeRequest, ctx context.Context) error {
	defer closeCnn(user_cnn)

	user, err := u.userRepo.GetUser(req.UserId, ctx)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	if !user.IsVip {
		return errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	if !utils.IsNumericString(req.VipCode) {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	isExist, err := u.userRepo.IsVipCodeExist(req.VipCode, ctx)
	if err != nil {
		return err
	}

	if isExist {
		return fmt.Errorf(noti.DATA_EXISTED_WARN_MSG, "Vip code", "vip code")
	}

	user.VipCode = &req.VipCode
	user.UpdatedAt = time.Now()

	return u.userRepo.UpdateUser(*user, ctx)
}

// ChangeUserStatus implements businesslogic.IUserService.
func (u *userService) ChangeUserStatus(req request.ChangeUserStatusRequest, ctx context.Context) (string, error) {
	defer closeCnn(user_cnn)

	// Empty status
	if req.Status == nil {
		return "", nil
	}

	// Validate actor and affected account
	if err := validateActorAndAffectedAccount(req.ActorId, req.UserId, u.userRepo, ctx); err != nil {
		return "", err
	}

	var status bool = *req.Status

	// Error while updating status
	if err := u.userRepo.ChangeUserStatus(req.UserId, status, ctx); err != nil {
		return "", err
	}

	// Locking account
	if !status {
		// User locking their own account
		if req.ActorId == req.UserId {
			return os.Getenv(auth.LOGIN_PAGE_URL), nil
		}
	}

	return "", nil
}

// CreateUser implements businesslogic.IUserService.
func (u *userService) CreateUser(req request.CreateUserRequest, ctx context.Context) (string, error) {
	defer closeCnn(user_cnn)

	// Admin create new account
	if req.ActorId != "" {
		// Verify admin account
		if err := validateActorAndAffectedAccount(req.ActorId, req.ActorId, u.userRepo, ctx); err != nil {
			return "", err
		}
	}

	// Email registered
	var check = isEntityExist(u.userRepo, email_type, req.Email, ctx)
	u.logger.Println(check)
	if check {
		return "", errors.New(noti.EMAIL_REGISTERED_WARN_MSG)
	}

	// Check password secure
	if !utils.IsPasswordSecure(req.Password) {
		return "", errors.New(noti.PASSWORD_NOT_SECURE_WARN_MSG)
	}

	if req.ProfileAvatar != "" {
		// if !InitializeImageService(req.ProfileAvatar).IsImageValid() {
		// 	return "", errors.New(noti.GENERIC_ERROR_WARN_MSG)
		// }
	}

	// Hash password
	hashPw, err := utils.ToHashString(req.Password, u.logger)
	if err != nil {
		return "", err
	}

	if req.RoleId == "" {
		req.RoleId = os.Getenv(env.USER_ROLE)
	}

	if req.Gender == "" {
		req.Gender = "Unknown"
	}

	// Belongs to last fail access of a new account
	var tmpTime time.Time = utils.GetPrimitiveTime()

	// Flag if account need to reset password (in case staff role creates)
	var isHaveToResetPw *bool = nil
	if req.ActorId != "" {
		var flag bool = true
		isHaveToResetPw = &flag
	}

	var fullName string = req.FullName
	if fullName == "" {
		fullName = req.Email
	}

	var userId string = utils.GenerateId()
	var curTime time.Time = time.Now()

	// Generate action token for verifying new account
	token, err := utils.GenerateActionToken(req.Email, userId, req.RoleId, u.logger)
	if err != nil {
		return "", err
	}

	if err := u.userRepo.CreateUser(entity.User{
		UserId:          userId,
		RoleId:          req.RoleId,
		FullName:        fullName,
		Email:           req.Email,
		Password:        hashPw,
		ProfileAvatar:   req.ProfileAvatar,
		Gender:          req.Gender,
		IsActive:        true,
		IsActivated:     false,
		IsHaveToResetPw: isHaveToResetPw,
		CreatedAt:       curTime,
		UpdatedAt:       curTime,
	}, ctx); err != nil {
		return "", err
	}

	// Create account security

	if err := u.userSecurityRepo.CreateUserSecurity(entity.UserSecurity{
		UserId:      userId,
		ActionToken: &token,
		FailAccess:  0,
		LastFail:    &tmpTime,
	}, ctx); err != nil {
		return "", err
	}

	// Send confirmation mail
	if err := utils.SendMail(request.SendMailRequest{
		Body: request.MailBody{ // Mail body
			Email:    req.Email,
			Password: req.Password,
			Subject:  noti.REGISTRATION_ACCOUNT_MAIL_SUBJECT,
			Url: utils.ToCombinedString([]string{ // Call back url when guest clicks to the confirmation, it will call back to the api endpoint which generate here to verify and finish the registration process
				os.Getenv(auth.PROCESS_ACTION_URL),
				token,
				userId,
				activateType,
			},
				mailSepChar),
		},

		TemplatePath: mail_const.ACCOUNT_REGISTRATION_MAIL_TEMPLATE, // Template path

		Logger: u.logger, // Logger
	}); err != nil {
		return "", err
	}

	var msg string = "success"
	if req.ActorId == "" {
		msg = noti.REGISTRATION_ACCOUNT_MESSAGE
	}

	return msg, nil
}

// GetAllUsers implements businesslogic.IUserService.
func (u *userService) GetAllUsers(pageNumber int, ctx context.Context) (response.PaginationDataResponse, error) {
	if pageNumber < 1 {
		pageNumber = 1
	}

	defer closeCnn(user_cnn)

	data, pages, err := u.userRepo.GetAllUsers(pageNumber, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	return response.PaginationDataResponse{
		Data:       data,
		PageNumber: pageNumber,
		TotalPages: pages,
	}, nil
}

// GetUser implements businesslogic.IUserService.
func (u *userService) GetUser(id string, ctx context.Context) (*entity.User, error) {
	defer closeCnn(user_cnn)
	return u.userRepo.GetUser(id, ctx)
}

// GetUsersByRole implements businesslogic.IUserService.
func (u *userService) GetUsersByRole(pageNumber int, role string, ctx context.Context) (response.PaginationDataResponse, error) {
	var data *[]entity.User
	var pages int
	var err error
	defer closeCnn(user_cnn)

	if role == "" {
		data, pages, err = u.userRepo.GetAllUsers(pageNumber, ctx)
	} else {
		data, pages, err = u.userRepo.GetUsersByRole(pageNumber, role, ctx)
	}

	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	return response.PaginationDataResponse{
		Data:       data,
		PageNumber: pageNumber,
		TotalPages: pages,
	}, nil
}

// GetUsersByStatus implements businesslogic.IUserService.
func (u *userService) GetUsersByStatus(pageNumber int, status *bool, ctx context.Context) (response.PaginationDataResponse, error) {
	var data *[]entity.User
	var pages int
	var err error
	defer closeCnn(user_cnn)

	if status == nil {
		data, pages, err = u.userRepo.GetAllUsers(pageNumber, ctx)
	} else {
		data, pages, err = u.userRepo.GetUsersByStatus(pageNumber, *status, ctx)
	}

	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	return response.PaginationDataResponse{
		Data:       data,
		PageNumber: pageNumber,
		TotalPages: pages,
	}, nil
}

// GetUsersByKeyword implements businesslogic.IUserService.
func (u *userService) GetUsersByKeyword(pageNumber int, keyword string, ctx context.Context) (response.PaginationDataResponse, error) {
	var data *[]entity.User
	var pages int
	var err error
	defer closeCnn(user_cnn)

	if keyword == "" {
		data, pages, err = u.userRepo.GetAllUsers(pageNumber, ctx)
	} else {
		data, pages, err = u.userRepo.GetUsersByRole(pageNumber, keyword, ctx)
	}

	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	return response.PaginationDataResponse{
		Data:       data,
		PageNumber: pageNumber,
		TotalPages: pages,
	}, nil
}

// Login implements businesslogic.IUserService.
func (u *userService) Login(req request.LoginRequest, ctx context.Context) (string, string, error) {
	var account entity.User
	defer closeCnn(user_cnn)

	// Verify account
	if err := verifyAccount(req.Email, email_validate, &account, u.userRepo, ctx); err != nil {
		return "", "", err
	}

	// Incorrect password
	if !utils.IsHashStringMatched(req.Password, account.Password) {
		processWrongCredentialsCase(account.UserId, u.userSecurityRepo, ctx)
		return "", "", errors.New(noti.WRONG_CREDENTIALS_WARN_MSG)
	}

	return processCorrectCredentialsCase(account, u.logger, u.userSecurityRepo, ctx)
}

// Logout implements businesslogic.IUserService.
func (u *userService) Logout(id string, ctx context.Context) error {
	defer closeCnn(user_cnn)
	return u.userSecurityRepo.Logout(id, ctx)
}

// ResetPassword implements businesslogic.IUserService.
func (u *userService) ResetPassword(newPass string, reNewPass string, token string, ctx context.Context) (string, error) {
	defer closeCnn(user_cnn)

	accountId, _, exp, err := utils.ExtractDataFromToken(token, u.logger)
	if err != nil {
		return os.Getenv(auth.LOGIN_PAGE_URL), err
	}

	var account entity.User
	if err := verifyAccount(accountId, id_validate, &account, u.userRepo, ctx); err != nil { // Verify account
		return os.Getenv(auth.LOGIN_PAGE_URL), err
	}

	// User state doesn't have to reset password
	if account.IsHaveToResetPw == nil {
		return os.Getenv(auth.LOGIN_PAGE_URL), errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	// Expired
	if utils.IsActionExpired(exp) {
		return os.Getenv(auth.LOGIN_PAGE_URL), errors.New("")
	}

	// Retrieve account security info
	usc, err := u.userSecurityRepo.GetUserSecurity(accountId, ctx)
	if err != nil {
		return os.Getenv(auth.LOGIN_PAGE_URL), err
	}

	// Not matched token
	if *usc.ActionToken != token {
		return os.Getenv(auth.LOGIN_PAGE_URL), errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	// Passwords not matched
	if newPass != reNewPass {
		return utils.ToCombinedString([]string{
			page_url.PROCESS_ENDPOINT,
			token,
		}, sepChar), errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	// Hash new password
	hashPw, err := utils.ToHashString(reNewPass, u.logger)
	if err != nil {
		return os.Getenv(auth.LOGIN_PAGE_URL), err
	}

	// Assign new data
	account.Password = hashPw
	account.IsHaveToResetPw = nil

	usc.ActionToken = nil

	// Process update
	var capturedErr error

	_, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(2)

	// Update user data
	go func() {
		defer wg.Done()

		if err := u.userRepo.UpdateUser(account, ctx); err != nil {
			mu.Lock()

			if capturedErr == nil {
				capturedErr = err
				cancel()
			}

			mu.Unlock()
		}
	}()

	// Update user security
	go func() {
		defer wg.Done()

		if err := u.userSecurityRepo.EditUserSecurity(*usc, ctx); err != nil {
			mu.Lock()

			if capturedErr == nil {
				capturedErr = err
				cancel()
			}

			mu.Unlock()
		}
	}()

	wg.Wait()

	return os.Getenv(os.Getenv(auth.LOGIN_PAGE_URL)), capturedErr
}

// UpdateUser implements businesslogic.IUserService.
func (u *userService) UpdateUser(req request.UpdateUserRequest, ctx context.Context) (string, error) {
	defer closeCnn(user_cnn)

	// Verify actor and affected account
	if err := validateActorAndAffectedAccount(req.ActorId, req.UserId, u.userRepo, ctx); err != nil {
		return "", err
	}

	// Retrieve original account info
	account, err := u.userRepo.GetUser(req.UserId, ctx)
	if err != nil {
		return "", err
	}

	var newEmail string
	var isHaveToSendMail bool = false
	if req.Email != "" { // New email inputted
		/// Normalize email
		newEmail = strings.TrimSpace(req.Email)
		if newEmail != account.Email {
			tmpAcc, err := u.userRepo.GetUserByEmail(newEmail, ctx)
			if err != nil {
				return "", err
			}

			// Email registered
			if tmpAcc != nil {
				return "", errors.New(noti.EMAIL_REGISTERED_WARN_MSG)
			}

			isHaveToSendMail = true
		}
	}

	if req.Password != "" {
		if utils.IsPasswordSecure(req.Password) {
			// Hash password
			hashedPassword, err := utils.ToHashString(req.Password, u.logger)
			// No error
			if err == nil {
				// Assign new hashed password
				account.Password = hashedPassword
			}
		} else {
			return "", errors.New(noti.PASSWORD_NOT_SECURE_WARN_MSG)
		}
	}

	if req.RoleId != "" {
		account.RoleId = req.RoleId
	}

	if req.FullName != "" {
		account.FullName = req.FullName
	}

	if req.ProfileAvatar != "" {
		account.ProfileAvatar = req.ProfileAvatar
	}

	if req.Gender != "" {
		account.Gender = req.Gender
	}

	// Must check field isVip
	// Implement later

	// Update new info to database
	if err := u.userRepo.UpdateUser(*account, ctx); err != nil {
		return "", err
	}

	var res string

	// If must be sended verification mail
	if isHaveToSendMail {
		// Generate token
		token, err := utils.GenerateActionToken(account.Email, req.UserId, req.RoleId, u.logger)
		if err != nil {
			return "", err
		}

		// Get user security data to update token for new action - update new email
		usc, err := u.userSecurityRepo.GetUserSecurity(req.UserId, ctx)
		if err != nil {
			return "", err
		}

		// Update user security to db
		usc.ActionToken = &token
		if err := u.userSecurityRepo.EditUserSecurity(*usc, ctx); err != nil {
			return res, err
		}

		// Send verification mail
		if utils.SendMail(request.SendMailRequest{
			Body: request.MailBody{
				Email:   newEmail,
				Subject: noti.UPDATE_EMAIL_MAIL_SUBJECT,
				Url: utils.ToCombinedString([]string{
					os.Getenv(auth.PROCESS_ACTION_URL),
					token,
					req.UserId,
					updateProfileType,
					newEmail,
				}, mailSepChar),
			},

			TemplatePath: mail_const.UPDATE_MAIL_TEMPLATE,

			Logger: u.logger,
		}) != nil {
			return "", errors.New("We have updated your other information. " + noti.GENERATE_MAIL_WARN_MSG)
		}

		res = noti.UPDATE_EMAIL_MESSAGE
	} else {
		res = "Success"
	}

	return res, nil
}

// VerifyAction implements businesslogic.IUserService.
func (u *userService) VerifyAction(rawToken string, ctx context.Context) (string, error) {
	var errRes error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	var cmps []string = utils.ToSliceString(rawToken, mailSepChar)
	defer closeCnn(user_cnn)

	if len(cmps) < min_length_verify_mail_component_combine { // Min length of combination of information in a call back url
		return "", errRes
	}

	var token string = cmps[0]
	var id string = cmps[1]
	var actionType string = cmps[2]

	usc, err := u.userSecurityRepo.GetUserSecurity(id, ctx)
	if err != nil {
		return "", errRes
	}

	if *usc.ActionToken != token {
		return "", errRes
	}

	// Extract data from token
	extractId, _, exp, err := utils.ExtractDataFromToken(token, u.logger)
	if err != nil {
		return "", err
	}

	if extractId != id {
		return "", errRes
	}

	// Expired
	if utils.IsActionExpired(exp) {
		return "", errors.New("Token expired")
	}

	// Invalid action type
	if actionType != activateType && actionType != resetPassType && actionType != updateProfileType && actionType != verifyType {
		return "", errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var res string

	if actionType == activateType || actionType == updateProfileType {
		user, err := u.userRepo.GetUser(id, ctx)
		if err != nil {
			return res, err
		}

		if actionType == activateType {
			user.IsActivated = true
		} else {
			user.Email = cmps[len(cmps)-1]
		}

		if err := u.userRepo.UpdateUser(*user, ctx); err != nil {
			return res, err
		}

		res = os.Getenv(auth.LOGIN_PAGE_URL)
	} else {
		// Setup prepare to reset password
		token, err := utils.GenerateActionToken("", usc.UserId, "", u.logger)
		if err != nil {
			return "", err
		}

		*usc.ActionToken = token

		res = utils.ToCombinedString([]string{
			page_url.PROCESS_ENDPOINT,
			token,
		}, sepChar)
	}

	return res, nil
}

// RefreshToken implements businesslogic.IUserService.
func (u *userService) RefreshToken(req request.RefreshTokenRequest, ctx context.Context) (string, error) {
	defer closeCnn(user_cnn)

	// Extract data from token
	id, _, exp, err := utils.ExtractDataFromToken(req.RefreshToken, u.logger)
	if err != nil {
		return "", err
	}

	// Get user security
	sc, err := u.userSecurityRepo.GetUserSecurity(id, ctx)
	if err != nil {
		return "", err
	}

	if *sc.RefreshToken != req.RefreshToken {
		return "", errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	if utils.IsActionExpired(exp) {
		return "", errors.New(noti.TOKEN_EXPIRED_MESSAGE)
	}

	token, _, err := utils.GenerateTokens("", id, "", u.logger)
	if err != nil {
		return "", err
	}

	sc.AccessToken = &token

	return token, u.userSecurityRepo.EditUserSecurity(*sc, ctx)
}
