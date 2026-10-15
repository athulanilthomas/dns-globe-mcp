import "@mcp-b/global";
import createGlobe from "cobe";
import { SAMPLE } from "./data.js";
import { registerTools } from "./webmcp-tools.js";

const STATUSES = ["resolved", "stale", "pending", "error", "nxdomain"];

const STATUS_RGB = {
  resolved: [0.09, 0.64, 0.29],
  stale: [0.92, 0.55, 0.05],
  pending: [0.45, 0.45, 0.5],
  error: [166 / 255, 25 / 255, 46 / 255],
  nxdomain: [246 / 255, 130 / 255, 31 / 255]
};

const STATUS_TEXT = {
  resolved: "Resolved",
  stale: "Stale",
  pending: "Pending",
  error: "Error",
  nxdomain: "Non Existant"
};

const $ = (id) => document.getElementById(id);
const globeEl = $("globe");
const canvas = $("cobe");

const el = (tag, className, text) => {
  const node = document.createElement(tag);
  if (className) node.className = className;
  if (text != null) node.textContent = text;
  return node;
};

const anchorTo = (node, anchor, visibleVar) => {
  node.style.positionAnchor = anchor;
  node.style.opacity = `var(${visibleVar}, 0)`;
  node.style.filter = `blur(calc((1 - var(${visibleVar}, 0)) * 8px))`;
};

const toRad = (deg) => (deg * Math.PI) / 180;

const dpr = Math.min(window.devicePixelRatio || 1, 2);
let size = globeEl.offsetWidth;
let phi = 0;
let theta = 0.25;
let velocityPhi = 0;
let velocityTheta = 0;
let pointer = null;
let focusTarget = null;
let pausedUntil = 0;
let currentResults = [];
let activeId = null;

const MARKER_SIZE = 0.03;
const MARKER_SIZE_ACTIVE = 0.07;
const PAUSE_MS = 6000;

const globe = createGlobe(canvas, {
  devicePixelRatio: dpr,
  width: size * dpr,
  height: size * dpr,
  phi,
  theta,
  dark: 0,
  diffuse: 1.2,
  mapSamples: 16000,
  mapBrightness: 6,
  mapBaseBrightness: 0,
  baseColor: [1, 1, 1],
  markerColor: [0.1, 0.1, 0.1],
  glowColor: [1, 1, 1],
  markerElevation: 0.02,
  markers: [],
});

new ResizeObserver(() => {
  size = globeEl.offsetWidth;
  globe.update({ width: size * dpr, height: size * dpr });
}).observe(globeEl);

canvas.addEventListener("pointerdown", (e) => {
  pointer = { x: e.clientX, y: e.clientY };
  velocityPhi = velocityTheta = 0;
  canvas.setPointerCapture(e.pointerId);
});
canvas.addEventListener("pointermove", (e) => {
  if (!pointer) return;
  velocityPhi = ((e.clientX - pointer.x) / size) * 3;
  velocityTheta = ((e.clientY - pointer.y) / size) * 3;
  pointer = { x: e.clientX, y: e.clientY };
  phi += velocityPhi;
  theta = Math.max(-1.2, Math.min(1.2, theta + velocityTheta));
});
const endDrag = () => (pointer = null);
canvas.addEventListener("pointerup", endDrag);
canvas.addEventListener("pointercancel", endDrag);

const wrapAngle = (a) => Math.atan2(Math.sin(a), Math.cos(a));

(function animate() {
  if (pointer) {
    focusTarget = null;
  } else if (focusTarget) {
    const dPhi = wrapAngle(focusTarget.phi - phi);
    const dTheta = focusTarget.theta - theta;
    phi += dPhi * 0.08;
    theta += dTheta * 0.08;
    if (Math.abs(dPhi) < 0.001 && Math.abs(dTheta) < 0.001) focusTarget = null;
  } else {
    const paused = performance.now() < pausedUntil;
    velocityPhi *= 0.95;
    velocityTheta *= 0.9;
    phi += (paused ? 0 : 0.003) + velocityPhi;
    theta += velocityTheta + (paused ? 0 : (0.25 - theta) * 0.02);
  }
  globe.update({ phi, theta });
  requestAnimationFrame(animate);
})();

requestAnimationFrame(() => canvas.classList.add("ready"));

let overlays = [];

const markersFor = (results) =>
  results.map((r) => ({
    id: r.id,
    location: [r.lat, r.lng],
    size: r.id === activeId ? MARKER_SIZE_ACTIVE : MARKER_SIZE,
    color: STATUS_RGB[r.status],
  }));

