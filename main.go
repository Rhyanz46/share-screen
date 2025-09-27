// Mac → iPhone Screen Share (no‑login, LAN only)
// ------------------------------------------------
// Minimal Go + WebRTC signaling over HTTP (no WebSockets, no auth/login).
// Use on the same LAN. One sender (Mac) mirrors screen to one viewer (iPhone Safari).
//
// How to run:
// 1) `go run main.go`
// 2) On your Mac: open http://localhost:8080/sender and click "Start Share".
//    The page will show a Viewer URL (with a one-time token).
// 3) On your iPhone: open the Viewer URL in Safari. Boom — mirrored.
//
// Notes:
// - Uses `getDisplayMedia` (you choose which screen/window to share).
// - Codec is whatever Safari negotiates (H.264/VP8). No audio, just video.
// - LAN only by default; NAT traversal via Google STUN for convenience.
// - Single viewer at a time per token. No persistence, no login, no tracking.
// - This is intentionally bare-bones; tweak constraints or add PIN if needed.
//
// Tested on: macOS (Chrome/Safari) as sender, iOS Safari as viewer.
// ------------------------------------------------

package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type sdp struct {
	Type string `json:"type"`
	Sdp  string `json:"sdp"`
}

type signalEntry struct {
	Offer     *sdp
	Answer    *sdp
	CreatedAt time.Time
}

type signalStore struct {
	mu   sync.Mutex
	data map[string]*signalEntry
}

func newStore() *signalStore { return &signalStore{data: make(map[string]*signalEntry)} }

func (s *signalStore) newToken() string {
	b := make([]byte, 9)
	if _, err := rand.Read(b); err != nil {
		log.Printf("❌ Error generating random token: %v", err)
		return ""
	}
	token := base64.RawURLEncoding.EncodeToString(b)
	log.Printf("🆕 New token generated: %s...", token[:8])
	return token
}

func (s *signalStore) putOffer(token string, offer *sdp) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[token] = &signalEntry{Offer: offer, CreatedAt: time.Now()}
	log.Printf("📤 Offer created for token: %s (type: %s)", token[:8]+"...", offer.Type)
}

func (s *signalStore) getOffer(token string) (*sdp, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.data[token]
	if !ok || e.Offer == nil {
		log.Printf("❌ Offer not found for token: %s", token[:8]+"...")
		return nil, false
	}
	log.Printf("📥 Offer retrieved for token: %s", token[:8]+"...")
	return e.Offer, true
}

func (s *signalStore) putAnswer(token string, answer *sdp) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.data[token]
	if !ok {
		log.Printf("❌ Token not found for answer: %s", token[:8]+"...")
		return false
	}
	if e.Answer != nil {
		log.Printf("⚠️  Answer already exists for token: %s", token[:8]+"...")
		return false
	}
	e.Answer = answer
	log.Printf("📤 Answer created for token: %s (type: %s)", token[:8]+"...", answer.Type)
	log.Printf("🎯 WebRTC handshake completed for token: %s", token[:8]+"...")
	return true
}

func (s *signalStore) getAnswer(token string) (*sdp, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.data[token]
	if !ok || e.Answer == nil {
		log.Printf("❌ Answer not ready for token: %s", token[:8]+"...")
		return nil, false
	}
	log.Printf("📥 Answer retrieved for token: %s", token[:8]+"...")
	return e.Answer, true
}

func (s *signalStore) gc(maxAge time.Duration) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	log.Printf("🗑️  Token garbage collector started (cleanup every 1 min, expiry: %v)", maxAge)

	for range ticker.C {
		cutoff := time.Now().Add(-maxAge)
		s.mu.Lock()
		deleted := 0
		var expiredTokens []string
		for k, v := range s.data {
			if v.CreatedAt.Before(cutoff) {
				expiredTokens = append(expiredTokens, k[:8]+"...")
				delete(s.data, k)
				deleted++
			}
		}
		activeTokens := len(s.data)
		s.mu.Unlock()

		if deleted > 0 {
			log.Printf("🗑️  GC: cleaned up %d expired tokens: %v (active: %d)", deleted, expiredTokens, activeTokens)
		}
	}
}

