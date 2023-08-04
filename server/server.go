package server

import (
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
)

var (
	mutex          = &sync.Mutex{}
	IPCountMap     = map[string]int{}
	IPChallengeMap = map[string]string{}
	IPSolutionMap  = map[string]string{}
)

const (
	difficulty     = 4
	retriesAllowed = 10
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
		solution = RandomSalt(difficulty)
		challenge = HashStr(IP + solution)

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
	return c.SendString(`<html>Verifying Your Browser ... Please Allow Up To 10 Seconds</html><script>let ip="` + IP + `",challenge="` + challenge + `",difficulty=` + fmt.Sprint(difficulty) + `,workerScript='importScripts("https://cdnjs.cloudflare.com/ajax/libs/crypto-js/4.0.0/crypto-js.min.js"),self.onmessage=function(t){resp={match:function t(a,o,e){if(e>4)return!0;for(let n in a){if("function"==typeof a[n])break;if("object"==typeof a[n])t(a[n],o[n],e+1);else if(a[n]!=o[n])return!1}return!0}(navigator,t.data.navigator,0),solution:""},t.data.arr.forEach(a=>{CryptoJS.MD5(t.data.ip+a)==t.data.challenge&&(resp.solution=a,self.postMessage(resp))})};',possibleStrings=[];function iterateStrings(e,t){let r="abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890";if(e.length===t){possibleStrings.push(e);return}for(let n=0;n<r.length;n++)iterateStrings(e+r[n],t)}function spawnWorker(e){console.log("Spawned Worker");let t=new Blob([workerScript],{type:"text/javascript"});var r=URL.createObjectURL(t),n=new Worker(r);n.onmessage=solved;let o={challenge:challenge,navigator:navigatorData,ip:ip,arr:e};n.postMessage(o)}function divideArray(e,t){let r=[],n=e.length,o=Math.floor(n/t);for(let a=0;a<t;a++){let i=o*a,s=a===t-1?n:i+o;r.push(e.slice(i,s))}return r}function solved(e){e.data.match?(console.log("\uD83E\uDD73 Heureka",e.data),document.cookie="POW-Solution="+e.data.solution+"; SameSite=Lax; path=/; Secure",window.location.reload()):console.log("\uD83D\uDD75️ Something's wrong")}function cloneObject(e,t){var r={};if(t>4)return r;for(var n in e)"object"!=typeof e[n]||null==e[n]||e[n]instanceof Function?"function"==typeof e[n]||e[n]instanceof HTMLElement||(r[n]=e[n]):r[n]=cloneObject(e[n],t+1);return r}navigatorData=cloneObject(navigator,0);let numWorkers=navigator.hardwareConcurrency;void 0==numWorkers&&(numWorkers=2),numWorkers>4&&(numWorkers=4),iterateStrings("",difficulty),console.log("\uD83E\uDD71 Indexed Strings");let arrs=divideArray(possibleStrings,numWorkers);console.log("\uD83D\uDCAA Bruteforcing"),arrs.forEach(e=>{spawnWorker(e)});</script>`)
}
