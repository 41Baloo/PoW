const ip = "1.1.1.1"
const challenge = "bafc7cafd278461aca805eff634b3973a957551d2e45605b4d51aab6cd716cc7"
const difficulty = 5000000
const publicSalt = "CHANGE_ME"
const disclaimer = document.getElementsByClassName("disclaimer")[0]
const disclaimerText = document.getElementById("disclaimer_text")
const indexing = document.getElementsByClassName("indexing")[0]
const indexRes = document.getElementById("index_res")
const bruteforcing = document.getElementsByClassName("bruteforcing")[0]
const bruteRes = document.getElementById("brute_res")
const solvedClass = document.getElementsByClassName("solved")[0]
const solvedText = document.getElementById("solved_text")
const solvedRes = document.getElementById("solved_res")
let startDate = undefined

console.log("🥱 Starting Workers")
indexing.style.visibility = "visible"

let workerScript = `

	importScripts('https://cdnjs.cloudflare.com/ajax/libs/crypto-js/4.0.0/crypto-js.min.js');

	self.onmessage = function(e) { 
	
		function compareObj(obj1, obj2, iteration){
			if(iteration > 4){
				return ""
			}
			for(let key in obj1){
				if(typeof obj1[key] == "function"){
					return ""
				}
				if(typeof obj1[key] == "object"){
					compareObj(obj1[key], obj2[key], iteration + 1)
				} else {
					if(obj1[key] != obj2[key]){
						return obj1[key].toString()
					}
				}
			}
			return ""
		}

		resp = {
			match: compareObj(navigator, e.data.navigator, 0),
			solution: "",
			access: ""
		}

		let start = e.data.start
		let end = e.data.end
		
		for(let i = start; i < end+1; i++){
			if(CryptoJS.SHA256(e.data.ip+e.data.publicSalt+i) == e.data.challenge){
				resp.solution = i
				resp.access = CryptoJS.SHA256(i+e.data.ip).toString()
				self.postMessage(resp)
				self.close()
			}
		}

		console.log("Worker Couldn't Find Hash ("+start+" - "+end+")")
		self.close()
	}
`

function checkElement(element, check){
	element.children[0].innerHTML = check
	element.classList.remove('blink')
	element.children[1].remove()
	element.children[1].remove()
}

function spawnWorker(start, end) {
	console.log("Spawned Worker")
	let blob = new Blob([workerScript], {
		type: 'text/javascript'
	});

	// Convert the Blob to a URL using URL.createObjectURL()
	var url = URL.createObjectURL(blob);

	// Create a new Worker using the Blob URL
	var worker = new Worker(url);

	// Listen for messages from the worker
	worker.onmessage = solved
	let workerMsg = {
		challenge: challenge,
		navigator: navigatorData,
		ip: ip,
		publicSalt: publicSalt,
		start: start,
		end: end
	}

	// Give the worker his array of strings to bruteforce
	worker.postMessage(workerMsg);
}

// A worker bruteforced it
function solved(res) {
	if (res.data.match == "") {
		let endDate = new Date
		console.log("🥳 Heureka", res.data)
		console.log("Solved In:", (endDate.getTime() - startDate.getTime())/1000)
		checkElement(bruteRes, "V")
		solvedText.style.visibility = "visible"
		document.cookie = "POW-Solution=" + res.data.access + "; SameSite=Lax; path=/; Secure";
		window.location.reload()
	} else {
		checkElement(bruteRes, "V")
		solvedText.innerHTML = res.data.match+" Mismatch"
		checkElement(solvedRes, "X")
		solvedClass.style.visibility = "visible"
		checkElement(disclaimer, "Blocked")
		console.log("🕵️ Something's wrong")
	}
}

function cloneObject(obj, iteration) {
	var clone = {};
	if (iteration > 4) {
		return clone
	}
	for (var i in obj) {
		if (typeof obj[i] == "object" && obj[i] != null && !(obj[i] instanceof Function))
			clone[i] = cloneObject(obj[i], iteration + 1);
		else if (typeof obj[i] !== 'function' && !(obj[i] instanceof HTMLElement))
			clone[i] = obj[i];
	}
	return clone;
}

// Clone navigator
navigatorData = cloneObject(navigator, 0);

// Calculate how many workers we can create
let numWorkers = navigator.hardwareConcurrency
if (numWorkers == undefined) {
	numWorkers = 2
}
if (numWorkers > 8) {
	numWorkers = 8
}

let divided = Math.ceil(difficulty/8)

for(let i = 0; i < difficulty; i = i + divided){
	spawnWorker(i, i+divided)
}

console.log("💪 Bruteforcing")
startDate = new Date
indexRes.children[0].innerHTML = "V"
indexRes.classList.remove('blink')
indexRes.children[1].remove()
indexRes.children[1].remove()
bruteforcing.style.visibility = "visible"