var store = newStore()

func getLANIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Printf("Error getting network interfaces: %v", err)
		return ""
	}
	for _, iface := range ifaces {
		if (iface.Flags & net.FlagUp) == 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			ipnet, ok := a.(*net.IPNet)
			if !ok || ipnet.IP == nil || ipnet.IP.IsLoopback() {
				continue
			}
			ipv4 := ipnet.IP.To4()
			if ipv4 == nil {
				continue
			}
			// pick typical private ranges
			if ipv4[0] == 10 || (ipv4[0] == 192 && ipv4[1] == 168) || (ipv4[0] == 172 && ipv4[1] >= 16 && ipv4[1] <= 31) {
				return ipv4.String()
			}
		}
	}
	return ""
}

// updateSTUNServer replaces the hardcoded STUN server in JavaScript code
func updateSTUNServer(jsCode, stunServer string) string {
	// Replace the hardcoded STUN server URL
	return fmt.Sprintf(jsCode, stunServer)
}

func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		return // .env file not found, continue with defaults
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if err := os.Setenv(key, value); err != nil {
				log.Printf("Error setting env var %s: %v", key, err)
			}
		}
	}
}

func main() {
	// Load .env file first
	loadEnv()

	port := flag.String("port", "8080", "Server port")
	stunServer := flag.String("stun", "stun:stun.l.google.com:19302", "STUN server URL")
	tokenExpiry := flag.Duration("token-expiry", 30*time.Minute, "Token expiry duration")
	enableHTTPS := flag.Bool("https", false, "Enable HTTPS")
	certFile := flag.String("cert", "certs/server.crt", "Path to TLS certificate file")
	keyFile := flag.String("key", "certs/server.key", "Path to TLS private key file")
	flag.Parse()

	// Allow environment variables to override flags
	if envPort := os.Getenv("PORT"); envPort != "" {
		*port = envPort
	}
	if envStun := os.Getenv("STUN_SERVER"); envStun != "" {
		*stunServer = envStun
	}
	if envExpiry := os.Getenv("TOKEN_EXPIRY"); envExpiry != "" {
		if duration, err := time.ParseDuration(envExpiry); err == nil {
			*tokenExpiry = duration
		}
	}
	if envHTTPS := os.Getenv("ENABLE_HTTPS"); envHTTPS != "" {
		*enableHTTPS = envHTTPS == "true"
	}
	if envCert := os.Getenv("TLS_CERT_FILE"); envCert != "" {
		*certFile = envCert
	}
	if envKey := os.Getenv("TLS_KEY_FILE"); envKey != "" {
		*keyFile = envKey
	}

	go store.gc(*tokenExpiry)

	// Update JavaScript with configurable STUN server
	updatedSenderJS := updateSTUNServer(senderJS, *stunServer)
	updatedViewerJS := updateSTUNServer(viewerJS, *stunServer)

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/sender", serveSender)
	http.HandleFunc("/viewer", serveViewer)
	http.HandleFunc("/api/new", apiNewToken)
	// Single handlers per path; they switch on r.Method internally
	http.HandleFunc("/api/offer", apiPostOffer)
	http.HandleFunc("/api/answer", apiPostAnswer)

	http.HandleFunc("/assets/sender.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Write([]byte(updatedSenderJS))
	})
	http.HandleFunc("/assets/viewer.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Write([]byte(updatedViewerJS))
	})
	http.HandleFunc("/assets/style.css", serveCSS)
	http.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{
			"host":  r.Host,
			"lanIP": getLANIP(),
		}); err != nil {
			log.Printf("Error encoding info response: %v", err)
			http.Error(w, "internal server error", 500)
		}
	})

	addr := ":" + *port
	protocol := "HTTP"
	if *enableHTTPS {
		protocol = "HTTPS"
	}
	log.Printf("%s Server listening on %s", protocol, addr)
	log.Printf("LAN IP: %s", getLANIP())
	log.Printf("STUN Server: %s", *stunServer)
	log.Printf("Token Expiry: %s", *tokenExpiry)

	var err error
	if *enableHTTPS {
		log.Printf("TLS Certificate: %s", *certFile)
		log.Printf("TLS Private Key: %s", *keyFile)
		err = http.ListenAndServeTLS(addr, *certFile, *keyFile, nil)
	} else {
		log.Printf("⚠️  Running in HTTP mode - consider enabling HTTPS for production")
		err = http.ListenAndServe(addr, nil)
	}

	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, indexHTML)
}

