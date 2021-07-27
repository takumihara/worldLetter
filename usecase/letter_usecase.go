package usecase

import "github.com/tacomea/worldLetter/domain"

type letterUsecase struct {
	letterRepo domain.LetterRepository
}

func NewLetterUsecase(lr domain.LetterRepository) domain.LetterUseCase {
	return &letterUsecase{
		letterRepo: lr,
	}
}

func (u *letterUsecase) Create(letter domain.Letter) error {
	return u.letterRepo.Create(letter)
}

func (u *letterUsecase) Read(id string) (domain.Letter, error) {
	return u.letterRepo.Read(id)
}

func (u *letterUsecase) Update(letter domain.Letter) error {
	return u.letterRepo.Update(letter)
}

func (u *letterUsecase) Delete(id string) error {
	return u.letterRepo.Delete(id)
}

func (u *letterUsecase) GetAll() ([]domain.Letter, error) {
	return u.letterRepo.GetAll()
}

func (u *letterUsecase) GetFirstUnsendLetter(authorID string) (domain.Letter, error) {
	return u.letterRepo.GetFirstUnsendLetter(authorID)
}

func (u *letterUsecase) GetLettersByAuthorID(authorID string) (string, error) {
	return u.letterRepo.GetLettersByAuthorID(authorID)
}

func (u *letterUsecase) GetLettersByReceiverID(receiverID string) (string, error) {
	return u.letterRepo.GetLettersByReceiverID(receiverID)
}