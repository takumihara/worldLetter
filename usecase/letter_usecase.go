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
	err := u.letterRepo.Create(letter)
	if err != nil {
		return err
	}
	return nil
}

func (u *letterUsecase) Read(id string) (domain.Letter, error) {
	user, err := u.letterRepo.Read(id)
	if err != nil {
		return domain.Letter{}, err
	}
	return user, nil
}

func (u *letterUsecase) Update(letter domain.Letter) error {
	err := u.letterRepo.Update(letter)
	if err != nil {
		return err
	}
	return nil
}

func (u *letterUsecase) Delete(id string) error {
	err := u.letterRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (u *letterUsecase) GetFirstUnsendLetter(authorID string) (domain.Letter, error) {
	letter, err := u.letterRepo.GetFirstUnsendLetter(authorID)
	if err != nil {
		return domain.Letter{}, err
	}
	return letter, nil
}
