// --- Initial UI Setup ---
const spinner = document.getElementById('loading-spinner');
const logContainer = document.getElementById('loading-logs');
if (spinner) spinner.style.display = 'block';
if (logContainer) logContainer.style.display = 'none';

let logEntries = [];
let lastUpdateTime = 0;
const UPDATE_THROTTLE = 10; // Minimum ms between UI updates (faster)

// Capture console.log and add to loading logs
const originalLog = console.log;
console.log = function(...args) {
    originalLog.apply(console, args);
    const message = args.join(' ');
    addLogEntry(message);
};

function addLogEntry(message) {
    // On first log, hide spinner and show log container
    if (logContainer && logContainer.style.display === 'none') {
        if (spinner) spinner.style.display = 'none';
        logContainer.style.display = 'block';
    }

    const logEntry = document.createElement('div');
    logEntry.className = 'log-entry';
    
    // Style different types of log messages
    if (message.includes('[PROGRESS]')) {
        logEntry.className += ' progress';
        
        // Parse progress from the message and update UI immediately
        // Format: [PROGRESS] Stage: Message (percentage%)
        const progressMatch = message.match(/\[PROGRESS\]\s+([^:]+):\s+([^(]+)\s+\((\d+)%\)/);
        if (progressMatch) {
            const [, stage, msg, percentage] = progressMatch;
            updateLoadingProgressThrottled(parseInt(percentage), stage, msg.trim());
        }
    } else if (message.includes('DEBUG:')) {
        logEntry.className += ' debug';
    }
    
    logEntry.textContent = new Date().toLocaleTimeString() + ' ' + message;
    logContainer.appendChild(logEntry);
    
    // Auto-scroll to bottom
    logContainer.scrollTop = logContainer.scrollHeight;
    
    // Keep only last 50 log entries
    while (logContainer.children.length > 50) {
        logContainer.removeChild(logContainer.firstChild);
    }
}

// Throttled version to prevent too frequent updates
function updateLoadingProgressThrottled(percentage, stage, message) {
    const now = performance.now();
    if (now - lastUpdateTime < UPDATE_THROTTLE && percentage < 100) {
        return; // Skip update if too soon (unless it's 100%)
    }
    lastUpdateTime = now;
    updateLoadingProgressInternal(percentage, stage, message);
}

// Global function to update loading progress (called from Go)
window.updateLoadingProgress = function(percentage, stage, message) {
    // Direct call from Go - no throttling needed
    updateLoadingProgressInternal(percentage, stage, message);
};

// Internal function that actually updates the UI
function updateLoadingProgressInternal(percentage, stage, message) {
    const loadingScreen = document.getElementById('loading-screen');
    const percentageEl = document.getElementById('loading-percentage');
    const barEl = document.getElementById('loading-bar');
    const stageEl = document.getElementById('loading-stage');
    const messageEl = document.getElementById('loading-message');
    
    const displayPercentage = Math.min(100, Math.max(0, Math.round(percentage)));

    if (percentageEl) percentageEl.textContent = displayPercentage + '%';
    if (barEl) {
        barEl.style.width = displayPercentage + '%';
        // Force a reflow to ensure the transition happens
        barEl.offsetHeight;
    }
    if (stageEl) stageEl.textContent = stage;
    if (messageEl) messageEl.textContent = message;
    
    // Ensure UI updates are immediately visible
    if (loadingScreen) {
        loadingScreen.offsetHeight; // Force reflow
    }
    
    if (percentage >= 100) {
        console.log('Loading complete! Hiding loading screen...');
        // Show spinner while finalizing
        const spinner = document.getElementById('loading-spinner');
        if (spinner) spinner.style.display = 'block';
        
        // Hide loading screen after a brief delay
        setTimeout(() => {
            if (loadingScreen) {
                loadingScreen.style.display = 'none';
                console.log('Loading screen hidden successfully');
            } else {
                console.error('Loading screen element not found!');
            }
        }, 1000);
    }
}

// --- WASM Memory Allocation Patch ---
// Improve WASM setup for a game: larger memory, better error handling, and performance hints
const WASM_INITIAL_MB = 512; // 512 MiB initial
const WASM_MAX_MB = 2048;    // 2 GiB max
const memory = new WebAssembly.Memory({ initial: WASM_INITIAL_MB * 16, maximum: WASM_MAX_MB * 16 }); // 1 page = 64KiB

// Ensure Go WASM runtime is available
if (typeof Go === 'undefined') {
    throw new Error('Go WASM runtime (Go) is not defined. Please include wasm_exec.js before this script.');
}
if (!window.go) window.go = new Go();
if (!go.importObject) go.importObject = {};
if (!go.importObject.env) go.importObject.env = {};
go.importObject.env['memory'] = memory;

document.addEventListener('DOMContentLoaded', () => {
    addLogEntry('Preparing WebAssembly environment...');
    // Preload WASM file for faster startup
    fetch('main.wasm').then(response => {
        if (!response.ok) throw new Error('Failed to fetch main.wasm');
        return response.arrayBuffer();
    }).then(bytes => {
        addLogEntry('Instantiating WebAssembly module...');
        return WebAssembly.instantiate(bytes, go.importObject);
    }).then(result => {
        addLogEntry('Starting WebAssembly module...');
        go.run(result.instance);
    }).catch(err => {
        addLogEntry('Error loading WebAssembly: ' + err.message);
        console.error('Failed to load WebAssembly:', err);
        alert('WebAssembly failed to load: ' + err.message);
    });
});

// Optionally, expose memory stats for debugging
window.getWasmMemoryUsage = function() {
    return {
        usedMB: memory.buffer.byteLength / (1024 * 1024),
        totalPages: memory.buffer.byteLength / 65536,
        maxPages: memory.maximum
    };
};
