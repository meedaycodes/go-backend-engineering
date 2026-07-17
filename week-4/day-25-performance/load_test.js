import http from 'k6/http';
import { sleep, check } from 'k6';

// 10 virtual users hammering the server for 30 seconds.
// Each VU runs the default function in a tight loop until duration expires.
export const options = {
    vus: 10,
    duration: '30s',
};

const BASE_URL = 'https://localhost:8081';

// k6 ignores TLS errors for self-signed certs when this is set.
export const tlsConfig = {
    insecureSkipTLSVerify: true,
};

export default function () {

    // Each VU+iteration combo produces a unique email so signup never
    // collides across virtual users or loop iterations.
    const email = `user_${__VU}_${__ITER}@test.com`;

    const signupRes = http.post(
        `${BASE_URL}/auth/signup`,
        JSON.stringify({ name: 'Habeeb', email: email, password: 'password123' }),
        {
            headers: { 'Content-Type': 'application/json' },
        }
    );

    check(signupRes, { 'signup 201': (r) => r.status === 201 });

    if (signupRes.status !== 201) {
        return;
    }

    const token = JSON.parse(signupRes.body).access_token;

    // Fetch all users — this hits the DB read path on every iteration.
    const listRes = http.get(
        `${BASE_URL}/users/`,
        {
            headers: { Authorization: `Bearer ${token}` },
        }
    );

    check(listRes, { 'list users 200': (r) => r.status === 200 });

    sleep(1);
}
