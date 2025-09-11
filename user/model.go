package user

type Credentials struct {
	Email    string `json:"email" firestore:"email"`
	Password string `json:"password" firestore:"password"`
}

type User struct {
	Uid      string `json:"uid" firestore:"uid"`
	Email    string `json:"email" firestore:"email"`
	Password string `json:"password" firestore:"password"`
}
