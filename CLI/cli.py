
#import msvcrt as m
#import codecs, csv
import json
import socket
import sys

HOST, PORT = "localhost", 8080

sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)



def cli():
    print "1 = save file. 2 = deletefile, 3 = getfile, 4 = pinfile, 5 = unpinfile"
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
    print("this is jsonObj: ")
    print(jsonObj)#send this content to server
    f.close()
    try:
        # Connect to server and send data
        sock.connect((HOST, PORT))
        sock.sendall(jsonObj)
        # Receive data from the server and shut down
        received = sock.recv(1024)
    finally:
        sock.close()

def deletefile():
    print "im in delete file"

def getfile():
    print "im in getfile"

def pinfile():
    print "im in pinfile"

def unpinfile():
    print "im in unpin file"

cli()


