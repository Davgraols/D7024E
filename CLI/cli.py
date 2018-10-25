
import json
import socket
import sys
import requests
import base64

HOST, PORT = "localhost", 8080

sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)


def cli():
    while True:
        print "1 = savefile. 2 = getfile, 3 = pinfile, 4 = unpinfile/delete, 5 = exit"
        command = input()
        if command == 1:
            savefile()
        elif command == 2:
            getfile()
        elif command == 3:
            pinfile()
        elif command == 4:
            unpinfile()
        elif command == 5:
            break

def savefile():
    print "type in the name as a strin on the file you want to save: "
    filename = input()
    f = open(filename, 'r')
    content = f.read()
    #jsonObj = json.dumps(content)
    f.close()
    headers = {
    "data": content,
    }
    jsonObj = json.dumps(headers)
    print jsonObj
    resp = requests.post(url = "http://localhost:8080/file", data = jsonObj)
    print resp.text

def deletefile():
    print "type in the name as a strin on the fileID you want to delete: "
    deleteID = input()
    resp = requests.delete(url = "http://localhost:8080/file/"+deleteID)
    print resp.text

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
    jsonObj = json.dumps(headers)
    resp = requests.patch(url = "http://localhost:8080/file/"+deleteID, data = jsonObj)
    print resp.text


def unpinfile():
    headers = {
    "Pintype": "unpin",
    }
    print "type in the name as a strin on the fileID you want to unpin: "
    deleteID = input()
    jsonObj = json.dumps(headers)
    resp = requests.patch(url = "http://localhost:8080/file/"+deleteID, data = jsonObj)
    print resp.text


cli()


