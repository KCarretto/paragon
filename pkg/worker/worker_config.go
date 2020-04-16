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

def run_linux(task, assetBlob):
	# Upload Interpreter
	interpreter = cdn.openFile(RG_NIX)
	intPath = "/tmp/"+str(env.rand())
	ssh_copy(interpreter, intPath, "0755")

	# Upload Assets
	assetPath = "/tmp/"+str(env.rand())
	ssh_write(assetBlob, assetPath, "0644")

	# Upload Task
	taskPath = "/tmp/"+str(env.rand())
	ssh_write(task, taskPath, "0644")

	# Run Task
	output, err = ssh.exec(intPath+"-f "+assetPath+" -t "+taskPath)
	print(output)
	assert.noError(err)

def worker_run(task, assetBlob):
    if env.isLinux():
		return run_linux(task, assetBlob)
    else:
		assert.noError("Unsupported Operating System")
`
