package main

import (
	"MD5_FIC/md5_calc"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"
)

func getFileMD5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5_calc.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func main() {
	// 指定8888端口监听
	ln, err := net.Listen("tcp", "0.0.0.0:8888")
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}
	fmt.Println("Server started, listening on port 8888...")
	defer ln.Close()

	// 若请求失败，则继续请求
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		} else {
			fmt.Println("--------------------------------------------------")
			fmt.Println("Client connected")
		}
		go handleConnection(conn)
	}
}

// 用goroutine处理连接，避免阻塞主线程
func handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Failed to close connection:", err)
		}
	}(conn)

	// 读取文件大小
	var fileSize int64
	err := binary.Read(conn, binary.BigEndian, &fileSize)
	if err != nil {
		fmt.Println("Failed to read file size:", err)
		return
	}

	// 读取MD5哈希值
	md5Hash := make([]byte, 32)
	_, err = io.ReadFull(conn, md5Hash)
	if err != nil {
		fmt.Println("Failed to read MD5 hash:", err)
		return
	}

	// 读取文件数据
	fileData := make([]byte, fileSize)
	_, err = io.ReadFull(conn, fileData)
	if err != nil {
		fmt.Println("Failed to read file data:", err)
		return
	}

	// 获取当前时间
	now := time.Now()
	// 将时间格式化为字符串
	timeStr := now.Format("2006-01-02 15:04:05.000")
	// 计算时间字符串的哈希值
	hash := md5_calc.Sum([]byte(timeStr))
	// 将哈希值转换为字符串
	hashStr := hex.EncodeToString(hash[:])
	// 将哈希值字符串与文件名拼接
	FileName := hashStr + ".txt"
	Md5Name := "[md5]" + FileName

	// 拼接文件路径
	filePath := filepath.Join("./receive_files/")
	fileMd5path := filepath.Join("./receive_md5/")

	folderPath := filepath.Join(filePath, FileName)
	folderpathMd5 := filepath.Join(fileMd5path, Md5Name)

	// 保存文件
	err = os.WriteFile(folderPath, fileData, 0644)
	err = os.WriteFile(folderpathMd5, md5Hash, 0644)
	if err != nil {
		fmt.Printf("Failed to save file: %v\n", err)
		return
	} else {
		fmt.Println("Save file success")
	}

	// 按标准计算文件内容的MD5
	expectedMD5, err := getFileMD5(folderPath)
	if err != nil {
		fmt.Printf("Failed to calculate expectedMD5: %v\n", err)
		return
	} else {
		fmt.Println("ExpectedMD5 calculate success")
	}

	// 将文件内容的MD5转换为字符串
	receivedMD5 := string(md5Hash)
	FileData := string(fileData)

	// 打印文件内容和MD5
	fmt.Println("Client File：" + "\n" + FileData + "\n")
	fmt.Println("Received Client MD5：" + "\n" + receivedMD5 + "\n")
	fmt.Println("Expected Sever MD5：" + "\n" + expectedMD5 + "\n")

	// 标准比较模块，将客户端发送的MD5与服务端标准计算的MD5进行比较
	fmt.Println("Check result：")
	if expectedMD5 == receivedMD5 {
		// 发送验证通过的响应消息
		_, err = fmt.Fprintln(conn, "File integrity verified.")
		if err != nil {
			fmt.Println("Failed to send response:", err)
			return
		}
		fmt.Println("File integrity verified.")
	} else {
		// 发送验证未通过的响应消息
		_, err = fmt.Fprintln(conn, "File integrity not verified.")
		if err != nil {
			fmt.Println("Failed to send response:", err)
			return
		}
		fmt.Println("File integrity not verified.")
		// 删除完整性验证失败文件
		err = os.Remove(folderPath)
		if err != nil {
			fmt.Printf("Failed to remove file: %v\n", err)
			return
		} else {
			fmt.Println("Remove the error file success")
		}
		// 删除完整性验证失败MD5文件
		err = os.Remove(folderpathMd5)
		if err != nil {
			fmt.Printf("Failed to remove file: %v\n", err)
			return
		} else {
			fmt.Println("Remove the md5 file success")
		}
	}
	//结束连接
	fmt.Println("Client disconnected")
	fmt.Println("--------------------------------------------------")

}