func serveSender(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, senderHTML)
}

func serveViewer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, viewerHTML)
}

func apiNewToken(w http.ResponseWriter, r *http.Request) {
	log.Printf("📞 API: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	token := store.newToken()
	if token == "" {
		http.Error(w, "failed to generate token", 500)
		return
	}
	log.Printf("🚀 Sender session started with token: %s...", token[:8])
	if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		log.Printf("Error encoding token response: %v", err)
		http.Error(w, "internal server error", 500)
	}
}

func apiPostOffer(w http.ResponseWriter, r *http.Request) {
	log.Printf("📞 API: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	switch r.Method {
	case http.MethodPost:
		var payload struct {
			Token string `json:"token"`
			SDP   sdp    `json:"sdp"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			log.Printf("❌ Invalid offer payload: %v", err)
			http.Error(w, err.Error(), 400)
			return
		}
		log.Printf("🔴 Sender posting offer for token: %s...", payload.Token[:8])
		store.putOffer(payload.Token, &payload.SDP)
		w.WriteHeader(204)
	case http.MethodGet:
		token := r.URL.Query().Get("token")
		log.Printf("🔵 Viewer requesting offer for token: %s...", token[:8])
		if off, ok := store.getOffer(token); ok {
			if err := json.NewEncoder(w).Encode(off); err != nil {
				log.Printf("Error encoding offer response: %v", err)
				http.Error(w, "internal server error", 500)
				return
			}
			return
		}
		http.Error(w, "offer not found", 404)
	default:
		http.Error(w, "method not allowed", 405)
	}
}

func apiPostAnswer(w http.ResponseWriter, r *http.Request) {
	log.Printf("📞 API: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	switch r.Method {
	case http.MethodPost:
		var payload struct {
			Token string `json:"token"`
			SDP   sdp    `json:"sdp"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			log.Printf("❌ Invalid answer payload: %v", err)
			http.Error(w, err.Error(), 400)
			return
		}
		log.Printf("🔵 Viewer posting answer for token: %s...", payload.Token[:8])
		if !store.putAnswer(payload.Token, &payload.SDP) {
			log.Printf("❌ Failed to store answer for token: %s...", payload.Token[:8])
			http.Error(w, "invalid token or answer exists", 400)
			return
		}
		w.WriteHeader(204)
	case http.MethodGet:
		token := r.URL.Query().Get("token")
		log.Printf("🔴 Sender requesting answer for token: %s...", token[:8])
		if ans, ok := store.getAnswer(token); ok {
			if err := json.NewEncoder(w).Encode(ans); err != nil {
				log.Printf("Error encoding answer response: %v", err)
				http.Error(w, "internal server error", 500)
				return
			}
			return
		}
		http.Error(w, "answer not found", 404)
	default:
		http.Error(w, "method not allowed", 405)
	}
}

// --- Static assets ---

const indexHTML = `<!doctype html>
<html><head>
<meta charset="utf-8"/>
<title>Mac → iPhone Screen Share</title>
<link rel="stylesheet" href="/assets/style.css"/>
</head><body>
<div class="wrap">
  <h1>Mac → iPhone Screen Share</h1>
  <p>Minimal, no-login, local network only.</p>
  <a class="btn" href="/sender">Start as Sender (Mac)</a>
</div>
</body></html>`

const senderHTML = `<!doctype html>
<html><head>
<meta charset="utf-8"/>
<meta name="viewport" content="width=device-width, initial-scale=1"/>
<title>Sender</title>
<link rel="stylesheet" href="/assets/style.css"/>
</head><body>
<div class="wrap">
  <h2>Sender (Mac)</h2>
  <button id="start" class="btn">Start Share</button>
  <div id="info" class="card" style="display:none"></div>
  <video id="preview" autoplay playsinline muted class="preview"></video>
</div>
<script src="/assets/sender.js"></script>
</body></html>`

