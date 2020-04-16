package worker

const DefaultConfig = `
RG_NIX = "renegade"
RG_WIN = "renegade.exe"

def ssh_copy(f, dstPath, perms):
	dst, err = ssh.openFile(dstPath)
    assert.noError(err)

	err = file.copy(dst, f)
	assert.noError(err)

	err = file.chmod(dst, perms)
	assert.noError(err)

def ssh_write(content, dstPath, perms):
	f, err = ssh.openFile(dstPath)
    assert.noError(err)

	err = file.write(f, content)
	assert.noError(err)

	err = file.chmod(f, perms)
	assert.noError(err)

def run_linux(bundle):
	# Upload Interpreter
	interpreter = cdn.openFile(RG_NIX)
	intPath = "/tmp/"+str(env.rand())
	ssh_copy(interpreter, intPath, "0755")

	# Upload Bundle
	bundlePath = "/tmp/"+str(env.rand())
	ssh_write(bundle, bundlePath, "0644")

	# Run Task
	output, err = ssh.exec(intPath+" --bundle "+bundlePath)
	print(output)
	assert.noError(err)

def worker_run(bundle):
    if env.isLinux():
		return run_linux(bundle)
    else:
		assert.noError("Unsupported Operating System")
`
