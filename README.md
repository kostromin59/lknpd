# LKNPD

Go-клиент для создания и аннулирования чеков, созданных в [«Мой налог» (Кабинет налогоплательщика НПД)](https://lknpd.nalog.ru).

**Важно! Клиент ещё в разработке и не тестировался. Использование до официального релиза на свой страх и риск!**

## Вход в аккаунт

Вход доступен с использованием ИНН и пароля, которые можно получить в ФНС или МВД на момент 2026 года. Дополнительно доступен способ через refresh token, а так же token. 

Token можно использовать до тех пор, пока он жив. Затем необходимо пересоздать клиент с новым токеном или любым другим способом авторизации.

Все способы можно комбинировать между собой в рамках одного клиента.

```go
// Через ИНН и пароль
lknpdClient := lknpd.New(lknpd.WithCredentials("inn", "password"))
loginResponse, err := lknpdClient.Login(context.TODO())
if err != nil {
  panic(err)
}

log.Printf("response: %+v", loginResponse)
log.Printf("token: %q", lknpdClient.Token())
log.Printf("refresh token: %q", lknpdClient.RefreshToken())

// На запросы до тех пор, пока token жив. 
// Потом необходимо пересоздать клиент с новым token, refresh token или ИНН и паролем
lknpdClient = lknpd.New(lknpd.WithToken("token"))

// Через Refresh Token при наличии
lknpdClient = lknpd.New(lknpd.WithRefreshToken("refreshToken"))
if err := lknpdClient.Refresh(context.TODO()); err != nil {
  panic(err)
}
```

## Создание и аннулирование чека

Методы для создания и аннулирования чека используют метод `RequestWithAuth`, который требует наличия хотя бы одного из токенов. Если оба токена пустые, то возвращает `lknpd.ErrUnauthorized`. Если token пустой или просрочен, то он будет обновлён. Только после этого отправляет запрос.
```go  
// Для физического лица
approvedReceiptUUID, err := lknpdClient.CreateIncome(
    context.TODO(),
    lknpd.IncomeClient{
        IncomeType: lknpd.IncomeClientTypeFromIndividual,
    },
    []lknpd.Income{
        {Name: "Услуга 1", Amount: 10, Quantity: 1},
        {Name: "Услуга 2", Amount: 3, Quantity: 2},
    }, 
    time.Now(),
)
if err != nil {
    panic(err)
}
log.Printf("approvedReceiptUUID: %q", approvedReceiptUUID)
  

// Для юридического лица
approvedReceiptUUID, err = lknpdClient.CreateIncome(
    context.TODO(),
    lknpd.IncomeClient{
        IncomeType:  lknpd.IncomeClientTypeFromLegalEntity,
        DisplayName: new("ООО \"Рога и Копыта\""),
        INN:         new("0123"),
    },
    []lknpd.Income{
        {Name: "Услуга 1", Amount: 10, Quantity: 1},
        {Name: "Услуга 2", Amount: 3, Quantity: 2},
    }, 
    time.Now(),
)
if err != nil {
    panic(err)
}
log.Printf("approvedReceiptUUID: %q", approvedReceiptUUID)  

// Аннулирование чека
if err := lknpdClient.CancelIncome(
    context.TODO(),
    approvedReceiptUUID,
    lknpd.CancelIncomeCommentMistake, // «Чек сформирован ошибочно»
); err != nil {
    panic(err)
}
```