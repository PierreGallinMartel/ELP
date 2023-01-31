#! /usr/bin/env node

const readline = require('readline');
var prompt = require('prompt');
var exec = require('child_process').exec;

console.log("Welcome to Scrim! Here is a list of everything you can do:")
console.log("launch _path_to_app_")
console.log("list (similar to ls)")
console.log("lp (simimlar to ps)")
console.log("ping _website_name_")
console.log("bing [-k,-p,-c] _pid_ (kills, pauses or continues process)")
console.log("Add ' !' at the end of your commands to run them in the background (don't forget the space before the !)")
console.log("Ctrl-p to exit the shell")

//If ctrl-p pressed then exit shell
readline.emitKeypressEvents(process.stdin);
process.stdin.on('keypress', (str, key) => 
  {
    if (key.ctrl && key.name === 'p') 
      {
        process.exit();
      } 
  });
  
getPrompt();

function getPrompt() 
  {
    prompt.start();
    prompt.get(['command'], (err, res) =>
      {
        if (err) 
          {
            console.error(`error : ${err}`);
            return;
          }
        executeCommand(res.command);
      }) //When command is received, go to execute function
  }

function executeCommand(text)
  {
    command = formatCommand(text) // transform our command into a shell command
    if(command != null)
      {
        console.log("I am running : " + command) // actual shell command that's running
        var childPc = exec(command);
        childPc.stdout.on('data', (data)=>{console.log(data)}); //print data as it comes
        childPc.on('exit', ()=>getPrompt());
      }
    else getPrompt() 
  }

function formatCommand(c)
  {
    if (c=="") return null; //if we only pressed on enter

    var cm = c.split(' ');
    var txt;
    if(cm.length>0)
      {
        background = ""; //create a background variable that will be added to all of our commands
        if (cm[cm.length-1]=="!") background=" &";//if there is "!" at the end of the command we set it to "&", otherwise it stays to "" and it has no effect
        switch (cm[0])
          {
            case "launch": 
              if (cm.length>1 & cm[1]!="") txt = "open " + cm[1] + background;
              else console.log("Incomplete command.")
              break;

            case "list":
              txt = "ls | nl" + background;
              break;
          
            case "lp":
              txt = "ps | nl" + background;
              break;
            
            case "ping":
              if(cm.length>1 & cm[1]!="") txt = "ping " + cm[1] + background;
              else console.log("Incomplete command.")
              break;
            
            case "bing":
              if(cm.length>2)
                {
                  if (cm[1]=="-k") txt = "kill " + cm[2];
                  if(cm[1]=="-p") txt = "kill -STOP " + cm[2];
                  if(cm[1]=="-c") txt = "kill -CONT " + cm[2];
                }
              else console.log("Incomplete command.")
              break;
            
            case "keep":
              console.log("Sorry, I still haven't that figured out :)") //Weird problem with the "job" thing not working on my computer

            default: 
                console.log("I don't know this command. Write correctly next time.");
          }
      }
    return txt;
  }

// process.on('SIGHUP', () => { console.log("SIGHUP")})
// process.on('SIGINT', () => { console.log("SIGINT")})
// //process.on('SIGKILL', () => { console.log("SIGKILL")})

// nohup = function(){
//   setTimeout(nohup)
// }

//nohup()

//This code did not seem to work on my computer. When running it in alone in a separate file, we couldn't
//exit the shell but when we clicked on the cross and closed the terminal window the process would stop. 