function renderGlobe(results) {
  overlays.forEach((n) => n.remove());
  overlays = [];

  for (const r of results) {
    const node = el("div", `site site--${r.status}`);
    const pyramid = el("div", "pyramid");
    for (let i = 0; i < 4; i++) pyramid.append(el("div", "pyramid-face"));
    const chip = el("span", "site-chip");
    chip.append(el("span", "dot"), el("span", "site-name", r.resolver.toLowerCase()), el("span", "site-status", r.status));
    node.append(pyramid, chip);
    anchorTo(node, `--cobe-${r.id}`, `--cobe-visible-${r.id}`);
    overlays.push(node);
  }

  globeEl.append(...overlays);

  globe.update({ markers: markersFor(results) });
}

function renderPanel(results) {
  const counts = Object.fromEntries(STATUSES.map((s) => [s, results.filter((r) => r.status === s).length]));
//   const total = results.length || 1; // TODO: We will take this later

  const progress = $("progress");
  progress.replaceChildren(
    ...STATUSES.filter((s) => counts[s]).map((s) => {
      const seg = el("span", `seg seg--${s}`);
      seg.style.flexGrow = counts[s];
      return seg;
    }),
  );

  $("summary").replaceChildren(
    ...STATUSES.map((s) => {
      const li = el("li", `badge badge--${s}`);
      li.append(el("span", "dot"), el("span", "badge-count", String(counts[s])), el("span", null, STATUS_TEXT[s]));
      return li;
    }),
  );
}

function focus(results) {
  const target = results.find((r) => r.status !== "resolved") ?? results[0];
  if (!target) return;
  phi = Math.PI - (toRad(target.lng) - Math.PI / 2);
}

function setResults(resolverResults) {
  const results = resolverResults.map((r, i) => ({ ...r, id: `resolver-${i}` }));
  currentResults = results;
  activeId = null;
  focusTarget = null;
  renderGlobe(results);
  renderPanel(results);
  focus(results);
}

function setInput({ domain, recordType }) {
  if (typeof domain !== "string" || !domain) return;
  const title = $("domain");
  title.replaceChildren(document.createTextNode(domain.slice(0, 253)));
  if (typeof recordType === "string" && recordType) title.append(el("span", "record-type", recordType.slice(0, 10)));
}

function render(sample) {
  setInput(sample);
  setResults(sample.results);
}

window.updateGlobe = render;

const injected = window.__DNS_DATA__;
render(Array.isArray(injected) ? { ...SAMPLE, results: injected } : Array.isArray(injected?.results) ? injected : SAMPLE);

const toolResult = (payload, isError = false) => ({
  content: [{ type: "text", text: JSON.stringify(payload) }],
  ...(isError && { isError: true }),
});

const resolverNames = () => currentResults.map((r) => r.resolver);

const faceLocation = (lat, lng) => [Math.PI - (toRad(lng) - Math.PI / 2), Math.max(-1, Math.min(1, toRad(lat)))];

function focusResolver(args) {
  const resolver = typeof args === "string" ? args : args?.resolver;
  const wanted = typeof resolver === "string" ? resolver.trim().toLowerCase() : "";
  const target = currentResults.find((r) => r.resolver.toLowerCase() === wanted);
  if (!target) {
    return toolResult({ error: `Unknown resolver ${JSON.stringify(resolver)}`, available: resolverNames() }, true);
  }

  const [targetPhi, targetTheta] = faceLocation(target.lat, target.lng);
  focusTarget = { phi: targetPhi, theta: targetTheta };
  pausedUntil = performance.now() + PAUSE_MS;
  velocityPhi = velocityTheta = 0;
  activeId = target.id;
  globe.update({ markers: markersFor(currentResults) });

  return toolResult({ resolver: target.resolver, status: target.status, records: target.records });
}

function getGlobeState() {
  const active = currentResults.find((r) => r.id === activeId);
  return toolResult({
    highlighted: active?.resolver ?? null,
    resolvers: currentResults.map(({ resolver, status, records, error }) => ({ resolver, status, records, error })),
  });
}

window.__focusResolver = focusResolver;
window.__getGlobeState = getGlobeState;

registerTools([
  {
    name: "focus_resolver",
    description: "Rotate the DNS globe to a resolver (e.g. Cloudflare, Google, Quad9, OpenDNS), highlight its marker, and return its status and records.",
    inputSchema: {
      type: "object",
      properties: { resolver: { type: "string", description: "Resolver name, case-insensitive" } },
      required: ["resolver"],
    },
    execute: focusResolver,
  },
  {
    name: "get_globe_state",
    description: "Get the highlighted resolver and the status, records and error of every resolver shown on the globe.",
    inputSchema: { type: "object", properties: {} },
    annotations: { readOnlyHint: true },
    execute: getGlobeState,
  },
]);
