import { Terminal } from "xterm";
import "xterm/css/xterm.css";
import { FitAddon } from "xterm-addon-fit";
//import { AttachAddon } from "xterm-addon-attach";

const SHELL_PROMPT = "> ";

function prompt(term) {
  term.write("\r\n" + SHELL_PROMPT);
}

function createTerminal(container){
    const term = new Terminal({
      convertEol: true,
      smoothScrollDuration: 0,
    });
  
    const fitAddon = new FitAddon();
    term.loadAddon(fitAddon);
  
    //const attachAddon = new AttachAddon(webSocket);
    //term.loadAddon(attachAddon);
  
    term.onKey((key) => {
      // Track the user input
      let input = "";
      console.log("k:", key);
      const e = key.domEvent;
      switch (e.key) {
        case "c":
          if(e.ctrlKey){
            prompt(term);
            input = "";
            return
          }
          break
          case "l": {
            if (e.ctrlKey) {
              term.clear();
              return;
            }
            break;
          }
        case "Backspace":
          input = handleBackspace(term, input)
          break;
        case "Enter":
          input = input.trim();
          if(input.length === 0){
            prompt(term);
            return;
          }
          term.writeln("");
          break;
        case "ArrowLeft":
          break;
      }
    });
  
    term.onData((key) => {
      console.log("onData key:", key);
      term.write(key);
    });
  
    term.onLineFeed(() => {
      console.log("new line");
    });
  
    term.open(document.getElementById(container));
  
    setTimeout(function(){
      fitAddon.fit();
    }, 10)
    return term
}

async function initTerminalSession(term, host){
  term.writeln("Connecting to " + host);
  await sleep(1300);
  term.write(SHELL_PROMPT);
}

const isPrintableKeyCode = (keyCode) => {
  return (
    keyCode === 32 ||
    (keyCode >= 48 && keyCode <= 90) ||
    (keyCode >= 96 && keyCode <= 111) ||
    (keyCode >= 186 && keyCode <= 222)
  );
};

function handleBackspace(term, input) {
  if (input.length === 0) return input;
  term.write("\b \b");
  return input.substring(0, input.length - 1);
}

export function sleep(delay) {
  return new Promise((resolve) => setTimeout(resolve, delay));
}

export {createTerminal, initTerminalSession}
