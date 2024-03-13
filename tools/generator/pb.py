import os
import sys

path = ''

insertTxt = '''
import "vendor/github.com/gogo/protobuf/gogoproto/gogo.proto";

option go_package                           = "pb";
option (gogoproto.goproto_enum_prefix_all)  = false;
option (gogoproto.goproto_unrecognized_all) = false;
option (gogoproto.goproto_unkeyed_all)      = false;
option (gogoproto.goproto_sizecache_all)    = false;
'''

def gogoFile(path):
    proto_files = []
    for root, dirs, files in os.walk(path):
        for file in files:
            if file.endswith('.proto'):
                proto_files.append(os.path.join(root, file))
    if os.path.exists(path + '/tmp'):
        os.system('rm -r {}/tmp'.format(path))
    os.mkdir(path + '/tmp')
    for file in proto_files:
        with open(file, 'r+') as f:
            txt = f.read()
            index = insertIndex(txt)
            if index == -1:
                continue
            if txt.find('gogoproto') == -1:
                txt = txt[:index] + insertTxt + txt[index:]
        with open(path + '/tmp/'+file.split('/')[-1], 'w') as f:
            f.write(txt)
    os.system('protoc --proto_path=./ --gofast_out=.  {}/tmp/*.proto'.format(path))
    os.system('mv {}/tmp/*.go {}/'.format(path, path))    
    os.system('rm -r {}/tmp'.format(path))

def insertIndex(txt):
    i1 = txt.find('message')
    i2 = txt.find('enum')
    index = 1E10
    if i1 != -1:
        index = i1
    if i2 != -1:
        index = min(index, i2)
    if index == 1E10:
        return -1
    return index

if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("cmd argv error")
        exit()
    for arg in sys.argv[1:]:
        key, value = arg.split("=")
        if key == "path":
            path = value
    gogoFile(path)
    print("generated successfully".format(path))

