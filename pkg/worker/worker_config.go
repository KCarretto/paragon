package worker

const DefaultConfig = `
RG_NIX = "renegade"
RG_BSD = "renegade_freebsd"
RG_WIN = "renegade.exe"

def ssh_copy(f, dstPath, perms):
    dst, err = ssh.openFile(dstPath)
    assert.noError(err)

    err = file.copy(dst, f)
    assert.noError(err)

    err = file.chmod(dst, perms)
    assert.noError(err)
    file.close(dst)

def ssh_write(content, dstPath, perms):
    f, err = ssh.openFile(dstPath)
    assert.noError(err)

    err = file.write(f, content)
    assert.noError(err)

    err = file.chmod(f, perms)
    assert.noError(err)
    file.close(f)

def encrypt_bundle(bundle):
    # Encrypt Bundle
    key, err = crypto.generateKey()
    assert.noError(err)
    encryptedBundle, err = crypto.encrypt(key, bundle)
    assert.noError(err)

    return encryptedBundle, key


def run_ssh(bundle, key, interpreterAsset):
    # Upload Bundle
    bundlePath = "/tmp/"+str(env.rand())
    bundleDst, err = ssh.file(bundlePath)
    assert.noError(err)
    file.write(bundleDst, bundle)

    # Upload Interpreter
    intpSrc, err = cdn.download(interpreterAsset)
    assert.noError(err)

    binPath = "/tmp/"+str(env.rand())
    intpDst, err = ssh.file(binPath)
    assert.noError(err)

    err = file.copy(intpSrc, intpDst)
    assert.noError(err)

    file.chmod(intpDst, "0755")

    # Run Task
    output, err = ssh.exec(binPath+" --bundle "+bundlePath+" --key "+str(key)+" --cleanup")
    print(output)
    assert.noError(err)

def worker_run(bundle):
    encryptedBundle, key = encrypt_bundle(bundle)

    if env.isLinux():
        return run_ssh(encryptedBundle, key, RG_NIX)
    elif env.OS() == "BSD":
        return run_ssh(encryptedBundle, key, RG_BSD)
    else:
        assert.noError("Unsupported Operating System")
`
