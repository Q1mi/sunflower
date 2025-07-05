package impl

import (
	"backend/internal/consts"
	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/model/entity"
	"backend/internal/service/user"
	"backend/utility/injection"
	"context"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sony/sonyflake/v2"
)

const (
	defaultAvatar = "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Mnx8dXNlciUyMHByb2ZpbGV8ZW58MHx8MHx8fDA=&auto=format&fit=crop&w=50&q=60"
)

// 定义一个结构体，实现 UserInfoService 接口
type UserInfo struct {
	snowflake *sonyflake.Sonyflake
}

func New() user.UserInfoService {
	return &UserInfo{
		snowflake: injection.MustInvoke[*sonyflake.Sonyflake](),
	}
}

// encryptPassword 加密密码
func (u *UserInfo) encryptPassword(password string) string {
	return gmd5.MustEncryptString(password)
}

// Create 创建用户
func (u *UserInfo) Create(ctx context.Context, input *model.CreateUserInput) (*model.CreateUserOutput, error) {
	// 1. 判断用户名是否已经存在（根据用户名查重）
	exist, err := dao.Userinfo.Ctx(ctx).Where(dao.Userinfo.Columns().Username, input.Username).Exist()
	if err != nil {
		g.Log().Errorf(ctx, "查询用户是否存在失败：%v", err)
		return nil, err
	}
	if exist {
		return nil, gerror.New("用户已存在")
	}
	// 2. 生成用户ID
	userId, err := u.snowflake.NextID()
	if err != nil {
		g.Log().Errorf(ctx, "生成用户ID失败：%v", err)
		return nil, gerror.Wrapf(err, "生成用户ID失败")
	}
	// 3. 创建用户
	// 创建用户，入库
	newUserInfo := entity.Userinfo{
		UserId:   uint64(userId), // 使用雪花算法生成唯一ID
		Username: input.Username,
		Password: u.encryptPassword(input.Password), // 是不是需要对用户输入的密码进行加密
		Email:    input.Email,
		Avatar:   defaultAvatar, // 简化注册流程，一般使用默认头像，后续支持用户在个人中心上传头像
	}
	_, err = dao.Userinfo.Ctx(ctx).InsertAndGetId(newUserInfo)
	if err != nil {
		g.Log().Errorf(ctx, "创建用户失败: %v", err)
		return nil, gerror.Wrapf(err, "创建用户失败")
	}
	// 4. 返回结果
	return &model.CreateUserOutput{
		UserId:   uint64(userId),
		Username: input.Username,
	}, nil
}

// Login 登录
func (u *UserInfo) Login(ctx context.Context, input *model.LoginInput) (*model.LoginOutput, error) {
	// 拿用户输入的用户名和密码，去数据库查询
	var user entity.Userinfo
	err := dao.Userinfo.Ctx(ctx).
		Where(dao.Userinfo.Columns().Username, input.Username).
		Where(dao.Userinfo.Columns().Password, u.encryptPassword(input.Password)).
		Scan(&user)
	if err != nil {
		g.Log().Errorf(ctx, "查询用户失败：%v", err)
		return nil, gerror.Wrapf(err, "查询用户失败")
	}
	// 生成 JWT Token
	tokenObj, err := genJwtByUserInfo(ctx, user.UserId, user.Username)
	if err != nil {
		g.Log().Errorf(ctx, "生成 JWT Token 失败：%v", err)
		return nil, gerror.Wrapf(err, "生成 JWT Token 失败")
	}

	// 返回结果
	return &model.LoginOutput{
		AccessToken:  tokenObj.AccessToken,
		RefreshToken: tokenObj.RefreshToken,
	}, nil
}

// genJwtByUserInfo 根据用户信息生成 JWT
func genJwtByUserInfo(ctx context.Context, userId uint64, username string) (*model.TokenOutput, error) {
	// 生成 access Token
	claims := &model.JWTClaims{
		UserId:   userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "liwenzhou.com",
			Subject:   "sunflower",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(consts.JWTTokenExpireSeconds * time.Second)), // 设置过期时间为1天
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedAccessToken, err := accessToken.SignedString([]byte(consts.JWTAccessTokenSecret))
	if err != nil {
		g.Log().Errorf(ctx, "生成 JWT Token 失败：%v", err)
		return nil, err
	}
	// 生成 refresh Token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.JWTClaims{
		UserId:   userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "liwenzhou.com",
			Subject:   "sunflower_refresh",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(consts.JWTRefreshExpireSeconds * time.Second)), // 设置过期时间为1周
		},
	})
	signedRefreshToken, err := refreshToken.SignedString([]byte(consts.JWTRefreshTokenSecret))
	if err != nil {
		g.Log().Errorf(ctx, "生成 JWT Token 失败：%v", err)
		return nil, err
	}
	return &model.TokenOutput{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
}

func (u *UserInfo) GetInfo(ctx context.Context, userId string) (*entity.Userinfo, error) {
	var user entity.Userinfo
	if err := dao.Userinfo.Ctx(ctx).
		Where(dao.Userinfo.Columns().UserId, userId).
		Scan(&user); err != nil {
		g.Log().Errorf(ctx, "查询用户失败：%v", err)
		return nil, err
	}
	return &user, nil
}

// RefreshToken 刷新 Token
func (u *UserInfo) RefreshToken(ctx context.Context, refreshToken string) (res *model.TokenOutput, err error) {
	// 1. 解析 refreshToken, 拿到 userid
	var claim model.JWTClaims
	token, err := jwt.ParseWithClaims(refreshToken, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(consts.JWTRefreshTokenSecret), nil
	})
	if err != nil || !token.Valid {
		g.Log().Errorf(ctx, "refresh token: %v, err:%+v", token, err)
		return nil, gerror.New("refresh token 无效")
	}
	// 2. 根据 userid 获取用户信息
	userInfo, err := u.GetInfo(ctx, strconv.FormatUint(claim.UserId, 10))
	if err != nil {
		g.Log().Errorf(ctx, "根据 userid 获取用户信息失败：%v", err)
		return nil, gerror.Wrapf(err, "根据 userid 获取用户信息失败")
	}
	// 3. 生成新的 accessToken 和 refreshToken
	// 4. 返回新的 token
	return genJwtByUserInfo(ctx, userInfo.UserId, userInfo.Username)
}
