
Welcome to Renegade's documentation!
====================================
The following doc page is automatically generated. Please be kind to it.


stdlib/assert
--------------------------------------

.. currentmodule:: assert

.. function:: equal(expected: starlark.Value,actual: starlark.Value) -> ()

	Equal will check if two values are equal. This function will result in a fatal error if the assertion is incorrect.

----

.. function:: noError(err: starlark.Value) -> ()

	NoError will check if the passed value is a starlark.NoneType, if not it will error out the script.  This function may cause a fatal error if the assertion is incorrect.

----

stdlib/assets
--------------------------------------

.. currentmodule:: assets

.. function:: openFile(path: String) -> (f: File,err: Error)

	OpenFile that was packed into the compiled binary. The resulting file does not support many operations such as Chown, Write, etc. but you may read it's contents or copy it to another file i.e. one opened by ssh or sys.

----

.. function:: drop(f: File,dstPath: String,perms: String) -> (err: Error)

	Drop will take a given file, copy it to disk, optionaly set the permissions move it to a given destination, and clean up the temp file created. The default perms are '0755'.

----

.. function:: require(filePath: String) -> ()

	Require will be used in the init function for the worker to specify which files you wish to include in the asset bundle which will be accessible on the target. Will fatal if error occurs.

----

.. function:: file(path: String) -> (f: File)

	Prepare a descriptor for a file that was packaged into the binary. The descriptor may be used with the file library.

----

stdlib/cdn
--------------------------------------

.. currentmodule:: cdn

.. function:: openFile(name: String) -> (f: File)

	OpenFile stored on the CDN. Writing to the file will cause an upload to the CDN, overwriting any previously stored contents. Reading the file will download it from the CDN. Since operations are performed lazily, openFile will never error however reading from or writing to the file may.

----

.. function:: upload(f: File) -> (err: Error)

	Upload a file to the CDN, overwriting any previously stored contents.

----

.. function:: download(name: String) -> (f: File,err: Error)

	Download a file from the CDN.

----

stdlib/file
--------------------------------------

.. currentmodule:: file

.. function:: move(f: File,dstPath: String) -> (err: Error)

	Move a file to the desired location.

----

.. function:: close(f: File) -> ()

	Close a file if possible, otherwise this operation is a no-op.

----

.. function:: name(f: File) -> (name: String)

	The name or path used to open the file.

----

.. function:: content(f: File) -> (content: String)

	Read and return the file's contents.

----

.. function:: write(f: File,content: String) -> ()

	Write sets the file's content, overwriting any previous value. It creates the file if it does not yet exist.

----

.. function:: copy(src: File,dst: File) -> (err: Error)

	Copy the file's content into a destination file, overwriting any previous value. It creates the destination file if it does not yet exist.

----

.. function:: remove(f: File) -> (err: Error)

	Remove the file

----

.. function:: chown(f: File,username: String,group: String) -> (err: Error)

	Chown modifies the file's ownership metadata. Passing an empty string for either the username or group parameter will result in a no-op. For example, file.chown(f, '', 'new_group') will change the file's group ownership to 'new_group' but will not affect the file's user ownership.

----

.. function:: chmod(f: File,mode: String) -> ()

	Chmod modifies the file's permission metadata. The strong passed is expected to be an octal representation of what os.FileMode you wish to set file to have (i.e. '0755').

----

.. function:: drop(src: File,dst: File,perms: ?String) -> (err: Error)

	Drop will:
	1. Copy a given file to a tempfile on disk
	2. Optionally set the permissions The default perms are '0755'.
	3. Move it to a given destination
	4. Clean up the temp file created.

----

stdlib/http
--------------------------------------

.. currentmodule:: http

.. function:: newRequest(url: String) -> (request: Request)

	NewRequest creates a new Request object to be passed around.

----

.. function:: setMethod(r: Request,method: String) -> ()

	SetMethod sets the http method on the request object.

----

.. function:: setHeader(r: Request,header: String,value: String) -> ()

	SetHeader sets the http header to the value passed on the request object.

----