const viewerHTML = `<!doctype html>
<html><head>
<meta charset="utf-8"/>
<meta name="viewport" content="width=device-width, initial-scale=1"/>
<title>Viewer</title>
<link rel="stylesheet" href="/assets/style.css"/>
</head><body>
<div class="wrap">
  <h2>Viewer (iPhone)</h2>
  <video id="view" autoplay playsinline class="viewer"></video>
</div>
<script src="/assets/viewer.js"></script>
</body></html>`

func serveCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	fmt.Fprint(w, css)
}

const css = `:root{font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Inter,Roboto,Arial,sans-serif}body{margin:0;background:#0b0b0c;color:#f2f3f5}.wrap{max-width:800px;margin:32px auto;padding:0 16px}.btn{background:#4b8bff;color:#fff;border:none;padding:10px 16px;border-radius:12px;font-weight:600;cursor:pointer}.btn:hover{opacity:.9}.card{background:#15161a;border:1px solid #26282e;padding:12px;border-radius:12px;margin-top:12px}.preview,.viewer{width:100%;max-height:70vh;background:#000;border-radius:12px}`


const senderJS = `const startBtn=document.getElementById('start');
const preview=document.getElementById('preview');
const info=document.getElementById('info');

async function postJSON(url, data){
  const res = await fetch(url,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(data)});
  if(!res.ok) throw new Error(await res.text());
  return res.json().catch(()=>({}));
}
async function getJSON(url){
  const res=await fetch(url);
  if(!res.ok) throw new Error(await res.text());
  return res.json();
}
function waitIce(pc){
  if (pc.iceGatheringState==='complete') return Promise.resolve();
  return new Promise(res=>{
    function check(){ if(pc.iceGatheringState==='complete'){ pc.removeEventListener('icegatheringstatechange',check); res(); } }
    pc.addEventListener('icegatheringstatechange',check);
  });
}

startBtn.onclick = async ()=>{
  try {
    startBtn.disabled=true;

    // Check if getDisplayMedia is supported
    console.log('navigator.mediaDevices:', navigator.mediaDevices);
    console.log('getDisplayMedia:', navigator.mediaDevices?.getDisplayMedia);

    if (!navigator.mediaDevices) {
      throw new Error('MediaDevices API not available. Make sure you are using HTTPS or localhost.');
    }

    if (!navigator.mediaDevices.getDisplayMedia) {
      throw new Error('getDisplayMedia not supported. Chrome 72+, Firefox 66+, or Safari 13+ required.');
    }

    // fetch server info to build a LAN URL (avoid localhost on iPhone)
    const infoRes = await getJSON('/api/info');
    const baseHost = infoRes.lanIP || (new URL(location.href)).hostname;
    const baseOrigin = location.protocol + '//' + baseHost + ':' + location.port;

    // 1) get token
    const {token} = await postJSON('/api/new',{});

    // 2) capture screen
    const stream = await navigator.mediaDevices.getDisplayMedia({
      video: { frameRate: { ideal: 30 }, width: { ideal: 1920 }, height: { ideal: 1080 } },
      audio: false
    });
    preview.srcObject = stream;

  // 3) WebRTC PC
  const pc = new RTCPeerConnection({iceServers:[{urls:'%s'}]});
  stream.getTracks().forEach(t=>pc.addTrack(t, stream));

  // Connection status monitoring
  pc.oniceconnectionstatechange = ()=>{
    const state = pc.iceConnectionState;
    console.log('ICE Connection State:', state);

    if (state === 'connected' || state === 'completed') {
      info.innerHTML += '<br/><span style="color: #4CAF50; font-weight: bold;">✅ Viewer Connected!</span>';
    } else if (state === 'disconnected' || state === 'failed') {
      info.innerHTML += '<br/><span style="color: #f44336; font-weight: bold;">❌ Viewer Disconnected</span>';
    } else if (state === 'connecting') {
      info.innerHTML += '<br/><span style="color: #ff9800;">🔄 Connecting to viewer...</span>';
    }
  };

  pc.onconnectionstatechange = ()=>{
    console.log('PC Connection State:', pc.connectionState);
  };

  const offer = await pc.createOffer({offerToReceiveVideo:false});
  await pc.setLocalDescription(offer);
  await waitIce(pc); // ensure non-trickle offer includes candidates

  await fetch('/api/offer',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({token, sdp: pc.localDescription})});

    // show viewer URL using LAN IP
    const viewerURL = baseOrigin + '/viewer?token=' + encodeURIComponent(token);
    info.style.display='block';
    info.innerHTML = '<b>Viewer URL:</b> <code>'+viewerURL+'</code><br/><small>Open on iPhone Safari (same Wi‑Fi)</small><br/><span style="color: #ff9800;">⏳ Waiting for viewer to connect...</span>';

  } catch (error) {
    startBtn.disabled = false;
    info.style.display='block';
    info.innerHTML = '<b style="color: red;">Error:</b> ' + error.message;
    console.error('Screen sharing error:', error);
  }
};
`


