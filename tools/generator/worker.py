import os
import sys

name = ""


def readFile(name):
    with open(name, "r") as f:
        return f.read()


def writeFile(name, text):
    if os.path.exists(name):
        os.remove(name)
    with open(name, "w") as f:
        f.write(text)
    os.system("go fmt " + name)

def genMain():
    txt = '''
    package main

    import (
        "east/core"
        "east/core/sys"
    )

    func main() {
        server := core.NewServer(
        )
        defer server.Exit()
        sys.WaitExitSignal()
    }
    '''
    writeFile("%s/main.go" % name, txt)

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("cmd argv error")
        exit()
    for arg in sys.argv[1:]:
        key, value = arg.split("=")
        if key == "name":
            name = value
    os.system("mkdir -p " + name + "/modules")
    genMain()
    