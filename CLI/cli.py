
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
    f.close()
    headers = {
    "data": content,
    }
    resp = requests.post(url = "http://localhost:8080/file/", data = headers)
    print resp

def deletefile():
    print "type in the name as a strin on the fileID you want to delete: "
    deleteID = input()
    resp = requests.delete(url = "http://localhost:8080/file/"+deleteID)
    print resp

def getfile():
    print "type in the name as a strin on the fileID you want to get: "
    deleteID = input()
    resp = requests.get(url = "http://localhost:8080/file/"+deleteID) 
    print resp.text

def pinfile():
    headers = {
    "Pintype": "pin",
}
    print "type in the name as a strin on the fileID you want to pin: "
    deleteID = input()
    jsonObj = json.dumps(deleteID)
    resp = requests.patch(url = "http://localhost:8080/file/"+deleteID, data = headers)
    print resp


def unpinfile():
    headers = {
    "Pintype": "unpin",
    }
    print "type in the name as a strin on the fileID you want to unpin: "
    deleteID = input()
    jsonObj = json.dumps(deleteID)
    resp = requests.patch(url = "http://localhost:8080/file/"+deleteID, data = headers)
    print resp


cli()


