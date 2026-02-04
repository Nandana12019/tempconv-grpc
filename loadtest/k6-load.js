/**
 * K6 load test for TempConv backend API
 * Simulates many frontends calling the conversion endpoints.
 *
 * Run against local backend:
 *   k6 run --vus 50 --duration 60s loadtest/k6-load.js
 *
 * Run against deployed URL (e.g. GKE Ingress):
 *   k6 run -e BASE_URL=https://YOUR_INGRESS_IP loadtest/k6-load.js --vus 100 --duration 120s
 */
import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:80';

export const options = {
  stages: [
    { duration: '30s', target: 20 },   // Ramp up to 20 VUs
    { duration: '1m', target: 50 },   // Stay at 50 VUs
    { duration: '30s', target: 100 }, // Spike to 100 VUs
    { duration: '1m', target: 100 },   // Hold 100 VUs
    { duration: '30s', target: 0 },   // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],  // 95% of requests under 500ms
    http_req_failed: ['rate<0.01'],    // Error rate under 1%
  },
};

export default function () {
  // Celsius to Fahrenheit
  const c2f = http.get(`${BASE_URL}/api/celsius-to-fahrenheit?c=100`);
  check(c2f, {
    'c2f status 200': (r) => r.status === 200,
    'c2f body has fahrenheit': (r) => {
      try {
        const j = JSON.parse(r.body);
        return j.fahrenheit === 212;
      } catch (_) {
        return false;
      }
    },
  });

  // Fahrenheit to Celsius
  const f2c = http.get(`${BASE_URL}/api/fahrenheit-to-celsius?f=212`);
  check(f2c, {
    'f2c status 200': (r) => r.status === 200,
    'f2c body has celsius': (r) => {
      try {
        const j = JSON.parse(r.body);
        return j.celsius === 100;
      } catch (_) {
        return false;
      }
    },
  });

  sleep(0.5 + Math.random());
}
