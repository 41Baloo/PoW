package server

import (
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	mutex       = &sync.RWMutex{}
	bypassMutex = &sync.Mutex{}
	IPMap       = map[string]IP_INFORMATION{}

	CurrTime      int64
	TotalRequests = 0
)

var (
	TimeValid         int
	Difficulty        int
	RetriesAllowed    int
	DynamicSaltLength int
)

func Middleware(c *fiber.Ctx) error {
	// Make immutable copy when using headers
	IP := c.IP() //string([]byte(c.GetReqHeaders()["Cf-Connecting-Ip"]))

	// Lock map to avoid race conditions, no need to write here tho. Just reading
	mutex.RLock()
	IPInfo, found := IPMap[IP]
	mutex.RUnlock()

	// The client didnt get challenged yet
	if !found || IPInfo.Solution == "" {
		return UAM(c, IP, IP_INFORMATION{})
	}
	// The client provided an invalid solution

	if c.Cookies("POW-Solution") != IPInfo.Solution {
		return UAM(c, IP, IPInfo)
	}

	// The client completed the PoW. Do whatever you want
	return Passed(c)
}

func Passed(c *fiber.Ctx) error {
	bypassMutex.Lock()
	TotalRequests++
	cTTR := TotalRequests
	bypassMutex.Unlock()
	switch c.Path() {
	case "/":
		return c.SendString("🥳 You Passed The POW")
	case "/dstat":
		c.Append("Content-Type", "text/html")
		return c.SendString(`<script>let lRs=0;setInterval(()=>{fetch("/info").then(e=>{e.text().then(e=>{let t=parseInt(e);t==NaN&&location.reload();let l=t-lRs;l>-1&&(document.body.innerHTML="[ Bypassing Requests Per Second ]: "+l),lRs=t})})},1e3);</script>`)
	case "/info":
		return c.SendString(strconv.Itoa(cTTR))
	default:
		return c.SendString("😢 Couldn't Find That Path")
	}
}

func UAM(c *fiber.Ctx, IP string, IPInfo IP_INFORMATION) error {

	if IPInfo == (IP_INFORMATION{}) {
		publicSalt := RandomStr(DynamicSaltLength)
		salt := strconv.Itoa(RandomNum(Difficulty))
		IPInfo.PublicSalt = publicSalt
		IPInfo.Solution = HashStr(salt + publicSalt)
		IPInfo.Challenge = HashStr(publicSalt + salt)
		IPInfo.Served = CurrTime
	}

	if IPInfo.Attempts > RetriesAllowed {
		return c.SendString("😢 You Failed The Challenge Too Many Times")
	}

	IPInfo.Attempts++

	mutex.Lock()
	IPMap[IP] = IPInfo
	mutex.Unlock()

	c.Append("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	c.Append("Content-Type", "text/html")
	return c.SendString(`<!doctypehtml><meta charset=utf-8><meta content="IE=edge"http-equiv=X-UA-Compatible><meta content="width=device-width,initial-scale=1"name=viewport><title>Attention !</title><script src=https://cdn.jsdelivr.net/gh/41Baloo/balooPow@latest/balooPow.js></script><link href=https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.15.1/css/all.min.css rel=stylesheet><style>html{background:#131516;color:#d8d4cf}body{display:flex;justify-content:center;align-items:center;height:100vh;margin:0;font-family:Helvetica,Arial,sans-serif}.container{width:230px;place-items:center;display:grid;transition:all .5s;position:relative}.disclaimer{margin-top:30px}.details{min-width:90%}.indexing{visibility:hidden}.bruteforcing{visibility:hidden}.solved{visibility:hidden}table{padding-left:20px;padding-right:20px}.w-full{text-align:left;width:100%}.blink span{opacity:1;animation:blinker 1.2s linear infinite}.blink span:nth-child(2){animation-delay:.4s}.blink span:nth-child(3){animation-delay:.8s}@keyframes blinker{50%{opacity:0}}</style><div class=container><i class="fa-5x fa-shield-alt fas"></i><div class="blink disclaimer"><a id=disclaimer_text>Checking Your Browser </a><span>.</span><span>.</span><span>.</span></div><table class=details><tr class=indexing><td class=w-full>Starting Workers<td class=blink id=index_res><span>.</span><span>.</span><span>.</span><tr class=bruteforcing><td class=w-full>Bruteforcing<td class=blink id=brute_res><span>.</span><span>.</span><span>.</span><tr class=solved><td class=w-full id=solved_text>Reloading<td class=blink id=solved_res><span>.</span><span>.</span><span>.</span></table></div><script>const pow=new BalooPow("` + IPInfo.PublicSalt + `",` + strconv.Itoa(Difficulty) + `,"` + IPInfo.Challenge + `"),disclaimer=document.getElementsByClassName("disclaimer")[0],disclaimerText=document.getElementById("disclaimer_text"),indexing=document.getElementsByClassName("indexing")[0],indexRes=document.getElementById("index_res"),bruteforcing=document.getElementsByClassName("bruteforcing")[0],bruteRes=document.getElementById("brute_res"),solvedClass=document.getElementsByClassName("solved")[0],solvedText=document.getElementById("solved_text"),solvedRes=document.getElementById("solved_res");console.log("\uD83E\uDD71 Starting Workers"),indexing.style.visibility="visible";function checkElement(e,i){e.children[0].innerHTML=i,e.classList.remove("blink"),e.children[1].remove(),e.children[1].remove()}console.log("\uD83D\uDCAA Bruteforcing"),indexRes.children[0].innerHTML="V",indexRes.classList.remove("blink"),indexRes.children[1].remove(),indexRes.children[1].remove(),bruteforcing.style.visibility="visible",pow.Solve().then(e=>{null!=e&&!0==e.match?(checkElement(bruteRes,"V"),solvedText.style.visibility="visible",document.cookie="POW-Solution="+e.access+"; SameSite=Lax; path=/",window.location.reload()):(checkElement(bruteRes,"V"),solvedText.innerHTML="Navigator Mismatch",checkElement(solvedRes,"X"),solvedClass.style.visibility="visible",checkElement(disclaimer,"Blocked"))});</script>`)
}

func ClearCache() {
	for {

		CurrTime = time.Now().Unix()

		mutex.Lock()
		for IP, IPInfo := range IPMap {
			if (CurrTime - IPInfo.Served) > int64(TimeValid) {
				//fmt.Println(IP + " Is No Longer Valid")
				delete(IPMap, IP)
			}
		}
		mutex.Unlock()
		time.Sleep(1 * time.Second)
	}
}
