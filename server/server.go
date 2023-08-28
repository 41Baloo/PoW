package server

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	mutex          = &sync.Mutex{}
	IPMap          = map[string]IP_INFORMATION{}
	IPCountMap     = map[string]int{}
	IPSaltMap      = map[string]string{}
	IPChallengeMap = map[string]string{}
	IPSolutionMap  = map[string]string{}

	CurrTime int64
)

const (
	timeValid      = 240
	difficulty     = 5000000
	retriesAllowed = 10
	publicSalt     = "wJs78inhgbztznsjksbgzbn3"
)

func Middleware(c *fiber.Ctx) error {
	// Make immutable copy
	IP := string([]byte(c.GetReqHeaders()["Cf-Connecting-Ip"]))

	// Lock map to avoid race conditions
	mutex.Lock()
	IPInfo := IPMap[IP]
	mutex.Unlock()

	// The client didnt get challenged yet
	if IPInfo.Solution == "" {
		//fmt.Println(IP + " Wasn't Challenged Yet")
		return UAM(c, IP, IP_INFORMATION{})
	}
	// The client provided an invalid solution
	if c.Cookies(IP+"_POW-Solution") != IPInfo.Solution {
		//fmt.Println(IP + " Provided An Invalid Solution")
		return UAM(c, IP, IPInfo)
	}

	// The client completed the PoW. Do whatever you want
	return Passed(c)
}

func Passed(c *fiber.Ctx) error {
	return c.SendString("🥳 You Passed The POW")
}

