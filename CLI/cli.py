
#import msvcrt as m
#import codecs, csv
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
    requests.post(url = "http://localhost/file/{id}", data = jsonObj)

def deletefile():
    print "type in the name as a strin on the fileID you want to delete: "
    deleteID = input()
    jsonObj = json.dumps(deleteID)
    r = requests.post(url = "localhost/file/{id}", data = jsonObj)
    #try:
        # Connect to server and send data
        #sock.connect((HOST, PORT))
        #r = requests.post(url = "localhost/file/{id}", data = jsonObj)
        #sock.sendall(jsonObj)
        # Receive data from the server and shut down
        #received = sock.recv(1024)
    #finally:
        #sock.close()

def getfile():
    print "type in the name as a strin on the fileID you want to get: "
    deleteID = input()
    jsonObj = json.dumps(deleteID)
    try:
        # Connect to server and send data
        sock.connect((HOST, PORT))
        sock.sendall(jsonObj)
        # Receive data from the server and shut down
        received = sock.recv(1024)
    finally:
        sock.close()

def pinfile():
    print "type in the name as a strin on the fileID you want to pin: "
    deleteID = input()
    jsonObj = json.dumps(deleteID)
    try:
        # Connect to server and send data
        sock.connect((HOST, PORT))
        sock.sendall(jsonObj)
        # Receive data from the server and shut down
        received = sock.recv(1024)
    finally:
        sock.close()

def unpinfile():
    print "type in the name as a strin on the fileID you want to unpin: "
    deleteID = input()
    jsonObj = json.dumps(deleteID)
    try:
        # Connect to server and send data
        sock.connect((HOST, PORT))
        sock.sendall(jsonObj)
        # Receive data from the server and shut down
        received = sock.recv(1024)
    finally:
        sock.close()

cli()


