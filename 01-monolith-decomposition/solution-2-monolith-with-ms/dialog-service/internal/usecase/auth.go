package usecase

type AuthUseCase struct {
	authRepo AuthRepository
}

func NewAuthUseCase(authRepo AuthRepository) *AuthUseCase {
	return &AuthUseCase{authRepo: authRepo}
}

func (uc *AuthUseCase) ValidateToken(token string) (string, error) {
	userID, err := uc.authRepo.GetUserIDByToken(token)
	if err != nil {
		return "", ErrUnauthorized
	}
	return userID, nil
}

// AuthRepository interface for token validation
type AuthRepository interface {
	GetUserIDByToken(token string) (string, error)
}
