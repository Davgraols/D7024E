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
	file, fileExist := files.getFile(fileId)
	if fileExist {
		if FileStoreDebug {
			fmt.Println("File Exists. sending true to nodeRepub")
		}
		file.nodeRepub <- true
	} else {
		file = NewFile(fileContent, owner)
		FileLock.Lock()
		files.fileData[*fileId] = file
		FileLock.Unlock()
		go files.republish(file.fileID)
		if !owner.ID.Equals(RT.me.ID) {
			go files.FileHeartbeat(file.fileID)
		} else {
			files.PinFile(fileId)
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

func (files *FileStore) republish(fileID *KademliaID) {
	file, exist := files.getFile(fileID)
	for exist {

		select {
		case repub := <-file.nodeRepub:
			if repub {
				if FileStoreDebug {
					fmt.Println("A FileRePub message received true, resetting repub timer. File: ", string(file.content))
				}
			} else {
				if FileStoreDebug {
					fmt.Println("A FileRePub message received false, deleting file. File: ", string(file.content))
				}
				files.DeleteFile(file.fileID)
			}

		case <-time.After(NodeRepublish):
			fmt.Println("No republishes received, republishing file: ", string(file.content))
			go KademliaObj.Store(file.content, &file.owner)
		}
	}
}

func (files *FileStore) FileHeartbeat(fileID *KademliaID) {

	file, exist := files.getFile(fileID)

	for exist {
		if FileStoreDebug {
			fmt.Println("Staring heartbeat timer for fileID: ", fileID.String())
		}
		time.Sleep(OwnerRepublish)
		file, exist = files.getFile(fileID)

		if file.republished {
			files.SetRepublished(fileID, false)
			if FileStoreDebug {
				fmt.Println("File republished is true. Staring heartbeat procedure for file: ", string(file.content))
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
						fmt.Println("Received response from owner. File exists: ", string(file.content))
					}
				} else {
					if FileStoreDebug {
						fmt.Println("Received response from owner. File is deleted: ", string(file.content))
					}
					file.nodeRepub <- false
					exist = false
				}
			case <-time.After(TimeOut):
				if FileStoreDebug {
					fmt.Println("Did not receive any response from owner. Keeping file: ", string(file.content))
				}
			}
		} else {
			if !file.pinned {
				if FileStoreDebug {
					fmt.Println("File has not been republished, deleting file. Heartbeat should stop")
				}
				file.nodeRepub <- false
				exist = false
			}
		}
	}
}

func (files *FileStore) PinFile(fileID *KademliaID) {
	file, exist := files.getFile(fileID)
	if exist {
		file.Pin()
		FileLock.Lock()
		files.fileData[*fileID] = file
		FileLock.Unlock()
	} else {
		if FileStoreDebug {
			fmt.Println("could not pin file. File does not exist.")
		}
	}
}

func (files *FileStore) UnpinFile(fileID *KademliaID) {
	file, exist := files.getFile(fileID)
	if exist {
		file.Unpin()
		FileLock.Lock()
		files.fileData[*fileID] = file
		FileLock.Unlock()
	} else {
		if FileStoreDebug {
			fmt.Println("could not unpin file. File does not exist.")
		}
	}
}

func (files *FileStore) SetRepublished(fileID *KademliaID, republished bool) {
	file, exist := files.getFile(fileID)
	if exist {
		file.SetRepublished(republished)
		FileLock.Lock()
		files.fileData[*file.fileID] = file
		FileLock.Unlock()
	} else {
		fmt.Println("File does not exist. SetRepublished does nothing")
	}
}

func (files *FileStore) IsPinned(fileID *KademliaID) bool {
	file, exist := files.getFile(fileID)
	if exist {
		return file.pinned
	} else {
		if FileStoreDebug {
			fmt.Println("File does not exist")
		}
		return false
	}
}

type File struct {
	fileID      *KademliaID
	content     []byte
	owner       Contact
	nodeRepub   chan bool
	republished bool
	pinned      bool
}

func NewFile(fileContent []byte, fileOwner *Contact) File {
	file := File{
		fileID:      NewRandomHash(string(fileContent)),
		content:     fileContent,
		owner:       *fileOwner,
		nodeRepub:   make(chan bool),
		republished: false,
	}
	return file
}

func (file *File) Pin() {
	file.pinned = true
}

func (file *File) Unpin() {
	file.pinned = false
}

func (file *File) SetRepublished(republished bool) {
	file.republished = republished
}

func (file *File) String() string {
	return fmt.Sprintf(`file("%s", "%s", "%s")`, file.fileID, file.content, file.owner.String())
}