.. function:: setBody(r: Request,value: String) -> ()

	SetBody sets the http body to the value passed on the request object.

----

.. function:: exec(r: Request) -> (response: String,err: Error)

	Exec sends the passed request object.

----

stdlib/process
--------------------------------------

.. currentmodule:: process

.. function:: kill(proc: Process) -> (err: Error)

	Kill a process (using SIGKILL).

----

.. function:: name(proc: Process) -> (name: String,err: Error)

	Name gets the name of the passed process.

----

stdlib/regex
--------------------------------------

.. currentmodule:: regex

.. function:: replace(oldString: String,pattern: String,newString: String) -> (replacedString: String,err: Error)

	Replace uses the golang regex lib to replace all occurences of the pattern in the old string into the new strong.

----

stdlib/ssh
--------------------------------------

.. currentmodule:: ssh

.. function:: setUser(user: String) -> ()

	SetUser sets the RemoteUser attribute to be used in the outgoing SSH Connection. WARNING: MUST BE CALLED BEFORE OTHER SSH CALLS TO WORK.

----

.. function:: exec(cmd: String,disown: ?Bool) -> (output: String,err: Error)

	Exec a command on the remote system using an underlying ssh session.

----

.. function:: openFile(path: String) -> (f: File,err: Error)

	OpenFile on the remote system using SFTP over SSH. The file is created if it does not yet exist.

----

.. function:: getRemoteHost() -> (host: String)

	GetRemoteHost will return the remote host being used by the worker to connect to.

----

.. function:: file(path: String) -> (f: File,err: Error)

	Prepare a descriptor for a file on the remote system using SFTP via SSH. The descriptor may be used with the file library.

----

stdlib/sys
--------------------------------------

.. currentmodule:: sys

.. function:: openFile(path: String) -> (f: File,err: Error)

	OpenFile uses os.Open to Open a file.

----

.. function:: detectOS() -> (os: String)

	DetectOS uses the GOOS variable to determine the OS.

----

.. function:: exec(executable: String,disown: ?Bool) -> (output: String,err: Error)

	Exec uses the os/exec.command to execute the passed executable/params. Disown will optionally spawn a process but prevent it's output from being returned.

----

.. function:: connections(parent: ?Process) -> (connections: []Connection)

	Connections uses the gopsutil/net to get all connections created by a process (or all by default).

----

.. function:: processes() -> (procs: []Process)

	Processes uses the gopsutil/process to get all processes.

----

.. function:: files() -> (files: []File)

	Files uses the ioutil.ReadDir to get all files in a given path.

----

.. function:: file(path: String) -> (f: File)

	Prepare a descriptor for a file on the system. The descriptor may be used with the file library.

----

stdlib/env
--------------------------------------

.. currentmodule:: env

.. function:: IP() -> ()

	IP returns the primary IP address.

----

.. function:: OS() -> ()

	OS returns the operating system.

----

.. function:: isLinux() -> ()

	IsLinux returns true if the operating system is linux.

----

.. function:: isWindows() -> ()

	IsWindows returns true if the operating system is windows.

----

.. function:: uid() -> (uid: string)

	UID returns the current user id. If not found, an empty string is returned.

----

.. function:: user() -> (username: string)

	User returns the current username. If not found, an empty string is returned.

----

.. function:: time() -> (i: int)

	Time returns the current number of seconds since the unix epoch.

----

.. function:: rand() -> (i: int)

	Rand returns a random int. Not cryptographically secure.

----

stdlib/crypto
--------------------------------------

.. currentmodule:: crypto

.. function:: generateKey() -> (key: Key,err: Error)

	GenerateKey creates a new Key object to be passed around.

----

.. function:: encrypt(key: Key,data: String) -> (ciphertext: String,err: Error)

	Encrypt takes a Key and some data and returns the AESGCM encrypted IV+ciphertext.

----

.. function:: decrypt(key: Key,data: String) -> (plaintext: String,err: Error)

	Decrypt takes a Key and AESGCM encrypted IV+ciphertext data and returns the plaintext.

----

