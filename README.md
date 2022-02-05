# 世界手紙交換アプリ

## アプリの概要

1通手紙を書くことで、代わりに1通手紙を受け取れる世界中の人との手紙交換サービスです。

## 使用言語

```
Go 1.15
```

## 機能

- Signup, Login, Logoutなどのユーザー認証
- 手紙の送受信
- 送受信した手紙の確認

## 注力した点

- クリーンアーキテクチャを元にレイヤー分割を行った
- Interfaceを使うことでインメモリ (`sync.Map`) とDBを環境変数で切り替えられるようにし、テストなどを容易に行えるようにした。
```shell
$ DB_DIALECT="postgres" go test ./repository
$ DB_DIALECT="map" go test ./repository
```
```go
// domain/user.go

type UserRepository interface {
	Create(user User) error
	Read(email string) (User, error)
	Update(user User) error
	Delete(email string) error
}
```

```go
// repository/user_postgres.go

func (u *userRepositoryPG) Create(user domain.User) error {
	result := u.db.Create(user)
	if result.Error != nil {
		return errors.New("unexpected error while creating user")
	}
	return nil
}
```

```go
// repository/user_map.go

func (u *userRepository) Create(user domain.User) error {
	u.m.Store(user.Email, user)
	_, ok := u.m.Load(user.Email)
	if !ok {
		return errors.New("user not created")
	}
	return nil
}
```

## 環境構築の手順

```bash
git clone https://github.com/tacomeet/worldLetter
cd worldLetter
docker compose up
```

## デモサイト

こちらのサイトにデプロイされています。
https://world-letter.herokuapp.com/


下記テストアカウントを用いてログイン可能です！

```
Email: test@exmaple.com
Passcode: passcode
```
