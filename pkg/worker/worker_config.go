package worker

const DefaultConfig = `
def run_linux(task, assetBlob):
	# Upload Interpreter
	interpreter = cdn.openFile("paragon_interpreter_linux")

    f, err = ssh.openFile("/tmp/intp")
    assert.noError(err)

	err = file.copy(interpreter, f)
	assert.noError(err)

	err = file.chmod(f, "0755")
	assert.noError(err)

	# Upload Assets
	assetDst, err = ssh.openFile("/tmp/12345")
	assert.noError(err)

	err = file.write(assetDst, assetBlob)
	assert.noError(err)

	err = file.chmod(assetDst, "0644")
	assert.noError(err)

	# Run Task
	output, err = ssh.exec("/tmp/intp -f /tmp12345 -t "+task)
	print(output)
	assert.noError(err)

def worker_run(task, assetBlob):
	host = ssh.getRemoteHost()
    if host == "10.0.0.1":
		return run_linux(task, assetBlob)
    else:
		assert.noError("Unsupported Operating System")
`
