package user

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/answerdev/answer/pkg/token"
	"github.com/segmentfault/pacman/log"
	"net/http"
	"time"

	"github.com/answerdev/answer/internal/base/data"
	"github.com/answerdev/answer/internal/base/reason"
	"github.com/answerdev/answer/internal/entity"
	"github.com/answerdev/answer/internal/service/config"
	usercommon "github.com/answerdev/answer/internal/service/user_common"
	"github.com/segmentfault/pacman/errors"
)

// userRepo user repository
type userRepo struct {
	data       *data.Data
	configRepo config.ConfigRepo
}

// NewUserRepo new repository
func NewUserRepo(data *data.Data, configRepo config.ConfigRepo) usercommon.UserRepo {
	return &userRepo{
		data:       data,
		configRepo: configRepo,
	}
}

// AddUser add user
func (ur *userRepo) AddUser(ctx context.Context, user *entity.User) (err error) {
	_, err = ur.data.DB.Insert(user)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// IncreaseAnswerCount increase answer count
func (ur *userRepo) IncreaseAnswerCount(ctx context.Context, userID string, amount int) (err error) {
	user := &entity.User{}
	_, err = ur.data.DB.Where("id = ?", userID).Incr("answer_count", amount).Update(user)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

// IncreaseQuestionCount increase question count
func (ur *userRepo) IncreaseQuestionCount(ctx context.Context, userID string, amount int) (err error) {
	user := &entity.User{}
	_, err = ur.data.DB.Where("id = ?", userID).Incr("question_count", amount).Update(user)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

// UpdateLastLoginDate update last login date
func (ur *userRepo) UpdateLastLoginDate(ctx context.Context, userID string) (err error) {
	user := &entity.User{LastLoginDate: time.Now()}
	_, err = ur.data.DB.Where("id = ?", userID).Cols("last_login_date").Update(user)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

// UpdateEmailStatus update email status
func (ur *userRepo) UpdateEmailStatus(ctx context.Context, userID string, emailStatus int) error {
	cond := &entity.User{MailStatus: emailStatus}
	_, err := ur.data.DB.Where("id = ?", userID).Cols("mail_status").Update(cond)
	if err != nil {
		return err
	}
	return nil
}

// UpdateNoticeStatus update notice status
func (ur *userRepo) UpdateNoticeStatus(ctx context.Context, userID string, noticeStatus int) error {
	cond := &entity.User{NoticeStatus: noticeStatus}
	_, err := ur.data.DB.Where("id = ?", userID).Cols("notice_status").Update(cond)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

func (ur *userRepo) UpdatePass(ctx context.Context, userID, pass string) error {
	_, err := ur.data.DB.Where("id = ?", userID).Cols("pass").Update(&entity.User{Pass: pass})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

func (ur *userRepo) UpdateEmail(ctx context.Context, userID, email string) (err error) {
	_, err = ur.data.DB.Where("id = ?", userID).Update(&entity.User{EMail: email})
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// UpdateInfo update user info
func (ur *userRepo) UpdateInfo(ctx context.Context, userInfo *entity.User) (err error) {
	_, err = ur.data.DB.Where("id = ?", userInfo.ID).
		Cols("username", "display_name", "avatar", "bio", "bio_html", "website", "location").Update(userInfo)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// GetByUserID get user info by user id
func (ur *userRepo) GetByUserID(ctx context.Context, userID string) (userInfo *entity.User, exist bool, err error) {
	userInfo = &entity.User{}
	exist, err = ur.data.DB.Where("id = ?", userID).Get(userInfo)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

func (ur *userRepo) BatchGetByID(ctx context.Context, ids []string) ([]*entity.User, error) {
	list := make([]*entity.User, 0)
	err := ur.data.DB.In("id", ids).Find(&list)
	if err != nil {
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return list, nil
}

// GetByUsername get user by username
func (ur *userRepo) GetByUsername(ctx context.Context, username string) (userInfo *entity.User, exist bool, err error) {
	userInfo = &entity.User{}
	exist, err = ur.data.DB.Where("username = ?", username).Get(userInfo)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// GetByEmail get user by email
func (ur *userRepo) GetByEmail(ctx context.Context, email string) (userInfo *entity.User, exist bool, err error) {
	userInfo = &entity.User{}
	exist, err = ur.data.DB.Where("e_mail = ?", email).Get(userInfo)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

type SamLogin struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	} `json:"data"`
}

type SamProfile struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		ID    int64  `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"data"`
}

// GetByUserID get user info by user id
func (ur *userRepo) GetBySamID(ctx context.Context, samID int64) (userInfo *entity.User, exist bool, err error) {
	userInfo = &entity.User{}
	exist, err = ur.data.DB.Where("sam_id = ?", samID).Get(userInfo)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// GetBySam get user by Sam
func (ur *userRepo) GetBySam(ctx context.Context, email string, pass string) (userInfo *entity.User, exist bool, err error) {
	values := map[string]string{"email": email, "password": pass}
	jsonData, err := json.Marshal(values)
	if err != nil {
		err = errors.InternalServer(reason.RequestFormatError).WithError(err).WithStack()
		return nil, false, err
	}
	samUri := ur.data.Sam.Uri
	//samUri, err := ur.configRepo.GetString("sam.uri")
	if err != nil {
		err = errors.InternalServer(reason.SamUriNotFound).WithError(err).WithStack()
		return nil, false, err
	}
	resp, err := http.Post(samUri+"/tess/login", "application/json",
		bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	var samLogin SamLogin
	err = json.NewDecoder(resp.Body).Decode(&samLogin)
	if err != nil {
		return nil, false, err
	}
	if samLogin.Code != 0 {
		err = errors.InternalServer(reason.SamReqError).WithError(err).WithStack()
		return nil, false, err
	}

	var bearer = "Bearer " + samLogin.Data.AccessToken
	req, err := http.NewRequest("GET", samUri+"/me", nil)
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	var samProfile SamProfile
	err = json.NewDecoder(resp.Body).Decode(&samProfile)
	if err != nil {
		return nil, false, err
	}

	userInfo, exist, err = ur.GetBySamID(ctx, samProfile.Data.ID)
	if err != nil || samProfile.Code != 0 {
		return nil, false, err
	}

	if exist {
		userInfo.AccessToken = samLogin.Data.AccessToken
		userInfo.ExpiredAt = time.Now().Add(time.Second * time.Duration(samLogin.Data.ExpiresIn))
		err = ur.UpdateSamLogin(ctx, userInfo)
		if err != nil {
			log.Error("UpdateSamLogin", err.Error())
		}
	} else {
		userInfo, err = ur.AddSamUser(ctx, samLogin, samProfile)
		if err != nil {
			return nil, false, err
		}
	}
	return userInfo, true, nil
}

// UpdateSamLogin update sam token data
func (ur *userRepo) UpdateSamLogin(ctx context.Context, userInfo *entity.User) (err error) {
	_, err = ur.data.DB.Where("id = ?", userInfo.ID).Cols("access_token", "expired_at").Update(userInfo)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

func (ur *userRepo) AddSamUser(ctx context.Context, samLogin SamLogin, samProfile SamProfile) (userInfo *entity.User, err error) {
	userInfo = &entity.User{}
	userInfo.SamId = samProfile.Data.ID
	userInfo.EMail = samProfile.Data.Email
	userInfo.AccessToken = samLogin.Data.AccessToken
	userInfo.ExpiredAt = time.Now().Add(time.Second * time.Duration(samLogin.Data.ExpiresIn))

	userInfo.Pass = token.GenerateToken()
	userInfo.Username = samProfile.Data.Name
	userInfo.DisplayName = samProfile.Data.Name
	userInfo.MailStatus = entity.EmailStatusAvailable
	userInfo.Status = entity.UserStatusAvailable

	err = ur.AddUser(ctx, userInfo)
	if err != nil {
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return userInfo, nil
}
