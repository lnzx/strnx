import { Terminal } from "xterm";
import "xterm/css/xterm.css";
import { FitAddon } from "xterm-addon-fit";
import { AttachAddon } from "xterm-addon-attach";

const SHELL_PROMPT = "> ";

function prompt(term) {
  term.write("\r\n" + SHELL_PROMPT);
}

function createTerminal(container){
    const term = new Terminal({
      convertEol: true,
      smoothScrollDuration: 0,
    });

    const socket = new WebSocket('ws://127.0.0.1:8080/attach');
    const attachAddon = new AttachAddon(socket);
    term.loadAddon(attachAddon);

    const fitAddon = new FitAddon();
    term.loadAddon(fitAddon);

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

export function sleep(delay) {
  return new Promise((resolve) => setTimeout(resolve, delay));
}

export {createTerminal, initTerminalSession}
