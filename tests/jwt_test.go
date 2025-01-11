package server_test

import (
	"fmt"
	"testing"

	"github.com/antfley/asyncapi/api/server"
	"github.com/antfley/asyncapi/config"
	"github.com/stretchr/testify/require"
)

func TestJWTManager(t *testing.T) {
	cfg, err := config.New()
	require.NoError(t, err)
	jwtManager := server.NewJWTManager(cfg)
	userID := uint(1)
	tokenPair, err := jwtManager.GenerateTokenPair(userID)
	require.NoError(t, err)

	require.True(t, jwtManager.IsAccessToken(tokenPair.AccessToken))
	require.False(t, jwtManager.IsAccessToken(tokenPair.RefreshToken))

	accessTokenSubject, err := tokenPair.AccessToken.Claims.GetSubject()
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("%d", userID), accessTokenSubject)

	accessTokenIssuer, err := tokenPair.AccessToken.Claims.GetIssuer()
	require.NoError(t, err)
	require.Equal(t, "http://localhost:"+cfg.PORT, accessTokenIssuer)

	refreshTokenSubject, err := tokenPair.RefreshToken.Claims.GetSubject()
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("%d", userID), refreshTokenSubject)

	refreshTokenIssuer, err := tokenPair.RefreshToken.Claims.GetIssuer()
	require.NoError(t, err)
	require.Equal(t, "http://localhost:"+cfg.PORT, refreshTokenIssuer)
	parsedAccessToken, err := jwtManager.Parse(tokenPair.AccessToken.Raw)
	require.NoError(t, err)
	require.Equal(t, tokenPair.AccessToken.Raw, parsedAccessToken.Raw)
	parsedRefreshToken, err := jwtManager.Parse(tokenPair.RefreshToken.Raw)
	require.NoError(t, err)
	require.Equal(t, tokenPair.RefreshToken.Raw, parsedRefreshToken.Raw)
}
