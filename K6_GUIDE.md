# k6 Load Testing Guide

k6 is a developer-centric load testing tool. Scripts are written in JavaScript,
results are structured and machine-readable, and it integrates cleanly into CI pipelines.

---

## Installation

```bash
brew install k6          # macOS
choco install k6         # Windows
sudo apt install k6      # Ubuntu/Debian
```

---

## Script Structure

Every k6 script has three parts:

```js
import http from 'k6/http';
import { sleep, check } from 'k6';

// 1. Options — configure load shape
export const options = {
    vus: 10,          // virtual users (concurrent)
    duration: '30s',  // how long to run
};

// 2. Default function — what each VU does in a loop
export default function () {
    const res = http.get('http://localhost:8080/healthz');
    check(res, { 'status 200': (r) => r.status === 200 });
    sleep(1);
}
```

---

## Running Tests

```bash
# Basic run
k6 run script.js

# Skip TLS verification (self-signed certs)
k6 run --insecure-skip-tls-verify script.js

# Override VUs and duration from command line
k6 run --vus 50 --duration 60s script.js

# Output results to a file
k6 run script.js --out json=results.json

# Run with summary output only (no live progress)
k6 run --quiet script.js
```

---

## Load Shapes (Options)

### Constant load
```js
export const options = {
    vus: 10,
    duration: '30s',
};
```

### Ramping load — gradual ramp up then down
```js
export const options = {
    stages: [
        { duration: '30s', target: 10 },   // ramp up to 10 VUs over 30s
        { duration: '1m',  target: 10 },   // hold at 10 VUs for 1 minute
        { duration: '10s', target: 0  },   // ramp down to 0
    ],
};
```

### Constant arrival rate — fixed requests per second regardless of VU count
```js
export const options = {
    scenarios: {
        constant_rps: {
            executor: 'constant-arrival-rate',
            rate: 100,              // 100 requests per second
            timeUnit: '1s',
            duration: '30s',
            preAllocatedVUs: 20,
        },
    },
};
```

---

## Making HTTP Requests

```js
import http from 'k6/http';

// GET
const res = http.get('https://api.example.com/users');

// POST with JSON body
const res = http.post(
    'https://api.example.com/auth/signup',
    JSON.stringify({ name: 'Habeeb', email: 'h@test.com', password: 'pass' }),
    { headers: { 'Content-Type': 'application/json' } }
);

// With auth header
const res = http.get('https://api.example.com/users', {
    headers: { Authorization: `Bearer ${token}` },
});

// With timeout
const res = http.get('https://api.example.com/users', { timeout: '5s' });
```

---

## Checks (Assertions)

Checks do not fail the test — they record pass/fail counts in the summary.
Use thresholds (below) to actually fail the test.

```js
import { check } from 'k6';

check(res, {
    'status is 200':        (r) => r.status === 200,
    'status is 201':        (r) => r.status === 201,
    'body contains id':     (r) => JSON.parse(r.body).id !== undefined,
    'duration < 500ms':     (r) => r.timings.duration < 500,
    'content-type is json': (r) => r.headers['Content-Type'] === 'application/json',
});
```

---

## Thresholds (Pass/Fail Criteria)

Thresholds cause k6 to exit with a non-zero code if breached — useful in CI.

```js
export const options = {
    vus: 10,
    duration: '30s',
    thresholds: {
        // 95th percentile response time under 200ms
        http_req_duration: ['p(95)<200'],

        // error rate under 1%
        http_req_failed: ['rate<0.01'],

        // all checks must pass
        checks: ['rate==1.0'],

        // specific named check
        'check_name{check:signup 201}': ['rate>0.95'],
    },
};
```

---

## Built-in Variables

```js
__VU      // current virtual user number (1-based)
__ITER    // current iteration number for this VU (0-based)
__ENV     // environment variables passed via -e flag

// Example: unique email per VU per iteration
const email = `user_${__VU}_${__ITER}@test.com`;
```

---

## Reading the Output

```
checks_total.......: 571     18.81/s
checks_succeeded...: 96.00%  548 out of 571
checks_failed......: 4.00%   23 out of 571

HTTP
http_req_duration....: avg=40ms  min=213µs  med=16ms  max=118ms  p(90)=86ms  p(95)=92ms
http_req_failed......: 1.92%  11 out of 571
http_reqs............: 571    18.81/s

EXECUTION
iteration_duration...: avg=1.04s  min=789µs  med=1.08s  max=1.11s
iterations...........: 291    9.58/s
vus..................: 10     min=10  max=10
```

