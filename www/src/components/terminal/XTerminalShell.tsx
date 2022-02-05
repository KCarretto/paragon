import React from 'react';
// import Terminal from 'react-console-emulator'
import { XTerm } from 'xterm-for-react';



const debugCallback = (command) => {
  console.log(command.command);
}


export const XTerminalShell = ({t, handleCallback=debugCallback, commandOutput}) => {
  const xtermRef = React.useRef(null)

  var lock = false;
  React.useEffect(() => {
    xtermRef.current.terminal.write(commandOutput);
    lock = false
  }, [commandOutput])

  React.useEffect(() => {
    var shellPrompt = "username@"+t+"~:$ ";

    // You can call any method in XTerm.js by using 'xterm xtermRef.current.terminal.[What you want to call]
    xtermRef.current.terminal.writeln("Hello, World!")
    xtermRef.current.terminal.write('\r\n' + shellPrompt);
    xtermRef.current.terminal.setOption('cursorBlink', true);

    var cmd = '';
    var cmdlen = 0;

    xtermRef.current.terminal.onKey(function (ev) {
      function prompt(){
        xtermRef.current.terminal.write('\r\n' + shellPrompt);
        xtermRef.current.terminal.setOption('cursorBlink', true);
      }
      var printable = (
        !ev.domEvent.altKey && !ev.altGraphKey && !ev.ctrlKey && !ev.metaKey
      );
      console.log(commandOutput)
      console.log(ev.domEvent)
      if (ev.domEvent.keyCode == 13) {
        if(cmd === 'clear')
        {
          xtermRef.current.terminal.clear();
          cmdlen = 0;
        }
        handleCallback(cmd);

        lock = true
        xtermRef.current.terminal.write('\r\n');

        cmd = "";
        cmdlen = 0;
        prompt();

      } else if (ev.domEvent.keyCode == 8) {
        // Do not delete the prompt
        if (cmdlen > 0 ){
          xtermRef.current.terminal.write('\b \b');
          cmd = cmd.slice(0, -1)
          cmdlen-=1;
        }

        // if (xtermRef.current.terminal.x > 2) {
        // }
      } else if (ev.domEvent.keyCode == 85 && ev.domEvent.ctrlKey == true) {
        for (let i = 0; i < cmdlen; i++){
          xtermRef.current.terminal.write('\b \b');
        }
        cmd = ""
        cmdlen = 0
      } else if (ev.domEvent.keyCode == 9) {
        console.log("trying to do autocomplete");
      } else if (ev.domEvent.keyCode == 86 && ev.domEvent.ctrlKey == true) {
        navigator.clipboard.readText()
        .then(text => {
          console.log('Pasted content: ', text);
          cmd += text
          cmdlen += text.length
          xtermRef.current.terminal.write(text)
        })
        .catch(err => {
          console.error('Failed to read clipboard contents: ', err);
        });
      } else if (printable) {
        cmdlen += 1;
        cmd += ev.key;
        xtermRef.current.terminal.write(ev.key);
      }
    });
  }, []);



  return (
      // Create a new terminal and set it's ref.
      <XTerm ref={xtermRef} />
  )
}

export default XTerminalShell;