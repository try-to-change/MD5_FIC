package main

import (
	"MD5_FIC/md5_calc"
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"io"
	"net"
	"os"
	"time"
)

// 计算文件的MD5 hash
func getFileMD5(path string) (string, error) {
	//打开指定路径文件
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	//计算MD5
	hash := md5_calc.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func main() {
	// 创建GUI应用
	a := app.New()
	win := a.NewWindow("Client")
	win.CenterOnScreen()

	// 选择文件按钮，选择文件后将文件路径显示在filePathEntry中
	filePathEntry := widget.NewLabel("File path...")
	browseButton := widget.NewButton("Browse", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				filePathEntry.SetText(reader.URI().Path())
			}
		}, win)
	})

	// 计算MD5按钮，计算MD5后将MD5值显示在md5Label中
	md5Label := widget.NewLabel("MD5 hash...")
	calculateButton := widget.NewButton("Calculate MD5", func() {
		// 获取文件路径并进行非空判断
		filePath := filePathEntry.Text
		if filePath == "" {
			return
		}
		// 计算MD5
		md5Hash, err := getFileMD5(filePath)
		if err != nil {
			fmt.Printf("Failed to calculate MD5: %v\n", err)
			return
		}

		md5Label.SetText(md5Hash)
	})

	// 发送按钮处理客户端与服务器的连接
	result := widget.NewLabel("Result...")
	status := widget.NewLabel("Server response...")
	sendButton := widget.NewButton("Send", func() {
		filePath := filePathEntry.Text
		if filePath == "" {
			return
		}

		md5Hash := md5Label.Text

		// 指定连接的服务器地址
		conn, err := net.Dial("tcp", "0.0.0.0:8888")
		if err != nil {
			result.SetText("Failed to connect to server" + "\n" + err.Error() + "\n")
			return
		}

		// 读取文件数据
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			result.SetText("Failed to read file" + "\n" + err.Error() + "\n")
			return
		}

		// 发送文件大小
		fileSize := int64(len(fileData))
		err = binary.Write(conn, binary.BigEndian, fileSize)
		if err != nil {
			result.SetText("Failed to send file size" + "\n" + err.Error() + "\n")
			return
		}

		// 发送MD5 hash
		_, err = conn.Write([]byte(md5Hash))
		if err != nil {
			result.SetText("Failed to send MD5 hash" + "\n" + err.Error() + "\n")
			return
		}

		// 发送文件数据
		_, err = conn.Write(fileData)
		if err != nil {
			result.SetText("Failed to send file data" + "\n" + err.Error() + "\n")
			return
		}

		result.SetText("Server received file successfully" + "\n")

		// 设置读取超时时间
		err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			status.SetText("Failed to set read deadline" + "\n" + err.Error() + "\n")
			return
		}
		// 读取服务器返回消息
		reader := bufio.NewReader(conn)
		statusMsg, err := reader.ReadString('\n')
		if err != nil {
			status.SetText("Failed to read server response" + "\n" + err.Error() + "\n")
			return
		}
		// 显示服务器返回消息
		status.SetText("Server response: " + "\n" + statusMsg + "\n")

	})

	//ui界面布局
	content := container.NewVBox(
		browseButton,
		filePathEntry,
		calculateButton,
		md5Label,
		sendButton,
		result,
		status,
	)

	win.SetContent(content)
	win.Resize(fyne.NewSize(1200, 800))
	win.ShowAndRun()
}