const viewerJS = `const v=document.getElementById('view');
const params=new URLSearchParams(location.search);
const token=params.get('token');
if(!token){
  document.body.innerHTML='<div class="wrap"><p>Missing token. Open link from Sender page.</p></div>';
} else {
  start().catch(e=>{
    document.body.innerHTML='<div class="wrap"><p>Error: '+e+'</p></div>';
  });
}

async function getJSON(url){ const r=await fetch(url); if(!r.ok) throw new Error(await r.text()); return r.json(); }
async function postJSON(url,data){ const r=await fetch(url,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(data)}); if(!r.ok) throw new Error(await r.text()); return r.json().catch(()=>({})); }
function waitIce(pc){
  if (pc.iceGatheringState==='complete') return Promise.resolve();
  return new Promise(res=>{
    function check(){ if(pc.iceGatheringState==='complete'){ pc.removeEventListener('icegatheringstatechange',check); res(); } }
    pc.addEventListener('icegatheringstatechange',check);
  });
}

async function start(){
  const statusDiv = document.createElement('div');
  statusDiv.className = 'card';
  statusDiv.style.marginTop = '12px';
  statusDiv.innerHTML = '<span style="color: #ff9800;">🔄 Connecting to sender...</span>';
  document.querySelector('.wrap').appendChild(statusDiv);

  const pc = new RTCPeerConnection({iceServers:[{urls:'%s'}]});

  // Connection monitoring
  pc.oniceconnectionstatechange = ()=>{
    const state = pc.iceConnectionState;
    console.log('Viewer ICE State:', state);

    if (state === 'connected' || state === 'completed') {
      statusDiv.innerHTML = '<span style="color: #4CAF50; font-weight: bold;">✅ Connected! Receiving screen share</span>';
    } else if (state === 'disconnected' || state === 'failed') {
      statusDiv.innerHTML = '<span style="color: #f44336; font-weight: bold;">❌ Connection lost</span>';
    } else if (state === 'connecting') {
      statusDiv.innerHTML = '<span style="color: #ff9800;">🔄 Connecting...</span>';
    }
  };

  pc.ontrack = (ev)=>{
    console.log('Received video track');
    v.srcObject = ev.streams[0];
    v.play().catch(()=>{
      // iOS may block autoplay; show a tap-to-start overlay
      const wrap=document.createElement('div');
      wrap.className='wrap';
      wrap.innerHTML='<button class="btn" id="tap">Tap to start</button>';
      document.body.appendChild(wrap);
      document.getElementById('tap').onclick=()=>{ v.play(); wrap.remove(); };
    });
  };

  // get offer
  const offer = await getJSON('/api/offer?token='+encodeURIComponent(token));
  await pc.setRemoteDescription(offer);

  const answer = await pc.createAnswer();
  await pc.setLocalDescription(answer);
  await waitIce(pc); // ensure non-trickle answer includes candidates

  await postJSON('/api/answer',{token, sdp: pc.localDescription});

  statusDiv.innerHTML = '<span style="color: #2196F3;">🔗 Handshake completed, waiting for video...</span>';
}
`
