let ip = "1.1.1.1"
let challenge = "1416a8c534344980685a0073b97cb7dc"
let difficulty = 6000
let publicSalt = "CHANGE_ME"
let startDate = undefined

let workerScript = `

	importScripts('https://cdnjs.cloudflare.com/ajax/libs/crypto-js/4.0.0/crypto-js.min.js');

    self.onmessage = function(e) { 
	
		function compareObj(obj1, obj2, iteration){
			if(iteration > 4){
				return true
		  	}
			for(let key in obj1){
				if(typeof obj1[key] == "function"){
					return true
				}
			  	if(typeof obj1[key] == "object"){
					compareObj(obj1[key], obj2[key], iteration + 1)
				} else {
					if(obj1[key] != obj2[key]){
						return false
			  		}
				}
			}
		  	return true
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
	if (res.data.match) {
		let endDate = new Date
		console.log("🥳 Heureka", res.data)
		console.log("Solved In:", (endDate.getTime() - startDate.getTime())/1000)
		document.cookie = "POW-Solution=" + res.data.access + "; SameSite=Lax; path=/; Secure";
		window.location.reload()
	} else {
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

console.log("🥱 Starting Workers")

for(let i = 0; i < difficulty; i = i + divided){
	spawnWorker(i, i+divided)
}

console.log("💪 Bruteforcing")
startDate = new Date