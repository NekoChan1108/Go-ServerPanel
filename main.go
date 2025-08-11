package main

import (
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
	"log"
	"os"
)

func main() {
	// 设置客户端请求参数
	config := &ssh.ClientConfig{
		User: "ubuntu",
		Auth: []ssh.AuthMethod{
			ssh.Password("******"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 忽略主机密钥
	}

	// 连接SSH服务器
	conn, err := ssh.Dial("tcp", "*.*.*.*:22", config)
	if err != nil {
		log.Fatal("unable to connect: ", err)
	}
	defer conn.Close()

	// 创建会话
	session, err := conn.NewSession()
	if err != nil {
		log.Fatal("unable to create session: ", err)
	}
	defer session.Close()

	// 设置会话的标准输出、错误输出、标准输入
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	// 获取终端尺寸
	fd := int(os.Stdin.Fd())
	width, height, err := term.GetSize(fd)
	if err != nil {
		width, height = 80, 40
	}

	// 设置终端为原始模式
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		log.Fatal("failed to make terminal raw: ", err)
	}
	defer term.Restore(fd, oldState)

	// 请求伪终端，使用更完整的终端模式
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // 启用回显
		ssh.TTY_OP_ISPEED: 14400, // 输入速度
		ssh.TTY_OP_OSPEED: 14400, // 输出速度
	}

	if err := session.RequestPty("xterm-256color", height, width, modes); err != nil {
		log.Fatal("failed to request pty: ", err)
	}

	// 启动远程Shell
	if err := session.Shell(); err != nil {
		log.Fatal("failed to start shell: ", err)
	}

	// 阻塞直至结束会话
	if err := session.Wait(); err != nil {
		log.Fatal("exit error: ", err)
	}
}
