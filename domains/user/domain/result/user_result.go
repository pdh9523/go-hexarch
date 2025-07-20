package result

type CheckNicknameResult struct {
	Nickname string
	Status   NicknameStatus
	Reason   string
}

type NicknameStatus int

const (
	NicknameAvailable NicknameStatus = iota
	NicknamePolicyViolated
	NicknameDuplicated
)
