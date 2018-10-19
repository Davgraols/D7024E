package main

import (
	"fmt"
	"time"
)

type FileStore struct {
	fileData map[KademliaID]File
}

func NewFileStore() FileStore {
	fs := FileStore{
		fileData: make(map[KademliaID]File),
	}

	return fs
}

func (files *FileStore) StoreFile(fileContent []byte, owner *Contact) {
	fileId := NewRandomHash(string(fileContent))
	FileLock.Lock()
	file, fileExist := files.fileData[*fileId]
	FileLock.Unlock()
	if fileExist {
		if FileStoreDebug {
			fmt.Println("File Exists. sending true to nodeRepub")
		}
		file.nodeRepub <- true
		fmt.Println("sent true")
	} else {
		file = NewFile(fileContent, owner)
		files.fileData[*fileId] = file
		go files.republish(file)
		if !owner.ID.Equals(RT.me.ID) {
			go files.FileHeartbeat(file)
		}
		if FileStoreDebug {
			fmt.Println("File does not exist. Storing file and starting republish procedures")
		}
	}
}

func (files *FileStore) getFile(fileID *KademliaID) (File, bool) {

	FileLock.Lock()
	file, fileExist := files.fileData[*fileID]
	FileLock.Unlock()

	return file, fileExist
}

func (files *FileStore) getFileContent(fileID *KademliaID) []byte {

	FileLock.Lock()
	file := files.fileData[*fileID]
	FileLock.Unlock()

	return file.content
}

func (files *FileStore) DeleteFile(fileID *KademliaID) {
	file, exists := files.getFile(fileID)
	file.nodeRepub <- false
	FileLock.Lock()
	if exists {
		delete(files.fileData, *fileID)
	} else {
		if FileStoreDebug {
			fmt.Println("File could not be removed, does not exist")
		}
	}
	FileLock.Unlock()
	if FileStoreDebug {
		fmt.Printf("Removed file with id: %s stopping republish\n", fileID.String())
	}
}

func (files *FileStore) republish(file File) {

	select {
	case repub := <-file.nodeRepub:
		if repub {
			if FileStoreDebug {
				fmt.Println("A FileRePub message received true, resetting repub timer. File: ", file.String())
			}
			files.republish(file)
		} else {
			if FileStoreDebug {
				fmt.Println("A FileRePub message received false, deleting file. File: ", file.String())
			}
			files.DeleteFile(file.fileID)
		}

	case <-time.After(NodeRepublish):
		fmt.Println("No republishes received, republishing file: ", file.String())
		go KademliaObj.Store(file.content, &file.owner)
		files.republish(file)
	}
}

func (files *FileStore) FileHeartbeat(file File) {
	if FileStoreDebug {
		fmt.Println("Staring heartbeat timer for file: ", file.String())
	}
	time.Sleep(OwnerRepublish)
	if FileStoreDebug {
		fmt.Println("Staring heartbeat procedure for file: ", file.String())
	}
	responseChannel := make(chan RPC)
	serial := NewRandomSerial()
	ConnectionLock.Lock()
	Connections[serial] = responseChannel
	ConnectionLock.Unlock()
	Net.SendFindDataMessage(file.fileID, &file.owner, serial)

	select {
	case responseRPC := <-responseChannel:
		if responseRPC.Value != nil {
			if FileStoreDebug {
				fmt.Println("Received response from owner. File exists: ", file.String())
			}
			files.FileHeartbeat(file)
		} else {
			if FileStoreDebug {
				fmt.Println("Received response from owner. File is deleted: ", file.String())
			}
			file.nodeRepub <- false
		}
	case <-time.After(TimeOut):
		if FileStoreDebug {
			fmt.Println("Did not receive any response from owner. Keeping file: ", file.String())
		}
		files.FileHeartbeat(file)
	}
}

type File struct {
	fileID    *KademliaID
	content   []byte
	owner     Contact
	nodeRepub chan bool
	pinned    bool
}

func NewFile(fileContent []byte, fileOwner *Contact) File {
	file := File{
		fileID:    NewRandomHash(string(fileContent)),
		content:   fileContent,
		owner:     *fileOwner,
		nodeRepub: make(chan bool),
	}
	return file
}

func (file *File) String() string {
	return fmt.Sprintf(`file("%s", "%s", "%s")`, file.fileID, file.content, file.owner.String())
}
