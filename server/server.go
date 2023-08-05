package server

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	mutex          = &sync.Mutex{}
	IPCountMap     = map[string]int{}
	IPSaltMap      = map[string]string{}
	IPChallengeMap = map[string]string{}
	IPSolutionMap  = map[string]string{}
)

const (
	difficulty     = 4
	retriesAllowed = 10
	publicSalt     = "wJs78inhgbztznsjksbgzbn3"
)

func Middleware(c *fiber.Ctx) error {
	IP := c.IP()

	// Lock map to avoid race conditions
	mutex.Lock()
	IPSolution := IPSolutionMap[IP]
	mutex.Unlock()

	if IPSolution == "" {
		return UAM(c, IP, "")
	}
	if c.Cookies("POW-Solution") != IPSolution {
		return UAM(c, IP, IPSolution)
	}

	return Passed(c)
}

func Passed(c *fiber.Ctx) error {
	return c.SendString("🥳 You Passed The POW")
}

func UAM(c *fiber.Ctx, IP string, solution string) error {

	var challenge string

	if solution == "" {
		salt := RandomSalt(difficulty)
		solution = HashStr(salt + IP)
		challenge = HashStr(IP + publicSalt + salt)

		// Lock map to avoid race condition. Stores the solution in a map for future reference, so we dont have to hash every time
		mutex.Lock()
		IPSolutionMap[IP] = solution
		IPChallengeMap[IP] = challenge
		mutex.Unlock()
	}

	if challenge == "" {
		mutex.Lock()
		challenge = IPChallengeMap[IP]
		RequestCount := IPCountMap[IP]
		mutex.Unlock()

		if RequestCount > retriesAllowed {
			return c.SendString("😢 You Failed The Challenge Too Many Times")
		}

		mutex.Lock()
		IPCountMap[IP]++
		mutex.Unlock()
	}

	c.Append("Content-Type", "text/html")
	return c.SendString(`<!DOCTYPE html><html><head><meta charset="utf-8"><meta http-equiv="X-UA-Compatible" content="IE=edge"><meta name="viewport" content="width=device-width, initial-scale=1"><title>Attention !</title><link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.15.1/css/all.min.css"><style>html{background: #131516;color: #d8d4cf;}body{display: flex;justify-content: center;align-items: center;height: 100vh;margin: 0;font-family:Helvetica,Arial,sans-serif}.container{width:230px;place-items: center;display: grid;transition:all .5s;position:relative}.disclaimer{margin-top: 30px;}.details{min-width: 90%;}.indexing{visibility:hidden}.bruteforcing{visibility:hidden}.solved{visibility:hidden}table{padding-left:20px;padding-right:20px}.w-full{text-align:left;width:100%}.blink span{opacity: 1;animation: blinker 1.2s linear infinite;}.blink span:nth-child(2){animation-delay: 0.4s;}.blink span:nth-child(3){animation-delay: 0.8s;}@keyframes blinker{50%{opacity:0}}</style></head><body><div class="container"><i class="fas fa-shield-alt fa-5x"></i><div class="disclaimer blink"><a>Checking Your Browser </a><span>.</span><span>.</span><span>.</span></div><table class="details"><tr class="indexing"><td class="w-full">Indexing</td><td id="index_res" class="blink"><span>.</span><span>.</span><span>.</span></td></tr><tr class="bruteforcing"><td class="w-full">Bruteforcing</td><td id="brute_res" class="blink"><span>.</span><span>.</span><span>.</span></td></tr><tr class="solved"><td class="w-full">Reloading</td><td class="blink"><span>.</span><span>.</span><span>.</span></td></tr></table></div></body><script>const ip="` + IP + `",challenge="` + challenge + `",difficulty=` + fmt.Sprint(difficulty) + `,publicSalt="` + publicSalt + `",indexing=document.getElementsByClassName("indexing")[0],indexRes=document.getElementById("index_res"),bruteforcing=document.getElementsByClassName("bruteforcing")[0],bruteRes=document.getElementById("brute_res"),solvedRes=document.getElementsByClassName("solved")[0];let workerScript='importScripts("https://cdnjs.cloudflare.com/ajax/libs/crypto-js/4.0.0/crypto-js.min.js"),self.onmessage=function(t){resp={match:function t(a,o,e){if(e>4)return!0;for(let n in a){if("function"==typeof a[n])break;if("object"==typeof a[n])t(a[n],o[n],e+1);else if(a[n]!=o[n])return!1}return!0}(navigator,t.data.navigator,0),solution:"",access:""},t.data.arr.forEach(a=>{CryptoJS.MD5(t.data.ip+t.data.publicSalt+a)==t.data.challenge&&(resp.solution=a,resp.access=CryptoJS.MD5(a+t.data.ip).toString(),self.postMessage(resp))})};',possibleStrings=[];function iterateStrings(e,t){let r="abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890";if(e.length===t){possibleStrings.push(e);return}for(let n=0;n<r.length;n++)iterateStrings(e+r[n],t)}function spawnWorker(e){console.log("Spawned Worker");let t=new Blob([workerScript],{type:"text/javascript"});var r=URL.createObjectURL(t),n=new Worker(r);n.onmessage=solved;let i={challenge:challenge,navigator:navigatorData,ip:ip,publicSalt:publicSalt,arr:e};n.postMessage(i)}function divideArray(e,t){let r=[],n=e.length,i=Math.floor(n/t);for(let s=0;s<t;s++){let o=i*s,a=s===t-1?n:o+i;r.push(e.slice(o,a))}return r}function solved(e){e.data.match?(console.log("\uD83E\uDD73 Heureka",e.data),bruteRes.children[0].innerHTML="V",bruteRes.classList.remove("blink"),bruteRes.children[1].remove(),bruteRes.children[1].remove(),solvedRes.style.visibility="visible",document.cookie="POW-Solution="+e.data.access+"; SameSite=Lax; path=/; Secure",window.location.reload()):console.log("\uD83D\uDD75️ Something's wrong")}function cloneObject(e,t){var r={};if(t>4)return r;for(var n in e)"object"!=typeof e[n]||null==e[n]||e[n]instanceof Function?"function"==typeof e[n]||e[n]instanceof HTMLElement||(r[n]=e[n]):r[n]=cloneObject(e[n],t+1);return r}indexing.style.visibility="visible",navigatorData=cloneObject(navigator,0);let numWorkers=navigator.hardwareConcurrency;void 0==numWorkers&&(numWorkers=2),numWorkers>8&&(numWorkers=8),iterateStrings("",4),console.log("\uD83E\uDD71 Indexed Strings"),indexRes.children[0].innerHTML="V",indexRes.classList.remove("blink"),indexRes.children[1].remove(),indexRes.children[1].remove(),bruteforcing.style.visibility="visible";let arrs=divideArray(possibleStrings,numWorkers);console.log("\uD83D\uDCAA Bruteforcing"),arrs.forEach(e=>{spawnWorker(e)});</script></html>`)
}

func ClearCache() {
	for {
		mutex.Lock()
		IPCountMap = map[string]int{}
		IPSaltMap = map[string]string{}
		IPChallengeMap = map[string]string{}
		IPSolutionMap = map[string]string{}
		mutex.Unlock()
		time.Sleep(10 * time.Minute)
	}
}
