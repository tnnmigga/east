import os
import sys

path = ""
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


def init():
    global path, name
    dirname = path + "/" + name.lower()
    if os.path.exists(dirname):
            print("useCase already exists")
    os.mkdir(dirname)


def genUseCase():
    global path, name
    text = '''
    package {0}

    import (
        "east/game/modules/play/domain"
        "east/game/modules/play/domain/api"
    )

    type useCase struct {
        *domain.Domain
    }

    func New(d *domain.Domain) api.I{1} {
        c := &useCase{
            Domain: d,
        }
        c.After(idef.ServerStateInit, c.afterInit)
        return c
    }

    func (c *useCase) afterInit() error {
        return nil
    }

    '''.format(name.lower(), name)
    dirname = path + name.lower()
    writeFile(dirname + "/usecase.go", text)

def genDomain():
    global path, name
    text = readFile(path + "/domain/domain.go")
    index = text.find("MaxCaseIndex", 400)
    text = text[:index] + name + "CaseIndex" + text[index:]
    writeFile(path + "domain/domain.go", text)


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("cmd argv error")
        exit()
    for arg in sys.argv[1:]:
        key, value = arg.split("=")
        if key == "name":
            name = value
        if key == "module":
            path = value.lower()
            if path[-1] != "/":
                path += "/"
    init()
    genUseCase()
    genDomain()
    print("useCase generated")