### What each metric means

| Metric | What it tells you |
|---|---|
| `http_req_duration` | End-to-end request time including DNS, TCP, TLS, send, wait, receive |
| `http_req_failed` | % of requests that returned a non-2xx/3xx status |
| `http_reqs` | Total requests made and rate per second |
| `iteration_duration` | Time for one full run of your default function (includes sleep) |
| `checks_succeeded` | % of `check()` assertions that passed |
| `vus` | Active virtual users |
| `p(90)` | 90th percentile — 90% of requests were faster than this |
| `p(95)` | 95th percentile — the SLO metric most teams use |

### Key percentiles explained

- **`avg`** — the mean. Misleading if there are outliers. Don't use alone.
- **`med`** — median (p50). Half of requests were faster. Better than avg.
- **`p(90)`** — 9 out of 10 requests were faster. Common SLO target.
- **`p(95)`** — 19 out of 20 requests were faster. Strictest common SLO target.
- **`max`** — the worst single request. Often an outlier, but worth watching.

---

## Interpreting Results — Common Patterns

### Healthy service
```
http_req_duration p(95) < 200ms
http_req_failed   rate  < 1%
checks_succeeded  rate  = 100%
```

### Rate limited (429s flooding in)
```
http_req_failed rate = 99%
http_req_duration avg = ~100µs   ← too fast, rejected before handler
```
Fix: increase rate limit for load test IPs, or distribute test traffic.

### CPU bottleneck (bcrypt, encryption)
```
http_req_duration avg = 40-80ms
http_reqs rate = low despite low VU count
```
Fix: profile with pprof to confirm, then architect around the bottleneck.

### Database bottleneck (connection pool exhausted)
```
http_req_duration p(95) spikes as VUs increase
http_req_failed rate increases at higher load
```
Fix: increase pool size, add read replicas, add caching.

### Memory leak
```
http_req_duration increases over time during a sustained run
```
Fix: capture memory profile with `pprof/heap` and look for growing allocations.

---

## Capturing pprof Profile During Load Test

Run both simultaneously — the load test generates pressure while pprof samples:

```bash
# Terminal 1 — capture 25s CPU profile
curl -s "http://localhost:6060/debug/pprof/profile?seconds=25" -o cpu.prof &

# Terminal 2 — run load test
k6 run --insecure-skip-tls-verify load_test.js
```

Then analyse:

```bash
# Top functions by CPU time
go tool pprof -top cpu.prof

# Interactive browser-based flame graph
go tool pprof -http=:8090 cpu.prof

# Memory profile
curl -s "http://localhost:6060/debug/pprof/heap" -o mem.prof
go tool pprof -top mem.prof
```

---

## Available pprof Endpoints

```
http://localhost:6060/debug/pprof/               ← index of all profiles
http://localhost:6060/debug/pprof/profile?seconds=30  ← CPU profile
http://localhost:6060/debug/pprof/heap           ← memory (heap) snapshot
http://localhost:6060/debug/pprof/goroutine      ← all goroutine stack traces
http://localhost:6060/debug/pprof/block          ← goroutine blocking events
http://localhost:6060/debug/pprof/mutex          ← mutex contention
http://localhost:6060/debug/pprof/trace?seconds=5 ← execution trace
```

---

## Reading a pprof CPU Profile

```
flat  flat%   sum%        cum   cum%  function
10.23s 84.27% 84.27%     10.52s 86.66%  blowfish.encryptBlock
 1.09s  8.98% 93.25%     11.62s 95.72%  blowfish.ExpandKey
```

| Column | Meaning |
|---|---|
| `flat` | Time spent directly in this function (not in callees) |
| `flat%` | Percentage of total CPU time spent directly here |
| `sum%` | Cumulative % — how much is accounted for by top N functions |
| `cum` | Total time in this function including all functions it called |
| `cum%` | Cumulative % including callees |

**`flat` = where time is actually spent.**
**`cum` = which call path led there.**

If `flat` is high, that function is the bottleneck.
If `cum` is high but `flat` is low, it's a coordinator — look at its callees.

---

## CI Integration

```yaml
- name: Load test
  run: |
    k6 run --insecure-skip-tls-verify load_test.js
  # k6 exits non-zero if thresholds are breached
  # Set thresholds in the script to gate deployments on performance
```
