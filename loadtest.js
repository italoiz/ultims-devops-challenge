import http from "k6/http";
import { check, sleep } from "k6";
import { randomIntBetween } from "https://jslib.k6.io/k6-utils/1.2.0/index.js";

const BASE_URL = __ENV.BASE_URL || "http://localhost:8080";

const events = [
  { type: "page_view", weight: 5 },
  { type: "add_to_cart", weight: 3 },
  { type: "purchase", weight: 1 },
  { type: "begin_checkout", weight: 1 },
];

// Weighted random pick - mirrors real e-commerce funnel distribution:
// lots of page views, fewer cart adds, even fewer purchases.
function pickEvent() {
  const total = events.reduce((s, e) => s + e.weight, 0);
  let r = Math.random() * total;
  for (const e of events) {
    r -= e.weight;
    if (r <= 0) return e.type;
  }
  return events[0].type;
}

function buildPayload(eventType) {
  const userId = `u${randomIntBetween(1, 500)}`;

  switch (eventType) {
    case "page_view":
      return {
        type: eventType,
        user_id: userId,
        data: { page: `/products/${randomIntBetween(1, 200)}` },
      };
    case "add_to_cart":
      return {
        type: eventType,
        user_id: userId,
        data: {
          product_id: `p${randomIntBetween(1, 200)}`,
          quantity: randomIntBetween(1, 5),
        },
      };
    case "purchase":
      return {
        type: eventType,
        user_id: userId,
        data: { amount: (Math.random() * 200 + 5).toFixed(2) },
      };
    case "begin_checkout":
      return {
        type: eventType,
        user_id: userId,
        data: { items: randomIntBetween(1, 8) },
      };
  }
}

export const options = {
  stages: [
    // Ramp up over 30s so the graphs show a clear traffic curve
    { duration: "30s", target: 20 },
    // Sustain for 2 minutes - enough for Prometheus to scrape several times
    { duration: "2m", target: 20 },
    // Spike to stress the API and make the latency graph interesting
    { duration: "15s", target: 50 },
    { duration: "30s", target: 50 },
    // Cool down
    { duration: "15s", target: 5 },
    { duration: "30s", target: 5 },
  ],
};

export default function () {
  const eventType = pickEvent();
  const payload = buildPayload(eventType);

  const res = http.post(`${BASE_URL}/events`, JSON.stringify(payload), {
    headers: { "Content-Type": "application/json" },
  });

  check(res, {
    "status is 202": (r) => r.status === 202,
  });

  // Also hit health/readiness so those show up on the dashboard
  if (Math.random() < 0.1) {
    http.get(`${BASE_URL}/health`);
  }

  sleep(randomIntBetween(50, 200) / 1000);
}
