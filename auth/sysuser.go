package auth

type sysuser struct{
	username   string
	uid        int
	gid        int
	homedir    string
	loginshell string
}