func UAM(c *fiber.Ctx, IP string, IPInfo IP_INFORMATION) error {

	if IPInfo == (IP_INFORMATION{}) {
		salt := fmt.Sprint(RandomNum(difficulty))
		IPInfo.Solution = HashStr(salt + IP)
		IPInfo.Challenge = HashStr(IP + publicSalt + salt)
		IPInfo.Served = CurrTime
	}

	if IPInfo.Attempts > retriesAllowed {
		return c.SendString("😢 You Failed The Challenge Too Many Times")
	}

	IPInfo.Attempts++

	mutex.Lock()
	IPMap[IP] = IPInfo
	mutex.Unlock()

	c.Append("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	c.Append("Content-Type", "text/html")
	return c.SendString(`<!DOCTYPE html><html><head><meta charset="utf-8"><meta http-equiv="X-UA-Compatible" content="IE=edge"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Attention !</title><link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.15.1/css/all.min.css"><style>html{background:#131516;color:#d8d4cf}body{display:flex;justify-content:center;align-items:center;height:100vh;margin:0;font-family:Helvetica,Arial,sans-serif}.container{width:230px;place-items:center;display:grid;transition:all .5s;position:relative}.disclaimer{margin-top:30px}.details{min-width:90%}.indexing{visibility:hidden}.bruteforcing{visibility:hidden}.solved{visibility:hidden}table{padding-left:20px;padding-right:20px}.w-full{text-align:left;width:100%}.blink span{opacity:1;animation:blinker 1.2s linear infinite}.blink span:nth-child(2){animation-delay:.4s}.blink span:nth-child(3){animation-delay:.8s}@keyframes blinker{50%{opacity:0}}</style></head><body><div class="container"><i class="fas fa-shield-alt fa-5x"></i><div class="disclaimer blink"><a id="disclaimer_text">Checking Your Browser</a><span>.</span><span>.</span><span>.</span></div><table class="details"><tr class="indexing"><td class="w-full">Starting Workers</td><td id="index_res" class="blink"><span>.</span><span>.</span><span>.</span></td></tr><tr class="bruteforcing"><td class="w-full">Bruteforcing</td><td id="brute_res" class="blink"><span>.</span><span>.</span><span>.</span></td></tr><tr class="solved"><td id="solved_text" class="w-full">Reloading</td><td id="solved_res" class="blink"><span>.</span><span>.</span><span>.</span></td></tr></table></div></body><script>const ip="` + IP + `",challenge="` + IPInfo.Challenge + `",difficulty=` + fmt.Sprint(difficulty) + `,publicSalt="` + publicSalt + `",disclaimer=document.getElementsByClassName("disclaimer")[0],disclaimerText=document.getElementById("disclaimer_text"),indexing=document.getElementsByClassName("indexing")[0],indexRes=document.getElementById("index_res"),bruteforcing=document.getElementsByClassName("bruteforcing")[0],bruteRes=document.getElementById("brute_res"),solvedClass=document.getElementsByClassName("solved")[0],solvedText=document.getElementById("solved_text"),solvedRes=document.getElementById("solved_res");let startDate;console.log("\uD83E\uDD71 Starting Workers"),indexing.style.visibility="visible";let workerScript='importScripts("https://cdnjs.cloudflare.com/ajax/libs/crypto-js/4.0.0/crypto-js.min.js"),self.onmessage=function(t){resp={match:function t(a,o,e){if(e>4)return"";for(let n in a){if("function"==typeof a[n])break;if("object"==typeof a[n])t(a[n],o[n],e+1);else if(a[n]!=o[n])return a[n].toString()}return""}(navigator,t.data.navigator,0),solution:"",access:""};let a=t.data.start,o=t.data.end;for(let e=a;e<o+1;e++)CryptoJS.SHA256(t.data.ip+t.data.publicSalt+e)==t.data.challenge&&(resp.solution=e,resp.access=CryptoJS.SHA256(e+t.data.ip).toString(),self.postMessage(resp),self.close());console.log("Worker Couldn\'t Find Hash ("+a+" - "+o+")"),self.close()};';function checkElement(e,t){e.children[0].innerHTML=t,e.classList.remove("blink"),e.children[1].remove(),e.children[1].remove()}function spawnWorker(e,t){console.log("Spawned Worker");let n=new Blob([workerScript],{type:"text/javascript"});var s=URL.createObjectURL(n),l=new Worker(s);l.onmessage=solved;let a={challenge:challenge,navigator:navigatorData,ip:ip,publicSalt:publicSalt,start:e,end:t};l.postMessage(a)}function solved(e){if(""==e.data.match){let t=new Date;console.log("\uD83E\uDD73 Heureka",e.data),console.log("Solved In:",(t.getTime()-startDate.getTime())/1e3),checkElement(bruteRes,"V"),solvedText.style.visibility="visible",document.cookie=ip+"_POW-Solution="+e.data.access+"; SameSite=Lax; path=/; Secure",window.location.reload()}else checkElement(bruteRes,"V"),solvedText.innerHTML=e.data.match+" Mismatch",checkElement(solvedRes,"X"),solvedClass.style.visibility="visible",checkElement(disclaimer,"Blocked"),console.log("\uD83D\uDD75️ Something's wrong")}function cloneObject(e,t){var n={};if(t>4)return n;for(var s in e)"object"!=typeof e[s]||null==e[s]||e[s]instanceof Function?"function"==typeof e[s]||e[s]instanceof HTMLElement||(n[s]=e[s]):n[s]=cloneObject(e[s],t+1);return n}navigatorData=cloneObject(navigator,0);let numWorkers=navigator.hardwareConcurrency;void 0==numWorkers&&(numWorkers=2),numWorkers>8&&(numWorkers=8);let divided=Math.ceil(625e3);for(let i=0;i<difficulty;i+=divided)spawnWorker(i,i+divided);console.log("\uD83D\uDCAA Bruteforcing"),startDate=new Date,indexRes.children[0].innerHTML="V",indexRes.classList.remove("blink"),indexRes.children[1].remove(),indexRes.children[1].remove(),bruteforcing.style.visibility="visible";</script></html>`)
}

func ClearCache() {
	for {

		CurrTime = time.Now().Unix()

		mutex.Lock()
		for IP, IPInfo := range IPMap {
			if (CurrTime - IPInfo.Served) > timeValid {
				//fmt.Println(IP + " Is No Longer Valid")
				delete(IPMap, IP)
			}
		}
		mutex.Unlock()
		time.Sleep(1 * time.Second)
	}
}
