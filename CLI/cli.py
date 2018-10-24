
import json
import socket
import sys
import requests
import base64

HOST, PORT = "localhost", 8080

sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)



def cli():
    print "1 = savefile. 2 = deletefile, 3 = getfile, 4 = pinfile, 5 = unpinfile"
    command = input()
    if command == 1:
        savefile()
    elif command == 2:
        deletefile()
    elif command == 3:
        getfile()
    elif command == 4:
        pinfile()
    elif command == 5:
        unpinfile()

def savefile():
    print "type in the name as a strin on the file you want to save: "
    filename = input()
    f = open(filename, 'r')
    content = f.read()
    jsonObj = json.dumps(content)
    jsonDict = {'SaveFile' : content}
    print type(content)
    print("this is jsonObj: ")
    print(jsonDict)#send this content to server
    f.close()
    resp = requests.post(url = "http://localhost/file/{id}", data = jsonObj)
    print resp

def deletefile():
    print "type in the name as a strin on the fileID you want to delete: "
    deleteID = input()
    jsonObj = json.dumps(deleteID)
    resp = requests.delete(url = "localhost/file/{id}", data = jsonObj)
    print resp

def getfile():
    print "type in the name as a strin on the fileID you want to get: "
    deleteID = input()
    jsonObj = json.dumps(deleteID)
    resp = requests.get(url = "localhost/file/{id}", data = jsonObj)
    print resp

def pinfile():
    headers = {
    "Pintype": "pin",
}
    print "type in the name as a strin on the fileID you want to pin: "
    deleteID = input()
    jsonObj = json.dumps(deleteID)
    resp = requests.patch(url = "localhost/file/{id}", data = jsonObj, headers = headers)
    print resp


def unpinfile():
    headers = {
    "Pintype": "unpin",
    }
    print "type in the name as a strin on the fileID you want to unpin: "
    deleteID = input()
    jsonObj = json.dumps(deleteID)
    resp = requests.patch(url = "localhost/file/{id}", data = jsonObj, headers = headers)
    print resp


cli()


